package config

import (
	"errors"
	"io/fs"
	"os"

	"github.com/kairos-io/security/internal/state"
	"gopkg.in/yaml.v3"
)

type ReposConfig struct {
	Repos   []state.Repo `yaml:"repos"`
	Exclude []string     `yaml:"exclude"`
}

type ModelCfg struct {
	Name    string `yaml:"name"`
	Gallery string `yaml:"gallery"`
	Quant   string `yaml:"quant"`
}

type LocalAICfg struct {
	Version        string   `yaml:"version"`
	Endpoint       string   `yaml:"endpoint"`
	StartupTimeout string   `yaml:"startupTimeout"`
	Model          ModelCfg `yaml:"model"`
}

type NibCfg struct {
	Version     string  `yaml:"version"`
	Mode        string  `yaml:"mode"`
	Model       string  `yaml:"model"`
	Endpoint    string  `yaml:"endpoint"`
	MaxTokens   int     `yaml:"maxTokens"`
	Temperature float64 `yaml:"temperature"`
}

type ReviewCfg struct {
	Enabled     bool     `yaml:"enabled"`
	AutoApprove bool     `yaml:"autoApprove"`
	MaxPerRun   int      `yaml:"maxPerRun"`
	Notify      []string `yaml:"notify"`
}

// ApplicabilityCfg controls the AI applicability classifier that runs during
// the correlate phase. Defaults derive from localai the same way NibCfg does,
// so a caller only has to flip `enabled: true`.
type ApplicabilityCfg struct {
	Enabled             bool    `yaml:"enabled"`
	Endpoint            string  `yaml:"endpoint"`
	Model               string  `yaml:"model"`
	Temperature         float64 `yaml:"temperature"`
	MaxTokens           int     `yaml:"maxTokens"`
	ConfidenceThreshold string  `yaml:"confidenceThreshold"` // "high" (default) | "medium"
}

type AIConfig struct {
	LocalAI       LocalAICfg       `yaml:"localai"`
	Nib           NibCfg           `yaml:"nib"`
	Review        ReviewCfg        `yaml:"review"`
	Applicability ApplicabilityCfg `yaml:"applicability"`
}

func readYAML[T any](path string, v *T) error {
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil // missing file → zero value, no error
	}
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

func LoadRepos(path string) (ReposConfig, error) {
	var cfg ReposConfig
	return cfg, readYAML(path, &cfg)
}

func LoadAI(path string) (AIConfig, error) {
	var cfg AIConfig
	if err := readYAML(path, &cfg); err != nil {
		return cfg, err
	}
	if v := os.Getenv("LOCALAI_URL"); v != "" {
		cfg.LocalAI.Endpoint = v
	}
	if v := os.Getenv("LOCALAI_MODEL"); v != "" {
		cfg.LocalAI.Model.Name = v
	}
	// nib defaults derive from localai so they cannot drift
	if cfg.Nib.Endpoint == "" {
		cfg.Nib.Endpoint = cfg.LocalAI.Endpoint
	}
	if cfg.Nib.Model == "" {
		cfg.Nib.Model = cfg.LocalAI.Model.Name
	}
	if cfg.Review.MaxPerRun <= 0 {
		cfg.Review.MaxPerRun = 20
	}
	// Applicability defaults follow the nib pattern: derive from localai
	// unless explicitly overridden, so a bare `applicability: {enabled: true}`
	// suffices.
	if cfg.Applicability.Endpoint == "" {
		cfg.Applicability.Endpoint = cfg.LocalAI.Endpoint
	}
	if cfg.Applicability.Model == "" {
		cfg.Applicability.Model = cfg.LocalAI.Model.Name
	}
	if cfg.Applicability.MaxTokens <= 0 {
		cfg.Applicability.MaxTokens = 1024
	}
	if cfg.Applicability.ConfidenceThreshold == "" {
		cfg.Applicability.ConfidenceThreshold = "high"
	}
	return cfg, nil
}
