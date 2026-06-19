package ghclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type PullRequest struct {
	Number int      `json:"number"`
	Title  string   `json:"title"`
	Author string   `json:"author"`
	URL    string   `json:"url"`
	Labels []string `json:"labels"`
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

type GitHub interface {
	ListOrgRepos(org string) ([]string, error)
	GetFile(repo, path, ref string) ([]byte, error)
	ListOpenPRs(repo string) ([]PullRequest, error)
	ListDependabotAlerts(repo string) ([]Alert, error)
	UpsertIssue(repo, marker, title, body string, labels []string) (int, error)
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
			return nil, fmt.Errorf("gh %v: %v: %s", args, err, errb.String())
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
		"--json", "number,title,author,url,labels",
		"-q", "[.[] | {number, title, author: .author.login, url, labels: [.labels[].name]}]")
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
		return nil, err
	}
	var alerts []Alert
	return alerts, json.Unmarshal(b, &alerts)
}

func (c *CLI) UpsertIssue(repo, marker, title, body string, labels []string) (int, error) {
	full := body + "\n\n" + marker
	// Find an existing issue containing the marker.
	listed, err := c.run("issue", "list", "-R", repo, "--state", "open", "--search", marker, "--limit", "1", "--json", "number", "-q", ".[].number")
	if err != nil {
		return 0, err
	}
	if lines := splitLines(listed); len(lines) > 0 {
		var n int
		fmt.Sscanf(lines[0], "%d", &n)
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
