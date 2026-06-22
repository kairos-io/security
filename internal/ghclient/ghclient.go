package ghclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type PullRequest struct {
	Number  int      `json:"number"`
	Title   string   `json:"title"`
	Author  string   `json:"author"`
	IsBot   bool     `json:"isBot"` // gh author.is_bot — authoritative bot signal
	URL     string   `json:"url"`
	HeadRef string   `json:"headRef"`
	Labels  []string `json:"labels"`
}

type Alert struct {
	Number       int    `json:"number"`
	CVEID        string `json:"cveID"`
	GHSA         string `json:"ghsa"`
	Package      string `json:"package"`
	Ecosystem    string `json:"ecosystem"`
	Severity     string `json:"severity"`
	URL          string `json:"url"`
	FixedVersion string `json:"fixedVersion"`
}

type ReviewComment struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Body   string `json:"body"`
}

type PRStatus struct {
	State          string `json:"state"`
	Mergeable      bool   `json:"mergeable"`
	ChecksPassing  bool   `json:"checksPassing"`
	ReviewDecision string `json:"reviewDecision"`
}

type GitHub interface {
	ListOrgRepos(org string) ([]string, error)
	GetFile(repo, path, ref string) ([]byte, error)
	ListOpenPRs(repo string) ([]PullRequest, error)
	ListDependabotAlerts(repo string) ([]Alert, error)
	UpsertIssue(repo, marker, title, body string, labels []string) (int, error)
	ListPRComments(repo string, pr int) ([]ReviewComment, error)
	PostPRComment(repo string, pr int, body string) error
	ClosePR(repo string, pr int, comment string) error
	PRStatusOf(repo string, pr int) (PRStatus, error)
	MergePR(repo string, pr int, auto bool) error
}

// CLI is the production GitHub client; it shells out to `gh`.
type CLI struct {
	run func(args ...string) ([]byte, error)
}

func NewCLI() *CLI {
	return &CLI{run: func(args ...string) ([]byte, error) {
		cmd := exec.Command("gh", args...)
		var out, errb bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &errb
		if err := cmd.Run(); err != nil {
			if stderr := strings.TrimSpace(errb.String()); stderr != "" {
				return nil, fmt.Errorf("gh %s: %s", args[0], stderr)
			}
			return nil, fmt.Errorf("gh %s: %v", args[0], err)
		}
		return out.Bytes(), nil
	}}
}

func (c *CLI) api(path string, jqOrFields ...string) ([]byte, error) {
	args := append([]string{"api", path}, jqOrFields...)
	return c.run(args...)
}

func (c *CLI) ListOrgRepos(org string) ([]string, error) {
	b, err := c.run("repo", "list", org, "--no-archived", "--limit", "1000", "--json", "nameWithOwner", "-q", ".[].nameWithOwner")
	if err != nil {
		return nil, err
	}
	return splitLines(b), nil
}

func (c *CLI) GetFile(repo, path, ref string) ([]byte, error) {
	// gh api returns the raw content with the proper Accept header.
	return c.run("api", fmt.Sprintf("repos/%s/contents/%s?ref=%s", repo, path, ref),
		"-H", "Accept: application/vnd.github.raw+json")
}

func (c *CLI) ListOpenPRs(repo string) ([]PullRequest, error) {
	b, err := c.run("pr", "list", "-R", repo, "--state", "open", "--limit", "200",
		"--json", "number,title,author,url,headRefName,labels",
		"-q", "[.[] | {number, title, author: .author.login, isBot: .author.is_bot, url, headRef: .headRefName, labels: [.labels[].name]}]")
	if err != nil {
		return nil, err
	}
	var prs []PullRequest
	return prs, json.Unmarshal(b, &prs)
}

func (c *CLI) ListDependabotAlerts(repo string) ([]Alert, error) {
	b, err := c.api(fmt.Sprintf("repos/%s/dependabot/alerts?state=open&per_page=100", repo),
		"-q", "[.[] | {number, cveID: (.security_advisory.cve_id // \"\"), ghsa: .security_advisory.ghsa_id, package: .dependency.package.name, ecosystem: .dependency.package.ecosystem, severity: .security_advisory.severity, url: .html_url, fixedVersion: (.security_vulnerability.first_patched_version.identifier // \"\")}]")
	if err != nil {
		// Dependabot may be disabled, the repo archived, or the token may
		// lack the required scope. GitHub answers all of these with 403;
		// treat them as "no alerts" rather than a collection error.
		msg := err.Error()
		for _, s := range []string{"403", "Dependabot alerts are disabled", "not available", "not authorized", "admin:repo_hook"} {
			if strings.Contains(msg, s) {
				return nil, nil
			}
		}
		return nil, err
	}
	var alerts []Alert
	return alerts, json.Unmarshal(b, &alerts)
}

