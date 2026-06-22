# Fork-Based PRs — Design

## Problem

Every PR-opening executor (`Open`, `Cascade`, `Toolchain`, `Supersede`) does `git push -u origin <branch>` to the **upstream** repo and `gh pr create -R <repo> --head <branch>`; the update paths (`Adjust`, `Repin`, `ResolveConflict`) do `git push --force origin <branch>`. This only works where `kairos-security-bot` has **write access** — i.e. the `kairos-io` org repos. For **external** repos (`mudler/*`, `mauromorales/*`), the bot cannot push to the upstream, so the entire remediation engine silently can't act there.

The fix is the standard OSS contribution flow: **fork the external repo, push the `ksec/` branch to the fork, and open a cross-repo PR** (`--head <forkOwner>:<branch>`). Decision (locked): **external → always fork; org → push-direct** (unchanged).

## Goals

- Make all six push/PR paths work on external repos via a fork, with one small fork-aware abstraction — not seven ad-hoc edits.
- Keep org-repo behavior byte-identical (push-direct).
- The PR still lives on upstream (so reconcile/merge/close/comment by PR number are unchanged).
- Safe + deterministic: dry-run does zero writes; the fork is the bot's own; foreign upstream branches are never pushed to.

## Design

### Fork-vs-direct decision

`GitExecutor` gains:
- `ForkOwner string` — the bot's GitHub login (owner of the forks), e.g. `kairos-security-bot`.
- `ShouldFork func(repo string) bool` — `true` for external repos. `main.go` builds it from `repos.json`: `ShouldFork(repo) = kindByRepo[repo] == "external"`. **Nil-safe: a nil `ShouldFork` never forks** (preserves current behavior and all existing tests).

So the predicate is data-driven from `repos.yaml`'s `kind`; org and unknown repos push-direct, external repos fork.

### Helpers (all in `git_executor.go`)

- `func forkSlug(forkOwner, repo string) string` (pure) → `forkOwner + "/" + path.Base(repo)` (e.g. `kairos-security-bot/edgevpn`). Tested.
- `func (g *GitExecutor) prHead(repo, branch string) string` → `branch` when not forking, else `g.ForkOwner + ":" + branch`. Tested.
- `func (g *GitExecutor) forkURL(repo string) string` → token-auth URL to `forkSlug(...)` (mirrors `cloneURL`).
- `func (g *GitExecutor) ensureFork(repo string) error` → `gh repo fork <repo> --clone=false` (idempotent: a pre-existing fork is a no-op). Dry-run: print, no-op.
- `func (g *GitExecutor) pushBranch(dir, repo, branch string, force bool) error` — the single push chokepoint:
  - **Not forking (org):** `git push [-u|--force] origin <branch>` (exactly today's behavior).
  - **Forking (external):** `ensureFork(repo)`; `git remote add fork <forkURL>` (idempotent — ignore "already exists"); `git push [--force] fork <branch>`.
  - Dry-run: print intended target, no writes.
- `func (g *GitExecutor) checkoutOwnBranch(dir, repo, branch string) error` — for the **update** paths whose branch already lives on the fork (external): add the `fork` remote, `git fetch fork`, `git checkout -b <branch> fork/<branch>`. Org: `git checkout <branch>` (today's behavior).

### Per-path application

- **New-PR paths** (`Open`, `Cascade`, `Toolchain`, `Supersede`): keep cloning **upstream** (branch is cut fresh from upstream HEAD). Replace `git push -u origin <branch>` → `g.pushBranch(dir, repo, branch, false)`. Replace `gh pr create … --head <branch>` → `--head g.prHead(repo, branch)`.
- **Update paths** (`Adjust`, `Repin`): the branch already exists on the fork (external) or origin (org). Replace the post-clone `git checkout <branch>` → `g.checkoutOwnBranch(dir, repo, branch)`; replace `git push --force origin <branch>` → `g.pushBranch(dir, repo, branch, true)`.
- **`ResolveConflict`** (our own `ksec/` branch): clone upstream; `g.checkoutOwnBranch` (fork for external) to get the branch; rebase onto `origin/HEAD` (upstream default — add upstream as the `origin` remote base is already the clone origin); on success `g.pushBranch(dir, repo, branch, true)` (force to fork/origin). The `ksec/` guard and "never touch a foreign branch" invariant hold — for external we only push our own branch on **our** fork.

The PR itself is always created/viewed against **upstream** (`gh pr create -R <upstream>`, `gh pr view -R <upstream>`), so `Reconcile`/`MergePR`/`ClosePR`/`PostPRComment` (keyed by upstream PR number) are unchanged.

### Wiring (`cmd/ksec/main.go`)

- Resolve `ForkOwner` once via `gh api user --jq .login` (the token's identity); fall back to a `KSEC_FORK_OWNER` env if the call fails. Set `ex.ForkOwner`.
- Build `kindByRepo` from `repos.json`; set `ex.ShouldFork = func(repo) bool { return kindByRepo[repo] == "external" }`.

## Out of scope

- Org-repo behavior (unchanged push-direct).
- Fork cleanup/garbage-collection (forks persist; harmless).
- Keeping the fork's default branch in sync (we branch from a fresh upstream clone / rebase onto upstream, so the PR diff is always against current upstream).
- Auto-detecting write access at runtime (the `kind`-based rule is explicit and predictable).

## Risks / notes

- **Bot identity must be able to own forks.** `gh repo fork` forks to the authenticated account; `kairos-security-bot` must be a user/org account with fork permission (a GitHub **App installation token** cannot own a user fork — confirm the bot is a machine-user PAT). If forking is not permitted, external pushes fail loudly as a per-repo executor error (surfaced on the dashboard) — no silent breakage.
- **First fork is asynchronous.** `gh repo fork` may return before the fork is fully ready; the immediate `git push fork` usually succeeds (GitHub provisions quickly), but a transient failure becomes a recorded executor error and retries next run.
- **Token scope:** the fork remote uses the same `x-access-token:<GH_TOKEN>` auth as `cloneURL`; the token needs push scope on the bot's own forks (it owns them) — no upstream write needed.

## Testing

- `forkSlug` (pure): `kairos-security-bot` + `mudler/edgevpn` → `kairos-security-bot/edgevpn`.
- `prHead` (pure, via a `GitExecutor` with `ShouldFork`/`ForkOwner` set): org repo → `branch`; external repo → `kairos-security-bot:branch`; nil `ShouldFork` → `branch`.
- `ShouldFork` decision from `kindByRepo` (a small main.go-level or helper unit test): external → true, org → false, unknown → false.
- The git/gh shelling (`pushBranch`/`ensureFork`/`checkoutOwnBranch`) is integration — verified by `go build`/`go vet` and the dry-run path (prints, no writes); covered behaviorally by the existing executor flows.
- Regression: with `ShouldFork == nil` (the default in all current executor tests), every push/clone/PR path is byte-identical to today.
- Manual: a live run opens an external-repo PR from `kairos-security-bot:ksec/...` against the upstream, and `Reconcile`/merge act on the upstream PR number.
