package remediate

import (
	"errors"
	"strings"
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

func TestNibPromptIsSingleLine(t *testing.T) {
	// nib's CLI reads one prompt per stdin line; a multi-line task must collapse
	// to a single line terminated by exactly one newline.
	out := nibPrompt("A bump broke the build.\nFix it so `go build ./...` compiles.\n\nDo not change versions.")
	assert.Equal(t, "A bump broke the build. Fix it so `go build ./...` compiles. Do not change versions.\n", out)
	assert.Equal(t, 1, strings.Count(out, "\n"))
	assert.True(t, strings.HasSuffix(out, "\n"))
}

func TestNibBaseURL(t *testing.T) {
	assert.Equal(t, "http://localhost:8080/v1", nibBaseURL("http://localhost:8080"))
	assert.Equal(t, "http://localhost:8080/v1", nibBaseURL("http://localhost:8080/"))
	assert.Equal(t, "http://localhost:8080/v1", nibBaseURL("http://localhost:8080/v1"))
	assert.Equal(t, "", nibBaseURL(""))
}
