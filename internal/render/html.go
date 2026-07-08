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
	Components          []htmlComponentRow
	Informational       []htmlInfoRow
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
	PR             int
	URL            string
	Icon           string
	Verdict        string
	Reasoning      string
	ChangesSummary string
	Trace          []string
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
	Repo                   string
	Crit, High, Med, Total int
	Status                 string
	StatusClass            string
}

// htmlComponentRow exposes a hadron component-manifest CVE to the template.
type htmlComponentRow struct {
	Package, Current, Fixed, Severity, Title, URL string
}

// htmlInfoRow exposes an informational (separated + uncounted) finding to the
// template. Why carries the ClassReason explaining the separation.
type htmlInfoRow struct {
	Package, Current, Fixed, Severity, Title, URL, Why string
}

// dashboardHTMLTmpl renders a self-contained page. All dynamic values are
// emitted through html/template's contextual escaping, so repo names, titles
// and messages containing <, & or | cannot break the markup or inject HTML.
var dashboardHTMLTmpl = template.Must(template.New("dashboard").Funcs(template.FuncMap{
	"join":     func(s []string) string { return strings.Join(s, ", ") },
	"sevClass": severityClass,
	"repoURL":  func(r string) string { return "https://github.com/" + r },
}).Parse(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Kairos Security Dashboard</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@500;600;700&family=IBM+Plex+Sans:wght@400;500;600&family=IBM+Plex+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
:root{
  --bg: oklch(0.99 0.004 285); --surface: oklch(1 0 0); --surface-2: oklch(0.97 0.006 285);
  --border: oklch(0.91 0.008 285); --line: oklch(0.94 0.006 285);
  --text: oklch(0.28 0.03 285); --muted: oklch(0.52 0.025 285);
  --accent: oklch(0.52 0.19 285); --accent-weak: oklch(0.95 0.03 285);
  --crit: oklch(0.55 0.2 27); --crit-bg: oklch(0.95 0.05 27);
  --high: oklch(0.58 0.15 55); --high-bg: oklch(0.95 0.05 65);
  --med: oklch(0.62 0.11 90); --med-bg: oklch(0.96 0.05 95);
  --low: oklch(0.55 0.015 285); --low-bg: oklch(0.95 0.006 285);
  --ok: oklch(0.55 0.13 155); --ok-bg: oklch(0.95 0.04 155);
  --warn: oklch(0.58 0.14 60); --warn-bg: oklch(0.96 0.05 70);
}
@media (prefers-color-scheme: dark){
  :root{
    --bg: oklch(0.2 0.02 285); --surface: oklch(0.24 0.02 285); --surface-2: oklch(0.27 0.02 285);
    --border: oklch(0.36 0.02 285); --line: oklch(0.31 0.02 285);
    --text: oklch(0.92 0.015 285); --muted: oklch(0.7 0.02 285);
    --accent: oklch(0.74 0.14 285); --accent-weak: oklch(0.31 0.05 285);
    --crit-bg: oklch(0.32 0.08 27); --high-bg: oklch(0.33 0.07 60); --med-bg: oklch(0.34 0.06 90);
    --low-bg: oklch(0.3 0.01 285); --ok-bg: oklch(0.32 0.06 155); --warn-bg: oklch(0.34 0.07 65);
  }
}
*{ box-sizing: border-box; }
html{ scroll-behavior: smooth; }
body{ font-family: "IBM Plex Sans", system-ui, sans-serif; margin: 0; color: var(--text); background: var(--bg);
  line-height: 1.55; font-size: 15px; -webkit-font-smoothing: antialiased; }
a{ color: var(--accent); text-decoration: none; }
a:hover{ text-decoration: underline; text-underline-offset: 2px; }
.wrap{ max-width: 1080px; margin: 0 auto; padding: 0 clamp(1rem, 4vw, 2.5rem); }
code, .mono{ font-family: "IBM Plex Mono", ui-monospace, monospace; font-size: 0.86em; }

/* Masthead */
header.mast{ position: sticky; top: 0; z-index: 20; background: color-mix(in oklch, var(--bg) 88%, transparent);
  backdrop-filter: saturate(1.4) blur(8px); border-bottom: 1px solid var(--border); }
.mast .wrap{ display: flex; align-items: baseline; gap: 1rem 1.25rem; flex-wrap: wrap; padding-top: 0.9rem; padding-bottom: 0.9rem; }
.brand{ font-family: "Space Grotesk", sans-serif; font-weight: 700; font-size: 1.15rem; letter-spacing: -0.01em; color: var(--text); }
.brand .dot{ color: var(--accent); }
.mast .upd{ color: var(--muted); font-size: 0.85rem; }
.mast .grow{ flex: 1; }
.ai-warn{ color: var(--crit); font-weight: 600; font-size: 0.85rem; }

/* Section nav */
nav.toc{ position: sticky; top: 53px; z-index: 19; background: var(--bg); border-bottom: 1px solid var(--line); }
nav.toc .wrap{ display: flex; gap: 0.35rem; flex-wrap: wrap; padding-top: 0.5rem; padding-bottom: 0.5rem; overflow-x: auto; }
nav.toc a{ color: var(--muted); font-size: 0.82rem; font-weight: 500; padding: 0.25rem 0.6rem; border-radius: 999px; white-space: nowrap; }
nav.toc a:hover{ color: var(--text); background: var(--surface-2); text-decoration: none; }

/* Hero / overview */
.hero{ padding: clamp(1.5rem,4vw,2.5rem) 0 0.5rem; }
.hero .why{ font-family: "Space Grotesk", sans-serif; font-weight: 600; font-size: clamp(1.25rem, 3.2vw, 1.7rem);
  letter-spacing: -0.015em; line-height: 1.25; max-width: 30ch; margin: 0 0 1.25rem; }
.stats{ display: flex; flex-wrap: wrap; gap: 0.5rem 2.25rem; align-items: baseline; }
.stat{ display: flex; flex-direction: column; gap: 0.1rem; }
.stat .n{ font-family: "Space Grotesk", sans-serif; font-weight: 700; font-size: 1.5rem; line-height: 1; letter-spacing: -0.02em; }
.stat .k{ color: var(--muted); font-size: 0.72rem; text-transform: uppercase; letter-spacing: 0.06em; }
.stat .sub{ color: var(--muted); font-size: 0.78rem; }
.stat.crit .n{ color: var(--crit); } .stat.high .n{ color: var(--high); }

/* Sections */
section{ padding: 1.5rem 0; border-top: 1px solid var(--line); scroll-margin-top: 96px; }
section h2{ font-family: "Space Grotesk", sans-serif; font-weight: 600; font-size: 1.05rem; letter-spacing: -0.01em;
  margin: 0 0 0.9rem; display: flex; align-items: center; gap: 0.45rem; }
section h2 .count{ color: var(--muted); font-weight: 500; font-size: 0.85rem; }
h3.repo{ font-family: "IBM Plex Sans", sans-serif; font-weight: 600; font-size: 0.92rem; margin: 1.1rem 0 0.4rem; }

/* Callout (needs attention) */
.callout{ border: 1px solid var(--border); border-left: 3px solid var(--warn); background: var(--warn-bg);
  border-radius: 8px; padding: 0.9rem 1.1rem; }
.callout.ok{ border-left-color: var(--ok); background: var(--ok-bg); }
.callout h2{ margin-top: 0; }
.callout ul{ margin: 0.4rem 0 0; padding-left: 1.1rem; } .callout li{ margin: 0.2rem 0; }

/* Tables */
.tbl{ overflow-x: auto; border: 1px solid var(--border); border-radius: 8px; }
table{ border-collapse: collapse; width: 100%; font-size: 0.88rem; }
thead th{ position: sticky; top: 0; background: var(--surface-2); text-align: left; font-weight: 600; color: var(--muted);
  font-size: 0.72rem; text-transform: uppercase; letter-spacing: 0.05em; padding: 0.55rem 0.7rem; border-bottom: 1px solid var(--border); }
td{ padding: 0.5rem 0.7rem; border-bottom: 1px solid var(--line); vertical-align: top; }
tbody tr:last-child td{ border-bottom: 0; }
tbody tr:nth-child(even){ background: color-mix(in oklch, var(--surface-2) 45%, transparent); }
td.num{ text-align: right; font-variant-numeric: tabular-nums; font-family: "IBM Plex Mono", monospace; }
td.num.zero{ color: var(--muted); opacity: 0.45; }

/* Pills */
.pill{ display: inline-block; padding: 0.05rem 0.45rem; border-radius: 999px; font-size: 0.72rem; font-weight: 600;
  letter-spacing: 0.02em; line-height: 1.5; }
.sev-critical{ background: var(--crit-bg); color: var(--crit); }
.sev-high{ background: var(--high-bg); color: var(--high); }
.sev-medium{ background: var(--med-bg); color: var(--med); }
.sev-low{ background: var(--low-bg); color: var(--low); }
.v-good{ background: var(--ok-bg); color: var(--ok); }
.v-bad{ background: var(--crit-bg); color: var(--crit); }
.v-needs_human_verification{ background: var(--warn-bg); color: var(--warn); }
.src{ color: var(--muted); font-size: 0.78rem; }
.status-errored{ color: var(--warn); font-weight: 600; }
.status-clean{ color: var(--muted); }

/* Lists */
ul.flat{ list-style: none; margin: 0; padding: 0; }
ul.flat > li{ padding: 0.45rem 0; border-bottom: 1px solid var(--line); }
ul.flat > li:last-child{ border-bottom: 0; }
.empty{ color: var(--muted); font-style: italic; }
.summary{ color: var(--muted); font-size: 0.85rem; margin: 0.15rem 0 0; }
details.trace{ margin-top: 0.35rem; } details.trace summary{ cursor: pointer; color: var(--muted); font-size: 0.78rem; }
details.trace ul{ margin: 0.3rem 0 0; padding-left: 1.1rem; } details.trace li{ color: var(--muted); font-size: 0.8rem; }
blockquote{ border-left: 3px solid var(--border); margin: 0; padding: 0.25rem 1rem; color: var(--muted); font-style: italic; }

footer{ border-top: 1px solid var(--line); padding: 1.5rem 0 3rem; color: var(--muted); font-size: 0.85rem;
  display: flex; gap: 1rem; flex-wrap: wrap; }
</style>
</head>
<body>
<header class="mast"><div class="wrap">
<span class="brand">Kairos <span class="dot">Security</span></span>
<span class="grow"></span>
<span class="upd">Updated {{.GeneratedAt}}</span>
{{if not .AIAvailable}}<span class="ai-warn">&#9888;&#65039; AI unavailable this run</span>{{end}}
</div></header>

<nav class="toc"><div class="wrap">
<a href="#attention">Attention</a>
<a href="#findings">Findings</a>
<a href="#prs">Open PRs</a>
<a href="#reviews">PR reviews</a>
<a href="#ledger">Ledger</a>
<a href="#fronts">Waterfall</a>
<a href="#run">This run</a>
</div></nav>

<main class="wrap">

<div class="hero">
<p class="why">{{.Activity.Why}}</p>
<div class="stats">
<div class="stat"><span class="n">{{.Activity.Repos}}</span><span class="k">Repos scanned</span>{{if or (gt .Activity.Skipped 0) (gt .Activity.Errored 0)}}<span class="sub">{{if gt .Activity.Skipped 0}}{{.Activity.Skipped}} skipped{{end}}{{if gt .Activity.Errored 0}}{{if gt .Activity.Skipped 0}} &middot; {{end}}{{.Activity.Errored}} errored{{end}}</span>{{end}}</div>
<div class="stat{{if gt .Activity.Crit 0}} crit{{else if gt .Activity.High 0}} high{{end}}"><span class="n">{{.Activity.Findings}}</span><span class="k">CVE findings</span><span class="sub">{{.Activity.Crit}} crit &middot; {{.Activity.High}} high &middot; {{.Activity.Med}} med</span></div>
<div class="stat"><span class="n">{{.Activity.PRs}}</span><span class="k">CVE-related PRs</span>{{if .Activity.SourceBreakdown}}<span class="sub">{{.Activity.SourceBreakdown}}</span>{{end}}</div>
<div class="stat"><span class="n">{{.Activity.LedgerOpen}}</span><span class="k">Bot PRs open</span><span class="sub">{{.Activity.Merged}} merged &middot; {{.Activity.NeedsHuman}} need human</span></div>
</div>
</div>

<section id="attention">
<h2>&#128680; Needs attention</h2>
{{if or .NeedsHuman .CollectErrors .Focus}}
{{if .NeedsHuman}}<div class="callout"><strong>{{len .NeedsHuman}} item{{if ne (len .NeedsHuman) 1}}s{{end}} need a human</strong>
<ul>{{range .NeedsHuman}}<li>{{.}}</li>{{end}}</ul></div>{{end}}
{{if .Focus}}<h3 class="repo">&#128293; Focus now</h3>
<ul class="flat">{{range .Focus}}<li>{{if .URL}}<a href="{{.URL}}" target="_blank" rel="noopener">{{.Title}}</a>{{else}}<strong>{{.Title}}</strong>{{end}}{{if .Summary}} &mdash; {{.Summary}}{{end}}</li>{{end}}</ul>{{end}}
{{if .CollectErrors}}<h3 class="repo">&#9888;&#65039; {{len .CollectErrors}} collection error{{if ne (len .CollectErrors) 1}}s{{end}}</h3>
<ul class="flat">{{range .CollectErrors}}<li><a href="{{repoURL .Repo}}" target="_blank" rel="noopener">{{.Repo}}</a> <span class="src">{{.Collector}}</span> &mdash; {{.Message}}</li>{{end}}</ul>{{end}}
{{else}}<div class="callout ok">All clear &mdash; no findings need a human, no scan errors.</div>{{end}}
</section>

<section id="findings">
<h2>&#128230; Per-repo findings <span class="count">{{len .Repos}} repos</span></h2>
<div class="tbl"><table>
<thead><tr><th>Repo</th><th>Crit</th><th>High</th><th>Med</th><th>Total</th><th>Status</th></tr></thead>
<tbody>
{{range .Repos}}<tr><td><a href="{{repoURL .Repo}}" target="_blank" rel="noopener">{{.Repo}}</a></td>
<td class="num{{if not .Crit}} zero{{end}}">{{.Crit}}</td><td class="num{{if not .High}} zero{{end}}">{{.High}}</td><td class="num{{if not .Med}} zero{{end}}">{{.Med}}</td><td class="num{{if not .Total}} zero{{end}}">{{.Total}}</td>
<td class="{{.StatusClass}}">{{.Status}}</td></tr>
{{else}}<tr><td colspan="6" class="empty">No repos tracked.</td></tr>
{{end}}</tbody>
</table></div>
</section>

{{if .Components}}<section id="hadron-components">
<h2>&#129513; Hadron component CVEs <span class="count">{{len .Components}}</span></h2>
<div class="tbl"><table>
<thead><tr><th>Package</th><th>Current</th><th>Fixed</th><th>Severity</th><th>CVE</th></tr></thead>
<tbody>
{{range .Components}}<tr><td>{{.Package}}</td><td>{{.Current}}</td><td>{{.Fixed}}</td><td class="{{sevClass .Severity}}">{{.Severity}}</td>
<td>{{if .URL}}<a href="{{.URL}}" target="_blank" rel="noopener">{{.Title}}</a>{{else}}{{.Title}}{{end}}</td></tr>
{{end}}</tbody>
</table></div>
</section>

{{end}}{{if .Informational}}<section id="informational">
<h2>&#8505;&#65039; Informational &mdash; not counted <span class="count">{{len .Informational}}</span></h2>
<p class="summary">These findings are separated from the counts above: CVEs we are already past, or components accepted as pinned risk.</p>
<div class="tbl"><table>
<thead><tr><th>Package</th><th>Current</th><th>Fixed</th><th>Severity</th><th>CVE</th><th>Why</th></tr></thead>
<tbody>
{{range .Informational}}<tr><td>{{.Package}}</td><td>{{.Current}}</td><td>{{.Fixed}}</td><td class="{{sevClass .Severity}}">{{.Severity}}</td>
<td>{{if .URL}}<a href="{{.URL}}" target="_blank" rel="noopener">{{.Title}}</a>{{else}}{{.Title}}{{end}}</td><td>{{.Why}}</td></tr>
{{end}}</tbody>
</table></div>
</section>

{{end}}<section id="prs">
<h2>&#128203; Open PRs <span class="count">CVE-related</span></h2>
{{if .OpenPRs}}{{range .OpenPRs}}<h3 class="repo"><a href="{{repoURL .Repo}}" target="_blank" rel="noopener">{{.Repo}}</a></h3>
<ul class="flat">
{{range .PRs}}<li>{{if .URL}}<a href="{{.URL}}" target="_blank" rel="noopener">#{{.Number}} {{.Title}}</a>{{else}}#{{.Number}} {{.Title}}{{end}} <span class="pill src">{{.Source}}</span> &mdash; {{.Status}}</li>
{{end}}</ul>
{{end}}{{else}}<p class="empty">No CVE-related PRs open.</p>{{end}}
</section>

{{if .Reviews}}<section id="reviews">
<h2>&#128270; Bot-PR reviews</h2>
{{range .Reviews}}<h3 class="repo"><a href="{{repoURL .Repo}}" target="_blank" rel="noopener">{{.Repo}}</a></h3>
<ul class="flat">
{{range .Reviews}}<li>{{if .URL}}<a href="{{.URL}}" target="_blank" rel="noopener">#{{.PR}}</a>{{else}}#{{.PR}}{{end}} <span class="pill v-{{.Verdict}}">{{.Icon}} {{.Verdict}}</span> &mdash; {{.Reasoning}}
{{if .ChangesSummary}}<p class="summary">&#8627; {{.ChangesSummary}}</p>{{end}}
{{if .Trace}}<details class="trace"><summary>review trace</summary><ul>{{range .Trace}}<li>{{.}}</li>{{end}}</ul></details>{{end}}</li>
{{end}}</ul>
{{end}}</section>{{else}}<section id="reviews"><h2>&#128270; Bot-PR reviews</h2><p class="empty">No bot PRs reviewed yet.</p></section>{{end}}

<section id="ledger">
<h2>&#129302; Bot PR ledger</h2>
{{if .Ledger}}<div class="tbl"><table>
<thead><tr><th>Repo</th><th>Bump</th><th>Kind</th><th>Source</th><th>State</th><th>PR</th></tr></thead>
<tbody>
{{range .Ledger}}<tr><td><a href="{{repoURL .Repo}}" target="_blank" rel="noopener">{{.Repo}}</a></td><td><code>{{.Bump}}</code></td><td>{{.Kind}}</td><td><span class="src">{{.Source}}</span></td><td>{{.State}}</td><td>{{if gt .PRNumber 0}}<a href="{{.PRURL}}" target="_blank" rel="noopener">#{{.PRNumber}}</a>{{else}}&mdash;{{end}}</td></tr>
{{end}}</tbody>
</table></div>{{else}}<p class="empty">No bot PRs yet.</p>{{end}}
</section>

<section id="fronts">
<h2>&#127754; Waterfall fronts <span class="count">shared root causes</span></h2>
{{if .Waterfall}}<div class="tbl"><table>
<thead><tr><th>Root cause</th><th>Severity</th><th>Suggested bump</th><th>Affected repos</th></tr></thead>
<tbody>
{{range .Waterfall}}<tr><td>{{.RootCause}}</td><td><span class="pill {{sevClass .Severity}}">{{.Severity}}</span></td><td><code>{{.SuggestedBump.Package}}@{{.SuggestedBump.To}}</code></td><td>{{join .AffectedRepos}}</td></tr>
{{end}}</tbody>
</table></div>{{else}}<p class="empty">None.</p>{{end}}
</section>

<section id="run">
<h2>&#128203; This run</h2>
{{if .Narrative}}<blockquote>{{.Narrative}}</blockquote>{{end}}
<ul class="flat">
<li><strong>Scanned</strong> {{.Activity.Repos}} repos{{if gt .Activity.Skipped 0}} &middot; {{.Activity.Skipped}} skipped{{end}}{{if gt .Activity.Errored 0}} &middot; {{.Activity.Errored}} errored{{end}}</li>
<li><strong>Findings</strong> {{.Activity.Findings}} &mdash; {{.Activity.Crit}} critical / {{.Activity.High}} high / {{.Activity.Med}} medium / {{.Activity.Low}} low / {{.Activity.Unknown}} unknown</li>
<li><strong>Remediation</strong> {{.Activity.LedgerOpen}} open &middot; {{.Activity.Superseded}} superseded &middot; {{.Activity.Merged}} merged &middot; {{.Activity.NeedsHuman}} need-human</li>
</ul>
{{if .CoordinationSummary}}<p class="summary">{{.CoordinationSummary}}</p>{{end}}
</section>

</main>
<footer class="wrap">
<a href="https://github.com/kairos-io/security" target="_blank" rel="noopener">kairos-io/security</a>
{{if .RunURL}}<a href="{{.RunURL}}" target="_blank" rel="noopener">Run log</a>{{end}}
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
			Repo: r.repo, Crit: r.crit, High: r.high, Med: r.med, Total: r.total,
			Status: repoStatus(r), StatusClass: repoStatusClass(r),
		})
	}
	var components []htmlComponentRow
	for _, f := range hadronComponentRows(in.Correlated.Findings) {
		fixed := f.FixedVersion
		if fixed == "" {
			fixed = "—"
		}
		title, url := focusTitleURL(f)
		components = append(components, htmlComponentRow{
			Package: f.Package, Current: f.CurrentVersion, Fixed: fixed,
			Severity: f.Severity, Title: title, URL: url,
		})
	}
	var informational []htmlInfoRow
	for _, f := range informationalFindings(in.Correlated.Findings) {
		fixed := f.FixedVersion
		if fixed == "" {
			fixed = "—"
		}
		title, url := focusTitleURL(f)
		informational = append(informational, htmlInfoRow{
			Package: f.Package, Current: f.CurrentVersion, Fixed: fixed,
			Severity: f.Severity, Title: title, URL: url, Why: f.ClassReason,
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
			PR: r.PR, URL: r.URL, Icon: verdictIcon(r.Verdict), Verdict: r.Verdict, Reasoning: r.Reasoning, ChangesSummary: r.ChangesSummary, Trace: r.Trace,
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
		Components:          components,
		Informational:       informational,
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
