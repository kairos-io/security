package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/stretchr/testify/assert"
)

func TestClassifySource(t *testing.T) {
	assert.Equal(t, "renovate", classifySource(ghclient.PullRequest{Author: "renovate[bot]"}))
	assert.Equal(t, "dependabot", classifySource(ghclient.PullRequest{Author: "dependabot[bot]"}))
	assert.Equal(t, "ksec", classifySource(ghclient.PullRequest{Author: "kairos-security-bot"}))
	assert.Equal(t, "ksec", classifySource(ghclient.PullRequest{Author: "alice", HeadRef: "ksec/bump"}))
	assert.Equal(t, "human", classifySource(ghclient.PullRequest{Author: "alice"}))
}

func TestMatchPR(t *testing.T) {
	prs := []ghclient.PullRequest{
		{Number: 1, Title: "Bump golang.org/x/net from 0.30.0 to 0.33.0", Author: "dependabot[bot]"},
		{Number: 2, Title: "Some feature", Author: "alice"},
	}
	pr, src, ok := MatchPR("golang.org/x/net", "", prs)
	assert.True(t, ok)
	assert.Equal(t, 1, pr.Number)
	assert.Equal(t, "dependabot", src)

	_, _, ok = MatchPR("golang.org/x/crypto", "", prs)
	assert.False(t, ok)
}

func TestMatchPRRequiresVersion(t *testing.T) {
	prs := []ghclient.PullRequest{
		{Number: 1, Title: "Bump golang.org/x/net from 0.30.0 to 0.33.0", Author: "dependabot[bot]"},
	}
	// pkg + matching version -> match
	pr, src, ok := MatchPR("golang.org/x/net", "0.33.0", prs)
	assert.True(t, ok)
	assert.Equal(t, 1, pr.Number)
	assert.Equal(t, "dependabot", src)
	// pkg present but version absent -> no match (avoids "remove x/net usage" false positives)
	_, _, ok = MatchPR("golang.org/x/net", "9.9.9", prs)
	assert.False(t, ok)
	// empty version disables the version requirement
	_, _, ok = MatchPR("golang.org/x/net", "", prs)
	assert.True(t, ok)
}

// A version-aware match accepts a title version >= the required minimum, so a
// dependabot/renovate bump that overshoots our floor is still adopted; a lower
// bump or a version-less mention is not.
func TestMatchPRAcceptsHigherVersion(t *testing.T) {
	cases := []struct {
		name     string
		title    string
		required string
		want     bool
	}{
		{"higher bump matches", "Bump golang.org/x/net from 0.30.0 to 0.36.0", "0.33.0", true},
		{"lower bump does not match", "Bump golang.org/x/net from 0.10.0 to 0.20.0", "0.33.0", false},
		{"no version does not match", "remove golang.org/x/net usage", "0.33.0", false},
		{"empty requirement matches", "remove golang.org/x/net usage", "", true},
		{"v-prefixed token matches", "chore: golang.org/x/net v0.36.0", "0.33.0", true},
		{"exact still matches", "Bump golang.org/x/net from 0.30.0 to 0.33.0", "0.33.0", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			prs := []ghclient.PullRequest{{Number: 1, Title: tc.title, Author: "dependabot[bot]"}}
			_, _, ok := MatchPR("golang.org/x/net", tc.required, prs)
			assert.Equal(t, tc.want, ok)
		})
	}
}

func TestIsOwnPRByBranch(t *testing.T) {
	_, src, ok := MatchPR("golang.org/x/net", "", []ghclient.PullRequest{
		{Number: 2, Title: "bump golang.org/x/net", Author: "someoneelse", HeadRef: "ksec/bump-golang-org-x-net-0-33-0"},
	})
	assert.True(t, ok)
	assert.Equal(t, "ksec", src, "a ksec/ branch is ours regardless of author")
}
