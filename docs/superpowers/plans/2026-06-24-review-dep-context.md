# Bot-PR Review — Dependency-Change Context + Comment Upsert Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Feed the bot-PR assessor the real dependency changes (PR changelog body + upstream source diff X..Y) so its verdict is informed and it returns a `changesSummary`; surface that summary; and **upsert** the review comment (one per PR, edited in place — no spam).

**Architecture:** ghclient gains `PullRequest.Body`, `CompareDiff` (GitHub compare, diff media type), and `UpsertPRComment` (find-by-marker edit-or-create). `internal/review` gains pure `parseBumps`/`moduleRepo`; `Run` assembles per-PR context (body + per-bump upstream diff + PR diff) and the `Assessor` returns `{verdict, reasoning, changesSummary}`. `state.PRReview` + render show the summary. Builds on the just-shipped review feature.

**Tech Stack:** Go 1.22, `stretchr/testify`, `gh` CLI (compare diff media type, REST issue comments), LocalAI forced tool call.

## Global Constraints

- Module `github.com/kairos-io/security`; Go 1.22.
- Forced tool-call verdict (enum stays grammar-constrained); every assessor error path → `("needs_human_verification", reason, "", nil)`.
- GitHub-compare is the source-diff mechanism; unresolvable modules / compare failures **degrade silently** to body+PR-diff (never fail the run). Context capped (~40KB/bump, ~60KB total).
- **Comment upsert**: exactly one `<!-- ksec:review -->` comment per PR, edited in place — never spam. Idempotency on head SHA still gates whether we re-assess.
- Dry-run zero writes (no compare fetch needed beyond read; no comment/approve). Deterministic dashboard.

---

## File structure

```
internal/ghclient/ghclient.go   # PullRequest.Body; CompareDiff; UpsertPRComment; interface (modify)
internal/ghclient/fake.go       # Fake CompareDiff/UpsertPRComment (modify)
internal/review/deps.go         # parseBumps + moduleRepo (create)
internal/review/deps_test.go    # (create)
internal/review/assessor.go     # Assessor sig (+context,+changesSummary); FakeAssessor (modify)
internal/review/openai.go       # context-aware prompt + changesSummary tool field (modify)
internal/review/openai_test.go  # (modify)
internal/review/run.go          # context assembly + upsert comment + changesSummary (modify)
internal/review/run_test.go     # (modify)
internal/state/types.go         # PRReview.ChangesSummary (modify)
internal/render/render.go       # reviews row shows changesSummary (modify)
internal/render/render_test.go  # (modify)
internal/render/html.go         # mirror (modify)
internal/render/testdata/       # regenerated goldens
```

---

### Task 1: ghclient — Body, CompareDiff, UpsertPRComment

**Files:** Modify `internal/ghclient/ghclient.go`, `internal/ghclient/fake.go`.

**Interfaces:** `PullRequest.Body string` (json `body`); `CompareDiff(repo, base, head string) ([]byte, error)`; `UpsertPRComment(repo string, pr int, marker, body string) error`. All on `GitHub` + `CLI` + `FakeGitHub`.

- [ ] **Step 1: PullRequest.Body + ListOpenPRs projection** — add `Body string \`json:"body"\`` to `PullRequest`; add `body` to `ListOpenPRs`'s `--json` and `body: .body` to the `-q` map.

- [ ] **Step 2: CLI methods**

```go
func (c *CLI) CompareDiff(repo, base, head string) ([]byte, error) {
	return c.run("api", fmt.Sprintf("repos/%s/compare/%s...%s", repo, base, head),
		"-H", "Accept: application/vnd.github.diff")
}

// UpsertPRComment edits our existing marker comment in place (no spam), or
// creates one. Uses REST issue comments so the numeric id matches the PATCH
// endpoint.
func (c *CLI) UpsertPRComment(repo string, pr int, marker, body string) error {
	listed, err := c.run("api", fmt.Sprintf("repos/%s/issues/%d/comments?per_page=100", repo, pr),
		"-q", "[.[] | {id, body}]")
	if err == nil && len(bytes.TrimSpace(listed)) > 0 {
		var cs []struct {
			ID   int64  `json:"id"`
			Body string `json:"body"`
		}
		_ = json.Unmarshal(listed, &cs)
		for _, cm := range cs {
			if strings.Contains(cm.Body, marker) {
				_, err := c.run("api", "-X", "PATCH",
					fmt.Sprintf("repos/%s/issues/comments/%d", repo, cm.ID), "-f", "body="+body)
				return err
			}
		}
	}
	_, err = c.run("api", fmt.Sprintf("repos/%s/issues/%d/comments", repo, pr), "-f", "body="+body)
	return err
}
```
Add both to the `GitHub` interface. (`bytes`/`encoding/json`/`fmt`/`strings` already imported.)

