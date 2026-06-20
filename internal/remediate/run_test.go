package remediate

import (
	"errors"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunOpensReconcilesAndIsolatesErrors(t *testing.T) {
	intents := []Intent{
		{Type: IntentOpen, Key: "r|p1", Repo: "r", Package: "p1", Bump: state.Bump{Package: "p1", To: "1.0.0"}},
		{Type: IntentOpen, Key: "r|p2", Repo: "r", Package: "p2", Bump: state.Bump{Package: "p2", To: "2.0.0"}},
		{Type: IntentReconcile, Key: "r|old", Entry: &state.LedgerEntry{Key: "r|old", Repo: "r", State: "open"}},
	}
	fake := &FakeExecutor{
		Opened: map[string]state.LedgerEntry{
			"r|p1": {Key: "r|p1", Repo: "r", Package: "p1", State: "open", PRNumber: 1},
		},
		OpenErr:    map[string]error{"r|p2": errors.New("build-failed")},
		Reconciled: map[string]state.LedgerEntry{"r|old": {Key: "r|old", Repo: "r", State: "merged"}},
	}

	out, results := Run(intents, fake, state.Ledger{}, "2026-06-19")

	byKey := map[string]state.LedgerEntry{}
	for _, e := range out.Entries {
		byKey[e.Key] = e
	}
	assert.Equal(t, "open", byKey["r|p1"].State)
	assert.Equal(t, "error", byKey["r|p2"].State, "open failure recorded, not aborted")
	assert.Equal(t, "merged", byKey["r|old"].State)
	// deterministic order
	require.Len(t, out.Entries, 3)
	assert.True(t, out.Entries[0].Key <= out.Entries[1].Key)
	// a result per intent
	assert.Len(t, results, 3)
}

func TestRunCascadeAndRepin(t *testing.T) {
	entry := state.LedgerEntry{Key: "c|m", Repo: "c", Package: "m", State: "open", Kind: "cascade", Pseudo: true}
	intents := []Intent{
		{Type: IntentCascade, Key: "c|m", Repo: "c", Package: "m", Ref: "main", CascadeFrom: "u|x"},
		{Type: IntentRepin, Key: "c|m", Entry: &entry},
	}
	fake := &FakeExecutor{
		Cascaded: map[string]state.LedgerEntry{"c|m": {Key: "c|m", Repo: "c", Package: "m", State: "open", Kind: "cascade", Pseudo: true}},
		Repinned: map[string]state.LedgerEntry{"c|m": {Key: "c|m", Repo: "c", Package: "m", State: "open", Kind: "cascade", Pseudo: false, PinTarget: "v1.0.0"}},
	}
	out, results := Run(intents, fake, state.Ledger{}, "2026-06-20")
	require.Len(t, out.Entries, 1)
	// repin ran after cascade (same key) -> pseudo cleared, pin set
	assert.False(t, out.Entries[0].Pseudo)
	assert.Equal(t, "v1.0.0", out.Entries[0].PinTarget)
	require.Len(t, results, 2)
}

// Repin must build on the reconciled entry, not the stale pre-run snapshot.
// Here reconcile flips the entry to merged and repin echoes its input; the
// final entry must be merged (proving repin read the reconciled entry).
func TestRunRepinBuildsOnReconciledEntry(t *testing.T) {
	entry := state.LedgerEntry{Key: "c|m", Repo: "c", Package: "m", State: "open", Kind: "cascade", Pseudo: true}
	intents := []Intent{
		{Type: IntentReconcile, Key: "c|m", Repo: "c", Entry: &entry},
		{Type: IntentRepin, Key: "c|m", Entry: &entry},
	}
	fake := &FakeExecutor{
		Reconciled: map[string]state.LedgerEntry{
			"c|m": {Key: "c|m", Repo: "c", Package: "m", State: "merged", Kind: "cascade", Pseudo: true},
		},
		// Repin echoes whatever it is handed.
	}
	out, _ := Run(intents, fake, state.Ledger{}, "2026-06-20")
	require.Len(t, out.Entries, 1)
	assert.Equal(t, "merged", out.Entries[0].State,
		"repin must operate on the reconciled (merged) entry, not the stale open snapshot")
}

func TestRunToolchain(t *testing.T) {
	intents := []Intent{
		{Type: IntentToolchain, Key: "r|go-toolchain", Repo: "r", ToolchainVersion: "1.22.5", Severity: "high"},
	}
	fake := &FakeExecutor{Toolchained: map[string]state.LedgerEntry{
		"r|go-toolchain": {Key: "r|go-toolchain", Repo: "r", Package: "go-toolchain", Kind: "toolchain", State: "open"},
	}}
	out, results := Run(intents, fake, state.Ledger{}, "2026-06-20")
	require.Len(t, out.Entries, 1)
	assert.Equal(t, "toolchain", out.Entries[0].Kind)
	assert.Equal(t, "open", out.Entries[0].State)
	require.Len(t, results, 1)
	assert.Equal(t, "toolchain", results[0].Action)
}

func TestRunToolchainIsolatesErrors(t *testing.T) {
	intents := []Intent{
		{Type: IntentToolchain, Key: "r|go-toolchain", Repo: "r", ToolchainVersion: "1.22.5", Severity: "high"},
	}
	fake := &FakeExecutor{ToolchainErr: map[string]error{"r|go-toolchain": errors.New("clone-failed")}}
	out, results := Run(intents, fake, state.Ledger{}, "2026-06-20")
	require.Len(t, out.Entries, 1)
	assert.Equal(t, "error", out.Entries[0].State, "toolchain failure recorded, not aborted")
	assert.Equal(t, "toolchain", out.Entries[0].Kind)
	require.Len(t, results, 1)
	assert.Equal(t, "toolchain", results[0].Action)
	assert.Equal(t, "error", results[0].State)
}

func TestRunAdopts(t *testing.T) {
	intents := []Intent{
		{Type: IntentAdopt, Key: "r|p", Repo: "r", Package: "p", PRNumber: 9, PRURL: "u9", Source: "dependabot"},
	}
	fake := &FakeExecutor{Adopted: map[string]state.LedgerEntry{
		"r|p": {Key: "r|p", Repo: "r", Package: "p", State: "open", PRNumber: 9, Source: "dependabot"},
	}}
	out, results := Run(intents, fake, state.Ledger{}, "2026-06-20")
	require.Len(t, out.Entries, 1)
	assert.Equal(t, "dependabot", out.Entries[0].Source)
	require.Len(t, results, 1)
	assert.Equal(t, "adopt", results[0].Action)
}