func (c *CLI) UpsertIssue(repo, marker, title, body string, labels []string) (int, error) {
	full := body + "\n\n" + marker
	// Ensure each label exists up front (best-effort, idempotent): the list
	// below filters by `--label`, and `gh issue create` fails outright on an
	// unknown label. `--force` updates an existing label instead of erroring.
	// Errors are ignored so a missing label-create permission degrades
	// silently rather than emitting a `'<label>' not found` warning.
	for _, l := range labels {
		_, _ = c.run("label", "create", l, "-R", repo, "--color", "ededed", "--description", "kairos-security", "--force")
	}
	// Find an existing issue deterministically by label + exact title.
	// Full-text search does not reliably match text inside HTML comments,
	// so we cannot rely on the marker for lookup; the marker remains in the
	// body only as an in-body sentinel.
	listArgs := []string{"issue", "list", "-R", repo, "--state", "open",
		"--search", title + " in:title", "--json", "number,title", "-q", "."}
	if len(labels) > 0 {
		// Filter by the last label (the bot label) for a tighter match.
		listArgs = append(listArgs, "--label", labels[len(labels)-1])
	}
	listed, err := c.run(listArgs...)
	if err != nil {
		return 0, err
	}
	var found []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}
	if len(bytes.TrimSpace(listed)) > 0 {
		if err := json.Unmarshal(listed, &found); err != nil {
			return 0, err
		}
	}
	// Guard against `in:title` fuzzy matches: require an exact title match.
	var matches []int
	for _, f := range found {
		if f.Title == title {
			matches = append(matches, f.Number)
		}
	}
	if len(matches) == 1 {
		n := matches[0]
		_, err := c.run("issue", "edit", fmt.Sprint(n), "-R", repo, "--body", full)
		return n, err
	}
	args := []string{"issue", "create", "-R", repo, "--title", title, "--body", full}
	for _, l := range labels {
		args = append(args, "--label", l)
	}
	out, err := c.run(args...)
	if err != nil {
		return 0, err
	}
	return parseIssueNumberFromURL(out), nil
}

func (c *CLI) ListPRComments(repo string, pr int) ([]ReviewComment, error) {
	b, err := c.run("pr", "view", fmt.Sprint(pr), "-R", repo, "--json", "comments",
		"-q", "[.comments[] | {id: (.id|tostring), author: .author.login, body: .body}]")
	if err != nil {
		return nil, err
	}
	var out []ReviewComment
	return out, json.Unmarshal(b, &out)
}

func (c *CLI) PostPRComment(repo string, pr int, body string) error {
	_, err := c.run("pr", "comment", fmt.Sprint(pr), "-R", repo, "--body", body)
	return err
}

func (c *CLI) ClosePR(repo string, pr int, comment string) error {
	_, err := c.run("pr", "close", fmt.Sprint(pr), "-R", repo, "--comment", comment)
	return err
}

func (c *CLI) PRStatusOf(repo string, pr int) (PRStatus, error) {
	b, err := c.run("pr", "view", fmt.Sprint(pr), "-R", repo,
		"--json", "state,mergeable,reviewDecision,statusCheckRollup",
		"-q", "{state: .state, mergeable: (.mergeable == \"MERGEABLE\"), reviewDecision: (.reviewDecision // \"\"), "+
			"checksPassing: ([.statusCheckRollup[]? | select((.conclusion // .state) as $s | $s != \"SUCCESS\" and $s != \"NEUTRAL\" and $s != \"SKIPPED\")] | length == 0)}")
	if err != nil {
		return PRStatus{}, err
	}
	var s PRStatus
	return s, json.Unmarshal(b, &s)
}

func (c *CLI) MergePR(repo string, pr int, auto bool) error {
	args := []string{"pr", "merge", fmt.Sprint(pr), "-R", repo, "--squash"}
	if auto {
		args = append(args, "--auto")
	}
	_, err := c.run(args...)
	return err
}

func splitLines(b []byte) []string {
	var out []string
	for _, line := range bytes.Split(bytes.TrimSpace(b), []byte("\n")) {
		if s := string(bytes.TrimSpace(line)); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func parseIssueNumberFromURL(b []byte) int {
	// `gh issue create` prints the new issue URL ending in /<number>.
	s := string(bytes.TrimSpace(b))
	var n int
	if i := bytes.LastIndexByte([]byte(s), '/'); i >= 0 {
		fmt.Sscanf(s[i+1:], "%d", &n)
	}
	return n
}
