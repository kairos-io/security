package collect

import (
	"sort"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

var botAuthors = map[string]bool{
	"renovate[bot]": true, "dependabot[bot]": true, "kairos-security-bot": true,
}
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
				Author: pr.Author, URL: pr.URL, Source: prSource(pr.Author),
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

func prSource(author string) string {
	switch author {
	case "renovate[bot]":
		return "renovate"
	case "dependabot[bot]":
		return "dependabot"
	case "kairos-security-bot":
		return "ksec"
	default:
		return "human"
	}
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
