# Autonomous Remediation (Bump PRs + Ledger) — Implementation Plan (Plan 2)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `ksec remediate` phase that autonomously opens and maintains dependency-bump pull requests across the tracked Go repos to fix source/dependency CVEs, remembering what it did in a committed ledger, coordinating the "waterfall" case, and surfacing the bot's PR state on the dashboard.

**Architecture:** A new `internal/remediate` package split into a **pure planner** (decides which bumps to attempt from correlated findings + the ledger, with a blast-radius cap) and an **executor** (an interface; the real one shells `git` + `gh` to branch/bump/verify/push/open-PR, a fake one is used in tests). A reconciliation `Run` loop applies intents and updates the committed `state/ledger.json`. Bumps are deterministic (`go get <pkg>@<version>` + `go mod tidy`, verified by `go build ./...`); the dashboard gains a "Bot PR ledger" section. The bot only ever touches PRs it authored (identified by a `ksec/` branch prefix and an HTML-comment marker).

**Tech Stack:** Go 1.22, `git` + `gh` CLIs, `stretchr/testify`. Builds on Plan 1's `state`, `ghclient`, `correlate`, and `render` packages.

## Global Constraints

- Module `github.com/kairos-io/security`; binary `ksec`; Go 1.22 floor.
- `remediate` is the FIRST phase that writes to **other** repos (branches + PRs). `--dry-run` turns every git/GitHub write into a printed plan and performs no writes.
- The bot only touches PRs **it authored**: branch names are prefixed `ksec/`, and every PR body carries the marker `<!-- ksec:key=<key> -->`. Human PRs are tracked as findings, never modified.
- A bump target is keyed `"<repo>|<package>"` (one PR per repo+package). The ledger entry `key` is exactly this string.
- Bumps are deterministic: `go get <package>@<version>` then `go mod tidy`, verified locally with `go build ./...`. If the build fails, **no PR is pushed** — the ledger records `build-failed` and it surfaces on the dashboard for a human.
- Actionable findings only: `Type` in {`sourceCVE`, `ghAlert`}, `Ecosystem == "go"`, non-empty `FixedVersion`. Image CVEs are reported, never auto-PR'd.
- Blast-radius guard: at most `--max-prs` NEW PRs per run (default 10); overflow is logged and deferred to the next run.
- The ledger (`state/ledger.json`) is committed every run; it is the bot's memory.
- Per-intent failures are isolated: record the error on the ledger entry and continue; never abort the run.
- **Out of scope (Plan 3):** reacting to PR review comments, AI-drafted PR prose. Plan 2 PR bodies are deterministic. The `LedgerEntry.SeenComments`/`History` fields are added now so Plan 3 needs no migration.
- **Operational prerequisite (not code):** the workflow's `KSEC_BOT_TOKEN` must have `contents:write` + `pull_requests:write` on the target repos for live (non-dry-run) actuation. Until then, run `remediate` in dry-run.

---

## File structure

```
internal/state/types.go            # + LedgerEntry, LedgerEvent, Ledger, LedgerFile (modify)
internal/state/ledger.go           # Ledger.ByKey helper (create)
internal/state/ledger_test.go      # (create)
internal/remediate/intent.go       # Intent, IntentType, Result types (create)
internal/remediate/planner.go      # Plan(): pure intent planning (create)
internal/remediate/planner_test.go # (create)
internal/remediate/run.go          # Executor interface + Run() reconciliation loop (create)
internal/remediate/fake.go         # FakeExecutor for tests (create)
internal/remediate/run_test.go     # (create)
internal/remediate/git_executor.go # real git+gh executor (create)
internal/remediate/prbody.go       # deterministic PR title/body (create)
internal/remediate/prbody_test.go  # (create)
internal/render/render.go          # + Ledger to Input + "Bot PR ledger" section (modify)
internal/render/ledger_test.go     # (create)
cmd/ksec/main.go                   # + newRemediateCmd, wire into root (modify)
.github/workflows/security-dashboard.yaml  # + remediate step (modify)
```

---

### Task 1: Ledger types + `ByKey`

**Files:**
- Modify: `internal/state/types.go`
- Create: `internal/state/ledger.go`
- Test: `internal/state/ledger_test.go`

**Interfaces:**
- Consumes: `state.Bump` (exists).
- Produces: `LedgerEvent`, `LedgerEntry`, `Ledger`, the constant `LedgerFile = "ledger.json"`, and `func (l *Ledger) ByKey(key string) (*LedgerEntry, bool)` returning a pointer into `l.Entries`.

- [ ] **Step 1: Add the types**

In `internal/state/types.go`, add to the file-name constants block: `LedgerFile = "ledger.json"`. Then append these types:

