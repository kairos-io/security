package remediate

import (
	"strings"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
)

func sampleIntent() Intent {
	return Intent{Type: IntentOpen, Key: "kairos-io/immucore|golang.org/x/net",
		Repo: "kairos-io/immucore", Package: "golang.org/x/net", Severity: "high",
		Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}}
}

func TestBranchAndTitleAndBody(t *testing.T) {
	in := sampleIntent()
	assert.Equal(t, "ksec/bump-golang-org-x-net-0-33-0", BranchName(in))
	assert.Equal(t, "chore(security): bump golang.org/x/net to 0.33.0", PRTitle(in))

	body := PRBody(in)
	assert.Contains(t, body, "golang.org/x/net")
	assert.Contains(t, body, "0.33.0")
	assert.Contains(t, body, "high")
	assert.Contains(t, body, "kairos-security")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(body), PRMarker(in.Key)),
		"marker must be the last line")
}
