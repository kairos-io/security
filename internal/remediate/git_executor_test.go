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

func TestRepinRefusesNonKsecBranch(t *testing.T) {
	// The ksec guard (after the State!="open" check, before `go list`/clone)
	// must reject a non-ksec branch without touching the network.
	g := &GitExecutor{}
	_, err := g.Repin(state.LedgerEntry{State: "open", Branch: "main", Package: "m"}, "run")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ksec")
}

func TestAdjustRefusesNonKsecBranch(t *testing.T) {
	// The non-ksec guard is the very first thing Adjust does, so this must fail
	// without touching the network (no clone/push for a real branch like main).
	g := &GitExecutor{}
	_, err := g.Adjust(state.LedgerEntry{Repo: "r", Branch: "main", PRNumber: 1}, "1.0.0", "run")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ksec")
}

func TestForkSlug(t *testing.T) {
	assert.Equal(t, "kairos-security-bot/edgevpn", forkSlug("kairos-security-bot", "mudler/edgevpn"))
	assert.Equal(t, "bot/kairos", forkSlug("bot", "kairos-io/kairos"))
}

func TestPRHead(t *testing.T) {
	ext := func(string) bool { return true }
	org := func(string) bool { return false }
	g := &GitExecutor{ForkOwner: "kairos-security-bot", ShouldFork: ext}
	assert.Equal(t, "kairos-security-bot:ksec/x", g.prHead("mudler/edgevpn", "ksec/x"))

	g2 := &GitExecutor{ForkOwner: "kairos-security-bot", ShouldFork: org}
	assert.Equal(t, "ksec/x", g2.prHead("kairos-io/kairos", "ksec/x"))

	g3 := &GitExecutor{} // nil ShouldFork -> never fork
	assert.Equal(t, "ksec/x", g3.prHead("mudler/edgevpn", "ksec/x"))
}

func TestForkURL(t *testing.T) {
	g := &GitExecutor{ForkOwner: "bot", Token: "tok"}
	assert.Equal(t, "https://x-access-token:tok@github.com/bot/edgevpn.git", g.forkURL("mudler/edgevpn"))
}

func TestPushBranchDryRunNoWrites(t *testing.T) {
	// Dry-run must not shell out, for both fork and non-fork repos.
	g := &GitExecutor{DryRun: true, ForkOwner: "bot", ShouldFork: func(string) bool { return true }}
	assert.NoError(t, g.pushBranch("/tmp/nonexistent", "mudler/edgevpn", "ksec/x", false))
	g2 := &GitExecutor{DryRun: true} // org/no-fork
	assert.NoError(t, g2.pushBranch("/tmp/nonexistent", "kairos-io/kairos", "ksec/x", true))
	assert.NoError(t, g.ensureFork("mudler/edgevpn"))
}
