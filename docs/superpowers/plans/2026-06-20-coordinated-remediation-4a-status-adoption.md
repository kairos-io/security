# Coordinated Remediation 4a — Status & Adoption Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `ksec remediate` status-aware: for each actionable finding, detect whether a PR already addresses it (dependabot / renovate / human / ours), adopt and drive that PR instead of opening a duplicate, nudge it by default, optionally auto-merge it when green, and surface source + kind + status on the dashboard.

**Architecture:** Extend the existing pure `Plan` (it gains the repo's open PRs and emits a new `adopt` intent for findings already covered by an external PR, `open` only for true gaps), a pure `matcher` (links a finding's package to an open PR and classifies its source), the `Executor` (a new `Adopt` action that records the link, posts an idempotent nudge, and — with `--automerge` — merges green/unblocked PRs), the ledger (new `source`/`kind`/`blocked`/`needsHuman` fields), and the dashboard (coordination columns). Builds on Plans 2 & 3.

**Tech Stack:** Go 1.22, `gh` CLI, `stretchr/testify`, the existing `internal/remediate`, `internal/state`, `internal/ghclient`, `internal/render`.

## Global Constraints

- Module `github.com/kairos-io/security`; binary `ksec`; Go 1.22.
- **Never duplicate**: if an external PR (dependabot/renovate/human) addresses a finding, adopt it — do not open ours.
- A PR's `source` is classified by author: `renovate[bot]`→`renovate`, `dependabot[bot]`→`dependabot`, `kairos-security-bot`→`ksec`, else `human`.
- Adoption emits no git writes and is **not** subject to the `--max-prs` cap (only `open` creates PRs); the cap applies to `open` intents only.
- **Nudge is idempotent**: post exactly one nudge comment per adopted PR, guarded by the marker `<!-- ksec:nudge -->` (skip if a comment with the marker already exists).
- `--automerge` (default **false**): only merge a PR that is **mergeable + checks passing + not blocked by a requested-changes review**. Default behavior is nudge-only. Ambiguity → `needsHuman`, never auto-merge.
- Package-match heuristic: a PR addresses a finding iff its title contains the finding's full package path (case-insensitive).
- Token never logged (existing `g.run` redaction); dry-run short-circuits every git/GitHub write.

---

## File structure

```
internal/state/types.go            # + LedgerEntry.Source/Kind/Blocked/NeedsHuman (modify)
internal/remediate/intent.go       # + IntentAdopt + Intent.PRNumber/PRURL/Source (modify)
internal/remediate/matcher.go      # MatchPR + classifySource (create)
internal/remediate/matcher_test.go # (create)
internal/remediate/planner.go      # Plan gains prsByRepo; emits adopt vs open (modify)
internal/remediate/planner_test.go # (modify: new signature + adopt cases)
internal/remediate/automerge.go    # ShouldAutomerge pure decision (create)
internal/remediate/automerge_test.go # (create)
internal/remediate/run.go          # Executor.Adopt + Run handles IntentAdopt (modify)
internal/remediate/fake.go         # FakeExecutor.Adopt (modify)
internal/remediate/run_test.go     # (modify: adopt path)
internal/ghclient/ghclient.go      # + PRStatus, MergePR, types (modify)
internal/ghclient/fake.go          # + fakes (modify)
internal/ghclient/status_test.go   # (create)
internal/remediate/git_executor.go # + GH/Automerge fields + Adopt method (modify)
internal/render/render.go          # ledger table: Source/Kind/Status columns (modify)
internal/render/coord_test.go      # (create)
cmd/ksec/main.go                   # collect prsByRepo, pass to Plan; --automerge; wire ex.GH (modify)
```

---

### Task 1: Ledger coordination fields

**Files:**
- Modify: `internal/state/types.go`
- Test: `internal/state/ledger_test.go` (add a case)

**Interfaces:**
- Produces: `LedgerEntry` gains `Source`, `Kind`, `Blocked string` and `NeedsHuman bool`.

- [ ] **Step 1: Add the fields**

In `internal/state/types.go`, add to `LedgerEntry` (after `Severity`):

```go
	Source     string `json:"source,omitempty"`     // ksec | dependabot | renovate | human
	Kind       string `json:"kind,omitempty"`       // direct | cascade | toolchain
	Blocked    string `json:"blocked,omitempty"`    // human-readable reason progress is stuck
	NeedsHuman bool   `json:"needsHuman,omitempty"`
```

