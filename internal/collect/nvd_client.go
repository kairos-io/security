package collect

import (
	"encoding/json"
	"strings"
)

// NVDQueryFunc performs a single NVD CPE-match query and returns the raw
// JSON response body. Real implementations GET
// https://services.nvd.nist.gov/rest/json/cves/2.0?cpeName=...; tests inject
// a fixture-returning fake.
type NVDQueryFunc func(vendor, product, version string) ([]byte, error)

type nvdCPEMatch struct {
	Vulnerable          bool   `json:"vulnerable"`
	VersionEndExcluding string `json:"versionEndExcluding"`
}

type nvdVulnerability struct {
	CVE struct {
		ID           string `json:"id"`
		Descriptions []struct {
			Lang  string `json:"lang"`
			Value string `json:"value"`
		} `json:"descriptions"`
		Metrics struct {
			CvssMetricV31 []struct {
				CvssData struct {
					BaseSeverity string `json:"baseSeverity"`
				} `json:"cvssData"`
			} `json:"cvssMetricV31"`
		} `json:"metrics"`
		Configurations []struct {
			Nodes []struct {
				CPEMatch []nvdCPEMatch `json:"cpeMatch"`
			} `json:"nodes"`
		} `json:"configurations"`
	} `json:"cve"`
}

type nvdResponse struct {
	Vulnerabilities []nvdVulnerability `json:"vulnerabilities"`
}

// NVDResult is one CVE hit from an NVD CPE-match query. VersionEndExcluding
// is the vulnerability boundary reported by NVD — a *candidate* fixed
// version, not yet validated against any package's update constraints (that
// validation is Plan 5b's job, done immediately before opening a bump PR).
type NVDResult struct {
	CVEID               string
	Severity            string
	VersionEndExcluding string
	Title               string
	URL                 string
	// Details is the English description (same source as Title). Kept
	// separate so the finding's Title can stay short while Details carries the
	// full text for the applicability classifier.
	Details string
	// AffectedRanges is the CVE's `configurations[]` re-marshalled to a JSON
	// string. Feeds the applicability classifier alongside Details.
	AffectedRanges string
}

func QueryNVD(query NVDQueryFunc, vendor, product, version string) ([]NVDResult, error) {
	raw, err := query(vendor, product, version)
	if err != nil {
		return nil, err
	}
	var resp nvdResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	out := make([]NVDResult, 0, len(resp.Vulnerabilities))
	for _, v := range resp.Vulnerabilities {
		title := ""
		for _, d := range v.CVE.Descriptions {
			if d.Lang == "en" {
				title = d.Value
				break
			}
		}
		sev := "unknown"
		if len(v.CVE.Metrics.CvssMetricV31) > 0 {
			sev = strings.ToLower(v.CVE.Metrics.CvssMetricV31[0].CvssData.BaseSeverity)
		}
		boundary := ""
		for _, c := range v.CVE.Configurations {
			for _, n := range c.Nodes {
				for _, m := range n.CPEMatch {
					if m.Vulnerable && m.VersionEndExcluding != "" {
						boundary = m.VersionEndExcluding
					}
				}
			}
		}
		affected := ""
		if len(v.CVE.Configurations) > 0 {
			if b, mErr := json.Marshal(v.CVE.Configurations); mErr == nil {
				affected = string(b)
			}
		}
		out = append(out, NVDResult{
			CVEID:               v.CVE.ID,
			Severity:            sev,
			VersionEndExcluding: boundary,
			Title:               title,
			URL:                 "https://nvd.nist.gov/vuln/detail/" + v.CVE.ID,
			Details:             title,
			AffectedRanges:      affected,
		})
	}
	return out, nil
}
