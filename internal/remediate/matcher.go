package remediate

import (
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
)

func classifySource(author string) string {
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

// MatchPR returns the first open PR whose title contains the package path
// (case-insensitive) and the PR's source. A non-empty pkg is required.
func MatchPR(pkg string, prs []ghclient.PullRequest) (ghclient.PullRequest, string, bool) {
	if pkg == "" {
		return ghclient.PullRequest{}, "", false
	}
	needle := strings.ToLower(pkg)
	for _, pr := range prs {
		if strings.Contains(strings.ToLower(pr.Title), needle) {
			return pr, classifySource(pr.Author), true
		}
	}
	return ghclient.PullRequest{}, "", false
}
