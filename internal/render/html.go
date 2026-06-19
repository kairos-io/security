package render

import (
	"html/template"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

// htmlData is the view model passed to the HTML template.
type htmlData struct {
	GeneratedAt   string
	AIAvailable   bool
	Narrative     string
	Focus         []htmlFocus
	Waterfall     []state.WaterfallGroup
	Repos         []htmlRepoRow
	CollectErrors []state.CollectionError
	Ledger        []htmlLedgerEntry
	RunURL        string
}

// htmlLedgerEntry exposes a bot PR ledger row to the template (exported fields).
type htmlLedgerEntry struct {
	Repo     string
	Bump     string
	State    string
	PRNumber int
	PRURL    string
}

type htmlFocus struct {
	ID      string
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
<h2>&#128293; Focus now</h2>
{{if .Focus}}<ol>
{{range .Focus}}<li><strong>{{.ID}}</strong>{{if .Summary}} &mdash; {{.Summary}}{{end}}</li>
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
<h2>&#129302; Bot PR ledger</h2>
{{if .Ledger}}<table>
<thead><tr><th>Repo</th><th>Bump</th><th>State</th><th>PR</th></tr></thead>
<tbody>
{{range .Ledger}}<tr><td>{{.Repo}}</td><td><code>{{.Bump}}</code></td><td>{{.State}}</td><td>{{if gt .PRNumber 0}}<a href="{{.PRURL}}">#{{.PRNumber}}</a>{{else}}&mdash;{{end}}</td></tr>
{{end}}</tbody>
</table>{{else}}<p class="empty">No bot PRs yet.</p>{{end}}
</section>

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
	focus := make([]htmlFocus, 0, len(in.Triage.Focus))
	for _, id := range in.Triage.Focus {
		focus = append(focus, htmlFocus{ID: id, Summary: in.Triage.Summaries[id]})
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
		ledger = append(ledger, htmlLedgerEntry{
			Repo: e.Repo, Bump: e.Bump.Package + "@" + e.Bump.To,
			State: e.State, PRNumber: e.PRNumber, PRURL: e.PRURL,
		})
	}
	data := htmlData{
		GeneratedAt:   in.Triage.GeneratedAt,
		AIAvailable:   in.Triage.AIAvailable,
		Narrative:     in.Triage.Narrative,
		Focus:         focus,
		Waterfall:     in.Correlated.Waterfall,
		Repos:         repos,
		CollectErrors: in.CollectErrors,
		Ledger:        ledger,
		RunURL:        in.RunURL,
	}
	var b strings.Builder
	if err := dashboardHTMLTmpl.Execute(&b, data); err != nil {
		// Template is static and fields are simple; execution should not fail.
		return "<!DOCTYPE html><html><body><h1>Kairos Security Dashboard</h1><p>render error: " +
			template.HTMLEscapeString(err.Error()) + "</p></body></html>"
	}
	return b.String()
}
