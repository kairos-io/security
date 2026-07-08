package render

import (
	"fmt"
	"strings"
)

type RunActivity struct {
	Repos, Skipped, Errored                    int
	Findings, Crit, High, Med, Low, Unknown    int
	Informational                              int
	PRs                                        int
	PRsBySource                                map[string]int
	LedgerOpen, NeedsHuman, Superseded, Merged int
	Why                                        string
}

func computeActivity(in Input) RunActivity {
	a := RunActivity{Repos: len(in.Repos), PRsBySource: map[string]int{}}
	for _, r := range in.Repos {
		if !r.SourceScanEnabled() {
			a.Skipped++
		}
	}
	erroredRepos := map[string]bool{}
	for _, e := range in.CollectErrors {
		erroredRepos[e.Repo] = true
	}
	a.Errored = len(erroredRepos)
	for _, f := range in.Correlated.Findings {
		if f.Class == "informational" {
			a.Informational++
			continue
		}
		a.Findings++
		switch f.Severity {
		case "critical":
			a.Crit++
		case "high":
			a.High++
		case "medium":
			a.Med++
		case "low":
			a.Low++
		default:
			a.Unknown++
		}
	}
	a.PRs = len(in.OpenPRs)
	for _, p := range in.OpenPRs {
		a.PRsBySource[p.Source]++
	}
	for _, e := range in.Ledger.Entries {
		if e.NeedsHuman {
			a.NeedsHuman++
		}
		if e.Supersedes != "" {
			a.Superseded++
		}
		switch e.State {
		case "open":
			a.LedgerOpen++
		case "merged":
			a.Merged++
		}
	}
	a.Why = activityWhy(a)
	return a
}

func activityWhy(a RunActivity) string {
	switch {
	case a.Findings == 0 && a.Errored == 0:
		return fmt.Sprintf("No CVEs found across %d repos — nothing to remediate.", a.Repos)
	case a.Findings == 0 && a.Errored > 0:
		return fmt.Sprintf("No CVEs found, but %d repo(s) could not be scanned — see collection errors.", a.Errored)
	case a.NeedsHuman > 0:
		return fmt.Sprintf("%d finding(s); %d PR(s) open, %d need a human.", a.Findings, a.LedgerOpen, a.NeedsHuman)
	default:
		return fmt.Sprintf("%d finding(s); %d PR(s) open.", a.Findings, a.LedgerOpen)
	}
}

// Markdown renders the "📋 This run" section body.
func (a RunActivity) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "- **Scanned:** %d repos", a.Repos)
	if a.Skipped > 0 {
		fmt.Fprintf(&b, " (%d skipped)", a.Skipped)
	}
	if a.Errored > 0 {
		fmt.Fprintf(&b, " · ⚠️ %d errored", a.Errored)
	}
	b.WriteString("\n")
	fmt.Fprintf(&b, "- **Findings:** %d (%d critical / %d high / %d medium / %d low / %d unknown)\n",
		a.Findings, a.Crit, a.High, a.Med, a.Low, a.Unknown)
	if a.Informational > 0 {
		fmt.Fprintf(&b, "- **Informational (not counted):** %d\n", a.Informational)
	}
	fmt.Fprintf(&b, "- **CVE-related PRs:** %d", a.PRs)
	if a.PRs > 0 {
		b.WriteString(" (" + sourceBreakdown(a.PRsBySource) + ")")
	}
	b.WriteString("\n")
	fmt.Fprintf(&b, "- **Remediation:** %d open · %d superseded · %d merged · %d need-human\n",
		a.LedgerOpen, a.Superseded, a.Merged, a.NeedsHuman)
	fmt.Fprintf(&b, "- **Why:** %s\n", a.Why)
	return b.String()
}

func sourceBreakdown(m map[string]int) string {
	var parts []string
	for _, src := range []string{"ksec", "dependabot", "renovate", "bot", "human"} {
		if n := m[src]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", n, src))
		}
	}
	return strings.Join(parts, ", ")
}
