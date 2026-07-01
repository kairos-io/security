package collect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const hadronManifestJSON = `{
  "ref": "main",
  "commit": "60a9683",
  "generated": "",
  "groups": {
    "Security Tools": {"openssl": "3.6.3", "openssh": "10.3p1"},
    "Compiler Tools": {"gcc": "15.3.0"}
  }
}`

func TestParseHadronManifest(t *testing.T) {
	comps, err := ParseHadronManifest([]byte(hadronManifestJSON))
	require.NoError(t, err)
	require.Len(t, comps, 3)
	// Sorted by group then package, for deterministic output.
	assert.Equal(t, HadronComponent{Group: "Compiler Tools", Package: "gcc", Version: "15.3.0"}, comps[0])
	assert.Equal(t, HadronComponent{Group: "Security Tools", Package: "openssh", Version: "10.3p1"}, comps[1])
	assert.Equal(t, HadronComponent{Group: "Security Tools", Package: "openssl", Version: "3.6.3"}, comps[2])
}

func TestParseHadronManifestInvalidJSON(t *testing.T) {
	_, err := ParseHadronManifest([]byte("not json"))
	assert.Error(t, err)
}
