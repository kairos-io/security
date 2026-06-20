package remediate

import (
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
)

func isOwnPR(pr ghclient.PullRequest) bool {
	return strings.HasPrefix(pr.HeadRef, "ksec/") || pr.Author == "kairos-security-bot"
}

func classifySource(pr ghclient.PullRequest) string {
	if isOwnPR(pr) {
		return "ksec"
	}
	switch pr.Author {
	case "renovate[bot]":
		return "renovate"
	case "dependabot[bot]":
		return "dependabot"
	default:
		return "human"
	}
}

// MatchPR returns the first open PR whose title contains the package path
// (case-insensitive) and, when version != "", the version (leading 'v'
// stripped). Requiring the version avoids matching PRs that merely mention the
// package (e.g. "remove golang.org/x/net usage").
func MatchPR(pkg, version string, prs []ghclient.PullRequest) (ghclient.PullRequest, string, bool) {
	if pkg == "" {
		return ghclient.PullRequest{}, "", false
	}
	pkgL := strings.ToLower(pkg)
	verL := strings.ToLower(strings.TrimPrefix(version, "v"))
	for _, pr := range prs {
		title := strings.ToLower(pr.Title)
		if !strings.Contains(title, pkgL) {
			continue
		}
		if verL != "" && !strings.Contains(title, verL) {
			continue
		}
		return pr, classifySource(pr), true
	}
	return ghclient.PullRequest{}, "", false
}
