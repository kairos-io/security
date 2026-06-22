# Fork-Based PRs Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Let the bot remediate **external** repos it can't push to: fork the repo, push the `ksec/` branch to the bot's fork, and open a cross-repo PR (`--head <forkOwner>:<branch>`). Org repos push direct, unchanged.

**Architecture:** One fork-aware abstraction on `GitExecutor` — `ForkOwner`/`ShouldFork` fields plus `forkSlug`/`prHead`/`forkURL`/`ensureFork`/`pushBranch`/`checkoutOwnBranch` helpers — routed through the 7 push sites and 4 `gh pr create --head` sites. PRs are always opened/viewed against upstream, so `Reconcile`/`MergePR`/`ClosePR`/`PostPRComment` (keyed by the upstream PR number) are untouched. A nil `ShouldFork` never forks, so current behavior and all existing tests stay byte-identical.

**Tech Stack:** Go 1.22, `stretchr/testify`, `gh`/`git` CLIs, existing `internal/remediate`, `internal/state`, `cmd/ksec`.

## Global Constraints

- Module `github.com/kairos-io/security`; binary `ksec`; Go 1.22.
- **external (`kind == "external"`) → fork; org → push-direct.** The decision is data-driven from `repos.json`'s `kind`. Unknown repos → push-direct (no surprise forks).
- **Nil `ShouldFork` ⇒ never fork** — every existing executor test (which constructs `GitExecutor` without the new fields) must behave exactly as today.
- Never push/force-push a foreign upstream branch. For external repos we only push our own `ksec/` branch to the **bot's fork**; the PR is opened against upstream.
- Dry-run performs zero writes (fork/push/PR all no-op + print). Token redacted via `g.run`. Verify-before-push unchanged.
- The fork remote uses the same `x-access-token:<GH_TOKEN>` auth as `cloneURL`; the bot is a machine-user PAT that owns its forks.

---

## File structure

```
internal/remediate/git_executor.go   # fields + helpers; route 7 push + 4 head sites (modify)
internal/remediate/git_executor_test.go  # pure-helper tests (create or modify)
internal/remediate/fork.go           # ForkByKind decision builder (create)
internal/remediate/fork_test.go      # (create)
cmd/ksec/main.go                     # resolve ForkOwner + wire ShouldFork (modify)
```

---

### Task 1: Fork fields + pure helpers

**Files:** Modify `internal/remediate/git_executor.go`; create/modify `internal/remediate/git_executor_test.go`.

**Interfaces:**
- `GitExecutor` gains `ForkOwner string` and `ShouldFork func(repo string) bool`.
- `func forkSlug(forkOwner, repo string) string` (pure) → `forkOwner + "/" + path.Base(repo)`.
- `func (g *GitExecutor) forking(repo string) bool` → `g.ShouldFork != nil && g.ShouldFork(repo)`.
- `func (g *GitExecutor) prHead(repo, branch string) string` → `branch`, or `g.ForkOwner + ":" + branch` when forking.
- `func (g *GitExecutor) forkURL(repo string) string` → token URL to `forkSlug(g.ForkOwner, repo)`.

- [ ] **Step 1: Write the failing tests**

Create `internal/remediate/git_executor_test.go` (or add if it exists):

```go
package remediate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForkSlug(t *testing.T) {
	assert.Equal(t, "kairos-security-bot/edgevpn", forkSlug("kairos-security-bot", "mudler/edgevpn"))
	assert.Equal(t, "bot/kairos", forkSlug("bot", "kairos-io/kairos"))
}

func TestPRHead(t *testing.T) {
	ext := func(string) bool { return true }
	org := func(string) bool { return false }
	g := &GitExecutor{ForkOwner: "kairos-security-bot", ShouldFork: ext}
	assert.Equal(t, "kairos-security-bot:ksec/x", g.prHead("mudler/edgevpn", "ksec/x"))

	g2 := &GitExecutor{ForkOwner: "kairos-security-bot", ShouldFork: org}
	assert.Equal(t, "ksec/x", g2.prHead("kairos-io/kairos", "ksec/x"))

	g3 := &GitExecutor{} // nil ShouldFork -> never fork
	assert.Equal(t, "ksec/x", g3.prHead("mudler/edgevpn", "ksec/x"))
}

func TestForkURL(t *testing.T) {
	g := &GitExecutor{ForkOwner: "bot", Token: "tok"}
	assert.Equal(t, "https://x-access-token:tok@github.com/bot/edgevpn.git", g.forkURL("mudler/edgevpn"))
}
```

