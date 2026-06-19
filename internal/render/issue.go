package render

import (
	"fmt"

	"github.com/kairos-io/security/internal/ghclient"
)

const (
	IssueMarker = "<!-- ksec:dashboard -->"
	IssueTitle  = "Kairos Security Dashboard"
)

var IssueLabels = []string{"security", "kairos-security-bot"}

func UpsertTrackingIssue(gh ghclient.GitHub, repo, body string, dryRun bool) (int, error) {
	if dryRun {
		fmt.Printf("[dry-run] would upsert tracking issue in %s (%d bytes)\n", repo, len(body))
		return 0, nil
	}
	return gh.UpsertIssue(repo, IssueMarker, IssueTitle, body, IssueLabels)
}
