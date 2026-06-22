# Run Activity Summary & Signal-Not-Noise — Design

## Problem

A live run produced a dashboard the maintainer found confusing and uninformative:
- **22 routine dependency PRs listed as if they mattered** — `OpenPRs` tracks *any* bot-authored or `dependencies`-labelled PR, so every renovate/dependabot bump shows up even when it has nothing to do with a CVE. The maintainer's verdict: **"we want only PRs tied to CVEs. No noise."**
- **The dashboard tells no story.** With 0 CVE findings, the findings table shows `0`, the AI summary says "No security findings," yet 22 PRs are listed — so the page reads as simultaneously empty and cluttered, with no explanation of **what the run did, what it didn't, and why**. The maintainer: *"I can't tell what happened at all… so we don't even have to look at the logs."*
- **The pipeline is silent.** Phases emit almost no per-run logging, so the CI log is equally opaque.
- A noisy CI warning: `tracking issue upsert failed: gh issue: could not add label: 'security' not found`.

## Goals

- Show **only PRs tied to a CVE** (no routine-bump noise).
- A **deterministic, always-present run-activity summary** on the dashboard: what was scanned, found, tracked, acted on, and **why** — self-sufficient, no log-reading required.
- **Per-phase logging** so the CI log is also legible.
- Kill the `security`-label warning.

## Design

### A. Signal-not-noise: Open PRs tied to CVEs only

`OpenPRs` gains the run's findings and keeps a PR only when it is **CVE-relevant**:
`func OpenPRs(repos []state.Repo, gh ghclient.GitHub, findings []state.Finding) ([]state.TrackedPR, []state.CollectionError)`.

A PR is kept iff any of:
1. **Tied to a finding** — the PR's repo has a finding whose `Package` appears in the PR title (case-insensitive substring, the existing match heuristic). This is the "addresses a CVE" signal.
2. **Explicit `security` label** — a human/maintainer-flagged security PR.
3. **Ours** — `ksec/` branch or `kairos-security-bot` author (our own remediation PR).

Dropped: the blanket "any bot author" and "`dependencies` label" tracking. A renovate/dependabot bump of a package with no finding is **not listed**. When findings are 0, the Open PRs list is correctly **empty** — the true state, no noise. `TrackedPR` keeps the corrected `Source` (dependabot/renovate/bot/human/ksec from the `is_bot` fix).

### B. Deterministic run-activity summary

A new top section **"📋 This run"** rendered from a computed `state.RunActivity`, produced deterministically (no AI, no churn, always present). `render` computes it from `Input` (findings, openPRs, ledger, repos, collectErrors):

- **Scanned:** `N repos — S source-scanned, I image-scanned, K skipped (not source-scannable), E errored`.
- **Findings:** `F findings (C critical / H high / M medium / L low / U unknown)`, or `none`.
- **CVE-related PRs:** `P tracked — <by source: dependabot/renovate/bot/human/ksec>`, or `none`.
- **Remediation:** `opened O · adopted A · superseded S · re-pinned R · needs-human K`, or `no action`.
- **Why** (the key line): a plain-language explanation derived from the numbers, e.g.:
  - 0 findings → *"No CVEs found across N repos this run — nothing to remediate."*
  - findings but dry-run → *"Dry-run: M intents computed, no writes."*
  - errored repos → *"E repos could not be scanned (see collection errors)."*
  - needs-human present → *"K items need a human (build-failed/conflict)."*

The existing AI narrative becomes an **optional flavor line** above the deterministic block (kept only when non-empty/non-generic); the factual summary never depends on the AI being available.

`SummarizeLedger` (the AI coordination call) is fed the activity facts too, so its prose (when present) has context instead of summarizing an empty ledger.

### C. Per-phase logging

Each phase logs one structured line to stderr at the end:
- `collect: 28 repos · 0 findings · 1 skipped · 1 error`
- `triage: 0 findings → no AI call (clean)` / `triage: AI ok, focus=N`
- `remediate: 0 findings → 0 intents` / `remediate: 5 intents (open=2 adopt=1 supersede=1 toolchain=1), N open PRs tied to CVEs`
- `render: dashboard.md (Xkb), site/index.html, issue #N`

So the CI log mirrors the dashboard's activity summary.

### D. Fix the `security`-label warning

Before the tracking-issue upsert, **ensure the labels exist** (`gh label create <name> --force`, best-effort, ignore errors) — or, if creation fails, upsert **without** the missing label rather than warning. Net: no spurious warning, and the label appears when creatable.

## Data flow (changes)

```
collect:  OpenPRs(repos, gh, findings) — CVE-tied filter; logs a summary line
triage/remediate/render: each logs a one-line phase summary
render:   compute RunActivity (deterministic) -> "📋 This run" section (+ optional AI line)
issue:    ensure labels exist before upsert (no warning)
```

## Out of scope

- Shepherding/auto-merging non-CVE dependency PRs (explicitly unwanted — "no noise").
- Changing what counts as a finding (govulncheck reachability + alerts + trivy unchanged).
- Replacing the AI narrative (kept as optional flavor).

## Testing

- `OpenPRs` filter: a dependabot PR whose package matches a finding → kept (source dependabot); a dependabot PR with no matching finding and no security label → dropped; a `security`-labelled human PR → kept; a `ksec/` PR → kept; findings=0 → empty list.
- `RunActivity` compute (pure): given findings/openPRs/ledger/repos/errors → correct counts and the right "why" line for each case (0 findings; findings+dry-run; errored repos; needs-human).
- render: "📋 This run" section present and accurate; deterministic (same input → byte-identical); golden regenerated.
- label fix: ensure-labels best-effort; upsert succeeds without warning when the label is missing (fake GH).
- Manual: a 0-CVE run shows an empty Open PRs list + a clear "no CVEs / nothing to remediate" activity summary; a run with a CVE shows the tied PR(s) and the action taken.
