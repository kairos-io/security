package config

// AcceptedComponent marks a whole package as accepted/pinned risk: all of its
// CVEs are classified informational (separated from the dashboard counts).
type AcceptedComponent struct {
	Reason string `yaml:"reason"`
}

// CVEPolicy is the parsed cve-policy.yaml. A zero value (e.g. from a missing
// file) is a valid empty policy.
type CVEPolicy struct {
	AcceptedComponents map[string]AcceptedComponent `yaml:"accepted-components"`
}

// Accepted reports whether pkg is an accepted component and its reason.
func (p CVEPolicy) Accepted(pkg string) (string, bool) {
	a, ok := p.AcceptedComponents[pkg]
	return a.Reason, ok
}

// LoadCVEPolicy reads cve-policy.yaml. A missing file yields an empty policy
// (readYAML treats a not-exist path as the zero value with no error).
func LoadCVEPolicy(path string) (CVEPolicy, error) {
	var cfg CVEPolicy
	return cfg, readYAML(path, &cfg)
}
