package e2e

import (
	"testing"

	"github.com/kairos-io/security/internal/correlate"
	"github.com/kairos-io/security/internal/render"
	"github.com/kairos-io/security/internal/state"
	"github.com/kairos-io/security/internal/triage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// failingAI forces the deterministic fallback so the test needs no model.
type failingAI struct{}

func (failingAI) Summarize(state.Correlated) ([]string, map[string]string, string, error) {
	return nil, nil, "", assert.AnError
}

func TestPipelineCorrelateTriageRenderProducesSurfaces(t *testing.T) {
	findings := state.Findings{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high", FixedVersion: "0.33.0", FirstSeen: "2026-06-01", LastSeen: "2026-06-19"},
		{ID: "b", Repo: "kairos-io/kairos-agent", Type: "sourceCVE", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high", FixedVersion: "0.33.0", FirstSeen: "2026-06-10", LastSeen: "2026-06-19"},
	}}

	c := correlate.Run(findings)
	require.Len(t, c.Waterfall, 1, "two repos sharing a go CVE form a waterfall front")

	tr, aiErr := triage.Run(c, failingAI{}, "test-model")
	assert.Error(t, aiErr, "failingAI must surface an error")
	assert.False(t, tr.AIAvailable)
	assert.NotEmpty(t, tr.Focus)

	md := render.DashboardMarkdown(render.Input{Correlated: c, Triage: tr})
	assert.Contains(t, md, "Waterfall fronts")
	assert.Contains(t, md, "golang.org/x/net@0.33.0")

	j, err := render.DashboardJSON(render.Input{Correlated: c, Triage: tr})
	require.NoError(t, err)
	assert.Contains(t, string(j), "CVE-2025-1")
}
