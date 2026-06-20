package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanOpensNewActionableTargetsDedupedAndCapped(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		// two CVEs in the same repo+package -> one target at the highest fixed version
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
		{ID: "b", Repo: "kairos-io/immucore", Type: "ghAlert", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.36.0", Severity: "critical"},
		// a different repo+package -> second target
		{ID: "c", Repo: "kairos-io/kairos", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/crypto", FixedVersion: "0.31.0", Severity: "high"},
		// not actionable: image CVE
		{ID: "d", Repo: "kairos-io/kairos", Type: "imageCVE", Package: "openssl", FixedVersion: "1.1.1w", Severity: "critical"},
		// not actionable: no fixed version
		{ID: "e", Repo: "kairos-io/kairos", Type: "sourceCVE", Ecosystem: "go", Package: "x/text", Severity: "low"},
	}}

	intents, deferred := Plan(c, state.Ledger{}, nil, 1) // cap to 1 new PR
	require.Len(t, intents, 1)
	assert.Equal(t, 1, deferred)
	in := intents[0]
	assert.Equal(t, IntentOpen, in.Type)
	// highest severity target first: immucore/x/net (critical) at the highest fixed version
	assert.Equal(t, "kairos-io/immucore|golang.org/x/net", in.Key)
	assert.Equal(t, "0.36.0", in.Bump.To)
	assert.Equal(t, "critical", in.Severity)
}

func TestPlanReconcilesExistingLedgerEntries(t *testing.T) {
	c := state.Correlated{}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", Repo: "kairos-io/immucore", State: "open"},
	}}
	intents, _ := Plan(c, led, nil, 10)
	require.Len(t, intents, 1)
	assert.Equal(t, IntentReconcile, intents[0].Type)
	require.NotNil(t, intents[0].Entry)
	assert.Equal(t, "open", intents[0].Entry.State)
}

func TestPlanSkipsTargetsAlreadyInLedger(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", State: "open"},
	}}
	intents, _ := Plan(c, led, nil, 10)
	// only the reconcile for the existing entry; no new open
	require.Len(t, intents, 1)
	assert.Equal(t, IntentReconcile, intents[0].Type)
}

// intentFor returns the first intent of the given type for a key, or nil.
func intentFor(intents []Intent, typ IntentType, key string) *Intent {
	for i := range intents {
		if intents[i].Type == typ && intents[i].Key == key {
			return &intents[i]
		}
	}
	return nil
}

// A non-live ledger state (planned, from a prior dry-run) must NOT permanently
// suppress the key: going live should re-open the PR while still reconciling.
func TestPlanReopensPlannedLedgerEntry(t *testing.T) {
	k := "kairos-io/immucore|golang.org/x/net"
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: k, Repo: "kairos-io/immucore", State: "planned", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	intents, _ := Plan(c, led, nil, 10)
	require.NotNil(t, intentFor(intents, IntentReconcile, k), "expected reconcile for the existing entry")
	open := intentFor(intents, IntentOpen, k)
	require.NotNil(t, open, "expected re-open for the planned entry")
	assert.Equal(t, "0.33.0", open.Bump.To)
}

// A transient build-failed entry must retry (re-open) on a later run.
func TestPlanReopensBuildFailedLedgerEntry(t *testing.T) {
	k := "kairos-io/immucore|golang.org/x/net"
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: k, Repo: "kairos-io/immucore", State: "build-failed", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	intents, _ := Plan(c, led, nil, 10)
	require.NotNil(t, intentFor(intents, IntentOpen, k), "expected re-open for the build-failed entry")
}

// An open entry has a live PR maintained via reconcile: no new open.
func TestPlanSkipsOpenLedgerEntryButReconciles(t *testing.T) {
	k := "kairos-io/immucore|golang.org/x/net"
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: k, Repo: "kairos-io/immucore", State: "open", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	intents, _ := Plan(c, led, nil, 10)
	require.NotNil(t, intentFor(intents, IntentReconcile, k))
	assert.Nil(t, intentFor(intents, IntentOpen, k), "open entry must not be re-opened")
}

// A merged entry re-opens only when a NEWER fixed version is later required.
func TestPlanReopensMergedOnlyForHigherVersion(t *testing.T) {
	k := "kairos-io/immucore|golang.org/x/net"

	// merged at 0.33.0, finding needs 0.36.0 -> re-open (re-bump).
	cHigher := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.36.0", Severity: "high"},
	}}
	ledLow := state.Ledger{Entries: []state.LedgerEntry{
		{Key: k, Repo: "kairos-io/immucore", State: "merged", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	intents, _ := Plan(cHigher, ledLow, nil, 10)
	open := intentFor(intents, IntentOpen, k)
	require.NotNil(t, open, "merged at lower version must re-open for the newer fix")
	assert.Equal(t, "0.36.0", open.Bump.To)

	// merged at 0.36.0, finding needs 0.33.0 -> skip (already addressed).
	cLower := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	ledHigh := state.Ledger{Entries: []state.LedgerEntry{
		{Key: k, Repo: "kairos-io/immucore", State: "merged", Bump: state.Bump{Package: "golang.org/x/net", To: "0.36.0"}},
	}}
	intents2, _ := Plan(cLower, ledHigh, nil, 10)
	assert.Nil(t, intentFor(intents2, IntentOpen, k), "merged at >= version must not re-open")
}

func TestPlanAdoptsExistingExternalPR(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go",
			Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	prs := map[string][]ghclient.PullRequest{
		"kairos-io/immucore": {{Number: 7, Title: "Bump golang.org/x/net to 0.33.0", Author: "renovate[bot]", URL: "u7"}},
	}
	intents, _ := Plan(c, state.Ledger{}, prs, 10)
	require.Len(t, intents, 1)
	assert.Equal(t, IntentAdopt, intents[0].Type)
	assert.Equal(t, 7, intents[0].PRNumber)
	assert.Equal(t, "renovate", intents[0].Source)
}

func TestPlanOpensWhenNoExternalPR(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go",
			Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	intents, _ := Plan(c, state.Ledger{}, nil, 10)
	require.Len(t, intents, 1)
	assert.Equal(t, IntentOpen, intents[0].Type)
}
