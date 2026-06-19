package collect

import (
	"fmt"
	"sort"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

type PRs struct {
	GH ghclient.GitHub
}

func (PRs) Name() string { return "prs" }

var botAuthors = map[string]bool{
	"renovate[bot]": true, "dependabot[bot]": true, "kairos-security-bot": true,
}
var secLabels = map[string]bool{"security": true, "dependencies": true}

func (c PRs) Collect(repo state.Repo) ([]state.Finding, error) {
	prs, err := c.GH.ListOpenPRs(repo.Repo)
	if err != nil {
		return nil, err
	}
	var out []state.Finding
	for _, pr := range prs {
		if !isSecurityPR(pr) {
			continue
		}
		out = append(out, state.Finding{
			ID:        FindingID(repo.Repo, "pr", fmt.Sprintf("#%d", pr.Number), ""),
			Repo:      repo.Repo,
			Type:      "pr",
			Severity:  "unknown",
			Source:    "github-pr",
			Title:     fmt.Sprintf("#%d %s (@%s)", pr.Number, pr.Title, pr.Author),
			URL:       pr.URL,
			FirstSeen: Today(),
			LastSeen:  Today(),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func isSecurityPR(pr ghclient.PullRequest) bool {
	if botAuthors[pr.Author] {
		return true
	}
	for _, l := range pr.Labels {
		if secLabels[l] {
			return true
		}
	}
	return false
}
