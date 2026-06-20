# Coordinated Remediation Engine — Design (Plan 4)

- **Date:** 2026-06-20
- **Status:** Approved design (pre-implementation)
- **Repo:** `kairos-io/security` (`kairos-security`)
- **Builds on:** Plan 1 (dashboard), Plan 2 (bump PRs + ledger), Plan 3 (comment reactions). Extends the existing `internal/remediate` package — NOT a new phase.

## 1. Problem

Today `ksec remediate` opens and maintains one dependency-bump PR per `repo|package` for direct Go `go.mod` dependencies that have a known fixed version, reconciles their state, and reacts to review comments. It does **not**:

- Notice that **dependabot / renovate / a human** already opened a PR for a CVE — so it can duplicate work instead of adopting and driving the existing PR.
- Help **get PRs merged**.
- Handle the **first-party dependency cascade**: a CVE fixed in `kairos-io/kairos-sdk` requires every repo that *consumes* the sdk to bump the sdk to pick up the fix — and their consumers in turn. A single fix fans out into a chain of PRs across the dependency graph.
- Bump the **Go toolchain** when a stdlib CVE is fixed by a Go release.
- **Repair** a bump that breaks `go build`, or resolve a PR conflict.
- Present a coherent **status** of which CVEs are open, which PR addresses each, and what is blocked / missing / needs a human.

We want remediation to act **proactively and status-aware**: track every open CVE and the PR that addresses it, drive PRs to merge, create what is missing, and propagate fixes down the first-party dependency graph — "things must flow," without waiting for human approval to progress.

## 2. Goals

1. **Findings ↔ PRs linkage.** For every open CVE, identify the PR that addresses it (ours / dependabot / renovate / human), record the link, source, and state in the ledger; create one only when none exists. Never duplicate.
2. **Drive to merge.** Nudge addressing PRs by default; with an opt-in flag, auto-merge green, unblocked PRs (ours or external).
3. **First-party cascade.** Derive the dependency graph from each tracked repo's `go.mod`; when a first-party module's fix is available on its default branch, immediately open **pseudo-version** bump PRs in its consumers (do not wait for a release tag), note in the PR body that a maintainer should tag, and **re-pin to the tag** once released. Recurse to consumers-of-consumers.
4. **Toolchain bumps.** For a stdlib CVE fixed in a Go release, bump the `go`/`toolchain` directive across affected repos.
5. **Agentic repair (nib).** When a bump/cascade/toolchain breaks `go build`, or a PR conflicts beyond a clean rebase, invoke `nib --yolo` to repair the code, verify `go build`, and push only if it builds (else record `needs-human`). nib also synthesizes the cross-repo coordination summary.
6. **Status.** The ledger and dashboard present the full picture per open CVE and per cascade front, including what is blocked / missing / needs-human.

## 3. Non-goals

- Non-Go ecosystems for cascade/toolchain (still reported via collectors).
- Reacting to comments on PRs the bot does not own (unchanged from Plan 3; adopted external PRs are driven via merge-assistance, not comment threads).
- A general dependency-update bot (we act on **security** findings + their cascades, not routine updates).
- Replacing dependabot/renovate — we adopt and drive their PRs, not compete.

## 4. Key decisions (from brainstorm)

| Dimension | Decision |
|---|---|
| Where it lives | Extend `internal/remediate` (planner + ledger + executor); not a new phase |
| Cascade staging | **Pseudo-version immediately** when the upstream fix is available; do NOT wait for a tag or human approval; note in PR body to tag; **re-pin to the tag** as a follow-up once released |
| Existing PRs | **Detect + adopt, never duplicate**; default **nudge** (comment + surface); `--automerge` flag merges green, unblocked addressing PRs (ours / dependabot / renovate) |
| Dependency graph | **Auto-derived** from each tracked repo's `go.mod` (`github.com/kairos-io/*` requires mapped to tracked repos) |
| nib agent role | Build-break **repair**, **conflict resolution**, **toolchain** bump breakage, and the **coordination summary** |
| Autonomy | Proactive, live by default (same dry-run / blast-radius / token model as Plan 2/3); flows without human gates; verify-before-push everywhere |
| Phasing | **4a** status & adoption → **4b** first-party cascade → **4c** nib agent. Spec covers all; plans are sequential. |

