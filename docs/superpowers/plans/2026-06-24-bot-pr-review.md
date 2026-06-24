# AI Bot-PR Review Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** `ksec review` — for each open bot-authored PR on a tracked repo, fetch the diff, get an AI supply-chain verdict (good / bad / needs_human_verification via a forced tool call), record it, comment on the PR (cc-ing configured handles), and optionally auto-approve a `good` PR. Idempotent on the PR head SHA.

**Architecture:** `ai.yaml` gains a `review` config; ghclient gains `PRDiff`/`ApprovePR` + `PullRequest.HeadSHA`; a new `internal/review` package holds the `Assessor` (forced-tool-call `OpenAIAssessor` + `FakeAssessor`) and `Run` orchestration producing `[]state.PRReview`; a `ksec review` command writes `state/reviews.json`; the dashboard renders a "🔎 Bot-PR reviews" section. Builds on `internal/config`, `internal/ghclient`, `internal/state`, `internal/triage` (tool-call pattern), `internal/render`, `cmd/ksec`.

**Tech Stack:** Go 1.22, `stretchr/testify`, `gh` CLI, LocalAI (OpenAI-compatible forced tool calling).

## Global Constraints

- Module `github.com/kairos-io/security`; Go 1.22.
- Forced tool call only — the verdict is grammar-constrained to `{good, bad, needs_human_verification}`; never trust freeform model JSON. Any AI error/unparseable → `needs_human_verification` (safe default, never a hard failure).
- **Idempotent on head SHA:** a PR is (re)assessed/commented only when its `HeadSHA` changes — no re-spam, no re-spend.
- Bots only (`pr.IsBot`). Auto-approve only on `good` and only when configured; `bad`/`needs_human` never approve. Auto-**merge** is out of scope.
- Dry-run performs zero writes (no comment, no approve) — prints intended actions. Token redacted via the existing `gh` client.
- Deterministic committed artifacts; no raw SHA-256 finding id in human output.

---

## File structure

```
internal/config/config.go            # ReviewCfg + AIConfig.Review + LoadAI default (modify)
internal/config/config_test.go       # (modify)
internal/ghclient/ghclient.go        # PullRequest.HeadSHA; PRDiff; ApprovePR; interface (modify)
internal/ghclient/fake.go            # FakeGitHub PRDiff/ApprovePR (modify)
internal/state/types.go              # PRReview + ReviewsFile (modify)
internal/state/ledger_test.go        # PRReview round-trip (modify)
internal/review/assessor.go          # Assessor + FakeAssessor (create)
internal/review/openai.go            # OpenAIAssessor (forced tool call) (create)
internal/review/openai_test.go       # (create)
internal/review/run.go               # Run orchestration (create)
internal/review/run_test.go          # (create)
internal/render/render.go            # "🔎 Bot-PR reviews" section + Input.Reviews (modify)
internal/render/render_test.go       # (modify)
internal/render/html.go              # mirror (modify)
internal/render/testdata/            # regenerated goldens
cmd/ksec/main.go                     # ksec review command (modify)
.github/workflows/security-dashboard.yaml  # ksec review step (modify)
```

---

### Task 1: Review config

**Files:** Modify `internal/config/config.go`, `internal/config/config_test.go`.

**Interfaces:** `type ReviewCfg struct { Enabled bool; AutoApprove bool; MaxPerRun int; Notify []string }`; `AIConfig.Review ReviewCfg` (yaml `review`); `LoadAI` defaults `MaxPerRun` to 20 when ≤0.

- [ ] **Step 1: Write the failing test** — add to `internal/config/config_test.go` (match the existing LoadAI test style; if AI config is loaded from a temp yaml, add a `review:` block):

```go
func TestLoadAIReview(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ai.yaml")
	require.NoError(t, os.WriteFile(p, []byte(`
localai:
  endpoint: http://localhost:8080
  model:
    name: m
review:
  enabled: true
  autoApprove: true
  notify: ["@team"]
`), 0o644))
	cfg, err := LoadAI(p)
	require.NoError(t, err)
	assert.True(t, cfg.Review.Enabled)
	assert.True(t, cfg.Review.AutoApprove)
	assert.Equal(t, []string{"@team"}, cfg.Review.Notify)
	assert.Equal(t, 20, cfg.Review.MaxPerRun) // defaulted
}
```