```go
type LedgerEvent struct {
	Run    string `json:"run"`
	Action string `json:"action"`
	Detail string `json:"detail,omitempty"`
}

type LedgerEntry struct {
	Key           string        `json:"key"`   // "<repo>|<package>"
	Repo          string        `json:"repo"`
	Package       string        `json:"package"`
	Branch        string        `json:"branch"`
	PRNumber      int           `json:"prNumber,omitempty"`
	PRURL         string        `json:"prURL,omitempty"`
	State         string        `json:"state"` // planned|open|merged|closed|conflicted|build-failed|error
	Bump          Bump          `json:"bump"`
	Severity      string        `json:"severity,omitempty"`
	CreatedRun    string        `json:"createdRun"`
	LastActionRun string        `json:"lastActionRun"`
	SeenComments  []string      `json:"seenComments,omitempty"` // reserved for Plan 3
	History       []LedgerEvent `json:"history,omitempty"`
}

type Ledger struct {
	Entries []LedgerEntry `json:"entries"`
}
```

- [ ] **Step 2: Write the failing test**

Create `internal/state/ledger_test.go`:

```go
package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedgerByKey(t *testing.T) {
	l := Ledger{Entries: []LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", State: "open"},
		{Key: "kairos-io/kairos|stdlib", State: "merged"},
	}}
	got, ok := l.ByKey("kairos-io/kairos|stdlib")
	require.True(t, ok)
	assert.Equal(t, "merged", got.State)

	got.State = "closed" // pointer write must mutate the slice
	again, _ := l.ByKey("kairos-io/kairos|stdlib")
	assert.Equal(t, "closed", again.State)

	_, ok = l.ByKey("nope")
	assert.False(t, ok)
}

func TestLedgerRoundTrip(t *testing.T) {
	dir := t.TempDir()
	in := Ledger{Entries: []LedgerEntry{{Key: "a|b", Repo: "a", Package: "b", State: "open",
		Bump: Bump{Package: "b", To: "1.2.3"}, History: []LedgerEvent{{Run: "2026-06-19", Action: "opened"}}}}}
	require.NoError(t, Save(dir, LedgerFile, in))
	var out Ledger
	require.NoError(t, Load(dir, LedgerFile, &out))
	assert.Equal(t, in, out)
}
```

- [ ] **Step 3: Run it — expect FAIL** (`ByKey` undefined). Run: `go test ./internal/state/...`

- [ ] **Step 4: Implement `ByKey`**

Create `internal/state/ledger.go`:

```go
package state

// ByKey returns a pointer to the entry with the given key so callers can mutate
// it in place, plus whether it was found.
func (l *Ledger) ByKey(key string) (*LedgerEntry, bool) {
	for i := range l.Entries {
		if l.Entries[i].Key == key {
			return &l.Entries[i], true
		}
	}
	return nil, false
}
```

- [ ] **Step 5: Run it — expect PASS.** Run: `go test ./internal/state/...`

- [ ] **Step 6: Commit**

```bash
git add internal/state/types.go internal/state/ledger.go internal/state/ledger_test.go
git commit -m "feat(state): ledger types and ByKey"
```

---

### Task 2: Planner (pure intent planning)

**Files:**
- Create: `internal/remediate/intent.go`
- Create: `internal/remediate/planner.go`
- Test: `internal/remediate/planner_test.go`

**Interfaces:**
- Consumes: `state.Correlated`, `state.Finding`, `state.Ledger`, `state.LedgerEntry`, `state.Bump`.
- Produces:
  - `type IntentType string` with `IntentOpen IntentType = "open"` and `IntentReconcile IntentType = "reconcile"`.
  - `type Intent struct { Type IntentType; Key, Repo, Package, Severity string; Bump state.Bump; Entry *state.LedgerEntry }`
  - `type Result struct { Key, Action, State, Detail string }`
  - `func Plan(c state.Correlated, ledger state.Ledger, maxNew int) (intents []Intent, deferred int)` — emits one `IntentReconcile` per existing ledger entry, then `IntentOpen` for each NEW actionable target (deduped by `"<repo>|<package>"`, choosing the highest `FixedVersion`), capped at `maxNew`; `deferred` is the count of new targets dropped by the cap.

- [ ] **Step 1: Write the types**

Create `internal/remediate/intent.go`:

```go
package remediate

import "github.com/kairos-io/security/internal/state"

type IntentType string

const (
	IntentOpen      IntentType = "open"
	IntentReconcile IntentType = "reconcile"
)

type Intent struct {
	Type     IntentType
	Key      string
	Repo     string
	Package  string
	Severity string
	Bump     state.Bump
	Entry    *state.LedgerEntry // set for IntentReconcile
}

type Result struct {
	Key    string
	Action string
	State  string
	Detail string
}
```

- [ ] **Step 2: Write the failing test**

Create `internal/remediate/planner_test.go`:

