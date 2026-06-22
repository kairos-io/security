# Run Activity Summary & Signal-Not-Noise Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the dashboard tell the story of the run and stop the noise: list **only CVE-tied PRs**, add a **deterministic "📋 This run" activity summary** (what was scanned/found/tracked/acted + why), add **per-phase logging**, and fix the spurious `security`-label warning.

**Architecture:** (1) `collect.OpenPRs` takes the run's findings and keeps a PR only when it's tied to a CVE (finding-package match), `security`-labelled, or ours. (2) `render` computes a pure `RunActivity` from the committed state and renders a top "📋 This run" section (AI narrative stays an optional flavor line). (3) Each phase logs one summary line. (4) `UpsertIssue` ensures labels exist before use. Builds on `internal/collect`, `internal/render`, `internal/ghclient`, `cmd/ksec`.

**Tech Stack:** Go 1.22, `stretchr/testify`, `gh` CLI, existing packages.

## Global Constraints

- Module `github.com/kairos-io/security`; Go 1.22.
- **Only CVE-tied PRs are tracked** — a routine bot bump with no matching finding and no `security` label is NOT listed. findings=0 ⇒ empty Open PRs list.
- The activity summary is **deterministic** (computed from committed state) — always present, byte-identical for identical state (no churn on the committed dashboard). The AI narrative is optional flavor only.
- Committed `dashboard.md`/`dashboard.json` stay deterministic; no raw SHA-256 id in human-facing output.
- Logging goes to **stderr** (never stdout, which carries rendered artifacts in some phases).

---

## File structure

```
internal/collect/prs.go              # OpenPRs(repos, gh, findings) CVE-tied filter (modify)
internal/collect/prs_test.go         # (modify)
cmd/ksec/main.go                     # pass findings to OpenPRs; per-phase log lines (modify)
internal/render/activity.go          # RunActivity + computeActivity (create)
internal/render/activity_test.go     # (create)
internal/render/render.go            # "📋 This run" section (modify)
internal/render/render_test.go       # (modify)
internal/render/html.go              # mirror (modify)
internal/render/testdata/            # regenerated goldens
internal/ghclient/ghclient.go        # UpsertIssue ensures labels exist (modify)
```

---

### Task 1: Open PRs tied to CVEs only

**Files:** Modify `internal/collect/prs.go`, `internal/collect/prs_test.go`, `cmd/ksec/main.go`.

**Interfaces:** `func OpenPRs(repos []state.Repo, gh ghclient.GitHub, findings []state.Finding) ([]state.TrackedPR, []state.CollectionError)`. A PR is kept iff CVE-tied: a finding in the PR's repo has a `Package` that is a case-insensitive substring of the PR title, OR the PR has a `security` label, OR it's ours (`ksec/` branch / `kairos-security-bot`).

- [ ] **Step 1: Write the failing test** — rewrite the tracking tests in `internal/collect/prs_test.go`:

