# Bot-PR Review — Upstream Diffs from PR-Body Compare Links Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fetch upstream source diffs for any ecosystem (npm/Docker/Actions/Go) by parsing `github.com/<owner>/<repo>/compare/<base>...<head>` links out of the renovate/dependabot PR body and `CompareDiff`-ing them, deduped/capped, alongside the existing Go `go.mod` path. Fixes AuroraBoot#566 (npm lucide bump) and records everything in the trace.

**Architecture:** Pure `parseCompareURLs` + `compareTargets` (unify Go-bump compares + body-link compares, dedup, cap) in `internal/review/deps.go`; `review.Run` iterates `compareTargets` to build the upstream context + trace. Body links use base/head verbatim (renovate already resolved the tag form); the Go path keeps `compareRef`.

**Tech Stack:** Go 1.22, `stretchr/testify`, existing `internal/review`.

## Global Constraints

- Module `github.com/kairos-io/security`; Go 1.22.
- Body compare links: base/head **verbatim** (no `v`-prefix munging). Go go.mod bumps: keep `compareRef` (v-tag / pseudo-version SHA).
- Dedup compares by `repo|base|head`; cap count at `maxCompares` (5); per-diff `maxBumpDiff` (40000); total context `maxContext` (60000).
- Idempotency (HeadSHA), MaxPerRun cap, comment upsert (no spam), dry-run no writes, degrade-still-assess — all preserved. Only the upstream-context assembly + trace lines change.
- Deterministic ordering (Go bumps first in parse order, then body links in match order).

---

## File structure

```
internal/review/deps.go        # CompareRef, parseCompareURLs, compareTargets (modify)
internal/review/deps_test.go   # (modify)
internal/review/run.go         # use compareTargets in context assembly (modify)
internal/review/run_test.go    # (modify)
```

---

### Task 1: parseCompareURLs + compareTargets

**Files:** Modify `internal/review/deps.go`, `internal/review/deps_test.go`.

**Interfaces:** `type CompareRef struct{ Repo, Base, Head, Label string }`; `func parseCompareURLs(body string) []CompareRef`; `func compareTargets(diff []byte, body string) []CompareRef`.

- [ ] **Step 1: Write the failing tests** — add to `deps_test.go`:

```go
func TestParseCompareURLs(t *testing.T) {
	body := "Release notes\n" +
		"[Compare Source](https://redirect.github.com/lucide-icons/lucide/compare/0.576.0...0.577.0)\n" +
		"Full Changelog: <https://github.com/lucide-icons/lucide/compare/0.468.0...0.577.0>\n" +
		"dup: https://github.com/lucide-icons/lucide/compare/0.576.0...0.577.0\n"
	got := parseCompareURLs(body)
	// deduped: two distinct compares
	assert.Len(t, got, 2)
	assert.Equal(t, CompareRef{Repo: "lucide-icons/lucide", Base: "0.576.0", Head: "0.577.0",
		Label: "lucide-icons/lucide 0.576.0..0.577.0 (PR body)"}, got[0])
	assert.Equal(t, "0.468.0", got[1].Base)
	assert.Equal(t, "0.577.0", got[1].Head)
}

func TestParseCompareURLsNone(t *testing.T) {
	assert.Empty(t, parseCompareURLs("no links here"))
}

func TestCompareTargetsUnifiesAndCaps(t *testing.T) {
	// a Go bump (via go.mod diff) + a body compare link → both, deduped
	diff := []byte("--- a/go.mod\n+++ b/go.mod\n" +
		"-\tgithub.com/foo/bar v1.2.0\n+\tgithub.com/foo/bar v1.3.0\n")
	body := "https://github.com/baz/qux/compare/v2.0.0...v2.1.0"
	got := compareTargets(diff, body)
	require.Len(t, got, 2)
	assert.Equal(t, "foo/bar", got[0].Repo) // Go path first
	assert.Equal(t, "v1.2.0", got[0].Base)  // compareRef adds v
	assert.Equal(t, "v1.3.0", got[0].Head)
	assert.Equal(t, "baz/qux", got[1].Repo) // body link
	assert.Equal(t, "v2.0.0", got[1].Base)  // verbatim from URL
}
```

- [ ] **Step 2: Run red.** `go test ./internal/review/...`

- [ ] **Step 3: Implement** — in `deps.go`:

```go
type CompareRef struct{ Repo, Base, Head, Label string }

const maxCompares = 5

var reCompare = regexp.MustCompile(`github\.com/([\w.\-]+)/([\w.\-]+)/compare/([\w.\-+/@]+?)\.\.\.([\w.\-+/@]+)`)

// parseCompareURLs extracts GitHub compare links (owner/repo + base...head) from
// a renovate/dependabot PR body. Base/head are taken verbatim — the bot already
// resolved the repo's real tag form (incl. monorepos). Works for any ecosystem.
func parseCompareURLs(body string) []CompareRef {
	var out []CompareRef
	seen := map[string]bool{}
	for _, m := range reCompare.FindAllStringSubmatch(body, -1) {
		repo := m[1] + "/" + m[2]
		ref := CompareRef{Repo: repo, Base: m[3], Head: m[4],
			Label: fmt.Sprintf("%s %s..%s (PR body)", repo, m[3], m[4])}
		k := ref.Repo + "|" + ref.Base + "|" + ref.Head
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, ref)
	}
	return out
}

// compareTargets unifies upstream comparisons from the Go go.mod bumps and the
// PR-body compare links, deduped by repo|base|head and capped at maxCompares.
func compareTargets(diff []byte, body string) []CompareRef {
	var out []CompareRef
	seen := map[string]bool{}
	add := func(r CompareRef) {
		if r.Repo == "" {
			return
		}
		k := r.Repo + "|" + r.Base + "|" + r.Head
		if seen[k] {
			return
		}
		seen[k] = true
		out = append(out, r)
	}
	for _, b := range parseBumps(diff) {
		if repo, ok := moduleRepo(b.Module); ok {
			add(CompareRef{Repo: repo, Base: compareRef(b.From), Head: compareRef(b.To),
				Label: fmt.Sprintf("%s %s→%s", b.Module, b.From, b.To)})
		}
	}
	for _, c := range parseCompareURLs(body) {
		add(c)
	}
	if len(out) > maxCompares {
		out = out[:maxCompares]
	}
	return out
}
```