(Use the imports the existing config_test uses — `os`, `path/filepath`, testify. If the file lacks them, add.)

- [ ] **Step 2: Run red.** `go test ./internal/config/...`

- [ ] **Step 3: Implement** — in `config.go`:

```go
type ReviewCfg struct {
	Enabled     bool     `yaml:"enabled"`
	AutoApprove bool     `yaml:"autoApprove"`
	MaxPerRun   int      `yaml:"maxPerRun"`
	Notify      []string `yaml:"notify"`
}
```
Add `Review ReviewCfg \`yaml:"review"\`` to `AIConfig`. In `LoadAI`, after the existing Nib defaulting, add:
```go
	if cfg.Review.MaxPerRun <= 0 {
		cfg.Review.MaxPerRun = 20
	}
```

- [ ] **Step 4: Run green + commit**

Run: `go test ./internal/config/... && go build ./...`
```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "feat(config): review config (enabled/autoApprove/maxPerRun/notify)"
```

---

### Task 2: ghclient — HeadSHA, PRDiff, ApprovePR

**Files:** Modify `internal/ghclient/ghclient.go`, `internal/ghclient/fake.go`.

**Interfaces:** `PullRequest.HeadSHA string` (json `headSHA`, from `headRefOid`); `PRDiff(repo string, pr int) ([]byte, error)`; `ApprovePR(repo string, pr int, body string) error`. All added to the `GitHub` interface, `CLI`, and `FakeGitHub`.

- [ ] **Step 1: Add `HeadSHA` to PullRequest + the ListOpenPRs projection**

In `ghclient.go`: add `HeadSHA string \`json:"headSHA"\`` to `PullRequest`; in `ListOpenPRs`, add `headRefOid` to the `--json` list and `headSHA: .headRefOid` to the `-q` projection.

- [ ] **Step 2: Add the interface methods + CLI impls**

In the `GitHub` interface add:
```go
	PRDiff(repo string, pr int) ([]byte, error)
	ApprovePR(repo string, pr int, body string) error
```
CLI impls:
```go
func (c *CLI) PRDiff(repo string, pr int) ([]byte, error) {
	return c.run("pr", "diff", fmt.Sprint(pr), "-R", repo)
}

func (c *CLI) ApprovePR(repo string, pr int, body string) error {
	_, err := c.run("pr", "review", fmt.Sprint(pr), "-R", repo, "--approve", "--body", body)
	return err
}
```

