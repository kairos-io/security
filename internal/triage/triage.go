package triage

import (
	"fmt"
	"sort"
	"time"

	"github.com/kairos-io/security/internal/state"
)

type AIClient interface {
	Summarize(c state.Correlated) (focus []string, summaries map[string]string, narrative string, err error)
}

// FocusLimit caps the deterministic-fallback focus list so it stays a short,
// human-readable shortlist rather than a dump of every finding.
const FocusLimit = 20

var sevRank = map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1, "unknown": 0, "": 0}

var nowFn = func() string { return time.Now().UTC().Format("2006-01-02") }

// Run produces the triage. When an AI client is supplied and succeeds, the
// returned Triage carries the model's output and AIAvailable is true. When the
// AI client fails, Run still returns a fully-populated deterministic Triage
// (so callers can render something) AND a non-nil error describing the AI
// failure — the caller decides whether to fail hard or fall back. A nil AI
// client is the deterministic path and is not an error.
func Run(c state.Correlated, ai AIClient, model string) (state.Triage, error) {
	t := state.Triage{GeneratedAt: nowFn(), Model: model, Summaries: map[string]string{}}
	// A clean scan (no findings, no waterfall groups) is a valid outcome, not an
	// AI failure: short-circuit before any AI call so --require-ai does not fail
	// on empty input, and so a small model is never asked to triage an empty
	// list (which it answers with prose that breaks the forced tool call).
	if len(c.Findings) == 0 && len(c.Waterfall) == 0 {
		t.AIAvailable = true
		t.Narrative = "No security findings to triage this run."
		return t, nil
	}
	if ai != nil {
		focus, summaries, narrative, err := ai.Summarize(c)
		if err == nil {
			t.AIAvailable = true
			t.Focus = focus
			t.Summaries = summaries
			t.Narrative = narrative
			return t, nil
		}
		applyFallback(&t, c)
		return t, fmt.Errorf("AI summarize failed: %w", err)
	}
	applyFallback(&t, c)
	return t, nil
}

func applyFallback(t *state.Triage, c state.Correlated) {
	t.AIAvailable = false
	t.Focus = deterministicFocus(c)
	t.Summaries = templatedSummaries(c)
	t.Narrative = fmt.Sprintf("AI unavailable this run. %d findings, %d waterfall groups, ordered by severity.",
		len(c.Findings), len(c.Waterfall))
}

func deterministicFocus(c state.Correlated) []string {
	type item struct {
		id  string
		sev string
	}
	var items []item
	for _, f := range c.Findings {
		items = append(items, item{f.ID, f.Severity})
	}
	for _, g := range c.Waterfall {
		items = append(items, item{g.ID, g.Severity})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if sevRank[items[i].sev] != sevRank[items[j].sev] {
			return sevRank[items[i].sev] > sevRank[items[j].sev]
		}
		return items[i].id < items[j].id
	})
	out := make([]string, 0, len(items))
	for _, it := range items {
		out = append(out, it.id)
	}
	if len(out) > FocusLimit {
		out = out[:FocusLimit]
	}
	return out
}

func templatedSummaries(c state.Correlated) map[string]string {
	out := map[string]string{}
	for _, f := range c.Findings {
		// Summarize every finding so each focus id resolves to readable text.
		out[f.ID] = fmt.Sprintf("%s %s in %s (%s)", f.Severity, f.CVEID, f.Repo, f.Package)
	}
	for _, g := range c.Waterfall {
		out[g.ID] = fmt.Sprintf("%s affects %d repos via %s", g.Severity, len(g.AffectedRepos), g.RootCause)
	}
	return out
}
