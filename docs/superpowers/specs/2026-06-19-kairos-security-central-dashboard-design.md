# Kairos Security — Central Dashboard & Autonomous Remediation Engine

- **Date:** 2026-06-19
- **Status:** Approved design (pre-implementation)
- **Repo:** `kairos-io/security` (`kairos-security`)

## 1. Problem

Today `kairos-security` is **image-centric**: it scans the framework container images
listed in `images.json` with trivy/grype/govulncheck via Earthly, fails CI on critical
CVEs, and pings Slack. It knows nothing about the *other* repositories we maintain, does
not aggregate, does not prioritize, and its only "dashboard" is the Actions tab. As a
result the work quickly becomes stale and unhelpful.

We have no single place that answers:

- What CVEs and security PRs exist **across all repos we maintain** (the `kairos-io` org,
  plus external deps like `mudler/*` and `mauromorales/*`)?
- Which of them are **high-severity and need our attention now**?
- Are the **PRs that fix them actually open, correct, and ready to merge**?
- When a single root cause (e.g. a Go stdlib CVE) **waterfalls** across many repos, who
  coordinates the bumps?

We also do not get notified about any of this in a useful, prioritized way.

## 2. Goals

1. A **central engine** in `kairos-security` that gives a single glimpse of the security
   state across every maintained repo.
2. Track four classes of finding per repo: **open security PRs**, **image/binary CVEs**,
   **source/dependency CVEs**, and **GitHub-native security alerts**.
3. **Prioritize** findings and tell us **what to focus on**, using a small self-hosted AI
   model (LocalAI) for the language/judgment layer only.
4. **Autonomously remediate**: open and *maintain* dependency-bump PRs across repos,
   including coordinated **waterfall** bumps when one root cause affects many repos.
5. **Own the PRs it opens**: keep them rebased, react to review comments, and remember
   what it has done across runs (durable memory).
6. Surface everything as a **committed dashboard** in this repo (canonical source of
   truth) and an **auto-updated tracking issue** in `kairos-io/kairos` (human surface).
7. Provide a **dry-run mode** and operate via a **scoped GitHub bot** identity.

## 3. Non-goals (v1)

- Slack notifications (the existing webhook is not extended in v1).
- Per-repo tracking issues (only the one central issue in `kairos-io/kairos`).
- A hosted GitHub Pages render of the dashboard (possible later; JSON is published so it
  is cheap to add).
- Modifying or "fixing" **human-authored** PRs. Human PRs are *tracked* as findings; the
  bot only ever modifies PRs it authored.
- Non-Go ecosystems for autonomous bumps in v1 (other ecosystems are still *reported*).

## 4. Key decisions (from brainstorm)

| Dimension | Decision |
|---|---|
| Engine location | `kairos-security` (brain + source of truth) |
| Architecture | Go CLI (`ksec`) of phased subcommands + thin GHA + `nib` for agentic parts |
| Tracking scope | open security PRs + image/binary CVEs + source/dep CVEs + GH security alerts |
| Repo list | **hybrid**: auto-discover (kairos-diff style) + `repos.yaml` overrides |
| Autonomy | **fully autonomous** bumps + PR maintenance + comment reactions |
| Dry-run | global flag; scheduled runs are **live by default**; PRs/forks force dry-run |
| AI | **self-hosted LocalAI (small model) + `nib --yolo`**, judgment/language only |
| Memory | **committed state files** (`state/*.json`) in this repo |
| Surfaces | committed dashboard (`dashboard.md` + `.json`) here + upserted tracking issue in `kairos-io/kairos` |
| Identity | scoped GitHub bot (App install token / fine-grained PAT) for cross-repo writes |

## 5. Architecture & data flow

A single Go CLI, `ksec`, run by a scheduled GHA workflow as a sequence of phases. Each
phase reads/writes committed **state**, so the pipeline is resumable, auditable in git
history, and every phase is independently testable and runnable locally.