- [ ] **Step 2: Run red.** `go test ./internal/remediate/...`

- [ ] **Step 3: Implement** — in `git_executor.go`, add the two struct fields (after `Agent`), import `path`, and add:

```go
func forkSlug(forkOwner, repo string) string { return forkOwner + "/" + path.Base(repo) }

func (g *GitExecutor) forking(repo string) bool { return g.ShouldFork != nil && g.ShouldFork(repo) }

func (g *GitExecutor) prHead(repo, branch string) string {
	if g.forking(repo) {
		return g.ForkOwner + ":" + branch
	}
	return branch
}

func (g *GitExecutor) forkURL(repo string) string {
	slug := forkSlug(g.ForkOwner, repo)
	if g.Token != "" {
		return "https://x-access-token:" + g.Token + "@github.com/" + slug + ".git"
	}
	return "https://github.com/" + slug + ".git"
}
```

- [ ] **Step 4: Run green + build + commit**

Run: `go test ./internal/remediate/... && go build ./...`
```bash
git add internal/remediate/git_executor.go internal/remediate/git_executor_test.go
git commit -m "feat(remediate): fork fields + prHead/forkURL/forkSlug helpers"
```

---

### Task 2: Fork shell helpers — ensureFork, pushBranch, checkoutOwnBranch

**Files:** Modify `internal/remediate/git_executor.go`, `internal/remediate/git_executor_test.go`.

**Interfaces:**
- `func (g *GitExecutor) ensureFork(repo string) error` — `gh repo fork <repo> --clone=false` (idempotent); dry-run no-op.
- `func (g *GitExecutor) pushBranch(dir, repo, branch string, force bool) error` — the single push chokepoint (origin for org; ensure-fork + `fork` remote + push for external); dry-run prints, no writes.
- `func (g *GitExecutor) checkoutOwnBranch(dir, repo, branch string) error` — org: `git checkout <branch>`; external: add `fork` remote, `git fetch fork`, `git checkout -b <branch> fork/<branch>`.

- [ ] **Step 1: Write the dry-run test** (the shelling itself is integration; the dry-run decision path is testable):

```go
func TestPushBranchDryRunNoWrites(t *testing.T) {
	// Dry-run must not shell out, for both fork and non-fork repos.
	g := &GitExecutor{DryRun: true, ForkOwner: "bot", ShouldFork: func(string) bool { return true }}
	assert.NoError(t, g.pushBranch("/tmp/nonexistent", "mudler/edgevpn", "ksec/x", false))
	g2 := &GitExecutor{DryRun: true} // org/no-fork
	assert.NoError(t, g2.pushBranch("/tmp/nonexistent", "kairos-io/kairos", "ksec/x", true))
	assert.NoError(t, g.ensureFork("mudler/edgevpn"))
}
```

- [ ] **Step 2: Run red.** `go test ./internal/remediate/...`

- [ ] **Step 3: Implement**

```go
func (g *GitExecutor) ensureFork(repo string) error {
	if g.DryRun {
		fmt.Printf("[dry-run] would ensure fork of %s under %s\n", repo, g.ForkOwner)
		return nil
	}
	// Idempotent: gh exits 0 and prints "already exists" when the fork is present.
	_, err := g.run("", "gh", "repo", "fork", repo, "--clone=false")
	return err
}

func (g *GitExecutor) pushBranch(dir, repo, branch string, force bool) error {
	args := []string{"push"}
	if force {
		args = append(args, "--force")
	}
	if !g.forking(repo) {
		if g.DryRun {
			fmt.Printf("[dry-run] would push %s to origin (%s)\n", branch, repo)
			return nil
		}
		_, err := g.run(dir, "git", append(args, "-u", "origin", branch)...)
		return err
	}
	if g.DryRun {
		fmt.Printf("[dry-run] would push %s to fork %s\n", branch, forkSlug(g.ForkOwner, repo))
		return nil
	}
	if err := g.ensureFork(repo); err != nil {
		return err
	}
	_, _ = g.run(dir, "git", "remote", "add", "fork", g.forkURL(repo)) // ignore "already exists"
	_, err := g.run(dir, "git", append(args, "fork", branch)...)
	return err
}

func (g *GitExecutor) checkoutOwnBranch(dir, repo, branch string) error {
	if !g.forking(repo) {
		_, err := g.run(dir, "git", "checkout", branch)
		return err
	}
	if g.DryRun {
		return nil
	}
	if err := g.ensureFork(repo); err != nil {
		return err
	}
	_, _ = g.run(dir, "git", "remote", "add", "fork", g.forkURL(repo))
	if _, err := g.run(dir, "git", "fetch", "fork"); err != nil {
		return err
	}
	_, err := g.run(dir, "git", "checkout", "-b", branch, "fork/"+branch)
	return err
}
```