```go
func TestOpenPRsTracksOnlyCVETied(t *testing.T) {
	findings := []state.Finding{
		{Repo: "o/r", Package: "golang.org/x/crypto", CVEID: "GO-1", Severity: "high"},
	}
	gh := fakePRGH{byRepo: map[string][]ghclient.PullRequest{
		"o/r": {
			{Number: 2, Title: "Bump golang.org/x/crypto from 0.39.0 to 0.45.0", Author: "app/dependabot", IsBot: true, URL: "u2"},
			{Number: 5, Title: "Bump github.com/foo/bar to 1.2.3", Author: "app/dependabot", IsBot: true, URL: "u5"}, // no matching finding -> noise, dropped
			{Number: 7, Title: "security hardening", Author: "alice", Labels: []string{"security"}, URL: "u7"},
			{Number: 9, Title: "ksec bump", Author: "someone", URL: "u9", HeadRef: "ksec/bump-x"},
		},
	}}
	prs, errs := OpenPRs([]state.Repo{{Repo: "o/r"}}, gh, findings)
	require.Empty(t, errs)
	nums := map[int]string{}
	for _, p := range prs {
		nums[p.Number] = p.Source
	}
	assert.Contains(t, nums, 2)          // tied to the x/crypto finding
	assert.Equal(t, "dependabot", nums[2])
	assert.NotContains(t, nums, 5)       // unrelated bump -> dropped (no noise)
	assert.Contains(t, nums, 7)          // security label
	assert.Contains(t, nums, 9)          // ours
}

func TestOpenPRsEmptyWhenNoFindings(t *testing.T) {
	gh := fakePRGH{byRepo: map[string][]ghclient.PullRequest{
		"o/r": {{Number: 2, Title: "Bump x", Author: "app/dependabot", IsBot: true}},
	}}
	prs, _ := OpenPRs([]state.Repo{{Repo: "o/r"}}, gh, nil)
	assert.Empty(t, prs) // 0 findings -> no CVE-tied PRs -> no noise
}
```

Delete the now-obsolete `TestOpenPRsTracksAndClassifies`/`TestOpenPRsTracksAnyBot` (their "track any bot" premise is gone) or fold their source-classification assertions into the new test.

- [ ] **Step 2: Run red.** `go test ./internal/collect/...`

- [ ] **Step 3: Implement** — in `prs.go`, replace `isSecurityPR` with a CVE-tied predicate and thread findings in:

```go
func OpenPRs(repos []state.Repo, gh ghclient.GitHub, findings []state.Finding) ([]state.TrackedPR, []state.CollectionError) {
	pkgsByRepo := map[string][]string{}
	for _, f := range findings {
		if f.Package != "" {
			pkgsByRepo[f.Repo] = append(pkgsByRepo[f.Repo], strings.ToLower(f.Package))
		}
	}
	var out []state.TrackedPR
	var errs []state.CollectionError
	for _, repo := range repos {
		prs, err := gh.ListOpenPRs(repo.Repo)
		if err != nil {
			errs = append(errs, state.CollectionError{Repo: repo.Repo, Collector: "prs", Message: err.Error()})
			continue
		}
		for _, pr := range prs {
			if !cveTied(pr, pkgsByRepo[repo.Repo]) {
				continue
			}
			out = append(out, state.TrackedPR{
				Repo: repo.Repo, Number: pr.Number, Title: pr.Title,
				Author: pr.Author, URL: pr.URL, Source: prSource(pr),
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Repo != out[j].Repo {
			return out[i].Repo < out[j].Repo
		}
		return out[i].Number < out[j].Number
	})
	return out, errs
}

// cveTied keeps a PR only when it is security-relevant: it bumps a package that
// has a finding (CVE) in this repo, OR carries a `security` label, OR is ours.
func cveTied(pr ghclient.PullRequest, findingPkgs []string) bool {
	if pr.Author == "kairos-security-bot" || strings.HasPrefix(pr.HeadRef, "ksec/") {
		return true
	}
	for _, l := range pr.Labels {
		if l == "security" {
			return true
		}
	}
	title := strings.ToLower(pr.Title)
	for _, pkg := range findingPkgs {
		if strings.Contains(title, pkg) {
			return true
		}
	}
	return false
}
```

Remove `isSecurityPR` and `secLabels` (no longer used). Keep `prSource`/`botName`. (`strings` stays imported.)

- [ ] **Step 4: Wire the collect command** — in `cmd/ksec/main.go`, change `collect.OpenPRs(repos, gh)` → `collect.OpenPRs(repos, gh, out.Findings)`.

- [ ] **Step 5: Run green + build + commit**

Run: `go test ./internal/collect/... && go build ./... && go test ./...`
```bash
git add internal/collect/prs.go internal/collect/prs_test.go cmd/ksec/main.go
git commit -m "feat(collect): track only CVE-tied open PRs (drop routine-bump noise)"
```

