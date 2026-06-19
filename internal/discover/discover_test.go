package discover

import (
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeAppliesDefaultsExcludeAndSort(t *testing.T) {
	cfg := config.ReposConfig{
		Repos: []state.Repo{
			{Repo: "kairos-io/kairos"}, // kairos-io: kind inferred "org", defaults applied
			{Repo: "mudler/edgevpn", Kind: "external", Branch: "master", Criticality: "high"}, // explicit fields kept
			{Repo: "kairos-io/archived"}, // dropped via Exclude
		},
		Exclude: []string{"kairos-io/archived"},
	}

	repos := Normalize(cfg)

	names := map[string]state.Repo{}
	for _, r := range repos {
		names[r.Repo] = r
	}

	// excluded repo dropped
	assert.NotContains(t, names, "kairos-io/archived")
	require.Len(t, repos, 2)

	// defaults applied to the kairos-io repo
	k := names["kairos-io/kairos"]
	assert.Equal(t, "org", k.Kind)
	assert.Equal(t, "main", k.Branch)
	assert.Equal(t, "medium", k.Criticality)

	// explicit fields preserved
	e := names["mudler/edgevpn"]
	assert.Equal(t, "external", e.Kind)
	assert.Equal(t, "master", e.Branch)
	assert.Equal(t, "high", e.Criticality)

	// sorted by repo
	assert.Equal(t, "kairos-io/kairos", repos[0].Repo)
	assert.Equal(t, "mudler/edgevpn", repos[1].Repo)
}

func TestSeedFromOrgEnumeratesOrgAndDeps(t *testing.T) {
	gh := ghclient.NewFake()
	gh.OrgRepos["kairos-io"] = []string{"kairos-io/kairos", "kairos-io/immucore"}
	gh.Files["kairos-io/kairos-init|Makefile|main"] = []byte("AGENT_VERSION?=v1\n")
	gh.Files["kairos-io/kairos-init|go.mod|main"] = []byte("require github.com/mudler/yip v1.0.0\n")

	repos, err := SeedFromOrg(gh, "kairos-io", "kairos-io/kairos-init", "main")
	require.NoError(t, err)

	names := map[string]state.Repo{}
	for _, r := range repos {
		names[r.Repo] = r
	}
	assert.Contains(t, names, "kairos-io/kairos")
	assert.Contains(t, names, "kairos-io/immucore")
	assert.Contains(t, names, "kairos-io/kairos-agent") // from Makefile
	assert.Contains(t, names, "mudler/yip")             // from go.mod
	assert.Equal(t, "org", names["kairos-io/kairos"].Kind)
	// sorted output
	assert.True(t, repos[0].Repo < repos[len(repos)-1].Repo)
}