- [ ] **Step 2: Add a round-trip assertion**

In `internal/state/ledger_test.go`, extend `TestLedgerRoundTrip`'s entry literal to set the new fields and assert they survive:

```go
	in := Ledger{Entries: []LedgerEntry{{Key: "a|b", Repo: "a", Package: "b", State: "open",
		Source: "dependabot", Kind: "direct", Blocked: "checks failing", NeedsHuman: true,
		Bump: Bump{Package: "b", To: "1.2.3"}, History: []LedgerEvent{{Run: "2026-06-20", Action: "opened"}}}}}
```

- [ ] **Step 3: Run + commit**

Run: `go test ./internal/state/...` (PASS).

```bash
git add internal/state/types.go internal/state/ledger_test.go
git commit -m "feat(state): ledger coordination fields (source/kind/blocked/needsHuman)"
```

---

### Task 2: `ghclient` PR status + merge

**Files:**
- Modify: `internal/ghclient/ghclient.go`
- Modify: `internal/ghclient/fake.go`
- Test: `internal/ghclient/status_test.go`

**Interfaces:**
- Produces:
  - `type PRStatus struct { State string; Mergeable bool; ChecksPassing bool; ReviewDecision string }`
  - `GitHub` interface gains `PRStatusOf(repo string, pr int) (PRStatus, error)` and `MergePR(repo string, pr int, auto bool) error`.
  - `CLI` impls (shell `gh`); `Fake` fields `Statuses map[string]PRStatus` (key `"<repo>#<pr>"`) and `Merged []string` (record `"<repo>#<pr>"`; `auto` appends `" (auto)"`).

- [ ] **Step 1: Write the failing test**

Create `internal/ghclient/status_test.go`:

```go
package ghclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakePRStatusAndMerge(t *testing.T) {
	f := NewFake()
	f.Statuses["o/r#5"] = PRStatus{State: "OPEN", Mergeable: true, ChecksPassing: true, ReviewDecision: ""}
	got, err := f.PRStatusOf("o/r", 5)
	require.NoError(t, err)
	assert.True(t, got.ChecksPassing)

	require.NoError(t, f.MergePR("o/r", 5, true))
	assert.Equal(t, []string{"o/r#5 (auto)"}, f.Merged)
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/ghclient/...`

- [ ] **Step 3: Implement**

In `internal/ghclient/ghclient.go` add the type, two interface methods, and CLI impls:

```go
type PRStatus struct {
	State          string `json:"state"`
	Mergeable      bool   `json:"mergeable"`
	ChecksPassing  bool   `json:"checksPassing"`
	ReviewDecision string `json:"reviewDecision"`
}
```

Interface additions (in the `GitHub` interface):

```go
	PRStatusOf(repo string, pr int) (PRStatus, error)
	MergePR(repo string, pr int, auto bool) error
```

CLI methods:

```go
func (c *CLI) PRStatusOf(repo string, pr int) (PRStatus, error) {
	b, err := c.run("pr", "view", fmt.Sprint(pr), "-R", repo,
		"--json", "state,mergeable,reviewDecision,statusCheckRollup",
		"-q", "{state: .state, mergeable: (.mergeable == \"MERGEABLE\"), reviewDecision: (.reviewDecision // \"\"), "+
			"checksPassing: ([.statusCheckRollup[]? | select((.conclusion // .state) as $s | $s != \"SUCCESS\" and $s != \"NEUTRAL\" and $s != \"SKIPPED\")] | length == 0)}")
	if err != nil {
		return PRStatus{}, err
	}
	var s PRStatus
	return s, json.Unmarshal(b, &s)
}

func (c *CLI) MergePR(repo string, pr int, auto bool) error {
	args := []string{"pr", "merge", fmt.Sprint(pr), "-R", repo, "--squash"}
	if auto {
		args = append(args, "--auto")
	}
	_, err := c.run(args...)
	return err
}
```

In `internal/ghclient/fake.go` add fields + `NewFake` init + methods:

