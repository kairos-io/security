package triage

import (
	"errors"
	"fmt"
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
	got, err := Run(sampleCorrelated, ai, "m")
	assert.NoError(t, err)
	assert.True(t, got.AIAvailable)
	assert.Equal(t, []string{"crit1"}, got.Focus)
	assert.Equal(t, "n", got.Narrative)
}

func TestRunEmptyFindingsSkipsAI(t *testing.T) {
	// A clean scan must not call the AI, so --require-ai cannot fail on empty
	// input. stubAI returns an error; if Run called it, Run would surface that.
	got, err := Run(state.Correlated{}, stubAI{err: errors.New("must not be called")}, "m")
	require.NoError(t, err)
	assert.True(t, got.AIAvailable)
	assert.Empty(t, got.Focus)
	assert.Contains(t, got.Narrative, "No security findings")
}

func TestRunFallsBackOnAIError(t *testing.T) {
	got, err := Run(sampleCorrelated, stubAI{err: errors.New("model down")}, "m")
	assert.Error(t, err)
	assert.False(t, got.AIAvailable)
	// deterministic severity ordering: critical, high, low
	require.Equal(t, []string{"crit1", "high1", "low1"}, got.Focus)
	assert.NotEmpty(t, got.Summaries["crit1"])
}

func TestRunFallbackExcludesInformationalFindings(t *testing.T) {
	// Informational findings (accepted/already-fixed) are separated and must
	// never reach the "focus now" shortlist, even in the deterministic path.
	c := state.Correlated{Findings: []state.Finding{
		{ID: "act1", Severity: "high"},
		{ID: "info1", Severity: "critical", Class: "informational"},
	}}

	got, err := Run(c, stubAI{err: errors.New("model down")}, "m")
	assert.Error(t, err)
	assert.False(t, got.AIAvailable)
	assert.Equal(t, []string{"act1"}, got.Focus, "informational finding must not appear in focus")
	assert.NotContains(t, got.Focus, "info1")
}

func TestRunFallbackCapsFocusAndSummarizesEveryFinding(t *testing.T) {
	// A few high-severity findings plus many low ones: with FocusLimit==20,
	// the high ones sort first and the remaining slots are filled by lows, so
	// at least one low-severity finding lands in Focus.
	var findings []state.Finding
	mkSev := func(i int) string {
		if i < 3 {
			return "critical"
		}
		return "low"
	}
	for i := 0; i < 40; i++ {
		findings = append(findings, state.Finding{
			ID:       fmt.Sprintf("f%02d", i),
			Severity: mkSev(i),
			CVEID:    fmt.Sprintf("CVE-2024-%04d", i),
			Repo:     "kairos-io/kairos",
			Package:  "pkg",
		})
	}
	c := state.Correlated{Findings: findings}

	got, err := Run(c, stubAI{err: errors.New("model down")}, "m")
	assert.Error(t, err)
	assert.False(t, got.AIAvailable)
	require.Len(t, got.Focus, FocusLimit, "focus must be capped at FocusLimit")

	// Every focus id resolves to a non-empty summary.
	for _, id := range got.Focus {
		assert.NotEmpty(t, got.Summaries[id], "every focus id must resolve to a summary")
	}

	// A low-severity finding present in Focus must have a non-empty summary.
	sevByID := map[string]string{}
	for _, f := range findings {
		sevByID[f.ID] = f.Severity
	}
	var lowInFocus string
	for _, id := range got.Focus {
		if sevByID[id] == "low" {
			lowInFocus = id
			break
		}
	}
	require.NotEmpty(t, lowInFocus, "expected a low-severity finding in Focus")
	assert.NotEmpty(t, got.Summaries[lowInFocus])
}
