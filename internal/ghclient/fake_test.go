package ghclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakeUpsertCreatesThenUpdates(t *testing.T) {
	f := NewFake()
	n, err := f.UpsertIssue("kairos-io/kairos", "<!-- ksec:dashboard -->", "Security", "body v1", []string{"security"})
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	n2, err := f.UpsertIssue("kairos-io/kairos", "<!-- ksec:dashboard -->", "Security", "body v2", []string{"security"})
	require.NoError(t, err)
	assert.Equal(t, 1, n2, "same marker reuses the issue")
	assert.Equal(t, "body v2", f.Issues["kairos-io/kairos"].Body)
}

func TestFakeListOrgRepos(t *testing.T) {
	f := NewFake()
	f.OrgRepos["kairos-io"] = []string{"kairos-io/kairos", "kairos-io/immucore"}
	got, err := f.ListOrgRepos("kairos-io")
	require.NoError(t, err)
	assert.Equal(t, []string{"kairos-io/kairos", "kairos-io/immucore"}, got)
}
