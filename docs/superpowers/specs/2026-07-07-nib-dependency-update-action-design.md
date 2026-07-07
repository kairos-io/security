# nib Dependency-Update Composite Action — Design

**Date:** 2026-07-07
**Status:** Draft (awaiting review)
**Repo:** `kairos-io/security`

## Problem

We have a working pattern in this repo for AI-assisted maintenance: a workflow
downloads and starts a self-hosted LocalAI binary, installs `nib` (an agentic
CLI), and drives `nib --cli --yolo` to make code edits (today: repairing build
breaks during CVE remediation — see `internal/remediate/nib_agent.go` and
`.github/workflows/security-dashboard.yaml`).

We want to reuse that pattern for a *different, simpler* job: **keep a
repository's dependencies up to date by having nib bump them and open a PR** —
and we want it packaged so *any* of our repos can adopt it with a few lines,
not copy-pasted per repo.

## Goal

A **reusable composite action**, hosted in this repo, that a target repository
calls from a small caller workflow. The action:

1. starts LocalAI (same best-effort binary pattern as the dashboard workflow),
2. has nib update all dependencies to their latest versions,
3. verifies the result builds,
4. opens (or updates) a pull request whose CI actually runs.

First adopter: **`kairos-io/kairos-installer`** (Go). The action is structured
so other Go repos are a ~15-line caller workflow, and non-Go ecosystems are a
future additive case (see Extensibility).

## Validation (done before writing this spec)

The core mechanic was tested end-to-end locally against the exact production
model, in isolated git worktrees, before committing to the design:

- **edgevpn** (`mudler/edgevpn`, Go 1.26): nib bumped `echo/v4` + dozens of
  indirect deps, ran `go mod tidy` (stable on re-run), correctly pruned three
  genuinely-unused direct requires, and `go build`/`go vet` passed.
- **kairos-installer** (`kairos-io/kairos-installer`, Go 1.26): nib bumped
  `kairos-sdk`, `yip`, `ginkgo`, `gomega` and dozens of indirect deps,
  **including a major module-path migration** `containerd` v1 → `containerd/v2`.
  `go build`/`go vet` passed, `go mod tidy` stable.

Both ran on the small production model `gemma-4-e2b-it` (as `ai.yaml` pins),
proving the small model is capable enough for the full multi-step agentic
update. This is the flow the action encodes.

## Non-goals

- CVE-driven or coordinated cross-repo remediation — that already exists in
  `ksec`; this action is deliberately the *simpler* "bump everything" job.
- Running the target repo's test suite inside the action (see Verify gate).
- Auto-merging. The action opens/updates a PR; humans/existing automation merge.
- Non-Go ecosystems in v1 (structural hooks only).

## Architecture

### Packaging: composite action

`.github/actions/update-deps/action.yml` in `kairos-io/security`. Consumers
reference it directly:

```yaml
# <target-repo>/.github/workflows/update-deps.yml
name: Update dependencies
on:
  schedule: [{ cron: "0 5 * * 1" }]   # weekly, Monday 05:00 UTC
  workflow_dispatch:
jobs:
  update:
    runs-on: oracle-vm-16cpu-64gb-x86-64   # needs RAM for LocalAI + model
    concurrency: { group: update-deps, cancel-in-progress: false }
    steps:
      - uses: actions/checkout@v4
        with: { fetch-depth: 0 }
      - uses: actions/create-github-app-token@v2
        id: app-token
        with:
          app-id: ${{ secrets.DEPS_BOT_APP_ID }}
          private-key: ${{ secrets.DEPS_BOT_APP_KEY }}
      - uses: kairos-io/security/.github/actions/update-deps@main
        with:
          language: go
          token: ${{ steps.app-token.outputs.token }}
```

The **caller checks out its own repo**; the action operates on
`github.workspace`. This keeps all reusable logic in one place and each
consumer's footprint minimal.

Rationale for composite action (over reusable `workflow_call` or a `ksec`
subcommand): it embeds directly into a job the consumer controls, needs no
`ksec` binary dependency, and bundles the LocalAI/nib/PR machinery as reusable
steps. (Decision made with the user.)

### Action steps (in order)

1. **setup-go** — `actions/setup-go`, input `go-version` (default `stable`).
2. **Install tooling** — `go install github.com/mudler/nib@<nib-version>`;
   download `yq` and the `local-ai` release binary into a workspace-local
   `bin/` on `PATH` (reuse the exact resolution/download bash from
   `security-dashboard.yaml`, incl. resolving `latest` via the releases API).
3. **Cache + start LocalAI** — `actions/cache` on `models/` keyed by the model
   name; start `local-ai run <model>` best-effort; wait until a real chat
   completion succeeds (bounded by `startup-timeout`, default 20m). If LocalAI
   never becomes ready → go to the fallback in step 5.
