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
	Repos         []state.Repo            `json:"repos"`
	Ledger        state.Ledger            `json:"ledger"`
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
	b.WriteString("| Repo | Critical | High | Medium | Low | Total | Status |\n|---|---|---|---|---|---|---|\n")
	for _, row := range perRepoRows(in.Repos, in.Correlated.Findings, in.CollectErrors) {
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %d | %d | %s |\n",
			row.repo, row.crit, row.high, row.med, row.low, row.total, repoStatus(row))
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

	// Bot PR ledger
	b.WriteString("## 🤖 Bot PR ledger\n\n")
	if len(in.Ledger.Entries) == 0 {
		b.WriteString("_No bot PRs yet._\n\n")
	} else {
		b.WriteString("| Repo | Bump | Kind | Source | State | PR |\n|---|---|---|---|---|---|\n")
		for _, e := range in.Ledger.Entries {
			pr := "—"
			if e.PRNumber > 0 {
				pr = fmt.Sprintf("[#%d](%s)", e.PRNumber, e.PRURL)
			}
			kind := e.Kind
			if kind == "" {
				kind = "direct"
			}
			source := e.Source
			if source == "" {
				source = "ksec"
			}
			st := e.State
			if e.NeedsHuman {
				st = "⚠️ needs-human"
			} else if e.Blocked != "" {
				st = "⛔ " + e.Blocked
			}
			fmt.Fprintf(&b, "| %s | %s@%s | %s | %s | %s | %s |\n", e.Repo, e.Bump.Package, e.Bump.To, kind, source, st, pr)
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
	errored                     bool
}

// perRepoRows enumerates the union of every tracked repo, every repo that
// produced a finding, and every repo that produced a collection error, so that
// clean and errored repos remain visible. It is backward compatible: passing a
// nil repos slice still yields a row for each repo seen in findings or errs.
func perRepoRows(repos []state.Repo, findings []state.Finding, errs []state.CollectionError) []repoRow {
	idx := map[string]*repoRow{}
	get := func(name string) *repoRow {
		r := idx[name]
		if r == nil {
			r = &repoRow{repo: name}
			idx[name] = r
		}
		return r
	}
	for _, repo := range repos {
		get(repo.Repo)
	}
	for _, f := range findings {
		r := get(f.Repo)
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
	for _, e := range errs {
		get(e.Repo).errored = true
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
		if rows[i].total != rows[j].total {
			return rows[i].total > rows[j].total
		}
		return rows[i].repo < rows[j].repo
	})
	return rows
}

// repoStatus classifies a per-repo row for display: errored repos take
// precedence, then repos with no findings are "clean", otherwise "ok".
func repoStatus(r repoRow) string {
	switch {
	case r.errored:
		return "⚠️ errors"
	case r.total == 0:
		return "clean"
	default:
		return "ok"
	}
}
