package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadHadronComponentsParsesOSVAndSkip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "hadron-components.yaml")
	require.NoError(t, os.WriteFile(p, []byte(`
components:
  openssl:
    osv: {ecosystem: Alpine, package: openssl}
  mussel:
    skip: true
`), 0o644))

	cfg, err := LoadHadronComponents(p)
	require.NoError(t, err)
	require.Contains(t, cfg.Components, "openssl")
	require.NotNil(t, cfg.Components["openssl"].OSV)
	assert.Equal(t, "Alpine", cfg.Components["openssl"].OSV.Ecosystem)
	assert.Equal(t, "openssl", cfg.Components["openssl"].OSV.Package)
	assert.False(t, cfg.Components["openssl"].Skip)

	assert.True(t, cfg.Components["mussel"].Skip)
	assert.Nil(t, cfg.Components["mussel"].OSV)
}

func TestLoadHadronComponentsMissingFile(t *testing.T) {
	cfg, err := LoadHadronComponents(filepath.Join(t.TempDir(), "nope.yaml"))
	require.NoError(t, err)
	assert.Empty(t, cfg.Components)
}