(Confirm `fmt` is imported — it is.)

- [ ] **Step 4: Run green + build + vet + commit**

Run: `go test ./internal/remediate/... && go build ./... && go vet ./...`
```bash
git add internal/remediate/git_executor.go internal/remediate/git_executor_test.go
git commit -m "feat(remediate): ensureFork/pushBranch/checkoutOwnBranch fork helpers"
```

---

### Task 3: Route the new-PR paths through the fork helpers

**Files:** Modify `internal/remediate/git_executor.go` (`Open`, `Cascade`, `Toolchain`, `Supersede`).

**Interfaces:** consumes `pushBranch`/`prHead` from Tasks 1-2. No new exported surface.

- [ ] **Step 1: Apply the edits** — in EACH of `Open`, `Cascade`, `Toolchain`, `Supersede` (they clone upstream and cut a fresh branch — keep the clone as-is):
  - Replace `if _, err := g.run(dir, "git", "push", "-u", "origin", branch); err != nil {` with `if err := g.pushBranch(dir, in.Repo, branch, false); err != nil {` (use the method's repo variable — `in.Repo` in all four).
  - Replace the `gh pr create … "--head", branch, …` argument `branch` with `g.prHead(in.Repo, branch)`.

(Each site keeps its existing error handling and the surrounding `git config`/commit/`gh pr create -R in.Repo` calls. Only the push call and the `--head` value change.)

- [ ] **Step 2: Build + vet + test + commit** — integration (the four methods are exercised by dry-run + build; existing tests use `FakeExecutor`, unaffected).

Run: `go build ./... && go vet ./... && go test ./...`
Confirm no `git push … origin` or bare `--head", branch` remains in the four new-PR methods (grep).
```bash
git add internal/remediate/git_executor.go
git commit -m "feat(remediate): new-PR paths push to fork + cross-repo head for external repos"
```

---

### Task 4: Route the update paths through the fork helpers

**Files:** Modify `internal/remediate/git_executor.go` (`Adjust`, `Repin`, `ResolveConflict`).

**Interfaces:** consumes `checkoutOwnBranch`/`pushBranch`. The branch already exists on the fork (external) or origin (org).

- [ ] **Step 1: Apply the edits** — in each of `Adjust`, `Repin`, `ResolveConflict` (they full-clone upstream via `cloneURL`):
  - Replace the post-clone `git checkout <branch>` (`entry.Branch`/`e.Branch`) with `g.checkoutOwnBranch(dir, <repo>, <branch>)` (`entry.Repo`/`e.Repo`). For `ResolveConflict`, keep the subsequent `git fetch origin` + `git rebase origin/HEAD` (the clone's `origin` is upstream, so the rebase base is correct for both org and external).
  - Replace `if _, err := g.run(dir, "git", "push", "--force", "origin", <branch>); err != nil {` with `if err := g.pushBranch(dir, <repo>, <branch>, true); err != nil {`.

(Variable names per method: `Adjust` uses `entry.Repo`/`entry.Branch`; `Repin` and `ResolveConflict` use `e.Repo`/`e.Branch`. The existing `ksec/` guard in `ResolveConflict`/`Repin` and verify-before-push are unchanged.)

- [ ] **Step 2: Build + vet + test + commit**

Run: `go build ./... && go vet ./... && go test ./...`
Confirm no `git push --force … origin` and no bare `git checkout <branch>` of an own-branch remains in these three methods (grep) — all route through the helpers.
```bash
git add internal/remediate/git_executor.go
git commit -m "feat(remediate): update paths fetch/force-push via the fork for external repos"
```

---

### Task 5: Wire ForkOwner + ShouldFork in the command

**Files:** Create `internal/remediate/fork.go`, `internal/remediate/fork_test.go`; modify `cmd/ksec/main.go`.

**Interfaces:**
- `func ForkByKind(repos []state.Repo) func(repo string) bool` — `true` iff the repo's `Kind == "external"`.
- `main.go` sets `ex.ForkOwner` (from `gh api user --jq .login`, falling back to `KSEC_FORK_OWNER`) and `ex.ShouldFork = remediate.ForkByKind(repos)`.

- [ ] **Step 1: Write the failing test** — `internal/remediate/fork_test.go`:

```go
package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestForkByKind(t *testing.T) {
	f := ForkByKind([]state.Repo{
		{Repo: "mudler/edgevpn", Kind: "external"},
		{Repo: "kairos-io/kairos", Kind: "org"},
	})
	assert.True(t, f("mudler/edgevpn"))   // external -> fork
	assert.False(t, f("kairos-io/kairos")) // org -> direct
	assert.False(t, f("unknown/repo"))     // unknown -> direct (no surprise fork)
}
```

- [ ] **Step 2: Run red.** `go test ./internal/remediate/...`

- [ ] **Step 3: Implement** — `internal/remediate/fork.go`:

```go
package remediate

import "github.com/kairos-io/security/internal/state"

// ForkByKind returns a predicate that reports whether a repo should be
// contributed to via a fork (external repos the bot can't push to directly).
// Org and unknown repos push direct.
func ForkByKind(repos []state.Repo) func(string) bool {
	external := map[string]bool{}
	for _, r := range repos {
		if r.Kind == "external" {
			external[r.Repo] = true
		}
	}
	return func(repo string) bool { return external[repo] }
}
```

In `cmd/ksec/main.go` `newRemediateCmd`, after `repos` is loaded and `ex` is constructed (it already has `Token`/`GH`/`Automerge`/`Prose`/`Agent`), add:

```go
			ex.ShouldFork = remediate.ForkByKind(repos)
			forkOwner := os.Getenv("KSEC_FORK_OWNER")
			if out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output(); err == nil {
				if login := strings.TrimSpace(string(out)); login != "" {
					forkOwner = login
				}
			}
			ex.ForkOwner = forkOwner
```

(Confirm `os`, `os/exec`, `strings` are imported in `main.go` — they are, used elsewhere.)

- [ ] **Step 4: Run green + build + vet + full test + smoke**

Run: `go test ./internal/remediate/... && go build ./... && go vet ./... && test -z "$(gofmt -l .)" && go test ./...`
Smoke: `go run ./cmd/ksec remediate --help` still works.
```bash
git add internal/remediate/fork.go internal/remediate/fork_test.go cmd/ksec/main.go
git commit -m "feat(remediate): wire ForkOwner (gh api user) + ShouldFork by repo kind"
```

---

## Self-review

**Spec coverage:**
- external→fork / org→direct decision (data-driven from kind) → Tasks 1 (`forking`), 5 (`ForkByKind`). ✓
- `ForkOwner`/`ShouldFork` fields; nil ⇒ never fork (regression-safe) → Task 1. ✓
- `forkSlug`/`prHead`/`forkURL`/`ensureFork`/`pushBranch`/`checkoutOwnBranch` → Tasks 1-2. ✓
- New-PR paths push-to-fork + `--head forkOwner:branch` → Task 3. ✓
- Update paths (Adjust/Repin/ResolveConflict) fetch-from-fork + force-push-to-fork; rebase base stays upstream → Task 4. ✓
- PR/Reconcile/Merge/Close/Comment keyed on upstream PR number, unchanged → (untouched; verified in Tasks 3-4 by only changing push + head). ✓
- Wiring: ForkOwner via `gh api user` + `KSEC_FORK_OWNER` fallback; ShouldFork from repos.json → Task 5. ✓
- Dry-run no writes; token redacted → Tasks 2-4. ✓

**Placeholder scan:** none — complete code for the pure helpers (Tasks 1, 5) and the shell helpers (Task 2); Tasks 3-4 are mechanical search-replace at enumerated sites with exact before/after.

**Type consistency:** `forking`/`prHead`/`forkURL`/`forkSlug` (Task 1) used by `pushBranch`/`checkoutOwnBranch` (Task 2), used by all six executors (Tasks 3-4). `ForkByKind` (Task 5) sets `ShouldFork` (Task 1 field). All push/clone helpers use the existing `g.run` (token-redacting).

---

## Operational notes

- `ForkOwner` resolves to the PAT's own login via `gh api user`; the forks live at `<login>/<repo>`. Set `KSEC_FORK_OWNER` to override (or as a fallback if `gh api user` is unavailable).
- The first fork of a repo may provision asynchronously; a transient `git push fork` failure is a recorded per-repo executor error and retries next run.
- Org-repo behavior is byte-identical (push-direct) — `ShouldFork` returns false for them.
- A repo missing from `repos.json` pushes direct (no surprise fork); add it to `repos.yaml` with `kind: external` to enable forking.
```
