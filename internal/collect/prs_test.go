package collect

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPRsCollectorFiltersSecurityPRs(t *testing.T) {
	nowFn = func() string { return "2026-06-19" }
	defer func() { nowFn = defaultNow }()

	gh := ghclient.NewFake()
	gh.PRs["kairos-io/immucore"] = []ghclient.PullRequest{
		{Number: 1, Title: "Bump x/net", Author: "renovate[bot]", URL: "u1"},
		{Number: 2, Title: "Feature", Author: "alice", Labels: []string{"enhancement"}},
		{Number: 3, Title: "Patch CVE", Author: "alice", Labels: []string{"security"}},
	}
	c := PRs{GH: gh}
	fs, err := c.Collect(state.Repo{Repo: "kairos-io/immucore"})
	require.NoError(t, err)
	require.Len(t, fs, 2) // PR 1 (renovate) + PR 3 (security label)
	assert.Equal(t, "pr", fs[0].Type)
}
