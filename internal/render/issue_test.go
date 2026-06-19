package render

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpsertTrackingIssueWrites(t *testing.T) {
	gh := ghclient.NewFake()
	n, err := UpsertTrackingIssue(gh, "kairos-io/kairos", "body", false)
	require.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, "body", gh.Issues["kairos-io/kairos"].Body)
	assert.Contains(t, gh.Issues["kairos-io/kairos"].Labels, "kairos-security-bot")
}

func TestUpsertTrackingIssueDryRunSkips(t *testing.T) {
	gh := ghclient.NewFake()
	n, err := UpsertTrackingIssue(gh, "kairos-io/kairos", "body", true)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Empty(t, gh.Issues)
}