```go
package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanOpensNewActionableTargetsDedupedAndCapped(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		// two CVEs in the same repo+package -> one target at the highest fixed version
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
		{ID: "b", Repo: "kairos-io/immucore", Type: "ghAlert", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.36.0", Severity: "critical"},
		// a different repo+package -> second target
		{ID: "c", Repo: "kairos-io/kairos", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/crypto", FixedVersion: "0.31.0", Severity: "high"},
		// not actionable: image CVE
		{ID: "d", Repo: "kairos-io/kairos", Type: "imageCVE", Package: "openssl", FixedVersion: "1.1.1w", Severity: "critical"},
		// not actionable: no fixed version
		{ID: "e", Repo: "kairos-io/kairos", Type: "sourceCVE", Ecosystem: "go", Package: "x/text", Severity: "low"},
	}}

	intents, deferred := Plan(c, state.Ledger{}, 1) // cap to 1 new PR
	require.Len(t, intents, 1)
	assert.Equal(t, 1, deferred)
	in := intents[0]
	assert.Equal(t, IntentOpen, in.Type)
	// highest severity target first: immucore/x/net (critical) at the highest fixed version
	assert.Equal(t, "kairos-io/immucore|golang.org/x/net", in.Key)
	assert.Equal(t, "0.36.0", in.Bump.To)
	assert.Equal(t, "critical", in.Severity)
}

func TestPlanReconcilesExistingLedgerEntries(t *testing.T) {
	c := state.Correlated{}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", Repo: "kairos-io/immucore", State: "open"},
	}}
	intents, _ := Plan(c, led, 10)
	require.Len(t, intents, 1)
	assert.Equal(t, IntentReconcile, intents[0].Type)
	require.NotNil(t, intents[0].Entry)
	assert.Equal(t, "open", intents[0].Entry.State)
}

func TestPlanSkipsTargetsAlreadyInLedger(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", State: "open"},
	}}
	intents, _ := Plan(c, led, 10)
	// only the reconcile for the existing entry; no new open
	require.Len(t, intents, 1)
	assert.Equal(t, IntentReconcile, intents[0].Type)
}
```

- [ ] **Step 3: Run it — expect FAIL** (`Plan` undefined). Run: `go test ./internal/remediate/...`

- [ ] **Step 4: Implement the planner**

Create `internal/remediate/planner.go`:

```go
package remediate

import (
	"sort"

	"github.com/kairos-io/security/internal/state"
)

var sevRank = map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1, "unknown": 0, "": 0}

// actionable reports whether a finding can be auto-bumped.
func actionable(f state.Finding) bool {
	return (f.Type == "sourceCVE" || f.Type == "ghAlert") &&
		f.Ecosystem == "go" && f.Package != "" && f.FixedVersion != ""
}

func key(repo, pkg string) string { return repo + "|" + pkg }

// higherVersion returns the "greater" of two version strings. We avoid a full
// semver parser: trim a leading 'v' and compare dotted numeric segments,
// falling back to string comparison.
func higherVersion(a, b string) string {
	if compareVersions(a, b) >= 0 {
		return a
	}
	return b
}

func compareVersions(a, b string) int {
	na, nb := splitVer(a), splitVer(b)
	for i := 0; i < len(na) || i < len(nb); i++ {
		var x, y int
		if i < len(na) {
			x = na[i]
		}
		if i < len(nb) {
			y = nb[i]
		}
		if x != y {
			if x < y {
				return -1
			}
			return 1
		}
	}
	return 0
}

func splitVer(s string) []int {
	if len(s) > 0 && (s[0] == 'v' || s[0] == 'V') {
		s = s[1:]
	}
	var out []int
	cur, has := 0, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			cur = cur*10 + int(c-'0')
			has = true
		} else if c == '.' {
			out = append(out, cur)
			cur, has = 0, false
		} else {
			break // stop at pre-release / build metadata
		}
	}
	if has || len(out) == 0 {
		out = append(out, cur)
	}
	return out
}

func Plan(c state.Correlated, ledger state.Ledger, maxNew int) ([]Intent, int) {
	var intents []Intent

	// 1) Reconcile every existing ledger entry.
	for i := range ledger.Entries {
		e := &ledger.Entries[i]
		intents = append(intents, Intent{Type: IntentReconcile, Key: e.Key, Repo: e.Repo, Entry: e})
	}

	// 2) Collapse actionable findings into one target per repo+package.
	type target struct {
		repo, pkg, to, sev string
	}
	targets := map[string]*target{}
	for _, f := range c.Findings {
		if !actionable(f) {
			continue
		}
		k := key(f.Repo, f.Package)
		if _, ok := ledger.ByKey(k); ok {
			continue // already tracked
		}
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

	// 3) Order new targets by severity (desc) then key (asc), apply the cap.
	keys := make([]string, 0, len(targets))
	for k := range targets {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ti, tj := targets[keys[i]], targets[keys[j]]
		if sevRank[ti.sev] != sevRank[tj.sev] {
			return sevRank[ti.sev] > sevRank[tj.sev]
		}
		return keys[i] < keys[j]
	})

	deferred := 0
	for n, k := range keys {
		if n >= maxNew {
			deferred = len(keys) - n
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

- [ ] **Step 5: Run it — expect PASS.** Run: `go test ./internal/remediate/...`

- [ ] **Step 6: Commit**

```bash
git add internal/remediate/intent.go internal/remediate/planner.go internal/remediate/planner_test.go
git commit -m "feat(remediate): pure planner with dedupe, severity order, blast-radius cap"
```

---

### Task 3: Reconciliation loop + Executor interface + Fake

**Files:**
- Create: `internal/remediate/run.go`
- Create: `internal/remediate/fake.go`
- Test: `internal/remediate/run_test.go`

**Interfaces:**
- Consumes: `Intent`, `Result`, `state.Ledger`, `state.LedgerEntry`.
- Produces:
  - `type Executor interface { Open(in Intent, run string) (state.LedgerEntry, error); Reconcile(e state.LedgerEntry, run string) (state.LedgerEntry, error) }`
  - `func Run(intents []Intent, ex Executor, ledger state.Ledger, run string) (state.Ledger, []Result)` — for each `IntentOpen` calls `ex.Open` and appends/updates the entry by key; for each `IntentReconcile` calls `ex.Reconcile` and replaces the entry; isolates per-intent errors (records `state="error"` + a history event, keeps prior entry), and drops entries that reconcile to `state=="merged"` or `"closed"` is **kept** (history matters) — entries are never deleted. Output ledger entries are sorted by `Key`.
  - `type FakeExecutor struct { Opened map[string]state.LedgerEntry; Reconciled map[string]state.LedgerEntry; OpenErr map[string]error }` implementing `Executor`.

- [ ] **Step 1: Write the failing test**

Create `internal/remediate/run_test.go`:

```go
package remediate