## 5. Architecture

All within `internal/remediate`, adding components alongside the existing planner / executor / ledger / reaction layers:

```
correlate (findings + waterfall)  +  collect.prs (open PRs)  +  repos.yaml
                       │
                       ▼
   ┌──────────────────────────────────────────────────────────────────┐
   │ remediate (extended)                                              │
   │                                                                   │
   │  depgraph   ── build first-party module graph from each repo go.mod
   │  matcher    ── link a finding to an existing PR (ours/dependabot/  │
   │               renovate/human) by package/CVE; classify source     │
   │  planner    ── status-aware intents:                              │
   │               • adopt existing PR (link, drive) — no new PR        │
   │               • open direct bump PR (gap: none exists)             │
   │               • cascade: pseudo-version bump consumers of a fix    │
   │               • re-pin: pseudo → released tag                      │
   │               • toolchain: bump go directive for stdlib CVEs       │
   │  executor   ── git+gh actions (open/adjust/cascade/re-pin/merge)   │
   │               + nib agent for repair/conflict/toolchain breakage   │
   │  reactor    ── comment reactions (Plan 3, unchanged)               │
   │  ledger     ── memory: finding↔PR links, source, kind, cascade     │
   │               edges, pin targets, status                          │
   └──────────────────────────────────────────────────────────────────┘
                       │
                       ▼
            ledger.json (committed)  +  dashboard "coordination" view
```

New/extended units, each independently testable:

- **`depgraph`** (4b) — pure: given the tracked repos and their fetched `go.mod`s, produce `consumers(modulePath) → []repo` and `repoModule(repo) → modulePath`. Acyclic (Go modules).
- **`matcher`** (4a) — pure: given a finding (repo, package, CVE) and the open PRs for that repo (from the `prs` collector / live `gh`), decide whether an existing PR addresses it and its `source` (ksec / dependabot / renovate / human), by matching the package name (and bot author).
- **`planner`** (extended) — emits richer intents: `adopt`, `open`, `cascade`, `repin`, `toolchain`, plus the existing `reconcile`. Pure over correlated + ledger + depgraph + matched PRs.
- **`executor`** (extended) — performs the new actions (cascade pseudo-version bump, re-pin to tag, merge/automerge) via git+gh, delegating repair/conflict/toolchain-breakage to the **nib agent**.
- **`nibagent`** (4c) — wraps `nib --yolo` against LocalAI in a repo clone: repair a build break, resolve a conflict, fix toolchain breakage, or write the coordination summary; always verify `go build` before allowing a push.
- **`ledger`** (extended) — see §6.

## 6. Data model (ledger extensions)

`state.LedgerEntry` gains fields (all back-compatible; existing entries default cleanly):

```jsonc
{
  // existing: Key, Repo, Package, Branch, PRNumber, PRURL, State, Bump,
  //           Severity, CreatedRun, LastActionRun, SeenComments, History
  "source": "ksec",        // ksec | dependabot | renovate | human — who owns the PR
  "kind":   "direct",      // direct | cascade | toolchain — what this bump is
  "cascadeFrom": "",       // for kind=cascade: the upstream ledger key that triggered it
  "pinTarget":  "",        // for kind=cascade pseudo-version: the tag to re-pin to once released ("" = still pseudo)
  "pseudo":     false,     // true while the bump points at a pseudo-version awaiting re-pin
  "blocked":    "",        // human-readable reason if progress is stuck (e.g. "checks failing", "needs-human")
  "needsHuman": false      // true when the bot recorded it cannot proceed autonomously
}
```

