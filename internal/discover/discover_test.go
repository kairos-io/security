package discover

import (
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunMergesOrgDepsAndConfig(t *testing.T) {
	gh := ghclient.NewFake()
	gh.OrgRepos["kairos-io"] = []string{"kairos-io/kairos", "kairos-io/immucore", "kairos-io/archived"}
	gh.Files["kairos-io/kairos-init|Makefile|main"] = []byte("AGENT_VERSION?=v1\n")
	gh.Files["kairos-io/kairos-init|go.mod|main"] = []byte("require github.com/mudler/yip v1.0.0\n")

	cfg := config.ReposConfig{
		Repos:   []state.Repo{{Repo: "mudler/edgevpn", Kind: "external", Branch: "master", Criticality: "high"}},
		Exclude: []string{"kairos-io/archived"},
	}

	repos, err := Run(gh, cfg, "kairos-io", "kairos-io/kairos-init", "main")
	require.NoError(t, err)

	names := map[string]state.Repo{}
	for _, r := range repos {
		names[r.Repo] = r
	}
	assert.Contains(t, names, "kairos-io/kairos")
	assert.Contains(t, names, "kairos-io/immucore")
	assert.Contains(t, names, "kairos-io/kairos-agent") // from Makefile
	assert.Contains(t, names, "mudler/yip")             // from go.mod
	assert.Contains(t, names, "mudler/edgevpn")         // from config
	assert.NotContains(t, names, "kairos-io/archived")  // excluded
	assert.Equal(t, "high", names["mudler/edgevpn"].Criticality)
	assert.Equal(t, "org", names["kairos-io/kairos"].Kind)
	// sorted output
	assert.True(t, repos[0].Repo < repos[len(repos)-1].Repo)
}