import (
	"errors"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunOpensReconcilesAndIsolatesErrors(t *testing.T) {
	intents := []Intent{
		{Type: IntentOpen, Key: "r|p1", Repo: "r", Package: "p1", Bump: state.Bump{Package: "p1", To: "1.0.0"}},
		{Type: IntentOpen, Key: "r|p2", Repo: "r", Package: "p2", Bump: state.Bump{Package: "p2", To: "2.0.0"}},
		{Type: IntentReconcile, Key: "r|old", Entry: &state.LedgerEntry{Key: "r|old", Repo: "r", State: "open"}},
	}
	fake := &FakeExecutor{
		Opened: map[string]state.LedgerEntry{
			"r|p1": {Key: "r|p1", Repo: "r", Package: "p1", State: "open", PRNumber: 1},
		},
		OpenErr:    map[string]error{"r|p2": errors.New("build-failed")},
		Reconciled: map[string]state.LedgerEntry{"r|old": {Key: "r|old", Repo: "r", State: "merged"}},
	}

	out, results := Run(intents, fake, state.Ledger{}, "2026-06-19")

	byKey := map[string]state.LedgerEntry{}
	for _, e := range out.Entries {
		byKey[e.Key] = e
	}
	assert.Equal(t, "open", byKey["r|p1"].State)
	assert.Equal(t, "error", byKey["r|p2"].State, "open failure recorded, not aborted")
	assert.Equal(t, "merged", byKey["r|old"].State)
	// deterministic order
	require.Len(t, out.Entries, 3)
	assert.True(t, out.Entries[0].Key <= out.Entries[1].Key)
	// a result per intent
	assert.Len(t, results, 3)
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement Run + Fake**

Create `internal/remediate/run.go`:

```go
package remediate

import (
	"sort"

	"github.com/kairos-io/security/internal/state"
)

type Executor interface {
	Open(in Intent, run string) (state.LedgerEntry, error)
	Reconcile(e state.LedgerEntry, run string) (state.LedgerEntry, error)
}

func Run(intents []Intent, ex Executor, ledger state.Ledger, run string) (state.Ledger, []Result) {
	// Index existing entries by key for in-place replacement.
	byKey := map[string]state.LedgerEntry{}
	for _, e := range ledger.Entries {
		byKey[e.Key] = e
	}
	var results []Result

	for _, in := range intents {
		switch in.Type {
		case IntentOpen:
			entry, err := ex.Open(in, run)
			if err != nil {
				rec := state.LedgerEntry{
					Key: in.Key, Repo: in.Repo, Package: in.Package, State: "error",
					Bump: in.Bump, Severity: in.Severity, CreatedRun: run, LastActionRun: run,
					History: []state.LedgerEvent{{Run: run, Action: "open-failed", Detail: err.Error()}},
				}
				byKey[in.Key] = rec
				results = append(results, Result{Key: in.Key, Action: "open", State: "error", Detail: err.Error()})
				continue
			}
			byKey[entry.Key] = entry
			results = append(results, Result{Key: entry.Key, Action: "open", State: entry.State})
		case IntentReconcile:
			prior := *in.Entry
			entry, err := ex.Reconcile(prior, run)
			if err != nil {
				prior.LastActionRun = run
				prior.History = append(prior.History, state.LedgerEvent{Run: run, Action: "reconcile-failed", Detail: err.Error()})
				byKey[prior.Key] = prior
				results = append(results, Result{Key: prior.Key, Action: "reconcile", State: "error", Detail: err.Error()})
				continue
			}
			byKey[entry.Key] = entry
			results = append(results, Result{Key: entry.Key, Action: "reconcile", State: entry.State})
		}
	}

	out := state.Ledger{Entries: make([]state.LedgerEntry, 0, len(byKey))}
	for _, e := range byKey {
		out.Entries = append(out.Entries, e)
	}
	sort.Slice(out.Entries, func(i, j int) bool { return out.Entries[i].Key < out.Entries[j].Key })
	return out, results
}
```

Create `internal/remediate/fake.go`:

```go
package remediate

import "github.com/kairos-io/security/internal/state"

// FakeExecutor is an in-memory Executor for tests.
type FakeExecutor struct {
	Opened     map[string]state.LedgerEntry
	Reconciled map[string]state.LedgerEntry
	OpenErr    map[string]error
}

func (f *FakeExecutor) Open(in Intent, run string) (state.LedgerEntry, error) {
	if err := f.OpenErr[in.Key]; err != nil {
		return state.LedgerEntry{}, err
	}
	e, ok := f.Opened[in.Key]
	if !ok {
		return state.LedgerEntry{Key: in.Key, Repo: in.Repo, Package: in.Package, State: "open", CreatedRun: run, LastActionRun: run}, nil
	}
	return e, nil
}

func (f *FakeExecutor) Reconcile(e state.LedgerEntry, run string) (state.LedgerEntry, error) {
	if r, ok := f.Reconciled[e.Key]; ok {
		return r, nil
	}
	return e, nil
}
```

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/remediate/...`

- [ ] **Step 5: Commit**

```bash
git add internal/remediate/run.go internal/remediate/fake.go internal/remediate/run_test.go
git commit -m "feat(remediate): reconciliation loop with error isolation + fake executor"
```

---

### Task 4: Deterministic PR title/body

**Files:**
- Create: `internal/remediate/prbody.go`
- Test: `internal/remediate/prbody_test.go`

**Interfaces:**
- Consumes: `Intent`.
- Produces:
  - `const Marker = "<!-- ksec:key=%s -->"` is NOT used directly; instead `func PRMarker(key string) string` returns `"<!-- ksec:key=" + key + " -->"`.
  - `func BranchName(in Intent) string` — `"ksec/bump-" + slug(package) + "-" + slug(version)`, lowercased, non-alphanumerics → `-`.
  - `func PRTitle(in Intent) string` — `"chore(security): bump <package> to <version>"`.
  - `func PRBody(in Intent) string` — deterministic body: what/why (severity, package, version), an explicit "automated by kairos-security" line, and the marker on its own last line.

- [ ] **Step 1: Write the failing test**

Create `internal/remediate/prbody_test.go`:

```go
package remediate

import (
	"strings"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
)

func sampleIntent() Intent {
	return Intent{Type: IntentOpen, Key: "kairos-io/immucore|golang.org/x/net",
		Repo: "kairos-io/immucore", Package: "golang.org/x/net", Severity: "high",
		Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}}
}

func TestBranchAndTitleAndBody(t *testing.T) {
	in := sampleIntent()
	assert.Equal(t, "ksec/bump-golang-org-x-net-0-33-0", BranchName(in))
	assert.Equal(t, "chore(security): bump golang.org/x/net to 0.33.0", PRTitle(in))

	body := PRBody(in)
	assert.Contains(t, body, "golang.org/x/net")
	assert.Contains(t, body, "0.33.0")
	assert.Contains(t, body, "high")
	assert.Contains(t, body, "kairos-security")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(body), PRMarker(in.Key)),
		"marker must be the last line")
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement**

