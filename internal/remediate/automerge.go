package remediate

import "github.com/kairos-io/security/internal/ghclient"

// ShouldAutomerge reports whether an addressing PR is safe to merge: it must be
// mergeable, have passing checks, and not be blocked by a requested-changes
// review.
func ShouldAutomerge(s ghclient.PRStatus) bool {
	return s.Mergeable && s.ChecksPassing && s.ReviewDecision != "CHANGES_REQUESTED"
}
