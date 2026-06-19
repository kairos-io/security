package collect

import (
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

type GHAlerts struct {
	GH ghclient.GitHub
}

func (GHAlerts) Name() string { return "ghAlerts" }

func (c GHAlerts) Collect(repo state.Repo) ([]state.Finding, error) {
	alerts, err := c.GH.ListDependabotAlerts(repo.Repo)
	if err != nil {
		return nil, err
	}
	out := make([]state.Finding, 0, len(alerts))
	for _, a := range alerts {
		cve := a.CVEID
		if cve == "" {
			cve = a.GHSA
		}
		out = append(out, state.Finding{
			ID:           FindingID(repo.Repo, "ghAlert", cve, a.Package),
			Repo:         repo.Repo,
			Type:         "ghAlert",
			CVEID:        a.CVEID,
			GHSA:         a.GHSA,
			Ecosystem:    strings.ToLower(a.Ecosystem),
			Package:      a.Package,
			FixedVersion: a.FixedVersion,
			Severity:     strings.ToLower(a.Severity),
			Source:       "dependabot",
			URL:          a.URL,
			FirstSeen:    Today(),
			LastSeen:     Today(),
		})
	}
	return out, nil
}
