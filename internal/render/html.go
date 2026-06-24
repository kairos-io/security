package render

import (
	"html/template"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

// htmlData is the view model passed to the HTML template.
type htmlData struct {
	GeneratedAt         string
	AIAvailable         bool
	Activity            htmlActivity
	Narrative           string
	CoordinationSummary string
	Focus               []htmlFocus
	Waterfall           []state.WaterfallGroup
	Repos               []htmlRepoRow
	CollectErrors       []state.CollectionError
	OpenPRs             []htmlOpenPRGroup
	Ledger              []htmlLedgerEntry
	NeedsHuman          []string
	Reviews             []htmlReviewGroup
	RunURL              string
}

// htmlReviewGroup exposes bot-PR reviews for a single repo to the template.
type htmlReviewGroup struct {
	Repo    string
	Reviews []htmlReview
}

// htmlReview exposes a single bot-PR review to the template (exported fields).
type htmlReview struct {
	PR        int
	URL       string
	Icon      string
	Verdict   string
	Reasoning string
}

// htmlOpenPRGroup exposes open PRs for a single repo to the template.
type htmlOpenPRGroup struct {
	Repo string
	PRs  []htmlOpenPR
}

// htmlOpenPR exposes a single tracked PR to the template (exported fields).
type htmlOpenPR struct {
	Number int
	Title  string
	URL    string
	Source string
	Status string
}

// htmlLedgerEntry exposes a bot PR ledger row to the template (exported fields).
type htmlLedgerEntry struct {
	Repo     string
	Bump     string
	Kind     string
	Source   string
	State    string
	PRNumber int
	PRURL    string
}

// htmlActivity exposes the deterministic "This run" summary to the template
// (exported fields). SourceBreakdown is the fixed-order PR-source string.
type htmlActivity struct {
	Repos, Skipped, Errored                    int
	Findings, Crit, High, Med, Low, Unknown    int
	PRs                                        int
	SourceBreakdown                            string
	LedgerOpen, Superseded, Merged, NeedsHuman int
	Why                                        string
}

type htmlFocus struct {
	Title   string
	URL     string
	Summary string
}

// htmlRepoRow exposes per-repo counts to the template (exported fields).
type htmlRepoRow struct {
	Repo                        string
	Crit, High, Med, Low, Total int
	Status                      string
	StatusClass                 string
}

// dashboardHTMLTmpl renders a self-contained page. All dynamic values are
// emitted through html/template's contextual escaping, so repo names, titles
// and messages containing <, & or | cannot break the markup or inject HTML.
var dashboardHTMLTmpl = template.Must(template.New("dashboard").Funcs(template.FuncMap{
	"join":     func(s []string) string { return strings.Join(s, ", ") },
	"sevClass": severityClass,
}).Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Kairos Security Dashboard</title>
<style>
:root { color-scheme: light dark; }
body { font-family: system-ui, -apple-system, Segoe UI, Roboto, sans-serif; margin: 0; padding: 2rem; line-height: 1.5; color: #1b1f23; background: #fff; }
h1 { margin-top: 0; }
.meta { color: #586069; font-size: 0.95rem; }
.ai-warn { color: #b00020; font-weight: 600; }
blockquote { border-left: 4px solid #d0d7de; margin: 1rem 0; padding: 0.25rem 1rem; color: #444; background: #f6f8fa; }
section { margin: 2rem 0; }
table { border-collapse: collapse; width: 100%; }
th, td { border: 1px solid #d0d7de; padding: 0.4rem 0.6rem; text-align: left; }
th { background: #f6f8fa; }
td.num { text-align: right; }
.sev-critical { background: #ffd6d6; color: #86181d; font-weight: 600; }
.sev-high { background: #ffe3c2; color: #8a4b00; font-weight: 600; }
.sev-medium { background: #fff7c2; color: #735c00; }
.sev-low { background: #e8eaed; color: #444; }
.status-errored { color: #8a4b00; font-weight: 600; }
.status-clean { color: #586069; }
.empty { color: #586069; font-style: italic; }
footer { margin-top: 2rem; border-top: 1px solid #d0d7de; padding-top: 1rem; color: #586069; font-size: 0.9rem; }
code { background: #f6f8fa; padding: 0.1rem 0.3rem; border-radius: 3px; }
</style>
</head>
<body>
<h1>Kairos Security Dashboard</h1>
<p class="meta">Updated {{.GeneratedAt}}{{if not .AIAvailable}} &mdash; <span class="ai-warn">&#9888;&#65039; AI unavailable this run</span>{{end}}</p>
{{if .Narrative}}<blockquote>{{.Narrative}}</blockquote>{{end}}

<section>
<h2>&#128203; This run</h2>
<ul>
<li><strong>Scanned:</strong> {{.Activity.Repos}} repos{{if gt .Activity.Skipped 0}} ({{.Activity.Skipped}} skipped){{end}}{{if gt .Activity.Errored 0}} &middot; &#9888;&#65039; {{.Activity.Errored}} errored{{end}}</li>
<li><strong>Findings:</strong> {{.Activity.Findings}} ({{.Activity.Crit}} critical / {{.Activity.High}} high / {{.Activity.Med}} medium / {{.Activity.Low}} low / {{.Activity.Unknown}} unknown)</li>
<li><strong>CVE-related PRs:</strong> {{.Activity.PRs}}{{if gt .Activity.PRs 0}} ({{.Activity.SourceBreakdown}}){{end}}</li>
<li><strong>Remediation:</strong> {{.Activity.LedgerOpen}} open &middot; {{.Activity.Superseded}} superseded &middot; {{.Activity.Merged}} merged &middot; {{.Activity.NeedsHuman}} need-human</li>
<li><strong>Why:</strong> {{.Activity.Why}}</li>
</ul>
</section>
{{- if .CoordinationSummary}}

<section>
<h2>&#129517; Coordination</h2>
<p>{{.CoordinationSummary}}</p>
</section>
{{- end}}

<section>
<h2>&#128293; Focus now</h2>
{{if .Focus}}<ol>
{{range .Focus}}<li>{{if .URL}}<a href="{{.URL}}">{{.Title}}</a>{{else}}<strong>{{.Title}}</strong>{{end}}{{if .Summary}} &mdash; {{.Summary}}{{end}}</li>
{{end}}</ol>{{else}}<p class="empty">Nothing flagged.</p>{{end}}
</section>

<section>
<h2>&#127754; Waterfall fronts</h2>
{{if .Waterfall}}<table>
<thead><tr><th>Root cause</th><th>Severity</th><th>Suggested bump</th><th>Affected repos</th></tr></thead>
<tbody>
{{range .Waterfall}}<tr><td>{{.RootCause}}</td><td class="{{sevClass .Severity}}">{{.Severity}}</td><td><code>{{.SuggestedBump.Package}}@{{.SuggestedBump.To}}</code></td><td>{{join .AffectedRepos}}</td></tr>
{{end}}</tbody>
</table>{{else}}<p class="empty">None.</p>{{end}}
</section>

<section>
<h2>&#128230; Per-repo findings</h2>
<table>
<thead><tr><th>Repo</th><th>Critical</th><th>High</th><th>Medium</th><th>Low</th><th>Total</th><th>Status</th></tr></thead>
<tbody>
{{range .Repos}}<tr><td>{{.Repo}}</td><td class="num sev-critical">{{.Crit}}</td><td class="num sev-high">{{.High}}</td><td class="num sev-medium">{{.Med}}</td><td class="num sev-low">{{.Low}}</td><td class="num">{{.Total}}</td><td class="{{.StatusClass}}">{{.Status}}</td></tr>
{{else}}<tr><td colspan="7" class="empty">No repos tracked.</td></tr>
{{end}}</tbody>
</table>
</section>

{{if .CollectErrors}}<section>
<h2>&#9888;&#65039; {{len .CollectErrors}} collection error{{if ne (len .CollectErrors) 1}}s{{end}}</h2>
<ul>
{{range .CollectErrors}}<li><code>{{.Repo}}</code> / {{.Collector}}: {{.Message}}</li>
{{end}}</ul>
</section>{{end}}

<section>
<h2>&#128203; Open PRs</h2>
{{if .OpenPRs}}{{range .OpenPRs}}<h3>{{.Repo}}</h3>
<ul>
{{range .PRs}}<li>{{if .URL}}<a href="{{.URL}}">#{{.Number}} {{.Title}}</a>{{else}}#{{.Number}} {{.Title}}{{end}} &mdash; {{.Source}} &mdash; {{.Status}}</li>
{{end}}</ul>
{{end}}{{else}}<p class="empty">None.</p>{{end}}
</section>

<section>
<h2>&#129302; Bot PR ledger</h2>
{{if .Ledger}}<table>
<thead><tr><th>Repo</th><th>Bump</th><th>Kind</th><th>Source</th><th>State</th><th>PR</th></tr></thead>
<tbody>
{{range .Ledger}}<tr><td>{{.Repo}}</td><td><code>{{.Bump}}</code></td><td>{{.Kind}}</td><td>{{.Source}}</td><td>{{.State}}</td><td>{{if gt .PRNumber 0}}<a href="{{.PRURL}}">#{{.PRNumber}}</a>{{else}}&mdash;{{end}}</td></tr>
{{end}}</tbody>
</table>{{else}}<p class="empty">No bot PRs yet.</p>{{end}}
</section>
{{if .NeedsHuman}}
<section>
<h2>&#128657; Needs human</h2>
<ul>
{{range .NeedsHuman}}<li>{{.}}</li>
{{end}}</ul>
</section>
{{end}}
{{- if .Reviews}}
<section>
<h2>&#128270; Bot-PR reviews</h2>
{{range .Reviews}}<h3>{{.Repo}}</h3>
<ul>
{{range .Reviews}}<li>{{if .URL}}<a href="{{.URL}}">#{{.PR}}</a>{{else}}#{{.PR}}{{end}} &mdash; {{.Icon}} <strong>{{.Verdict}}</strong> &mdash; {{.Reasoning}}</li>
{{end}}</ul>
{{end}}</section>
{{- end}}
<footer>
{{if .RunURL}}<a href="{{.RunURL}}">Run log</a>{{else}}Kairos central security dashboard{{end}}
</footer>
</body>
</html>
`))

// repoStatusClass returns the CSS class for a repo's status cell: errored rows
// are amber/red, clean rows are muted gray, ok rows use the default style.
func repoStatusClass(r repoRow) string {
	switch {
	case r.errored:
		return "status-errored"
	case r.total == 0:
		return "status-clean"
	default:
		return ""
	}
}

func severityClass(sev string) string {
	switch sev {
	case "critical":
		return "sev-critical"
	case "high":
		return "sev-high"
	case "medium":
		return "sev-medium"
	case "low":
		return "sev-low"
	default:
		return ""
	}
}

// DashboardHTML renders a self-contained HTML dashboard page. It is
// deterministic for a given Input and escapes all dynamic text.
func DashboardHTML(in Input) string {
	byID := map[string]state.Finding{}
	for _, f := range in.Correlated.Findings {
		byID[f.ID] = f
	}
	focus := make([]htmlFocus, 0, len(in.Triage.Focus))
	for _, id := range in.Triage.Focus {
		hf := htmlFocus{}
		if f, ok := byID[id]; ok {
			hf.Title, hf.URL = focusTitleURL(f)
			if s := in.Triage.Summaries[id]; s != "" && !strings.HasPrefix(s, "Finding in ") {
				hf.Summary = s
			}
		} else if s, ok := in.Triage.Summaries[id]; ok {
			hf.Title = s
		} else {
			continue
		}
		focus = append(focus, hf)
	}
	rows := perRepoRows(in.Repos, in.Correlated.Findings, in.CollectErrors)
	repos := make([]htmlRepoRow, 0, len(rows))
	for _, r := range rows {
		repos = append(repos, htmlRepoRow{
			Repo: r.repo, Crit: r.crit, High: r.high, Med: r.med, Low: r.low, Total: r.total,
			Status: repoStatus(r), StatusClass: repoStatusClass(r),
		})
	}
	ledger := make([]htmlLedgerEntry, 0, len(in.Ledger.Entries))
	for _, e := range in.Ledger.Entries {
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
		bump := e.Bump.Package + "@" + e.Bump.To
		if e.Pseudo {
			bump += " (pseudo)"
		}
		if e.Supersedes != "" {
			bump += " ↳ supersedes " + e.Supersedes
		}
		ledger = append(ledger, htmlLedgerEntry{
			Repo: e.Repo, Bump: bump,
			Kind: kind, Source: source,
			State: st, PRNumber: e.PRNumber, PRURL: e.PRURL,
		})
	}
	supersededBy, conflictedURLs := correlateOpenPRs(in.Ledger.Entries)
	var openPRs []htmlOpenPRGroup
	for _, pr := range in.OpenPRs {
		if len(openPRs) == 0 || openPRs[len(openPRs)-1].Repo != pr.Repo {
			openPRs = append(openPRs, htmlOpenPRGroup{Repo: pr.Repo})
		}
		g := &openPRs[len(openPRs)-1]
		g.PRs = append(g.PRs, htmlOpenPR{
			Number: pr.Number, Title: pr.Title, URL: pr.URL, Source: pr.Source,
			Status: openPRStatus(pr, supersededBy, conflictedURLs),
		})
	}
	var reviews []htmlReviewGroup
	for _, r := range in.Reviews {
		if len(reviews) == 0 || reviews[len(reviews)-1].Repo != r.Repo {
			reviews = append(reviews, htmlReviewGroup{Repo: r.Repo})
		}
		g := &reviews[len(reviews)-1]
		g.Reviews = append(g.Reviews, htmlReview{
			PR: r.PR, URL: r.URL, Icon: verdictIcon(r.Verdict), Verdict: r.Verdict, Reasoning: r.Reasoning,
		})
	}
	act := computeActivity(in)
	data := htmlData{
		GeneratedAt: in.Triage.GeneratedAt,
		AIAvailable: in.Triage.AIAvailable,
		Activity: htmlActivity{
			Repos: act.Repos, Skipped: act.Skipped, Errored: act.Errored,
			Findings: act.Findings, Crit: act.Crit, High: act.High, Med: act.Med, Low: act.Low, Unknown: act.Unknown,
			PRs: act.PRs, SourceBreakdown: sourceBreakdown(act.PRsBySource),
			LedgerOpen: act.LedgerOpen, Superseded: act.Superseded, Merged: act.Merged, NeedsHuman: act.NeedsHuman,
			Why: act.Why,
		},
		Narrative:           in.Triage.Narrative,
		CoordinationSummary: in.CoordinationSummary,
		Focus:               focus,
		Waterfall:           in.Correlated.Waterfall,
		Repos:               repos,
		CollectErrors:       in.CollectErrors,
		OpenPRs:             openPRs,
		Ledger:              ledger,
		NeedsHuman:          needsHumanRows(in.Ledger.Entries),
		Reviews:             reviews,
		RunURL:              in.RunURL,
	}
	var b strings.Builder
	if err := dashboardHTMLTmpl.Execute(&b, data); err != nil {
		// Template is static and fields are simple; execution should not fail.
		return "<!DOCTYPE html><html><body><h1>Kairos Security Dashboard</h1><p>render error: " +
			template.HTMLEscapeString(err.Error()) + "</p></body></html>"
	}
	return b.String()
}