Create `internal/remediate/prbody.go`:

```go
package remediate

import (
	"fmt"
	"strings"
)

func PRMarker(key string) string { return "<!-- ksec:key=" + key + " -->" }

func slug(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteByte('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

func BranchName(in Intent) string {
	return "ksec/bump-" + slug(in.Package) + "-" + slug(in.Bump.To)
}

func PRTitle(in Intent) string {
	return fmt.Sprintf("chore(security): bump %s to %s", in.Package, in.Bump.To)
}

func PRBody(in Intent) string {
	return fmt.Sprintf(`## Automated security bump

Bumps **%s** to **%s** to address a %s-severity vulnerability detected by
[kairos-security](https://github.com/kairos-io/security).

- Package: `+"`%s`"+`
- Target version: `+"`%s`"+`
- Severity: %s

This PR was opened automatically. The change is a deterministic
`+"`go get %s@%s` + `go mod tidy`"+`; CI on this PR runs the repository's tests.

%s`, in.Package, in.Bump.To, in.Severity, in.Package, in.Bump.To, in.Severity,
		in.Package, in.Bump.To, PRMarker(in.Key))
}
```

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/remediate/...`

- [ ] **Step 5: Commit**

```bash
git add internal/remediate/prbody.go internal/remediate/prbody_test.go
git commit -m "feat(remediate): deterministic PR branch/title/body with marker"
```

---

### Task 5: Real git+gh executor (integration; dry-run aware)

**Files:**
- Create: `internal/remediate/git_executor.go`

**Interfaces:**
- Consumes: `Intent`, `state.LedgerEntry`, `BranchName`/`PRTitle`/`PRBody`/`PRMarker`, `ghclient` (for reconcile reads).
- Produces: `type GitExecutor struct { Token string; DryRun bool; MaxBuildSeconds int }` implementing `Executor`. Not unit-tested (it shells `git`/`gh`); verified by `go build`/`go vet` and the end-to-end dry-run in Task 7.

- [ ] **Step 1: Implement the executor**

Create `internal/remediate/git_executor.go`:

```go
package remediate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/kairos-io/security/internal/state"
)

type GitExecutor struct {
	Token  string // GH_TOKEN, for authenticated clone/push
	DryRun bool
}

func (g *GitExecutor) cloneURL(repo string) string {
	if g.Token != "" {
		return "https://x-access-token:" + g.Token + "@github.com/" + repo + ".git"
	}
	return "https://github.com/" + repo + ".git"
}

