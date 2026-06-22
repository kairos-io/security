package remediate

import (
	"regexp"
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
)

// reTitleVersion matches version tokens like "0.36.0", "v1.2", or "0.30.0" in a
// PR title.
var reTitleVersion = regexp.MustCompile(`v?\d+\.\d+(?:\.\d+)?`)

// titleSatisfiesVersion reports whether the title advertises a version >= the
// required minimum. An empty requirement is always satisfied. A title with no
// version token (e.g. "remove golang.org/x/net usage") never satisfies a
// non-empty requirement, and a bump to a version BELOW required does not match.
func titleSatisfiesVersion(title, required string) bool {
	if required == "" {
		return true
	}
	for _, tok := range reTitleVersion.FindAllString(title, -1) {
		if compareVersions(tok, required) >= 0 {
			return true
		}
	}
	return false
}

func isOwnPR(pr ghclient.PullRequest) bool {
	return strings.HasPrefix(pr.HeadRef, "ksec/") || pr.Author == "kairos-security-bot"
}

// botName strips gh's bot-author decorations to the underlying name: the
// "app/" prefix gh puts on GitHub App authors (e.g. "app/dependabot",
// "app/renovate") and the "[bot]" suffix. Both forms -> "dependabot"/"renovate".
func botName(login string) string {
	return strings.TrimSuffix(strings.TrimPrefix(login, "app/"), "[bot]")
}

func classifySource(pr ghclient.PullRequest) string {
	if isOwnPR(pr) {
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

// MatchPR returns the first open PR whose title contains the package path
// (case-insensitive) and, when version != "", advertises a version token that
// is >= the required version. Accepting a higher version lets us adopt a
// dependabot/renovate bump that overshoots our minimum fix, while a title that
// merely mentions the package (e.g. "remove golang.org/x/net usage") or bumps
// below the required version still does not match.
func MatchPR(pkg, version string, prs []ghclient.PullRequest) (ghclient.PullRequest, string, bool) {
	if pkg == "" {
		return ghclient.PullRequest{}, "", false
	}
	pkgL := strings.ToLower(pkg)
	for _, pr := range prs {
		title := strings.ToLower(pr.Title)
		if !strings.Contains(title, pkgL) {
			continue
		}
		if !titleSatisfiesVersion(pr.Title, version) {
			continue
		}
		return pr, classifySource(pr), true
	}
	return ghclient.PullRequest{}, "", false
}