- An **adopted** entry (`source != "ksec"`) records and drives an external PR; the bot never edits its branch, only merges (with `--automerge`) or nudges.
- A **cascade** entry (`kind=="cascade"`) links to its upstream via `cascadeFrom`, carries `pseudo=true` + `pinTarget` until re-pinned, and can itself be the upstream of further cascade entries (recursion through the DAG).
- The ledger remains the single committed memory; entries are keyed and sorted as today (`"<repo>|<package>"`), with cascade keys distinguished by the module being bumped.

## 7. Engine behavior

### 7.1 Status & adoption (Plan 4a)
Per actionable finding, the `matcher` checks the repo's open PRs:
- **An addressing PR exists** (dependabot/renovate/human/ours): record/refresh a ledger entry linking the finding to that PR with its `source` and live `state`; do **not** open a duplicate. Default action: **nudge** (a single, idempotent comment surfacing the CVE + that it's awaiting merge) and surface on the dashboard. With `--automerge`: if the PR is green (checks pass) and unblocked (no requested-changes review), merge it (or enable GitHub auto-merge).
- **No addressing PR**: this is a gap → open a direct bump PR as today (Plan 2), `source=ksec`.

Matching heuristic: a PR addresses a finding when the PR modifies the same module/package (dependabot/renovate encode the package in the title/branch; the bot matches the package path; author classifies the source). Ambiguous matches are surfaced as `needs-human`, never auto-merged.

### 7.2 First-party cascade (Plan 4b)
`depgraph` maps each tracked repo to the first-party module it provides and the first-party modules it consumes. When a **first-party module's fix is available on its default branch** — the upstream ledger entry reaches `merged` (the bot's fix PR merged, or an adopted PR merged) — the planner enumerates `consumers(module)` and, for each, emits a **cascade** intent: bump that module to the **pseudo-version of the fixed default-branch commit** (`go get <module>@<commit>`), `kind=cascade`, `pseudo=true`, `cascadeFrom=<upstream key>`. The PR body explains it's an unreleased pseudo-version and asks a maintainer to cut a tag. This does **not** wait for a tag.

**Re-pin (follow-up):** each run, for `pseudo` cascade entries, the planner checks whether the module has since published a tag `>=` the fix (`go list -m -versions` / tags). If so, emit a **repin** intent: re-bump the consumer from the pseudo-version to the tag (`go get <module>@<tag>`), set `pseudo=false`, force-push the existing branch. Cascade recurses: a cascade bump in a consumer is itself a fix its own consumers need.

### 7.3 Toolchain bumps (Plan 4c)
For a stdlib finding (`Package=="stdlib"`) whose fix is a Go release, emit a **toolchain** intent: bump the repo's `go`/`toolchain` directive in `go.mod` to the fixed Go version. Deterministic for the directive edit; the nib agent fixes any resulting build breakage.