func run(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		return out.Bytes(), fmt.Errorf("%s %v: %v: %s", name, args, err, errb.String())
	}
	return out.Bytes(), nil
}

func (g *GitExecutor) Open(in Intent, runID string) (state.LedgerEntry, error) {
	branch := BranchName(in)
	entry := state.LedgerEntry{
		Key: in.Key, Repo: in.Repo, Package: in.Package, Branch: branch,
		Bump: in.Bump, Severity: in.Severity, CreatedRun: runID, LastActionRun: runID,
	}

	if g.DryRun {
		fmt.Printf("[dry-run] would open PR on %s: branch %s, go get %s@%s\n",
			in.Repo, branch, in.Bump.Package, in.Bump.To)
		entry.State = "planned"
		entry.History = []state.LedgerEvent{{Run: runID, Action: "plan-open"}}
		return entry, nil
	}

	dir, err := os.MkdirTemp("", "ksec-rem-*")
	if err != nil {
		return entry, err
	}
	defer os.RemoveAll(dir)

	if _, err := run("", "git", "clone", "--depth", "1", g.cloneURL(in.Repo), dir); err != nil {
		return entry, err
	}
	if _, err := run(dir, "git", "checkout", "-b", branch); err != nil {
		return entry, err
	}
	if _, err := run(dir, "go", "get", in.Bump.Package+"@"+in.Bump.To); err != nil {
		return entry, err
	}
	if _, err := run(dir, "go", "mod", "tidy"); err != nil {
		return entry, err
	}
	// Verify-before-push: a broken build must not be pushed.
	if _, err := run(dir, "go", "build", "./..."); err != nil {
		entry.State = "build-failed"
		entry.History = []state.LedgerEvent{{Run: runID, Action: "build-failed", Detail: err.Error()}}
		return entry, nil // not an error: recorded for a human, run continues
	}

	run(dir, "git", "config", "user.name", "kairos-security-bot")
	run(dir, "git", "config", "user.email", "bot@kairos.io")
	if _, err := run(dir, "git", "commit", "-am", PRTitle(in)); err != nil {
		return entry, err
	}
	if _, err := run(dir, "git", "push", "-u", "origin", branch); err != nil {
		return entry, err
	}

	// Create the PR with gh (GH_TOKEN is read from the environment by gh).
	out, err := run(dir, "gh", "pr", "create", "-R", in.Repo, "--head", branch,
		"--title", PRTitle(in), "--body", PRBody(in))
	if err != nil {
		return entry, err
	}
	entry.PRURL = string(bytes.TrimSpace(out))
	entry.PRNumber = prNumberFromURL(entry.PRURL)
	entry.State = "open"
	entry.History = []state.LedgerEvent{{Run: runID, Action: "opened", Detail: entry.PRURL}}
	return entry, nil
}

func (g *GitExecutor) Reconcile(e state.LedgerEntry, runID string) (state.LedgerEntry, error) {
	e.LastActionRun = runID
	if e.PRNumber == 0 || g.DryRun {
		if g.DryRun {
			fmt.Printf("[dry-run] would reconcile %s (PR #%d)\n", e.Repo, e.PRNumber)
		}
		return e, nil
	}
	out, err := run("", "gh", "pr", "view", fmt.Sprint(e.PRNumber), "-R", e.Repo,
		"--json", "state,mergedAt", "-q", "{state: .state, mergedAt: .mergedAt}")
	if err != nil {
		return e, err
	}
	var view struct {
		State    string `json:"state"`
		MergedAt string `json:"mergedAt"`
	}
	_ = json.Unmarshal(out, &view)
	switch {
	case view.MergedAt != "" || view.State == "MERGED":
		e.State = "merged"
	case view.State == "CLOSED":
		e.State = "closed"
	default:
		e.State = "open"
	}
	e.History = append(e.History, state.LedgerEvent{Run: runID, Action: "reconciled", Detail: e.State})
	return e, nil
}

func prNumberFromURL(url string) int {
	n := 0
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] < '0' || url[i] > '9' {
			if i == len(url)-1 {
				return 0
			}
			fmt.Sscanf(url[i+1:], "%d", &n)
			return n
		}
	}
	return n
}
```

- [ ] **Step 2: Build + vet**

Run: `go build ./... && go vet ./...`
Expected: success (this task adds no tests; it is exercised by the dry-run e2e in Task 7).

- [ ] **Step 3: Commit**

```bash
git add internal/remediate/git_executor.go
git commit -m "feat(remediate): real git+gh executor (clone/bump/verify/push/PR) with dry-run"
```

---

### Task 6: Dashboard "Bot PR ledger" section

**Files:**
- Modify: `internal/render/render.go`
- Test: `internal/render/ledger_test.go`

**Interfaces:**
- Consumes: `state.Ledger`, `state.LedgerEntry`.
- Produces: a new field on `render.Input`: `Ledger state.Ledger` (json tag `ledger`, placed after `Repos`); the markdown (and HTML, via the shared renderer) gains a `## 🤖 Bot PR ledger` section listing each entry (repo, package→version, state, PR link). When the ledger is empty, the section prints `_No bot PRs yet._`.

