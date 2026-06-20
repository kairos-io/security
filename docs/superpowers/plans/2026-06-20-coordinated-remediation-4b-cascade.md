# Coordinated Remediation 4b — First-Party Cascade Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Propagate a security fix down the first-party dependency graph: when a fix merges into a kairos-io module (e.g. `kairos-sdk`), immediately open **pseudo-version** bump PRs in every repo that consumes it (don't wait for a tag), recursing to consumers-of-consumers, and **re-pin** each to a real tag once the module cuts a release — plus fold in two Plan-4a follow-ups (detect our own PRs by branch, not login; match existing PRs by package **and** version).

**Architecture:** A new pure `depgraph` (built from each tracked repo's `go.mod`) maps module→repo and consumers(module). The pure `Plan` gains the graph and emits `cascade` intents (bump a first-party module to its default-branch pseudo-version in each consumer of a merged fix) and `repin` intents (re-pin a pseudo cascade to a published tag). The `Executor` gains `Cascade` (clone/`go get @branch`/verify/push/PR) and `Repin` (tag-check/`go get @tag`/force-push); the ledger gains `cascadeFrom`/`pinTarget`/`pseudo`. Build-break repair is deferred to 4c — a cascade that breaks `go build` records `build-failed`/`needsHuman`, never a broken push. Builds on Plan 2/3/4a.

**Tech Stack:** Go 1.22, `gh` + `git` + `go` CLIs, `stretchr/testify`, existing `internal/remediate`, `internal/state`, `internal/ghclient`, `internal/discover` (regex go.mod parsing pattern).

## Global Constraints

- Module `github.com/kairos-io/security`; binary `ksec`; Go 1.22.
- **Cascade triggers on availability, not approval:** when a first-party module's fix is on its default branch (an upstream ledger entry in that module repo reaches `merged`), open consumer bump PRs immediately to the module's **default-branch pseudo-version** (`go get <module>@<defaultBranch>`). The PR body notes it's unreleased and asks a maintainer to tag. Do NOT wait for a tag or human approval.
- **Re-pin follow-up:** each run, for `pseudo` cascade entries, if the module has published a tag, re-pin (`go get <module>@<tag>`, force-push the same branch, set `pseudo=false`, `pinTarget=<tag>`). Until then, record `awaiting-release` (no error).
- **Recursion:** a merged cascade bump is itself an upstream fix — its consumers cascade in turn. Termination is guaranteed by the acyclic Go module graph + the "cascade key already exists" skip.
- **Cascade key:** `"<consumerRepo>|<moduleImportPath>"` (the existing `key()` format, `Package` = the module import path).
- **Blast-radius cap:** `--max-prs` applies to **new-PR-creating intents collectively** — direct `open` + `cascade`. `repin` (force-pushes an existing branch) and `adopt` are not capped.
- **Verify-before-push** on cascade and re-pin (`go build ./...`); a non-building tree is never pushed (records `build-failed`).
- Dry-run short-circuits every git/GitHub write; token never logged (existing `g.run` redaction).
- **Follow-up A:** a PR is "ours" iff its head branch starts with `ksec/` OR its author is `kairos-security-bot`; classify such PRs `ksec` (never adopt them).
- **Follow-up B:** `MatchPR` requires the PR title to contain the package path **and** the target version (normalized, leading `v` stripped); empty version disables the version requirement.

---

## File structure

```
internal/state/types.go              # + LedgerEntry.CascadeFrom/PinTarget/Pseudo (modify)
internal/remediate/depgraph.go       # BuildGraph + DepGraph methods (create)
internal/remediate/depgraph_test.go  # (create)
internal/ghclient/ghclient.go        # PullRequest.HeadRef + ListOpenPRs query (modify)
internal/remediate/matcher.go        # isOwnPR (branch-based); MatchPR(pkg, version, prs) (modify)
internal/remediate/matcher_test.go   # (modify)
internal/remediate/intent.go         # IntentCascade/IntentRepin + Intent.Ref/CascadeFrom (modify)
internal/remediate/planner.go        # Plan gains graph; cascade+repin intents; combined cap (modify)
internal/remediate/planner_test.go   # (modify)
internal/remediate/run.go            # Executor.Cascade/Repin + Run cases (modify)
internal/remediate/fake.go           # FakeExecutor.Cascade/Repin (modify)
internal/remediate/run_test.go       # (modify)
internal/remediate/prbody.go         # CascadeBranchName + CascadePRBody (modify)
internal/remediate/prbody_test.go    # (modify)
internal/remediate/git_executor.go   # GitExecutor.Cascade + Repin (modify)
internal/render/render.go            # ledger: show pseudo + cascadeFrom (modify)
internal/render/coord_test.go        # (modify)
cmd/ksec/main.go                     # build depgraph (fetch go.mods), pass to Plan (modify)
```

---

### Task 1: Ledger cascade fields

**Files:** Modify `internal/state/types.go`; Test `internal/state/ledger_test.go`.

**Interfaces:** `LedgerEntry` gains `CascadeFrom string`, `PinTarget string`, `Pseudo bool`.

- [ ] **Step 1: Add the fields** — in `LedgerEntry`, after `NeedsHuman`:

```go
	CascadeFrom string `json:"cascadeFrom,omitempty"` // upstream ledger key that triggered this cascade bump
	PinTarget   string `json:"pinTarget,omitempty"`   // for a pseudo cascade: the tag to re-pin to ("" while still pseudo)
	Pseudo      bool   `json:"pseudo,omitempty"`      // true while the bump points at a pseudo-version awaiting re-pin
```

- [ ] **Step 2: Extend the round-trip test** — set `CascadeFrom: "kairos-io/kairos-sdk|golang.org/x/net"`, `PinTarget: "v0.8.1"`, `Pseudo: true` on the entry in `TestLedgerRoundTrip`; the existing `assert.Equal(t, in, out)` covers them.

- [ ] **Step 3: Run + commit**

Run: `go test ./internal/state/...` (PASS).
```bash
git add internal/state/types.go internal/state/ledger_test.go
git commit -m "feat(state): ledger cascade fields (cascadeFrom/pinTarget/pseudo)"
```

---

### Task 2: `depgraph` — first-party dependency graph

**Files:** Create `internal/remediate/depgraph.go`, `internal/remediate/depgraph_test.go`.

**Interfaces:**
- Consumes: `state.Repo`.
- Produces:
  - `type DepGraph struct { ... }`
  - `func BuildGraph(repos []state.Repo, gomodByRepo map[string][]byte) *DepGraph`
  - `func (g *DepGraph) ModuleOf(repo string) string` — the repo's module import path (or "").
  - `func (g *DepGraph) RepoOf(module string) (string, bool)` — the tracked repo providing that module.
  - `func (g *DepGraph) Consumers(module string) []string` — tracked repos whose go.mod requires that module (sorted).
  - `func (g *DepGraph) BranchOf(repo string) string` — the repo's default branch (from `repos`, default "main").

- [ ] **Step 1: Write the failing test**

Create `internal/remediate/depgraph_test.go`:

```go
package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildGraph(t *testing.T) {
	repos := []state.Repo{
		{Repo: "kairos-io/kairos-sdk", Branch: "main"},
		{Repo: "kairos-io/immucore", Branch: "master"},
		{Repo: "kairos-io/kairos-agent"},
	}
	gomod := map[string][]byte{
		"kairos-io/kairos-sdk": []byte("module github.com/kairos-io/kairos-sdk\ngo 1.22\n"),
		"kairos-io/immucore": []byte("module github.com/kairos-io/immucore\nrequire github.com/kairos-io/kairos-sdk v0.7.0\n"),
		"kairos-io/kairos-agent": []byte("module github.com/kairos-io/kairos-agent\nrequire (\n\tgithub.com/kairos-io/kairos-sdk v0.7.0\n\tgithub.com/kairos-io/immucore v0.5.0\n)\n"),
	}
	g := BuildGraph(repos, gomod)

	assert.Equal(t, "github.com/kairos-io/kairos-sdk", g.ModuleOf("kairos-io/kairos-sdk"))
	r, ok := g.RepoOf("github.com/kairos-io/kairos-sdk")
	require.True(t, ok)
	assert.Equal(t, "kairos-io/kairos-sdk", r)

	// consumers of the sdk module: immucore and kairos-agent, sorted
	assert.Equal(t, []string{"kairos-io/immucore", "kairos-io/kairos-agent"}, g.Consumers("github.com/kairos-io/kairos-sdk"))
	// consumers of immucore: kairos-agent
	assert.Equal(t, []string{"kairos-io/kairos-agent"}, g.Consumers("github.com/kairos-io/immucore"))
	// branch lookups
	assert.Equal(t, "master", g.BranchOf("kairos-io/immucore"))
	assert.Equal(t, "main", g.BranchOf("kairos-io/kairos-agent")) // default
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement**

Create `internal/remediate/depgraph.go`:

```go
package remediate

import (
	"regexp"
	"sort"

	"github.com/kairos-io/security/internal/state"
)

type DepGraph struct {
	moduleOf  map[string]string   // repo -> module import path
	repoOf    map[string]string   // module import path -> repo
	consumers map[string][]string // module import path -> repos requiring it
	branchOf  map[string]string   // repo -> default branch
}

var (
	reModuleLine  = regexp.MustCompile(`(?m)^module\s+(\S+)`)
	reRequireLine = regexp.MustCompile(`(?m)^\s*(?:require\s+)?(github\.com/\S+)\s+v\S+`)
)

// BuildGraph parses each tracked repo's go.mod to map module<->repo and to find,
// for each first-party module, the tracked repos that require it.
func BuildGraph(repos []state.Repo, gomodByRepo map[string][]byte) *DepGraph {
	g := &DepGraph{
		moduleOf:  map[string]string{},
		repoOf:    map[string]string{},
		consumers: map[string][]string{},
		branchOf:  map[string]string{},
	}
	for _, r := range repos {
		b := r.Branch
		if b == "" {
			b = "main"
		}
		g.branchOf[r.Repo] = b
		if m := reModuleLine.FindSubmatch(gomodByRepo[r.Repo]); m != nil {
			mod := string(m[1])
			g.moduleOf[r.Repo] = mod
			g.repoOf[mod] = r.Repo
		}
	}
	// Second pass: requires, keeping only first-party modules (those we map to a repo).
	for repo, mod := range g.moduleOf {
		_ = mod
		for _, m := range reRequireLine.FindAllSubmatch(gomodByRepo[repo], -1) {
			req := string(m[1])
			if _, ok := g.repoOf[req]; ok {
				g.consumers[req] = append(g.consumers[req], repo)
			}
		}
	}
	for k := range g.consumers {
		sort.Strings(g.consumers[k])
	}
	return g
}

func (g *DepGraph) ModuleOf(repo string) string { return g.moduleOf[repo] }
func (g *DepGraph) RepoOf(module string) (string, bool) {
	r, ok := g.repoOf[module]
	return r, ok
}
func (g *DepGraph) Consumers(module string) []string { return g.consumers[module] }
func (g *DepGraph) BranchOf(repo string) string {
	if b, ok := g.branchOf[repo]; ok {
		return b
	}
	return "main"
}
```

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/remediate/...`

- [ ] **Step 5: Commit**

```bash
git add internal/remediate/depgraph.go internal/remediate/depgraph_test.go
git commit -m "feat(remediate): first-party dependency graph from go.mod"
```

---

### Task 3: Follow-ups — own-PR by branch + version-aware matching

**Files:** Modify `internal/ghclient/ghclient.go`, `internal/remediate/matcher.go`, `internal/remediate/matcher_test.go`.

**Interfaces:**
- `ghclient.PullRequest` gains `HeadRef string` (json `headRefName`); `CLI.ListOpenPRs` selects it.
- `func isOwnPR(pr ghclient.PullRequest) bool` — `strings.HasPrefix(pr.HeadRef, "ksec/") || pr.Author == "kairos-security-bot"`.
- `classifySource` returns `ksec` for own PRs (branch or author).
- `MatchPR(pkg, version string, prs []ghclient.PullRequest) (ghclient.PullRequest, string, bool)` — title must contain `pkg` AND (if `version != ""`) the version with a leading `v` stripped.

- [ ] **Step 1: Add `HeadRef` to ghclient**

In `internal/ghclient/ghclient.go`: add `HeadRef string \`json:"headRefName"\`` to `PullRequest`; in `CLI.ListOpenPRs`, add `headRefName` to the `--json` list and `headRefName` to the jq projection (`headRef: .headRefName`)... actually keep field name aligned: change the `-q` map to include `headRef: .headRefName` and the struct tag to `json:"headRef"`. Use:

```go
type PullRequest struct {
	Number  int      `json:"number"`
	Title   string   `json:"title"`
	Author  string   `json:"author"`
	URL     string   `json:"url"`
	HeadRef string   `json:"headRef"`
	Labels  []string `json:"labels"`
}
```

and in `ListOpenPRs` the gh call:

```go
	b, err := c.run("pr", "list", "-R", repo, "--state", "open", "--limit", "200",
		"--json", "number,title,author,url,headRefName,labels",
		"-q", "[.[] | {number, title, author: .author.login, url, headRef: .headRefName, labels: [.labels[].name]}]")
```

- [ ] **Step 2: Update the matcher test**

In `internal/remediate/matcher_test.go`: every `MatchPR(...)` call gains a version arg; add cases:

```go
func TestMatchPRRequiresVersion(t *testing.T) {
	prs := []ghclient.PullRequest{
		{Number: 1, Title: "Bump golang.org/x/net from 0.30.0 to 0.33.0", Author: "dependabot[bot]"},
	}
	// pkg + matching version -> match
	pr, src, ok := MatchPR("golang.org/x/net", "0.33.0", prs)
	assert.True(t, ok)
	assert.Equal(t, 1, pr.Number)
	assert.Equal(t, "dependabot", src)
	// pkg present but version absent -> no match (avoids "remove x/net usage" false positives)
	_, _, ok = MatchPR("golang.org/x/net", "9.9.9", prs)
	assert.False(t, ok)
	// empty version disables the version requirement
	_, _, ok = MatchPR("golang.org/x/net", "", prs)
	assert.True(t, ok)
}

func TestIsOwnPRByBranch(t *testing.T) {
	_, src, ok := MatchPR("golang.org/x/net", "", []ghclient.PullRequest{
		{Number: 2, Title: "bump golang.org/x/net", Author: "someoneelse", HeadRef: "ksec/bump-golang-org-x-net-0-33-0"},
	})
	assert.True(t, ok)
	assert.Equal(t, "ksec", src, "a ksec/ branch is ours regardless of author")
}
```

Also update the existing `MatchPR(...)` call in `TestMatchPR` to pass `""` for version.

- [ ] **Step 3: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 4: Implement**

Rewrite `internal/remediate/matcher.go`:

```go
package remediate

import (
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
)

func isOwnPR(pr ghclient.PullRequest) bool {
	return strings.HasPrefix(pr.HeadRef, "ksec/") || pr.Author == "kairos-security-bot"
}

func classifySource(pr ghclient.PullRequest) string {
	if isOwnPR(pr) {
		return "ksec"
	}
	switch pr.Author {
	case "renovate[bot]":
		return "renovate"
	case "dependabot[bot]":
		return "dependabot"
	default:
		return "human"
	}
}

// MatchPR returns the first open PR whose title contains the package path
// (case-insensitive) and, when version != "", the version (leading 'v'
// stripped). Requiring the version avoids matching PRs that merely mention the
// package (e.g. "remove golang.org/x/net usage").
func MatchPR(pkg, version string, prs []ghclient.PullRequest) (ghclient.PullRequest, string, bool) {
	if pkg == "" {
		return ghclient.PullRequest{}, "", false
	}
	pkgL := strings.ToLower(pkg)
	verL := strings.ToLower(strings.TrimPrefix(version, "v"))
	for _, pr := range prs {
		title := strings.ToLower(pr.Title)
		if !strings.Contains(title, pkgL) {
			continue
		}
		if verL != "" && !strings.Contains(title, verL) {
			continue
		}
		return pr, classifySource(pr), true
	}
	return ghclient.PullRequest{}, "", false
}
```

- [ ] **Step 5: Run it — expect FAIL** (the planner still calls the old `MatchPR(pkg, prs)`). Update the planner call site in `planner.go` to `MatchPR(t.pkg, t.to, prsByRepo[t.repo])`. Re-run `go test ./internal/remediate/...` — expect PASS.

- [ ] **Step 6: Build the whole module + commit**

Run: `go build ./... && go test ./...`
```bash
git add internal/ghclient/ghclient.go internal/remediate/matcher.go internal/remediate/matcher_test.go internal/remediate/planner.go
git commit -m "fix(remediate): own-PR by branch + version-aware PR matching (4a follow-ups)"
```

---

### Task 4: Cascade + repin intents (planner)

**Files:** Modify `internal/remediate/intent.go`, `internal/remediate/planner.go`, `internal/remediate/planner_test.go`.

**Interfaces:**
- `IntentCascade IntentType = "cascade"`, `IntentRepin IntentType = "repin"`; `Intent` gains `Ref string` (the module's default branch for the pseudo `go get`) and `CascadeFrom string`.
- `Plan(c state.Correlated, ledger state.Ledger, prsByRepo map[string][]ghclient.PullRequest, graph *DepGraph, maxNew int) ([]Intent, int)` — adds cascade + repin. The `--max-prs` cap is shared across `open` + `cascade`.

- [ ] **Step 1: Update intent types**

In `intent.go`: add `IntentCascade IntentType = "cascade"` and `IntentRepin IntentType = "repin"` to the const block; add `Ref string` and `CascadeFrom string` to `Intent`.

- [ ] **Step 2: Update + extend the planner test**

In `planner_test.go`: every `Plan(...)` call gains a `graph` argument (pass `nil` for the existing direct/adopt tests — guard nil in the implementation so cascade is skipped when graph is nil). Add:

```go
func TestPlanCascadesMergedFirstPartyFix(t *testing.T) {
	repos := []state.Repo{
		{Repo: "kairos-io/kairos-sdk", Branch: "main"},
		{Repo: "kairos-io/immucore", Branch: "master"},
	}
	gomod := map[string][]byte{
		"kairos-io/kairos-sdk": []byte("module github.com/kairos-io/kairos-sdk\n"),
		"kairos-io/immucore":   []byte("module github.com/kairos-io/immucore\nrequire github.com/kairos-io/kairos-sdk v0.7.0\n"),
	}
	g := BuildGraph(repos, gomod)
	// A merged fix in the sdk repo -> cascade a pseudo bump into immucore.
	ledger := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/kairos-sdk|golang.org/x/net", Repo: "kairos-io/kairos-sdk", State: "merged",
			Kind: "direct", Severity: "high", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	intents, _ := Plan(state.Correlated{}, ledger, nil, g, 10)

	var cas *Intent
	for i := range intents {
		if intents[i].Type == IntentCascade {
			cas = &intents[i]
		}
	}
	require.NotNil(t, cas, "expected a cascade intent")
	assert.Equal(t, "kairos-io/immucore", cas.Repo)
	assert.Equal(t, "github.com/kairos-io/kairos-sdk", cas.Package)
	assert.Equal(t, "main", cas.Ref) // sdk's default branch for the pseudo go get
	assert.Equal(t, "kairos-io/kairos-sdk|golang.org/x/net", cas.CascadeFrom)
}

func TestPlanRepinsPseudoCascade(t *testing.T) {
	ledger := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|github.com/kairos-io/kairos-sdk", Repo: "kairos-io/immucore",
			Package: "github.com/kairos-io/kairos-sdk", State: "open", Kind: "cascade", Pseudo: true},
	}}
	intents, _ := Plan(state.Correlated{}, ledger, nil, nil, 10)
	var found bool
	for _, in := range intents {
		if in.Type == IntentRepin && in.Key == "kairos-io/immucore|github.com/kairos-io/kairos-sdk" {
			found = true
		}
	}
	assert.True(t, found, "expected a repin intent for the pseudo cascade entry")
}
```

- [ ] **Step 3: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 4: Implement the cascade + repin logic**

In `planner.go`, change the signature to add `graph *DepGraph`, and BEFORE the final return add cascade + repin emission, and fold cascade opens into the cap. Insert after the reconcile loop and after computing `openKeys` (replace the open-emission tail). The complete new tail:

```go
	// Repin: every pseudo cascade entry is a repin candidate (the executor
	// decides whether a tag is available yet).
	for i := range ledger.Entries {
		e := &ledger.Entries[i]
		if e.Kind == "cascade" && e.Pseudo && e.State == "open" {
			intents = append(intents, Intent{Type: IntentRepin, Key: e.Key, Repo: e.Repo, Entry: e})
		}
	}

	// Cascade: a merged fix in a first-party module repo means that module's
	// default branch has the fix; bump it in each consumer that isn't already
	// tracked for it. Cascade PRs share the maxNew cap with direct opens.
	type newPR struct {
		intent Intent
		sev    string
	}
	var pool []newPR
	for _, k := range openKeys { // direct gaps from the earlier 4a logic
		t := targets[k]
		pool = append(pool, newPR{
			intent: Intent{Type: IntentOpen, Key: k, Repo: t.repo, Package: t.pkg, Severity: t.sev,
				Bump: state.Bump{Package: t.pkg, To: t.to}},
			sev: t.sev,
		})
	}
	if graph != nil {
		for i := range ledger.Entries {
			e := &ledger.Entries[i]
			mod := graph.ModuleOf(e.Repo)
			if mod == "" || e.State != "merged" {
				continue
			}
			for _, consumer := range graph.Consumers(mod) {
				ck := key(consumer, mod)
				if ce, ok := ledger.ByKey(ck); ok {
					if ce.State == "open" || ce.State == "conflicted" || ce.State == "merged" {
						continue // already cascading / done
					}
				}
				pool = append(pool, newPR{
					intent: Intent{Type: IntentCascade, Key: ck, Repo: consumer, Package: mod,
						Ref: graph.BranchOf(e.Repo), CascadeFrom: e.Key, Severity: e.Severity},
					sev: e.Severity,
				})
			}
		}
	}

	sort.SliceStable(pool, func(i, j int) bool {
		if sevRank[pool[i].sev] != sevRank[pool[j].sev] {
			return sevRank[pool[i].sev] > sevRank[pool[j].sev]
		}
		return pool[i].intent.Key < pool[j].intent.Key
	})
	deferred := 0
	for n := range pool {
		if n >= maxNew {
			deferred = len(pool) - n
			break
		}
		intents = append(intents, pool[n].intent)
	}
	return intents, deferred
}
```

Remove the OLD `openKeys` sort + emit + `deferred` tail (the `keys`/`deferred` loop from 4a that emitted `IntentOpen` and computed `deferred`) — it is replaced by the `pool` above. Keep the earlier code that POPULATES `openKeys` (the adopt-vs-gap decision). The adopt emission (IntentAdopt) stays as-is, before this tail.

- [ ] **Step 5: Run it — expect PASS.** Run: `go test ./internal/remediate/...`
Also fix the `main.go` `Plan(...)` call to pass `nil` for `graph` temporarily (Task 9 wires the real graph): `remediate.Plan(c, ledger, prsByRepo, nil, maxPRs)` with `// TODO(plan-4b task 9): pass real depgraph`.

- [ ] **Step 6: Build + commit**

Run: `go build ./... && go test ./...`
```bash
git add internal/remediate/intent.go internal/remediate/planner.go internal/remediate/planner_test.go cmd/ksec/main.go
git commit -m "feat(remediate): planner cascade + repin intents (shared cap)"
```

---

### Task 5: Run loop + Executor.Cascade/Repin + fakes

**Files:** Modify `internal/remediate/run.go`, `internal/remediate/fake.go`, `internal/remediate/run_test.go`.

**Interfaces:** `Executor` gains `Cascade(in Intent, run string) (state.LedgerEntry, error)` and `Repin(e state.LedgerEntry, run string) (state.LedgerEntry, error)`; `Run` handles `IntentCascade` (like `Open`, error-isolated) and `IntentRepin` (like `Reconcile`, takes `*in.Entry`). `FakeExecutor` gains `Cascaded`/`Repinned` maps + methods. A temporary `GitExecutor.Cascade`/`Repin` stub keeps the build green until Tasks 7-8 (the `var _ Executor = (*GitExecutor)(nil)` assertion requires them).

- [ ] **Step 1: Update the failing test** — add to `run_test.go`:

```go
func TestRunCascadeAndRepin(t *testing.T) {
	entry := state.LedgerEntry{Key: "c|m", Repo: "c", Package: "m", State: "open", Kind: "cascade", Pseudo: true}
	intents := []Intent{
		{Type: IntentCascade, Key: "c|m", Repo: "c", Package: "m", Ref: "main", CascadeFrom: "u|x"},
		{Type: IntentRepin, Key: "c|m", Entry: &entry},
	}
	fake := &FakeExecutor{
		Cascaded: map[string]state.LedgerEntry{"c|m": {Key: "c|m", Repo: "c", Package: "m", State: "open", Kind: "cascade", Pseudo: true}},
		Repinned: map[string]state.LedgerEntry{"c|m": {Key: "c|m", Repo: "c", Package: "m", State: "open", Kind: "cascade", Pseudo: false, PinTarget: "v1.0.0"}},
	}
	out, results := Run(intents, fake, state.Ledger{}, "2026-06-20")
	require.Len(t, out.Entries, 1)
	// repin ran after cascade (same key) -> pseudo cleared, pin set
	assert.False(t, out.Entries[0].Pseudo)
	assert.Equal(t, "v1.0.0", out.Entries[0].PinTarget)
	require.Len(t, results, 2)
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement** — add the two `Executor` methods, two `Run` switch cases (cascade mirrors `IntentOpen` with `Kind:"cascade"`/`CascadeFrom`/`Pseudo:true` on the error entry; repin mirrors `IntentReconcile` using `*in.Entry`), `FakeExecutor.Cascaded`/`Repinned` + methods (return the mapped entry or a sensible default), and a temporary `GitExecutor.Cascade`/`Repin` stub each returning `fmt.Errorf("... not implemented")` with a `// TODO(plan-4b tasks 7-8)` comment.

(Mirror the exact structure of the existing `IntentOpen`/`IntentReconcile` cases and the Task-6/Plan-4a `Adopt` stub pattern. The cascade error entry: `state.LedgerEntry{Key: in.Key, Repo: in.Repo, Package: in.Package, State: "error", Kind: "cascade", CascadeFrom: in.CascadeFrom, Pseudo: true, Severity: in.Severity, CreatedRun: run, LastActionRun: run, History: []state.LedgerEvent{{Run: run, Action: "cascade-failed", Detail: err.Error()}}}`.)

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/remediate/...` then `go build ./...`

- [ ] **Step 5: Commit**

```bash
git add internal/remediate/run.go internal/remediate/fake.go internal/remediate/run_test.go internal/remediate/git_executor.go
git commit -m "feat(remediate): Run handles cascade/repin; Executor methods + fakes (stubs)"
```

---

### Task 6: Cascade/repin PR body

**Files:** Modify `internal/remediate/prbody.go`, `internal/remediate/prbody_test.go`.

**Interfaces:**
- `func CascadeBranchName(in Intent) string` — `"ksec/cascade-" + slug(module) + "-pseudo"`.
- `func CascadePRBody(in Intent) string` — explains this consumes an unreleased fix in the upstream module (`in.Package`) via a pseudo-version, asks a maintainer to tag a release, and ends with `PRMarker(in.Key)` as the last line.

- [ ] **Step 1: Write the failing test** — `prbody_test.go`:

```go
func TestCascadePRBodyAndBranch(t *testing.T) {
	in := Intent{Type: IntentCascade, Key: "kairos-io/immucore|github.com/kairos-io/kairos-sdk",
		Repo: "kairos-io/immucore", Package: "github.com/kairos-io/kairos-sdk", CascadeFrom: "kairos-io/kairos-sdk|x", Severity: "high"}
	assert.Equal(t, "ksec/cascade-github-com-kairos-io-kairos-sdk-pseudo", CascadeBranchName(in))
	body := CascadePRBody(in)
	assert.Contains(t, body, "github.com/kairos-io/kairos-sdk")
	assert.Contains(t, body, "pseudo")
	assert.Contains(t, strings.ToLower(body), "tag")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(body), PRMarker(in.Key)))
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement** — append to `prbody.go`:

```go
func CascadeBranchName(in Intent) string {
	return "ksec/cascade-" + slug(in.Package) + "-pseudo"
}

func CascadePRBody(in Intent) string {
	return fmt.Sprintf(`## Automated security cascade

This bumps **%s** to a pseudo-version of its latest default-branch commit, which
contains an unreleased security fix. Once a maintainer cuts a release tag for
that module, kairos-security will re-pin this PR to the tagged version.

- Module: `+"`%s`"+`
- Severity: %s
- Please **tag a release** of the upstream module so this can be pinned cleanly.

This PR was opened automatically by kairos-security. CI on this PR runs the tests.

%s`, in.Package, in.Package, in.Severity, PRMarker(in.Key))
}
```

- [ ] **Step 4: Run + commit**

Run: `go test ./internal/remediate/...` (PASS).
```bash
git add internal/remediate/prbody.go internal/remediate/prbody_test.go
git commit -m "feat(remediate): cascade PR branch/body (pseudo-version + tag request)"
```

---

### Task 7: `GitExecutor.Cascade` (integration)

**Files:** Modify `internal/remediate/git_executor.go` (replace the Task-5 cascade stub).

**Interfaces:** `func (g *GitExecutor) Cascade(in Intent, runID string) (state.LedgerEntry, error)` — clone the consumer (`in.Repo`), branch `CascadeBranchName(in)`, `go get <in.Package>@<in.Ref>` (pseudo-version of the module's default branch), `go mod tidy`, `go build ./...` verify (fail → `build-failed`/no push), commit, push, `gh pr create` with `CascadePRBody(in)`; entry `Kind:"cascade"`, `Pseudo:true`, `CascadeFrom:in.CascadeFrom`, `State:"open"`, PR url/number. Dry-run prints + returns a planned cascade entry. Uses `g.run` (token-redacting).

- [ ] **Step 1: Implement** (replace the stub; mirror `Open` but with `go get <module>@<Ref>` and the cascade branch/body). Reuse `g.cloneURL`, `g.run`, `prNumberFromURL`.

```go
func (g *GitExecutor) Cascade(in Intent, runID string) (state.LedgerEntry, error) {
	branch := CascadeBranchName(in)
	entry := state.LedgerEntry{
		Key: in.Key, Repo: in.Repo, Package: in.Package, Branch: branch, Kind: "cascade",
		CascadeFrom: in.CascadeFrom, Pseudo: true, Severity: in.Severity, CreatedRun: runID, LastActionRun: runID,
		Bump: state.Bump{Package: in.Package, To: in.Ref},
	}
	if g.DryRun {
		fmt.Printf("[dry-run] would cascade %s: branch %s, go get %s@%s (pseudo)\n", in.Repo, branch, in.Package, in.Ref)
		entry.State = "planned"
		entry.History = []state.LedgerEvent{{Run: runID, Action: "plan-cascade"}}
		return entry, nil
	}
	dir, err := os.MkdirTemp("", "ksec-cas-*")
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
	if _, err := g.run(dir, "go", "get", in.Package+"@"+in.Ref); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "go", "mod", "tidy"); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "go", "build", "./..."); err != nil {
		entry.State = "build-failed"
		entry.NeedsHuman = true
		entry.History = []state.LedgerEvent{{Run: runID, Action: "cascade-build-failed", Detail: err.Error()}}
		return entry, nil
	}
	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	if _, err := g.run(dir, "git", "commit", "-am", "chore(security): cascade-bump "+in.Package); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "git", "push", "-u", "origin", branch); err != nil {
		return entry, err
	}
	out, err := g.run(dir, "gh", "pr", "create", "-R", in.Repo, "--head", branch,
		"--title", "chore(security): cascade-bump "+in.Package, "--body", CascadePRBody(in))
	if err != nil {
		return entry, err
	}
	entry.PRURL = strings.TrimSpace(string(out))
	entry.PRNumber = prNumberFromURL(entry.PRURL)
	entry.State = "open"
	entry.History = []state.LedgerEvent{{Run: runID, Action: "cascade-opened", Detail: entry.PRURL}}
	return entry, nil
}
```

- [ ] **Step 2: Build + vet + test** — `go build ./... && go vet ./... && go test ./...` (PASS; integration, not unit-tested).
- [ ] **Step 3: Commit** — `git commit -am "feat(remediate): GitExecutor.Cascade (pseudo-version bump + PR)"` (only `git_executor.go`).

---

### Task 8: `GitExecutor.Repin` (integration)

**Files:** Modify `internal/remediate/git_executor.go` (replace the Task-5 repin stub).

**Interfaces:** `func (g *GitExecutor) Repin(e state.LedgerEntry, runID string) (state.LedgerEntry, error)` — if a published tag exists for the module (`go list -m -versions <module>` returns a tagged version newer than the current pin/pseudo), re-pin: clone, checkout `e.Branch`, `go get <module>@<tag>`, tidy, build, force-push; set `Pseudo:false`, `PinTarget:<tag>`. If no tag yet, record `awaiting-release` and return unchanged (no error). Dry-run prints + returns unchanged.

- [ ] **Step 1: Implement** (replace the stub). Use `go list -m -versions <module>` (run in a temp module or with `GOFLAGS=-mod=mod`); pick the highest non-pseudo tag via `compareVersions`. If none, `awaiting-release`. Otherwise clone, checkout `e.Branch`, `go get module@tag`, tidy, `go build` verify (fail → `build-failed`/no push), `git push --force`. Reuse `g.run`, the `compareVersions` helper from `planner.go`, and the `Adjust` force-push pattern.

```go
func (g *GitExecutor) Repin(e state.LedgerEntry, runID string) (state.LedgerEntry, error) {
	module := e.Package
	if g.DryRun {
		fmt.Printf("[dry-run] would check %s for a release tag to re-pin %s\n", module, e.Repo)
		return e, nil
	}
	// Find the latest published tag for the module.
	out, err := g.run("", "go", "list", "-m", "-versions", module)
	tag := latestTag(out) // highest vN.N.N token on the line; "" if none
	if err != nil || tag == "" {
		e.Blocked = "awaiting-release"
		return e, nil
	}
	dir, err := os.MkdirTemp("", "ksec-pin-*")
	if err != nil {
		return e, err
	}
	defer os.RemoveAll(dir)
	if _, err := g.run("", "git", "clone", g.cloneURL(e.Repo), dir); err != nil {
		return e, err
	}
	if _, err := g.run(dir, "git", "checkout", e.Branch); err != nil {
		return e, err
	}
	if _, err := g.run(dir, "go", "get", module+"@"+tag); err != nil {
		return e, err
	}
	if _, err := g.run(dir, "go", "mod", "tidy"); err != nil {
		return e, err
	}
	if _, err := g.run(dir, "go", "build", "./..."); err != nil {
		e.State = "build-failed"
		e.NeedsHuman = true
		e.History = append(e.History, state.LedgerEvent{Run: runID, Action: "repin-build-failed", Detail: err.Error()})
		return e, nil
	}
	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	if out, _ := g.run(dir, "git", "status", "--porcelain"); len(bytes.TrimSpace(out)) == 0 {
		e.Pseudo = false
		e.PinTarget = tag
		e.Blocked = ""
		return e, nil // already at the tag
	}
	if _, err := g.run(dir, "git", "commit", "-am", "chore(security): re-pin "+module+" to "+tag); err != nil {
		return e, err
	}
	if _, err := g.run(dir, "git", "push", "--force", "origin", e.Branch); err != nil {
		return e, err
	}
	e.Pseudo = false
	e.PinTarget = tag
	e.Blocked = ""
	e.Bump.To = tag
	e.LastActionRun = runID
	e.History = append(e.History, state.LedgerEvent{Run: runID, Action: "repinned", Detail: tag})
	return e, nil
}

// latestTag returns the highest vN.N.N token found in `go list -m -versions`
// output (a single space-separated line: "<module> v1 v1.0.1 ..."), or "".
func latestTag(b []byte) string {
	best := ""
	for _, tok := range strings.Fields(string(b)) {
		if !strings.HasPrefix(tok, "v") {
			continue
		}
		if best == "" || compareVersions(tok, best) > 0 {
			best = tok
		}
	}
	return best
}
```

- [ ] **Step 2: Build + vet + test** — `go build ./... && go vet ./... && go test ./...` (PASS).
- [ ] **Step 3: Commit** — `git commit -am "feat(remediate): GitExecutor.Repin (re-pin pseudo cascade to a release tag)"`.

---

### Task 9: Wire depgraph + dashboard cascade display

**Files:** Modify `cmd/ksec/main.go`, `internal/render/render.go`, `internal/render/coord_test.go`.

**Interfaces:** the `remediate` command builds the `DepGraph` (fetch each tracked repo's `go.mod` via `gh.GetFile(repo, "go.mod", branch)`), passes it to `Plan`; the dashboard ledger marks `pseudo` bumps and shows `cascadeFrom`.

- [ ] **Step 1: Wire the graph in `main.go`**

In `newRemediateCmd`, after building `prsByRepo` and loading `repos`, build the graph:

```go
			gomodByRepo := map[string][]byte{}
			for _, r := range repos {
				if b, err := gh.GetFile(r.Repo, "go.mod", r.Branch); err == nil {
					gomodByRepo[r.Repo] = b
				}
			}
			graph := remediate.BuildGraph(repos, gomodByRepo)
```

Replace the `Plan(...)` call's `nil` graph shim with `graph` (remove the `// TODO(plan-4b task 9)`): `intents, deferred := remediate.Plan(c, ledger, prsByRepo, graph, maxPRs)`.

- [ ] **Step 2: Dashboard — mark pseudo + cascadeFrom**

In `internal/render/render.go` ledger row: when `e.Pseudo`, render the Bump as `pkg@<ref> (pseudo)`; when `e.CascadeFrom != ""`, append a small `↳ from <cascadeFrom>` note in the Kind cell (`cascade ↳`). Add a `coord_test.go` assertion:

```go
func TestDashboardMarkdownShowsPseudoCascade(t *testing.T) {
	in := Input{Ledger: state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|github.com/kairos-io/kairos-sdk", Repo: "kairos-io/immucore",
			Package: "github.com/kairos-io/kairos-sdk", State: "open", Kind: "cascade", Pseudo: true,
			CascadeFrom: "kairos-io/kairos-sdk|x", Bump: state.Bump{Package: "github.com/kairos-io/kairos-sdk", To: "main"}}}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "cascade")
	assert.Contains(t, md, "pseudo")
}
```

Mirror in `html.go` if the HTML ledger is present; regenerate goldens (`UPDATE_GOLDEN=1 go test ./internal/render/...`), eyeball, re-run.

- [ ] **Step 3: Build + vet + full test + smoke**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: all pass; no leftover `nil` graph shim / TODO. `go run ./cmd/ksec remediate --help` still works.

- [ ] **Step 4: Commit**

```bash
git add cmd/ksec/main.go internal/render/render.go internal/render/coord_test.go internal/render/html.go internal/render/testdata/
git commit -m "feat(remediate): build depgraph and wire cascade; dashboard shows pseudo/cascade"
```

---

## Self-review

**Spec coverage** (§7.2 cascade + §6 cascade ledger fields + the two 4a follow-ups):
- depgraph from go.mod (module↔repo, consumers) → Task 2. ✓
- Cascade on merged first-party fix → pseudo bump in consumers → Tasks 4 (plan), 7 (execute). ✓
- Pseudo-version (`go get @defaultBranch`), PR notes tag request → Tasks 6 (body), 7. ✓
- Re-pin to tag once released; `awaiting-release` until then → Tasks 4 (plan), 8 (execute). ✓
- Recursion (merged cascade triggers further cascade) → Task 4 (planner keys on any merged first-party-module entry). ✓
- Cascade ledger fields cascadeFrom/pinTarget/pseudo → Task 1. ✓
- Shared blast-radius cap (open + cascade) → Task 4. ✓
- Verify-before-push; build-failed/needsHuman, no broken push; dry-run no writes; token redaction → Tasks 7, 8. ✓
- Dashboard shows pseudo + cascadeFrom → Task 9. ✓
- **Follow-up A** (own-PR by branch, not login) → Task 3. ✓
- **Follow-up B** (version-aware MatchPR) → Task 3. ✓

**Deferred to 4c (correctly absent):** nib build-break repair (cascade/repin breakage records `build-failed`/`needsHuman` here), conflict resolution, toolchain bumps, AI coordination summary.

**Placeholder scan:** none — complete code/commands in every step.

**Type consistency:** `DepGraph`/`BuildGraph`/`ModuleOf`/`Consumers`/`BranchOf` (Task 2) used by Tasks 4, 9. `PullRequest.HeadRef` + `isOwnPR` + `MatchPR(pkg,version,prs)` (Task 3) used by the planner (Task 4 call site updated in Task 3). `IntentCascade`/`IntentRepin` + `Intent.Ref`/`CascadeFrom` (Task 4) used by Tasks 5, 7, 8. `Executor.Cascade`/`Repin` (Task 5) implemented by `FakeExecutor` (Task 5) and `GitExecutor` (Tasks 7, 8). `CascadeBranchName`/`CascadePRBody` (Task 6) used by Task 7. `compareVersions` (planner.go, Plan 2) reused by `latestTag` (Task 8). Ledger fields (Task 1) rendered in Task 9.

---

## Operational notes

- Live cascade/re-pin needs `KSEC_BOT_TOKEN` write scope on the consumer repos; without it, cascade records `error` per intent (logged, isolated).
- Pseudo-version `go get <module>@<branch>` requires module-graph network access on the runner (same as the existing source-CVE clone+build).
- A cascade that breaks the consumer build is recorded `build-failed` + `needsHuman` and surfaced — **4c's nib agent** will attempt to repair these.
- Re-pin uses `go list -m -versions`; a module with no tags stays `awaiting-release` indefinitely (the PR body asks maintainers to tag), which is the intended "flow + follow-up" behavior.
