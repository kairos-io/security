package collect

import (
	"sort"
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

var secLabels = map[string]bool{"security": true, "dependencies": true}

// OpenPRs lists the tracked open PRs (bot-authored or security/dependencies
// labelled) across repos for the dashboard's PR list. These are remediation
// artifacts, NOT security findings, so they no longer enter the findings set.
func OpenPRs(repos []state.Repo, gh ghclient.GitHub) ([]state.TrackedPR, []state.CollectionError) {
	var out []state.TrackedPR
	var errs []state.CollectionError
	for _, repo := range repos {
		prs, err := gh.ListOpenPRs(repo.Repo)
		if err != nil {
			errs = append(errs, state.CollectionError{Repo: repo.Repo, Collector: "prs", Message: err.Error()})
			continue
		}
		for _, pr := range prs {
			if !isSecurityPR(pr) {
				continue
			}
			out = append(out, state.TrackedPR{
				Repo: repo.Repo, Number: pr.Number, Title: pr.Title,
				Author: pr.Author, URL: pr.URL, Source: prSource(pr),
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Repo != out[j].Repo {
			return out[i].Repo < out[j].Repo
		}
		return out[i].Number < out[j].Number
	})
	return out, errs
}

// botName strips gh's bot-author decorations: the "app/" prefix gh puts on
// GitHub App authors (e.g. "app/dependabot", "app/renovate") and the "[bot]"
// suffix. Both forms collapse to the underlying name.
func botName(login string) string {
	return strings.TrimSuffix(strings.TrimPrefix(login, "app/"), "[bot]")
}

func prSource(pr ghclient.PullRequest) string {
	if pr.Author == "kairos-security-bot" || strings.HasPrefix(pr.HeadRef, "ksec/") {
		return "ksec"
	}
	switch botName(pr.Author) {
	case "renovate":
		return "renovate"
	case "dependabot":
		return "dependabot"
	}
	if pr.IsBot {
		return "bot"
	}
	return "human"
}

func isSecurityPR(pr ghclient.PullRequest) bool {
	if pr.IsBot || pr.Author == "kairos-security-bot" {
		return true
	}
	for _, l := range pr.Labels {
		if secLabels[l] {
			return true
		}
	}
	return false
}
