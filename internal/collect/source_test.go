package collect

import (
	"strings"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Two govulncheck JSON-line messages: one "osv" (advisory metadata) and one
// "finding" referencing it at module scope.
const govulnJSON = `
{"osv":{"id":"GO-2025-1234","aliases":["CVE-2025-1234"],"summary":"x/net flaw","affected":[{"package":{"name":"golang.org/x/net"},"ranges":[{"type":"SEMVER","events":[{"introduced":"0"},{"fixed":"0.33.0"}]}]}]}}
{"finding":{"osv":"GO-2025-1234","trace":[{"module":"golang.org/x/net","version":"v0.30.0","function":"Parse"}]}}
`

func TestSourceCVEParse(t *testing.T) {
	nowFn = func() string { return "2026-06-19" }
	defer func() { nowFn = defaultNow }()

	c := SourceCVE{Runner: func(state.Repo) ([]byte, error) { return []byte(govulnJSON), nil }}
	fs, err := c.Collect(state.Repo{Repo: "kairos-io/immucore"})
	require.NoError(t, err)
	require.Len(t, fs, 1)
	f := fs[0]
	assert.Equal(t, "sourceCVE", f.Type)
	assert.Equal(t, "CVE-2025-1234", f.CVEID)
	assert.Equal(t, "golang.org/x/net", f.Package)
	assert.Equal(t, "v0.30.0", f.CurrentVersion)
	assert.Equal(t, "0.33.0", f.FixedVersion)
	assert.Equal(t, "go", f.Ecosystem)
	assert.Equal(t, FindingID("kairos-io/immucore", "sourceCVE", "CVE-2025-1234", "golang.org/x/net"), f.ID)
	assert.Equal(t, "2026-06-19", f.FirstSeen)
}

func TestSourceCVEReachabilityAndSeverity(t *testing.T) {
	// One reachable HIGH finding (trace has a function) and one non-reachable
	// finding (no function) for a different module — only the reachable one survives.
	lines := []string{
		`{"osv":{"id":"GO-2024-1","summary":"reachable bug","aliases":["CVE-2024-1"],"database_specific":{"severity":"HIGH"},"affected":[{"package":{"name":"example.com/m"},"ranges":[{"events":[{"fixed":"1.2.3"}]}]}]}}`,
		`{"osv":{"id":"GO-2024-2","summary":"imported only","database_specific":{"severity":"LOW"}}}`,
		`{"finding":{"osv":"GO-2024-1","trace":[{"module":"example.com/m","version":"1.0.0","function":"Vuln"}]}}`,
		`{"finding":{"osv":"GO-2024-2","trace":[{"module":"example.com/other","version":"2.0.0"}]}}`,
	}
	c := SourceCVE{Runner: func(state.Repo) ([]byte, error) { return []byte(strings.Join(lines, "\n")), nil }}
	out, err := c.Collect(state.Repo{Repo: "o/r"})
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "example.com/m", out[0].Package)
	assert.Equal(t, "high", out[0].Severity)
	assert.Equal(t, "1.2.3", out[0].FixedVersion)
	// The advisory URL is built from the GO id (not the CVE alias) so it
	// resolves on pkg.go.dev/vuln, which only serves GO-… paths.
	assert.Equal(t, "https://pkg.go.dev/vuln/GO-2024-1", out[0].URL)
}

func TestSeverityFromOSV(t *testing.T) {
	assert.Equal(t, "critical", severityFromOSV("CRITICAL"))
	assert.Equal(t, "high", severityFromOSV("HIGH"))
	assert.Equal(t, "medium", severityFromOSV("MODERATE"))
	assert.Equal(t, "medium", severityFromOSV("MEDIUM"))
	assert.Equal(t, "low", severityFromOSV("LOW"))
	assert.Equal(t, "high", severityFromOSV("")) // reachable default
}
