package remediate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

var _ Executor = (*GitExecutor)(nil)

type GitExecutor struct {
	Token     string // GH_TOKEN, for authenticated clone/push
	DryRun    bool
	Prose     ProseClient     // optional; nil -> deterministic PR body
	GH        ghclient.GitHub // used by Adopt for comment/status/merge
	Automerge bool
	Agent     Agent // optional; when set, attempts to repair a broken build before giving up
	// ForkOwner is the account that owns the bot's forks; when forking, PRs are
	// opened cross-fork (head "ForkOwner:branch") and pushes go to the fork.
	ForkOwner string
	// ShouldFork decides whether a given repo must be remediated via a fork
	// (e.g. EXTERNAL repos the bot can't push to). nil means never fork.
	ShouldFork func(repo string) bool
}

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

// verifyOrRepair runs `go build ./...`; on failure it asks the agent (if any)
// to repair the code and re-verifies. Returns true iff the tree builds (which
// the caller requires before pushing).
func (g *GitExecutor) verifyOrRepair(dir, task, runID string) bool {
	if _, err := g.run(dir, "go", "build", "./..."); err == nil {
		return true
	} else if g.Agent == nil {
		return false
	} else {
		_ = g.Agent.Repair(dir, RepairTask(err.Error()))
	}
	_, err := g.run(dir, "go", "build", "./...")
	return err == nil
}

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

func (g *GitExecutor) cloneURL(repo string) string {
	if g.Token != "" {
		return "https://x-access-token:" + g.Token + "@github.com/" + repo + ".git"
	}
	return "https://github.com/" + repo + ".git"
}

func (g *GitExecutor) run(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		msg := fmt.Sprintf("%s %v: %v: %s", name, args, err, errb.String())
		// Never leak the live token into errors that flow to the ledger or CI logs.
		if g.Token != "" {
			msg = strings.ReplaceAll(msg, g.Token, "***")
		}
		return out.Bytes(), errors.New(msg)
	}
	return out.Bytes(), nil
}

func (g *GitExecutor) Open(in Intent, runID string) (state.LedgerEntry, error) {
	branch := BranchName(in)
	entry := state.LedgerEntry{
		Key: in.Key, Repo: in.Repo, Package: in.Package, Branch: branch,
		Bump: in.Bump, Severity: in.Severity, CreatedRun: runID, LastActionRun: runID,
	}

	if g.DryRun {
		fmt.Printf("[dry-run] would open PR on %s: branch %s, go get %s@%s\n",
			in.Repo, branch, in.Bump.Package, in.Bump.To)
		entry.State = "planned"
		entry.History = []state.LedgerEvent{{Run: runID, Action: "plan-open"}}
		return entry, nil
	}

	dir, err := os.MkdirTemp("", "ksec-rem-*")
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
	if _, err := g.run(dir, "go", "get", in.Bump.Package+"@"+in.Bump.To); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "go", "mod", "tidy"); err != nil {
		return entry, err
	}
	// Verify-before-push: a broken build must not be pushed.
	if !g.verifyOrRepair(dir, "open "+in.Package, runID) {
		entry.State = "build-failed"
		entry.NeedsHuman = true
		entry.History = []state.LedgerEvent{{Run: runID, Action: "build-failed"}}
		return entry, nil // not an error: recorded for a human, run continues
	}

	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	if _, err := g.run(dir, "git", "commit", "-am", PRTitle(in)); err != nil {
		return entry, err
	}
	if err := g.pushBranch(dir, in.Repo, branch, false); err != nil {
		return entry, err
	}

	// Create the PR with gh (GH_TOKEN is read from the environment by gh).
	out, err := g.run(dir, "gh", "pr", "create", "-R", in.Repo, "--head", g.prHead(in.Repo, branch),
		"--title", PRTitle(in), "--body", PRBodyWith(in, g.Prose))
	if err != nil {
		return entry, err
	}
	entry.PRURL = string(bytes.TrimSpace(out))
	entry.PRNumber = prNumberFromURL(entry.PRURL)
	entry.State = "open"
	entry.History = []state.LedgerEvent{{Run: runID, Action: "opened", Detail: entry.PRURL}}
	return entry, nil
}

const nudgeMarker = "<!-- ksec:nudge -->"

