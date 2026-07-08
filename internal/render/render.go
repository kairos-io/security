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
	// CoordinationSummary is an AI-generated cross-repo coordination narrative.
	CoordinationSummary string `json:"coordinationSummary,omitempty"`
	// OpenPRs are the tracked remediation-relevant pull requests, grouped by repo.
	OpenPRs []state.TrackedPR `json:"openPRs,omitempty"`
	// Reviews are the AI bot-PR reviews, grouped by repo in the dashboard.
	Reviews []state.PRReview `json:"reviews,omitempty"`
}

// verdictIcon maps a bot-PR review verdict to a display icon.
func verdictIcon(verdict string) string {
	switch verdict {
	case "good":
		return "✅"
	case "bad":
		return "⛔"
	case "needs_human_verification":
		return "⚠️"
	default:
		return ""
	}
}

// findingLink renders a finding as a markdown title→URL link. For sourceCVE
// findings with no URL it synthesizes an advisory link from the CVE/GHSA id.
// It never emits the raw finding id: the title (or package, or "<repo>
// finding") is used as the link text. When the AI applicability classifier
// suspects the CVE does not affect us, a ⚠️ marker is appended (details are
// listed in the dedicated "AI-flagged" section further down).
func findingLink(f state.Finding) string {
	title, url := focusTitleURL(f)
	base := title
	if url != "" {
		base = fmt.Sprintf("[%s](%s)", title, url)
	}
	if suspectedNotApplicable(f) {
		base += " ⚠️"
	}
	return base
}

// nonEmptyMD returns v or "—" when v is blank, so table cells never render
// as trailing whitespace.
func nonEmptyMD(v string) string {
	if strings.TrimSpace(v) == "" {
		return "—"
	}
	return v
}

// suspectedNotApplicable reports whether the AI classifier attached a
// non-applicable verdict to this finding. Renderers surface a warning + reason
// on these but do NOT hide them (fail-visible — the verdict is advisory only).
func suspectedNotApplicable(f state.Finding) bool {
	return f.AIApplicability != nil && !f.AIApplicability.Applicable
}

// applicabilityWarnings returns the findings the AI classifier flagged as
// possibly-not-applicable, sorted by severity then repo/package for stable
// display.
func applicabilityWarnings(findings []state.Finding) []state.Finding {
	var out []state.Finding
	for _, f := range findings {
		if suspectedNotApplicable(f) {
			out = append(out, f)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		si, sj := severityRank(out[i].Severity), severityRank(out[j].Severity)
		if si != sj {
			return si > sj
		}
		if out[i].Repo != out[j].Repo {
			return out[i].Repo < out[j].Repo
		}
		return out[i].Package < out[j].Package
	})
	return out
}

// focusTitleURL derives the human-facing title and (possibly synthesized)
// advisory URL for a finding, shared by the markdown and HTML renderers. It
// never returns the raw finding id.
func focusTitleURL(f state.Finding) (title, url string) {
	url = f.URL
	if url == "" && f.Type == "sourceCVE" {
		id := f.CVEID
		if id == "" {
			id = f.GHSA
		}
		// pkg.go.dev/vuln only serves GO-… paths; a CVE/GHSA id there 400s.
		// Synthesize only for GO-prefixed ids; otherwise leave url empty
		// (bare title) rather than emit a dead link.
		if strings.HasPrefix(id, "GO-") {
			url = "https://pkg.go.dev/vuln/" + id
		}
	}
	title = f.Title
	if title == "" {
		title = f.CVEID
	}
	if title == "" {
		title = f.GHSA
	}
	if title == "" {
		title = f.Package
	}
	if title == "" {
		title = f.Repo + " finding"
	}
	return title, url
}

