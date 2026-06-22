package collect

import (
	"sort"
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

// OpenPRs lists the CVE-tied open PRs across repos for the dashboard's PR list.
// A PR is kept only when it is security-relevant (see cveTied), so routine
// dependabot/renovate bumps with no matching finding are dropped as noise.
// These are remediation artifacts, NOT security findings, so they no longer
// enter the findings set.
func OpenPRs(repos []state.Repo, gh ghclient.GitHub, findings []state.Finding) ([]state.TrackedPR, []state.CollectionError) {
	pkgsByRepo := map[string][]string{}
	for _, f := range findings {
		if f.Package != "" {
			pkgsByRepo[f.Repo] = append(pkgsByRepo[f.Repo], strings.ToLower(f.Package))
		}
	}
	var out []state.TrackedPR
	var errs []state.CollectionError
	for _, repo := range repos {
		prs, err := gh.ListOpenPRs(repo.Repo)
		if err != nil {
			errs = append(errs, state.CollectionError{Repo: repo.Repo, Collector: "prs", Message: err.Error()})
			continue
		}
		for _, pr := range prs {
			if !cveTied(pr, pkgsByRepo[repo.Repo]) {
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

// cveTied keeps a PR only when it is security-relevant: it bumps a package that
// has a finding (CVE) in this repo, OR carries a `security` label, OR is ours.
func cveTied(pr ghclient.PullRequest, findingPkgs []string) bool {
	if pr.Author == "kairos-security-bot" || strings.HasPrefix(pr.HeadRef, "ksec/") {
		return true
	}
	for _, l := range pr.Labels {
		if l == "security" {
			return true
		}
	}
	title := strings.ToLower(pr.Title)
	for _, pkg := range findingPkgs {
		if strings.Contains(title, pkg) {
			return true
		}
	}
	return false
}
