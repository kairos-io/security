# Coordinated Remediation 4c — nib Agent Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an agentic layer to `ksec remediate`: when a bump/cascade/toolchain leaves `go build` failing, an injected agent (`nib --cli --yolo` against LocalAI, run in the repo clone) attempts to **repair** the code; resolve PR **conflicts**; perform Go **toolchain** bumps for stdlib CVEs; and produce a cross-repo **coordination summary** for the dashboard — all best-effort, degrading to today's deterministic behavior when the agent is unavailable.

**Architecture:** A small `Agent` interface (`Repair(dir, task) error`) with a `NibAgent` implementation and a `FakeAgent` for tests. `GitExecutor` gains an `Agent` field; the four `go build ./...` verify sites collapse into one `g.verifyOrRepair(dir, task, runID) bool` helper that, on failure, runs the agent and re-verifies — pushing only if it then builds (else `build-failed`/`needsHuman`, never a broken push). A new `toolchain` intent/executor bumps the `go` directive for stdlib findings. Conflict resolution rebases an owned conflicted PR and uses the agent on conflicts. The coordination summary is generated from the committed ledger via the LocalAI **chat** endpoint (the right tool for text-from-data; nib is reserved for code edits) and shown on the dashboard. Builds on Plans 2/3/4a/4b.

**Tech Stack:** Go 1.22, `nib` + `git` + `go` + `gh` CLIs, LocalAI, `stretchr/testify`. Existing `internal/remediate`, `internal/triage` (OpenAI client pattern), `internal/render`, `internal/config`.

## Global Constraints

- Module `github.com/kairos-io/security`; binary `ksec`; Go 1.22.
- **The agent is best-effort.** If `Agent` is nil or `Repair` fails, behavior is exactly today's: `go build` failure → `build-failed` + `needsHuman`, no push. The agent NEVER lets a non-building tree be pushed — every agent attempt is followed by a fresh `go build ./...` and only a green tree is pushed.
- **nib invocation:** `nib --cli --yolo` (CLI mode, auto-approve) run in the clone dir with the task on stdin, inheriting the process env. nib's model/endpoint config is supplied by the workflow (env / config file pointing at LocalAI) — confirm against nib's docs; treat a misconfigured nib as "agent unavailable" (best-effort).
- Toolchain bumps: a finding with `Package == "stdlib"`, `Ecosystem == "go"`, non-empty `FixedVersion` → bump the repo's `go` directive to that Go version (`go mod edit -go=<ver>`). One toolchain bump per repo per run (highest version). Counts toward `--max-prs`.
- Coordination summary uses the LocalAI chat endpoint (reuse the triage OpenAI-client pattern), NOT nib; best-effort with a deterministic fallback.
- Verify-before-push, dry-run no writes, token redaction, `ksec/` force-push guards — all unchanged and apply to the new toolchain/conflict paths.
- Flags: `--repair` (default true) gates the agent; `--automerge`/`--max-prs`/`--ai-pr-prose` unchanged.

---

## File structure

```
internal/remediate/agent.go          # Agent interface + FakeAgent (create)
internal/remediate/nib_agent.go      # NibAgent (nib --cli --yolo) (create)
internal/remediate/agent_test.go     # (create)
internal/remediate/git_executor.go   # + Agent field + verifyOrRepair; refactor 4 build sites; Conflict rebase (modify)
internal/remediate/repair.go         # repair/conflict/toolchain prompt builders (create)
internal/remediate/repair_test.go    # (create)
internal/remediate/intent.go         # IntentToolchain + Intent.ToolchainVersion (modify)
internal/remediate/planner.go        # toolchain intents (stdlib findings) (modify)
internal/remediate/planner_test.go   # (modify)
internal/remediate/run.go            # Run handles IntentToolchain; Executor.Toolchain (modify)
internal/remediate/fake.go           # FakeExecutor.Toolchain (modify)
internal/remediate/run_test.go       # (modify)
internal/remediate/summary.go        # SummarizeLedger via LocalAI chat (create)
internal/remediate/summary_test.go   # (create)
internal/render/render.go            # show CoordinationSummary (modify)
internal/render/render_test.go       # (modify)
cmd/ksec/main.go                     # wire Agent + --repair; toolchain; summary into render (modify)
.github/workflows/security-dashboard.yaml  # install nib; nib LocalAI config (modify)
```