4. **Run nib (primary path)** — invoke `nib --cli --yolo` in the workspace with
   env `MODEL` / `BASE_URL` (with `/v1` suffix) / `API_KEY=sk-localai` pointed
   at LocalAI, feeding a single-line task (nib's CLI is line-at-a-time):

   > Update all Go dependencies in this repository to their latest versions by
   > running `go get -u ./...` and `go mod tidy`, then run `go build ./...` to
   > confirm it compiles and fix any errors caused by the updates.

   This is the validated prompt/flow.
5. **Deterministic fallback** *(see Open decision)* — if LocalAI never loaded or
   nib exited without producing a buildable change, run
   `go get -u ./... && go mod tidy` directly so a PR still opens when the model
   is unavailable. nib remains the primary path; this is a safety net.
6. **Verify gate** — `go build ./... && go vet ./...`. On failure, one nib
   repair retry ("fix the build after the dependency update"), then re-verify.
   Tests are **not** run here (kairos-installer/edgevpn suites are heavy/
   networked); the repo's own CI runs them on the PR.
7. **Decide**:
   - No diff to `go.mod`/`go.sum` → exit success, no PR ("already up to date").
   - Verify still failing after retry → **fail the action, open no PR** (never
     push broken deps).
8. **PR** — commit to branch `<branch>` (default `chore/update-deps`) as the bot
   identity; if an open PR from this action already exists on that head, force-
   update its branch (no duplicate); else `git push` + `gh pr create` against
   `<base>` (default: the repo's default branch), using the input `token`.

### Token model

The action takes a **`token` input** and is agnostic about its source. It must
be a token whose pushes/PRs **trigger the target repo's CI**.

- The built-in `GITHUB_TOKEN` is rejected as a default: GitHub suppresses
  workflow runs triggered by `GITHUB_TOKEN` (anti-recursion), so a PR it opens
  would have empty checks.
- **Recommended default: a GitHub App installation token** via
  `actions/create-github-app-token@v2` — an org-owned `kairos-deps-bot` App with
  **Contents: write** + **Pull requests: write**, installed on the target
  repos. App-authored events are *not* under the recursion guard, so CI runs.
  The token is minted fresh per run (JWT signed with the App private key →
  installation token) and expires in ~1h. No personal PAT; not tied to a person;
  scoped per-repo.
- A user PAT is documented as a fallback but discouraged.

Setup (one-time, org admin): create the App, generate a private key, install it
on target repos, store `APP_ID` + `APP_PRIVATE_KEY` as org secrets. For adopting
`kairos-installer`, the App must be installed on the **`kairos-io`** org.

> Note: this repo already uses `secrets.KSEC_BOT_TOKEN` for cross-repo writes in
> the dashboard workflow. If that credential is an App/org-bot token that
> triggers CI, a consumer may reuse it instead of provisioning a new App. If it
> is the bare `GITHUB_TOKEN`, it will not trigger CI and must not be used here.
> (Its exact type is to be confirmed separately; the action stays agnostic.)

### Inputs

| Input | Default | Purpose |
|---|---|---|
| `token` | *(required)* | CI-triggering token (App token recommended; PAT fallback). |
| `language` | `go` | Ecosystem selector. Only `go` implemented in v1. |
| `model` | `gemma-4-e2b-it` | LocalAI model driving nib. Validated. Overridable. |
| `localai-url` | `http://localhost:8080` | If pointed at an already-running server, the action skips download/startup. |
| `localai-version` | `latest` | Pin the LocalAI binary release. |
| `nib-version` | `latest` | Pin nib. |
| `go-version` | `stable` | Go toolchain. |
| `startup-timeout` | `20m` | Max wait for model download + load. |
| `branch` | `chore/update-deps` | PR head branch. |
| `base` | *(repo default branch)* | PR target branch. |
| `pr-title` | `chore(deps): update dependencies` | PR title. |
| `pr-labels` | `dependencies` | Comma-separated labels applied to the PR. |
| `dry-run` | `false` | Do everything except push/open PR (print intended writes). |

### Failure model

- **LocalAI is functionally required** for the primary path (nib *is* the
  engine — unlike triage, there is no deterministic triage fallback). If LocalAI
  cannot load and the deterministic fallback (step 5) is disabled or also fails,
  the action fails clearly rather than opening an empty PR.
- Verify failure after the repair retry → action fails, no PR.
- No dependency changes → success, no PR.

## Extensibility (structural only, not built in v1)

`language` selects a small per-ecosystem bundle: the nib task string and the
verify command(s). Go is implemented:

- task: `go get -u ./... && go mod tidy`; verify: `go build ./... && go vet ./...`.

Future cases (e.g. `node` → `npm update` / `npm run build`; `python` → uv/pip)
are added as new branches of that switch plus their toolchain setup. No such
case is implemented in v1 (YAGNI).

## Open decision (for reviewer)

**Deterministic fallback (step 5): include or not?** This design includes it —
if LocalAI/nib is unavailable, fall back to `go get -u ./... && go mod tidy` so a
PR still opens. It is cheap and makes the action robust to a model outage, while
nib stays the primary path. Recommendation: **include**. Alternative: drop it and
have the action hard-fail when nib can't run (stricter "nib drives everything").

## Files

- `.github/actions/update-deps/action.yml` — the composite action.
- `.github/workflows/update-deps.yml` in each consumer repo — caller (example
  above); first adopter is `kairos-io/kairos-installer`.
- Docs: a short "Adopting dependency updates" section (README or docs/) with the
  caller snippet and the App-token setup steps.
