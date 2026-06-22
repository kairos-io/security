# Proactive Reactivity & Status Surfacing — Design

## Problem

Three gaps surfaced from a live run:

1. **Bot PRs misclassified as human.** `cluster-api-provider-kairos#38` was opened by a bot but flagged `human`. `classifySource`/`prSource` recognize only a hardcoded allowlist (`renovate[bot]`, `dependabot[bot]`, `kairos-security-bot`); any other GitHub App bot falls through to `human`. The same allowlist gates `isSecurityPR`, so non-allowlisted bots may not even be tracked.

2. **No reaction to a stuck external PR.** `#38` is conflicted, but `ResolveConflict` (Plan 4c) only acts on our own `ksec/` branches (we must never force-push someone else's branch). So an adopted external PR that goes conflicted/stale gets endless nudges and no resolution. The bot should instead judge whether the fix still applies and, if so, **supersede it with a fresh `ksec/` PR**.

3. **One repo's failure makes the dashboard unreadable.** `kairos-must-burn` is a GTK4/cgo app; govulncheck can't build it headless (missing `glib`/`gtk4` system libs) and emits a ~50KB error wall that floods both the dashboard and the committed `findings.json`.

The unifying direction: **the bot acts proactively where it is safe, and surfaces clear status everywhere so a human can decide what to pick.**

## Goals

- Recognize any bot, not three.
- Proactively supersede a conflicted/stale external PR that is still relevant, on our own branch, transparently.
- Make every PR/finding carry a readable **status + recommended/taken action** so a human can choose.
- Keep an un-scannable repo an honest one-line status, not an error wall.

## Design

### A. Generalize bot detection (proactive recognition)

In `internal/remediate/matcher.go` and `internal/collect/prs.go`:
- `func isBotLogin(login string) bool` → `strings.HasSuffix(login, "[bot]")` OR a known-bot membership. GitHub App bots all carry the `[bot]` suffix, so this is authoritative for App-authored PRs.
- `classifySource`/`prSource` return: `ksec` (own), `renovate`, `dependabot` (the two we special-case for messaging), `bot` (any other `[bot]`), else `human`.
- `isSecurityPR` treats any `isBotLogin` author as a bot (so non-allowlisted bots are tracked), keeping the `security`/`dependencies` label path for human PRs.

### B. Proactive supersede of a stuck external PR

When an external PR we adopted (tracked, `source != ksec`) is **conflicted** and the underlying fix is **still relevant**, the bot opens a clean `ksec/` PR that supersedes it.

- **Detection (Reconcile):** for an adopted ledger entry, `Reconcile` already fetches the PR view; extend it to read `mergeable`. A `CONFLICTING` external PR sets the entry `Blocked: "upstream-conflict"` (status, surfaced) — never force-pushing the foreign branch.
- **Relevance (planner, deterministic):** the underlying finding/target still exists AND the required fixed version still exceeds what the repo ships (`compareVersions`). If the finding is gone or already satisfied, the adoption is moot → mark `stale`, stop nudging.
- **Supersede (planner → executor):** for a still-relevant entry `Blocked: "upstream-conflict"`, emit `IntentSupersede{Repo, Package, Bump, Supersedes: <extPR#/URL>}`. The executor opens a fresh `ksec/` bump PR (clone → `go get`→ `verifyOrRepair` → push → `gh pr create`), posts a comment on the external PR ("Superseded by <ksec PR URL> — the original had unresolved conflicts"), and records the new entry with `Supersedes` set and the old adoption marked `superseded`. If the clean bump fails to build (and the agent can't repair) → `needs-human` (no PR, surfaced).
- **Safety:** supersede only ever creates/pushes our own `ksec/` branch and *comments* on the foreign PR; it never edits or force-pushes the foreign branch. Capped by `--max-prs` like other new PRs. Dry-run prints, no writes.

New ledger fields: `Supersedes string` (ext PR URL this entry replaces) and a `superseded`/`stale` state value; `Blocked: "upstream-conflict"`.

### C. Surface status so a human can pick

The dashboard makes the *state and the action* explicit:
- **📋 Open PRs** gains a **Status/Action** column per PR: `tracked`, `conflicted → superseding`, `superseded by #X`, `stale`, `needs human`, plus the corrected source (`bot`/`renovate`/`dependabot`/`human`/`ksec`).
- **🤖 Bot PR ledger** shows, per entry, the state incl. `superseded`/`stale`/`upstream-conflict`/`needs-human` and any `Supersedes`/`PinTarget` link, so a human sees what the bot did and what still needs a decision.
- A compact **"needs human"** roll-up lists exactly the items awaiting a human pick (conflicts the bot couldn't resolve, build-failed, ambiguous staleness).

### D. Readable errors + honest un-scannable status

- **Truncate** the govulncheck error in `ClassifyGovulncheck`: collapse whitespace, keep the first meaningful line, cap ~240 chars + `… (truncated)`. Keeps `findings.json` small and the dashboard readable. Render defensively caps any collection-error message too.
- **Per-repo source-scan opt-out:** `repos.yaml` gains `scan: {source: false}` (default true). `kairos-must-burn` sets it (cgo/GTK4, needs system libs). The collector skips govulncheck for such repos and the dashboard shows status **`skipped: not source-scannable`** rather than an error — an honest, readable status the maintainer can act on (install deps / accept).

## Data flow (additions)

```
collect:  prs.go classifies bots correctly; source-scan opt-out honored; truncated errors
reconcile: adopted ext PR CONFLICTING -> Blocked="upstream-conflict"
plan:      relevant + upstream-conflict -> IntentSupersede ; irrelevant -> stale
execute:   IntentSupersede -> open ksec/ PR + comment on ext PR + mark old superseded
render:    Open PRs status/action column + ledger states + needs-human roll-up + skipped status
```

## Out of scope

- Editing or force-pushing foreign PR branches (we only comment + supersede on our own branch).
- Auto-merging the superseding PR (existing `--automerge` rules still apply).
- Installing system libraries in CI for cgo/GUI repos (opt-out + honest status instead).
- Author-type via GitHub API (the `[bot]` suffix is sufficient; revisit if a bot ever lacks it).

## Testing

- `isBotLogin`/`classifySource`/`prSource`: table tests incl. a novel `foo[bot]` → `bot` (not `human`), and tracking of a non-allowlisted bot.
- planner: adopted entry with `Blocked:"upstream-conflict"` + still-relevant target → `IntentSupersede` (capped); irrelevant → `stale`, no supersede.
- Reconcile: a `CONFLICTING` adopted ext PR → `Blocked:"upstream-conflict"`, no foreign push.
- executor `Supersede`: dry-run no writes; live opens ksec PR + comments ext PR + marks old `superseded`; ksec-branch guard; build-fail → needs-human.
- `ClassifyGovulncheck` truncation: a 50KB stderr → concise capped message.
- collect opt-out: `scan.source:false` repo is not govulncheck'd and renders `skipped`.
- render: Open PRs status/action column; needs-human roll-up; superseded/stale/skipped states; goldens.
- Manual: re-run surfaces `#38` correctly classified as a bot, conflicted, and (if still relevant) superseded by a ksec PR.