---

### Task 2: Deterministic run-activity summary

**Files:** Create `internal/render/activity.go`, `internal/render/activity_test.go`; modify `internal/render/render.go`, `internal/render/render_test.go`, `internal/render/html.go`, `internal/render/testdata/`.

**Interfaces:**
- `type RunActivity struct { Repos, Skipped, Errored, Findings, Crit, High, Med, Low, Unknown, PRs int; PRsBySource map[string]int; LedgerOpen, NeedsHuman, Superseded, Merged int; Why string }`
- `func computeActivity(in Input) RunActivity` (pure).
- `DashboardMarkdown`/HTML render a "📋 This run" section from it, near the top.

- [ ] **Step 1: Write the failing tests** — `internal/render/activity_test.go`:

```go
package render

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestComputeActivityNoFindings(t *testing.T) {
	f := false
	a := computeActivity(Input{
		Repos: []state.Repo{{Repo: "o/a"}, {Repo: "o/b", Scan: state.ScanConfig{Source: &f}}},
	})
	assert.Equal(t, 2, a.Repos)
	assert.Equal(t, 1, a.Skipped)
	assert.Equal(t, 0, a.Findings)
	assert.Contains(t, a.Why, "No CVEs")
}

func TestComputeActivityWithFindingsAndPRs(t *testing.T) {
	a := computeActivity(Input{
		Repos:         []state.Repo{{Repo: "o/a"}},
		Correlated:    state.Correlated{Findings: []state.Finding{{Repo: "o/a", Severity: "high"}, {Repo: "o/a", Severity: "low"}}},
		OpenPRs:       []state.TrackedPR{{Repo: "o/a", Source: "dependabot"}},
		Ledger:        state.Ledger{Entries: []state.LedgerEntry{{State: "open"}, {State: "build-failed", NeedsHuman: true}}},
		CollectErrors: []state.CollectionError{{Repo: "o/a", Collector: "sourceCVE", Message: "boom"}},
	})
	assert.Equal(t, 2, a.Findings)
	assert.Equal(t, 1, a.High)
	assert.Equal(t, 1, a.PRs)
	assert.Equal(t, 1, a.PRsBySource["dependabot"])
	assert.Equal(t, 1, a.LedgerOpen)
	assert.Equal(t, 1, a.NeedsHuman)
	assert.Equal(t, 1, a.Errored)
	assert.Contains(t, a.Why, "need a human")
}

func TestDashboardShowsThisRun(t *testing.T) {
	md := DashboardMarkdown(Input{Repos: []state.Repo{{Repo: "o/a"}}})
	assert.Contains(t, md, "📋 This run")
	assert.Contains(t, md, "No CVEs")
}
```

- [ ] **Step 2: Run red.** `go test ./internal/render/...`

- [ ] **Step 3: Implement** — `internal/render/activity.go`:

