package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunRedactsToken(t *testing.T) {
	g := &GitExecutor{Token: "SUPERSECRETTOKEN"}
	// An unknown git subcommand fails fast (no network) and the token-bearing
	// arg appears in the wrapped error — it must be redacted.
	_, err := g.run("", "git", "definitely-not-a-subcommand-SUPERSECRETTOKEN")
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "SUPERSECRETTOKEN")
	assert.Contains(t, err.Error(), "***")
}

func TestAdjustRefusesNonKsecBranch(t *testing.T) {
	// The non-ksec guard is the very first thing Adjust does, so this must fail
	// without touching the network (no clone/push for a real branch like main).
	g := &GitExecutor{}
	_, err := g.Adjust(state.LedgerEntry{Repo: "r", Branch: "main", PRNumber: 1}, "1.0.0", "run")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ksec")
}