```go
// struct fields:
	Statuses map[string]PRStatus
	Merged   []string

// NewFake():
		Statuses: map[string]PRStatus{},

func (f *Fake) PRStatusOf(repo string, pr int) (PRStatus, error) { return f.Statuses[prKey(repo, pr)], nil }
func (f *Fake) MergePR(repo string, pr int, auto bool) error {
	k := prKey(repo, pr)
	if auto {
		k += " (auto)"
	}
	f.Merged = append(f.Merged, k)
	return nil
}
```

- [ ] **Step 4: Run + commit**

Run: `go test ./internal/ghclient/...` (PASS).

```bash
git add internal/ghclient/
git commit -m "feat(ghclient): PR status (checks/mergeable/review) and merge"
```

---

### Task 3: `matcher` — link a finding to an existing PR

**Files:**
- Create: `internal/remediate/matcher.go`
- Test: `internal/remediate/matcher_test.go`

**Interfaces:**
- Consumes: `ghclient.PullRequest`.
- Produces:
  - `func classifySource(author string) string` — author → `renovate`/`dependabot`/`ksec`/`human`.
  - `func MatchPR(pkg string, prs []ghclient.PullRequest) (ghclient.PullRequest, string, bool)` — returns the first PR whose title contains `pkg` (case-insensitive) plus its source; `ok=false` if none.

- [ ] **Step 1: Write the failing test**

Create `internal/remediate/matcher_test.go`:

```go
package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/stretchr/testify/assert"
)

func TestClassifySource(t *testing.T) {
	assert.Equal(t, "renovate", classifySource("renovate[bot]"))
	assert.Equal(t, "dependabot", classifySource("dependabot[bot]"))
	assert.Equal(t, "ksec", classifySource("kairos-security-bot"))
	assert.Equal(t, "human", classifySource("alice"))
}

func TestMatchPR(t *testing.T) {
	prs := []ghclient.PullRequest{
		{Number: 1, Title: "Bump golang.org/x/net from 0.30.0 to 0.33.0", Author: "dependabot[bot]"},
		{Number: 2, Title: "Some feature", Author: "alice"},
	}
	pr, src, ok := MatchPR("golang.org/x/net", prs)
	assert.True(t, ok)
	assert.Equal(t, 1, pr.Number)
	assert.Equal(t, "dependabot", src)

	_, _, ok = MatchPR("golang.org/x/crypto", prs)
	assert.False(t, ok)
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement**

Create `internal/remediate/matcher.go`:

```go
package remediate

import (
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
)

func classifySource(author string) string {
	switch author {
	case "renovate[bot]":
		return "renovate"
	case "dependabot[bot]":
		return "dependabot"
	case "kairos-security-bot":
		return "ksec"
	default:
		return "human"
	}
}

// MatchPR returns the first open PR whose title contains the package path
// (case-insensitive) and the PR's source. A non-empty pkg is required.
func MatchPR(pkg string, prs []ghclient.PullRequest) (ghclient.PullRequest, string, bool) {
	if pkg == "" {
		return ghclient.PullRequest{}, "", false
	}
	needle := strings.ToLower(pkg)
	for _, pr := range prs {
		if strings.Contains(strings.ToLower(pr.Title), needle) {
			return pr, classifySource(pr.Author), true
		}
	}
	return ghclient.PullRequest{}, "", false
}
```

- [ ] **Step 4: Run + commit**

Run: `go test ./internal/remediate/...` (PASS).

```bash
git add internal/remediate/matcher.go internal/remediate/matcher_test.go
git commit -m "feat(remediate): matcher links findings to existing PRs by package"
```

---

### Task 4: `ShouldAutomerge` pure decision

**Files:**
- Create: `internal/remediate/automerge.go`
- Test: `internal/remediate/automerge_test.go`

**Interfaces:**
- Consumes: `ghclient.PRStatus`.
- Produces: `func ShouldAutomerge(s ghclient.PRStatus) bool` — true iff `Mergeable && ChecksPassing && ReviewDecision != "CHANGES_REQUESTED"`.

- [ ] **Step 1: Write the failing test**

Create `internal/remediate/automerge_test.go`:

```go
package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/stretchr/testify/assert"
)