```
                kairos-security repo (the engine + source of truth)
┌─────────────────────────────────────────────────────────────────────┐
│  GHA (scheduled daily) ── runs ksec phases ── commits state+surfaces  │
│                                                                       │
│  1. discover  ─► repo list      (hybrid: auto + repos.yaml override)  │
│  2. collect   ─► raw findings   per repo:                             │
│        • open security PRs (gh API)                                   │
│        • image/binary CVEs (trivy/grype/govulncheck on artifacts)     │
│        • source/dep CVEs   (govulncheck + GH advisories)              │
│        • GH security alerts (Dependabot/code-scan/secret-scan API)    │
│  3. correlate ─► dedupe + build "waterfall" graph                     │
│  4. triage    ─► [nib + LocalAI] severity narrative, "focus on" list  │
│  5. remediate ─► go get -u / mod tidy / bump; open or UPDATE PRs;     │
│        react to review comments  (dry-run gated, blast-radius capped) │
│  6. render    ─► dashboard (md+json) here + upsert tracking issue     │
└─────────────────────────────────────────────────────────────────────┘
         state/ (committed JSON)  ◄──── memory across runs ────►
```

Properties:

- **Phases are independent subcommands** (`ksec <phase> --state-dir ./state [--dry-run]`),
  each consuming the previous phase's state file.
- **State is the memory.** `state/` holds canonical findings plus the bot's PR ledger.
- **GitHub is the actuation target; state is the truth.** Each run reconciles GitHub
  against the ledger (re-read owned PRs, react to new comments, close/reopen as needed).
- **Dry-run** turns every write in `remediate`/`render` into a printed plan; reads still
  happen so the plan is realistic.

## 6. Components

The CLI is one package per concern with narrow, mockable interfaces; `cmd/ksec` wires
phases to subcommands. Suggested layout (`golang-project-layout` conventions):

```
cmd/ksec/                 # CLI entrypoint, flag wiring, phase dispatch
internal/discover/        # repo list (org enum + kairos-init parse + repos.yaml merge)
internal/collect/         # Collector interface + prs/imageCVE/sourceCVE/ghAlerts
internal/correlate/       # pure dedupe + waterfall grouping
internal/triage/          # AI client (nib/LocalAI) + deterministic fallback
internal/remediate/       # reconciliation loop, bumper, PR/comment lifecycle
internal/render/          # dashboard + tracking-issue renderers (golden-tested)
internal/state/           # typed load/save of state/*.json (stable, sorted output)
internal/ghclient/        # GitHub API wrapper (read + scoped write), token redaction
internal/config/          # repos.yaml + ai.yaml + env/flags
```

### 6.1 `discover`
Enumerate `kairos-io` org repos via the GitHub API; parse `kairos-init`'s Makefile +
`go.mod` for pinned `kairos-io/*`, `mudler/*`, `mauromorales/*` deps (the
`kairos-diff.sh` logic, reimplemented in Go); merge with `repos.yaml` (add externals,
pin/exclude, attach per-repo metadata). Output `state/repos.json`.

### 6.2 `collect`
Fan-out across repos. Four collectors behind one `Collector` interface:
- `prs` — gh API: open PRs authored/labeled by renovate, dependabot, and the bot.
- `imageCVE` — trivy/grype on each declared image artifact (reuses today's Earthfile
  toolchain).
- `sourceCVE` — govulncheck + GitHub advisory API on `go.mod`.
- `ghAlerts` — Dependabot / code-scanning / secret-scanning APIs.

A failing collector for one repo records a `collectionError` on that repo and continues.
Output `state/findings.json` (normalized `Finding` rows).

### 6.3 `correlate`
Pure function (no I/O). Collapses the same CVE seen via multiple sources into one finding
(by stable `id`); builds the **waterfall map** grouping findings by shared root cause
(e.g. a `golang.org/x/...` CVE → set of repos whose `go.mod` pulls it) with a
`suggestedBump`. Output `state/correlated.json`. Fully unit-testable.

### 6.4 `triage` (AI/judgment layer)
Sends correlated findings to `nib --yolo` against the small LocalAI model to produce a
prioritized **"what to focus on"** list, a human summary per high/critical item, and a
one-line rationale per waterfall group. **Language/judgment only** — it never decides
version numbers or edits code. Output `state/triage.json`. Best-effort: on AI failure,
falls back to deterministic severity sort + templated summaries and flags "AI
unavailable this run" (see §10).

