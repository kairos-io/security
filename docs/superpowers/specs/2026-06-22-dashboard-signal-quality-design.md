# Dashboard Signal Quality — Design

## Problem

The committed dashboard is dominated by noise and misses real vulnerabilities. Root causes, found by inspecting committed state (`state/findings.json`: 14 findings, all `type:"pr"`, `errors:null`) and scanning the repos directly:

1. **The source scanner silently scans nothing for modern repos.** Tracked repos (e.g. `mudler/edgevpn`, `mudler/entities`) declare `go 1.26`; CI pins `setup-go: "1.22"`. govulncheck must *build* the code, so Go 1.22 can't analyze a `go 1.26` module. `govulncheckRunner` only treats an **empty** stdout as failure (`if err != nil && len(out) == 0`), but govulncheck emits config/progress JSON before a build failure — so a build failure returns `(out, nil)` and is parsed as **zero findings with no error**. A direct govulncheck run (Go 1.26) finds **9 reachable vulnerabilities in edgevpn** (8 Go-stdlib, fixed in go 1.26.2–1.26.4; plus `GO-2024-3218` in `go-libp2p-kad-dht@v0.40.0`, no fix). All were reported as nothing.

2. **Open PRs are rendered as security findings.** The only collector producing output is `PRs` (`type:"pr"`, hardcoded `severity:"unknown"`). Routine/stale dependency PRs (`bump codecov-action 5.5.0→5.5.1`; `entities #11` bumping `x/net` to `0.7.0` when the repo is already on `v0.54.0`) fill the findings list and Focus, with empty severity buckets (`0 0 0 0` but Total `11`).

3. **Human-facing output shows opaque ids.** Focus renders the SHA-256 finding id + a generic AI line (`Finding in mudler/edgevpn`) instead of the finding's `title`/`url`, which are present in the data.

4. **govulncheck findings have no severity** (hardcoded `"unknown"`) → buckets never populate, and once the scanner works it would also emit **non-reachable** "required but not called" vulns (entities had 15), re-creating noise.

## Goals

Make the dashboard a true security signal: real scan results with severities, routine PRs moved to their own list, and human-readable links instead of hashes.

## Design

### A. Fix the source scanner (collection correctness)

- **CI Go toolchain:** the workflow uses `go-version: "stable"` (latest). Go is backward-compatible, so a newer toolchain builds every tracked repo (incl. `go 1.26`) and ksec itself (`go 1.22`). This is the single change that makes govulncheck actually analyze modern repos.
- **Loud build failures:** `govulncheckRunner` captures stderr and decides via a pure, testable helper `classifyGovulncheck(stdout, stderr []byte, runErr error) ([]byte, error)`:
  - `runErr == nil` → return `(stdout, nil)`.
  - `runErr != nil` and stdout contains at least one govulncheck record (`"osv"` or `"finding"` JSON object) → vulns-found is normal → return `(stdout, nil)`.
  - `runErr != nil` and stdout has no records → a real failure (build/load) → return `(nil, fmt.Errorf("govulncheck: %v: %s", runErr, stderr))`, surfaced as a `CollectionError` on the dashboard.

### B. Reachability + severity (collection quality)

In `SourceCVE.Collect` (`internal/collect/source.go`):
- **Reachability filter:** keep only findings govulncheck marks as *called* — the trace's most-specific frame has a function symbol (`trace[0].function != ""`). This is the "Your code is affected by N vulnerabilities" set; it drops the "required but not called" noise.
- **Severity:** parse the OSV record's severity via a pure helper `severityFromOSV(databaseSpecificSeverity string) string`: normalize a present value (`CRITICAL→critical`, `HIGH→high`, `MODERATE|MEDIUM→medium`, `LOW→low`); when absent, default to **`high`** (a reachable vulnerability with no severity data is actionable, not "unknown"). stdlib findings (`Package == "stdlib"`) continue to flow to the existing Plan-4c toolchain-bump path.

### C. PRs become their own list, not findings

- **Stop feeding findings:** remove `collect.PRs` from the findings collector list. No more `type:"pr"` findings, no Focus pollution, no fake buckets.
- **New tracked-PR state:** `state.TrackedPR{Repo string; Number int; Title, Author, URL, Source string}`, persisted as `state/openprs.json` (sorted by repo then number — deterministic). `Source` is the existing classification (`renovate|dependabot|ksec|human`).
- **Gather:** `collect.OpenPRs(repos []state.Repo, gh ghclient.GitHub) ([]state.TrackedPR, []state.CollectionError)` lists open PRs per repo and keeps the tracked set (the existing `isSecurityPR` filter: bot authors or `security`/`dependencies` labels). The `collect` command writes `openprs.json`.
- **Render section:** the dashboard gains **"📋 Open PRs"**, grouped by repo, each line `- [#N title](url) — source`. This is distinct from "🤖 Bot PR ledger" (which tracks ksec's *own* remediation PRs from `ledger.json`).

### D. Render links, not hashes

In `internal/render` (markdown + HTML):
- Build a `byID` map from `in.Correlated.Findings`. **Focus now** renders each focus id as the finding's **`title` linked to its `url`** (`- [title](url)`). When a finding has no `url` but is a `sourceCVE` with a `CVEID`/advisory id, link `https://pkg.go.dev/vuln/<id>`. Fallback: the title text (never the raw id). An optional AI summary may follow as ` — <summary>` only when it is non-generic.
- The opaque SHA-256 id never appears in human-facing markdown/HTML.

## Data flow

```
discover → collect → correlate → triage → remediate → render
                │                                          │
                ├─ findings.json (sourceCVE/ghAlert/imageCVE, severities)   ┘ (PRs removed)
                └─ openprs.json  (tracked open PRs)  ───────────────────────→ "📋 Open PRs"
```

## Out of scope

- CVSS-vector → severity parsing (Go vuln DB rarely ships clean scores; the `high` default covers it).
- Per-repo toolchain pinning (latest `stable` suffices; backward compatibility holds).
- Changing the triage AI prompt (its generic summaries are a symptom; once findings carry titles/URLs and PRs leave the findings set, Focus is readable without the AI line).

## Testing

- `classifyGovulncheck` and `severityFromOSV`: pure unit tests (table-driven).
- `SourceCVE.Collect`: fixture govulncheck JSON with reachable + non-reachable + severity-bearing OSV → asserts filter + severity.
- `collect.OpenPRs`: fake GitHub → asserts tracked filter, source classification, deterministic order.
- `render`: Focus shows `[title](url)` not the id; "📋 Open PRs" section appears; goldens regenerated.
- Manual proof: after the CI Go fix, a live scan surfaces edgevpn's `GO-2024-3218` + the stdlib/toolchain set.
