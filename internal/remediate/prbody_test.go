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

func TestCascadePRBodyAndBranch(t *testing.T) {
	in := Intent{Type: IntentCascade, Key: "kairos-io/immucore|github.com/kairos-io/kairos-sdk",
		Repo: "kairos-io/immucore", Package: "github.com/kairos-io/kairos-sdk", CascadeFrom: "kairos-io/kairos-sdk|x", Severity: "high"}
	assert.Equal(t, "ksec/cascade-github-com-kairos-io-kairos-sdk-pseudo", CascadeBranchName(in))
	body := CascadePRBody(in)
	assert.Contains(t, body, "github.com/kairos-io/kairos-sdk")
	assert.Contains(t, body, "pseudo")
	assert.Contains(t, strings.ToLower(body), "tag")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(body), PRMarker(in.Key)))
}
