package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/stretchr/testify/assert"
)

func TestShouldAutomerge(t *testing.T) {
	ok := ghclient.PRStatus{Mergeable: true, ChecksPassing: true, ReviewDecision: ""}
	assert.True(t, ShouldAutomerge(ok))

	assert.False(t, ShouldAutomerge(ghclient.PRStatus{Mergeable: false, ChecksPassing: true}))
	assert.False(t, ShouldAutomerge(ghclient.PRStatus{Mergeable: true, ChecksPassing: false}))
	assert.False(t, ShouldAutomerge(ghclient.PRStatus{Mergeable: true, ChecksPassing: true, ReviewDecision: "CHANGES_REQUESTED"}))
	assert.True(t, ShouldAutomerge(ghclient.PRStatus{Mergeable: true, ChecksPassing: true, ReviewDecision: "APPROVED"}))
}
