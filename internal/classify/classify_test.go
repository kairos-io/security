package classify

import (
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/state"
)

func TestApply(t *testing.T) {
	pol := config.CVEPolicy{AcceptedComponents: map[string]config.AcceptedComponent{
		"openssl-fips": {Reason: "FIPS pinned"},
	}}
	in := []state.Finding{
		{ID: "a", Package: "openssl-fips", CurrentVersion: "3.1.2", FixedVersion: "3.5.7"}, // accepted wins
		{ID: "b", Package: "glib", CurrentVersion: "2.86.2", FixedVersion: "2.66.6"},       // already-fixed
		{ID: "c", Package: "openssl", CurrentVersion: "3.1.2", FixedVersion: "3.3.1"},      // actionable
		{ID: "d", Package: "curl", CurrentVersion: "8.5.0"},                                // no fixed -> actionable
		{ID: "e", Package: "musl", CurrentVersion: "1.2.5", FixedVersion: "0"},             // "0" placeholder -> actionable, NOT already-fixed
	}
	out := Apply(in, pol)
	byID := map[string]state.Finding{}
	for _, f := range out {
		byID[f.ID] = f
	}
	if byID["a"].Class != "informational" || byID["a"].ClassReason == "" {
		t.Errorf("openssl-fips should be accepted informational: %+v", byID["a"])
	}
	if byID["b"].Class != "informational" || byID["b"].ClassReason != "already-fixed" {
		t.Errorf("glib should be already-fixed: %+v", byID["b"])
	}
	if byID["c"].Class != "" {
		t.Errorf("openssl should stay actionable: %+v", byID["c"])
	}
	if byID["d"].Class != "" {
		t.Errorf("no-fixed should stay actionable: %+v", byID["d"])
	}
	// "0" fixed is the OSV Alpine placeholder for "not fixed yet". Must NOT
	// be treated as already-fixed or every unpatched Alpine CVE would silently
	// vanish into the informational section.
	if byID["e"].Class != "" || byID["e"].ClassReason != "" {
		t.Errorf(`fixed "0" placeholder must stay actionable, got %+v`, byID["e"])
	}
}

// TestApplyGovulncheckFindingStaysActionable exercises the call site for a
// sourceCVE finding exactly as collect.SourceCVE builds one: govulncheck
// writes the Go module version with its "v" ("v0.30.0") and the stdlib version
// with its "go" ("go1.25.1"), while the advisory's fixed version is bare
// ("0.33.0"). Ordering the prefixed version above the bare one classified
// every reachable Go CVE as already-fixed, which drops it out of the
// dashboard's actionable counts and into the informational section.
func TestApplyGovulncheckFindingStaysActionable(t *testing.T) {
	in := []state.Finding{
		// The fixture in collect/source_test.go, verbatim.
		{ID: "net", Type: "sourceCVE", Package: "golang.org/x/net",
			CurrentVersion: "v0.30.0", FixedVersion: "0.33.0"},
		{ID: "std", Type: "sourceCVE", Package: "stdlib",
			CurrentVersion: "go1.25.1", FixedVersion: "1.25.3"},
		// Genuinely past the fix: must still be filed as already-fixed.
		{ID: "past", Type: "sourceCVE", Package: "golang.org/x/net",
			CurrentVersion: "v0.34.0", FixedVersion: "0.33.0"},
	}
	byID := map[string]state.Finding{}
	for _, f := range Apply(in, config.CVEPolicy{}) {
		byID[f.ID] = f
	}
	for _, id := range []string{"net", "std"} {
		if got := byID[id]; got.Class != "" {
			t.Errorf("%s: a vulnerable Go module must stay actionable, got class=%q reason=%q",
				id, got.Class, got.ClassReason)
		}
	}
	if got := byID["past"]; got.ClassReason != "already-fixed" {
		t.Errorf("past: want already-fixed, got class=%q reason=%q", got.Class, got.ClassReason)
	}
}