(Add `"fmt"` to `deps.go` imports if not present.)

- [ ] **Step 4: Run green + commit**

Run: `go test ./internal/review/... && go build ./...`
```bash
git add internal/review/deps.go internal/review/deps_test.go
git commit -m "feat(review): parse PR-body compare links; unify upstream compare targets"
```

---

### Task 2: Use compareTargets in Run

**Files:** Modify `internal/review/run.go`, `internal/review/run_test.go`.

**Interfaces:** `Run`'s upstream-context assembly uses `compareTargets(diff, pr.Body)`.

- [ ] **Step 1: Replace the bump loop** — in `Run`, swap the `parseBumps`-only loop for `compareTargets`:

```go
		var trace []string
		targets := compareTargets(diff, pr.Body)
		if len(targets) == 0 {
			trace = append(trace, "no upstream comparisons available (no go.mod bumps or compare links in the PR body)")
		}
		for _, t := range targets {
			if ctx.Len() > maxContext {
				break
			}
			ud, uerr := gh.CompareDiff(t.Repo, t.Base, t.Head)
			if uerr != nil || len(ud) == 0 {
				trace = append(trace, fmt.Sprintf("%s: compare %s...%s failed/empty (no upstream diff)", t.Label, t.Base, t.Head))
				continue
			}
			if len(ud) > maxBumpDiff {
				ud = ud[:maxBumpDiff]
			}
			fmt.Fprintf(&ctx, "Upstream %s (%s...%s):\n%s\n\n", t.Label, t.Base, t.Head, ud)
			trace = append(trace, fmt.Sprintf("%s: compare %s...%s ✓ %d bytes", t.Label, t.Base, t.Head, len(ud)))
		}
		ctx.WriteString("PR diff:\n" + string(diff))
		trace = append(trace, fmt.Sprintf("context: %d bytes", ctx.Len()))
```

(Keep the PR body prepend before this block, the `Assess` call, `rv` construction with `Trace: trace`, idempotency carry-forward, MaxPerRun cap, PRDiff fetch+error, dry-run print+continue, and UpsertPRComment/ApprovePR exactly as they are. `parseBumps`/`moduleRepo`/`compareRef` are now only referenced via `compareTargets`.)

- [ ] **Step 2: Update run_test.go** — add a case: a bot PR with NO go.mod bump but a `pr.Body` containing `https://github.com/lucide-icons/lucide/compare/0.576.0...0.577.0`, and the fake `CompareDiff` keyed for `lucide-icons/lucide 0.576.0 0.577.0` returning a diff → assert the fake `CompareDiff` was called with those verbatim refs, the context contains `Upstream lucide-icons/lucide …`, and the trace has the `✓` line. Keep/adjust the existing pseudo-version Go test (still works via compareTargets). Ensure the fake's `CompareDiff` records (repo, base, head).

- [ ] **Step 3: Run green + build + vet + commit**

Run: `go test ./internal/review/... && go build ./... && go vet ./...`
```bash
git add internal/review/run.go internal/review/run_test.go
git commit -m "feat(review): fetch upstream diffs from PR-body compare links (npm/etc.)"
```

---

## Self-review

**Spec coverage:**
- Parse body compare links (any ecosystem) + fetch → Tasks 1 (`parseCompareURLs`), 2 (use in Run). ✓
- Unify with Go go.mod path, deduped + capped → Task 1 (`compareTargets`). ✓
- Body links verbatim base/head; Go path keeps `compareRef` → Task 1. ✓
- Trace reflects each compare (or "no upstream comparisons available") → Task 2. ✓
- Idempotency / cap / upsert / dry-run / degrade preserved → Task 2 (only the loop changes). ✓

**Placeholder scan:** none — full code for the helpers and the Run loop.

**Type consistency:** `CompareRef`/`parseCompareURLs`/`compareTargets` (Task 1) consumed by `Run` (Task 2); `compareTargets` internally uses the existing `parseBumps`/`moduleRepo`/`compareRef`. The render trace (prior feature) shows the new lines unchanged.

---

## Operational notes

- AuroraBoot#566 (lucide npm bump) now fetches `lucide-icons/lucide` compares straight from the renovate body — npm, Docker, and Actions bumps get upstream diffs the same way, with no per-ecosystem code.
- Only `github.com` compares are fetched; other hosts are ignored (the body changelog still informs the verdict).
- A monorepo's huge range compare is byte-capped; the per-version compare links in the body (also parsed) give tighter diffs, and the cap of 5 keeps the prompt bounded.
