package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildGraph(t *testing.T) {
	repos := []state.Repo{
		{Repo: "kairos-io/kairos-sdk", Branch: "main"},
		{Repo: "kairos-io/immucore", Branch: "master"},
		{Repo: "kairos-io/kairos-agent"},
	}
	gomod := map[string][]byte{
		"kairos-io/kairos-sdk":   []byte("module github.com/kairos-io/kairos-sdk\ngo 1.22\n"),
		"kairos-io/immucore":     []byte("module github.com/kairos-io/immucore\nrequire github.com/kairos-io/kairos-sdk v0.7.0\n"),
		"kairos-io/kairos-agent": []byte("module github.com/kairos-io/kairos-agent\nrequire (\n\tgithub.com/kairos-io/kairos-sdk v0.7.0\n\tgithub.com/kairos-io/immucore v0.5.0\n)\n"),
	}
	g := BuildGraph(repos, gomod)

	assert.Equal(t, "github.com/kairos-io/kairos-sdk", g.ModuleOf("kairos-io/kairos-sdk"))
	r, ok := g.RepoOf("github.com/kairos-io/kairos-sdk")
	require.True(t, ok)
	assert.Equal(t, "kairos-io/kairos-sdk", r)

	// consumers of the sdk module: immucore and kairos-agent, sorted
	assert.Equal(t, []string{"kairos-io/immucore", "kairos-io/kairos-agent"}, g.Consumers("github.com/kairos-io/kairos-sdk"))
	// consumers of immucore: kairos-agent
	assert.Equal(t, []string{"kairos-io/kairos-agent"}, g.Consumers("github.com/kairos-io/immucore"))
	// branch lookups
	assert.Equal(t, "master", g.BranchOf("kairos-io/immucore"))
	assert.Equal(t, "main", g.BranchOf("kairos-io/kairos-agent")) // default
}
