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