---

### Task 1: `Agent` interface + `FakeAgent` + `NibAgent`

**Files:** Create `internal/remediate/agent.go`, `internal/remediate/nib_agent.go`, `internal/remediate/agent_test.go`.

**Interfaces:**
- `type Agent interface { Repair(dir, task string) error }`
- `type FakeAgent struct { Calls []string; Err error; Edit func(dir string) }` implementing it (records the task; optional `Edit` callback simulates file edits; returns `Err`).
- `type NibAgent struct { cfg config.AIConfig; run func(dir, task string) error }`; `func NewNibAgent(cfg config.AIConfig) *NibAgent` — `run` shells `nib --cli --yolo` in `dir` with `task` on stdin.

- [ ] **Step 1: Write the failing test**

Create `internal/remediate/agent_test.go`:

```go
package remediate

import (
	"errors"
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakeAgentRecordsAndErrors(t *testing.T) {
	f := &FakeAgent{}
	require.NoError(t, f.Repair("/tmp/x", "fix the build"))
	assert.Equal(t, []string{"fix the build"}, f.Calls)

	f2 := &FakeAgent{Err: errors.New("nope")}
	assert.Error(t, f2.Repair("/tmp/x", "t"))
}

func TestNewNibAgent(t *testing.T) {
	a := NewNibAgent(config.AIConfig{})
	require.NotNil(t, a)
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement**

Create `internal/remediate/agent.go`:

```go
package remediate

// Agent performs an in-repo code edit task (e.g. fix a build break) in the
// working directory dir. Returning nil means the agent ran; the caller MUST
// re-verify the build before trusting the result.
type Agent interface {
	Repair(dir, task string) error
}

// FakeAgent is a test double.
type FakeAgent struct {
	Calls []string
	Err   error
	Edit  func(dir string) // optional: simulate file edits
}

func (f *FakeAgent) Repair(dir, task string) error {
	f.Calls = append(f.Calls, task)
	if f.Edit != nil {
		f.Edit(dir)
	}
	return f.Err
}
```

Create `internal/remediate/nib_agent.go`:

```go
package remediate

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/kairos-io/security/internal/config"
)

// NibAgent runs `nib --cli --yolo` in a repo clone to perform a code-edit task.
// nib reads its model/endpoint from the environment/config provided by the
// workflow; a misconfigured or missing nib makes Repair return an error, which
// the caller treats as "no repair" (best-effort).
type NibAgent struct {
	cfg config.AIConfig
	run func(dir, task string) error
}

func NewNibAgent(cfg config.AIConfig) *NibAgent {
	return &NibAgent{cfg: cfg, run: func(dir, task string) error {
		cmd := exec.Command("nib", "--cli", "--yolo")
		cmd.Dir = dir
		cmd.Stdin = bytes.NewBufferString(task)
		var errb bytes.Buffer
		cmd.Stderr = &errb
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("nib: %v: %s", err, errb.String())
		}
		return nil
	}}
}