- [ ] **Step 3: FakeGitHub** — implement both. `CompareDiff` returns a configurable `Compares map["repo:base...head"][]byte` (nil ok). `UpsertPRComment` records `repo#pr` into an `Upserted []string` and, to let tests assert no-spam, key by `repo#pr` in a map so a second upsert overwrites rather than appends (e.g. `Upserts map[string]string` keyed `repo#pr` → body). Match the Fake's existing style.

- [ ] **Step 4: Build + vet + test + commit**

Run: `go build ./... && go vet ./... && go test ./...` (the new interface methods require the Fake to implement them).
```bash
git add internal/ghclient/ghclient.go internal/ghclient/fake.go
git commit -m "feat(ghclient): PullRequest.Body + CompareDiff + UpsertPRComment"
```

---

### Task 2: review pure helpers — parseBumps + moduleRepo

**Files:** Create `internal/review/deps.go`, `internal/review/deps_test.go`.

**Interfaces:** `type DepBump struct{ Module, From, To string }`; `func parseBumps(diff []byte) []DepBump`; `func moduleRepo(mod string) (string, bool)`.

- [ ] **Step 1: Write the failing tests** — `deps_test.go`:

```go
package review

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBumps(t *testing.T) {
	diff := []byte("--- a/go.mod\n+++ b/go.mod\n" +
		"-\tgolang.org/x/net v0.30.0\n+\tgolang.org/x/net v0.33.0\n" +
		"-\tgithub.com/foo/bar v1.2.0 // indirect\n+\tgithub.com/foo/bar v1.3.0 // indirect\n")
	bumps := parseBumps(diff)
	assert.Equal(t, []DepBump{
		{Module: "golang.org/x/net", From: "0.30.0", To: "0.33.0"},
		{Module: "github.com/foo/bar", From: "1.2.0", To: "1.3.0"},
	}, bumps)
}

func TestModuleRepo(t *testing.T) {
	cases := map[string]struct {
		repo string
		ok   bool
	}{
		"github.com/foo/bar":        {"foo/bar", true},
		"github.com/foo/bar/v2":     {"foo/bar", true},
		"github.com/foo/bar/sub":    {"foo/bar", true},
		"golang.org/x/net":          {"golang/net", true},
		"k8s.io/api":                {"kubernetes/api", true},
		"sigs.k8s.io/yaml":          {"kubernetes-sigs/yaml", true},
		"example.com/vanity/thing":  {"", false},
	}
	for mod, want := range cases {
		got, ok := moduleRepo(mod)
		assert.Equal(t, want.ok, ok, mod)
		assert.Equal(t, want.repo, got, mod)
	}
}
```

- [ ] **Step 2: Run red.** `go test ./internal/review/...`

- [ ] **Step 3: Implement** — `deps.go`:

```go
package review

import (
	"regexp"
	"strings"
)

type DepBump struct{ Module, From, To string }

var reModLine = regexp.MustCompile(`^[+-]\s+(\S+)\s+v(\S+)`)

// parseBumps extracts {module, from, to} from a PR's go.mod diff by pairing the
// "-" old and "+" new version lines for the same module.
func parseBumps(diff []byte) []DepBump {
	from := map[string]string{}
	var order []string
	to := map[string]string{}
	for _, line := range strings.Split(string(diff), "\n") {
		m := reModLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		mod, ver := m[1], strings.TrimSuffix(m[2], " //") // tolerate trailing tokens
		ver = strings.Fields(ver)[0]
		if strings.HasPrefix(line, "-") {
			from[mod] = ver
		} else {
			if _, seen := to[mod]; !seen {
				order = append(order, mod)
			}
			to[mod] = ver
		}
	}
	var out []DepBump
	for _, mod := range order {
		if f, ok := from[mod]; ok && f != to[mod] {
			out = append(out, DepBump{Module: mod, From: f, To: to[mod]})
		}
	}
	return out
}

// moduleRepo maps a Go module path to a GitHub "owner/name", handling the
// common hosts. Unresolvable (vanity) paths return ok=false (degrade).
func moduleRepo(mod string) (string, bool) {
	// strip a /vN major-version suffix
	parts := strings.Split(mod, "/")
	switch {
	case parts[0] == "github.com" && len(parts) >= 3:
		return parts[1] + "/" + parts[2], true
	case mod == "golang.org/x/"+strings.TrimPrefix(mod, "golang.org/x/") && strings.HasPrefix(mod, "golang.org/x/") && len(parts) >= 3:
		return "golang/" + parts[2], true
	case strings.HasPrefix(mod, "k8s.io/") && len(parts) >= 2:
		return "kubernetes/" + parts[1], true
	case strings.HasPrefix(mod, "sigs.k8s.io/") && len(parts) >= 2:
		return "kubernetes-sigs/" + parts[1], true
	}
	return "", false
}
```

