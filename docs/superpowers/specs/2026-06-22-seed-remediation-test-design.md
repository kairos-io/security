# Seed Remediation Test — Design

## Problem

There's no way to exercise the remediation pipeline on demand. A real run only acts when govulncheck/alerts produce a CVE finding, so when scans are clean (0 findings) the operator can't validate that plan → fork → PR (and supersede/toolchain) actually works. They asked for a "test for the remediation pipeline."

## Goal

A flag that injects a **synthetic finding** so the planner produces a real bump intent and the executor shows (dry-run) or opens (live) the exact PR it would create — exercising the full path on demand, honoring `--dry-run`, `--max-prs`, and the fork rules.

## Design

`ksec remediate` gains a repeatable flag:

```
ksec remediate --seed <owner/repo>=<package>@<version> [--seed …] [--dry-run]
# e.g. ksec remediate --seed mudler/edgevpn=golang.org/x/net@0.33.0 --dry-run
```

- **`ParseSeed(spec string) (state.Finding, error)`** (pure, in `internal/remediate/seed.go`): parses `owner/repo=package@version` (split on the first `=`, then the first `@`) into a synthetic finding the planner will action:
  `Finding{Repo, Package, FixedVersion: version, Type: "sourceCVE", Ecosystem: "go", Severity: "high", Source: "seed", CVEID: "SEED", Title: "synthetic seed finding (remediation test)", ID: "seed:<repo>|<package>"}`.
  These are exactly the fields `actionable()` requires (`sourceCVE` + `go` + non-stdlib package + fixed version). Malformed specs (missing `=`/`@`, empty parts) return a clear error.
- **Wiring** (`cmd/ksec/main.go` `newRemediateCmd`): after loading the correlated state and before `Plan`, parse each `--seed` and append it to `c.Findings`; log `remediate: injected N seed finding(s)`. Everything downstream is unchanged — `Plan` turns the seed into an `open` (or `adopt`/`supersede`) intent, capped by `--max-prs`; `Run` executes via the existing executor, which already honors `--dry-run` (prints, no writes) and forks external repos.

No special-casing: a seed is just a finding injected before planning, so it flows through the real engine.

## Behavior

- **Dry-run** (`--dry-run`): prints the intent(s) and the PR it *would* open (a fork PR for an external repo like `mudler/edgevpn`) — zero writes. The primary safe test.
- **Live**: opens one real demonstration PR (subject to `--max-prs`), on the bot's fork for external repos.
- The seed repo should be tracked in `repos.yaml` so fork-detection (`ShouldFork`) and PR-matching work; seeding an untracked repo push-directs and may fail (documented).

## Out of scope

- Severity/type options on the seed (always `high` `sourceCVE`) — keep the flag simple.
- Persisting seeds to state (they're transient, planning-time only).
- A separate test subcommand (a flag on `remediate` is enough).

## Testing

- `ParseSeed`: valid `mudler/edgevpn=golang.org/x/net@0.33.0` → finding with the right Repo/Package/FixedVersion/Type/Ecosystem/Severity; package paths with slashes preserved; malformed (`no-eq`, `repo=pkg` without `@`, empty repo/pkg/version) → error.
- The injected seed satisfies `actionable()` (so the planner produces a target) — assert via a planner test that `Plan` with a seed finding yields an `open`/`adopt` intent.
- Manual: `ksec remediate --seed mudler/edgevpn=golang.org/x/net@0.33.0 --dry-run --state-dir state` prints the planned fork PR.
