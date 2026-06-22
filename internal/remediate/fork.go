package remediate

import "github.com/kairos-io/security/internal/state"

// ForkByKind returns a predicate that reports whether a repo should be
// contributed to via a fork (external repos the bot can't push to directly).
// Org and unknown repos push direct.
func ForkByKind(repos []state.Repo) func(string) bool {
	external := map[string]bool{}
	for _, r := range repos {
		if r.Kind == "external" {
			external[r.Repo] = true
		}
	}
	return func(repo string) bool { return external[repo] }
}
