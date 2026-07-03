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

// Real-world shape of an Alpine OSV-converted advisory: the real CVE id lives in
// "upstream" (not "aliases", which is absent), and the long free-text lives in
// "details" (not "summary", which is absent). QueryOSV must surface the upstream
// CVE id, not the Alpine-internal "ALPINE-CVE-…" id.
const osvAlpineUpstreamJSON = `{
  "vulns": [
    {
      "id": "ALPINE-CVE-2023-5363",
      "upstream": ["CVE-2023-5363"],
      "details": "A bug has been identified in the processing of key and initialisation vector (IV) lengths. This can lead to potential truncation or overruns during the initialisation of some symmetric ciphers. A truncation in the IV can result in non-uniqueness, which could result in loss of confidentiality for some cipher modes.",
      "database_specific": {"severity": "HIGH"},
      "affected": [
        {"ranges": [{"type": "ECOSYSTEM", "events": [{"introduced": "0"}, {"fixed": "3.1.4-r0"}]}]}
      ]
    }
  ]
}`

func TestQueryOSVUsesUpstreamCVEIDForAlpineRecords(t *testing.T) {
	results, err := QueryOSV(func(string, string, string) ([]byte, error) {
		return []byte(osvAlpineUpstreamJSON), nil
	}, "Alpine", "openssl", "3.1.3")
	require.NoError(t, err)
	require.Len(t, results, 1)
	// The real CVE id from "upstream", not the Alpine-internal "ALPINE-CVE-…" id.
	assert.Equal(t, "CVE-2023-5363", results[0].CVEID)
	assert.Equal(t, "high", results[0].Severity)
	assert.Equal(t, "3.1.4", results[0].FixedVersion)
}

func TestQueryOSVNoHits(t *testing.T) {
	results, err := QueryOSV(func(string, string, string) ([]byte, error) {
		return []byte(`{"vulns": []}`), nil
	}, "Alpine", "bash", "5.3")
	require.NoError(t, err)
	assert.Empty(t, results)
}

// Real-world shape of ALPINE-CVE-2025-66199: no database_specific.severity, but a
// top-level severity array carrying a CVSS v3.1 vector that computes to 5.9 (Medium).
// Regression test for the bug where severityFromOSV("") defaulted to "high".
const osvCVSSOnlyJSON = `{
  "vulns": [
    {
      "id": "ALPINE-CVE-2025-66199",
      "aliases": ["CVE-2025-66199"],
      "summary": "openssl TLS 1.3 certificate decompression DoS",
      "severity": [{"type": "CVSS_V3", "score": "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H"}],
      "affected": [
        {"ranges": [{"type": "ECOSYSTEM", "events": [{"introduced": "0"}, {"fixed": "3.6.4-r0"}]}]}
      ]
    }
  ]
}`

func TestQueryOSVDerivesSeverityFromCVSSVector(t *testing.T) {
	results, err := QueryOSV(func(string, string, string) ([]byte, error) {
		return []byte(osvCVSSOnlyJSON), nil
	}, "Alpine", "openssl", "3.6.3")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "CVE-2025-66199", results[0].CVEID)
	assert.Equal(t, "medium", results[0].Severity) // 5.9 -> medium, NOT the old "high" default
}

// When database_specific.severity is present it is an explicit, human/tooling-assigned
// label and must win over any CVSS-computed value, even a conflicting one.
const osvSeverityPrecedenceJSON = `{
  "vulns": [
    {
      "id": "GHSA-aaaa-bbbb-cccc",
      "aliases": ["CVE-2025-9999"],
      "summary": "conflicting severity sources",
      "database_specific": {"severity": "HIGH"},
      "severity": [{"type": "CVSS_V3", "score": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:C/C:H/I:H/A:H"}],
      "affected": [{"ranges": [{"type": "ECOSYSTEM", "events": [{"fixed": "1.0.0"}]}]}]
    }
  ]
}`

func TestQueryOSVDatabaseSpecificSeverityWinsOverCVSS(t *testing.T) {
	results, err := QueryOSV(func(string, string, string) ([]byte, error) {
		return []byte(osvSeverityPrecedenceJSON), nil
	}, "Alpine", "pkg", "0.9.0")
	require.NoError(t, err)
	require.Len(t, results, 1)
	// CVSS vector would compute to critical (10.0); explicit "HIGH" must win.
	assert.Equal(t, "high", results[0].Severity)
}

// Neither an explicit label nor a parseable CVSS_V3 entry: be honest, return "unknown".
const osvNoUsableSeverityJSON = `{
  "vulns": [
    {
      "id": "ALPINE-CVE-2025-00000",
      "aliases": ["CVE-2025-00000"],
      "summary": "no severity data at all",
      "affected": [{"ranges": [{"type": "ECOSYSTEM", "events": [{"fixed": "2.0.0"}]}]}]
    }
  ]
}`

func TestQueryOSVUnknownWhenNoSeverityData(t *testing.T) {
	results, err := QueryOSV(func(string, string, string) ([]byte, error) {
		return []byte(osvNoUsableSeverityJSON), nil
	}, "Alpine", "pkg", "1.0.0")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "unknown", results[0].Severity)
}
