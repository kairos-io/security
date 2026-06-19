package correlate

import (
	"fmt"
	"sort"

	"github.com/kairos-io/security/internal/state"
)

var sevRank = map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1, "unknown": 0, "": 0}

func worse(a, b string) string {
	if sevRank[a] >= sevRank[b] {
		return a
	}
	return b
}

func Run(in state.Findings) state.Correlated {
	// 1) dedupe by (repo, cveID, package); PR findings (no CVE) pass through.
	merged := map[string]state.Finding{}
	var order []string
	for _, f := range in.Findings {
		key := f.Repo + "|" + f.CVEID + "|" + f.Package
		if f.CVEID == "" {
			key = f.ID // PRs and CVE-less findings never merge
		}
		cur, ok := merged[key]
		if !ok {
			merged[key] = f
			order = append(order, key)
			continue
		}
		cur.Severity = worse(cur.Severity, f.Severity)
		if cur.FixedVersion == "" {
			cur.FixedVersion = f.FixedVersion
		}
		if cur.FirstSeen == "" || (f.FirstSeen != "" && f.FirstSeen < cur.FirstSeen) {
			cur.FirstSeen = f.FirstSeen
		}
		merged[key] = cur
	}

	findings := make([]state.Finding, 0, len(merged))
	for _, k := range order {
		findings = append(findings, merged[k])
	}
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })

	// 2) waterfall: group go-ecosystem CVEs by (cveID, package) across repos.
	type agg struct {
		repos    map[string]bool
		severity string
		fixed    string
	}
	groups := map[string]*agg{}
	for _, f := range findings {
		if f.Ecosystem != "go" || f.CVEID == "" || f.Package == "" {
			continue
		}
		gk := f.CVEID + "|" + f.Package
		g := groups[gk]
		if g == nil {
			g = &agg{repos: map[string]bool{}}
			groups[gk] = g
		}
		g.repos[f.Repo] = true
		g.severity = worse(g.severity, f.Severity)
		if g.fixed == "" {
			g.fixed = f.FixedVersion
		}
	}

	var waterfall []state.WaterfallGroup
	for gk, g := range groups {
		if len(g.repos) < 2 {
			continue
		}
		repos := make([]string, 0, len(g.repos))
		for r := range g.repos {
			repos = append(repos, r)
		}
		sort.Strings(repos)
		cve, pkg := splitKey(gk)
		waterfall = append(waterfall, state.WaterfallGroup{
			ID:            "go-" + cve + "-" + pkg,
			RootCause:     fmt.Sprintf("%s (%s)", pkg, cve),
			Ecosystem:     "go",
			Severity:      g.severity,
			AffectedRepos: repos,
			SuggestedBump: state.Bump{Package: pkg, To: g.fixed},
		})
	}
	sort.Slice(waterfall, func(i, j int) bool { return waterfall[i].ID < waterfall[j].ID })

	return state.Correlated{Findings: findings, Waterfall: waterfall}
}

func splitKey(k string) (cve, pkg string) {
	for i := 0; i < len(k); i++ {
		if k[i] == '|' {
			return k[:i], k[i+1:]
		}
	}
	return k, ""
}
