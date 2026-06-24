# Bot-PR Review — Dependency-Change Context + Comment Upsert — Design

## Problem

The AI bot-PR review (just shipped) assesses a PR from only its **own diff** (`go.mod`/`go.sum` version+hash lines) — it has no view of what actually changed in the dependency between X and Y, so its verdict is shallow and it can't summarize the change. Separately, it posts a **new** comment on each (re)assessment instead of updating the existing one, which can accumulate comments over a PR's life.

## Goals

- Feed the assessor the **real dependency changes** (PR changelog body + the upstream source diff X..Y) so the verdict is informed and it can produce a **summary of the dependency changes**.
- Surface that summary in the PR comment and on the dashboard.
- **Upsert** the review comment (edit in place by marker) — exactly one review comment per PR, never spam.
- Keep the reliable **forced tool-call** verdict (enum can't drift); GitHub-compare is the bounded source-diff mechanism, degrading gracefully for unresolvable modules.

## Design

### ghclient additions
- `PullRequest.Body string` (json `body`) — populated by adding `body` to `ListOpenPRs`'s `--json`/`-q`. Holds dependabot/renovate's compiled Release-notes/Changelog/Commits.
- `CompareDiff(repo, base, head string) ([]byte, error)` → `gh api repos/<repo>/compare/<base>...<head> -H "Accept: application/vnd.github.diff"` — the real unified source diff between two tags.
- `UpsertPRComment(repo string, pr int, marker, body string) error` — lists the PR's issue comments (`gh api repos/<repo>/issues/<pr>/comments`, REST numeric ids), finds one whose body contains `marker`, **PATCHes** it (`gh api -X PATCH repos/<repo>/issues/comments/<id>`); else POSTs a new one. Encapsulates the no-spam upsert. (`PostPRComment` stays for the remediation nudges.)
- All three on the `GitHub` interface + `CLI` + `FakeGitHub`.

### review pure helpers (`internal/review`)
- `parseBumps(diff []byte) []DepBump{Module, From, To}` — from the PR's `go.mod` `-`/`+` lines (`- mod vX` / `+ mod vY`), paired by module. Versions normalized without the leading `v`.
- `moduleRepo(mod string) (repo string, ok bool)` — module path → `owner/name`: `github.com/owner/name[/v2][/sub]` → `owner/name` (strip major-version suffix + subpath); `golang.org/x/<n>` → `golang/<n>`; `k8s.io/<n>` → `kubernetes/<n>`; `sigs.k8s.io/<n>` → `kubernetes-sigs/<n>`. Unresolvable → `ok=false` (degrade).

### Assessor + context assembly
- Interface changes to: `Assess(pr ghclient.PullRequest, context string) (verdict, reasoning, changesSummary string, err error)`.
- `review.Run` assembles `context` per PR (it has `gh`): the PR **body** (changelog), then for each bump (cap ~5) where `moduleRepo` resolves, `CompareDiff(repo, "v"+From, "v"+To)` capped (~40KB each, total context cap ~60KB) labelled `Upstream <module> <From>..<To>:`; then the PR's own diff. CompareDiff failures degrade silently (skip that bump's source diff, note it).
- `OpenAIAssessor`: tool `assess_pr` schema gains `changesSummary` (string); prompt presents the changelog + upstream source diffs + PR diff and asks for `{verdict, reasoning, changesSummary}` — `changesSummary` = a 1-3 sentence summary of what the dependency change introduces (features/fixes/breaking/security). All existing degrade paths still return `("needs_human_verification", reason, "", nil)`.
- `FakeAssessor` returns canned `{Verdict, Reasoning, ChangesSummary, Err}`.

### Comment upsert + state + render
- `state.PRReview.ChangesSummary string` (json `changesSummary,omitempty`).
- `review.Run` builds the comment body = `🔎 kairos-security review: **<verdict>** — <reasoning>` + `\n\n**Dependency changes:** <changesSummary>` + `cc <notify>` + the `<!-- ksec:review -->` marker, and posts via `gh.UpsertPRComment(repo, pr, reviewMarker, body)` (edit-or-create) — replacing the old `PostPRComment`. Idempotency on head SHA still gates *whether* we re-assess; the upsert guarantees one comment regardless.
- Render "🔎 Bot-PR reviews": append the `changesSummary` under each row (when present).

## Out of scope

- nib-based assessment (decided against — forced tool-call is the engine).
- Forge-agnostic `go mod download` diffing (GitHub-compare is the bounded source-diff; non-GitHub modules degrade to body+PR-diff).
- Diffing more than ~5 bumped modules per PR / unbounded context (capped).

## Testing

- ghclient: `Body` projected; `CompareDiff` shells the compare endpoint with the diff media type; `UpsertPRComment` edits an existing marker comment vs creates (Fake records edit-vs-create). 
- `parseBumps`: a go.mod diff with one and multiple bumps → correct `{module,from,to}`; non-bump lines ignored.
- `moduleRepo`: github.com/owner/repo(/v2)(/sub), golang.org/x/net, k8s.io/api, sigs.k8s.io/x, and an unresolvable vanity path → ok=false.
- `review.Run` (fakes): context includes body + upstream diff (when resolvable) + PR diff; assessor receives it; `changesSummary` recorded; comment **upserted** (second run with a changed head edits, doesn't duplicate); CompareDiff error degrades (still assesses); dry-run no writes.
- `OpenAIAssessor`: forced tool call returns `{verdict,reasoning,changesSummary}`; degrade paths → needs_human + empty summary.
- render: reviews section shows the changes summary; deterministic; golden regenerated.
- Manual: dry-run shows the verdict + a real dependency-change summary; live upserts a single comment per PR.
