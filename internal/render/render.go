package render

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

type Input struct {
	Correlated    state.Correlated        `json:"correlated"`
	Triage        state.Triage            `json:"triage"`
	CollectErrors []state.CollectionError `json:"collectErrors"`
	RunURL        string                  `json:"runURL"`
}

func DashboardJSON(in Input) ([]byte, error) {
	b, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

func DashboardMarkdown(in Input) string {
	var b strings.Builder
	b.WriteString("# Kairos Security Dashboard\n\n")
	fmt.Fprintf(&b, "_Updated %s", in.Triage.GeneratedAt)
	if !in.Triage.AIAvailable {
		b.WriteString(" — ⚠️ AI unavailable this run")
	}
	b.WriteString("._\n\n")

	if in.Triage.Narrative != "" {
		b.WriteString("> " + in.Triage.Narrative + "\n\n")
	}

	// Focus now
	b.WriteString("## 🔥 Focus now\n\n")
	if len(in.Triage.Focus) == 0 {
		b.WriteString("_Nothing flagged._\n\n")
	} else {
		for _, id := range in.Triage.Focus {
			if s, ok := in.Triage.Summaries[id]; ok {
				fmt.Fprintf(&b, "- **%s** — %s\n", id, s)
			} else {
				fmt.Fprintf(&b, "- **%s**\n", id)
			}
		}
		b.WriteString("\n")
	}

	// Waterfall fronts
	b.WriteString("## 🌊 Waterfall fronts\n\n")
	if len(in.Correlated.Waterfall) == 0 {
		b.WriteString("_None._\n\n")
	} else {
		b.WriteString("| Root cause | Severity | Bump | Affected repos |\n|---|---|---|---|\n")
		for _, g := range in.Correlated.Waterfall {
			fmt.Fprintf(&b, "| %s | %s | %s@%s | %s |\n",
				g.RootCause, g.Severity, g.SuggestedBump.Package, g.SuggestedBump.To,
				strings.Join(g.AffectedRepos, ", "))
		}
		b.WriteString("\n")
	}

	// Per-repo table
	b.WriteString("## 📦 Per-repo findings\n\n")
	b.WriteString("| Repo | Critical | High | Medium | Low | Total |\n|---|---|---|---|---|---|\n")
	for _, row := range perRepoRows(in.Correlated.Findings) {
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %d | %d |\n",
			row.repo, row.crit, row.high, row.med, row.low, row.total)
	}
	b.WriteString("\n")

	// Collection errors
	if len(in.CollectErrors) > 0 {
		fmt.Fprintf(&b, "## ⚠️ %d collection errors\n\n", len(in.CollectErrors))
		for _, e := range in.CollectErrors {
			fmt.Fprintf(&b, "- `%s` / %s: %s\n", e.Repo, e.Collector, e.Message)
		}
		b.WriteString("\n")
	}

	if in.RunURL != "" {
		fmt.Fprintf(&b, "---\n[Run log](%s)\n", in.RunURL)
	}
	return b.String()
}

type repoRow struct {
	repo                        string
	crit, high, med, low, total int
}

func perRepoRows(findings []state.Finding) []repoRow {
	idx := map[string]*repoRow{}
	for _, f := range findings {
		r := idx[f.Repo]
		if r == nil {
			r = &repoRow{repo: f.Repo}
			idx[f.Repo] = r
		}
		r.total++
		switch f.Severity {
		case "critical":
			r.crit++
		case "high":
			r.high++
		case "medium":
			r.med++
		case "low":
			r.low++
		}
	}
	rows := make([]repoRow, 0, len(idx))
	for _, r := range idx {
		rows = append(rows, *r)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].crit != rows[j].crit {
			return rows[i].crit > rows[j].crit
		}
		if rows[i].high != rows[j].high {
			return rows[i].high > rows[j].high
		}
		return rows[i].repo < rows[j].repo
	})
	return rows
}
