# Kairos Security

[![Security Dashboard](https://github.com/kairos-io/security/actions/workflows/security-dashboard.yaml/badge.svg?branch=main)](https://github.com/kairos-io/security/actions/workflows/security-dashboard.yaml)

Central place for the security state across the repositories we maintain — the
`kairos-io` organization plus external dependencies (`mudler/*`,
`mauromorales/*`). It aggregates open security PRs, image/binary CVEs,
source/dependency CVEs, and GitHub security alerts into a single dashboard, and
tells us what to focus on.

## What it does

`ksec` is a Go CLI run on a schedule by GitHub Actions as a sequence of phases.
Each phase reads/writes committed JSON state under `state/`, so every run is
auditable in git history.

```
discover  → build the tracked-repo list (kairos-io org + kairos-init deps + repos.yaml), dropping archived repos
collect   → per repo: open security PRs, image CVEs (trivy), source CVEs (govulncheck), GitHub alerts, hadron component-manifest CVEs (OSV.dev + NVD)
correlate → dedupe findings + build the "waterfall" graph (one Go CVE → the set of affected repos)
triage    → prioritize + write the "focus now" summary (self-hosted LocalAI via nib; deterministic fallback)
render    → dashboard.md + dashboard.json (committed), site/index.html (GitHub Pages), tracking issue in kairos-io/kairos
```

This repository is **read-only** against every other repo: the only GitHub
write is upserting the single tracking issue in `kairos-io/kairos`. A
`--dry-run` flag turns every write into a printed plan.

> Autonomous remediation (coordinated dependency bumps, PR creation, and
> reacting to review comments — including the "waterfall" case where one Go CVE
> needs bumps across many repos) is a planned follow-up (Plan 2) and is not part
> of this pipeline yet.

## Running locally

```sh
go build -o ksec ./cmd/ksec
ksec discover  --state-dir state
ksec collect   --state-dir state
ksec correlate --state-dir state
ksec triage    --state-dir state
ksec render    --state-dir state --dry-run   # --dry-run prints intended writes
```

Each phase consumes the previous phase's state file, so any phase can be re-run
in isolation against committed state.

## Configuration

- **`repos.yaml`** — hybrid repo overrides. Auto-discovery (the `kairos-io` org
  plus dependencies parsed from `kairos-init`) is the base; this file *adds*
  external repos, *excludes* repos, and attaches per-repo metadata (artifacts to
  scan, branch, criticality).
- **`ai.yaml`** — LocalAI + `nib` handles: which small model to run/preload, how
  `nib` points at the LocalAI endpoint, and the pinned tool versions. Overridable
  via `LOCALAI_URL`, `LOCALAI_MODEL`, `LOCALAI_VERSION`, `NIB_VERSION`.
- **`hadron-components.yaml`** — maps `kairos-io/hadron`'s published component
  manifest packages to how they're checked for CVEs (OSV.dev ecosystem/package,
  optional NVD CPE fallback). NVD lookups are optionally authenticated via
  `NVD_API_KEY` (raises the rate limit from 5 to 50 requests/30s; unset works,
  just slower).

## Pipeline

`.github/workflows/security-dashboard.yaml` runs the five phases on a schedule
(live by default; dry-run on `workflow_dispatch` input or fork PRs), starts
LocalAI as a runner service for triage, commits the updated state + dashboards
back to `main`, and publishes the HTML dashboard to GitHub Pages.

## Design docs

- Spec: [`docs/superpowers/specs/2026-06-19-kairos-security-central-dashboard-design.md`](docs/superpowers/specs/2026-06-19-kairos-security-central-dashboard-design.md)
- Plan: [`docs/superpowers/plans/2026-06-19-central-dashboard-read-only-pipeline.md`](docs/superpowers/plans/2026-06-19-central-dashboard-read-only-pipeline.md)
