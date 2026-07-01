package collect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const osvResponseJSON = `{
  "vulns": [
    {
      "id": "GHSA-xxxx-yyyy-zzzz",
      "aliases": ["CVE-2025-1234"],
      "summary": "openssl heap overflow",
      "database_specific": {"severity": "HIGH"},
      "affected": [
        {"ranges": [{"type": "ECOSYSTEM", "events": [{"introduced": "0"}, {"fixed": "3.6.4-r0"}]}]}
      ]
    }
  ]
}`

func TestQueryOSVParsesHitWithAlpineFixedSuffixStripped(t *testing.T) {
	results, err := QueryOSV(func(ecosystem, pkg, version string) ([]byte, error) {
		assert.Equal(t, "Alpine", ecosystem)
		assert.Equal(t, "openssl", pkg)
		assert.Equal(t, "3.6.3", version)
		return []byte(osvResponseJSON), nil
	}, "Alpine", "openssl", "3.6.3")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "CVE-2025-1234", results[0].CVEID)
	assert.Equal(t, "high", results[0].Severity)
	assert.Equal(t, "3.6.4", results[0].FixedVersion) // "-r0" Alpine revision suffix stripped
	assert.Equal(t, "openssl heap overflow", results[0].Title)
	assert.Equal(t, "https://osv.dev/vulnerability/GHSA-xxxx-yyyy-zzzz", results[0].URL)
}

func TestQueryOSVNoHits(t *testing.T) {
	results, err := QueryOSV(func(string, string, string) ([]byte, error) {
		return []byte(`{"vulns": []}`), nil
	}, "Alpine", "bash", "5.3")
	require.NoError(t, err)
	assert.Empty(t, results)
}
