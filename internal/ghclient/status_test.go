package ghclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakePRStatusAndMerge(t *testing.T) {
	f := NewFake()
	f.Statuses["o/r#5"] = PRStatus{State: "OPEN", Mergeable: true, ChecksPassing: true, ReviewDecision: ""}
	got, err := f.PRStatusOf("o/r", 5)
	require.NoError(t, err)
	assert.True(t, got.ChecksPassing)

	require.NoError(t, f.MergePR("o/r", 5, true))
	assert.Equal(t, []string{"o/r#5 (auto)"}, f.Merged)
}
