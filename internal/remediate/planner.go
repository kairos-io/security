package remediate

import (
	"sort"
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

var sevRank = map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1, "unknown": 0, "": 0}

// actionable reports whether a finding can be auto-bumped. Informational
// findings (accepted components / already-fixed) are separated: even though an
// accepted or already-fixed go finding still carries a FixedVersion, we must
// never open a bump PR for it.
func actionable(f state.Finding) bool {
	return f.Class != "informational" &&
		(f.Type == "sourceCVE" || f.Type == "ghAlert") &&
		f.Ecosystem == "go" && f.Package != "" && f.Package != "stdlib" && f.FixedVersion != ""
}

func key(repo, pkg string) string { return repo + "|" + pkg }

// toolchainKey is the ledger key for a repo's Go toolchain bump intent.
func toolchainKey(repo string) string { return repo + "|go-toolchain" }

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

func Plan(c state.Correlated, ledger state.Ledger, prsByRepo map[string][]ghclient.PullRequest, graph *DepGraph, maxNew int) ([]Intent, int) {
	var intents []Intent

	// 1) Reconcile every existing ledger entry.
	for i := range ledger.Entries {
		e := &ledger.Entries[i]
		intents = append(intents, Intent{Type: IntentReconcile, Key: e.Key, Repo: e.Repo, Entry: e})
	}

	// 2) Collapse ALL actionable findings into one target per repo+package
	// (highest fixed version + worst severity). We do NOT skip keys already in
	// the ledger here: the skip decision is made per state/version below.
	type target struct{ repo, pkg, to, sev string }
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

	// 3) Decide per target. Targets already covered by one of our live PRs are
	// skipped (reconcile handles them). Otherwise: adopt an external PR if one
	// addresses it, else mark it a gap to open.
	// Cascade PRs share the maxNew cap with direct opens. Declared here so the
	// per-target loop below can also append supersede intents into the pool.
	type newPR struct {
		intent Intent
		sev    string
	}
	var pool []newPR

	var openKeys []string
	for k, t := range targets {
		if e, ok := ledger.ByKey(k); ok {
			// A ksec entry awaiting a human (build failed, errored, or explicitly
			// flagged) is terminal — do not re-adopt/re-open/re-supersede it.
			if e.Source == "ksec" && (e.NeedsHuman || e.State == "build-failed" || e.State == "error") {
				continue
			}
			// A conflicted adopted PR we can't rebase: supersede it with our own.
			if e.Source != "ksec" && e.Blocked == "upstream-conflict" && e.State == "open" {
				pool = append(pool, newPR{
					intent: Intent{Type: IntentSupersede, Key: k, Repo: t.repo, Package: t.pkg, Severity: t.sev,
						Bump: state.Bump{Package: t.pkg, To: t.to}, PRNumber: e.PRNumber, PRURL: e.PRURL},
					sev: t.sev,
				})
				continue
			}
			if e.State == "open" || e.State == "conflicted" {
				continue
			}
			if (e.State == "merged" || e.State == "closed") && compareVersions(e.Bump.To, t.to) >= 0 {
				continue
			}
		}
		if pr, source, ok := MatchPR(t.pkg, t.to, prsByRepo[t.repo]); ok && source != "ksec" {
			intents = append(intents, Intent{
				Type: IntentAdopt, Key: k, Repo: t.repo, Package: t.pkg, Severity: t.sev,
				Bump: state.Bump{Package: t.pkg, To: t.to}, PRNumber: pr.Number, PRURL: pr.URL, Source: source,
			})
			continue
		}
		openKeys = append(openKeys, k)
	}

	// Repin: every pseudo cascade entry is a repin candidate (the executor
	// decides whether a tag is available yet).
	for i := range ledger.Entries {
		e := &ledger.Entries[i]
		if e.Kind == "cascade" && e.Pseudo && e.State == "open" {
			intents = append(intents, Intent{Type: IntentRepin, Key: e.Key, Repo: e.Repo, Entry: e})
		}
	}

	// Cascade: a merged fix in a first-party module repo means that module's
	// default branch has the fix; bump it in each consumer that isn't already
	// tracked for it.
	for _, k := range openKeys { // direct gaps from the earlier 4a logic
		t := targets[k]
		pool = append(pool, newPR{
			intent: Intent{Type: IntentOpen, Key: k, Repo: t.repo, Package: t.pkg, Severity: t.sev,
				Bump: state.Bump{Package: t.pkg, To: t.to}},
			sev: t.sev,
		})
	}
	if graph != nil {
		seen := map[string]bool{}
		for i := range ledger.Entries {
			e := &ledger.Entries[i]
			mod := graph.ModuleOf(e.Repo)
			if mod == "" || e.State != "merged" {
				continue
			}
			for _, consumer := range graph.Consumers(mod) {
				ck := key(consumer, mod)
				if seen[ck] {
					continue // already cascaded this consumer this run (intra-run dedup)
				}
				if ce, ok := ledger.ByKey(ck); ok {
					// open/conflicted/merged: already cascading or done.
					// closed: a maintainer closed it; don't fight them by re-creating.
					if ce.State == "open" || ce.State == "conflicted" || ce.State == "merged" || ce.State == "closed" {
						continue
					}
				}
				seen[ck] = true
				pool = append(pool, newPR{
					intent: Intent{Type: IntentCascade, Key: ck, Repo: consumer, Package: mod,
						Ref: graph.BranchOf(e.Repo), CascadeFrom: e.Key, Severity: e.Severity},
					sev: e.Severity,
				})
			}
		}
	}

	// Toolchain: a stdlib finding (Package=="stdlib") can't be `go get`-ed; it
	// needs a Go toolchain bump. Collapse stdlib findings into one bump per repo
	// (highest fixed Go version, worst severity) and emit into the same capped
	// pool, skipping repos already covered by a live/closed ledger entry.
	type tc struct{ ver, sev string }
	tcByRepo := map[string]*tc{}
	for _, f := range c.Findings {
		if f.Class == "informational" || f.Package != "stdlib" || f.Ecosystem != "go" || f.FixedVersion == "" {
			continue
		}
		ver := strings.TrimPrefix(f.FixedVersion, "go")
		t := tcByRepo[f.Repo]
		if t == nil {
			tcByRepo[f.Repo] = &tc{ver: ver, sev: f.Severity}
			continue
		}
		t.ver = higherVersion(t.ver, ver)
		if sevRank[f.Severity] > sevRank[t.sev] {
			t.sev = f.Severity
		}
	}
	for repo, t := range tcByRepo {
		if e, ok := ledger.ByKey(toolchainKey(repo)); ok {
			if e.State == "open" || e.State == "conflicted" || e.State == "merged" || e.State == "closed" {
				continue
			}
		}
		pool = append(pool, newPR{
			intent: Intent{Type: IntentToolchain, Key: toolchainKey(repo), Repo: repo,
				ToolchainVersion: t.ver, Severity: t.sev},
			sev: t.sev,
		})
	}

	sort.SliceStable(pool, func(i, j int) bool {
		if sevRank[pool[i].sev] != sevRank[pool[j].sev] {
			return sevRank[pool[i].sev] > sevRank[pool[j].sev]
		}
		return pool[i].intent.Key < pool[j].intent.Key
	})
	deferred := 0
	for n := range pool {
		if n >= maxNew {
			deferred = len(pool) - n
			break
		}
		intents = append(intents, pool[n].intent)
	}
	return intents, deferred
}
