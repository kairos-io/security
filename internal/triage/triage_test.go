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

	got := Run(c, stubAI{err: errors.New("model down")}, "m")
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
