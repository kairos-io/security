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
}
