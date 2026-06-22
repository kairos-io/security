package collect

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

type SourceCVE struct {
	Runner func(repo state.Repo) ([]byte, error)
}

func (SourceCVE) Name() string { return "sourceCVE" }

type govulnLine struct {
	OSV *struct {
		ID       string   `json:"id"`
		Aliases  []string `json:"aliases"`
		Summary  string   `json:"summary"`
		Affected []struct {
			Package struct {
				Name string `json:"name"`
			} `json:"package"`
			Ranges []struct {
				Events []struct {
					Fixed string `json:"fixed"`
				} `json:"events"`
			} `json:"ranges"`
		} `json:"affected"`
		DatabaseSpecific *struct {
			Severity string `json:"severity"`
		} `json:"database_specific"`
	} `json:"osv"`
	Finding *struct {
		OSV   string `json:"osv"`
		Trace []struct {
			Module   string `json:"module"`
			Version  string `json:"version"`
			Function string `json:"function"`
		} `json:"trace"`
	} `json:"finding"`
}

func (c SourceCVE) Collect(repo state.Repo) ([]state.Finding, error) {
	raw, err := c.Runner(repo)
	if err != nil {
		return nil, err
	}
	type adv struct {
		cve, fixed, summary, severity string
	}
	advisories := map[string]adv{}
	var findings []govulnLine

	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		var gl govulnLine
		if err := json.Unmarshal(line, &gl); err != nil {
			continue // tolerate non-JSON progress lines
		}
		if gl.OSV != nil {
			a := adv{summary: gl.OSV.Summary}
			if gl.OSV.DatabaseSpecific != nil {
				a.severity = gl.OSV.DatabaseSpecific.Severity
			}
			for _, al := range gl.OSV.Aliases {
				if len(al) > 3 && al[:3] == "CVE" {
					a.cve = al
				}
			}
			for _, af := range gl.OSV.Affected {
				for _, rg := range af.Ranges {
					for _, ev := range rg.Events {
						if ev.Fixed != "" {
							a.fixed = ev.Fixed
						}
					}
				}
			}
			advisories[gl.OSV.ID] = a
		}
		if gl.Finding != nil {
			findings = append(findings, gl)
		}
	}

	out := map[string]state.Finding{}
	for _, gl := range findings {
		if len(gl.Finding.Trace) == 0 {
			continue
		}
		if gl.Finding.Trace[0].Function == "" {
			continue // not reachable: required but not called
		}
		t := gl.Finding.Trace[0]
		a := advisories[gl.Finding.OSV]
		cve := a.cve
		if cve == "" {
			cve = gl.Finding.OSV
		}
		f := state.Finding{
			ID:             FindingID(repo.Repo, "sourceCVE", cve, t.Module),
			Repo:           repo.Repo,
			Type:           "sourceCVE",
			CVEID:          cve,
			Ecosystem:      "go",
			Package:        t.Module,
			CurrentVersion: t.Version,
			FixedVersion:   a.fixed,
			Severity:       severityFromOSV(a.severity),
			Source:         "govulncheck",
			Title:          a.summary,
			FirstSeen:      Today(),
			LastSeen:       Today(),
		}
		out[f.ID] = f // dedupe module-level findings
	}
	res := make([]state.Finding, 0, len(out))
	for _, f := range out {
		res = append(res, f)
	}
	return res, nil
}

func severityFromOSV(s string) string {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "CRITICAL":
		return "critical"
	case "HIGH":
		return "high"
	case "MODERATE", "MEDIUM":
		return "medium"
	case "LOW":
		return "low"
	default:
		return "high" // reachable vuln with no severity data is actionable
	}
}
