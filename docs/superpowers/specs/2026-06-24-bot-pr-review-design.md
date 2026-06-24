# AI Bot-PR Review — Design

## Problem

Bot-authored dependency PRs (dependabot/renovate; CVE and routine bumps) land constantly, and a human has to eyeball each to judge whether the change is *genuine* (a clean version bump) or hides something (unexpected source edits, supply-chain red flags). The maintainer wants the bot to **scan the actual changes a PR introduces and emit a gated verdict** — good / bad / needs_human_verification — using the diff as context, surfaced on the dashboard and as a PR comment, with configurable notification and auto-approval.

## Goal

A `ksec review` phase that, for each open bot-authored PR on a tracked repo, fetches the diff, asks LocalAI (forced tool call) for a supply-chain verdict, records it, comments on the PR (cc-ing configured handles), and — when configured — auto-approves a `good` PR. Idempotent on the PR head SHA so it neither re-spams nor re-spends the model.

## Design

### Config (`ai.yaml`)
```yaml
review:
  enabled: true                        # also requires LocalAI configured
  autoApprove: false                   # verdict==good -> gh pr review --approve
  maxPerRun: 20                        # cap NEW assessments/comments per run
  notify: ["@kairos-io/maintainers"]   # @handles/teams cc'd in the comment
```
`config.AIConfig` gains `Review ReviewCfg{Enabled bool; AutoApprove bool; MaxPerRun int; Notify []string}` (yaml-tagged). `LoadAI` defaults `MaxPerRun` to 20 when ≤0.

### Assessor (`internal/review`)
- `type Assessor interface { Assess(diff []byte, pr ghclient.PullRequest) (verdict, reasoning string, err error) }`.
- `OpenAIAssessor` (mirrors `internal/triage/openai.go`): a **forced tool call** `assess_pr` whose JSON-schema parameters constrain `verdict` to the enum `["good","bad","needs_human_verification"]` plus a `reasoning` string. Prompt: review an automated dependency-bump PR for supply-chain safety given its title + diff — `good` = genuine/expected bump (version/lockfile changes consistent with the title, no suspicious source edits); `bad` = clear red flags (unexpected code changes in a bump PR, malicious-looking additions, version/identity mismatch); `needs_human_verification` = anything uncertain. The diff is truncated (cap ~60KB) with a note when truncated. On any AI error/unparseable result → `needs_human_verification` (safe default), never a hard failure.
- `FakeAssessor` for tests (canned verdict/reasoning/err).

### Orchestration (`review.Run`)
`func Run(repos []state.Repo, gh ghclient.GitHub, a Assessor, cfg config.ReviewCfg, prev []state.PRReview, runID string, dryRun bool) ([]state.PRReview, []state.CollectionError)`:
- For each repo, `gh.ListOpenPRs` → keep `pr.IsBot` (all bot PRs — CVE or not).
- **Idempotency:** index `prev` by `repo|pr`. If `prev[key].HeadSHA == pr.HeadSHA`, carry the existing review forward unchanged (no AI, no comment, no approve). This bounds cost and prevents re-spam.
- Otherwise, while under `cfg.MaxPerRun` new assessments: `diff = gh.PRDiff(repo, pr)`; `verdict, reasoning = a.Assess(diff, pr)`; record `PRReview{Repo, PR, URL, HeadSHA, Verdict, Reasoning, ReviewedRun: runID}`.
  - **Comment** (unless dryRun): `gh.PostPRComment` with `🔎 kairos-security review: **<verdict>** — <reasoning>` + `cc <notify…>` + a `<!-- ksec:review -->` marker. dryRun prints the intended comment.
  - **Auto-approve** (unless dryRun): if `cfg.AutoApprove && verdict=="good"` → `gh.ApprovePR(repo, pr, body)`. `bad`/`needs_human_verification` never approve.
- Reviews over the cap (not reassessed this run) keep their carried-forward `prev` entry so the dashboard still shows them. Per-repo `ListOpenPRs`/`PRDiff` errors become `CollectionError`s (surfaced), not failures.
- Output sorted by repo then PR (deterministic).

### ghclient additions
- `PullRequest.HeadSHA` (json `headSHA`) populated from `headRefOid` in `ListOpenPRs`.
- `PRDiff(repo string, pr int) ([]byte, error)` → `gh pr diff <pr> -R <repo>`.
- `ApprovePR(repo string, pr int, body string) error` → `gh pr review <pr> -R <repo> --approve --body <body>`.
- All three added to the `GitHub` interface + `CLI` + the test `FakeGitHub`.

### State + render
- `state.PRReview{Repo string; PR int; URL, HeadSHA, Verdict, Reasoning, ReviewedRun string}`; `state/reviews.json` = `[]PRReview`; `ReviewsFile` const.
- `render.Input` gains `Reviews []state.PRReview`; a **"🔎 Bot-PR reviews"** section groups by repo: `- [#N title-or-pr](url) — **<verdict>** — <reasoning>`. Verdicts get an icon (good ✅ / bad ⛔ / needs_human ⚠️). Deterministic; omitted when empty.

### Command + workflow
- `ksec review --state-dir state [--dry-run]`: load repos + `ai.yaml` + prev `reviews.json`; skip entirely (log) when `!review.enabled` or no LocalAI endpoint; build `OpenAIAssessor`; `review.Run(...)`; save `reviews.json`; log a summary (`review: N bot PRs · M assessed · good/bad/needs-human counts`).
- Workflow "Run pipeline" gains `ksec review --state-dir state $REMEDIATE_DRYRUN` after `remediate`, before `render` (it writes via gh like remediate, so it honors the same dry-run gate). `render` loads `reviews.json`.

## Out of scope

- Auto-**merge** from a verdict (auto-approve only; merge stays under the existing `--automerge` rules).
- Reviewing human-authored PRs (bots only).
- Deep static analysis beyond what the model infers from the diff.
- Re-review on a fixed cadence (only on head-SHA change).

## Testing

- `config`: `review:` block loads into `AIConfig.Review`; `MaxPerRun` defaults to 20.
- `Assessor`/`OpenAIAssessor`: httptest forced-tool-call → parses `{verdict, reasoning}`; AI error / bad JSON / no endpoint → `needs_human_verification`.
- `review.Run` (FakeAssessor + fake GH): bot PRs assessed, non-bot skipped; idempotent (same HeadSHA → carried, no new Assess/comment); cap respected; comment posted with verdict + cc; autoApprove approves only `good`; dryRun makes zero writes; errors → CollectionError; deterministic order.
- `state`: PRReview round-trip.
- render: "🔎 Bot-PR reviews" section accurate + deterministic + omitted when empty; golden regenerated; no raw 64-hex id.
- Manual: a dispatch (dry-run) prints the verdict it would post for AuroraBoot's open bot PRs; live posts the comment.
