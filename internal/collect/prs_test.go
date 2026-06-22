package collect

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakePRGH struct {
	byRepo map[string][]ghclient.PullRequest
	ghclient.GitHub
}

func (f fakePRGH) ListOpenPRs(repo string) ([]ghclient.PullRequest, error) {
	return f.byRepo[repo], nil
}

func TestOpenPRsTracksOnlyCVETied(t *testing.T) {
	findings := []state.Finding{
		{Repo: "o/r", Package: "golang.org/x/crypto", CVEID: "GO-1", Severity: "high"},
	}
	gh := fakePRGH{byRepo: map[string][]ghclient.PullRequest{
		"o/r": {
			{Number: 2, Title: "Bump golang.org/x/crypto from 0.39.0 to 0.45.0", Author: "app/dependabot", IsBot: true, URL: "u2"},
			{Number: 5, Title: "Bump github.com/foo/bar to 1.2.3", Author: "app/dependabot", IsBot: true, URL: "u5"}, // no matching finding -> noise, dropped
			{Number: 7, Title: "security hardening", Author: "alice", Labels: []string{"security"}, URL: "u7"},
			{Number: 9, Title: "ksec bump", Author: "someone", URL: "u9", HeadRef: "ksec/bump-x"},
		},
	}}
	prs, errs := OpenPRs([]state.Repo{{Repo: "o/r"}}, gh, findings)
	require.Empty(t, errs)
	nums := map[int]string{}
	for _, p := range prs {
		nums[p.Number] = p.Source
	}
	assert.Contains(t, nums, 2) // tied to the x/crypto finding
	assert.Equal(t, "dependabot", nums[2])
	assert.NotContains(t, nums, 5) // unrelated bump -> dropped (no noise)
	assert.Contains(t, nums, 7)    // security label
	assert.Contains(t, nums, 9)    // ours
}

func TestOpenPRsEmptyWhenNoFindings(t *testing.T) {
	gh := fakePRGH{byRepo: map[string][]ghclient.PullRequest{
		"o/r": {{Number: 2, Title: "Bump x", Author: "app/dependabot", IsBot: true}},
	}}
	prs, _ := OpenPRs([]state.Repo{{Repo: "o/r"}}, gh, nil)
	assert.Empty(t, prs) // 0 findings -> no CVE-tied PRs -> no noise
}
