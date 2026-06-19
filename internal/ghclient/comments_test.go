package ghclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakeCommentOps(t *testing.T) {
	f := NewFake()
	f.PRComments["kairos-io/immucore#412"] = []ReviewComment{
		{ID: "c1", Author: "maintainer", Body: "please pin to 0.36.0"},
	}
	got, err := f.ListPRComments("kairos-io/immucore", 412)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "c1", got[0].ID)

	require.NoError(t, f.PostPRComment("kairos-io/immucore", 412, "on it"))
	assert.Equal(t, []string{"kairos-io/immucore#412: on it"}, f.Posted)

	require.NoError(t, f.ClosePR("kairos-io/immucore", 412, "superseded"))
	assert.Equal(t, []string{"kairos-io/immucore#412"}, f.Closed)
}
