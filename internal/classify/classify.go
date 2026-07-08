// Package classify labels findings as actionable (default) or informational
// (separated from the dashboard counts): accepted/pinned components and CVEs we
// are already past (current version >= fixed).
package classify

import (
	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/state"
	"github.com/kairos-io/security/internal/version"
)

// Apply returns findings with Class/ClassReason set. Precedence:
//  1. accepted-component (whole package accepted) — reason from policy.
//  2. already-fixed (both versions parse and current >= fixed).
//  3. actionable (Class left empty).
func Apply(findings []state.Finding, policy config.CVEPolicy) []state.Finding {
	out := make([]state.Finding, len(findings))
	copy(out, findings)
	for i := range out {
		f := &out[i]
		if reason, ok := policy.Accepted(f.Package); ok {
			f.Class = "informational"
			f.ClassReason = "accepted-component: " + reason
			continue
		}
		// "already-fixed" only when we have a real, orderable fixed version.
		// OSV Alpine records sometimes carry "0" as a placeholder meaning
		// "no known fix yet"; treating current >= "0" as fixed would hide
		// unpatched vulns in the informational section.
		if f.CurrentVersion != "" && f.FixedVersion != "" && f.FixedVersion != "0" &&
			version.Compare(f.CurrentVersion, f.FixedVersion) >= 0 {
			f.Class = "informational"
			f.ClassReason = "already-fixed"
			continue
		}
		// actionable: leave Class/ClassReason empty (and clear any stale value)
		f.Class = ""
		f.ClassReason = ""
	}
	return out
}