### 6.5 `remediate` (autonomous, stateful actuator — §8)
Reconciliation loop over actionable findings/waterfall groups: compute deterministic bump
(`go get -u <pkg>@<fixed>` + `go mod tidy`), branch, open **or update** a PR via the
scoped bot, reconcile against the **PR ledger**, and react to new review comments on
owned PRs. Verifies build/test before any push. Dry-run gated; blast-radius capped.
Updates `state/ledger.json`.

### 6.6 `render`
Projects state to both surfaces (§9). Writes `dashboard.md` + `dashboard.json` here;
upserts the single tracking issue in `kairos-io/kairos`. Pure formatting over state;
dry-run prints instead of writing.

## 7. Data model (committed state = memory)

All under `state/`, pretty-printed JSON with stable key order and sorted arrays so diffs
are clean. Each phase owns its file; later phases never mutate earlier files.

```jsonc
// state/repos.json  — from discover
[{ "repo": "kairos-io/immucore", "kind": "dep",
   "branch": "master", "criticality": "high",
   "artifacts": [{ "type": "image", "ref": "quay.io/kairos/immucore:..." },
                 { "type": "go", "modpath": "." }] }]

// state/findings.json  — from collect (one row per observation)
[{ "id": "sha256(repo+type+cve+pkg)",
   "repo": "kairos-io/immucore", "type": "sourceCVE",
   "cveID": "CVE-2025-1234", "ghsa": "GHSA-...", "ecosystem": "go",
   "package": "golang.org/x/net", "currentVersion": "0.30.0",
   "fixedVersion": "0.33.0", "severity": "high",
   "source": "govulncheck", "links": ["..."],
   "firstSeen": "2026-06-19", "lastSeen": "2026-06-19" }]

// state/correlated.json  — from correlate
{ "findings": [ /* enriched findings, each with optional waterfallGroup ref */ ],
  "waterfall": [{ "id": "go-stdlib-CVE-2025-1234",
    "rootCause": "golang.org/x/net < 0.33.0", "ecosystem": "go",
    "severity": "high",
    "affectedRepos": ["kairos-io/immucore","kairos-io/kairos-agent"],
    "suggestedBump": { "package": "golang.org/x/net", "to": "0.33.0" } }] }

// state/triage.json  — from triage (AI output, text only)
{ "generatedAt": "2026-06-19", "model": "<local-model-id>", "aiAvailable": true,
  "focus": ["go-stdlib-CVE-2025-1234"],
  "summaries": { "go-stdlib-CVE-2025-1234": "x/net CVE affects 2 repos; ..." },
  "narrative": "This cycle the dominant risk is ..." }

// state/ledger.json  — PR ledger == the bot's long-term memory
[{ "key": "kairos-io/immucore:go-stdlib-CVE-2025-1234",
   "repo": "kairos-io/immucore", "branch": "ksec/bump-xnet-0.33.0",
   "prNumber": 412, "state": "open|merged|closed|conflicted|checks-failing",
   "createdRun": "2026-06-10", "lastActionRun": "2026-06-19",
   "bump": { "package": "golang.org/x/net", "to": "0.33.0" },
   "seenComments": ["IC_comment_id_1","IC_comment_id_2"],
   "history": [{ "run":"2026-06-10","action":"opened" },
               { "run":"2026-06-15","action":"replied","commentId":"IC_..." },
               { "run":"2026-06-19","action":"rebased" }] }]
```

- `firstSeen`/`lastSeen` drive **CVE aging** (long-open findings escalate in `focus`).
- `seenComments` makes comment reactions **idempotent** (react once per comment).
- `history` is the human-auditable trail.

## 8. Remediation + comment-reaction lifecycle

`remediate` is a **reconciliation loop**, not a fire-once script. Per run, for each
actionable finding / waterfall group:

```
actionable item
  └─ ledger entry exists?
       ├─ no  → OPEN: branch, go get -u <pkg>@<fixed>, go mod tidy,
       │        build/test, push, create PR (marker + AI summary) → record
       └─ yes → reconcile PR state (re-read via API)
                 ├─ merged?      → mark finding resolved, close ledger entry
                 ├─ conflicted / base moved? → rebase: redo bump on fresh base, force-push
                 └─ new review comments (not in seenComments)?
                      → [nib+LocalAI] classify: request-change | question | nack | approve
                      → apply: adjust bump (explicit version) & push,
                               OR post drafted reply,
                               OR close + record if nack
                 → append history, update seenComments, write ledger.json
```

