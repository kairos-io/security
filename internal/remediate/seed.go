package remediate

import (
	"fmt"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

// ParseSeed turns "owner/repo=package@version" into a synthetic sourceCVE
// finding, so an operator can exercise the remediation pipeline on demand:
//
//	ksec remediate --seed mudler/edgevpn=golang.org/x/net@0.33.0 --dry-run
//
// The fields are exactly those actionable() requires, so the planner turns the
// seed into a real bump intent. Seeds are planning-time only (never persisted).
func ParseSeed(spec string) (state.Finding, error) {
	repo, rest, ok := strings.Cut(spec, "=")
	if !ok || repo == "" {
		return state.Finding{}, fmt.Errorf("seed %q: want owner/repo=package@version", spec)
	}
	pkg, version, ok := strings.Cut(rest, "@")
	if !ok || pkg == "" || version == "" {
		return state.Finding{}, fmt.Errorf("seed %q: want owner/repo=package@version", spec)
	}
	return state.Finding{
		ID:           "seed:" + repo + "|" + pkg,
		Repo:         repo,
		Type:         "sourceCVE",
		Ecosystem:    "go",
		Package:      pkg,
		FixedVersion: version,
		Severity:     "high",
		Source:       "seed",
		CVEID:        "SEED",
		Title:        "synthetic seed finding (remediation test)",
	}, nil
}
