package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadReposParsesAndDefaultsMissing(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "repos.yaml")
	require.NoError(t, os.WriteFile(p, []byte(`
repos:
  - repo: mudler/edgevpn
    kind: external
    branch: master
    criticality: high
    artifacts:
      - type: go
        modpath: .
exclude:
  - kairos-io/some-archive
`), 0o644))

	cfg, err := LoadRepos(p)
	require.NoError(t, err)
	require.Len(t, cfg.Repos, 1)
	assert.Equal(t, "mudler/edgevpn", cfg.Repos[0].Repo)
	assert.Equal(t, "external", cfg.Repos[0].Kind)
	assert.Equal(t, []string{"kairos-io/some-archive"}, cfg.Exclude)

	missing, err := LoadRepos(filepath.Join(dir, "nope.yaml"))
	require.NoError(t, err)
	assert.Empty(t, missing.Repos)
}

func TestLoadAIAppliesEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ai.yaml")
	require.NoError(t, os.WriteFile(p, []byte(`
localai:
  endpoint: http://localhost:8080
  model:
    name: base-model
nib:
  mode: yolo
`), 0o644))
	t.Setenv("LOCALAI_URL", "http://override:9000")
	t.Setenv("LOCALAI_MODEL", "override-model")

	cfg, err := LoadAI(p)
	require.NoError(t, err)
	assert.Equal(t, "http://override:9000", cfg.LocalAI.Endpoint)
	assert.Equal(t, "override-model", cfg.LocalAI.Model.Name)
	assert.Equal(t, "yolo", cfg.Nib.Mode)
}

func TestLoadAIReview(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ai.yaml")
	require.NoError(t, os.WriteFile(p, []byte(`
localai:
  endpoint: http://localhost:8080
  model:
    name: m
review:
  enabled: true
  autoApprove: true
  notify: ["@team"]
`), 0o644))
	cfg, err := LoadAI(p)
	require.NoError(t, err)
	assert.True(t, cfg.Review.Enabled)
	assert.True(t, cfg.Review.AutoApprove)
	assert.Equal(t, []string{"@team"}, cfg.Review.Notify)
	assert.Equal(t, 20, cfg.Review.MaxPerRun) // defaulted
}
