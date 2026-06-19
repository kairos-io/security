package collect

import (
	"encoding/json"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

type ImageCVE struct {
	Runner func(ref string) ([]byte, error)
}

func (ImageCVE) Name() string { return "imageCVE" }

type trivyReport struct {
	Results []struct {
		Vulnerabilities []struct {
			VulnerabilityID  string `json:"VulnerabilityID"`
			PkgName          string `json:"PkgName"`
			InstalledVersion string `json:"InstalledVersion"`
			FixedVersion     string `json:"FixedVersion"`
			Severity         string `json:"Severity"`
			PrimaryURL       string `json:"PrimaryURL"`
			Title            string `json:"Title"`
		} `json:"Vulnerabilities"`
	} `json:"Results"`
}

func (c ImageCVE) Collect(repo state.Repo) ([]state.Finding, error) {
	out := map[string]state.Finding{}
	for _, art := range repo.Artifacts {
		if art.Type != "image" {
			continue
		}
		raw, err := c.Runner(art.Ref)
		if err != nil {
			return nil, err
		}
		var rep trivyReport
		if err := json.Unmarshal(raw, &rep); err != nil {
			return nil, err
		}
		for _, res := range rep.Results {
			for _, v := range res.Vulnerabilities {
				f := state.Finding{
					ID:             FindingID(repo.Repo, "imageCVE", v.VulnerabilityID, v.PkgName),
					Repo:           repo.Repo,
					Type:           "imageCVE",
					CVEID:          v.VulnerabilityID,
					Package:        v.PkgName,
					CurrentVersion: v.InstalledVersion,
					FixedVersion:   v.FixedVersion,
					Severity:       strings.ToLower(v.Severity),
					Source:         "trivy",
					Title:          v.Title,
					URL:            v.PrimaryURL,
					FirstSeen:      Today(),
					LastSeen:       Today(),
				}
				out[f.ID] = f
			}
		}
	}
	res := make([]state.Finding, 0, len(out))
	for _, f := range out {
		res = append(res, f)
	}
	return res, nil
}