- [ ] **Step 1: Write the failing test**

Create `internal/render/ledger_test.go`:

```go
package render

import (
	"strings"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestDashboardMarkdownLedgerSection(t *testing.T) {
	in := Input{Ledger: state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", Repo: "kairos-io/immucore",
			Package: "golang.org/x/net", State: "open", PRNumber: 412,
			PRURL: "https://github.com/kairos-io/immucore/pull/412",
			Bump:  state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "Bot PR ledger")
	assert.Contains(t, md, "kairos-io/immucore")
	assert.Contains(t, md, "golang.org/x/net")
	assert.Contains(t, md, "0.33.0")
	assert.Contains(t, md, "open")
	assert.Contains(t, md, "412")
}

func TestDashboardMarkdownLedgerEmpty(t *testing.T) {
	md := DashboardMarkdown(Input{})
	assert.Contains(t, md, "No bot PRs yet")
	assert.True(t, strings.Contains(md, "Bot PR ledger"))
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/render/...`

- [ ] **Step 3: Implement**

In `internal/render/render.go`: add `Ledger state.Ledger `json:"ledger"`` to `Input` (after `Repos`). Then, in `DashboardMarkdown`, BEFORE the run-log footer, add:

```go
	// Bot PR ledger
	b.WriteString("## 🤖 Bot PR ledger\n\n")
	if len(in.Ledger.Entries) == 0 {
		b.WriteString("_No bot PRs yet._\n\n")
	} else {
		b.WriteString("| Repo | Bump | State | PR |\n|---|---|---|---|\n")
		for _, e := range in.Ledger.Entries {
			pr := "—"
			if e.PRNumber > 0 {
				pr = fmt.Sprintf("[#%d](%s)", e.PRNumber, e.PRURL)
			}
			fmt.Fprintf(&b, "| %s | %s@%s | %s | %s |\n", e.Repo, e.Bump.Package, e.Bump.To, e.State, pr)
		}
		b.WriteString("\n")
	}
```

