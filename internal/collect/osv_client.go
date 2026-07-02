package collect

import (
	"encoding/json"
	"strconv"
	"strings"
)

// OSVQueryFunc performs a single OSV.dev query and returns the raw JSON
// response body. Real implementations POST to https://api.osv.dev/v1/query;
// tests inject a fixture-returning fake.
type OSVQueryFunc func(ecosystem, pkg, version string) ([]byte, error)

type osvVuln struct {
	ID               string   `json:"id"`
	Aliases          []string `json:"aliases"`
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
				Fixed string `json:"fixed"`
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
		fixed := ""
		for _, a := range v.Affected {
			for _, rg := range a.Ranges {
				for _, ev := range rg.Events {
					if ev.Fixed != "" {
						fixed = stripAlpineRevisionSuffix(ev.Fixed)
					}
				}
			}
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
