package remediate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeProse struct {
	body string
	err  error
}

func (f fakeProse) DraftPRBody(Intent) (string, error) { return f.body, f.err }

func TestPRBodyWithFallsBackOnError(t *testing.T) {
	in := sampleIntent()
	out := PRBodyWith(in, fakeProse{err: assertErrP()})
	assert.Equal(t, PRBody(in), out, "error -> deterministic body")
}

func TestPRBodyWithAppendsAIAndKeepsMarkerLast(t *testing.T) {
	in := sampleIntent()
	out := PRBodyWith(in, fakeProse{body: "AI: this is safe."})
	assert.Contains(t, out, "AI: this is safe.")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(out), PRMarker(in.Key)))
}

func assertErrP() error { return assert.AnError }
