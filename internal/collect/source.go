package collect

import (
	"bufio"
	"bytes"
	"encoding/json"

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
	} `json:"osv"`
	Finding *struct {
		OSV   string `json:"osv"`
		Trace []struct {
			Module  string `json:"module"`
			Version string `json:"version"`
		} `json:"trace"`
	} `json:"finding"`
}

func (c SourceCVE) Collect(repo state.Repo) ([]state.Finding, error) {
	raw, err := c.Runner(repo)
	if err != nil {
		return nil, err
	}
	type adv struct {
		cve, fixed, summary string
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
			Severity:       "unknown",
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