func (g *GitExecutor) Adopt(in Intent, runID string) (state.LedgerEntry, error) {
	entry := state.LedgerEntry{
		Key: in.Key, Repo: in.Repo, Package: in.Package, Source: in.Source, Kind: "direct",
		PRNumber: in.PRNumber, PRURL: in.PRURL, Bump: in.Bump, Severity: in.Severity,
		State: "open", CreatedRun: runID, LastActionRun: runID,
	}
	if g.DryRun || g.GH == nil {
		if g.DryRun {
			fmt.Printf("[dry-run] would adopt %s PR #%d (%s): nudge%s\n", in.Repo, in.PRNumber, in.Source,
				map[bool]string{true: " + automerge-if-green", false: ""}[g.Automerge])
		}
		entry.History = []state.LedgerEvent{{Run: runID, Action: "adopt", Detail: in.Source}}
		return entry, nil
	}

	// Refresh live PR state.
	if st, err := g.GH.PRStatusOf(in.Repo, in.PRNumber); err == nil {
		switch st.State {
		case "MERGED":
			entry.State = "merged"
		case "CLOSED":
			entry.State = "closed"
		}
		// Optional automerge.
		if g.Automerge && entry.State == "open" && ShouldAutomerge(st) {
			if err := g.GH.MergePR(in.Repo, in.PRNumber, true); err == nil {
				entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "automerge-requested"})
			}
		}
	}

	// Idempotent nudge: only if we haven't commented the marker yet.
	if entry.State == "open" {
		nudged := false
		if comments, err := g.GH.ListPRComments(in.Repo, in.PRNumber); err == nil {
			for _, c := range comments {
				if strings.Contains(c.Body, nudgeMarker) {
					nudged = true
					break
				}
			}
		}
		if !nudged {
			body := fmt.Sprintf("This PR addresses a %s-severity security finding (%s). Tracked by kairos-security.\n\n%s",
				in.Severity, in.Package, nudgeMarker)
			_ = g.GH.PostPRComment(in.Repo, in.PRNumber, body)
			entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "nudged"})
		}
	}
	return entry, nil
}

func (g *GitExecutor) Reconcile(e state.LedgerEntry, runID string) (state.LedgerEntry, error) {
	if e.PRNumber == 0 || g.DryRun {
		if g.DryRun {
			fmt.Printf("[dry-run] would reconcile %s (PR #%d)\n", e.Repo, e.PRNumber)
		}
		return e, nil
	}
	out, err := g.run("", "gh", "pr", "view", fmt.Sprint(e.PRNumber), "-R", e.Repo,
		"--json", "state,mergedAt,mergeable,headRefName",
		"-q", "{state: .state, mergedAt: .mergedAt, mergeable: .mergeable, headRef: .headRefName}")
	if err != nil {
		return e, err
	}
	var view struct {
		State     string `json:"state"`
		MergedAt  string `json:"mergedAt"`
		Mergeable string `json:"mergeable"`
		HeadRef   string `json:"headRef"`
	}
	_ = json.Unmarshal(out, &view)
	prior := e.State
	switch {
	case view.MergedAt != "" || view.State == "MERGED":
		e.State = "merged"
	case view.State == "CLOSED":
		e.State = "closed"
	default:
		e.State = "open"
	}
	// Only record a change when the state actually changed. An unchanged PR
	// must leave the ledger byte-identical so the volatile run id doesn't churn
	// it across runs.
	if e.State != prior {
		e.LastActionRun = runID
		e.History = append(e.History, state.LedgerEvent{Run: runID, Action: "reconciled", Detail: e.State})
	}
	// An owned, open PR reporting conflicts is a resolution candidate: rebase it
	// onto the base branch (agent-assisted) and force-push a building tree.
	branch := e.Branch
	if branch == "" {
		branch = view.HeadRef
	}
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
}

var _ Adjuster = (*GitExecutor)(nil)

func (g *GitExecutor) Adjust(entry state.LedgerEntry, toVersion, runID string) (state.LedgerEntry, error) {
	// Hard invariant: we only ever force-push bot-managed branches. Guard first,
	// before any clone or push, so a corrupted ledger can't make us rewrite a
	// real branch (e.g. main).
	if !strings.HasPrefix(entry.Branch, "ksec/") {
		return entry, fmt.Errorf("refusing to adjust non-ksec branch %q", entry.Branch)
	}
	entry.LastActionRun = runID
	if g.DryRun {
		fmt.Printf("[dry-run] would adjust %s PR #%d: go get %s@%s, force-push %s\n",
			entry.Repo, entry.PRNumber, entry.Package, toVersion, entry.Branch)
		entry.Bump.To = toVersion
		return entry, nil
	}
	dir, err := os.MkdirTemp("", "ksec-adj-*")
	if err != nil {
		return entry, err
	}
	defer os.RemoveAll(dir)

	if _, err := g.run("", "git", "clone", g.cloneURL(entry.Repo), dir); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "git", "checkout", entry.Branch); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "go", "get", entry.Package+"@"+toVersion); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "go", "mod", "tidy"); err != nil {
		return entry, err
	}
	if !g.verifyOrRepair(dir, "adjust "+entry.Package, runID) {
		entry.State = "build-failed"
		entry.NeedsHuman = true
		entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "adjust-build-failed"})
		return entry, nil
	}
	// If the branch is already at the requested version, go get/tidy produce no
	// diff. A `git commit -am` would then fail and we'd re-clone and retry every
	// run forever. Detect the no-op and return without committing or pushing.
	porcelain, err := g.run(dir, "git", "status", "--porcelain")
	if err != nil {
		return entry, err
	}
	if len(bytes.TrimSpace(porcelain)) == 0 {
		entry.Bump.To = toVersion
		entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "already-current"})
		return entry, nil
	}
	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	if _, err := g.run(dir, "git", "commit", "-am", "chore(security): adjust bump to "+toVersion); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "git", "push", "--force", "origin", entry.Branch); err != nil {
		return entry, err
	}
	entry.Bump.To = toVersion
	entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "adjusted", Detail: "to " + toVersion})
	return entry, nil
}

