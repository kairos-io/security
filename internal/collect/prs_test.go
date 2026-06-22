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

func TestOpenPRsTracksAndClassifies(t *testing.T) {
	gh := fakePRGH{byRepo: map[string][]ghclient.PullRequest{
		"o/r": {
			{Number: 2, Title: "bump y", Author: "dependabot[bot]", URL: "u2"},
			{Number: 1, Title: "feature", Author: "alice"}, // not tracked (no bot, no label)
			{Number: 3, Title: "sec fix", Author: "bob", Labels: []string{"security"}, URL: "u3"},
		},
	}}
	prs, errs := OpenPRs([]state.Repo{{Repo: "o/r"}}, gh)
	require.Empty(t, errs)
	require.Len(t, prs, 2)
	// sorted by repo then number
	assert.Equal(t, 2, prs[0].Number)
	assert.Equal(t, "dependabot", prs[0].Source)
	assert.Equal(t, "u2", prs[0].URL)
	assert.Equal(t, 3, prs[1].Number)
	assert.Equal(t, "human", prs[1].Source)
}
