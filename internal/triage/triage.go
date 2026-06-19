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

var sevRank = map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1, "unknown": 0, "": 0}

var nowFn = func() string { return time.Now().UTC().Format("2006-01-02") }

func Run(c state.Correlated, ai AIClient, model string) state.Triage {
	t := state.Triage{GeneratedAt: nowFn(), Model: model, Summaries: map[string]string{}}
	if ai != nil {
		focus, summaries, narrative, err := ai.Summarize(c)
		if err == nil {
			t.AIAvailable = true
			t.Focus = focus
			t.Summaries = summaries
			t.Narrative = narrative
			return t
		}
	}
	t.AIAvailable = false
	t.Focus = deterministicFocus(c)
	t.Summaries = templatedSummaries(c)
	t.Narrative = fmt.Sprintf("AI unavailable this run. %d findings, %d waterfall groups, ordered by severity.",
		len(c.Findings), len(c.Waterfall))
	return t
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
	return out
}

func templatedSummaries(c state.Correlated) map[string]string {
	out := map[string]string{}
	for _, f := range c.Findings {
		if sevRank[f.Severity] >= sevRank["high"] {
			out[f.ID] = fmt.Sprintf("%s %s in %s (%s)", f.Severity, f.CVEID, f.Repo, f.Package)
		}
	}
	for _, g := range c.Waterfall {
		out[g.ID] = fmt.Sprintf("%s affects %d repos via %s", g.Severity, len(g.AffectedRepos), g.RootCause)
	}
	return out
}
