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

// TestQueryOSV_PreReleaseIntroducedBoundary exercises the call site for a Go
// advisory. golang/vulndb writes the "introduced" boundary of a whole minor
// line as "<line>-0", the lowest possible pre-release of it (GO-2026-4601:
// {"introduced":"1.26.0-0"},{"fixed":"1.26.1"}). Ordering "1.26.0" below that
// boundary dropped the range, and with no other range applicable the vuln was
// reported as not affecting a version that is squarely inside it.
func TestQueryOSV_PreReleaseIntroducedBoundary(t *testing.T) {
	fixture := `{"vulns":[{"id":"GO-2026-4601","aliases":["CVE-2026-4601"],
	  "affected":[{"package":{"ecosystem":"Go","name":"stdlib"},
	  "ranges":[{"type":"SEMVER","events":[
	    {"introduced":"1.26.0-0"},{"fixed":"1.26.1"}]}]}]}]}`
	q := func(_, _, _ string) ([]byte, error) { return []byte(fixture), nil }

	for _, tc := range []struct {
		version, wantFixed string
		wantHit            bool
	}{
		{"1.26.0", "1.26.1", true}, // the release itself: inside the window
		{"1.26.0-rc.1", "1.26.1", true},
		{"1.26.1", "1.26.1", true}, // past the fix: surfaced, classified later
		{"1.25.9", "", false},      // genuinely below the line: still dropped
	} {
		got, err := QueryOSV(q, "Go", "stdlib", tc.version)
		if err != nil {
			t.Fatal(err)
		}
		if !tc.wantHit {
			if len(got) != 0 {
				t.Errorf("%s: want dropped below the introduced boundary, got %+v", tc.version, got)
			}
			continue
		}
		if len(got) != 1 {
			t.Fatalf("%s: want the vuln surfaced, got %+v", tc.version, got)
		}
		if got[0].FixedVersion != tc.wantFixed {
			t.Errorf("%s: want fixed %q, got %q", tc.version, tc.wantFixed, got[0].FixedVersion)
		}
	}
}

// TestQueryOSV_MultiIntervalRange: an OSV range's events are a timeline, and a
// CVE patched on several release branches carries every window in one range.
// Reading only the last introduced/fixed pair collapses them into the newest
// window, so a version inside an older window is reported as unaffected.
//
// The fixture is the real shape of golang/vulndb GO-2026-4340 (vulnerable in
// [0, 1.24.12) and in [1.25.0, 1.25.6)); 868 of the 4474 records in that
// database use a multi-window range.
func TestQueryOSV_MultiIntervalRange(t *testing.T) {
	fixture := `{"vulns":[{"id":"GO-2026-4340","affected":[{"ranges":[{"type":"SEMVER","events":[
	  {"introduced":"0"},{"fixed":"1.24.12"},{"introduced":"1.25.0"},{"fixed":"1.25.6"}]}]}]}]}`
	q := func(_, _, _ string) ([]byte, error) { return []byte(fixture), nil }

	for _, tc := range []struct {
		name    string
		version string
		fixed   string
	}{
		{"inside the first window", "1.24.5", "1.24.12"},
		{"inside the second window", "1.25.2", "1.25.6"},
		{"between the windows, already fixed", "1.24.12", "1.24.12"},
		{"past every window, already fixed", "1.25.9", "1.25.6"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := QueryOSV(q, "Go", "stdlib", tc.version)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 {
				t.Fatalf("version %s: want the vuln surfaced, got %+v", tc.version, got)
			}
			if got[0].FixedVersion != tc.fixed {
				t.Fatalf("version %s: want FixedVersion %q, got %q", tc.version, tc.fixed, got[0].FixedVersion)
			}
		})
	}

	// Below every window's introduced boundary is still not vulnerable.
	fixture2 := `{"vulns":[{"id":"GO-x","affected":[{"ranges":[{"type":"SEMVER","events":[
	  {"introduced":"1.24.0"},{"fixed":"1.24.12"},{"introduced":"1.25.0"},{"fixed":"1.25.6"}]}]}]}]}`
	q2 := func(_, _, _ string) ([]byte, error) { return []byte(fixture2), nil }
	got, err := QueryOSV(q2, "Go", "stdlib", "1.23.1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("version below every introduced boundary should be omitted, got %+v", got)
	}
}

// TestRangeIntervals covers the timeline walk on its own, including the
// event shapes the schema permits but Go's database does not emit.
func TestRangeIntervals(t *testing.T) {
	for _, tc := range []struct {
		name   string
		events []osvRangeEvent
		want   []osvInterval
	}{
		{
			name:   "single window",
			events: []osvRangeEvent{{Introduced: "0"}, {Fixed: "2.66.6"}},
			want:   []osvInterval{{introduced: "0", fixed: "2.66.6"}},
		},
		{
			name: "two windows",
			events: []osvRangeEvent{{Introduced: "0"}, {Fixed: "1.24.12"},
				{Introduced: "1.25.0"}, {Fixed: "1.25.6"}},
			want: []osvInterval{{introduced: "0", fixed: "1.24.12"},
				{introduced: "1.25.0", fixed: "1.25.6"}},
		},
		{
			name:   "unfixed window stays open",
			events: []osvRangeEvent{{Introduced: "1.1.1"}},
			want:   []osvInterval{{introduced: "1.1.1"}},
		},
		{
			name: "an introduced closes an unfixed predecessor",
			events: []osvRangeEvent{{Introduced: "1.0"}, {Introduced: "2.0"},
				{Fixed: "2.5"}},
			want: []osvInterval{{introduced: "1.0"}, {introduced: "2.0", fixed: "2.5"}},
		},
		{
			name:   "a leading fix starts at 0",
			events: []osvRangeEvent{{Fixed: "3.2"}},
			want:   []osvInterval{{introduced: "0", fixed: "3.2"}},
		},
		{
			name:   "alpine revision suffixes are stripped from the fix",
			events: []osvRangeEvent{{Introduced: "0"}, {Fixed: "3.6.4-r0"}},
			want:   []osvInterval{{introduced: "0", fixed: "3.6.4"}},
		},
		{
			name:   "no usable events means all versions",
			events: nil,
			want:   []osvInterval{{introduced: "0"}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := rangeIntervals(tc.events)
			if len(got) != len(tc.want) {
				t.Fatalf("want %d interval(s) %+v, got %d %+v", len(tc.want), tc.want, len(got), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("interval %d: want %+v, got %+v", i, tc.want[i], got[i])
				}
			}
		})
	}
}
