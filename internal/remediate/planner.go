package remediate

import (
	"sort"

	"github.com/kairos-io/security/internal/state"
)

var sevRank = map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1, "unknown": 0, "": 0}

// actionable reports whether a finding can be auto-bumped.
func actionable(f state.Finding) bool {
	return (f.Type == "sourceCVE" || f.Type == "ghAlert") &&
		f.Ecosystem == "go" && f.Package != "" && f.FixedVersion != ""
}

func key(repo, pkg string) string { return repo + "|" + pkg }

// higherVersion returns the "greater" of two version strings. We avoid a full
// semver parser: trim a leading 'v' and compare dotted-numeric segments,
// ignoring any pre-release / build metadata. So "1.2.0" and "1.2.0-rc1"
// compare equal.
func higherVersion(a, b string) string {
	if compareVersions(a, b) >= 0 {
		return a
	}
	return b
}

func compareVersions(a, b string) int {
	na, nb := splitVer(a), splitVer(b)
	for i := 0; i < len(na) || i < len(nb); i++ {
		var x, y int
		if i < len(na) {
			x = na[i]
		}
		if i < len(nb) {
			y = nb[i]
		}
		if x != y {
			if x < y {
				return -1
			}
			return 1
		}
	}
	return 0
}

func splitVer(s string) []int {
	if len(s) > 0 && (s[0] == 'v' || s[0] == 'V') {
		s = s[1:]
	}
	var out []int
	cur, has := 0, false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			cur = cur*10 + int(c-'0')
			has = true
		} else if c == '.' {
			out = append(out, cur)
			cur, has = 0, false
		} else {
			break // stop at pre-release / build metadata
		}
	}
	if has || len(out) == 0 {
		out = append(out, cur)
	}
	return out
}

func Plan(c state.Correlated, ledger state.Ledger, maxNew int) ([]Intent, int) {
	var intents []Intent

	// 1) Reconcile every existing ledger entry.
	for i := range ledger.Entries {
		e := &ledger.Entries[i]
		intents = append(intents, Intent{Type: IntentReconcile, Key: e.Key, Repo: e.Repo, Entry: e})
	}

	// 2) Collapse ALL actionable findings into one target per repo+package
	// (highest fixed version + worst severity). We do NOT skip keys already in
	// the ledger here: the skip decision is made per state/version below.
	type target struct {
		repo, pkg, to, sev string
	}
	targets := map[string]*target{}
	for _, f := range c.Findings {
		if !actionable(f) {
			continue
		}
		k := key(f.Repo, f.Package)
		t := targets[k]
		if t == nil {
			targets[k] = &target{repo: f.Repo, pkg: f.Package, to: f.FixedVersion, sev: f.Severity}
			continue
		}
		t.to = higherVersion(t.to, f.FixedVersion)
		if sevRank[f.Severity] > sevRank[t.sev] {
			t.sev = f.Severity
		}
	}

	// 3) Decide which targets actually need a new IntentOpen, based on the
	// existing ledger entry's STATE and version rather than mere presence:
	//   - open/conflicted     -> skip (a live PR is already maintained; the
	//                            IntentReconcile above covers it).
	//   - merged/closed at an  -> skip (already addressed at >= our version).
	//     equal/higher version
	//   - everything else (no entry; planned/error/build-failed; or
	//     merged/closed at a LOWER version) -> emit IntentOpen.
	keys := make([]string, 0, len(targets))
	for k, t := range targets {
		if e, ok := ledger.ByKey(k); ok {
			if e.State == "open" || e.State == "conflicted" {
				continue
			}
			if (e.State == "merged" || e.State == "closed") && compareVersions(e.Bump.To, t.to) >= 0 {
				continue
			}
		}
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		ti, tj := targets[keys[i]], targets[keys[j]]
		if sevRank[ti.sev] != sevRank[tj.sev] {
			return sevRank[ti.sev] > sevRank[tj.sev]
		}
		return keys[i] < keys[j]
	})

	deferred := 0
	for n, k := range keys {
		if n >= maxNew {
			deferred = len(keys) - n
			break
		}
		t := targets[k]
		intents = append(intents, Intent{
			Type: IntentOpen, Key: k, Repo: t.repo, Package: t.pkg, Severity: t.sev,
			Bump: state.Bump{Package: t.pkg, To: t.to},
		})
	}
	return intents, deferred
}
