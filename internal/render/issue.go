package render

import (
	"fmt"

	"github.com/kairos-io/security/internal/ghclient"
)

const (
	IssueMarker = "<!-- ksec:dashboard -->"
	IssueTitle  = "Kairos Security Dashboard"
)

// IssueLabels are asked for in the tracking repo, most-general first: the
// existing-issue lookup filters on the last one. `area/security` is the label
// kairos-io/kairos actually uses; a plain `security` label does not exist
// there, and asking for one that the bot also cannot create meant the upsert
// failed outright on every run. ghclient drops whatever the repo does not
// have, so this list is a request, not a requirement.
var IssueLabels = []string{"area/security", "kairos-security-bot"}

func UpsertTrackingIssue(gh ghclient.GitHub, repo, body string, dryRun bool) (int, error) {
	if dryRun {
		fmt.Printf("[dry-run] would upsert tracking issue in %s (%d bytes)\n", repo, len(body))
		return 0, nil
	}
	return gh.UpsertIssue(repo, IssueMarker, IssueTitle, body, IssueLabels)
}