Safety / idempotency rules:

- **Identity by marker + branch.** Every bot PR carries an HTML-comment marker
  `<!-- ksec:key=... -->` and a `ksec/`-prefixed branch. The bot only touches PRs it
  owns; human PRs are tracked as findings, never modified.
- **Bump is deterministic; only the *reaction* is AI.** The model classifies comment
  intent and drafts prose and may *suggest* a version when a maintainer says "pin to X",
  but the change itself is `go get <pkg>@<explicit-version>`, verified by build/test.
- **Verify before push.** Every bump runs the repo's build/test. On failure the bot does
  **not** push a broken PR — it records `checks-failing` and surfaces it for humans.
- **Waterfall coordination.** One group fans out to N per-repo PRs sharing the group key,
  cross-linked in their bodies and in the tracking issue, so the whole front is visible.
- **Blast-radius guard.** Configurable cap on new PRs per run (default 10); overflow is
  logged and deferred to the next run.
- **Dry-run** short-circuits every git/GitHub write to a printed plan.

## 9. Surfaces

### 9.1 Dashboard in this repo (canonical)
- `dashboard.json` — full machine-readable snapshot (correlated state + ledger status +
  triage focus). The consumable API.
- `dashboard.md` — human view generated from the JSON: **🔥 Focus now** (triage order),
  **Waterfall fronts** (grouped CVEs + per-repo PR status), **Per-repo table** (repo ·
  open CVEs by severity · open security PRs · oldest finding age), **Bot PR ledger**.
  Committed each run → `git log dashboard.md` is the security timeline.

### 9.2 Tracking issue in `kairos-io/kairos`
- A single issue located by a hidden marker `<!-- ksec:dashboard -->`; create if absent,
  else rewrite body. Body = `dashboard.md` + "last updated / run link" footer. Labels
  `security`, `kairos-security-bot`.
- Body-only rewrite (no comment spam). The bot posts a **comment** only on
  notification-worthy state changes (new critical, new waterfall front, a PR went green
  and is ready to merge). Threshold configurable.
- The issue is permanent and self-updating; never auto-closed.

## 10. Error handling

Degrade, never crash the whole run:

- **Per-repo / per-collector isolation** — one failure records a `collectionError` and
  continues; dashboard shows "⚠️ N repos failed to scan."
- **AI is best-effort** — on unreachable/garbage LocalAI output, `triage` falls back to
  deterministic severity sort + templated summaries (`aiAvailable: false`); `remediate`
  comment-reaction is skipped and recorded `needs-human` rather than guessing. The
  pipeline never blocks on the model.
- **Remediation is transactional per item** — failed build/test → no push, `checks-failing`;
  failed GitHub write → ledger keeps prior state, error logged, retried next run.
- **State is the recovery mechanism** — every phase commits state; a mid-run failure
  leaves consistent partial state and the next run reconciles forward. GHA `concurrency`
  group prevents overlap (no external locks).
- **Secrets never logged** — token redaction in all error/log output.

## 11. GHA orchestration

`.github/workflows/security-dashboard.yaml`, scheduled daily + `workflow_dispatch`:

```
jobs:
  run:
    1. checkout (full history — state lives in git)
    2. start LocalAI as a service (small model from config); wait for readiness
    3. install pinned tools: ksec, nib, trivy, grype, govulncheck, go
    4. ksec discover  --state-dir state
    5. ksec collect   --state-dir state          # fan-out, the slow phase
    6. ksec correlate --state-dir state
    7. ksec triage    --state-dir state --ai-url $LOCALAI_URL
    8. ksec remediate --state-dir state $DRYRUN   # cross-repo writes
    9. ksec render    --state-dir state $DRYRUN
   10. commit & push state/ + dashboard.* back to main (idempotent: skip if unchanged)
```

- `workflow_dispatch` exposes a `dry_run` input. **Scheduled runs are live by default;**
  PRs/forks force dry-run. `DRYRUN` env propagates `--dry-run` to write phases.
- LocalAI runs as a runner service; the model is pulled/cached. Only `triage` and
  `remediate`'s comment-judgment hit `LOCALAI_URL`; everything else is offline.