(If the `golang.org/x` case above reads awkwardly, simplify to `strings.HasPrefix(mod, "golang.org/x/") && len(parts) >= 3 → "golang/"+parts[2]`. Keep the behavior in the test.)

- [ ] **Step 4: Run green + commit**

Run: `go test ./internal/review/... && go build ./...`
```bash
git add internal/review/deps.go internal/review/deps_test.go
git commit -m "feat(review): parse go.mod bumps + map module to GitHub repo"
```

---

### Task 3: Context assembly + assessor changesSummary + comment upsert

**Files:** Modify `internal/review/assessor.go`, `internal/review/openai.go`, `internal/review/openai_test.go`, `internal/review/run.go`, `internal/review/run_test.go`.

**Interfaces:** `Assessor.Assess(pr ghclient.PullRequest, context string) (verdict, reasoning, changesSummary string, err error)`; `Run` builds the context (body + per-bump upstream diff + PR diff) and upserts the comment with the summary.

- [ ] **Step 1: Update the Assessor interface + FakeAssessor** (`assessor.go`):

```go
type Assessor interface {
	Assess(pr ghclient.PullRequest, context string) (verdict, reasoning, changesSummary string, err error)
}

type FakeAssessor struct {
	Verdict, Reasoning, ChangesSummary string
	Err                                error
	GotContext                         string // records the assembled context for assertions
}

func (f *FakeAssessor) Assess(_ ghclient.PullRequest, context string) (string, string, string, error) {
	f.GotContext = context
	return f.Verdict, f.Reasoning, f.ChangesSummary, f.Err
}
```

- [ ] **Step 2: OpenAIAssessor** (`openai.go`): add `changesSummary` to the `assess_pr` schema + `assessArgs`; change `Assess` to the new signature; build the prompt from the passed `context` (no longer the raw diff) instructing the model to also return a 1-3 sentence `changesSummary` of the dependency change. Every degrade path returns `(verdictNeeds, reason, "", nil)`. Update `openai_test.go` for the 4-return signature + a `changesSummary` assertion in the happy path.

- [ ] **Step 3: Run context assembly + upsert** (`run.go`): add a context cap consts (`maxBumpDiff = 40000`, `maxContext = 60000`). For each assessed PR:

```go
		// Assemble the assessment context: changelog (PR body) + upstream source
		// diffs for each bump + the PR's own diff.
		var ctx strings.Builder
		if strings.TrimSpace(pr.Body) != "" {
			ctx.WriteString("PR description / changelog:\n" + pr.Body + "\n\n")
		}
		diff, derr := gh.PRDiff(repo.Repo, pr.Number)
		if derr != nil {
			errs = append(errs, state.CollectionError{Repo: repo.Repo, Collector: "review", Message: derr.Error()})
			continue
		}
		for _, b := range parseBumps(diff) {
			if ctx.Len() > maxContext {
				break
			}
			gr, ok := moduleRepo(b.Module)
			if !ok {
				continue
			}
			ud, uerr := gh.CompareDiff(gr, "v"+b.From, "v"+b.To)
			if uerr != nil || len(ud) == 0 {
				continue // degrade: no upstream source diff for this bump
			}
			if len(ud) > maxBumpDiff {
				ud = ud[:maxBumpDiff]
			}
			fmt.Fprintf(&ctx, "Upstream %s %s..%s:\n%s\n\n", b.Module, b.From, b.To, ud)
		}
		ctx.WriteString("PR diff:\n" + string(diff))
		verdict, reasoning, summary, _ := a.Assess(pr, ctx.String())
		rv := state.PRReview{Repo: repo.Repo, PR: pr.Number, URL: pr.URL, HeadSHA: pr.HeadSHA,
			Verdict: verdict, Reasoning: reasoning, ChangesSummary: summary, ReviewedRun: runID}
		out = append(out, rv)
		if dryRun {
			fmt.Printf("[dry-run] would comment on %s#%d: %s — %s\n", repo.Repo, pr.Number, verdict, summary)
			continue
		}
		_ = gh.UpsertPRComment(repo.Repo, pr.Number, reviewMarker, comment(rv, cfg.Notify))
		if cfg.AutoApprove && verdict == "good" {
			_ = gh.ApprovePR(repo.Repo, pr.Number, "kairos-security: automated review verdict good")
		}
```

