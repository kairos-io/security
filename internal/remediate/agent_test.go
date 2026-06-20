package remediate

import (
	"errors"
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakeAgentRecordsAndErrors(t *testing.T) {
	f := &FakeAgent{}
	require.NoError(t, f.Repair("/tmp/x", "fix the build"))
	assert.Equal(t, []string{"fix the build"}, f.Calls)

	f2 := &FakeAgent{Err: errors.New("nope")}
	assert.Error(t, f2.Repair("/tmp/x", "t"))
}

func TestNewNibAgent(t *testing.T) {
	a := NewNibAgent(config.AIConfig{})
	require.NotNil(t, a)
}
