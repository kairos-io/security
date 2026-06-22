package remediate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSeed(t *testing.T) {
	f, err := ParseSeed("mudler/edgevpn=golang.org/x/net@0.33.0")
	require.NoError(t, err)
	assert.Equal(t, "mudler/edgevpn", f.Repo)
	assert.Equal(t, "golang.org/x/net", f.Package) // slashes preserved
	assert.Equal(t, "0.33.0", f.FixedVersion)
	assert.Equal(t, "sourceCVE", f.Type)
	assert.Equal(t, "go", f.Ecosystem)
	assert.Equal(t, "high", f.Severity)
	assert.True(t, actionable(f)) // the planner will turn it into a target
}

func TestParseSeedErrors(t *testing.T) {
	for _, bad := range []string{"", "no-eq", "repo=pkgNoVersion", "=pkg@1", "r=@1", "r=pkg@"} {
		_, err := ParseSeed(bad)
		assert.Error(t, err, "spec %q should error", bad)
	}
}
