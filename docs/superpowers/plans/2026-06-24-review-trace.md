# Bot-PR Review — Pseudo-Version Compare + Traceability Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the bot-PR review fetch the upstream diff for **pseudo-version** bumps (use the embedded commit SHA as the compare ref) and record a **trace** of what the assessor did (per-bump resolution/compare result + context size), surfaced in the PR comment (collapsible) and the dashboard.

**Architecture:** A pure `compareRef` maps a version to its compare ref (SHA for pseudo-versions, `v`-tag otherwise); `review.Run` uses it and builds a `Trace []string` during context assembly; the comment gets a `<details>review trace</details>` block; the dashboard shows trace sub-lines. Builds on the just-shipped review dep-context feature.

**Tech Stack:** Go 1.22, `stretchr/testify`, existing `internal/review`, `internal/state`, `internal/render`.

## Global Constraints

- Module `github.com/kairos-io/security`; Go 1.22.
- Comment stays an **upsert** (one per PR, edited; `<details>` block included). Idempotency on head SHA preserved; degrade paths unchanged (still assess, trace records the skip).
- Dry-run zero writes (trace still computed + printed). Deterministic dashboard; trace is ordered.
- The commit SHAs that appear in the trace are intended provenance, not the SHA-256 finding id (the no-raw-id rule is about finding ids).

---

## File structure

```
internal/state/types.go         # PRReview.Trace []string (modify)
internal/review/deps.go         # compareRef (modify)
internal/review/deps_test.go    # compareRef tests (modify)
internal/review/run.go          # use compareRef; build Trace; comment <details> (modify)
internal/review/run_test.go     # (modify)
internal/render/render.go       # render Trace under review rows (modify)
internal/render/render_test.go  # (modify)
internal/render/html.go         # mirror (modify)
internal/render/testdata/       # regenerated goldens
```

---

### Task 1: PRReview.Trace + compareRef

**Files:** Modify `internal/state/types.go`, `internal/review/deps.go`, `internal/review/deps_test.go`.

**Interfaces:** `state.PRReview.Trace []string` (json `trace,omitempty`); `func compareRef(version string) string`.

- [ ] **Step 1: Add the state field** — in `state.PRReview`, add `Trace []string \`json:"trace,omitempty"\``.

- [ ] **Step 2: Write the failing test** — add to `internal/review/deps_test.go`:

```go
func TestCompareRef(t *testing.T) {
	assert.Equal(t, "v0.33.0", compareRef("0.33.0"))
	assert.Equal(t, "fab4fdf2f2f3", compareRef("0.0.0-20241017190036-fab4fdf2f2f3"))
	assert.Equal(t, "abcdef123456", compareRef("1.2.3-0.20240101000000-abcdef123456"))
	assert.Equal(t, "v2.0.0+incompatible", compareRef("2.0.0+incompatible")) // not a pseudo-version
}
```

- [ ] **Step 3: Run red.** `go test ./internal/review/...`

- [ ] **Step 4: Implement** — in `deps.go`:

```go
// rePseudo matches a Go pseudo-version's trailing "-<14-digit-timestamp>-<sha>".
var rePseudo = regexp.MustCompile(`-\d{14}-([0-9a-f]{12,})$`)

// compareRef maps a (v-stripped) module version to the ref to compare against
// upstream: a pseudo-version compares by its embedded commit SHA, a real
// release compares by its "v"-prefixed tag.
func compareRef(version string) string {
	if m := rePseudo.FindStringSubmatch(version); m != nil {
		return m[1]
	}
	return "v" + version
}
```

- [ ] **Step 5: Run green + commit**

Run: `go test ./internal/review/... ./internal/state/... && go build ./...`
```bash
git add internal/state/types.go internal/review/deps.go internal/review/deps_test.go
git commit -m "feat(review): PRReview.Trace + compareRef (pseudo-version -> commit SHA)"
```

---

### Task 2: Use compareRef + build trace + comment details

**Files:** Modify `internal/review/run.go`, `internal/review/run_test.go`.

**Interfaces:** `Run` uses `compareRef` for `CompareDiff`, builds `rv.Trace`, and `comment()` appends a `<details>review trace</details>` block.

- [ ] **Step 1: Update the context-assembly + trace in `Run`** — in the per-PR assembly, build a `trace []string` alongside the context. Replace the bump loop with one that records each outcome and uses `compareRef`:

```go
		var trace []string
		bumps := parseBumps(diff)
		if len(bumps) == 0 {
			trace = append(trace, "no go.mod dependency bumps parsed from the PR diff")
		}
		for _, b := range bumps {
			if ctx.Len() > maxContext {
				break
			}
			gr, ok := moduleRepo(b.Module)
			if !ok {
				trace = append(trace, fmt.Sprintf("%s %s→%s: module not resolvable to a GitHub repo (skipped)", b.Module, b.From, b.To))
				continue
			}
			baseRef, headRef := compareRef(b.From), compareRef(b.To)
			ud, uerr := gh.CompareDiff(gr, baseRef, headRef)
			if uerr != nil || len(ud) == 0 {
				trace = append(trace, fmt.Sprintf("%s %s→%s: compare %s...%s failed: %v (no upstream diff)", b.Module, b.From, b.To, baseRef, headRef, uerr))
				continue
			}
			if len(ud) > maxBumpDiff {
				ud = ud[:maxBumpDiff]
			}
			fmt.Fprintf(&ctx, "Upstream %s %s..%s:\n%s\n\n", b.Module, b.From, b.To, ud)
			trace = append(trace, fmt.Sprintf("%s %s→%s: compare %s...%s ✓ %d bytes", b.Module, b.From, b.To, baseRef, headRef, len(ud)))
		}
		ctx.WriteString("PR diff:\n" + string(diff))
		trace = append(trace, fmt.Sprintf("context: %d bytes", ctx.Len()))
		verdict, reasoning, summary, _ := a.Assess(pr, ctx.String())
		rv := state.PRReview{Repo: repo.Repo, PR: pr.Number, URL: pr.URL, HeadSHA: pr.HeadSHA,
			Verdict: verdict, Reasoning: reasoning, ChangesSummary: summary, Trace: trace, ReviewedRun: runID}
```

