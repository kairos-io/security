package collect

import (
	"errors"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubCollector struct {
	name string
	out  []state.Finding
	err  error
}

func (s stubCollector) Name() string                                { return s.name }
func (s stubCollector) Collect(state.Repo) ([]state.Finding, error) { return s.out, s.err }

func TestRunIsolatesErrorsAndPreservesFirstSeen(t *testing.T) {
	nowFn = func() string { return "2026-06-19" }
	defer func() { nowFn = defaultNow }()

	repos := []state.Repo{{Repo: "kairos-io/immucore"}}
	good := stubCollector{name: "good", out: []state.Finding{
		{ID: "x", Repo: "kairos-io/immucore", Type: "sourceCVE", FirstSeen: "2026-06-19", LastSeen: "2026-06-19"},
	}}
	bad := stubCollector{name: "bad", err: errors.New("rate limited")}

	prev := state.Findings{Findings: []state.Finding{{ID: "x", FirstSeen: "2026-06-01"}}}
	got := Run(repos, []Collector{good, bad}, prev)

	require.Len(t, got.Findings, 1)
	assert.Equal(t, "2026-06-01", got.Findings[0].FirstSeen, "aging preserved")
	assert.Equal(t, "2026-06-19", got.Findings[0].LastSeen)
	require.Len(t, got.Errors, 1)
	assert.Equal(t, "bad", got.Errors[0].Collector)
	assert.Contains(t, got.Errors[0].Message, "rate limited")
}