- The bot's own state commits are excluded from triggering other workflows.

## 12. AI configuration handles (LocalAI + nib)

AI is fully self-hosted and **configurable without code changes**, via `ai.yaml` in this
repo with env/secret overrides. This covers: which LocalAI model to run, how `nib` is
pointed at it, and how `nib`/LocalAI themselves are installed.

```yaml
# ai.yaml
localai:
  version: "vX.Y.Z"            # pinned LocalAI release (overridable: LOCALAI_VERSION)
  model:
    name: "<small-model-id>"   # e.g. a small instruct model from the LocalAI gallery
    gallery: "<gallery-ref>"   # gallery entry / model URI to preload
    quant: "<quant-tag>"       # optional quantization selector
  endpoint: "http://localhost:8080"   # overridable: LOCALAI_URL
  startupTimeout: "5m"
nib:
  version: "vX.Y.Z"            # pinned nib release (overridable: NIB_VERSION)
  mode: "yolo"                 # --yolo harness
  model: "<small-model-id>"    # which served model nib targets (defaults to localai.model.name)
  endpoint: "${localai.endpoint}"
  maxTokens: 4096
  temperature: 0.2
```

Handles required:

1. **Model selection / preload** — `localai.model.*` selects and preloads the small
   model into LocalAI at workflow start (gallery name + optional quant). Changing the
   model is a one-line config edit; CI override via `LOCALAI_MODEL`.
2. **nib → LocalAI wiring** — `nib.endpoint` + `nib.model` point the harness at the
   served model; defaults derive from the `localai` block so they cannot drift.
3. **Install/pin nib and LocalAI** — both versions are pinned in `ai.yaml` and consumed
   by the workflow install step (and a `renovate` custom-manager entry can track them, as
   the Earthfile already does for trivy). Overridable via `NIB_VERSION` / `LOCALAI_VERSION`.
4. **Endpoint override** — `LOCALAI_URL` lets a run target an external/already-running
   LocalAI instead of the in-runner service (useful for local dev).

`triage` and `remediate` read this config through `internal/config`; the actual model
calls go through `nib --yolo`. The Go code treats the AI as an interface
(`AIClient`) so it is mocked in tests and swappable.

## 13. Configuration files (summary)

| File | Purpose |
|---|---|
| `repos.yaml` | hybrid repo overrides: add externals, pin/exclude, per-repo artifacts/criticality |
| `ai.yaml` | LocalAI model + nib handles (§12) |
| `state/*.json` | committed memory (repos, findings, correlated, triage, ledger) |
| `dashboard.md` / `dashboard.json` | committed canonical dashboard surface |

## 14. Testing strategy

- **Pure phases** unit-tested with table-driven tests + golden files: `correlate`
  (dedupe/waterfall), `render` (state → md/json golden), `discover`
  (Makefile/go.mod parsing fixtures). No network.
- **Collectors** tested against recorded GitHub API fixtures + sample
  trivy/grype/govulncheck JSON (interface-mocked).
- **`remediate`** tested as a state machine: feed a ledger + simulated PR/comment states,
  assert intended actions (in dry-run, actions are data → trivial assertions). AI call
  mocked behind `AIClient`.
- **End-to-end dry-run** in CI over a small fixture repo set: full pipeline produces a
  coherent plan + valid surfaces without writing anywhere.
- `testify` for assertions; golden-file pattern for renderers.

## 15. Migration / compatibility

- The existing image-CVE scanning (Earthfile trivy/grype/govulncheck, `images.json`) is
  **reused** by the `imageCVE` collector rather than discarded; framework images become
  one set of artifacts in `repos.yaml`.
- The current `scan.yaml` / `autobump.yaml` / `automerge.yaml` workflows remain until the
  new pipeline reaches parity, then are retired or folded in (decided during planning).

## 16. Open items for the implementation plan

- Exact GitHub App vs fine-grained PAT choice and the minimal permission set per target
  repo (`contents:write`, `pull_requests:write`, `issues:write`).
- Concrete small model id + gallery reference and its resource footprint on the runner.
- Whether non-patch bumps require an extra human-approval label before merge.
- Retirement plan for the legacy workflows once parity is confirmed.