```go
package render

import (
	"fmt"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

type RunActivity struct {
	Repos, Skipped, Errored                 int
	Findings, Crit, High, Med, Low, Unknown int
	PRs                                     int
	PRsBySource                             map[string]int
	LedgerOpen, NeedsHuman, Superseded, Merged int
	Why                                     string
}

func computeActivity(in Input) RunActivity {
	a := RunActivity{Repos: len(in.Repos), PRsBySource: map[string]int{}}
	for _, r := range in.Repos {
		if !r.SourceScanEnabled() {
			a.Skipped++
		}
	}
	erroredRepos := map[string]bool{}
	for _, e := range in.CollectErrors {
		erroredRepos[e.Repo] = true
	}
	a.Errored = len(erroredRepos)
	for _, f := range in.Correlated.Findings {
		a.Findings++
		switch f.Severity {
		case "critical":
			a.Crit++
		case "high":
			a.High++
		case "medium":
			a.Med++
		case "low":
			a.Low++
		default:
			a.Unknown++
		}
	}
	a.PRs = len(in.OpenPRs)
	for _, p := range in.OpenPRs {
		a.PRsBySource[p.Source]++
	}
	for _, e := range in.Ledger.Entries {
		if e.NeedsHuman {
			a.NeedsHuman++
		}
		if e.Supersedes != "" {
			a.Superseded++
		}
		switch e.State {
		case "open":
			a.LedgerOpen++
		case "merged":
			a.Merged++
		}
	}
	a.Why = activityWhy(a)
	return a
}

func activityWhy(a RunActivity) string {
	switch {
	case a.Findings == 0 && a.Errored == 0:
		return fmt.Sprintf("No CVEs found across %d repos — nothing to remediate.", a.Repos)
	case a.Findings == 0 && a.Errored > 0:
		return fmt.Sprintf("No CVEs found, but %d repo(s) could not be scanned — see collection errors.", a.Errored)
	case a.NeedsHuman > 0:
		return fmt.Sprintf("%d finding(s); %d PR(s) open, %d need a human.", a.Findings, a.LedgerOpen, a.NeedsHuman)
	default:
		return fmt.Sprintf("%d finding(s); %d PR(s) open.", a.Findings, a.LedgerOpen)
	}
}

// ActivityMarkdown renders the "📋 This run" section body.
func (a RunActivity) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "- **Scanned:** %d repos", a.Repos)
	if a.Skipped > 0 {
		fmt.Fprintf(&b, " (%d skipped)", a.Skipped)
	}
	if a.Errored > 0 {
		fmt.Fprintf(&b, " · ⚠️ %d errored", a.Errored)
	}
	b.WriteString("\n")
	fmt.Fprintf(&b, "- **Findings:** %d (%d critical / %d high / %d medium / %d low / %d unknown)\n",
		a.Findings, a.Crit, a.High, a.Med, a.Low, a.Unknown)
	fmt.Fprintf(&b, "- **CVE-related PRs:** %d", a.PRs)
	if a.PRs > 0 {
		b.WriteString(" (" + sourceBreakdown(a.PRsBySource) + ")")
	}
	b.WriteString("\n")
	fmt.Fprintf(&b, "- **Remediation:** %d open · %d superseded · %d merged · %d need-human\n",
		a.LedgerOpen, a.Superseded, a.Merged, a.NeedsHuman)
	fmt.Fprintf(&b, "- **Why:** %s\n", a.Why)
	return b.String()
}

func sourceBreakdown(m map[string]int) string {
	var parts []string
	for _, src := range []string{"ksec", "dependabot", "renovate", "bot", "human"} {
		if n := m[src]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", n, src))
		}
	}
	return strings.Join(parts, ", ")
}
```

In `render.go` `DashboardMarkdown`, after the `_Updated …_` line and before the AI narrative/Focus, insert:

```go
	a := computeActivity(in)
	b.WriteString("## 📋 This run\n\n")
	b.WriteString(a.Markdown() + "\n")
```

Mirror in `html.go` (a `<section>` with the same bullet lines, computed from `computeActivity(in)`, escaped). Regenerate goldens.

- [ ] **Step 4: Run green, regenerate goldens, build, vet, gofmt, commit**

