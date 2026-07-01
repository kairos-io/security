package render

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDashboardHTMLGolden(t *testing.T) {
	got := DashboardHTML(sampleInput())
	golden := filepath.Join("testdata", "dashboard.html.golden")
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		require.NoError(t, os.WriteFile(golden, []byte(got), 0o644))
	}
	want, err := os.ReadFile(golden)
	require.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestDashboardHTMLIsStable(t *testing.T) {
	a := DashboardHTML(sampleInput())
	b := DashboardHTML(sampleInput())
	assert.Equal(t, a, b)
}

// TestDashboardHTMLEscaping confirms that dynamic text containing markup is
// HTML-escaped and cannot break the page or inject elements.
func TestDashboardHTMLEscaping(t *testing.T) {
	in := Input{
		Correlated: state.Correlated{
			Findings: []state.Finding{
				{ID: "x", Repo: "kairos-io/x<script>alert(1)</script>", Severity: "critical"},
			},
			Waterfall: []state.WaterfallGroup{
				{RootCause: "evil <b>root</b> & co", Severity: "high",
					AffectedRepos: []string{"kairos-io/<i>repo</i>"},
					SuggestedBump: state.Bump{Package: "pkg<x>", To: "1<2"}},
			},
		},
		Triage: state.Triage{
			GeneratedAt: "2026-06-19", AIAvailable: true,
			Focus:     []string{"x"},
			Summaries: map[string]string{"x": "summary with <em>tag</em> & amp"},
			Narrative: "narrative <script>bad</script>",
		},
		CollectErrors: []state.CollectionError{
			{Repo: "kairos-io/x", Collector: "prs", Message: "broke <b>&</b> | pipe"},
		},
		RunURL: "https://example.com/run?a=1&b=2",
	}
	got := DashboardHTML(in)

	// Raw, unescaped markup must never appear in the output.
	for _, raw := range []string{
		"<script>alert(1)</script>",
		"<b>root</b>",
		"<i>repo</i>",
		"<em>tag</em>",
		"<script>bad</script>",
		"<b>&</b>",
	} {
		assert.NotContains(t, got, raw, "raw markup leaked into HTML output")
	}
	// And it must be escaped instead.
	assert.Contains(t, got, "&lt;script&gt;alert(1)&lt;/script&gt;")
}

// TestDashboardHTMLSections sanity-checks the major sections are present.
func TestDashboardHTMLSections(t *testing.T) {
	got := DashboardHTML(sampleInput())
	for _, want := range []string{
		"<!DOCTYPE html>",
		"<title>Kairos Security Dashboard</title>",
		`<span class="brand">Kairos `,
		"Needs attention",
		"Focus now",
		"Waterfall fronts",
		"Per-repo findings",
		"collection error",
		"actions/runs/1",
		// repo names are links, and every link opens in a new tab
		`href="https://github.com/kairos-io/kairos"`,
		`target="_blank" rel="noopener"`,
	} {
		assert.True(t, strings.Contains(got, want), "expected output to contain %q", want)
	}
}

func TestDashboardHTMLHadronComponentSection(t *testing.T) {
	in := Input{Correlated: state.Correlated{Findings: []state.Finding{
		{ID: "h1", Repo: "kairos-io/hadron", Type: "componentCVE", Package: "openssl", CVEID: "CVE-2025-1", CurrentVersion: "3.6.3", FixedVersion: "3.6.4", Severity: "high"},
	}}}
	got := DashboardHTML(in)
	assert.Contains(t, got, "Hadron component CVEs")
	assert.Contains(t, got, "<td>openssl</td>")
	assert.Contains(t, got, "<td>3.6.4</td>")
}
