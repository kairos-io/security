package remediate

import (
	"testing"

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