Run: `go test ./internal/render/...`; `UPDATE_GOLDEN=1 go test ./internal/render/...` (eyeball the "📋 This run" section + that it's deterministic); then `go build ./... && go vet ./... && test -z "$(gofmt -l .)" && go test ./...`.
```bash
git add internal/render/activity.go internal/render/activity_test.go internal/render/render.go internal/render/render_test.go internal/render/html.go internal/render/testdata/
git commit -m "feat(render): deterministic 'This run' activity summary"
```

---

### Task 3: Per-phase logging + label fix

**Files:** Modify `cmd/ksec/main.go`, `internal/ghclient/ghclient.go`.

**Interfaces:** each phase command logs one stderr summary line; `UpsertIssue` ensures its labels exist (best-effort) before use, removing the `'security' not found` warning.

- [ ] **Step 1: Add per-phase log lines** in `cmd/ksec/main.go` (stderr, at the end of each phase's RunE, before the final save/return):
  - **collect:** `fmt.Fprintf(os.Stderr, "collect: %d repos · %d findings · %d errors · %d PRs tied to CVEs\n", len(repos), len(out.Findings), len(out.Errors), len(prs))`
  - **triage:** after the existing AI ok/fail logging, add the finding count: `fmt.Fprintf(os.Stderr, "triage: %d findings, focus=%d\n", len(c.Findings), len(out.Focus))`
  - **remediate:** after `Run`, summarize the results by action: `fmt.Fprintf(os.Stderr, "remediate: %d intents → %s\n", len(intents), actionCounts(results))` where `actionCounts` tallies `r.Action` (open/adopt/supersede/cascade/toolchain/repin/reconcile/needs-human). Add a small `actionCounts(results []remediate.Result) string` helper in main.go (or inline a map tally).
  - **render:** `fmt.Fprintf(os.Stderr, "render: dashboard.md + site + issue\n")` (or include the issue number if returned).

(Logging is to stderr only; stdout still carries any rendered output. These are observational — no behavior change.)

- [ ] **Step 2: Fix the label warning** — in `internal/ghclient/ghclient.go` `UpsertIssue`, before the create/edit that applies `--label`, ensure each label exists (best-effort, idempotent):

```go
	for _, l := range labels {
		_, _ = c.run("label", "create", l, "-R", repo, "--force") // idempotent; ignore errors
	}
```

(Place it right after entering `UpsertIssue`, before the list/create/edit calls. `--force` updates if it exists. Errors are ignored so a missing label-create permission degrades silently rather than warning.)

- [ ] **Step 3: Build + vet + gofmt + test + smoke + commit**

Run: `go build ./... && go vet ./... && test -z "$(gofmt -l .)" && go test ./...`; smoke `go run ./cmd/ksec render --help`.
```bash
git add cmd/ksec/main.go internal/ghclient/ghclient.go
git commit -m "feat: per-phase log summaries; ensure issue labels exist (no warning)"
```

---

## Self-review

**Spec coverage:**
- A (Open PRs tied to CVEs only; drop bot/dependencies noise; empty when 0 findings) → Task 1. ✓
- B (deterministic "📋 This run" summary: scanned/findings/CVE-PRs/remediation/why; AI optional) → Task 2. ✓
- C (per-phase stderr logging) → Task 3. ✓
- D (label warning fixed) → Task 3. ✓
- Determinism (activity computed from committed state; no churn) → Task 2. ✓

**Placeholder scan:** none — full code for the pure parts (Tasks 1-2); Task 3 is enumerated stderr lines + a best-effort label-ensure loop with exact code.

**Type consistency:** `OpenPRs(…, findings)` + `cveTied`/`prSource` (Task 1) called by the collect command (`out.Findings`). `RunActivity`/`computeActivity`/`Markdown`/`sourceBreakdown` (Task 2) used by render md+html. `actionCounts(results)` (Task 3) consumes `remediate.Result` (existing). `UpsertIssue` label-ensure uses the existing `c.run`.

---

## Operational notes

- A 0-CVE run now renders: *"📋 This run — Scanned 28 repos (1 skipped) · Findings: 0 · CVE-related PRs: 0 · no remediation · Why: No CVEs found across 28 repos — nothing to remediate."* — clear and noise-free.
- The Open PRs list only populates when a CVE finding exists and a PR addresses it; routine dependabot/renovate bumps never appear.
- The activity summary reflects committed-state posture; per-run deltas are in the phase logs.
- If `gh label create` lacks permission, the upsert proceeds without the label (no warning, no failure).
```
