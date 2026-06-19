package collect

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const trivyJSON = `{"Results":[{"Target":"x","Vulnerabilities":[
{"VulnerabilityID":"CVE-2025-9999","PkgName":"openssl","InstalledVersion":"1.1.1","FixedVersion":"1.1.1w","Severity":"CRITICAL","PrimaryURL":"https://x/CVE-2025-9999","Title":"openssl flaw"}]}]}`

func TestImageCVEParse(t *testing.T) {
	nowFn = func() string { return "2026-06-19" }
	defer func() { nowFn = defaultNow }()

	c := ImageCVE{Runner: func(string) ([]byte, error) { return []byte(trivyJSON), nil }}
	repo := state.Repo{Repo: "kairos-io/kairos", Artifacts: []state.Artifact{{Type: "image", Ref: "quay.io/kairos/x:latest"}}}
	fs, err := c.Collect(repo)
	require.NoError(t, err)
	require.Len(t, fs, 1)
	f := fs[0]
	assert.Equal(t, "imageCVE", f.Type)
	assert.Equal(t, "CVE-2025-9999", f.CVEID)
	assert.Equal(t, "openssl", f.Package)
	assert.Equal(t, "critical", f.Severity)
	assert.Equal(t, "1.1.1w", f.FixedVersion)
}
