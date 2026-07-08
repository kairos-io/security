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
				// Informational: separated + uncounted. Both are severity-bearing
				// on a tracked repo, so the golden proves they never leak into the
				// headline counts or the per-repo / hadron actionable tables.
				{ID: "info-fixed", Repo: "kairos-io/kairos", Type: "componentCVE", CVEID: "CVE-2024-0001", Package: "libxml2", CurrentVersion: "2.13.0", FixedVersion: "2.12.0", Severity: "critical", Class: "informational", ClassReason: "already fixed: current 2.13.0 is past fixed 2.12.0"},
				{ID: "info-accepted", Repo: "kairos-io/kairos", Type: "componentCVE", CVEID: "CVE-2023-5678", Package: "openssl-fips", CurrentVersion: "3.0.8", Severity: "medium", Class: "informational", ClassReason: "accepted: pinned FIPS build"},
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
// repos beyond those with findings: a clean repo shows "clean (no crit/high/med)"
// and an errored repo shows the "⚠️ errors" status.
func TestDashboardMarkdownRepoStatus(t *testing.T) {
	got := DashboardMarkdown(sampleInput())
	assert.Contains(t, got, "| [kairos-io/kairos-sdk](https://github.com/kairos-io/kairos-sdk) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |")
	assert.Contains(t, got, "| [kairos-io/x](https://github.com/kairos-io/x) | 0 | 0 | 0 | 0 | ⚠️ errors |")
	assert.Contains(t, got, "| [kairos-io/kairos](https://github.com/kairos-io/kairos) | 1 | 0 | 0 | 1 | ok |")
}

func TestDashboardMarkdownCoordinationSummary(t *testing.T) {
	got := DashboardMarkdown(Input{CoordinationSummary: "X cascading"})
	assert.Contains(t, got, "Coordination")
	assert.Contains(t, got, "X cascading")
}

func TestFocusShowsTitleLinkNotID(t *testing.T) {
	in := Input{
		Correlated: state.Correlated{Findings: []state.Finding{
			{ID: "abc123", Repo: "o/r", Title: "x/net rapid reset", URL: "https://github.com/o/r/pull/9"},
			{ID: "def456", Repo: "o/r", Type: "sourceCVE", CVEID: "GO-2024-3218", Title: "kad-dht issue"},
		}},
		Triage: state.Triage{Focus: []string{"abc123", "def456"}},
	}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "[x/net rapid reset](https://github.com/o/r/pull/9)")
	assert.Contains(t, md, "[kad-dht issue](https://pkg.go.dev/vuln/GO-2024-3218)")
	assert.NotContains(t, md, "abc123") // raw id never shown
}

// TestFocusSourceCVEWithCVEIDNoURL guards the synthesis fallback: a sourceCVE
// finding whose only id is a CVE alias (not GO-prefixed) and which carries no
// URL must render as a bare title, never a pkg.go.dev/vuln/CVE-… link (which
// 400s, since pkg.go.dev/vuln only serves GO-… paths).
func TestFocusSourceCVEWithCVEIDNoURL(t *testing.T) {
	in := Input{
		Correlated: state.Correlated{Findings: []state.Finding{
			{ID: "ghi789", Repo: "o/r", Type: "sourceCVE", CVEID: "CVE-2023-39325", Title: "http2 rapid reset"},
		}},
		Triage: state.Triage{Focus: []string{"ghi789"}},
	}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "- http2 rapid reset") // bare title, no link
	assert.NotContains(t, md, "pkg.go.dev/vuln/CVE-")
}

