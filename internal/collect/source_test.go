package collect

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Two govulncheck JSON-line messages: one "osv" (advisory metadata) and one
// "finding" referencing it at module scope.
const govulnJSON = `
{"osv":{"id":"GO-2025-1234","aliases":["CVE-2025-1234"],"summary":"x/net flaw","affected":[{"package":{"name":"golang.org/x/net"},"ranges":[{"type":"SEMVER","events":[{"introduced":"0"},{"fixed":"0.33.0"}]}]}]}}
{"finding":{"osv":"GO-2025-1234","trace":[{"module":"golang.org/x/net","version":"v0.30.0"}]}}
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
