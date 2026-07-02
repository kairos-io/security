package collect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const nvdResponseJSON = `{
  "vulnerabilities": [
    {
      "cve": {
        "id": "CVE-2025-5678",
        "descriptions": [{"lang": "en", "value": "Buffer overflow in libfoo"}],
        "metrics": {"cvssMetricV31": [{"cvssData": {"baseSeverity": "HIGH"}}]},
        "configurations": [
          {"nodes": [{"cpeMatch": [{"criteria": "cpe:2.3:a:foo:libfoo:*:*:*:*:*:*:*:*", "vulnerable": true, "versionEndExcluding": "2.5.0"}]}]}
        ]
      }
    }
  ]
}`

func TestQueryNVDParsesHit(t *testing.T) {
	results, err := QueryNVD(func(vendor, product, version string) ([]byte, error) {
		assert.Equal(t, "foo", vendor)
		assert.Equal(t, "libfoo", product)
		assert.Equal(t, "2.4.0", version)
		return []byte(nvdResponseJSON), nil
	}, "foo", "libfoo", "2.4.0")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "CVE-2025-5678", results[0].CVEID)
	assert.Equal(t, "high", results[0].Severity)
	assert.Equal(t, "2.5.0", results[0].VersionEndExcluding)
	assert.Equal(t, "Buffer overflow in libfoo", results[0].Title)
	assert.Equal(t, "https://nvd.nist.gov/vuln/detail/CVE-2025-5678", results[0].URL)
}

func TestQueryNVDNoHits(t *testing.T) {
	results, err := QueryNVD(func(string, string, string) ([]byte, error) {
		return []byte(`{"vulnerabilities": []}`), nil
	}, "gnu", "bash", "5.3")
	require.NoError(t, err)
	assert.Empty(t, results)
}