### 7.4 Agentic repair (Plan 4c, nib)
The executor invokes the **nib agent** (`nib --yolo` against LocalAI, run in the repo clone) at these points, always **verify-before-push**:
- **Build-break repair:** a bump/cascade/toolchain leaves `go build ./...` failing → nib reads the error, edits the consuming code to adapt (e.g. an sdk API changed), re-runs `go build`; push only if it builds, else record `needsHuman` (no broken push).
- **Conflict resolution:** a bot PR conflicts with its base beyond a clean rebase → nib attempts to resolve, verify, push; else `needsHuman`.
- **Coordination summary:** nib synthesizes the cross-repo narrative (open CVEs, cascade fronts, what's blocked / waiting-release / needs-human) for the dashboard/issue — the "depth" the templated summaries lack.

nib is best-effort and gated like the rest: dry-run prints intended agent actions; a nib failure degrades to the deterministic path (`build-failed`/`needs-human`), never a broken push, never the whole run.

## 8. Surfaces / status

The dashboard's "Bot PR ledger" section becomes a **coordination view**: per entry — repo, bump (`pkg@version`, pseudo flagged), `kind` (direct/cascade/toolchain), `source` (ksec/dependabot/renovate/human), PR link + state, and status (open / awaiting-merge / waiting-release / blocked / needs-human). Cascade entries show their `cascadeFrom` so a front is visible end-to-end. The triage narrative is augmented by nib's coordination summary when available (deterministic fallback otherwise). The committed `ledger.json` is the machine-readable status across runs.

## 9. Error handling

- Per-repo / per-intent isolation (as today): a failure records on the entry and continues; never aborts the run.
- **Verify-before-push everywhere** — direct, cascade, re-pin, toolchain, and nib-repaired changes all gate on `go build ./...`; a non-building tree is never pushed (recorded `build-failed`/`needs-human`).
- nib unavailable / fails → deterministic fallback (`build-failed`/`needs-human`); never a broken push.
- Ambiguous PR matches or unresolved conflicts → `needsHuman`, surfaced, never auto-merged.
- Secrets never logged (token redaction from Plan 2 applies to all new git shell-outs).
- Blast-radius cap applies to all new-PR-creating intents (direct + cascade + toolchain) collectively.

## 10. Testing

- **`depgraph`, `matcher`, `planner`** are pure → table-driven + fixture tests (go.mod fixtures, PR-list fixtures, ledger states). The cascade/re-pin/recursion logic is exercised here without git/network.
- **executor / nibagent** are integration (shell git/gh/nib) → verified by `go build`/`go vet`, dry-run plans, and fakes for the orchestration seams (a `FakeAgent`, the existing `ghclient.Fake`).
- **End-to-end dry-run** over fixtures: a fix in a mock sdk cascades to a mock consumer (pseudo), then re-pins on a mock tag — asserted as a plan, no writes.
- Adoption: a dependabot PR fixture is linked (not duplicated) and, under `--automerge`, merged when green.

## 11. Phasing (implementation plans)

- **Plan 4a — Status & adoption.** ledger fields (`source`/`kind`/`blocked`/`needsHuman`), `matcher`, planner `adopt` intent + gap-detection (no duplicates), nudge, `--automerge` flag + green/unblocked merge, dashboard coordination view. *Foundation; high value, lower risk.*
- **Plan 4b — First-party cascade.** `depgraph`, cascade pseudo-version bumps, re-pin to tag, recursion, cascade ledger fields (`cascadeFrom`/`pinTarget`/`pseudo`). *The hard coordination.*
- **Plan 4c — nib agent.** `nibagent` wrapper, build-break repair, conflict resolution, toolchain bumps, coordination summary. *Agentic.*

Each phase produces working, shippable software and is planned + executed independently.

## 12. Operational prerequisites (not code)

- `KSEC_BOT_TOKEN` write scope (`contents:write` + `pull_requests:write`) on the target repos — already required for Plan 2/3 live runs; cascade + adoption + automerge need the same.
- For `--automerge`, the bot's identity must be permitted to merge (or the repos must allow GitHub auto-merge); branch protections still apply.
- nib + LocalAI available on the runner (Plan-1 workflow already provisions LocalAI; nib is re-introduced for 4c as the agent harness, run in repo clones — distinct from the dropped nib-as-triage-client).

## 13. Open items for the plans

- Exact pseudo-version derivation (`go get <module>@<commit>` vs `@<default-branch>`), and how the "fix is available on the default branch" signal is read (upstream entry `merged` vs a direct `git ls-remote` check).
- PR-matching precision (package-path match vs dependabot alert metadata) and the false-match guard threshold before `needsHuman`.
- `--automerge` mechanics (`gh pr merge --auto --squash` vs direct merge) and the green/unblocked checks to require.
- nib invocation contract in CI (model, time budget, working dir) and how its diffs are verified before push.
- Cascade termination/dedup keys so a diamond dependency doesn't double-bump a repo in one run.
