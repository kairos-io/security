package render

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestComputeActivityNoFindings(t *testing.T) {
	f := false
	a := computeActivity(Input{
		Repos: []state.Repo{{Repo: "o/a"}, {Repo: "o/b", Scan: state.ScanConfig{Source: &f}}},
	})
	assert.Equal(t, 2, a.Repos)
	assert.Equal(t, 1, a.Skipped)
	assert.Equal(t, 0, a.Findings)
	assert.Contains(t, a.Why, "No CVEs")
}

func TestComputeActivityWithFindingsAndPRs(t *testing.T) {
	a := computeActivity(Input{
		Repos:         []state.Repo{{Repo: "o/a"}},
		Correlated:    state.Correlated{Findings: []state.Finding{{Repo: "o/a", Severity: "high"}, {Repo: "o/a", Severity: "low"}, {Repo: "o/a", Severity: "critical", Class: "informational"}}},
		OpenPRs:       []state.TrackedPR{{Repo: "o/a", Source: "dependabot"}},
		Ledger:        state.Ledger{Entries: []state.LedgerEntry{{State: "open"}, {State: "build-failed", NeedsHuman: true}}},
		CollectErrors: []state.CollectionError{{Repo: "o/a", Collector: "sourceCVE", Message: "boom"}},
	})
	assert.Equal(t, 2, a.Findings)      // informational excluded from the count
	assert.Equal(t, 1, a.Informational) // …counted separately
	assert.Equal(t, 0, a.Crit)          // informational critical does not bump Crit
	assert.Equal(t, 1, a.High)
	assert.Equal(t, 1, a.PRs)
	assert.Equal(t, 1, a.PRsBySource["dependabot"])
	assert.Equal(t, 1, a.LedgerOpen)
	assert.Equal(t, 1, a.NeedsHuman)
	assert.Equal(t, 1, a.Errored)
	assert.Contains(t, a.Why, "need a human")
}

func TestDashboardShowsThisRun(t *testing.T) {
	md := DashboardMarkdown(Input{Repos: []state.Repo{{Repo: "o/a"}}})
	assert.Contains(t, md, "📋 This run")
	assert.Contains(t, md, "No CVEs")
}