func DashboardJSON(in Input) ([]byte, error) {
	b, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

// repoLink renders a repo slug as a markdown link to its GitHub page.
func repoLink(repo string) string {
	return fmt.Sprintf("[%s](https://github.com/%s)", repo, repo)
}

func DashboardMarkdown(in Input) string {
	var b strings.Builder
	b.WriteString("# Kairos Security Dashboard\n\n")
	fmt.Fprintf(&b, "_Updated %s", in.Triage.GeneratedAt)
	if !in.Triage.AIAvailable {
		b.WriteString(" — ⚠️ AI unavailable this run")
	}
	b.WriteString("._\n\n")
	b.WriteString("🌐 **[Live dashboard](https://kairos-io.github.io/security/)** — the published board with clickable links.\n\n")

	// Deterministic run-activity summary, computed from committed state.
	a := computeActivity(in)
	b.WriteString("## 📋 This run\n\n")
	b.WriteString(a.Markdown() + "\n")

	if in.Triage.Narrative != "" {
		b.WriteString("> " + in.Triage.Narrative + "\n\n")
	}

	// Cross-repo coordination narrative
	if in.CoordinationSummary != "" {
		b.WriteString("## 🧭 Coordination\n\n")
		b.WriteString(in.CoordinationSummary + "\n\n")
	}

	// Focus now
	byID := map[string]state.Finding{}
	for _, f := range in.Correlated.Findings {
		byID[f.ID] = f
	}
	b.WriteString("## 🔥 Focus now\n\n")
	if len(in.Triage.Focus) == 0 {
		b.WriteString("_Nothing flagged._\n\n")
	} else {
		for _, id := range in.Triage.Focus {
			if f, ok := byID[id]; ok {
				line := findingLink(f)
				if s := in.Triage.Summaries[id]; s != "" && !strings.HasPrefix(s, "Finding in ") {
					line += " — " + s
				}
				fmt.Fprintf(&b, "- %s\n", line)
			} else if s, ok := in.Triage.Summaries[id]; ok {
				fmt.Fprintf(&b, "- %s\n", s)
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
	b.WriteString("| Repo | Critical | High | Medium | Total | Status |\n|---|---|---|---|---|---|\n")
	for _, row := range perRepoRows(in.Repos, in.Correlated.Findings, in.CollectErrors) {
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %d | %s |\n",
			repoLink(row.repo), row.crit, row.high, row.med, row.total, repoStatus(row))
	}
	b.WriteString("\n")

	// Hadron component CVEs
	if rows := hadronComponentRows(in.Correlated.Findings); len(rows) > 0 {
		b.WriteString("## 🧩 Hadron component CVEs\n\n")
		b.WriteString("| Package | Current | Fixed | Severity | CVE |\n|---|---|---|---|---|\n")
		for _, f := range rows {
			fixed := f.FixedVersion
			if fixed == "" {
				fixed = "—"
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", f.Package, f.CurrentVersion, fixed, f.Severity, findingLink(f))
		}
		b.WriteString("\n")
	}

	// AI-flagged possibly-not-applicable findings. Big prominent section: the
	// finding is STILL counted and shown in every other table, this section
	// exists so an operator can see the model's reasoning at a glance and
	// override / silence it if they disagree.
	if warns := applicabilityWarnings(in.Correlated.Findings); len(warns) > 0 {
		fmt.Fprintf(&b, "## ⚠️ %d finding(s) possibly not applicable (AI)\n\n", len(warns))
		b.WriteString("> These findings are still counted and listed above. The AI applicability check thinks they may not affect us — verify the reasoning below and, if you agree, silence via `cve-policy.yaml`.\n\n")
		for _, f := range warns {
			ai := f.AIApplicability
			title, url := focusTitleURL(f)
			label := title
			if url != "" {
				label = fmt.Sprintf("[%s](%s)", title, url)
			}
			fmt.Fprintf(&b, "<details>\n<summary>⚠️ %s — %s (%s / confidence: %s)</summary>\n\n", label, repoLink(f.Repo), f.Package, ai.Confidence)
			fmt.Fprintf(&b, "**Reason:** %s\n\n", ai.Reasoning)
			if f.CVEID != "" {
				fmt.Fprintf(&b, "- CVE: `%s`\n", f.CVEID)
			}
			fmt.Fprintf(&b, "- Current: `%s`\n", nonEmptyMD(f.CurrentVersion))
			fmt.Fprintf(&b, "- Fixed: `%s`\n", nonEmptyMD(f.FixedVersion))
			if ai.Model != "" {
				fmt.Fprintf(&b, "- Checked by: `%s`", ai.Model)
				if ai.CheckedAt != "" {
					fmt.Fprintf(&b, " on %s", ai.CheckedAt)
				}
				b.WriteString("\n")
			}
			b.WriteString("\n</details>\n\n")
		}
	}

	// Informational — not counted
	if info := informationalFindings(in.Correlated.Findings); len(info) > 0 {
		b.WriteString("## Informational — not counted\n\n")
		b.WriteString("These findings are separated from the counts above: CVEs we are already past, or components accepted as pinned risk.\n\n")
		b.WriteString("| Package | Current | Fixed | Severity | CVE | Why |\n|---|---|---|---|---|---|\n")
		for _, f := range info {
			fixed := f.FixedVersion
			if fixed == "" {
				fixed = "—"
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n",
				f.Package, f.CurrentVersion, fixed, f.Severity, findingLink(f), f.ClassReason)
		}
		b.WriteString("\n")
	}

	// Collection errors
	if len(in.CollectErrors) > 0 {
		fmt.Fprintf(&b, "## ⚠️ %d collection errors\n\n", len(in.CollectErrors))
		for _, e := range in.CollectErrors {
			fmt.Fprintf(&b, "- %s / %s: %s\n", repoLink(e.Repo), e.Collector, e.Message)
		}
		b.WriteString("\n")
	}

	// Open PRs
	supersededBy, conflictedURLs := correlateOpenPRs(in.Ledger.Entries)
	b.WriteString("## 📋 Open PRs\n\n")
	if len(in.OpenPRs) == 0 {
		b.WriteString("_None._\n\n")
	} else {
		repo := ""
		for _, pr := range in.OpenPRs {
			if pr.Repo != repo {
				repo = pr.Repo
				fmt.Fprintf(&b, "**%s**\n\n", repoLink(repo))
			}
			link := fmt.Sprintf("#%d %s", pr.Number, pr.Title)
			if pr.URL != "" {
				link = fmt.Sprintf("[#%d %s](%s)", pr.Number, pr.Title, pr.URL)
			}
			fmt.Fprintf(&b, "- %s — %s — %s\n", link, pr.Source, openPRStatus(pr, supersededBy, conflictedURLs))
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
			if e.CascadeFrom != "" {
				kind = "cascade ↳ from " + e.CascadeFrom
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
			bump := fmt.Sprintf("%s@%s", e.Bump.Package, e.Bump.To)
			if e.Pseudo {
				bump += " (pseudo)"
			}
			if e.Supersedes != "" {
				bump += " ↳ supersedes " + e.Supersedes
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n", repoLink(e.Repo), bump, kind, source, st, pr)
		}
		b.WriteString("\n")
	}

	// Needs human roll-up
	if rows := needsHumanRows(in.Ledger.Entries); len(rows) > 0 {
		b.WriteString("## 🚑 Needs human\n\n")
		for _, r := range rows {
			fmt.Fprintf(&b, "- %s\n", r)
		}
		b.WriteString("\n")
	}

	// Bot-PR reviews
	if len(in.Reviews) > 0 {
		b.WriteString("## 🔎 Bot-PR reviews\n\n")
		repo := ""
		for _, r := range in.Reviews {
			if r.Repo != repo {
				repo = r.Repo
				fmt.Fprintf(&b, "**%s**\n\n", repoLink(repo))
			}
			link := fmt.Sprintf("#%d", r.PR)
			if r.URL != "" {
				link = fmt.Sprintf("[#%d](%s)", r.PR, r.URL)
			}
			fmt.Fprintf(&b, "- %s — %s **%s** — %s\n", link, verdictIcon(r.Verdict), r.Verdict, r.Reasoning)
			if r.ChangesSummary != "" {
				fmt.Fprintf(&b, "  ↳ %s\n", r.ChangesSummary)
			}
			for _, tl := range r.Trace {
				fmt.Fprintf(&b, "    - %s\n", tl)
			}
		}
		b.WriteString("\n")
	}

	if in.RunURL != "" {
		fmt.Fprintf(&b, "---\n[Run log](%s)\n", in.RunURL)
	}
	return b.String()
}

type repoRow struct {
	repo                   string
	crit, high, med, total int
	errored                bool
	skipped                bool // source scanning opted out for this repo
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
		get(repo.Repo).skipped = !repo.SourceScanEnabled()
	}
	for _, f := range findings {
		if f.Class == "informational" {
			continue // separated + uncounted; see informationalFindings
		}
		r := get(f.Repo)
		switch f.Severity {
		case "critical":
			r.crit++
			r.total++
		case "high":
			r.high++
			r.total++
		case "medium":
			r.med++
			r.total++
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

// hadronComponentRows returns componentCVE findings sorted by severity
// (critical first) then package name.
func hadronComponentRows(findings []state.Finding) []state.Finding {
	var rows []state.Finding
	for _, f := range findings {
		if f.Class == "informational" {
			continue // separated into the informational section, not counted here
		}
		if f.Type == "componentCVE" {
			rows = append(rows, f)
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		si, sj := severityRank(rows[i].Severity), severityRank(rows[j].Severity)
		if si != sj {
			return si > sj
		}
		return rows[i].Package < rows[j].Package
	})
	return rows
}

// informationalFindings returns findings classed "informational" — CVEs we are
// already past, or components accepted as pinned risk — in input order. They are
// separated from the actionable tables and excluded from every count.
func informationalFindings(findings []state.Finding) []state.Finding {
	var info []state.Finding
	for _, f := range findings {
		if f.Class == "informational" {
			info = append(info, f)
		}
	}
	return info
}

func severityRank(s string) int {
	switch s {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}

// correlateOpenPRs builds two lookup maps from the ledger in a single pass:
// supersededBy maps a foreign PR URL to the ksec entry that supersedes it
// (keyed on e.Supersedes), and conflictedURLs marks foreign PR URLs that have a
// ledger entry blocked on an upstream conflict (keyed on e.PRURL). Both are
// derived from the already-sorted ledger, so no map-order leak reaches output.
func correlateOpenPRs(entries []state.LedgerEntry) (supersededBy map[string]state.LedgerEntry, conflictedURLs map[string]bool) {
	supersededBy = map[string]state.LedgerEntry{}
	conflictedURLs = map[string]bool{}
	for _, e := range entries {
		if e.Supersedes != "" {
			supersededBy[e.Supersedes] = e
		}
		if e.PRURL != "" && e.Blocked == "upstream-conflict" {
			conflictedURLs[e.PRURL] = true
		}
	}
	return supersededBy, conflictedURLs
}

// openPRStatus reports the human-facing status/action for a tracked PR by
// correlating it with the ledger: a superseding ksec PR, an upstream conflict
// being superseded, or plain "tracked".
func openPRStatus(pr state.TrackedPR, supersededBy map[string]state.LedgerEntry, conflictedURLs map[string]bool) string {
	if e, ok := supersededBy[pr.URL]; ok {
		return fmt.Sprintf("superseded by [#%d](%s)", e.PRNumber, e.PRURL)
	}
	if conflictedURLs[pr.URL] {
		return "conflicted → superseding"
	}
	return "tracked"
}

// needsHumanRows lists, in ledger order, each entry flagged NeedsHuman as
// "<repo> <bump> — <Blocked or state>".
func needsHumanRows(entries []state.LedgerEntry) []string {
	var rows []string
	for _, e := range entries {
		if !e.NeedsHuman {
			continue
		}
		reason := e.Blocked
		if reason == "" {
			reason = e.State
		}
		rows = append(rows, fmt.Sprintf("%s %s@%s — %s", e.Repo, e.Bump.Package, e.Bump.To, reason))
	}
	return rows
}

// repoStatus classifies a per-repo row for display: errored repos take
// precedence, then source-scan opt-outs with no findings are "skipped: not
// source-scannable", then repos with no critical/high/medium findings are
// "clean (no crit/high/med)" — note low-severity findings may still exist for
// that repo; they're just not counted in this table. Otherwise "ok".
func repoStatus(r repoRow) string {
	switch {
	case r.errored:
		return "⚠️ errors"
	case r.total == 0 && r.skipped:
		return "skipped: not source-scannable"
	case r.total == 0:
		return "clean (no crit/high/med)"
	default:
		return "ok"
	}
}
