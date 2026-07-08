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

func TestQueryOSV_RangeApplicability(t *testing.T) {
	// One vuln, two branch ranges: introduced 0 fixed 2.66.6, and introduced
	// 2.80 fixed 2.86.0. Queried version 2.86.2 is past both fixes.
	fixture := `{"vulns":[{"id":"CVE-x","affected":[{"ranges":[{"events":[
	  {"introduced":"0"},{"fixed":"2.66.6"}]},{"events":[
	  {"introduced":"2.80"},{"fixed":"2.86.0"}]}]}]}]}`
	q := func(_, _, _ string) ([]byte, error) { return []byte(fixture), nil }

	// 2.86.2 is >= the applicable fix (2.86.0) -> still returned, FixedVersion=2.86.0.
	got, err := QueryOSV(q, "Alpine:v3.22", "glib", "2.86.2")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].FixedVersion != "2.86.0" {
		t.Fatalf("want 1 result fixed 2.86.0, got %+v", got)
	}

	// A version below every introduced is not yet vulnerable -> omitted.
	// (introduced values here are "0" so nothing is below; use a fixture with a
	// higher introduced.)
	fixture2 := `{"vulns":[{"id":"CVE-y","affected":[{"ranges":[{"events":[
	  {"introduced":"3.0"},{"fixed":"3.2"}]}]}]}]}`
	q2 := func(_, _, _ string) ([]byte, error) { return []byte(fixture2), nil }
	got2, err := QueryOSV(q2, "Alpine:v3.22", "pkg", "2.9")
	if err != nil {
		t.Fatal(err)
	}
	if len(got2) != 0 {
		t.Fatalf("version below introduced should be omitted, got %+v", got2)
	}
}

// TestQueryOSV_ZeroFixedTreatedAsUnfixed: OSV Alpine records use "0" as a
// placeholder meaning "no known fix yet". Returning that literally as
// FixedVersion made the deterministic classifier compare `current >= "0"` and
// mark every such finding as already-fixed, silently hiding unpatched vulns.
// Must be treated the same as an empty fixed: applicable, no target version.
func TestQueryOSV_ZeroFixedTreatedAsUnfixed(t *testing.T) {
	// Alpine range with "fixed":"0" placeholder on a version we clearly have
	// not fixed yet (introduced 1.1.1, current 3.6.3).
	fixture := `{"vulns":[{"id":"CVE-zero","affected":[{"ranges":[{"events":[
	  {"introduced":"1.1.1"},{"fixed":"0"}]}]}]}]}`
	q := func(_, _, _ string) ([]byte, error) { return []byte(fixture), nil }
	got, err := QueryOSV(q, "Alpine:v3.22", "openssl", "3.6.3")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("want 1 result surfaced, got %+v", got)
	}
	if got[0].FixedVersion != "" {
		t.Fatalf(`fixed "0" must be treated as unfixed (empty FixedVersion), got %q`, got[0].FixedVersion)
	}
}

// TestQueryOSV_UnparseableVersionFailsOpen: a non-numeric queried version can't
// be ordered against the range boundaries, so instead of silently dropping the
// vuln (which would hide it), the matcher must fail OPEN and surface it.
func TestQueryOSV_UnparseableVersionFailsOpen(t *testing.T) {
	fixture := `{"vulns":[{"id":"CVE-z","affected":[{"ranges":[{"events":[
	  {"introduced":"3.0"},{"fixed":"3.2"}]}]}]}]}`
	q := func(_, _, _ string) ([]byte, error) { return []byte(fixture), nil }
	// Both a word-leading ("unknown") and a punctuation-leading ("+incompatible")
	// version are non-numeric; neither can be ordered against "3.0", so both must
	// surface the vuln rather than being dropped below the introduced boundary.
	for _, q0 := range []string{"unknown", "+incompatible"} {
		got, err := QueryOSV(q, "Alpine:v3.22", "pkg", q0)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 {
			t.Fatalf("unparseable version %q must fail open (vuln surfaced), got %+v", q0, got)
		}
	}
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
