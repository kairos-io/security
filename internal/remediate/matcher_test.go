package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/stretchr/testify/assert"
)

func TestClassifySource(t *testing.T) {
	assert.Equal(t, "renovate", classifySource("renovate[bot]"))
	assert.Equal(t, "dependabot", classifySource("dependabot[bot]"))
	assert.Equal(t, "ksec", classifySource("kairos-security-bot"))
	assert.Equal(t, "human", classifySource("alice"))
}

func TestMatchPR(t *testing.T) {
	prs := []ghclient.PullRequest{
		{Number: 1, Title: "Bump golang.org/x/net from 0.30.0 to 0.33.0", Author: "dependabot[bot]"},
		{Number: 2, Title: "Some feature", Author: "alice"},
	}
	pr, src, ok := MatchPR("golang.org/x/net", prs)
	assert.True(t, ok)
	assert.Equal(t, 1, pr.Number)
	assert.Equal(t, "dependabot", src)

	_, _, ok = MatchPR("golang.org/x/crypto", prs)
	assert.False(t, ok)
}
