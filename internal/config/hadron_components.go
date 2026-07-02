package config

// HadronOSVSource identifies a package for an OSV.dev query.
type HadronOSVSource struct {
	Ecosystem string `yaml:"ecosystem"`
	Package   string `yaml:"package"`
}

// HadronCPESource identifies a package's NVD CPE vendor/product pair, used as
// a fallback when OSV has no advisory. Left unset until curated (see
// hadron-components.yaml's header comment).
type HadronCPESource struct {
	Vendor  string `yaml:"vendor"`
	Product string `yaml:"product"`
}

// HadronComponentEntry configures how one hadron manifest package is matched
// against CVE sources. Skip==true means this "package" isn't real, versioned,
// CVE-bearing upstream software (e.g. a build-tool commit pin or a data
// bundle) and is never queried.
type HadronComponentEntry struct {
	OSV  *HadronOSVSource `yaml:"osv,omitempty"`
	CPE  *HadronCPESource `yaml:"cpe,omitempty"`
	Skip bool             `yaml:"skip,omitempty"`
}

type HadronComponentsConfig struct {
	Components map[string]HadronComponentEntry `yaml:"components"`
}

func LoadHadronComponents(path string) (HadronComponentsConfig, error) {
	var cfg HadronComponentsConfig
	return cfg, readYAML(path, &cfg)
}