// ResolveConflict rebases an owned, conflicted PR branch onto the repo's
// default branch, using the agent to resolve conflicts if needed, then
// force-pushes a building tree. Only ever touches ksec/ branches.
func (g *GitExecutor) ResolveConflict(e state.LedgerEntry, runID string) (state.LedgerEntry, error) {
	// Hard invariant: we only ever force-push bot-managed branches. Guard first,
	// before any network, so a corrupted ledger can't make us rewrite a real
	// branch (e.g. main).
	if !strings.HasPrefix(e.Branch, "ksec/") {
		return e, fmt.Errorf("refusing to resolve non-ksec branch %q", e.Branch)
	}
	if g.DryRun {
		fmt.Printf("[dry-run] would resolve conflict on %s (PR #%d)\n", e.Repo, e.PRNumber)
		return e, nil
	}
	dir, err := os.MkdirTemp("", "ksec-cfl-*")
	if err != nil {
		return e, err
	}
	defer os.RemoveAll(dir)

	// Full clone (no --depth) so origin/HEAD resolves to the default branch and
	// a rebase has the history it needs.
	if _, err := g.run("", "git", "clone", g.cloneURL(e.Repo), dir); err != nil {
		return e, err
	}
	if _, err := g.run(dir, "git", "checkout", e.Branch); err != nil {
		return e, err
	}
	// Set a committer identity before any rebase: `git rebase origin/HEAD` and
	// `git rebase --continue` create commits, and a fresh clone has no identity
	// configured in bare CI.
	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	if _, err := g.run(dir, "git", "fetch", "origin"); err != nil {
		return e, err
	}
	// Rebase onto the default branch. origin/HEAD points at it after a full clone.
	if _, rebaseErr := g.run(dir, "git", "rebase", "origin/HEAD"); rebaseErr != nil {
		// Conflicts (or no agent): try agent-assisted resolution.
		resolved := false
		if g.Agent != nil {
			if aerr := g.Agent.Repair(dir, ConflictTask()); aerr == nil {
				_, _ = g.run(dir, "git", "add", "-A")
				if _, cerr := g.run(dir, "git", "-c", "core.editor=true", "rebase", "--continue"); cerr == nil {
					resolved = true
				}
			}
		}
		if !resolved {
			_, _ = g.run(dir, "git", "rebase", "--abort") // best-effort
			e.NeedsHuman = true
			e.Blocked = "conflict"
			e.History = append(e.History, state.LedgerEvent{Run: runID, Action: "conflict-unresolved"})
			return e, nil
		}
	}
	// Verify-before-push: never force-push a non-building tree.
	if !g.verifyOrRepair(dir, "resolve conflict "+e.Repo, runID) {
		e.NeedsHuman = true
		e.Blocked = "conflict-build-failed"
		e.History = append(e.History, state.LedgerEvent{Run: runID, Action: "conflict-build-failed"})
		return e, nil
	}
	if _, err := g.run(dir, "git", "push", "--force", "origin", e.Branch); err != nil {
		return e, err
	}
	e.State = "open"
	e.Blocked = ""
	e.NeedsHuman = false
	e.LastActionRun = runID
	e.History = append(e.History, state.LedgerEvent{Run: runID, Action: "conflict-resolved"})
	return e, nil
}

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
	if !g.verifyOrRepair(dir, "cascade "+in.Package, runID) {
		entry.State = "build-failed"
		entry.NeedsHuman = true
		entry.History = []state.LedgerEvent{{Run: runID, Action: "cascade-build-failed"}}
		return entry, nil
	}
	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	if _, err := g.run(dir, "git", "commit", "-am", "chore(security): cascade-bump "+in.Package); err != nil {
		return entry, err
	}
	if err := g.pushBranch(dir, in.Repo, branch, false); err != nil {
		return entry, err
	}
	out, err := g.run(dir, "gh", "pr", "create", "-R", in.Repo, "--head", g.prHead(in.Repo, branch),
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

func (g *GitExecutor) Repin(e state.LedgerEntry, runID string) (state.LedgerEntry, error) {
	module := e.Package
	if g.DryRun {
		fmt.Printf("[dry-run] would check %s for a release tag to re-pin %s\n", module, e.Repo)
		return e, nil
	}
	// Never re-pin / force-push a merged or closed PR: only an open cascade PR
	// is a live re-pin candidate.
	if e.State != "open" {
		return e, nil
	}
	// Hard invariant: we only ever force-push bot-managed branches. Guard before
	// any go list / clone / push so a corrupted ledger can't make us rewrite a
	// real branch (e.g. main).
	if !strings.HasPrefix(e.Branch, "ksec/") {
		return e, fmt.Errorf("refusing to repin non-ksec branch %q", e.Branch)
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
	if !g.verifyOrRepair(dir, "repin "+module, runID) {
		e.State = "build-failed"
		e.NeedsHuman = true
		e.History = append(e.History, state.LedgerEvent{Run: runID, Action: "repin-build-failed"})
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

func (g *GitExecutor) Toolchain(in Intent, runID string) (state.LedgerEntry, error) {
	branch := "ksec/toolchain-" + slug(in.ToolchainVersion)
	entry := state.LedgerEntry{
		Key: in.Key, Repo: in.Repo, Package: "go-toolchain", Branch: branch, Kind: "toolchain",
		Severity: in.Severity, CreatedRun: runID, LastActionRun: runID,
		Bump: state.Bump{Package: "go", To: in.ToolchainVersion},
	}
	if g.DryRun {
		fmt.Printf("[dry-run] would bump go toolchain in %s to %s\n", in.Repo, in.ToolchainVersion)
		entry.State = "planned"
		return entry, nil
	}
	dir, err := os.MkdirTemp("", "ksec-tc-*")
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
	if _, err := g.run(dir, "go", "mod", "edit", "-go="+in.ToolchainVersion); err != nil {
		return entry, err
	}
	_, _ = g.run(dir, "go", "mod", "tidy")
	if !g.verifyOrRepair(dir, "go toolchain bump to "+in.ToolchainVersion, runID) {
		entry.State = "build-failed"
		entry.NeedsHuman = true
		entry.History = []state.LedgerEvent{{Run: runID, Action: "toolchain-build-failed"}}
		return entry, nil
	}
	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	// No-op guard: if go.mod was already at this go directive, there is nothing
	// to commit. Committing an empty diff fails and forces a re-clone every run,
	// so mark the entry already-satisfied and return without committing/pushing.
	if out, _ := g.run(dir, "git", "status", "--porcelain"); len(bytes.TrimSpace(out)) == 0 {
		entry.State = "merged"
		entry.History = []state.LedgerEvent{{Run: runID, Action: "toolchain-already-current"}}
		return entry, nil
	}
	if _, err := g.run(dir, "git", "commit", "-am", "chore(security): bump go toolchain to "+in.ToolchainVersion); err != nil {
		return entry, err
	}
	if err := g.pushBranch(dir, in.Repo, branch, false); err != nil {
		return entry, err
	}
	out, err := g.run(dir, "gh", "pr", "create", "-R", in.Repo, "--head", g.prHead(in.Repo, branch),
		"--title", "chore(security): bump go toolchain to "+in.ToolchainVersion,
		"--body", "Bumps the Go toolchain to "+in.ToolchainVersion+" to address a stdlib vulnerability. "+PRMarker(in.Key))
	if err != nil {
		return entry, err
	}
	entry.PRURL = strings.TrimSpace(string(out))
	entry.PRNumber = prNumberFromURL(entry.PRURL)
	entry.State = "open"
	entry.History = []state.LedgerEvent{{Run: runID, Action: "toolchain-opened", Detail: entry.PRURL}}
	return entry, nil
}

func (g *GitExecutor) Supersede(in Intent, runID string) (state.LedgerEntry, error) {
	branch := BranchName(in) // fresh deterministic ksec/ bump branch
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
	if err := g.pushBranch(dir, in.Repo, branch, false); err != nil {
		return entry, err
	}
	out, err := g.run(dir, "gh", "pr", "create", "-R", in.Repo, "--head", g.prHead(in.Repo, branch),
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

func prNumberFromURL(url string) int {
	n := 0
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] < '0' || url[i] > '9' {
			if i == len(url)-1 {
				return 0
			}
			fmt.Sscanf(url[i+1:], "%d", &n)
			return n
		}
	}
	return n
}
