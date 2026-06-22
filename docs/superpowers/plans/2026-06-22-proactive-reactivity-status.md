# Proactive Reactivity & Status Surfacing Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Recognize every bot, proactively supersede a conflicted-but-still-relevant external PR with a fresh `ksec/` PR, surface readable per-PR status/action so a human can pick, and keep an un-scannable repo an honest one-line status instead of an error wall.

**Architecture:** Five threads on the existing pipeline: (A) generalize bot detection to any `[bot]` author; (B) truncate scanner error walls; (C) a per-repo source-scan opt-out with a `skipped` status; (D) detect an adopted external PR going `CONFLICTING` and mark it `upstream-conflict`; (E) plan + execute a proactive `supersede` (open a clean `ksec/` PR, comment on the foreign PR, link it via a new `Supersedes` ledger field) and surface all of it on the dashboard. Builds on `internal/remediate`, `internal/collect`, `internal/render`, `internal/state`, `cmd/ksec`.

**Tech Stack:** Go 1.22, `stretchr/testify`, `gh`/`git`/`go`/`govulncheck` CLIs, `gopkg.in/yaml` (repos.yaml → `[]state.Repo`, lowercased field names).

## Global Constraints

- Module `github.com/kairos-io/security`; binary `ksec`; Go 1.22.
- **Never edit or force-push a foreign PR branch.** Supersede only creates/pushes our own `ksec/` branch and *comments* on the foreign PR.
- Verify-before-push on the superseding PR; build failure (agent can't repair) → `needs-human`, no PR. Dry-run performs zero writes. Token redacted via `g.run`.
- Supersede counts toward the `--max-prs` cap (it opens a PR).
- Committed `dashboard.md`/`dashboard.json` stay deterministic (sorted; no volatile fields on the committed copy).
- Human-facing output never shows a raw SHA-256 id; truncated errors never exceed the cap.
- A reachable govulncheck finding default severity stays `high` (unchanged). stdlib keeps `Package=="stdlib"` (unchanged).

---

## File structure

```
internal/remediate/matcher.go        # isBotLogin; classifySource adds "bot" (modify)
internal/remediate/matcher_test.go   # (modify)
internal/collect/prs.go              # prSource adds "bot"; isSecurityPR uses isBotLogin (modify)
internal/collect/prs_test.go         # (modify)
internal/collect/govulncheck_result.go      # truncErr; cap the error (modify)
internal/collect/govulncheck_result_test.go # (modify)
internal/state/types.go              # Repo.Scan (ScanConfig); LedgerEntry.Supersedes (modify)
cmd/ksec/main.go                     # govulncheckRunner honors opt-out (modify)
internal/remediate/git_executor.go   # Reconcile upstream-conflict; GitExecutor.Supersede (modify)
internal/remediate/intent.go         # IntentSupersede (modify)
internal/remediate/planner.go        # supersede planning (modify)
internal/remediate/planner_test.go   # (modify)
internal/remediate/run.go            # Executor.Supersede + Run case (modify)
internal/remediate/fake.go           # FakeExecutor.Supersede (modify)
internal/remediate/run_test.go       # (modify)
internal/render/render.go            # Open PRs status/action; supersedes link; needs-human roll-up; skipped status (modify)
internal/render/render_test.go       # (modify)
internal/render/html.go              # mirror (modify)
internal/render/testdata/            # regenerated goldens
repos.yaml                           # kairos-must-burn scan.source:false (modify)
```

---

### Task 1: Generalize bot detection

**Files:** Modify `internal/remediate/matcher.go`, `internal/remediate/matcher_test.go`, `internal/collect/prs.go`, `internal/collect/prs_test.go`.

**Interfaces:**
- Produces: `func isBotLogin(login string) bool` (in `matcher.go`); `classifySource` returns `ksec|renovate|dependabot|bot|human`.
- `internal/collect/prs.go`: `prSource` returns the same set; `isSecurityPR` treats any bot login as a bot.

- [ ] **Step 1: Write the failing tests**

In `internal/remediate/matcher_test.go` add:

```go
func TestClassifySourceRecognizesAnyBot(t *testing.T) {
	// a novel GitHub App bot must classify as "bot", not "human"
	assert.Equal(t, "bot", classifySource(ghclient.PullRequest{Author: "kairos-io-bot[bot]"}))
	assert.Equal(t, "renovate", classifySource(ghclient.PullRequest{Author: "renovate[bot]"}))
	assert.Equal(t, "dependabot", classifySource(ghclient.PullRequest{Author: "dependabot[bot]"}))
	assert.Equal(t, "ksec", classifySource(ghclient.PullRequest{HeadRef: "ksec/x"}))
	assert.Equal(t, "human", classifySource(ghclient.PullRequest{Author: "alice"}))
}
```

In `internal/collect/prs_test.go` add (alongside the existing `TestOpenPRsTracksAndClassifies`):

```go
func TestOpenPRsTracksAnyBot(t *testing.T) {
	gh := fakePRGH{byRepo: map[string][]ghclient.PullRequest{
		"o/r": {{Number: 38, Title: "bump x", Author: "kairos-io-bot[bot]", URL: "u38"}},
	}}
	prs, errs := OpenPRs([]state.Repo{{Repo: "o/r"}}, gh)
	require.Empty(t, errs)
	require.Len(t, prs, 1) // tracked even though not in the renovate/dependabot/ksec allowlist
	assert.Equal(t, "bot", prs[0].Source)
}
```

- [ ] **Step 2: Run red.** `go test ./internal/remediate/... ./internal/collect/...`

- [ ] **Step 3: Implement**

In `matcher.go`, add `isBotLogin` and use it in `classifySource`:

```go
func isBotLogin(login string) bool { return strings.HasSuffix(login, "[bot]") }

func classifySource(pr ghclient.PullRequest) string {
	if isOwnPR(pr) {
		return "ksec"
	}
	switch pr.Author {
	case "renovate[bot]":
		return "renovate"
	case "dependabot[bot]":
		return "dependabot"
	}
	if isBotLogin(pr.Author) {
		return "bot"
	}
	return "human"
}
```

In `internal/collect/prs.go`: replace the `botAuthors` membership in `isSecurityPR` with `isBotLogin`, and add `bot` to `prSource`. (Add an `isBotLogin` to this package too — `strings.HasSuffix(login, "[bot]")` — to avoid a cross-package import; keep the named renovate/dependabot/kairos-security-bot mapping for `prSource`.) Concretely:

```go
func isBotLogin(login string) bool { return strings.HasSuffix(login, "[bot]") }

func prSource(author string) string {
	switch author {
	case "renovate[bot]":
		return "renovate"
	case "dependabot[bot]":
		return "dependabot"
	case "kairos-security-bot":
		return "ksec"
	}
	if isBotLogin(author) {
		return "bot"
	}
	return "human"
}

func isSecurityPR(pr ghclient.PullRequest) bool {
	if isBotLogin(pr.Author) || pr.Author == "kairos-security-bot" {
		return true
	}
	for _, l := range pr.Labels {
		if secLabels[l] {
			return true
		}
	}
	return false
}
```

Remove the now-unused `botAuthors` map (and confirm `strings` is imported in `prs.go`).

- [ ] **Step 4: Run green + build + commit**

Run: `go test ./internal/remediate/... ./internal/collect/... && go build ./...`
```bash
git add internal/remediate/matcher.go internal/remediate/matcher_test.go internal/collect/prs.go internal/collect/prs_test.go
git commit -m "feat(remediate): recognize any [bot] author, not a 3-login allowlist"
```

---

### Task 2: Truncate scanner error walls

**Files:** Modify `internal/collect/govulncheck_result.go`, `internal/collect/govulncheck_result_test.go`.

**Interfaces:** Produces `func truncErr(b []byte, max int) string` (collapse whitespace, cap with an ellipsis note). `ClassifyGovulncheck`'s error uses it.

- [ ] **Step 1: Write the failing test** — add to `govulncheck_result_test.go`:

```go
func TestClassifyGovulncheckTruncatesHugeError(t *testing.T) {
	huge := []byte(strings.Repeat("glib-2.0 not found; ", 5000)) // ~100KB, many lines
	_, err := ClassifyGovulncheck([]byte(`{"config":{}}`+"\n"+`{"progress":{}}`), huge, errors.New("exit status 1"))
	require.Error(t, err)
	assert.LessOrEqual(t, len(err.Error()), 320) // capped, not a 100KB wall
	assert.Contains(t, err.Error(), "truncated")
}

func TestTruncErr(t *testing.T) {
	assert.Equal(t, "abc", truncErr([]byte("  abc \n"), 240))           // trimmed, under cap
	long := truncErr([]byte(strings.Repeat("x", 1000)), 240)
	assert.Equal(t, 240+len(" … (truncated)"), len(long))
	assert.True(t, strings.HasSuffix(long, "… (truncated)"))
}
```

Ensure `strings` is imported in the test.

- [ ] **Step 2: Run red.** `go test ./internal/collect/...`

- [ ] **Step 3: Implement** — in `govulncheck_result.go`:

```go
// truncErr renders tool stderr as a one-line, length-capped summary so a
// build-failure (e.g. a cgo app missing system libs) doesn't flood the
// dashboard or the committed findings.json with a multi-KB wall.
func truncErr(b []byte, max int) string {
	s := strings.Join(strings.Fields(string(b)), " ") // collapse all whitespace/newlines
	if len(s) <= max {
		return s
	}
	return s[:max] + " … (truncated)"
}
```

Change the failure return to `return nil, fmt.Errorf("govulncheck: %v: %s", runErr, truncErr(stderr, 240))`. Add `"strings"` to the imports.

- [ ] **Step 4: Run green + commit**

Run: `go test ./internal/collect/... && go build ./...`
```bash
git add internal/collect/govulncheck_result.go internal/collect/govulncheck_result_test.go
git commit -m "fix(collect): cap govulncheck error output (no more dashboard error walls)"
```

---

### Task 3: Per-repo source-scan opt-out

**Files:** Modify `internal/state/types.go`, `cmd/ksec/main.go`, `repos.yaml`, `internal/render/render.go`, `internal/render/render_test.go`.

**Interfaces:**
- `state.ScanConfig{Source *bool}`; `Repo.Scan ScanConfig`; `func (Repo) SourceScanEnabled() bool` (nil → true).
- `govulncheckRunner` skips (returns `nil, nil`) when `!SourceScanEnabled()`.
- Render shows `skipped: not source-scannable` for opted-out repos with no findings.

- [ ] **Step 1: Add the state type + method, write the failing test**

In `internal/state/types.go`:

```go
type ScanConfig struct {
	Source *bool `json:"source,omitempty" yaml:"source"`
}

func (r Repo) SourceScanEnabled() bool { return r.Scan.Source == nil || *r.Scan.Source }
```

Add `Scan ScanConfig \`json:"scan,omitempty" yaml:"scan"\`` to `Repo`.

In `internal/state/types_test.go` (create if absent) or an existing state test:

```go
func TestSourceScanEnabled(t *testing.T) {
	assert.True(t, Repo{}.SourceScanEnabled()) // default
	f := false
	assert.False(t, Repo{Scan: ScanConfig{Source: &f}}.SourceScanEnabled())
	tr := true
	assert.True(t, Repo{Scan: ScanConfig{Source: &tr}}.SourceScanEnabled())
}
```

- [ ] **Step 2: Run red.** `go test ./internal/state/...`

- [ ] **Step 3: Implement the skip + repos.yaml + render status**

In `cmd/ksec/main.go` `govulncheckRunner`, at the very top (before the go.mod check):

```go
	if !r.SourceScanEnabled() {
		return nil, nil // explicitly opted out of source scanning
	}
```

In `repos.yaml`, set the opt-out on `kairos-must-burn` (it's a GTK4/cgo app that can't build headless):

```yaml
  - repo: kairos-io/kairos-must-burn
    kind: org
    criticality: low
    scan:
      source: false
```

In `internal/render/render.go`, the per-repo status computation: when a repo has no findings and `!repo.SourceScanEnabled()`, render its status as `skipped: not source-scannable` instead of `clean`. Add a `render_test.go` assertion:

```go
func TestDashboardMarksSkippedRepo(t *testing.T) {
	f := false
	in := Input{Repos: []state.Repo{{Repo: "o/r", Scan: state.ScanConfig{Source: &f}}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "skipped: not source-scannable")
}
```

(Find the existing per-repo table loop; it currently emits `clean`/`ok`/`⚠️ errors`. Add the `skipped` branch keyed on `!repo.SourceScanEnabled()` with zero findings and no error.)

- [ ] **Step 4: Run green, regenerate goldens, build, commit**

Run: `go test ./internal/state/... ./internal/render/...`; if the per-repo table golden changed, `UPDATE_GOLDEN=1 go test ./internal/render/...`, eyeball, re-run. Then `go build ./... && go test ./...`. Validate YAML: `python3 -c "import yaml; yaml.safe_load(open('repos.yaml'))" && echo OK`.
```bash
git add internal/state/types.go internal/state/types_test.go cmd/ksec/main.go repos.yaml internal/render/render.go internal/render/render_test.go internal/render/testdata/
git commit -m "feat: per-repo source-scan opt-out; kairos-must-burn marked not source-scannable"
```

---

### Task 4: Reconcile detects upstream-conflict on adopted external PRs

**Files:** Modify `internal/remediate/git_executor.go`.

**Interfaces:** `Reconcile` sets `Blocked: "upstream-conflict"` when an adopted (`source != ksec`, non-`ksec/` branch) open PR reports `mergeable == CONFLICTING`. (The existing `ksec/`-branch case still routes to `ResolveConflict`.)

- [ ] **Step 1: Implement** — in `git_executor.go` `Reconcile`, after the existing merged/closed/open switch and the `ResolveConflict` block, add the foreign-conflict branch. Replace the final tail:

```go
	if e.State == "open" && view.Mergeable == "CONFLICTING" && strings.HasPrefix(branch, "ksec/") {
		return g.ResolveConflict(e, runID)
	}
	// A foreign (adopted) PR we don't own that is conflicting: we cannot rebase
	// it, so flag it for the planner to supersede with our own PR.
	if e.State == "open" && view.Mergeable == "CONFLICTING" && e.Source != "ksec" {
		if e.Blocked != "upstream-conflict" {
			e.Blocked = "upstream-conflict"
			e.LastActionRun = runID
			e.History = append(e.History, state.LedgerEvent{Run: runID, Action: "upstream-conflict"})
		}
		return e, nil
	}
	// Conflict cleared on a previously-blocked entry.
	if e.Blocked == "upstream-conflict" && view.Mergeable != "CONFLICTING" {
		e.Blocked = ""
	}
	return e, nil
```

- [ ] **Step 2: Build + vet + test + commit** — integration (no unit test; the conflict path shells `gh`).

Run: `go build ./... && go vet ./... && go test ./...`
```bash
git add internal/remediate/git_executor.go
git commit -m "feat(remediate): flag a conflicting adopted PR as upstream-conflict"
```

---

### Task 5: Supersede intent + planner + ledger field

**Files:** Modify `internal/remediate/intent.go`, `internal/remediate/planner.go`, `internal/remediate/planner_test.go`, `internal/state/types.go`.

**Interfaces:**
- `IntentSupersede IntentType = "supersede"`. `Intent` already carries `PRNumber`/`PRURL` (reused as the foreign PR to supersede).
- `LedgerEntry.Supersedes string` (the foreign PR URL this entry replaced).
- Planner: in the per-target loop, an existing adopted entry (`Source != "ksec"`) with `Blocked == "upstream-conflict"` and still relevant (we are iterating it because the finding still exists) → emit `IntentSupersede` into the capped pool.

- [ ] **Step 1: Add the intent type + ledger field**

In `intent.go`: add `IntentSupersede IntentType = "supersede"`.
In `state/types.go` `LedgerEntry`: add `Supersedes string \`json:"supersedes,omitempty"\`` (after `PinTarget`).

- [ ] **Step 2: Write the failing planner test** — in `planner_test.go`:

```go
func TestPlanSupersedesConflictedAdoptedPR(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "f1", Repo: "o/r", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net",
			FixedVersion: "0.33.0", Severity: "high"},
	}}
	ledger := state.Ledger{Entries: []state.LedgerEntry{{
		Key: "o/r|golang.org/x/net", Repo: "o/r", Package: "golang.org/x/net", State: "open",
		Source: "bot", Blocked: "upstream-conflict", PRNumber: 38, PRURL: "https://github.com/o/r/pull/38",
		Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"},
	}}}
	intents, _ := Plan(c, ledger, nil, nil, 10)
	var sup *Intent
	for i := range intents {
		if intents[i].Type == IntentSupersede {
			sup = &intents[i]
		}
	}
	require.NotNil(t, sup)
	assert.Equal(t, "o/r|golang.org/x/net", sup.Key)
	assert.Equal(t, 38, sup.PRNumber)
	assert.Equal(t, "https://github.com/o/r/pull/38", sup.PRURL)
}
```

- [ ] **Step 3: Run red.** `go test ./internal/remediate/...`

- [ ] **Step 4: Implement** — in `planner.go` step 3, before the `if e.State == "open" || "conflicted" { continue }` skip, intercept the upstream-conflict case. Restructure the existing-entry check:

```go
		if e, ok := ledger.ByKey(k); ok {
			// A conflicted adopted PR we can't rebase: supersede it with our own.
			if e.Source != "ksec" && e.Blocked == "upstream-conflict" && e.State == "open" {
				pool = append(pool, newPR{
					intent: Intent{Type: IntentSupersede, Key: k, Repo: t.repo, Package: t.pkg, Severity: t.sev,
						Bump: state.Bump{Package: t.pkg, To: t.to}, PRNumber: e.PRNumber, PRURL: e.PRURL},
					sev: t.sev,
				})
				continue
			}
			if e.State == "open" || e.State == "conflicted" {
				continue
			}
			if (e.State == "merged" || e.State == "closed") && compareVersions(e.Bump.To, t.to) >= 0 {
				continue
			}
		}
```

NOTE: the `pool` is declared further down in the current code. Move the `type newPR struct{...}` + `var pool []newPR` declaration to ABOVE this target loop so supersede intents can be appended here (then the later direct/cascade/toolchain appends use the same `pool`). Adjust accordingly; the final sort/cap over `pool` is unchanged.

- [ ] **Step 5: Run green + build + commit**

Run: `go test ./internal/remediate/... && go build ./...`
```bash
git add internal/remediate/intent.go internal/remediate/planner.go internal/remediate/planner_test.go internal/state/types.go
git commit -m "feat(remediate): plan a supersede for a conflicted adopted PR"
```

---

### Task 6: Run + Executor.Supersede + GitExecutor.Supersede

**Files:** Modify `internal/remediate/run.go`, `internal/remediate/fake.go`, `internal/remediate/run_test.go`, `internal/remediate/git_executor.go`.

**Interfaces:** `Executor` gains `Supersede(in Intent, run string) (state.LedgerEntry, error)`; `Run` handles `IntentSupersede` (error-isolated, error entry `Blocked:"supersede-failed"`); `FakeExecutor.Superseded` map + method; real `GitExecutor.Supersede`.

- [ ] **Step 1: Add the test** — in `run_test.go`:

```go
func TestRunSupersede(t *testing.T) {
	intents := []Intent{{Type: IntentSupersede, Key: "o/r|p", Repo: "o/r", Package: "p",
		PRNumber: 38, PRURL: "https://github.com/o/r/pull/38", Bump: state.Bump{Package: "p", To: "1.2.3"}}}
	fake := &FakeExecutor{Superseded: map[string]state.LedgerEntry{
		"o/r|p": {Key: "o/r|p", Repo: "o/r", Package: "p", State: "open", Source: "ksec",
			Branch: "ksec/p", Supersedes: "https://github.com/o/r/pull/38", PRNumber: 77},
	}}
	out, results := Run(intents, fake, state.Ledger{}, "2026-06-22")
	require.Len(t, out.Entries, 1)
	assert.Equal(t, "https://github.com/o/r/pull/38", out.Entries[0].Supersedes)
	assert.Equal(t, 77, out.Entries[0].PRNumber)
	require.Len(t, results, 1)
}
```

- [ ] **Step 2: Run red.** `go test ./internal/remediate/...`

- [ ] **Step 3: Implement Executor method, Run case, Fake**

- `run.go`: add `Supersede(in Intent, run string) (state.LedgerEntry, error)` to `Executor`; add a `case IntentSupersede:` mirroring `IntentOpen` (error-isolated; on error a synthetic entry `State:"error"`, `Blocked:"supersede-failed"`, history `supersede-failed`; success stores by `entry.Key`; one `Result`).
- `fake.go`: add `Superseded map[string]state.LedgerEntry`; `Supersede` returns the mapped entry or a sensible default (`{Key:in.Key, Repo:in.Repo, Package:in.Package, State:"open", Source:"ksec", Supersedes:in.PRURL}`).

- [ ] **Step 4: Implement `GitExecutor.Supersede`** (replaces nothing; new method). Mirror `Open`, then comment on the foreign PR:

```go
func (g *GitExecutor) Supersede(in Intent, runID string) (state.LedgerEntry, error) {
	branch := BranchName(in) // existing ksec/ bump branch name
	entry := state.LedgerEntry{
		Key: in.Key, Repo: in.Repo, Package: in.Package, Branch: branch, Source: "ksec", Kind: "direct",
		Severity: in.Severity, Supersedes: in.PRURL, CreatedRun: runID, LastActionRun: runID,
		Bump: state.Bump{Package: in.Package, To: in.Bump.To},
	}
	if g.DryRun {
		fmt.Printf("[dry-run] would supersede %s PR %s with a fresh ksec PR (bump %s@%s)\n",
			in.Repo, in.PRURL, in.Package, in.Bump.To)
		entry.State = "planned"
		return entry, nil
	}
	dir, err := os.MkdirTemp("", "ksec-sup-*")
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
	if _, err := g.run(dir, "go", "get", in.Package+"@v"+strings.TrimPrefix(in.Bump.To, "v")); err != nil {
		return entry, err
	}
	_, _ = g.run(dir, "go", "mod", "tidy")
	if !g.verifyOrRepair(dir, "supersede "+in.Package, runID) {
		entry.State = "build-failed"
		entry.NeedsHuman = true
		entry.History = []state.LedgerEvent{{Run: runID, Action: "supersede-build-failed"}}
		return entry, nil
	}
	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	if _, err := g.run(dir, "git", "commit", "-am", "chore(security): bump "+in.Package+" to "+in.Bump.To); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "git", "push", "-u", "origin", branch); err != nil {
		return entry, err
	}
	out, err := g.run(dir, "gh", "pr", "create", "-R", in.Repo, "--head", branch,
		"--title", "chore(security): bump "+in.Package+" to "+in.Bump.To,
		"--body", fmt.Sprintf("Supersedes %s, which had unresolved conflicts. %s", in.PRURL, PRMarker(in.Key)))
	if err != nil {
		return entry, err
	}
	entry.PRURL = strings.TrimSpace(string(out))
	entry.PRNumber = prNumberFromURL(entry.PRURL)
	entry.State = "open"
	entry.History = []state.LedgerEvent{{Run: runID, Action: "superseded", Detail: in.PRURL}}
	// Comment on the foreign PR (best-effort; never edit/force-push its branch).
	if in.PRNumber > 0 {
		_ = g.GH.PostPRComment(in.Repo, in.PRNumber,
			fmt.Sprintf("Superseded by %s — the original had unresolved conflicts. Tracked by kairos-security.", entry.PRURL))
	}
	return entry, nil
}
```

(Confirm `BranchName`, `prNumberFromURL`, `PRMarker`, `g.GH` exist — they do, from Plans 2-4. If `BranchName` needs a target version, pass `in` as the other executors do.)

- [ ] **Step 5: Build + vet + test + commit**

Run: `go test ./internal/remediate/... && go build ./... && go vet ./...`
```bash
git add internal/remediate/run.go internal/remediate/fake.go internal/remediate/run_test.go internal/remediate/git_executor.go
git commit -m "feat(remediate): GitExecutor.Supersede opens a ksec PR + comments the conflicted one"
```

---

### Task 7: Surface status so a human can pick

**Files:** Modify `internal/render/render.go`, `internal/render/render_test.go`, `internal/render/html.go`, `internal/render/testdata/`.

**Interfaces:** Open PRs gains a Status/Action column (correlating each `TrackedPR` with the ledger); the bot ledger shows `Supersedes`; a new "🚑 Needs human" roll-up lists `NeedsHuman` entries.

- [ ] **Step 1: Write the failing tests** — in `render_test.go`:

```go
func TestOpenPRShowsSupersededStatus(t *testing.T) {
	in := Input{
		OpenPRs: []state.TrackedPR{{Repo: "o/r", Number: 38, Title: "bump x", URL: "https://github.com/o/r/pull/38", Source: "bot"}},
		Ledger: state.Ledger{Entries: []state.LedgerEntry{
			{Key: "o/r|x", Repo: "o/r", Source: "ksec", State: "open", PRNumber: 77,
				PRURL: "https://github.com/o/r/pull/77", Supersedes: "https://github.com/o/r/pull/38"},
		}},
	}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "superseded by")          // status/action surfaced
	assert.Contains(t, md, "/pull/77")               // links the ksec PR
}

func TestNeedsHumanRollup(t *testing.T) {
	in := Input{Ledger: state.Ledger{Entries: []state.LedgerEntry{
		{Key: "o/r|x", Repo: "o/r", State: "build-failed", NeedsHuman: true, Bump: state.Bump{Package: "x", To: "1.0"}},
	}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "🚑 Needs human")
	assert.Contains(t, md, "o/r")
}
```

- [ ] **Step 2: Run red.** `go test ./internal/render/...`

- [ ] **Step 3: Implement** in `render.go`:
- Build `supersededBy := map[string]state.LedgerEntry{}` keyed on `e.Supersedes` (the foreign PR URL) → the ksec entry. In the Open PRs loop, append a status: if `supersededBy[pr.URL]` exists → `superseded by [#N](url)`; else if a ledger entry has `PRURL == pr.URL && Blocked == "upstream-conflict"` → `conflicted → superseding`; else `tracked`. Render as `- [#N title](url) — <source> — <status>`.
- In the bot ledger table, when `e.Supersedes != ""` append `↳ supersedes <url>` to the Bump or State cell.
- Add a "🚑 Needs human" section (before the run-log footer) listing each entry with `NeedsHuman` as `- <repo> <bump> — <Blocked or state>`; omit the section when none.
Mirror the Open PRs status + needs-human section in `html.go` (escaped).

- [ ] **Step 4: Regenerate goldens, build, vet, gofmt, test, commit**

Run: `UPDATE_GOLDEN=1 go test ./internal/render/...` (eyeball: Open PRs show status, needs-human roll-up appears only when populated, no raw ids), re-run; then `go build ./... && go vet ./... && test -z "$(gofmt -l .)" && go test ./...`.
```bash
git add internal/render/render.go internal/render/render_test.go internal/render/html.go internal/render/testdata/
git commit -m "feat(render): Open PRs status/action column, supersedes links, needs-human roll-up"
```

---

## Self-review

**Spec coverage:**
- A (any `[bot]` author) → Task 1. ✓
- B (truncate error walls) → Task 2. ✓
- C (per-repo source-scan opt-out + honest skipped status) → Task 3. ✓
- D (Reconcile flags upstream-conflict, never touches foreign branch) → Task 4. ✓
- E (plan + execute proactive supersede; comment, link via Supersedes; needs-human on build fail) → Tasks 5, 6. ✓
- Surface status (Open PRs status/action, supersedes links, needs-human roll-up, skipped) → Tasks 3, 7. ✓
- Safety (own `ksec/` branch only; verify-before-push; dry-run no writes; capped; token redacted) → Task 6. ✓

**Placeholder scan:** none — complete code/commands per step; integration touch points (Reconcile, Supersede) reference existing helpers (`g.run`, `cloneURL`, `verifyOrRepair`, `BranchName`, `prNumberFromURL`, `PRMarker`, `g.GH.PostPRComment`).

**Type consistency:** `isBotLogin`/`classifySource`/`prSource` (Task 1). `truncErr` (Task 2). `ScanConfig`/`Repo.Scan`/`SourceScanEnabled` (Task 3) used by `govulncheckRunner` + render. `Blocked:"upstream-conflict"` (Task 4) read by the planner (Task 5). `IntentSupersede` + `LedgerEntry.Supersedes` (Task 5) implemented by `Executor.Supersede` (Task 6) and rendered (Task 7). `FakeExecutor.Superseded` (Task 6). `state.TrackedPR` correlated with the ledger in render (Task 7).

---

## Operational notes

- Supersede is bounded by `--max-prs`; a burst of conflicted adopted PRs won't open more than the cap per run (the rest stay `upstream-conflict`, surfaced for a human).
- The foreign PR is only *commented*, never pushed to — the superseding work is always on a fresh `ksec/` branch we own.
- `kairos-must-burn` now renders `skipped: not source-scannable`; revisit if CI ever gains the GTK toolchain.
- Once a maintainer closes the superseded foreign PR, Reconcile records it `closed`; our superseding PR proceeds through the normal lifecycle (reconcile/automerge).