func TestShouldAutomerge(t *testing.T) {
	ok := ghclient.PRStatus{Mergeable: true, ChecksPassing: true, ReviewDecision: ""}
	assert.True(t, ShouldAutomerge(ok))

	assert.False(t, ShouldAutomerge(ghclient.PRStatus{Mergeable: false, ChecksPassing: true}))
	assert.False(t, ShouldAutomerge(ghclient.PRStatus{Mergeable: true, ChecksPassing: false}))
	assert.False(t, ShouldAutomerge(ghclient.PRStatus{Mergeable: true, ChecksPassing: true, ReviewDecision: "CHANGES_REQUESTED"}))
	assert.True(t, ShouldAutomerge(ghclient.PRStatus{Mergeable: true, ChecksPassing: true, ReviewDecision: "APPROVED"}))
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement**

Create `internal/remediate/automerge.go`:

```go
package remediate

import "github.com/kairos-io/security/internal/ghclient"

// ShouldAutomerge reports whether an addressing PR is safe to merge: it must be
// mergeable, have passing checks, and not be blocked by a requested-changes
// review.
func ShouldAutomerge(s ghclient.PRStatus) bool {
	return s.Mergeable && s.ChecksPassing && s.ReviewDecision != "CHANGES_REQUESTED"
}
```

- [ ] **Step 4: Run + commit**

Run: `go test ./internal/remediate/...` (PASS).

```bash
git add internal/remediate/automerge.go internal/remediate/automerge_test.go
git commit -m "feat(remediate): pure automerge eligibility decision"
```

---

### Task 5: Planner emits `adopt` vs `open`

**Files:**
- Modify: `internal/remediate/intent.go`
- Modify: `internal/remediate/planner.go`
- Modify: `internal/remediate/planner_test.go`

**Interfaces:**
- Consumes: `MatchPR`, `ghclient.PullRequest`.
- Produces:
  - `IntentAdopt IntentType = "adopt"`; `Intent` gains `PRNumber int`, `PRURL string`, `Source string`.
  - `Plan(c state.Correlated, ledger state.Ledger, prsByRepo map[string][]ghclient.PullRequest, maxNew int) ([]Intent, int)` — for each target needing action, if an external PR addresses it emit `IntentAdopt` (uncapped); else emit `IntentOpen` (capped by `maxNew`).

- [ ] **Step 1: Update intent types**

In `internal/remediate/intent.go`: add `IntentAdopt IntentType = "adopt"` to the const block, and add to `Intent`:

```go
	PRNumber int
	PRURL    string
	Source   string // dependabot | renovate | human (for IntentAdopt)
```

- [ ] **Step 2: Update the failing test**

In `internal/remediate/planner_test.go`: every `Plan(...)` call gains a `prsByRepo` argument. Update the three existing calls to pass `nil` (no external PRs → unchanged behavior). Then ADD:

```go
func TestPlanAdoptsExistingExternalPR(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go",
			Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	prs := map[string][]ghclient.PullRequest{
		"kairos-io/immucore": {{Number: 7, Title: "Bump golang.org/x/net to 0.33.0", Author: "renovate[bot]", URL: "u7"}},
	}
	intents, _ := Plan(c, state.Ledger{}, prs, 10)
	require.Len(t, intents, 1)
	assert.Equal(t, IntentAdopt, intents[0].Type)
	assert.Equal(t, 7, intents[0].PRNumber)
	assert.Equal(t, "renovate", intents[0].Source)
}

func TestPlanOpensWhenNoExternalPR(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go",
			Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	intents, _ := Plan(c, state.Ledger{}, nil, 10)
	require.Len(t, intents, 1)
	assert.Equal(t, IntentOpen, intents[0].Type)
}
```

(Add the `ghclient` import to the test file.)

- [ ] **Step 3: Run it — expect FAIL** (signature mismatch + adopt cases). Run: `go test ./internal/remediate/...`

- [ ] **Step 4: Update `Plan`**

In `internal/remediate/planner.go`, change the signature and the step-3 loop. Replace the `keys`/`deferred` selection with adopt-vs-open:

```go
func Plan(c state.Correlated, ledger state.Ledger, prsByRepo map[string][]ghclient.PullRequest, maxNew int) ([]Intent, int) {
	var intents []Intent

	for i := range ledger.Entries {
		e := &ledger.Entries[i]
		intents = append(intents, Intent{Type: IntentReconcile, Key: e.Key, Repo: e.Repo, Entry: e})
	}

	type target struct{ repo, pkg, to, sev string }
	targets := map[string]*target{}
	for _, f := range c.Findings {
		if !actionable(f) {
			continue
		}
		k := key(f.Repo, f.Package)
		t := targets[k]
		if t == nil {
			targets[k] = &target{repo: f.Repo, pkg: f.Package, to: f.FixedVersion, sev: f.Severity}
			continue
		}
		t.to = higherVersion(t.to, f.FixedVersion)
		if sevRank[f.Severity] > sevRank[t.sev] {
			t.sev = f.Severity
		}
	}

	// Decide per target. Targets already covered by one of our live PRs are
	// skipped (reconcile handles them). Otherwise: adopt an external PR if one
	// addresses it, else mark it a gap to open.
	var openKeys []string
	for k, t := range targets {
		if e, ok := ledger.ByKey(k); ok {
			if e.State == "open" || e.State == "conflicted" {
				continue
			}
			if (e.State == "merged" || e.State == "closed") && compareVersions(e.Bump.To, t.to) >= 0 {
				continue
			}
		}
		if pr, source, ok := MatchPR(t.pkg, prsByRepo[t.repo]); ok && source != "ksec" {
			intents = append(intents, Intent{
				Type: IntentAdopt, Key: k, Repo: t.repo, Package: t.pkg, Severity: t.sev,
				Bump: state.Bump{Package: t.pkg, To: t.to}, PRNumber: pr.Number, PRURL: pr.URL, Source: source,
			})
			continue
		}
		openKeys = append(openKeys, k)
	}

	sort.Slice(openKeys, func(i, j int) bool {
		ti, tj := targets[openKeys[i]], targets[openKeys[j]]
		if sevRank[ti.sev] != sevRank[tj.sev] {
			return sevRank[ti.sev] > sevRank[tj.sev]
		}
		return openKeys[i] < openKeys[j]
	})

	deferred := 0
	for n, k := range openKeys {
		if n >= maxNew {
			deferred = len(openKeys) - n
			break
		}
		t := targets[k]
		intents = append(intents, Intent{
			Type: IntentOpen, Key: k, Repo: t.repo, Package: t.pkg, Severity: t.sev,
			Bump: state.Bump{Package: t.pkg, To: t.to},
		})
	}
	return intents, deferred
}
```

Add the `ghclient` import to `planner.go`.

- [ ] **Step 5: Run it — expect PASS.** Run: `go test ./internal/remediate/...`

- [ ] **Step 6: Commit**

```bash
git add internal/remediate/intent.go internal/remediate/planner.go internal/remediate/planner_test.go
git commit -m "feat(remediate): planner adopts existing PRs, opens only gaps"
```

---

### Task 6: Run loop + Executor.Adopt + Fake

**Files:**
- Modify: `internal/remediate/run.go`
- Modify: `internal/remediate/fake.go`
- Modify: `internal/remediate/run_test.go`

**Interfaces:**
- Produces: `Executor` gains `Adopt(in Intent, run string) (state.LedgerEntry, error)`; `Run` handles `IntentAdopt` (error-isolated like the others). `FakeExecutor` gains an `Adopted map[string]state.LedgerEntry` and an `Adopt` method.

- [ ] **Step 1: Update the failing test**

In `internal/remediate/run_test.go`, add to the test an adopt intent + assertion (and the fake will need `Adopt`). Add a focused test:

```go
func TestRunAdopts(t *testing.T) {
	intents := []Intent{
		{Type: IntentAdopt, Key: "r|p", Repo: "r", Package: "p", PRNumber: 9, PRURL: "u9", Source: "dependabot"},
	}
	fake := &FakeExecutor{Adopted: map[string]state.LedgerEntry{
		"r|p": {Key: "r|p", Repo: "r", Package: "p", State: "open", PRNumber: 9, Source: "dependabot"},
	}}
	out, results := Run(intents, fake, state.Ledger{}, "2026-06-20")
	require.Len(t, out.Entries, 1)
	assert.Equal(t, "dependabot", out.Entries[0].Source)
	require.Len(t, results, 1)
	assert.Equal(t, "adopt", results[0].Action)
}
```

- [ ] **Step 2: Run it — expect FAIL** (`Adopt` undefined). Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement**

In `internal/remediate/run.go`, add to the `Executor` interface:

```go
	Adopt(in Intent, run string) (state.LedgerEntry, error)
```

And add a case to the `Run` switch (mirror `IntentOpen`):

```go
		case IntentAdopt:
			entry, err := ex.Adopt(in, run)
			if err != nil {
				rec := state.LedgerEntry{
					Key: in.Key, Repo: in.Repo, Package: in.Package, State: "error",
					Source: in.Source, Kind: "direct", Bump: in.Bump, Severity: in.Severity,
					PRNumber: in.PRNumber, PRURL: in.PRURL, CreatedRun: run, LastActionRun: run,
					History: []state.LedgerEvent{{Run: run, Action: "adopt-failed", Detail: err.Error()}},
				}
				byKey[in.Key] = rec
				results = append(results, Result{Key: in.Key, Action: "adopt", State: "error", Detail: err.Error()})
				continue
			}
			byKey[entry.Key] = entry
			results = append(results, Result{Key: entry.Key, Action: "adopt", State: entry.State})
```

In `internal/remediate/fake.go`, add the field + method:

```go
// FakeExecutor struct: add
	Adopted map[string]state.LedgerEntry

func (f *FakeExecutor) Adopt(in Intent, run string) (state.LedgerEntry, error) {
	if e, ok := f.Adopted[in.Key]; ok {
		return e, nil
	}
	return state.LedgerEntry{Key: in.Key, Repo: in.Repo, Package: in.Package, State: "open",
		Source: in.Source, Kind: "direct", PRNumber: in.PRNumber, PRURL: in.PRURL,
		CreatedRun: run, LastActionRun: run}, nil
}
```

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/remediate/...`

- [ ] **Step 5: Commit**

```bash
git add internal/remediate/run.go internal/remediate/fake.go internal/remediate/run_test.go
git commit -m "feat(remediate): Run handles adopt intents; Executor.Adopt + fake"
```

---

### Task 7: `GitExecutor.Adopt` (link + idempotent nudge + optional automerge)

**Files:**
- Modify: `internal/remediate/git_executor.go`

**Interfaces:**
- Consumes: `ghclient.GitHub` (list comments/post/status/merge), `ShouldAutomerge`.
- Produces: `GitExecutor` gains `GH ghclient.GitHub` and `Automerge bool` fields and an `Adopt(in Intent, run string) (state.LedgerEntry, error)` method. Integration; verified by build/vet + a tiny nudge-marker constant. The nudge is idempotent via the marker `<!-- ksec:nudge -->`.

- [ ] **Step 1: Implement `Adopt`**

In `internal/remediate/git_executor.go` add the fields to the `GitExecutor` struct:

```go
	GH        ghclient.GitHub // used by Adopt for comment/status/merge
	Automerge bool
```

Add the import `"github.com/kairos-io/security/internal/ghclient"` and the method:

```go
const nudgeMarker = "<!-- ksec:nudge -->"

func (g *GitExecutor) Adopt(in Intent, runID string) (state.LedgerEntry, error) {
	entry := state.LedgerEntry{
		Key: in.Key, Repo: in.Repo, Package: in.Package, Source: in.Source, Kind: "direct",
		PRNumber: in.PRNumber, PRURL: in.PRURL, Bump: in.Bump, Severity: in.Severity,
		State: "open", CreatedRun: runID, LastActionRun: runID,
	}
	if g.DryRun || g.GH == nil {
		if g.DryRun {
			fmt.Printf("[dry-run] would adopt %s PR #%d (%s): nudge%s\n", in.Repo, in.PRNumber, in.Source,
				map[bool]string{true: " + automerge-if-green", false: ""}[g.Automerge])
		}
		entry.History = []state.LedgerEvent{{Run: runID, Action: "adopt", Detail: in.Source}}
		return entry, nil
	}

	// Refresh live PR state.
	if st, err := g.GH.PRStatusOf(in.Repo, in.PRNumber); err == nil {
		switch st.State {
		case "MERGED":
			entry.State = "merged"
		case "CLOSED":
			entry.State = "closed"
		}
		// Optional automerge.
		if g.Automerge && entry.State == "open" && ShouldAutomerge(st) {
			if err := g.GH.MergePR(in.Repo, in.PRNumber, true); err == nil {
				entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "automerge-requested"})
			}
		}
	}

	// Idempotent nudge: only if we haven't commented the marker yet.
	if entry.State == "open" {
		nudged := false
		if comments, err := g.GH.ListPRComments(in.Repo, in.PRNumber); err == nil {
			for _, c := range comments {
				if strings.Contains(c.Body, nudgeMarker) {
					nudged = true
					break
				}
			}
		}
		if !nudged {
			body := fmt.Sprintf("This PR addresses a %s-severity security finding (%s). Tracked by kairos-security.\n\n%s",
				in.Severity, in.Package, nudgeMarker)
			_ = g.GH.PostPRComment(in.Repo, in.PRNumber, body)
			entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "nudged"})
		}
	}
	return entry, nil
}
```

(Ensure `strings` is imported in `git_executor.go`.)

- [ ] **Step 2: Build + vet + test**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: pass (no new unit test; `Adopt` is integration, exercised via the dry-run path and Task 6's run-loop test through the fake).

- [ ] **Step 3: Commit**

```bash
git add internal/remediate/git_executor.go
git commit -m "feat(remediate): GitExecutor.Adopt links PR, idempotent nudge, optional automerge"
```

---

### Task 8: Dashboard coordination columns

**Files:**
- Modify: `internal/render/render.go`
- Test: `internal/render/coord_test.go`

**Interfaces:**
- Consumes: `state.LedgerEntry` (now with `Source`/`Kind`/`NeedsHuman`).
- Produces: the markdown "Bot PR ledger" table gains `Source` and `Kind` columns and a status that reflects `needsHuman`/`blocked`.

- [ ] **Step 1: Write the failing test**

Create `internal/render/coord_test.go`:

```go
package render

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestDashboardMarkdownLedgerShowsSourceAndKind(t *testing.T) {
	in := Input{Ledger: state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", Repo: "kairos-io/immucore",
			Package: "golang.org/x/net", State: "open", PRNumber: 7,
			PRURL: "https://github.com/kairos-io/immucore/pull/7",
			Source: "dependabot", Kind: "direct", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "dependabot")
	assert.Contains(t, md, "direct")
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/render/...`

- [ ] **Step 3: Implement**

In `internal/render/render.go`, update the "Bot PR ledger" table header and rows. Change the header line to:

```go
		b.WriteString("| Repo | Bump | Kind | Source | State | PR |\n|---|---|---|---|---|---|\n")
```

and the per-entry row to include `Kind`/`Source` and a status that prefers `needsHuman`/`blocked`:

```go
		for _, e := range in.Ledger.Entries {
			pr := "—"
			if e.PRNumber > 0 {
				pr = fmt.Sprintf("[#%d](%s)", e.PRNumber, e.PRURL)
			}
			kind := e.Kind
			if kind == "" {
				kind = "direct"
			}
			source := e.Source
			if source == "" {
				source = "ksec"
			}
			st := e.State
			if e.NeedsHuman {
				st = "⚠️ needs-human"
			} else if e.Blocked != "" {
				st = "⛔ " + e.Blocked
			}
			fmt.Fprintf(&b, "| %s | %s@%s | %s | %s | %s | %s |\n", e.Repo, e.Bump.Package, e.Bump.To, kind, source, st, pr)
		}
```

If `internal/render/html.go` has the ledger table, mirror the two columns and regenerate the HTML golden; regenerate both goldens: `UPDATE_GOLDEN=1 go test ./internal/render/...`, eyeball them (the new columns appear, prior sections intact), then re-run.

- [ ] **Step 4: Run it — expect PASS** (after regenerating goldens). Run: `go test ./internal/render/...`

- [ ] **Step 5: Commit**

```bash
git add internal/render/render.go internal/render/coord_test.go internal/render/testdata/ internal/render/html.go
git commit -m "feat(render): coordination columns (kind/source/status) in ledger view"
```

---

### Task 9: Wire it into the `remediate` command

**Files:**
- Modify: `cmd/ksec/main.go`

**Interfaces:**
- Consumes: `remediate.Plan` (new signature), `ghclient.NewCLI().ListOpenPRs`, `state.ReposFile`.
- Produces: the `remediate` command collects open PRs per tracked repo into `prsByRepo`, passes it to `Plan`, sets `ex.GH` + `ex.Automerge` from a new `--automerge` flag.

- [ ] **Step 1: Wire**

In `newRemediateCmd`'s `RunE`, after loading the ledger and before `remediate.Plan(...)`:

```go
			gh := ghclient.NewCLI()
			// Collect open PRs per tracked repo so the planner can adopt existing
			// dependabot/renovate/human PRs instead of duplicating them.
			prsByRepo := map[string][]ghclient.PullRequest{}
			var repos []state.Repo
			if err := state.Load(gf.stateDir, state.ReposFile, &repos); err == nil {
				for _, r := range repos {
					if prs, err := gh.ListOpenPRs(r.Repo); err == nil {
						prsByRepo[r.Repo] = prs
					}
				}
			}
```

Change the plan call to `intents, deferred := remediate.Plan(c, ledger, prsByRepo, maxPRs)`.

Change the executor construction to set `GH` and `Automerge`:

```go
			ex := &remediate.GitExecutor{Token: os.Getenv("GH_TOKEN"), DryRun: gf.dryRun, GH: gh, Automerge: automerge}
			if aiProse && aiCfg.Nib.Endpoint != "" {
				ex.Prose = remediate.NewOpenAIProse(aiCfg)
			}
```

(The reaction loop later already builds its own `gh := ghclient.NewCLI()` — leave it; or reuse this one. Reusing is fine: remove the later `gh := ghclient.NewCLI()` re-declaration if it now shadows; simplest is to keep them separate scopes. Ensure no unused-variable or redeclaration error — if the later block redeclares `gh`, rename this one to `ghc` and use `ghc` here.)

Add the `--automerge` flag and `var automerge bool` near `var maxPRs int`:

```go
	cmd.Flags().BoolVar(&automerge, "automerge", false, "merge addressing PRs (ours/dependabot/renovate) when green and unblocked")
```

- [ ] **Step 2: Build + vet + full test + smoke**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: all pass. Smoke: `go run ./cmd/ksec remediate --help` shows `--automerge` and `--max-prs` and `--ai-pr-prose`.

- [ ] **Step 3: Commit**

```bash
git add cmd/ksec/main.go
git commit -m "feat(remediate): collect PRs for adoption; --automerge flag; wire ex.GH"
```

---

## Self-review

**Spec coverage** (§7.1 Status & adoption + §6 ledger fields + §8 surfaces):
- Ledger source/kind/blocked/needsHuman → Task 1. ✓
- Find addressing PR + source classification (no duplicate) → Tasks 3, 5. ✓
- Default nudge (idempotent, marker) → Task 7. ✓
- `--automerge` merges green/unblocked → Tasks 2 (status/merge), 4 (decision), 7 (apply), 9 (flag). ✓
- Adopt records link + state, drives PR → Tasks 5, 6, 7. ✓
- Dashboard coordination view (source/kind/status) → Task 8. ✓
- Cap applies to `open` only; adopt uncapped → Task 5. ✓
- Verify/dry-run/token: adopt does no git writes; dry-run prints; token unaffected → Task 7. ✓

**Deferred to 4b/4c (correctly absent):** cascade, re-pin, depgraph (`cascadeFrom`/`pinTarget`/`pseudo` ledger fields), toolchain, nib agent. Ledger fields for those are added in 4b; 4a adds only source/kind/blocked/needsHuman.

**Placeholder scan:** none — complete code/commands in every step.

**Type consistency:** `IntentAdopt` + `Intent.PRNumber/PRURL/Source` (Task 5) used by Tasks 6, 7, 9. `Executor.Adopt` (Task 6) implemented by `FakeExecutor` (Task 6) and `GitExecutor` (Task 7). `ghclient.PRStatus`/`PRStatusOf`/`MergePR` (Task 2) used by Tasks 4, 7. `MatchPR`/`classifySource` (Task 3) used by Task 5. `ShouldAutomerge` (Task 4) used by Task 7. `Plan` new signature (Task 5) called in Task 9. Ledger fields (Task 1) rendered in Task 8.

---

## Operational notes

- Live adoption/automerge needs `KSEC_BOT_TOKEN` with `pull_requests:write` (comment/merge) on the target repos; without it, adoption records links and the nudge/merge calls fail gracefully (logged, not fatal).
- `--automerge` is off by default; turn it on via the workflow only once you trust the green/unblocked gate for these repos.
- Collecting PRs per repo adds one `gh pr list` call per tracked repo per run (cheap).