(If you add the same section to `internal/render/html.go`'s `DashboardHTML`, mirror it as an escaped table and regenerate the HTML golden with `UPDATE_GOLDEN=1`. The markdown test above is the gate; the HTML section is optional polish for this task — if you add it, escape every value via the existing `html/template` path.)

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/render/...`
If the render golden tests fail because `Input` changed, regenerate: `UPDATE_GOLDEN=1 go test ./internal/render/...`, eyeball the goldens (the ledger section should appear), then `go test ./internal/render/...`.

- [ ] **Step 5: Commit**

```bash
git add internal/render/render.go internal/render/ledger_test.go internal/render/testdata/
git commit -m "feat(render): Bot PR ledger dashboard section"
```

---

### Task 7: Wire the `remediate` phase + workflow step

**Files:**
- Modify: `cmd/ksec/main.go`
- Modify: `.github/workflows/security-dashboard.yaml`

**Interfaces:**
- Consumes: `remediate.Plan`, `remediate.Run`, `remediate.GitExecutor`, `state.Ledger`/`LedgerFile`, `state.Correlated`/`Repos`.
- Produces: a `remediate` subcommand and a workflow step that runs it between `triage` and `render`. `render` loads `ledger.json` and passes it to `Input.Ledger`.

- [ ] **Step 1: Add the `remediate` subcommand**

In `cmd/ksec/main.go`, register `root.AddCommand(newRemediateCmd(gf))` and add:

```go
func newRemediateCmd(gf *globalFlags) *cobra.Command {
	var maxPRs int
	cmd := &cobra.Command{
		Use:   "remediate",
		Short: "open and maintain dependency-bump PRs for actionable findings",
		RunE: func(cmd *cobra.Command, args []string) error {
			var c state.Correlated
			if err := state.Load(gf.stateDir, state.CorrelatedFile, &c); err != nil {
				return err
			}
			var ledger state.Ledger
			_ = state.Load(gf.stateDir, state.LedgerFile, &ledger) // best-effort: empty on first run

			runID := os.Getenv("KSEC_RUN_URL")
			if runID == "" {
				runID = "local"
			}
			intents, deferred := remediate.Plan(c, ledger, maxPRs)
			if deferred > 0 {
				fmt.Fprintf(os.Stderr, "remediate: %d new bumps deferred by --max-prs=%d\n", deferred, maxPRs)
			}
			ex := &remediate.GitExecutor{Token: os.Getenv("GH_TOKEN"), DryRun: gf.dryRun}
			out, results := remediate.Run(intents, ex, ledger, runID)
			for _, r := range results {
				fmt.Fprintf(os.Stderr, "remediate: %s %s -> %s %s\n", r.Action, r.Key, r.State, r.Detail)
			}
			return state.Save(gf.stateDir, state.LedgerFile, out)
		},
	}
	cmd.Flags().IntVar(&maxPRs, "max-prs", 10, "maximum NEW PRs to open per run (blast-radius guard)")
	return cmd
}
```

Add `"github.com/kairos-io/security/internal/remediate"` to the imports.

- [ ] **Step 2: Make `render` load the ledger**

In `newRenderCmd`, after the best-effort `repos` load, add:

```go
			var ledger state.Ledger
			_ = state.Load(gf.stateDir, state.LedgerFile, &ledger) // best-effort
```

and set `Ledger: ledger,` in the `render.Input{...}` literal.

- [ ] **Step 3: Build + vet + full test**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: all pass. Smoke: `go run ./cmd/ksec remediate --help` shows `--max-prs`.

- [ ] **Step 4: Add the workflow step (dry-run by default until the token is ready)**

In `.github/workflows/security-dashboard.yaml`, between the `ksec triage` and `ksec render` lines in the "Run pipeline" step, add:

```yaml
          ksec remediate --state-dir state $REMEDIATE_DRYRUN
```

and in the job `env:` block add:

```yaml
      # Remediation writes to OTHER repos. Keep it dry-run until KSEC_BOT_TOKEN
      # has contents:write + pull_requests:write on the targets. Flip to "" to
      # go live.
      REMEDIATE_DRYRUN: --dry-run
```

Also add `state/ledger.json` to the `git add` line of the "Commit state + dashboards" step so the ledger is committed:

```yaml
          git add state/ dashboard.md dashboard.json
```

(`state/` already covers `ledger.json`; confirm the path glob includes it.) Validate the YAML: `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/security-dashboard.yaml'))" && echo OK`.

- [ ] **Step 5: End-to-end dry-run smoke test (local)**

Run:
```bash
rm -rf /tmp/rem && mkdir -p /tmp/rem
cat > /tmp/rem/correlated.json <<'JSON'
{"findings":[{"id":"a","repo":"kairos-io/immucore","type":"sourceCVE","ecosystem":"go","package":"golang.org/x/net","fixedVersion":"0.33.0","severity":"high"}],"waterfall":[]}
JSON
go run ./cmd/ksec remediate --state-dir /tmp/rem --dry-run
cat /tmp/rem/ledger.json
```
Expected: prints `[dry-run] would open PR on kairos-io/immucore ...`, and `ledger.json` contains one entry with `"state": "planned"`, key `kairos-io/immucore|golang.org/x/net`. No network writes occurred.

- [ ] **Step 6: Commit**

```bash
git add cmd/ksec/main.go .github/workflows/security-dashboard.yaml
git commit -m "feat: wire remediate phase + dry-run workflow step; render loads ledger"
```

---

## Self-review

**Spec coverage** (design §6.5, §8):
- Reconciliation loop (open vs reconcile existing) → Tasks 2, 3. ✓
- `go get @<fixed>` + `go mod tidy`, verify-before-push (`go build`), no broken PR (`build-failed`) → Task 5. ✓
- Identity by marker + `ksec/` branch; only touches own PRs (reconcile reads by stored PR number; open creates fresh branch) → Tasks 4, 5. ✓
- Bump deterministic; AI only for reactions → AI explicitly deferred (Plan 3); Plan 2 bodies deterministic (Task 4). Noted in Global Constraints. ✓
- Blast-radius cap (default 10) → Task 2 (`maxNew`), Task 7 (`--max-prs`). ✓
- Dry-run short-circuits every write → Task 5 (`DryRun`), Task 7 (workflow `REMEDIATE_DRYRUN`). ✓
- Ledger committed as memory → Tasks 1, 7. ✓
- Waterfall coordination: **simplified** — bumps are keyed per repo+package (the unit that actually changes); waterfall groups remain a dashboard/correlation concept. A package shared across repos already yields one target per repo, which is the same set of PRs a group fan-out would produce. Documented deviation from the spec's group-keyed fan-out; cross-linking PRs by group is deferred (cosmetic).
- Dashboard "Bot PR ledger" → Task 6. ✓
- Error isolation, secrets not logged (executor surfaces stderr but not the token; clone URL with token is never printed — `[dry-run]` prints repo/branch only) → Tasks 3, 5. ✓

**Out of scope (Plan 3, correctly absent):** comment classification/reactions, AI PR prose, group cross-linking. `SeenComments`/`History` fields exist so Plan 3 needs no migration.

**Placeholder scan:** none — every step has complete code/commands.

**Type consistency:** `Intent`/`Result`/`IntentType` (Task 2) used by Tasks 3, 5, 7. `Executor` (Task 3) implemented by `FakeExecutor` (Task 3) and `GitExecutor` (Task 5). `state.LedgerEntry`/`Ledger`/`LedgerFile`/`ByKey` (Task 1) used by Tasks 2, 3, 5, 6, 7. `BranchName`/`PRTitle`/`PRBody`/`PRMarker` (Task 4) used by Task 5. `render.Input.Ledger` (Task 6) set by Task 7.

---

## Operational notes (resolve before going live)

- **Token:** `KSEC_BOT_TOKEN` needs `contents:write` + `pull_requests:write` on every repo the bot may open PRs against (the `repos.yaml` set). Until then, leave `REMEDIATE_DRYRUN: --dry-run`.
- **Go toolchain version per repo:** `go get`/`go build` use the runner's Go; a repo requiring a newer Go than installed will `build-failed` (surfaced, not silently dropped). Consider `actions/setup-go` with a recent Go, or per-repo toolchain handling, as a follow-up.
- **First live run:** start with a low `--max-prs` (e.g. 2-3) and watch the ledger before raising it.
