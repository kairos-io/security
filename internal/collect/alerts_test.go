package collect

import (
	"errors"
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGHAlertsCollector(t *testing.T) {
	nowFn = func() string { return "2026-06-19" }
	defer func() { nowFn = defaultNow }()

	gh := ghclient.NewFake()
	gh.Alerts["kairos-io/immucore"] = []ghclient.Alert{
		{Number: 7, CVEID: "CVE-2025-1234", GHSA: "GHSA-aaa", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high", URL: "u", FixedVersion: "0.33.0"},
	}
	c := GHAlerts{GH: gh}
	fs, err := c.Collect(state.Repo{Repo: "kairos-io/immucore"})
	require.NoError(t, err)
	require.Len(t, fs, 1)
	assert.Equal(t, "ghAlert", fs[0].Type)
	assert.Equal(t, "CVE-2025-1234", fs[0].CVEID)
	assert.Equal(t, "high", fs[0].Severity)
	assert.Equal(t, "0.33.0", fs[0].FixedVersion)
}

// A dashboard that could not read the alerts must not look like one that read
// them and found none. The refusal has to reach findings.json's errors, which
// is what the dashboard's "collection errors" block renders.
func TestGHAlertsReportsARefusalAsACollectionError(t *testing.T) {
	gh := ghclient.NewFake()
	gh.AlertsErr["kairos-io/AuroraBoot"] = errors.New(
		"cannot read Dependabot alerts: the token needs the security_events scope")

	out := Run(
		[]state.Repo{{Repo: "kairos-io/AuroraBoot"}},
		[]Collector{GHAlerts{GH: gh}},
		state.Findings{},
	)

	assert.Empty(t, out.Findings)
	require.Len(t, out.Errors, 1)
	assert.Equal(t, "ghAlerts", out.Errors[0].Collector)
	assert.Equal(t, "kairos-io/AuroraBoot", out.Errors[0].Repo)
	assert.Contains(t, out.Errors[0].Message, "security_events")
}