- [ ] **Step 3: Add to `FakeGitHub`** (`internal/ghclient/fake.go`) — fields + methods so tests can drive them, e.g.:
```go
	Diffs       map[string][]byte // "repo#pr" -> diff
	Approved    []string          // "repo#pr" recorded
```
`PRDiff` returns `Diffs[key]`/nil; `ApprovePR` appends to `Approved`. (Match the Fake's existing construction/style; if it uses a different recording pattern, follow it.)

- [ ] **Step 4: Build + vet + test + commit**

Run: `go build ./... && go vet ./... && go test ./...` (existing ghclient/consumers still compile — the new interface methods require the Fake to implement them).
```bash
git add internal/ghclient/ghclient.go internal/ghclient/fake.go
git commit -m "feat(ghclient): PullRequest.HeadSHA + PRDiff + ApprovePR"
```

---

### Task 3: state.PRReview

**Files:** Modify `internal/state/types.go`, `internal/state/ledger_test.go`.

**Interfaces:** `state.PRReview{Repo string; PR int; URL, HeadSHA, Verdict, Reasoning, ReviewedRun string}`; `ReviewsFile = "reviews.json"` const.

- [ ] **Step 1: Add the type + const + round-trip test**

In `types.go`:
```go
type PRReview struct {
	Repo        string `json:"repo"`
	PR          int    `json:"pr"`
	URL         string `json:"url,omitempty"`
	HeadSHA     string `json:"headSHA"`
	Verdict     string `json:"verdict"` // good | bad | needs_human_verification
	Reasoning   string `json:"reasoning,omitempty"`
	ReviewedRun string `json:"reviewedRun,omitempty"`
}
```
Add `ReviewsFile = "reviews.json"` in the file-name const block (alongside `OpenPRsFile`).

In `ledger_test.go` (or a state test), add a round-trip:
```go
func TestPRReviewRoundTrip(t *testing.T) {
	in := []state.PRReview{{Repo: "o/r", PR: 5, URL: "u", HeadSHA: "abc", Verdict: "good", Reasoning: "clean bump", ReviewedRun: "2026-06-24"}}
	dir := t.TempDir()
	require.NoError(t, state.Save(dir, state.ReviewsFile, in))
	var out []state.PRReview
	require.NoError(t, state.Load(dir, state.ReviewsFile, &out))
	assert.Equal(t, in, out)
}
```

- [ ] **Step 2: Run + commit**

Run: `go test ./internal/state/... && go build ./...`
```bash
git add internal/state/types.go internal/state/ledger_test.go
git commit -m "feat(state): PRReview + reviews.json"
```

---

### Task 4: review package — Assessor + OpenAIAssessor + Run

**Files:** Create `internal/review/assessor.go`, `internal/review/openai.go`, `internal/review/openai_test.go`, `internal/review/run.go`, `internal/review/run_test.go`.

**Interfaces:**
- `type Assessor interface { Assess(diff []byte, pr ghclient.PullRequest) (verdict, reasoning string, err error) }`; `FakeAssessor{Verdict, Reasoning string; Err error}`.
- `func NewOpenAIAssessor(cfg config.AIConfig) *OpenAIAssessor` — forced `assess_pr` tool call; on any error/unparseable → returns `("needs_human_verification", <reason>, nil)`.
- `func Run(repos []state.Repo, gh ghclient.GitHub, a Assessor, cfg config.ReviewCfg, prev []state.PRReview, runID string, dryRun bool) ([]state.PRReview, []state.CollectionError)`.

- [ ] **Step 1: Write the failing tests**

`internal/review/run_test.go` (FakeAssessor + a fake GitHub):

```go
package review

import (
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeGH struct {
	ghclient.GitHub
	prs       map[string][]ghclient.PullRequest
	diffs     map[string][]byte
	comments  []string // "repo#pr: body"
	approved  []string // "repo#pr"
}

func (f *fakeGH) ListOpenPRs(repo string) ([]ghclient.PullRequest, error) { return f.prs[repo], nil }
func (f *fakeGH) PRDiff(repo string, pr int) ([]byte, error) {
	return f.diffs[repo], nil
}
func (f *fakeGH) PostPRComment(repo string, pr int, body string) error {
	f.comments = append(f.comments, repo+"#"+itoa(pr)+": "+body)
	return nil
}
func (f *fakeGH) ApprovePR(repo string, pr int, body string) error {
	f.approved = append(f.approved, repo+"#"+itoa(pr))
	return nil
}
func itoa(n int) string { return string(rune('0'+n)) } // tiny helper for single-digit test PRs

func TestRunAssessesBotPRsIdempotently(t *testing.T) {
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 2, Title: "bump x", Author: "app/dependabot", IsBot: true, HeadSHA: "sha2", URL: "u2"},
			{Number: 3, Title: "human pr", Author: "alice", HeadSHA: "sha3"}, // not a bot -> skipped
		}},
		diffs: map[string][]byte{"o/r": []byte("go.mod bump")},
	}
	a := &FakeAssessor{Verdict: "good", Reasoning: "clean"}
	cfg := config.ReviewCfg{Enabled: true, AutoApprove: true, MaxPerRun: 20, Notify: []string{"@team"}}

	// First run: assesses the bot PR, comments (with cc), approves (good+autoApprove).
	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, out, 1)
	assert.Equal(t, "good", out[0].Verdict)
	assert.Equal(t, "sha2", out[0].HeadSHA)
	require.Len(t, gh.comments, 1)
	assert.Contains(t, gh.comments[0], "good")
	assert.Contains(t, gh.comments[0], "@team")
	require.Len(t, gh.approved, 1)

	// Second run with the SAME head SHA in prev: carried forward, no new comment/approve.
	a2 := &FakeAssessor{Err: assertFail} // would error if called
	out2, _ := Run([]state.Repo{{Repo: "o/r"}}, gh, a2, cfg, out, "run2", false)
	require.Len(t, out2, 1)
	assert.Equal(t, "good", out2[0].Verdict)
	assert.Len(t, gh.comments, 1) // unchanged — idempotent
}

var assertFail = &fakeErr{}

type fakeErr struct{}

func (*fakeErr) Error() string { return "assessor must not be called for unchanged head" }
```

(If the controller's brief specifies a cleaner fake-GH approach reusing an existing test fake, follow that; otherwise this self-contained fake is fine. The `itoa` shim only needs single-digit PR numbers in the test.)

Add `internal/review/openai_test.go`: an httptest server returning a forced tool call with `{"verdict":"bad","reasoning":"…"}` → `Assess` returns `("bad", …, nil)`; an empty-endpoint or error server → `("needs_human_verification", …, nil)`.

- [ ] **Step 2: Run red.** `go test ./internal/review/...`

- [ ] **Step 3: Implement**

`assessor.go`:
```go
package review

import "github.com/kairos-io/security/internal/ghclient"

type Assessor interface {
	Assess(diff []byte, pr ghclient.PullRequest) (verdict, reasoning string, err error)
}

type FakeAssessor struct {
	Verdict, Reasoning string
	Err                error
}

func (f *FakeAssessor) Assess([]byte, ghclient.PullRequest) (string, string, error) {
	return f.Verdict, f.Reasoning, f.Err
}
```

`openai.go`: mirror `internal/triage/openai.go`'s chat/tool-call plumbing (chatRequest/toolDef/toolFunctionDef/chatMessage/chatResponse/ToolChoice + the `/v1/chat/completions` POST). Tool name `assess_pr`, parameters schema constraining `verdict` to the enum + `reasoning`. Build the prompt from the PR title + a truncated diff (cap ~60000 bytes; append a "[diff truncated]" note when cut). On empty endpoint, non-200, decode error, no tool call, or a verdict outside the enum → return `("needs_human_verification", <reason>, nil)` (never a hard error). `NewOpenAIAssessor(cfg config.AIConfig)` uses `cfg.Nib.Endpoint`/`Model`/`Temperature` like the triage client.

`run.go`:
```go
package review

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

const reviewMarker = "<!-- ksec:review -->"

func Run(repos []state.Repo, gh ghclient.GitHub, a Assessor, cfg config.ReviewCfg, prev []state.PRReview, runID string, dryRun bool) ([]state.PRReview, []state.CollectionError) {
	prior := map[string]state.PRReview{}
	for _, r := range prev {
		prior[key(r.Repo, r.PR)] = r
	}
	var out []state.PRReview
	var errs []state.CollectionError
	assessed := 0
	for _, repo := range repos {
		prs, err := gh.ListOpenPRs(repo.Repo)
		if err != nil {
			errs = append(errs, state.CollectionError{Repo: repo.Repo, Collector: "review", Message: err.Error()})
			continue
		}
		for _, pr := range prs {
			if !pr.IsBot {
				continue
			}
			k := key(repo.Repo, pr.Number)
			// Idempotent: unchanged head -> carry the prior review forward.
			if p, ok := prior[k]; ok && p.HeadSHA == pr.HeadSHA {
				out = append(out, p)
				continue
			}
			if assessed >= cfg.MaxPerRun {
				// Over budget this run: keep any prior review so the dashboard
				// still shows it; otherwise skip until a future run.
				if p, ok := prior[k]; ok {
					out = append(out, p)
				}
				continue
			}
			assessed++
			diff, derr := gh.PRDiff(repo.Repo, pr.Number)
			if derr != nil {
				errs = append(errs, state.CollectionError{Repo: repo.Repo, Collector: "review", Message: derr.Error()})
				continue
			}
			verdict, reasoning, _ := a.Assess(diff, pr) // assessor never hard-errors (defaults needs_human)
			rv := state.PRReview{Repo: repo.Repo, PR: pr.Number, URL: pr.URL, HeadSHA: pr.HeadSHA,
				Verdict: verdict, Reasoning: reasoning, ReviewedRun: runID}
			out = append(out, rv)
			body := comment(rv, cfg.Notify)
			if dryRun {
				fmt.Printf("[dry-run] would comment on %s#%d: %s\n", repo.Repo, pr.Number, verdict)
				continue
			}
			_ = gh.PostPRComment(repo.Repo, pr.Number, body)
			if cfg.AutoApprove && verdict == "good" {
				_ = gh.ApprovePR(repo.Repo, pr.Number, "kairos-security: automated review verdict good")
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Repo != out[j].Repo {
			return out[i].Repo < out[j].Repo
		}
		return out[i].PR < out[j].PR
	})
	return out, errs
}

func key(repo string, pr int) string { return fmt.Sprintf("%s#%d", repo, pr) }

func comment(r state.PRReview, notify []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "🔎 kairos-security review: **%s** — %s", r.Verdict, r.Reasoning)
	if len(notify) > 0 {
		fmt.Fprintf(&b, "\n\ncc %s", strings.Join(notify, " "))
	}
	b.WriteString("\n\n" + reviewMarker)
	return b.String()
}
```

- [ ] **Step 4: Run green + build + commit**

Run: `go test ./internal/review/... && go build ./...`
```bash
git add internal/review/
git commit -m "feat(review): bot-PR assessor + idempotent Run (verdict/comment/approve)"
```

---

### Task 5: `ksec review` command + workflow step

**Files:** Modify `cmd/ksec/main.go`, `.github/workflows/security-dashboard.yaml`.

**Interfaces:** `newReviewCmd(gf)` registered on root; workflow runs `ksec review` after `remediate`, before `render`.

- [ ] **Step 1: Add the command** — `newReviewCmd`:
```go
func newReviewCmd(gf *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "review",
		Short: "AI-assess open bot PRs and post a verdict",
		RunE: func(cmd *cobra.Command, args []string) error {
			aiCfg, err := config.LoadAI("ai.yaml")
			if err != nil {
				return err
			}
			if !aiCfg.Review.Enabled || aiCfg.Nib.Endpoint == "" {
				fmt.Fprintln(os.Stderr, "review: disabled or no AI endpoint — skipping")
				return nil
			}
			var repos []state.Repo
			if err := state.Load(gf.stateDir, state.ReposFile, &repos); err != nil {
				return err
			}
			var prev []state.PRReview
			_ = state.Load(gf.stateDir, state.ReviewsFile, &prev) // best-effort
			gh := ghclient.NewCLI()
			reviews, errs := review.Run(repos, gh, review.NewOpenAIAssessor(aiCfg), aiCfg.Review, prev, collect.Today(), gf.dryRun)
			counts := map[string]int{}
			for _, r := range reviews {
				counts[r.Verdict]++
			}
			fmt.Fprintf(os.Stderr, "review: %d reviews (good=%d bad=%d needs-human=%d) · %d errors\n",
				len(reviews), counts["good"], counts["bad"], counts["needs_human_verification"], len(errs))
			return state.Save(gf.stateDir, state.ReviewsFile, reviews)
		},
	}
}
```
Register it: `root.AddCommand(newReviewCmd(gf))`. (Confirm imports: `review`, `collect` for `Today()`, `config`, `ghclient`, `state`, `os`, `fmt` — add `internal/review` import.) If `gf.dryRun` is the global dry-run flag, use it; otherwise add a local `--dry-run` consistent with other write commands.

- [ ] **Step 2: Workflow step** — in the "Run pipeline" step, after the `ksec remediate …` line and before `ksec render …`:
```yaml
          ksec review    --state-dir state $REMEDIATE_DRYRUN
```

- [ ] **Step 3: Build + vet + gofmt + test + YAML + smoke + commit**

Run: `go build ./... && go vet ./... && test -z "$(gofmt -l .)" && go test ./...`; `python3 -c "import yaml; yaml.safe_load(open('.github/workflows/security-dashboard.yaml'))" && echo OK`; `go run ./cmd/ksec review --help`.
```bash
git add cmd/ksec/main.go .github/workflows/security-dashboard.yaml
git commit -m "feat(cmd): ksec review command + pipeline step"
```

---

### Task 6: Dashboard "🔎 Bot-PR reviews" section

**Files:** Modify `internal/render/render.go`, `internal/render/render_test.go`, `internal/render/html.go`, `internal/render/testdata/`; wire `render` command to load `reviews.json`.

**Interfaces:** `render.Input` gains `Reviews []state.PRReview`; a "🔎 Bot-PR reviews" section grouped by repo with a verdict icon.

- [ ] **Step 1: Write the failing test** — `render_test.go`:
```go
func TestDashboardShowsBotPRReviews(t *testing.T) {
	md := DashboardMarkdown(Input{Reviews: []state.PRReview{
		{Repo: "kairos-io/AuroraBoot", PR: 566, URL: "https://github.com/kairos-io/AuroraBoot/pull/566", Verdict: "good", Reasoning: "clean go.mod bump"},
		{Repo: "kairos-io/AuroraBoot", PR: 567, URL: "u567", Verdict: "needs_human_verification", Reasoning: "touches source"},
	}})
	assert.Contains(t, md, "🔎 Bot-PR reviews")
	assert.Contains(t, md, "[#566")
	assert.Contains(t, md, "good")
	assert.Contains(t, md, "needs_human_verification")
}
```

- [ ] **Step 2: Run red.** `go test ./internal/render/...`

- [ ] **Step 3: Implement** — add `Reviews []state.PRReview \`json:"reviews,omitempty"\`` to `Input`. In `DashboardMarkdown`, when `len(in.Reviews) > 0`, render a section (place it after the Open PRs / bot ledger area):
```go
	if len(in.Reviews) > 0 {
		b.WriteString("## 🔎 Bot-PR reviews\n\n")
		repo := ""
		for _, r := range in.Reviews {
			if r.Repo != repo {
				repo = r.Repo
				fmt.Fprintf(&b, "**%s**\n\n", repo)
			}
			link := fmt.Sprintf("#%d", r.PR)
			if r.URL != "" {
				link = fmt.Sprintf("[#%d](%s)", r.PR, r.URL)
			}
			fmt.Fprintf(&b, "- %s — %s **%s** — %s\n", link, verdictIcon(r.Verdict), r.Verdict, r.Reasoning)
		}
		b.WriteString("\n")
	}
```
Add `verdictIcon` (good → ✅, bad → ⛔, needs_human_verification → ⚠️, else ""). Mirror in `html.go` (escaped). Regenerate goldens.

- [ ] **Step 4: Wire the render command** — in `newRenderCmd`, load reviews best-effort and set on the base Input:
```go
			var reviews []state.PRReview
			_ = state.Load(gf.stateDir, state.ReviewsFile, &reviews)
```
and add `Reviews: reviews` to the `render.Input{...}` literal.

- [ ] **Step 5: Regenerate goldens, build, vet, gofmt, test, commit**

Run: `UPDATE_GOLDEN=1 go test ./internal/render/...` (eyeball the reviews section + determinism + no raw ids), re-run; then `go build ./... && go vet ./... && test -z "$(gofmt -l .)" && go test ./...`.
```bash
git add internal/render/ cmd/ksec/main.go
git commit -m "feat(render): Bot-PR reviews section"
```

---

## Self-review

**Spec coverage:**
- review config (enabled/autoApprove/maxPerRun/notify) → Task 1. ✓
- diff fetch + head SHA + approve → Task 2. ✓
- PRReview state → Task 3. ✓
- forced-tool-call assessor (verdict enum; needs_human default) + idempotent Run (comment+cc, autoApprove good only, dry-run no writes, cap, errors→CollectionError, deterministic) → Task 4. ✓
- `ksec review` command + pipeline step (gated on enabled+endpoint; dry-run) → Task 5. ✓
- dashboard reviews section → Task 6. ✓
- bots only; auto-merge out of scope → Tasks 4 (IsBot filter; no merge). ✓

**Placeholder scan:** none — full code for config, the Run orchestration + comment builder, the command, and the render section; Task 4's `openai.go` and Task 2's CLI methods are described against the exact existing patterns (triage forced tool call; `c.run` shell-outs).

**Type consistency:** `ReviewCfg`/`AIConfig.Review` (Task 1) consumed by `review.Run`/the command (Tasks 4-5). `PullRequest.HeadSHA`/`PRDiff`/`ApprovePR` (Task 2) used by `review.Run` + fakes (Task 4). `state.PRReview`/`ReviewsFile` (Task 3) produced by Run, persisted by the command, rendered (Task 6). `Assessor`/`NewOpenAIAssessor` (Task 4) wired in the command (Task 5).

---

## Operational notes

- First live run comments on up to `maxPerRun` existing bot PRs, then spreads the rest over subsequent runs (head-SHA idempotency means each PR is commented once per change).
- Dry-run (`$REMEDIATE_DRYRUN`) prints the verdict it would post without commenting/approving — the safe way to preview.
- Auto-approve is off unless `review.autoApprove: true`; it only ever approves a `good` verdict (never merges).
- The assessor degrades to `needs_human_verification` whenever LocalAI is unavailable or returns anything unexpected — never blocks the pipeline.
