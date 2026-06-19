package triage

import (
	"errors"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubAI struct {
	focus     []string
	summaries map[string]string
	narrative string
	err       error
}

func (s stubAI) Summarize(state.Correlated) ([]string, map[string]string, string, error) {
	return s.focus, s.summaries, s.narrative, s.err
}

var sampleCorrelated = state.Correlated{
	Findings: []state.Finding{
		{ID: "low1", Severity: "low"},
		{ID: "crit1", Severity: "critical"},
		{ID: "high1", Severity: "high"},
	},
}

func TestRunUsesAIWhenAvailable(t *testing.T) {
	ai := stubAI{focus: []string{"crit1"}, summaries: map[string]string{"crit1": "bad"}, narrative: "n"}
	got := Run(sampleCorrelated, ai, "m")
	assert.True(t, got.AIAvailable)
	assert.Equal(t, []string{"crit1"}, got.Focus)
	assert.Equal(t, "n", got.Narrative)
}

func TestRunFallsBackOnAIError(t *testing.T) {
	got := Run(sampleCorrelated, stubAI{err: errors.New("model down")}, "m")
	assert.False(t, got.AIAvailable)
	// deterministic severity ordering: critical, high, low
	require.Equal(t, []string{"crit1", "high1", "low1"}, got.Focus)
	assert.NotEmpty(t, got.Summaries["crit1"])
}
