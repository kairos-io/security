package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleInput() Input {
	return Input{
		Correlated: state.Correlated{
			Findings: []state.Finding{
				{ID: "crit1", Repo: "kairos-io/kairos", Type: "imageCVE", CVEID: "CVE-2025-9999", Package: "openssl", Severity: "critical", FirstSeen: "2026-06-01", LastSeen: "2026-06-19"},
			},
			Waterfall: []state.WaterfallGroup{
				{ID: "go-CVE-2025-1-golang.org/x/net", RootCause: "golang.org/x/net (CVE-2025-1)", Severity: "high", AffectedRepos: []string{"kairos-io/immucore", "kairos-io/kairos-agent"}, SuggestedBump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
			},
		},
		Triage: state.Triage{
			GeneratedAt: "2026-06-19", AIAvailable: true,
			Focus:     []string{"crit1", "go-CVE-2025-1-golang.org/x/net"},
			Summaries: map[string]string{"crit1": "Critical openssl CVE in kairos image"},
			Narrative: "Focus on the openssl critical first.",
		},
		Repos: []state.Repo{
			{Repo: "kairos-io/kairos"},     // has the critical finding -> "ok"
			{Repo: "kairos-io/kairos-sdk"}, // no findings, no error -> "clean"
			{Repo: "kairos-io/x"},          // appears in CollectErrors -> "⚠️ errors"
		},
		CollectErrors: []state.CollectionError{{Repo: "kairos-io/x", Collector: "prs", Message: "rate limited"}},
		Ledger: state.Ledger{Entries: []state.LedgerEntry{
			{Key: "kairos-io/immucore|golang.org/x/net", Repo: "kairos-io/immucore", Package: "golang.org/x/net", State: "open", PRNumber: 412, PRURL: "https://github.com/kairos-io/immucore/pull/412", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
		}},
		RunURL: "https://github.com/kairos-io/security/actions/runs/1",
	}
}

func TestDashboardMarkdownGolden(t *testing.T) {
	got := DashboardMarkdown(sampleInput())
	golden := filepath.Join("testdata", "dashboard.md.golden")
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		require.NoError(t, os.WriteFile(golden, []byte(got), 0o644))
	}
	want, err := os.ReadFile(golden)
	require.NoError(t, err)
	assert.Equal(t, string(want), got)
}

// TestDashboardMarkdownRepoStatus confirms the per-repo table lists tracked
// repos beyond those with findings: a clean repo shows "clean" and an errored
// repo shows the "⚠️ errors" status.
func TestDashboardMarkdownRepoStatus(t *testing.T) {
	got := DashboardMarkdown(sampleInput())
	assert.Contains(t, got, "| kairos-io/kairos-sdk | 0 | 0 | 0 | 0 | 0 | clean |")
	assert.Contains(t, got, "| kairos-io/x | 0 | 0 | 0 | 0 | 0 | ⚠️ errors |")
	assert.Contains(t, got, "| kairos-io/kairos | 1 | 0 | 0 | 0 | 1 | ok |")
}

func TestDashboardMarkdownCoordinationSummary(t *testing.T) {
	got := DashboardMarkdown(Input{CoordinationSummary: "X cascading"})
	assert.Contains(t, got, "Coordination")
	assert.Contains(t, got, "X cascading")
}

func TestDashboardJSONIsStable(t *testing.T) {
	a, err := DashboardJSON(sampleInput())
	require.NoError(t, err)
	b, err := DashboardJSON(sampleInput())
	require.NoError(t, err)
	assert.Equal(t, string(a), string(b))
}