Update `comment(r, notify)` to include the summary: after the verdict/reasoning line, `if r.ChangesSummary != "" { b.WriteString("\n\n**Dependency changes:** " + r.ChangesSummary) }`, then cc + marker. (Remove the old `gh.PostPRComment` review path.)

- [ ] **Step 4: Update run_test.go** — fakeGH gains `PRDiff`/`CompareDiff`/`UpsertPRComment` (record upserts in a map keyed `repo#pr` so a re-upsert overwrites — proving no-spam); FakeAssessor uses the 4-return form; assert: context contains the PR body + the upstream-diff label when CompareDiff returns data; `ChangesSummary` recorded; the comment is **upserted once** per PR (second assessment with a changed head edits, map still has one entry); CompareDiff error still assesses (degrade); dry-run no upsert.

- [ ] **Step 5: Run green + build + commit**

Run: `go test ./internal/review/... && go build ./... && go vet ./...`
```bash
git add internal/review/assessor.go internal/review/openai.go internal/review/openai_test.go internal/review/run.go internal/review/run_test.go
git commit -m "feat(review): assess with dep-change context (body+upstream diff); changesSummary; upsert comment"
```

---

### Task 4: PRReview.ChangesSummary + dashboard

**Files:** Modify `internal/state/types.go`, `internal/render/render.go`, `internal/render/render_test.go`, `internal/render/html.go`, `internal/render/testdata/`.

- [ ] **Step 1: state** — add `ChangesSummary string \`json:"changesSummary,omitempty"\`` to `state.PRReview`.

- [ ] **Step 2: render** — in the "🔎 Bot-PR reviews" rows, when `r.ChangesSummary != ""`, append it (e.g. a sub-line `  ↳ <changesSummary>` or extend the line). Add a `render_test.go` assertion that a review with `ChangesSummary` renders it. Mirror in `html.go` (escaped).

- [ ] **Step 3: Regenerate goldens, build, vet, gofmt, test, commit**

Run: `UPDATE_GOLDEN=1 go test ./internal/render/...` (eyeball: summary shown; deterministic; no raw 64-hex id), re-run; then `go build ./... && go vet ./... && test -z "$(gofmt -l .)" && go test ./...`.
```bash
git add internal/state/types.go internal/render/ 
git commit -m "feat(render): show dependency-change summary in Bot-PR reviews"
```

---

## Self-review

**Spec coverage:**
- PR body + upstream source diff X..Y fed to the assessor → Tasks 1 (Body, CompareDiff), 2 (parse/resolve), 3 (assembly). ✓
- `changesSummary` from the forced tool call → Task 3. ✓
- Comment **upsert** (one per PR, edited; no spam) → Tasks 1 (UpsertPRComment), 3 (use it). ✓
- Summary in comment + dashboard → Tasks 3, 4. ✓
- Degrade (unresolvable module / compare failure / non-GitHub) → Task 3 (skip bump, still assess). ✓
- Dry-run no writes; verdict enum + needs_human default preserved → Task 3. ✓

**Placeholder scan:** none — full code for ghclient methods, the pure helpers, the context-assembly + upsert in Run; openai.go's prompt/schema change is described against the existing forced-tool-call structure.

**Type consistency:** `PullRequest.Body`/`CompareDiff`/`UpsertPRComment` (Task 1) used by `Run` (Task 3). `parseBumps`/`moduleRepo` (Task 2) used by `Run` (Task 3). `Assessor` 4-return signature (Task 3) implemented by `OpenAIAssessor`+`FakeAssessor`. `PRReview.ChangesSummary` (Task 4) set by `Run`, rendered by Task 4.

---

## Operational notes

- The assessor now sees the real X..Y source diff for GitHub-hosted deps (≈all here); a vanity/non-GitHub module degrades to changelog + PR-diff context (the `changesSummary` will say so).
- One review comment per PR, edited in place — re-assessment (on head change) updates it; no comment pile-up.
- Context is capped (~40KB/bump, ~60KB total) so a large dependency diff can't blow up the prompt.
- Dry-run prints `verdict — summary` and writes nothing.