(Keep the surrounding idempotency carry-forward, MaxPerRun cap, the `PRDiff` fetch + its error handling, the `out = append(out, rv)`, the dry-run print, and the upsert/approve exactly as they are — only the bump loop + `rv` construction change to add the trace and `compareRef`.)

Update the dry-run print to include the trace count if helpful: keep `fmt.Printf("[dry-run] would comment on %s#%d: %s — %s\n", ...)` as-is.

- [ ] **Step 2: Add the `<details>` block to `comment()`** — after the dependency-changes line (and cc), before the marker:

```go
	if len(r.Trace) > 0 {
		b.WriteString("\n\n<details><summary>review trace</summary>\n\n```\n")
		for _, line := range r.Trace {
			b.WriteString(line + "\n")
		}
		b.WriteString("```\n\n</details>")
	}
	b.WriteString("\n\n" + reviewMarker)
```

(Adjust to the existing `comment()` structure — ensure the marker stays last.)

- [ ] **Step 3: Update `run_test.go`** — assert: a pseudo-version bump (e.g. `0.0.0-20241017190036-fab4fdf2f2f3` → `…-d29549a44f29`) makes the fake `CompareDiff` receive the SHA refs (`fab4fdf2f2f3`, `d29549a44f29`) and the trace records `✓ <n> bytes`; an unresolvable module records the "not resolvable" trace line; no-bumps records the "no bumps" line; `rv.Trace` is non-empty; the comment contains `<details><summary>review trace</summary>`. (Have the fakeGH `CompareDiff` record the (repo, base, head) it was called with so the SHA-ref assertion is possible.)

- [ ] **Step 4: Run green + build + vet + commit**

Run: `go test ./internal/review/... && go build ./... && go vet ./...`
```bash
git add internal/review/run.go internal/review/run_test.go
git commit -m "feat(review): pseudo-version compare via commit SHA; record + comment the trace"
```

---

### Task 3: Dashboard shows the trace

**Files:** Modify `internal/render/render.go`, `internal/render/render_test.go`, `internal/render/html.go`, `internal/render/testdata/`.

- [ ] **Step 1: Write the failing test** — add to `render_test.go` (extend the reviews test or add one): a review with `Trace: []string{"foo/bar 1.0→1.1: compare v1.0.0...v1.1.0 ✓ 1234 bytes"}` renders that trace line in the "🔎 Bot-PR reviews" section.

- [ ] **Step 2: Implement** — in `render.go`'s reviews section, after the `ChangesSummary` sub-line, when `len(r.Trace) > 0` render each trace line as a deeper indented sub-item, e.g.:

```go
		for _, tl := range r.Trace {
			fmt.Fprintf(&b, "    - %s\n", tl)
		}
```

Mirror in `html.go` (a nested `<ul>`/lines, escaped via the template). Keep the section omitted when no reviews.

- [ ] **Step 3: Regenerate goldens, build, vet, gofmt, test, commit**

Run: `UPDATE_GOLDEN=1 go test ./internal/render/...` (eyeball: trace lines under a review; deterministic; section omitted when no reviews), re-run; then `go build ./... && go vet ./... && test -z "$(gofmt -l .)" && go test ./...`.
```bash
git add internal/render/ 
git commit -m "feat(render): show review trace under Bot-PR reviews"
```

---

## Self-review

**Spec coverage:**
- Pseudo-version → commit-SHA compare ref → Tasks 1 (`compareRef`), 2 (use it). ✓
- `Trace` recorded (per-bump resolve/compare result + context size + no-bumps) → Tasks 1 (field), 2 (build). ✓
- Trace in the PR comment `<details>` + dashboard → Tasks 2, 3. ✓
- Upsert/idempotency/degrade/dry-run preserved → Task 2 (only the bump loop + rv + comment change). ✓

**Placeholder scan:** none — full code for `compareRef`, the trace-building bump loop, the comment block, and the render lines.

**Type consistency:** `compareRef` (Task 1) used in `Run` (Task 2). `PRReview.Trace` (Task 1) set in `Run` (Task 2), rendered (Task 3). `comment()` change uses `r.Trace`.

---

## Operational notes

- `go-ukify#31` (and all renovate digest / Go pseudo-version bumps) now resolve the compare to the commit SHAs, so the assessor sees the real upstream change; the trace line proves it (`compare fab4fdf2f2f3...d29549a44f29 ✓ N bytes`).
- When upstream still can't be fetched (vanity module, non-semver tag, monorepo subdir tag), the trace says exactly why instead of silently degrading.
- The trace lives in `reviews.json`, the PR comment's collapsible block, and the dashboard — three places to see what the agent did.
