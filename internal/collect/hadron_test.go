package collect

import (
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponentManifestCollectSkipsReposWithoutArtifact(t *testing.T) {
	c := ComponentManifest{
		FetchManifest: func() ([]byte, error) { t.Fatal("should not fetch"); return nil, nil },
	}
	fs, err := c.Collect(state.Repo{Repo: "kairos-io/kairos"}) // no component-manifest artifact
	require.NoError(t, err)
	assert.Empty(t, fs)
}

func TestComponentManifestCollectOSVHitThenSkipsNVD(t *testing.T) {
	nowFn = func() string { return "2026-07-01" }
	defer func() { nowFn = defaultNow }()

	manifest := []byte(`{"ref":"main","commit":"abc","groups":{"Security Tools":{"openssl":"3.6.3"}}}`)
	c := ComponentManifest{
		FetchManifest: func() ([]byte, error) { return manifest, nil },
		Components: map[string]config.HadronComponentEntry{
			"openssl": {OSV: &config.HadronOSVSource{Ecosystem: "Alpine", Package: "openssl"}},
		},
		QueryOSV: func(ecosystem, pkg, version string) ([]byte, error) {
			return []byte(`{"vulns":[{"id":"CVE-2025-1","summary":"x","database_specific":{"severity":"HIGH"},"affected":[{"ranges":[{"events":[{"fixed":"3.6.4-r0"}]}]}]}]}`), nil
		},
		QueryNVD: func(string, string, string) ([]byte, error) {
			t.Fatal("NVD should not be queried when OSV already hit")
			return nil, nil
		},
	}
	repo := state.Repo{Repo: "kairos-io/hadron", Artifacts: []state.Artifact{{Type: "component-manifest", Ref: "https://hadron-linux.io/components/main.json"}}}
	fs, err := c.Collect(repo)
	require.NoError(t, err)
	require.Len(t, fs, 1)
	f := fs[0]
	assert.Equal(t, "componentCVE", f.Type)
	assert.Equal(t, "kairos-io/hadron", f.Repo)
	assert.Equal(t, "CVE-2025-1", f.CVEID)
	assert.Equal(t, "openssl", f.Package)
	assert.Equal(t, "3.6.3", f.CurrentVersion)
	assert.Equal(t, "3.6.4", f.FixedVersion)
	assert.Equal(t, "high", f.Severity)
	assert.Equal(t, "hadron", f.Ecosystem)
	assert.Equal(t, "osv", f.Source)
}

func TestComponentManifestCollectFallsBackToNVDWhenOSVEmpty(t *testing.T) {
	nowFn = func() string { return "2026-07-01" }
	defer func() { nowFn = defaultNow }()

	manifest := []byte(`{"ref":"main","commit":"abc","groups":{"Other":{"libfoo":"2.4.0"}}}`)
	c := ComponentManifest{
		FetchManifest: func() ([]byte, error) { return manifest, nil },
		Components: map[string]config.HadronComponentEntry{
			"libfoo": {
				OSV: &config.HadronOSVSource{Ecosystem: "Alpine", Package: "libfoo"},
				CPE: &config.HadronCPESource{Vendor: "foo", Product: "libfoo"},
			},
		},
		QueryOSV: func(string, string, string) ([]byte, error) { return []byte(`{"vulns":[]}`), nil },
		QueryNVD: func(vendor, product, version string) ([]byte, error) {
			return []byte(`{"vulnerabilities":[{"cve":{"id":"CVE-2025-9","descriptions":[{"lang":"en","value":"y"}],"metrics":{"cvssMetricV31":[{"cvssData":{"baseSeverity":"MEDIUM"}}]},"configurations":[{"nodes":[{"cpeMatch":[{"vulnerable":true,"versionEndExcluding":"2.5.0"}]}]}]}}]}`), nil
		},
	}
	repo := state.Repo{Repo: "kairos-io/hadron", Artifacts: []state.Artifact{{Type: "component-manifest"}}}
	fs, err := c.Collect(repo)
	require.NoError(t, err)
	require.Len(t, fs, 1)
	assert.Equal(t, "CVE-2025-9", fs[0].CVEID)
	assert.Equal(t, "nvd", fs[0].Source)
	assert.Equal(t, "medium", fs[0].Severity)
	assert.Equal(t, "2.5.0", fs[0].FixedVersion)
}

// TestComponentManifestCollectFiltersNonVulnerableNVDMatch guards against a
// false positive: QueryNVD (Task 5) emits one NVDResult per NVD
// "vulnerabilities[]" entry regardless of whether any cpeMatch was actually
// vulnerable==true for the queried CPE (NVD configurations legitimately
// include vulnerable:false entries for platform/AND conditions). A CVE whose
// only cpeMatch is vulnerable:false must not become a Finding.
func TestComponentManifestCollectFiltersNonVulnerableNVDMatch(t *testing.T) {
	manifest := []byte(`{"ref":"main","commit":"abc","groups":{"Other":{"libfoo":"2.4.0"}}}`)
	c := ComponentManifest{
		FetchManifest: func() ([]byte, error) { return manifest, nil },
		Components: map[string]config.HadronComponentEntry{
			"libfoo": {
				OSV: &config.HadronOSVSource{Ecosystem: "Alpine", Package: "libfoo"},
				CPE: &config.HadronCPESource{Vendor: "foo", Product: "libfoo"},
			},
		},
		QueryOSV: func(string, string, string) ([]byte, error) { return []byte(`{"vulns":[]}`), nil },
		QueryNVD: func(vendor, product, version string) ([]byte, error) {
			return []byte(`{"vulnerabilities":[{"cve":{"id":"CVE-2025-9","descriptions":[{"lang":"en","value":"y"}],"metrics":{"cvssMetricV31":[{"cvssData":{"baseSeverity":"MEDIUM"}}]},"configurations":[{"nodes":[{"cpeMatch":[{"vulnerable":false,"versionEndExcluding":"2.5.0"}]}]}]}}]}`), nil
		},
	}
	repo := state.Repo{Repo: "kairos-io/hadron", Artifacts: []state.Artifact{{Type: "component-manifest"}}}
	fs, err := c.Collect(repo)
	require.NoError(t, err)
	assert.Empty(t, fs)
}

// TestComponentManifestCollectQualifiesAlpineEcosystemForOSV pins the fix for
// the Task 9 defect: OSV.dev requires Alpine ecosystem strings to be
// release-qualified (e.g. "Alpine:v3.22"), not the bare "Alpine" family name.
// The bare value returns zero results for every package, silently breaking the
// OSV-first matching strategy. Collect must rewrite the per-entry "Alpine"
// ecosystem to the qualified branch before handing it to QueryOSV, without
// touching any Finding field.
func TestComponentManifestCollectQualifiesAlpineEcosystemForOSV(t *testing.T) {
	manifest := []byte(`{"ref":"main","commit":"abc","groups":{"Other":{"openssl":"3.6.3"}}}`)
	var gotEcosystem string
	c := ComponentManifest{
		FetchManifest: func() ([]byte, error) { return manifest, nil },
		Components: map[string]config.HadronComponentEntry{
			"openssl": {OSV: &config.HadronOSVSource{Ecosystem: "Alpine", Package: "openssl"}},
		},
		QueryOSV: func(ecosystem, pkg, version string) ([]byte, error) {
			gotEcosystem = ecosystem
			return []byte(`{"vulns":[]}`), nil
		},
	}
	repo := state.Repo{Repo: "kairos-io/hadron", Artifacts: []state.Artifact{{Type: "component-manifest"}}}
	_, err := c.Collect(repo)
	require.NoError(t, err)
	assert.Equal(t, "Alpine:v3.22", gotEcosystem)
}

func TestComponentManifestCollectSkipsMarkedPackages(t *testing.T) {
	manifest := []byte(`{"ref":"main","commit":"abc","groups":{"Other":{"mussel":"deadbeef"}}}`)
	c := ComponentManifest{
		FetchManifest: func() ([]byte, error) { return manifest, nil },
		Components:    map[string]config.HadronComponentEntry{"mussel": {Skip: true}},
		QueryOSV:      func(string, string, string) ([]byte, error) { t.Fatal("should not query a skipped package"); return nil, nil },
	}
	repo := state.Repo{Repo: "kairos-io/hadron", Artifacts: []state.Artifact{{Type: "component-manifest"}}}
	fs, err := c.Collect(repo)
	require.NoError(t, err)
	assert.Empty(t, fs)
}
