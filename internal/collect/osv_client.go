package collect

import (
	"encoding/json"
	"strconv"
	"strings"

	ver "github.com/kairos-io/security/internal/version"
)

// OSVQueryFunc performs a single OSV.dev query and returns the raw JSON
// response body. Real implementations POST to https://api.osv.dev/v1/query;
// tests inject a fixture-returning fake.
type OSVQueryFunc func(ecosystem, pkg, version string) ([]byte, error)

type osvVuln struct {
	ID               string   `json:"id"`
	Aliases          []string `json:"aliases"`
	Upstream         []string `json:"upstream"`
	Summary          string   `json:"summary"`
	DatabaseSpecific struct {
		Severity string `json:"severity"`
	} `json:"database_specific"`
	Severity []struct {
		Type  string `json:"type"`
		Score string `json:"score"`
	} `json:"severity"`
	Affected []struct {
		Ranges []struct {
			Events []struct {
				Introduced string `json:"introduced"`
				Fixed      string `json:"fixed"`
			} `json:"events"`
		} `json:"ranges"`
	} `json:"affected"`
}

type osvQueryResponse struct {
	Vulns []osvVuln `json:"vulns"`
}

// OSVResult is one CVE hit from an OSV.dev query, normalized to ksec's finding shape.
type OSVResult struct {
	CVEID        string
	Severity     string // critical|high|medium|low|unknown, via osvSeverity
	FixedVersion string
	Title        string
	URL          string
}

// QueryOSV queries OSV.dev for (ecosystem, pkg, version) and normalizes every
// hit. For the Alpine ecosystem, OSV's "fixed" version carries Alpine's own
// package-revision suffix (e.g. "3.6.4-r0"); stripAlpineRevisionSuffix
// removes it so the result approximates hadron's upstream version numbering.
func QueryOSV(query OSVQueryFunc, ecosystem, pkg, version string) ([]OSVResult, error) {
	raw, err := query(ecosystem, pkg, version)
	if err != nil {
		return nil, err
	}
	var resp osvQueryResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	out := make([]OSVResult, 0, len(resp.Vulns))
	for _, v := range resp.Vulns {
		cve := v.ID
		for _, alias := range v.Aliases {
			if strings.HasPrefix(alias, "CVE-") {
				cve = alias
			}
		}
		// Alpine OSV-converted advisories carry the real CVE id in "upstream"
		// (and never populate "aliases"); honor it with the same last-match-wins
		// semantics as the alias loop above.
		for _, up := range v.Upstream {
			if strings.HasPrefix(up, "CVE-") {
				cve = up
			}
		}
		fixed, applicable := osvApplicableFix(v, ver.Compare, version)
		if !applicable {
			continue // version is below every introduced boundary — not vulnerable
		}
		out = append(out, OSVResult{
			CVEID:        cve,
			Severity:     osvSeverity(v),
			FixedVersion: fixed,
			Title:        v.Summary,
			URL:          "https://osv.dev/vulnerability/" + v.ID,
		})
	}
	return out, nil
}

// osvApplicableFix picks the fix version for the range that applies to the
// queried version v, and reports whether v is at/after some range's
// "introduced" boundary (i.e. potentially vulnerable). Alpine revision suffixes
// are stripped for comparison. When v sits inside a range (introduced <= v <
// fixed) that range's fix is returned; when v is past a range's fix, the
// largest such fix is returned (already-fixed, classified later). When v is
// below every introduced boundary, applicable=false and the vuln is dropped.
func osvApplicableFix(v osvVuln, cmp func(string, string) int, queried string) (string, bool) {
	q := stripAlpineRevisionSuffix(queried)
	parseable := looksNumeric(q)
	bestFix := ""
	applicable := false
	for _, a := range v.Affected {
		for _, rg := range a.Ranges {
			introduced := "0"
			fixed := ""
			for _, ev := range rg.Events {
				if ev.Introduced != "" {
					introduced = ev.Introduced
				}
				if ev.Fixed != "" {
					fixed = stripAlpineRevisionSuffix(ev.Fixed)
				}
			}
			// Only trust the introduced-boundary drop when q is orderable. A
			// non-numeric version can't be compared meaningfully, so we fail
			// OPEN (keep the vuln visible) instead of dropping it below "0".
			if parseable && cmp(q, stripAlpineRevisionSuffix(introduced)) < 0 {
				continue // below this range's introduced boundary
			}
			applicable = true
			if fixed == "" {
				continue
			}
			// Inside the range (q < fixed): this is the fix we want.
			if cmp(q, fixed) < 0 {
				return fixed, true
			}
			// Past the fix: remember the largest fix seen (already-fixed).
			if bestFix == "" || cmp(fixed, bestFix) > 0 {
				bestFix = fixed
			}
		}
	}
	return bestFix, applicable
}

// looksNumeric reports whether v begins with a digit (so version.Compare can
// order it meaningfully). Non-numeric versions fail OPEN — we keep the vuln
// visible rather than silently dropping it.
func looksNumeric(v string) bool {
	return v != "" && v[0] >= '0' && v[0] <= '9'
}

// osvSeverity derives a Finding's severity from an OSV vuln record, in
// precedence order:
//  1. An explicit database_specific.severity label (CRITICAL/HIGH/MODERATE/
//     MEDIUM/LOW) is trusted directly — it's a human/tooling-assigned rating.
//  2. Otherwise the first parseable CVSS_V3 entry in the top-level severity
//     array is computed into a base score and mapped to a band.
//  3. Otherwise "unknown" — being honest about missing data beats guessing.
func osvSeverity(v osvVuln) string {
	switch strings.ToUpper(strings.TrimSpace(v.DatabaseSpecific.Severity)) {
	case "CRITICAL":
		return "critical"
	case "HIGH":
		return "high"
	case "MODERATE", "MEDIUM":
		return "medium"
	case "LOW":
		return "low"
	}
	for _, s := range v.Severity {
		if s.Type != "CVSS_V3" {
			continue
		}
		score, err := cvssV31BaseScore(s.Score)
		if err != nil {
			continue // try any remaining entries
		}
		return cvssSeverityLabel(score)
	}
	return "unknown"
}

// stripAlpineRevisionSuffix strips a trailing "-rN" Alpine package-revision
// suffix (e.g. "3.6.4-r0" -> "3.6.4"), leaving other version strings
// untouched.
func stripAlpineRevisionSuffix(v string) string {
	i := strings.LastIndex(v, "-r")
	if i <= 0 {
		return v
	}
	if _, err := strconv.Atoi(v[i+2:]); err != nil {
		return v
	}
	return v[:i]
}
