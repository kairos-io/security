package remediate

import (
	"testing"

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

	intents, deferred := Plan(c, state.Ledger{}, 1) // cap to 1 new PR
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
	intents, _ := Plan(c, led, 10)
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
	intents, _ := Plan(c, led, 10)
	// only the reconcile for the existing entry; no new open
	require.Len(t, intents, 1)
	assert.Equal(t, IntentReconcile, intents[0].Type)
}
