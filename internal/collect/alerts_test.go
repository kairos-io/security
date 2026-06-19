package collect

import (
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
