package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestForkByKind(t *testing.T) {
	f := ForkByKind([]state.Repo{
		{Repo: "mudler/edgevpn", Kind: "external"},
		{Repo: "kairos-io/kairos", Kind: "org"},
	})
	assert.True(t, f("mudler/edgevpn"))    // external -> fork
	assert.False(t, f("kairos-io/kairos")) // org -> direct
	assert.False(t, f("unknown/repo"))     // unknown -> direct (no surprise fork)
}
