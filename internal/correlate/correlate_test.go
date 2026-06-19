package correlate

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCorrelateDedupesAndBuildsWaterfall(t *testing.T) {
	in := state.Findings{Findings: []state.Finding{
		// same CVE in immucore seen by two sources → dedupe to 1, severity high wins
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "unknown", FixedVersion: "0.33.0"},
		{ID: "b", Repo: "kairos-io/immucore", Type: "ghAlert", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high"},
		// same CVE/package in a second repo → waterfall group of 2 repos
		{ID: "c", Repo: "kairos-io/kairos-agent", Type: "sourceCVE", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high", FixedVersion: "0.33.0"},
	}}

	out := Run(in)

	// dedupe: immucore CVE-2025-1 collapses to one finding, severity "high"
	count := 0
	for _, f := range out.Findings {
		if f.Repo == "kairos-io/immucore" && f.CVEID == "CVE-2025-1" {
			count++
			assert.Equal(t, "high", f.Severity)
			assert.Equal(t, "0.33.0", f.FixedVersion)
		}
	}
	assert.Equal(t, 1, count)

	require.Len(t, out.Waterfall, 1)
	g := out.Waterfall[0]
	assert.ElementsMatch(t, []string{"kairos-io/immucore", "kairos-io/kairos-agent"}, g.AffectedRepos)
	assert.Equal(t, "golang.org/x/net", g.SuggestedBump.Package)
	assert.Equal(t, "0.33.0", g.SuggestedBump.To)
	assert.Equal(t, "high", g.Severity)
}