func (a *NibAgent) Repair(dir, task string) error { return a.run(dir, task) }
```

- [ ] **Step 4: Run + commit** — `go test ./internal/remediate/...` PASS.
```bash
git add internal/remediate/agent.go internal/remediate/nib_agent.go internal/remediate/agent_test.go
git commit -m "feat(remediate): Agent interface + NibAgent + FakeAgent"
```

---

### Task 2: Build-break repair helper + refactor the 4 verify sites

**Files:** Create `internal/remediate/repair.go`, `internal/remediate/repair_test.go`; modify `internal/remediate/git_executor.go`.

**Interfaces:**
- `func RepairTask(buildErr string) string` (pure) — the prompt: explain the `go build ./...` failure and ask the agent to fix the code so it compiles, changing as little as possible, without altering the dependency versions.
- `GitExecutor` gains `Agent Agent` field and `func (g *GitExecutor) verifyOrRepair(dir, task, runID string) bool` — runs `go build ./...`; on failure, if `g.Agent != nil` runs `g.Agent.Repair(dir, RepairTask(buildErr))` and re-runs `go build ./...`; returns whether the tree now builds. The four existing `go build ./...` blocks (Open:83, Adjust:235, Cascade:295, Repin:362) are refactored to call it.

- [ ] **Step 1: Write the failing test (the pure prompt builder)**

Create `internal/remediate/repair_test.go`:

```go
package remediate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepairTask(t *testing.T) {
	task := RepairTask("./foo.go:10: undefined: Bar")
	assert.Contains(t, task, "undefined: Bar")
	assert.Contains(t, strings.ToLower(task), "go build")
	assert.Contains(t, strings.ToLower(task), "compile")
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement**

Create `internal/remediate/repair.go`:

```go
package remediate

import "fmt"

// RepairTask is the prompt handed to the agent when `go build` fails after a
// dependency bump. It must fix the code to compile WITHOUT changing the
// dependency versions (the security bump is deterministic and must stand).
func RepairTask(buildErr string) string {
	return fmt.Sprintf(`A dependency security bump in this repository broke the build.
Make the minimal source-code changes needed so that `+"`go build ./...`"+` compiles again.
Do NOT change any dependency versions in go.mod/go.sum — only adapt the calling code
(e.g. to a changed API). The compiler error was:

%s`, buildErr)
}
```

In `internal/remediate/git_executor.go`: add `Agent Agent` to the struct, add the helper, and replace each of the four `if _, err := g.run(dir, "go", "build", "./..."); err != nil { ...build-failed... }` blocks with a call to it. The helper:

```go
// verifyOrRepair runs `go build ./...`; on failure it asks the agent (if any)
// to repair the code and re-verifies. Returns true iff the tree builds (which
// the caller requires before pushing).
func (g *GitExecutor) verifyOrRepair(dir, task, runID string) bool {
	if _, err := g.run(dir, "go", "build", "./..."); err == nil {
		return true
	} else if g.Agent == nil {
		return false
	} else {
		_ = g.Agent.Repair(dir, RepairTask(err.Error()))
	}
	_, err := g.run(dir, "go", "build", "./...")
	return err == nil
}
```

Refactor pattern (apply at each site; example for Open):

```go
	if !g.verifyOrRepair(dir, "open "+in.Package, runID) {
		entry.State = "build-failed"
		entry.NeedsHuman = true
		entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "build-failed"})
		return entry, nil
	}
```

(Keep each site's existing `entry`/`e` variable and history-event action name; only the build+repair check changes. Set `NeedsHuman = true` consistently on the failed path at all four sites.)

- [ ] **Step 4: Run + build + commit** — `go test ./internal/remediate/...`, `go build ./... && go vet ./...` (PASS).
```bash
git add internal/remediate/repair.go internal/remediate/repair_test.go internal/remediate/git_executor.go
git commit -m "feat(remediate): agent build-break repair (verifyOrRepair) at all bump sites"
```

---

### Task 3: Toolchain intent (planner)

**Files:** Modify `internal/remediate/intent.go`, `internal/remediate/planner.go`, `internal/remediate/planner_test.go`.

**Interfaces:**
- `IntentToolchain IntentType = "toolchain"`; `Intent` gains `ToolchainVersion string`.
- `Plan` emits, for each repo with a stdlib finding (`Package=="stdlib"`, `Ecosystem=="go"`, `FixedVersion != ""`), one `IntentToolchain{Repo, ToolchainVersion: <highest fixed Go version, leading "go" stripped>, Severity}` — deduped per repo, added to the shared new-PR `pool` (capped). Skip if a ledger entry keyed `"<repo>|go-toolchain"` is already open/conflicted/merged/closed.

- [ ] **Step 1: Update intent + test**

In `intent.go`: add `IntentToolchain IntentType = "toolchain"` and `ToolchainVersion string` to `Intent`. In `planner_test.go` add:

```go
func TestPlanToolchainForStdlib(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "s", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go",
			Package: "stdlib", FixedVersion: "go1.22.5", Severity: "high"},
	}}
	intents, _ := Plan(c, state.Ledger{}, nil, nil, 10)
	var tc *Intent
	for i := range intents {
		if intents[i].Type == IntentToolchain {
			tc = &intents[i]
		}
	}
	require.NotNil(t, tc)
	assert.Equal(t, "kairos-io/immucore", tc.Repo)
	assert.Equal(t, "1.22.5", tc.ToolchainVersion) // leading "go" stripped
	assert.Equal(t, "kairos-io/immucore|go-toolchain", tc.Key)
}
```

- [ ] **Step 2: Run red.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement** — in `planner.go`, add a toolchain collection alongside the targets loop and emit into the `pool`. Add helper `func toolchainKey(repo string) string { return repo + "|go-toolchain" }`. Collect the highest `FixedVersion` (strip a leading "go") per repo over stdlib findings; for each, if no covering ledger entry (open/conflicted/merged/closed) for `toolchainKey(repo)`, append `newPR{intent: Intent{Type: IntentToolchain, Key: toolchainKey(repo), Repo: repo, ToolchainVersion: ver, Severity: sev}, sev: sev}` to the `pool` (before the sort/cap). `actionable()` already excludes stdlib from the direct-bump path (its package is "stdlib", not a module path you can `go get`), so there's no overlap. (If `actionable()` currently would treat stdlib as a normal target, add `&& f.Package != "stdlib"` to it.)

- [ ] **Step 4: Run + commit** — PASS.
```bash
git add internal/remediate/intent.go internal/remediate/planner.go internal/remediate/planner_test.go
git commit -m "feat(remediate): toolchain bump intents for stdlib findings"
```

---

### Task 4: Run loop + Executor.Toolchain + fake + GitExecutor.Toolchain

**Files:** Modify `internal/remediate/run.go`, `internal/remediate/fake.go`, `internal/remediate/run_test.go`, `internal/remediate/git_executor.go`.

**Interfaces:** `Executor` gains `Toolchain(in Intent, run string) (state.LedgerEntry, error)`; `Run` handles `IntentToolchain` (mirrors `IntentOpen`, error-isolated, error entry `Kind:"toolchain"`); `FakeExecutor.Toolchained` map + method; real `GitExecutor.Toolchain` (clone, branch `ksec/toolchain-<ver>`, `go mod edit -go=<ver>`, `go mod tidy`, `verifyOrRepair`, push, PR; entry `Kind:"toolchain"`; dry-run prints).

- [ ] **Step 1: Test** — add `TestRunToolchain` to `run_test.go` (an `IntentToolchain` → `FakeExecutor.Toolchain` returns an entry with `Kind:"toolchain"`; one Result action "toolchain").

- [ ] **Step 2: Run red.** `go test ./internal/remediate/...`

- [ ] **Step 3: Implement** — `Executor.Toolchain` interface method; `Run` `case IntentToolchain` (mirror `IntentOpen`; error entry sets `Kind:"toolchain"`, `ToolchainVersion` info in history); `FakeExecutor.Toolchained` + method; and the real `GitExecutor.Toolchain`:

```go
func (g *GitExecutor) Toolchain(in Intent, runID string) (state.LedgerEntry, error) {
	branch := "ksec/toolchain-" + slug(in.ToolchainVersion)
	entry := state.LedgerEntry{
		Key: in.Key, Repo: in.Repo, Package: "go-toolchain", Branch: branch, Kind: "toolchain",
		Severity: in.Severity, CreatedRun: runID, LastActionRun: runID,
		Bump: state.Bump{Package: "go", To: in.ToolchainVersion},
	}
	if g.DryRun {
		fmt.Printf("[dry-run] would bump go toolchain in %s to %s\n", in.Repo, in.ToolchainVersion)
		entry.State = "planned"
		return entry, nil
	}
	dir, err := os.MkdirTemp("", "ksec-tc-*")
	if err != nil {
		return entry, err
	}
	defer os.RemoveAll(dir)
	if _, err := g.run("", "git", "clone", "--depth", "1", g.cloneURL(in.Repo), dir); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "git", "checkout", "-b", branch); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "go", "mod", "edit", "-go="+in.ToolchainVersion); err != nil {
		return entry, err
	}
	_, _ = g.run(dir, "go", "mod", "tidy")
	if !g.verifyOrRepair(dir, "go toolchain bump to "+in.ToolchainVersion, runID) {
		entry.State = "build-failed"
		entry.NeedsHuman = true
		entry.History = []state.LedgerEvent{{Run: runID, Action: "toolchain-build-failed"}}
		return entry, nil
	}
	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	if _, err := g.run(dir, "git", "commit", "-am", "chore(security): bump go toolchain to "+in.ToolchainVersion); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "git", "push", "-u", "origin", branch); err != nil {
		return entry, err
	}
	out, err := g.run(dir, "gh", "pr", "create", "-R", in.Repo, "--head", branch,
		"--title", "chore(security): bump go toolchain to "+in.ToolchainVersion,
		"--body", "Bumps the Go toolchain to "+in.ToolchainVersion+" to address a stdlib vulnerability. "+PRMarker(in.Key))
	if err != nil {
		return entry, err
	}
	entry.PRURL = strings.TrimSpace(string(out))
	entry.PRNumber = prNumberFromURL(entry.PRURL)
	entry.State = "open"
	entry.History = []state.LedgerEvent{{Run: runID, Action: "toolchain-opened", Detail: entry.PRURL}}
	return entry, nil
}
```

- [ ] **Step 4: Run + build + commit** — `go test ./internal/remediate/...`, `go build ./... && go vet ./...` PASS.
```bash
git add internal/remediate/run.go internal/remediate/fake.go internal/remediate/run_test.go internal/remediate/git_executor.go
git commit -m "feat(remediate): Run + GitExecutor.Toolchain (go directive bump, agent-repaired)"
```

---

### Task 5: Conflict resolution (rebase + agent)

**Files:** Modify `internal/remediate/git_executor.go`.

**Interfaces:** `func (g *GitExecutor) ResolveConflict(e state.LedgerEntry, runID string) (state.LedgerEntry, error)` — for an owned (`ksec/` branch) conflicted PR: clone, checkout `e.Branch`, `git rebase origin/<defaultBranch>`; if the rebase reports conflicts, run the agent with `ConflictTask()` then `git rebase --continue`; verify `go build` (via `verifyOrRepair`), force-push; set `State:"open"` (cleared conflict) or `needsHuman` if unresolved. Called from `Reconcile` when it detects a conflicted PR. Dry-run prints; `ksec/` guard; token redacted.

- [ ] **Step 1: Add `ConflictTask` (pure)** in `repair.go`:

```go
func ConflictTask() string {
	return "This branch has git merge conflicts after rebasing onto the base branch. " +
		"Resolve every conflict marker so the code is correct and `go build ./...` compiles, " +
		"keeping the dependency-version change from this branch."
}
```

Add a test in `repair_test.go` asserting `ConflictTask()` mentions "conflict".

- [ ] **Step 2: Implement `ResolveConflict`** and call it from `Reconcile` when the reconciled PR is conflicted (extend `Reconcile`'s `PRStatusOf`/view to detect `mergeable==false` + a conflict state, then `return g.ResolveConflict(e, runID)`). Guard: only own `ksec/` branches; dry-run prints + returns unchanged. Build the rebase+agent+verify+force-push flow with `g.run` (token-redacting), mirroring `Adjust`/`Repin` push discipline. On unresolved conflict or build failure → `e.NeedsHuman = true`, `e.Blocked = "conflict"`, no push.

- [ ] **Step 3: Build + vet + test + commit** — integration; `go build ./... && go vet ./... && go test ./...` PASS.
```bash
git add internal/remediate/git_executor.go internal/remediate/repair.go internal/remediate/repair_test.go
git commit -m "feat(remediate): agent-assisted conflict resolution on owned PRs"
```

---

### Task 6: Coordination summary (LocalAI chat)

**Files:** Create `internal/remediate/summary.go`, `internal/remediate/summary_test.go`; modify `internal/render/render.go`, `internal/render/render_test.go`.

**Interfaces:**
- `func SummarizeLedger(cfg config.AIConfig, led state.Ledger) (string, error)` — posts a compact view of the ledger (per entry: repo, package, kind, source, state, pseudo) to `<endpoint>/v1/chat/completions` and returns a 2-4 sentence coordination narrative ("what's open, cascading, blocked, needs-human"). Empty endpoint → error (caller falls back). Reuse the `internal/triage/openai.go` HTTP shape.
- `render.Input` gains `CoordinationSummary string`; `DashboardMarkdown` renders it under a `## 🧭 Coordination` heading (when non-empty).

- [ ] **Step 1: Write the failing tests** — `summary_test.go` (httptest: a ledger → POST to `/v1/chat/completions` → returns content; assert the returned string; empty-endpoint → error) and a `render_test.go` assertion that `DashboardMarkdown` shows the summary under "Coordination" when `Input.CoordinationSummary != ""`.

- [ ] **Step 2: Run red.** `go test ./internal/remediate/... ./internal/render/...`

- [ ] **Step 3: Implement** — `SummarizeLedger` (mirror `internal/triage/openai.go`'s chat POST; a compact prompt over the ledger entries; return trimmed content); add `CoordinationSummary` to `render.Input` and a section in `DashboardMarkdown` (and `html.go`, escaped; regenerate goldens). Deterministic: same input → same output (the render side is pure; the AI call is in the command).

- [ ] **Step 4: Run + commit** — PASS (regenerate render goldens if needed).
```bash
git add internal/remediate/summary.go internal/remediate/summary_test.go internal/render/render.go internal/render/render_test.go internal/render/html.go internal/render/testdata/
git commit -m "feat(remediate): coordination summary via LocalAI; dashboard section"
```

---

### Task 7: Wire the agent + summary + workflow

**Files:** Modify `cmd/ksec/main.go`, `.github/workflows/security-dashboard.yaml`.

**Interfaces:** the `remediate` command sets `ex.Agent` (a `NibAgent` when `--repair` and an AI endpoint are configured) and the planner already emits toolchain intents; the `render` command best-effort-calls `SummarizeLedger` and sets `Input.CoordinationSummary`. The workflow re-installs `nib` and points it at LocalAI.

- [ ] **Step 1: Wire the agent in `newRemediateCmd`** — add `var repair bool` + `cmd.Flags().BoolVar(&repair, "repair", true, "use the nib agent to repair build breaks / conflicts")`; after constructing `ex` (which already has `GH`/`Automerge`/`Prose`), set:

```go
			if repair && aiCfg.Nib.Endpoint != "" {
				ex.Agent = remediate.NewNibAgent(aiCfg)
			}
```

- [ ] **Step 2: Wire the summary in `newRenderCmd`** — after loading the ledger, best-effort:

```go
			summary := ""
			if aiCfg, err := config.LoadAI("ai.yaml"); err == nil {
				if s, err := remediate.SummarizeLedger(aiCfg, ledger); err == nil {
					summary = s
				}
			}
```

and set `CoordinationSummary: summary` on the `render.Input`. (Add the `remediate` import to main.go if not present.)

- [ ] **Step 3: Workflow** — re-add a nib install step and its LocalAI config. After the LocalAI start step, add:

```yaml
      - name: Install nib
        run: go install github.com/mudler/nib@${NIB_VERSION:-latest}
```

and export the env nib needs to reach LocalAI (confirm the exact variable names against nib's docs; this is the integration point):

```yaml
      # in the job env: block
      NIB_MODEL: ${{ vars.LOCALAI_MODEL }}     # or read from ai.yaml at runtime
      OPENAI_API_BASE: http://localhost:8080/v1
      OPENAI_API_KEY: sk-localai
```

- [ ] **Step 4: Build + vet + full test + YAML validate + smoke**

Run: `go build ./... && go vet ./... && go test ./...` and `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/security-dashboard.yaml'))" && echo OK`.
Smoke: `go run ./cmd/ksec remediate --help` shows `--repair`.

- [ ] **Step 5: Commit**

```bash
git add cmd/ksec/main.go .github/workflows/security-dashboard.yaml
git commit -m "feat(remediate): wire nib agent (--repair) + coordination summary; install nib in CI"
```

---

## Self-review

**Spec coverage** (§7.3 toolchain + §7.4 nib agent):
- Build-break repair (nib, verify-before-push, else build-failed/needsHuman) → Tasks 1, 2. ✓
- Conflict resolution (rebase + agent on owned PRs) → Task 5. ✓
- Toolchain bumps for stdlib CVEs (go directive + agent-repaired) → Tasks 3, 4. ✓
- Coordination summary (cross-repo narrative on the dashboard) → Task 6 (via LocalAI chat — nib is for edits, chat for text-from-data; documented deviation). ✓
- Best-effort agent, degrades to deterministic behavior; dry-run/verify/token/ksec-guard preserved → Tasks 1, 2, 4, 5. ✓
- nib re-installed + configured in CI → Task 7. ✓
- Flags: `--repair` default true → Task 7. ✓

**Placeholder scan:** the nib env-var names in Task 7 Step 3 are flagged as an integration point to confirm against nib's docs (not a code placeholder); everything else is complete code/commands.

**Type consistency:** `Agent`/`FakeAgent`/`NibAgent` (Task 1) used by Tasks 2, 7. `RepairTask`/`ConflictTask` (Tasks 2, 5) used by `verifyOrRepair`/`ResolveConflict`. `GitExecutor.Agent` + `verifyOrRepair` (Task 2) used by Tasks 4, 5. `IntentToolchain`/`Intent.ToolchainVersion` (Task 3) used by Tasks 4, 7. `Executor.Toolchain` (Task 4) implemented by `FakeExecutor` + `GitExecutor`. `SummarizeLedger` (Task 6) + `render.Input.CoordinationSummary` used by Task 7.

---

## Operational notes / known risks

- **nib model wiring is the integration risk.** nib has no `-model`/`-endpoint` flags (we learned this earlier); it reads config from env or a config file. Task 7 sets `OPENAI_API_BASE`/`OPENAI_API_KEY`/`NIB_MODEL` as the likely mechanism — confirm against nib's docs and adjust. If nib can't reach the model, `Repair` errors and the system records `build-failed`/`needsHuman` exactly as today (no regression).
- **Small-model repair quality is uncertain.** A small LocalAI model may not fix non-trivial API breaks; the verify-before-push gate guarantees it can only ever push a building tree, so a failed repair is safe (recorded for a human).
- **Toolchain version mapping:** stdlib `FixedVersion` from govulncheck is normalized by stripping a leading "go"; confirm the format in a live run and adjust if govulncheck reports a bare `1.22.5`.
- The coordination summary is best-effort; on AI failure the dashboard simply omits the section.
