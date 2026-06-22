# Dashboard Signal Quality Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `ksec`'s dashboard a real security signal — fix the silently-empty source scanner, give findings real severities, move routine PRs to their own list, and render links instead of SHA-256 ids.

**Architecture:** Four focused changes: (A) CI uses latest Go + `govulncheckRunner` surfaces build failures as collection errors via a pure classifier; (B) `SourceCVE.Collect` keeps only reachable vulns and maps OSV severity; (C) the PR collector leaves the findings set for a new `state/openprs.json` tracked-PR list; (D) render shows finding `title`→`url` links and a "📋 Open PRs" section. Builds on the existing `internal/collect`, `internal/render`, `internal/state`, `cmd/ksec` and the `.github` workflow.

**Tech Stack:** Go 1.22, `stretchr/testify`, `gh`/`git`/`govulncheck` CLIs, GitHub Actions, `internal/remediate.classifySource`-style PR classification.

## Global Constraints

- Module `github.com/kairos-io/security`; binary `ksec`; Go 1.22 (the module's own floor — unchanged; only the *CI runner* toolchain moves to latest).
- Committed artifacts (`dashboard.md`/`dashboard.json`) stay deterministic (the no-op commit guard). Findings and tracked PRs are sorted before persisting.
- Human-facing markdown/HTML must never show a raw SHA-256 finding id — always a `title` (linked to a URL when one exists).
- A reachable govulncheck finding with no severity data is severity `high`, never `unknown`.
- `stdlib` findings keep `Package == "stdlib"` so the existing Plan-4c toolchain-bump path still fires.

---

## File structure

```
.github/workflows/security-dashboard.yaml   # setup-go 1.22 -> stable (modify)
cmd/ksec/main.go                             # govulncheckRunner uses classifier; collect writes openprs; render loads it (modify)
internal/collect/govulncheck_result.go       # classifyGovulncheck (pure) (create)
internal/collect/govulncheck_result_test.go  # (create)
internal/collect/source.go                   # reachability filter + severityFromOSV (modify)
internal/collect/source_test.go              # (modify)
internal/collect/prs.go                      # OpenPRs() -> []state.TrackedPR; drop Collector role (modify)
internal/collect/prs_test.go                 # (modify)
internal/state/types.go                      # TrackedPR (modify)
internal/state/files.go                      # OpenPRsFile const (modify)
internal/render/render.go                    # Focus links + Open PRs section + Input.OpenPRs (modify)
internal/render/render_test.go               # (modify)
internal/render/html.go                      # mirror (modify)
internal/render/testdata/                     # regenerated goldens
```

---

### Task 1: Loud govulncheck failures + latest CI Go

**Files:** Create `internal/collect/govulncheck_result.go`, `internal/collect/govulncheck_result_test.go`; modify `cmd/ksec/main.go`, `.github/workflows/security-dashboard.yaml`.

**Interfaces:**
- Produces: `func classifyGovulncheck(stdout, stderr []byte, runErr error) ([]byte, error)` — the build-failure-vs-vulns-found decision, pure and testable.

- [ ] **Step 1: Write the failing test**

Create `internal/collect/govulncheck_result_test.go`:

```go
package collect

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassifyGovulncheck(t *testing.T) {
	// Success: no run error -> stdout passed through.
	out, err := classifyGovulncheck([]byte(`{"config":{}}`), nil, nil)
	require.NoError(t, err)
	assert.Equal(t, `{"config":{}}`, string(out))

	// Non-zero exit but stdout has an osv/finding record -> vulns found, normal.
	stdout := []byte(`{"config":{}}` + "\n" + `{"finding":{"osv":"GO-1"}}`)
	out, err = classifyGovulncheck(stdout, []byte("some progress"), errors.New("exit status 3"))
	require.NoError(t, err)
	assert.Equal(t, stdout, out)

	// Non-zero exit, only config/progress on stdout, build error on stderr -> real failure.
	_, err = classifyGovulncheck([]byte(`{"config":{}}`+"\n"+`{"progress":{}}`),
		[]byte("go: updates to go.mod needed; module requires go >= 1.26"), errors.New("exit status 1"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "go >= 1.26")
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/collect/...`

- [ ] **Step 3: Implement**

Create `internal/collect/govulncheck_result.go`:

```go
package collect

import (
	"bytes"
	"fmt"
)

// classifyGovulncheck decides whether a govulncheck run that exited non-zero
// failed for real (build/load error) or merely found vulnerabilities. In -json
// mode govulncheck emits config/progress objects on stdout before analysis, so
// a non-empty stdout does NOT mean it succeeded — only the presence of an
// "osv"/"finding" record does. A non-zero exit with no such record is a real
// failure and must surface (it is how a Go-toolchain mismatch was silently
// reported as zero vulnerabilities).
func classifyGovulncheck(stdout, stderr []byte, runErr error) ([]byte, error) {
	if runErr == nil {
		return stdout, nil
	}
	if bytes.Contains(stdout, []byte(`"osv"`)) || bytes.Contains(stdout, []byte(`"finding"`)) {
		return stdout, nil // vulnerabilities found: non-zero exit is expected
	}
	return nil, fmt.Errorf("govulncheck: %v: %s", runErr, bytes.TrimSpace(stderr))
}
```

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/collect/...`

- [ ] **Step 5: Wire the runner + bump CI Go**

In `cmd/ksec/main.go` `govulncheckRunner`, replace the `cmd.Output()` block (the `out, err := cmd.Output(); if err != nil && len(out) == 0 { return nil, err }; return out, nil`) with stderr capture + the classifier:

```go
	cmd := exec.Command("govulncheck", "-json", "./...")
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, runErr := cmd.Output()
	return classifyGovulncheck(out, stderr.Bytes(), runErr)
```

Ensure `bytes` is imported in `main.go`.

In `.github/workflows/security-dashboard.yaml`, change the setup-go version:

```yaml
      - uses: actions/setup-go@v5
        with:
          go-version: "stable"
```

- [ ] **Step 6: Build, vet, test, YAML-validate, commit**

Run: `go build ./... && go vet ./... && go test ./...` and
`python3 -c "import yaml; yaml.safe_load(open('.github/workflows/security-dashboard.yaml'))" && echo OK`.
```bash
git add internal/collect/govulncheck_result.go internal/collect/govulncheck_result_test.go cmd/ksec/main.go .github/workflows/security-dashboard.yaml
git commit -m "fix(collect): surface govulncheck build failures; CI uses latest Go"
```

---

### Task 2: Reachability filter + OSV severity

**Files:** Modify `internal/collect/source.go`, `internal/collect/source_test.go`.

**Interfaces:**
- Produces: `func severityFromOSV(databaseSpecificSeverity string) string` — normalizes an OSV severity string; empty → `"high"`.
- Behavior change in `SourceCVE.Collect`: keep only findings whose `trace[0].function != ""` (reachable/called); set `Severity` from the finding's OSV severity.

- [ ] **Step 1: Extend the govulncheck JSON structs**

In `source.go`, add `Severity` capture to the OSV struct and `Function` to the trace frame. In `govulnLine.OSV`, add:

```go
		DatabaseSpecific *struct {
			Severity string `json:"severity"`
		} `json:"database_specific"`
```

In `govulnLine.Finding.Trace[]`, add `Function string \`json:"function"\``.

- [ ] **Step 2: Write the failing test**

Add to `source_test.go` (uses the existing test pattern — a Runner returning fixture JSON lines):

```go
func TestSourceCVEReachabilityAndSeverity(t *testing.T) {
	// One reachable HIGH finding (trace has a function) and one non-reachable
	// finding (no function) for a different module — only the reachable one survives.
	lines := []string{
		`{"osv":{"id":"GO-2024-1","summary":"reachable bug","aliases":["CVE-2024-1"],"database_specific":{"severity":"HIGH"},"affected":[{"package":{"name":"example.com/m"},"ranges":[{"events":[{"fixed":"1.2.3"}]}]}]}}`,
		`{"osv":{"id":"GO-2024-2","summary":"imported only","database_specific":{"severity":"LOW"}}}`,
		`{"finding":{"osv":"GO-2024-1","trace":[{"module":"example.com/m","version":"1.0.0","function":"Vuln"}]}}`,
		`{"finding":{"osv":"GO-2024-2","trace":[{"module":"example.com/other","version":"2.0.0"}]}}`,
	}
	c := SourceCVE{Runner: func(state.Repo) ([]byte, error) { return []byte(strings.Join(lines, "\n")), nil }}
	out, err := c.Collect(state.Repo{Repo: "o/r"})
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "example.com/m", out[0].Package)
	assert.Equal(t, "high", out[0].Severity)
	assert.Equal(t, "1.2.3", out[0].FixedVersion)
}

func TestSeverityFromOSV(t *testing.T) {
	assert.Equal(t, "critical", severityFromOSV("CRITICAL"))
	assert.Equal(t, "high", severityFromOSV("HIGH"))
	assert.Equal(t, "medium", severityFromOSV("MODERATE"))
	assert.Equal(t, "medium", severityFromOSV("MEDIUM"))
	assert.Equal(t, "low", severityFromOSV("LOW"))
	assert.Equal(t, "high", severityFromOSV("")) // reachable default
}
```

Ensure `strings` is imported in `source_test.go`.

- [ ] **Step 3: Run it — expect FAIL.** Run: `go test ./internal/collect/...`

- [ ] **Step 4: Implement**

In `source.go`:
- Capture severity into the `adv` struct: add a `severity` field; when reading `gl.OSV`, set `a.severity = ""`; if `gl.OSV.DatabaseSpecific != nil` set `a.severity = gl.OSV.DatabaseSpecific.Severity`.
- In the findings loop, **skip non-reachable**: `if gl.Finding.Trace[0].Function == "" { continue }` (after the existing `len(...)==0` guard).
- Set `Severity: severityFromOSV(a.severity)` on the `state.Finding` (replacing the hardcoded `"unknown"`).

Add the helper:

```go
func severityFromOSV(s string) string {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "CRITICAL":
		return "critical"
	case "HIGH":
		return "high"
	case "MODERATE", "MEDIUM":
		return "medium"
	case "LOW":
		return "low"
	default:
		return "high" // reachable vuln with no severity data is actionable
	}
}
```

Add `"strings"` to `source.go` imports.

- [ ] **Step 5: Run it — expect PASS, then build/test all.** Run: `go test ./internal/collect/... && go build ./... && go test ./...`

- [ ] **Step 6: Commit**

```bash
git add internal/collect/source.go internal/collect/source_test.go
git commit -m "feat(collect): keep only reachable govulncheck findings with real severities"
```

---

### Task 3: Tracked-PR list (PRs out of findings)

**Files:** Modify `internal/state/types.go`, `internal/state/files.go`, `internal/collect/prs.go`, `internal/collect/prs_test.go`, `cmd/ksec/main.go`.

**Interfaces:**
- Produces: `state.TrackedPR{Repo string; Number int; Title, Author, URL, Source string}`; `state.OpenPRsFile` const (`"openprs.json"`); `func collect.OpenPRs(repos []state.Repo, gh ghclient.GitHub) ([]state.TrackedPR, []state.CollectionError)`.
- `collect.PRs` (the finding `Collector`) is removed from the findings collector list.

- [ ] **Step 1: Add the state type + file const**

In `internal/state/types.go`:

```go
type TrackedPR struct {
	Repo   string `json:"repo"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"url"`
	Source string `json:"source"` // renovate|dependabot|ksec|human
}
```

In `internal/state/files.go`, add alongside the other file-name consts:

```go
	OpenPRsFile = "openprs.json"
```

(Match the existing const style — confirm whether they are grouped in a `const (...)` block and add accordingly.)

- [ ] **Step 2: Write the failing test**

Rewrite `internal/collect/prs_test.go` to test `OpenPRs` (replacing the old `Collect` test):

```go
package collect

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakePRGH struct {
	byRepo map[string][]ghclient.PullRequest
	ghclient.GitHub
}

func (f fakePRGH) ListOpenPRs(repo string) ([]ghclient.PullRequest, error) {
	return f.byRepo[repo], nil
}

func TestOpenPRsTracksAndClassifies(t *testing.T) {
	gh := fakePRGH{byRepo: map[string][]ghclient.PullRequest{
		"o/r": {
			{Number: 2, Title: "bump y", Author: "dependabot[bot]", URL: "u2"},
			{Number: 1, Title: "feature", Author: "alice"},                       // not tracked (no bot, no label)
			{Number: 3, Title: "sec fix", Author: "bob", Labels: []string{"security"}, URL: "u3"},
		},
	}}
	prs, errs := OpenPRs([]state.Repo{{Repo: "o/r"}}, gh)
	require.Empty(t, errs)
	require.Len(t, prs, 2)
	// sorted by repo then number
	assert.Equal(t, 2, prs[0].Number)
	assert.Equal(t, "dependabot", prs[0].Source)
	assert.Equal(t, "u2", prs[0].URL)
	assert.Equal(t, 3, prs[1].Number)
	assert.Equal(t, "human", prs[1].Source)
}
```

- [ ] **Step 3: Run it — expect FAIL.** Run: `go test ./internal/collect/...`

- [ ] **Step 4: Implement**

Rewrite `internal/collect/prs.go`: drop the `Collector` methods (`Name`/`Collect`) and add `OpenPRs`. Keep `isSecurityPR` and the author/label sets. Classify source inline (renovate[bot]→renovate, dependabot[bot]→dependabot, kairos-security-bot→ksec, else human):

```go
package collect

import (
	"sort"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

var botAuthors = map[string]bool{
	"renovate[bot]": true, "dependabot[bot]": true, "kairos-security-bot": true,
}
var secLabels = map[string]bool{"security": true, "dependencies": true}

// OpenPRs lists the tracked open PRs (bot-authored or security/dependencies
// labelled) across repos for the dashboard's PR list. These are remediation
// artifacts, NOT security findings, so they no longer enter the findings set.
func OpenPRs(repos []state.Repo, gh ghclient.GitHub) ([]state.TrackedPR, []state.CollectionError) {
	var out []state.TrackedPR
	var errs []state.CollectionError
	for _, repo := range repos {
		prs, err := gh.ListOpenPRs(repo.Repo)
		if err != nil {
			errs = append(errs, state.CollectionError{Repo: repo.Repo, Collector: "prs", Message: err.Error()})
			continue
		}
		for _, pr := range prs {
			if !isSecurityPR(pr) {
				continue
			}
			out = append(out, state.TrackedPR{
				Repo: repo.Repo, Number: pr.Number, Title: pr.Title,
				Author: pr.Author, URL: pr.URL, Source: prSource(pr.Author),
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

func prSource(author string) string {
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

func isSecurityPR(pr ghclient.PullRequest) bool {
	if botAuthors[pr.Author] {
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

- [ ] **Step 5: Wire the collect command**

In `cmd/ksec/main.go` `newCollectCmd`: remove `collect.PRs{GH: gh}` from the `collectors` slice; after `state.Save(... FindingsFile, out)`, gather and persist PRs:

```go
			if err := state.Save(gf.stateDir, state.FindingsFile, out); err != nil {
				return err
			}
			prs, prErrs := collect.OpenPRs(repos, gh)
			out.Errors = append(out.Errors, prErrs...)
			if len(prErrs) > 0 {
				_ = state.Save(gf.stateDir, state.FindingsFile, out) // include PR-list errors
			}
			return state.Save(gf.stateDir, state.OpenPRsFile, prs)
```

- [ ] **Step 6: Run + build + commit**

Run: `go test ./internal/collect/... ./internal/state/... && go build ./... && go test ./...`
```bash
git add internal/state/types.go internal/state/files.go internal/collect/prs.go internal/collect/prs_test.go cmd/ksec/main.go
git commit -m "feat(collect): track open PRs in openprs.json; remove PRs from findings"
```

---

### Task 4: Render links + Open PRs section

**Files:** Modify `internal/render/render.go`, `internal/render/render_test.go`, `internal/render/html.go`, `cmd/ksec/main.go`, `internal/render/testdata/`.

**Interfaces:**
- Consumes: `state.TrackedPR`, `state.Correlated.Findings` (for the id→finding lookup).
- `render.Input` gains `OpenPRs []state.TrackedPR \`json:"openPRs,omitempty"\``.
- `DashboardMarkdown` renders Focus as title→URL links and adds a "📋 Open PRs" section.

- [ ] **Step 1: Write the failing test**

Add to `internal/render/render_test.go`:

```go
func TestFocusShowsTitleLinkNotID(t *testing.T) {
	in := Input{
		Correlated: state.Correlated{Findings: []state.Finding{
			{ID: "abc123", Repo: "o/r", Title: "x/net rapid reset", URL: "https://github.com/o/r/pull/9"},
			{ID: "def456", Repo: "o/r", Type: "sourceCVE", CVEID: "GO-2024-3218", Title: "kad-dht issue"},
		}},
		Triage: state.Triage{Focus: []string{"abc123", "def456"}},
	}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "[x/net rapid reset](https://github.com/o/r/pull/9)")
	assert.Contains(t, md, "[kad-dht issue](https://pkg.go.dev/vuln/GO-2024-3218)")
	assert.NotContains(t, md, "abc123") // raw id never shown
}

func TestOpenPRsSection(t *testing.T) {
	in := Input{OpenPRs: []state.TrackedPR{
		{Repo: "o/r", Number: 7, Title: "bump foo", URL: "https://github.com/o/r/pull/7", Source: "dependabot"},
	}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "📋 Open PRs")
	assert.Contains(t, md, "[#7 bump foo](https://github.com/o/r/pull/7)")
	assert.Contains(t, md, "dependabot")
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/render/...`

- [ ] **Step 3: Implement the Focus link + Open PRs section**

In `render.go` `DashboardMarkdown`:
- Before the Focus loop, build `byID := map[string]state.Finding{}` from `in.Correlated.Findings`.
- Replace the Focus loop body with a link renderer using this helper (add it to `render.go`):

```go
func findingLink(f state.Finding) string {
	url := f.URL
	if url == "" && f.Type == "sourceCVE" {
		id := f.CVEID
		if id == "" {
			id = f.GHSA
		}
		if id != "" {
			url = "https://pkg.go.dev/vuln/" + id
		}
	}
	title := f.Title
	if title == "" {
		title = f.Package
	}
	if title == "" {
		title = f.Repo + " finding"
	}
	if url != "" {
		return fmt.Sprintf("[%s](%s)", title, url)
	}
	return title
}
```

Focus loop:

```go
		for _, id := range in.Triage.Focus {
			if f, ok := byID[id]; ok {
				line := findingLink(f)
				if s := in.Triage.Summaries[id]; s != "" && !strings.HasPrefix(s, "Finding in ") {
					line += " — " + s
				}
				fmt.Fprintf(&b, "- %s\n", line)
			} else if s, ok := in.Triage.Summaries[id]; ok {
				fmt.Fprintf(&b, "- %s\n", s)
			}
		}
```

Add the "📋 Open PRs" section (after the per-repo table, before the "🤖 Bot PR ledger" section):

```go
	b.WriteString("## 📋 Open PRs\n\n")
	if len(in.OpenPRs) == 0 {
		b.WriteString("_None._\n\n")
	} else {
		repo := ""
		for _, pr := range in.OpenPRs {
			if pr.Repo != repo {
				repo = pr.Repo
				fmt.Fprintf(&b, "**%s**\n\n", repo)
			}
			link := fmt.Sprintf("#%d %s", pr.Number, pr.Title)
			if pr.URL != "" {
				link = fmt.Sprintf("[#%d %s](%s)", pr.Number, pr.Title, pr.URL)
			}
			fmt.Fprintf(&b, "- %s — %s\n", link, pr.Source)
		}
		b.WriteString("\n")
	}
```

Add `OpenPRs []state.TrackedPR \`json:"openPRs,omitempty"\`` to `Input`.

- [ ] **Step 4: Mirror in HTML**

In `internal/render/html.go`, add the equivalent Focus-link rendering (use `findingLink` to build the same text; in the template, render the resulting links — keep the existing escaping approach for any user text). Add an "📋 Open PRs" block mirroring the markdown, grouped by repo. (Match the file's existing template/section style.)

- [ ] **Step 5: Wire the render command**

In `cmd/ksec/main.go` `newRenderCmd`: load the tracked PRs best-effort and set them on the base `Input` (so all surfaces get them):

```go
			var openPRs []state.TrackedPR
			_ = state.Load(gf.stateDir, state.OpenPRsFile, &openPRs) // best-effort
```

and add `OpenPRs: openPRs,` to the `render.Input{...}` literal (the base `in`; the committed `fileIn` copy keeps it — PR data is deterministic, unlike the AI summary).

- [ ] **Step 6: Regenerate goldens, build, test, commit**

Run: `UPDATE_GOLDEN=1 go test ./internal/render/...` (eyeball the new Focus links + Open PRs section + that no raw ids appear), then `go build ./... && go vet ./... && test -z "$(gofmt -l .)" && go test ./...`.
```bash
git add internal/render/render.go internal/render/render_test.go internal/render/html.go internal/render/testdata/ cmd/ksec/main.go
git commit -m "feat(render): link finding titles in Focus; add Open PRs section"
```

---

## Self-review

**Spec coverage:**
- A (CI Go + loud build failures) → Task 1. ✓
- B (reachability filter + OSV severity) → Task 2. ✓
- C (PRs → openprs.json, out of findings) → Task 3. ✓
- D (Focus title/URL links, no hashes) + 📋 Open PRs section → Task 4. ✓
- Determinism (sorted findings/PRs; PR data on committed copy) → Tasks 3, 4. ✓
- stdlib → toolchain path preserved (severity change doesn't touch `Package` derivation) → Task 2. ✓

**Placeholder scan:** none — every step has concrete code/commands. The two integration touch points (govulncheckRunner shell, html.go template) reference exact existing structures.

**Type consistency:** `classifyGovulncheck` (Task 1) used by the runner. `severityFromOSV` + the trace `Function`/OSV `DatabaseSpecific` fields (Task 2) used in `SourceCVE.Collect`. `state.TrackedPR` + `OpenPRsFile` + `collect.OpenPRs` (Task 3) consumed by the render command + `Input.OpenPRs` (Task 4). `findingLink` + `byID` (Task 4) use `state.Finding`/`state.Correlated` fields that exist today. `collect.PRs` removal (Task 3) is the only deletion — confirmed its sole caller is the collectors slice in `newCollectCmd`.

---

## Operational notes

- After Task 1 ships, the next scheduled run will (correctly) start reporting Go-toolchain/stdlib findings for repos built against an older toolchain than latest `stable`; those flow to Plan-4c's toolchain-bump path. Expect edgevpn's `GO-2024-3218` (no upstream fix → needs-human) plus the stdlib set.
- `severityFromOSV` defaulting to `high` may over-state a few low-severity reachable vulns; this is deliberate (reachable = actionable) and preferable to hiding them under `unknown`. Revisit with CVSS parsing only if it proves noisy.
- The "📋 Open PRs" list uses the `isSecurityPR` filter (bot/labelled), so it shows remediation-relevant PRs, not every open PR; widen the filter later if a fuller list is wanted.
