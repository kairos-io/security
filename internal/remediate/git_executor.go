package remediate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

var _ Executor = (*GitExecutor)(nil)

type GitExecutor struct {
	Token  string // GH_TOKEN, for authenticated clone/push
	DryRun bool
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
	if _, err := g.run(dir, "go", "build", "./..."); err != nil {
		entry.State = "build-failed"
		entry.History = []state.LedgerEvent{{Run: runID, Action: "build-failed", Detail: err.Error()}}
		return entry, nil // not an error: recorded for a human, run continues
	}

	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	if _, err := g.run(dir, "git", "commit", "-am", PRTitle(in)); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "git", "push", "-u", "origin", branch); err != nil {
		return entry, err
	}

	// Create the PR with gh (GH_TOKEN is read from the environment by gh).
	out, err := g.run(dir, "gh", "pr", "create", "-R", in.Repo, "--head", branch,
		"--title", PRTitle(in), "--body", PRBody(in))
	if err != nil {
		return entry, err
	}
	entry.PRURL = string(bytes.TrimSpace(out))
	entry.PRNumber = prNumberFromURL(entry.PRURL)
	entry.State = "open"
	entry.History = []state.LedgerEvent{{Run: runID, Action: "opened", Detail: entry.PRURL}}
	return entry, nil
}

func (g *GitExecutor) Reconcile(e state.LedgerEntry, runID string) (state.LedgerEntry, error) {
	e.LastActionRun = runID
	if e.PRNumber == 0 || g.DryRun {
		if g.DryRun {
			fmt.Printf("[dry-run] would reconcile %s (PR #%d)\n", e.Repo, e.PRNumber)
		}
		return e, nil
	}
	out, err := g.run("", "gh", "pr", "view", fmt.Sprint(e.PRNumber), "-R", e.Repo,
		"--json", "state,mergedAt", "-q", "{state: .state, mergedAt: .mergedAt}")
	if err != nil {
		return e, err
	}
	var view struct {
		State    string `json:"state"`
		MergedAt string `json:"mergedAt"`
	}
	_ = json.Unmarshal(out, &view)
	switch {
	case view.MergedAt != "" || view.State == "MERGED":
		e.State = "merged"
	case view.State == "CLOSED":
		e.State = "closed"
	default:
		e.State = "open"
	}
	e.History = append(e.History, state.LedgerEvent{Run: runID, Action: "reconciled", Detail: e.State})
	return e, nil
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