// TestFindingLinkFallsBackToCVEIDNotPackage guards the title-fallback chain: a
// finding with no Title but a CVEID must render its CVE id as the link text, not
// the bare package name — otherwise many distinct CVEs on the same package all
// render identically (the "looks duplicated" bug).
func TestFindingLinkFallsBackToCVEIDNotPackage(t *testing.T) {
	in := Input{Correlated: state.Correlated{Findings: []state.Finding{
		{ID: "h1", Repo: "kairos-io/hadron", Type: "componentCVE", Package: "expat", CVEID: "CVE-2026-56408", CurrentVersion: "2.8.1", FixedVersion: "2.8.2", Severity: "high"},
	}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "CVE-2026-56408")
	assert.NotContains(t, md, "| expat | 2.8.1 | 2.8.2 | high | expat |") // NOT the bare package name as link text
}

// TestFindingLinkFallsBackToGHSAWhenNoCVEID: with Title and CVEID both empty but
// a GHSA present, the GHSA id is the link text (still better than bare package).
func TestFindingLinkFallsBackToGHSAWhenNoCVEID(t *testing.T) {
	in := Input{Correlated: state.Correlated{Findings: []state.Finding{
		{ID: "h1", Repo: "kairos-io/hadron", Type: "componentCVE", Package: "expat", GHSA: "GHSA-abcd-1234-wxyz", CurrentVersion: "2.8.1", FixedVersion: "2.8.2", Severity: "high"},
	}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "GHSA-abcd-1234-wxyz")
}

func TestOpenPRsSection(t *testing.T) {
	in := Input{OpenPRs: []state.TrackedPR{
		{Repo: "o/r", Number: 7, Title: "bump foo", URL: "https://github.com/o/r/pull/7", Source: "dependabot"},
	}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "📋 Open PRs")
	assert.Contains(t, md, "[#7 bump foo](https://github.com/o/r/pull/7)")
	assert.Contains(t, md, "dependabot")
}

func TestOpenPRShowsSupersededStatus(t *testing.T) {
	in := Input{
		OpenPRs: []state.TrackedPR{{Repo: "o/r", Number: 38, Title: "bump x", URL: "https://github.com/o/r/pull/38", Source: "bot"}},
		Ledger: state.Ledger{Entries: []state.LedgerEntry{
			{Key: "o/r|x", Repo: "o/r", Source: "ksec", State: "open", PRNumber: 77,
				PRURL: "https://github.com/o/r/pull/77", Supersedes: "https://github.com/o/r/pull/38"},
		}},
	}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "superseded by") // status/action surfaced
	assert.Contains(t, md, "/pull/77")      // links the ksec PR
}

func TestNeedsHumanRollup(t *testing.T) {
	in := Input{Ledger: state.Ledger{Entries: []state.LedgerEntry{
		{Key: "o/r|x", Repo: "o/r", State: "build-failed", NeedsHuman: true, Bump: state.Bump{Package: "x", To: "1.0"}},
	}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "🚑 Needs human")
	assert.Contains(t, md, "o/r")
}

func TestDashboardShowsBotPRReviews(t *testing.T) {
	md := DashboardMarkdown(Input{Reviews: []state.PRReview{
		{Repo: "kairos-io/AuroraBoot", PR: 566, URL: "https://github.com/kairos-io/AuroraBoot/pull/566", Verdict: "good", Reasoning: "clean go.mod bump", ChangesSummary: "bumps golang.org/x/net v0.30.0 → v0.31.0; bugfixes only", Trace: []string{"foo/bar 1.0→1.1: compare v1.0.0...v1.1.0 ✓ 1234 bytes"}},
		{Repo: "kairos-io/AuroraBoot", PR: 567, URL: "u567", Verdict: "needs_human_verification", Reasoning: "touches source"},
	}})
	assert.Contains(t, md, "🔎 Bot-PR reviews")
	assert.Contains(t, md, "[#566")
	assert.Contains(t, md, "good")
	assert.Contains(t, md, "needs_human_verification")
	assert.Contains(t, md, "  ↳ bumps golang.org/x/net v0.30.0 → v0.31.0; bugfixes only")
	assert.Contains(t, md, "    - foo/bar 1.0→1.1: compare v1.0.0...v1.1.0 ✓ 1234 bytes")
}

func TestDashboardJSONIsStable(t *testing.T) {
	a, err := DashboardJSON(sampleInput())
	require.NoError(t, err)
	b, err := DashboardJSON(sampleInput())
	require.NoError(t, err)
	assert.Equal(t, string(a), string(b))
}

func TestDashboardMarksSkippedRepo(t *testing.T) {
	f := false
	in := Input{Repos: []state.Repo{{Repo: "o/r", Scan: state.ScanConfig{Source: &f}}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "skipped: not source-scannable")
}

func TestHadronComponentCVESection(t *testing.T) {
	in := Input{Correlated: state.Correlated{Findings: []state.Finding{
		{ID: "h1", Repo: "kairos-io/hadron", Type: "componentCVE", Package: "openssl", CVEID: "CVE-2025-1", CurrentVersion: "3.6.3", FixedVersion: "3.6.4", Severity: "high", Source: "osv", URL: "https://osv.dev/vulnerability/CVE-2025-1"},
		{ID: "h2", Repo: "kairos-io/hadron", Type: "componentCVE", Package: "busybox", CVEID: "CVE-2025-2", CurrentVersion: "1.37.0", Severity: "medium", Source: "nvd"},
	}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "🧩 Hadron component CVEs")
	assert.Contains(t, md, "| openssl | 3.6.3 | 3.6.4 | high |")
	assert.Contains(t, md, "| busybox | 1.37.0 | — | medium |") // no fixed version yet -> em dash, not blank
}

func TestHadronComponentCVESectionOmittedWhenEmpty(t *testing.T) {
	md := DashboardMarkdown(Input{})
	assert.NotContains(t, md, "Hadron component CVEs")
}
