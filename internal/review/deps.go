package review

import (
	"fmt"
	"regexp"
	"strings"
)

type DepBump struct{ Module, From, To string }

var reModLine = regexp.MustCompile(`^[+-]\s+(\S+)\s+v(\S+)`)

// rePseudo matches a Go pseudo-version's trailing "-<14-digit-timestamp>-<sha>".
var rePseudo = regexp.MustCompile(`[-.]\d{14}-([0-9a-f]{12,})$`)

// compareRef maps a (v-stripped) module version to the ref to compare against
// upstream: a pseudo-version compares by its embedded commit SHA, a real
// release compares by its "v"-prefixed tag.
func compareRef(version string) string {
	if m := rePseudo.FindStringSubmatch(version); m != nil {
		return m[1]
	}
	return "v" + version
}

// parseBumps extracts {module, from, to} from a PR's go.mod diff by pairing the
// "-" old and "+" new version lines for the same module.
func parseBumps(diff []byte) []DepBump {
	from := map[string]string{}
	to := map[string]string{}
	var order []string
	for _, line := range strings.Split(string(diff), "\n") {
		m := reModLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		mod, ver := m[1], strings.Fields(m[2])[0] // tolerate trailing " // indirect"
		if strings.HasPrefix(line, "-") {
			from[mod] = ver
		} else {
			if _, seen := to[mod]; !seen {
				order = append(order, mod)
			}
			to[mod] = ver
		}
	}
	var out []DepBump
	for _, mod := range order {
		if f, ok := from[mod]; ok && f != to[mod] {
			out = append(out, DepBump{Module: mod, From: f, To: to[mod]})
		}
	}
	return out
}

// moduleRepo maps a Go module path to a GitHub "owner/name", handling the
// common hosts. Unresolvable (vanity) paths return ok=false (degrade).
func moduleRepo(mod string) (string, bool) {
	parts := strings.Split(mod, "/")
	switch {
	case parts[0] == "github.com" && len(parts) >= 3:
		return parts[1] + "/" + parts[2], true
	case strings.HasPrefix(mod, "golang.org/x/") && len(parts) >= 3:
		return "golang/" + parts[2], true
	case strings.HasPrefix(mod, "k8s.io/") && len(parts) >= 2:
		return "kubernetes/" + parts[1], true
	case strings.HasPrefix(mod, "sigs.k8s.io/") && len(parts) >= 2:
		return "kubernetes-sigs/" + parts[1], true
	}
	return "", false
}

type CompareRef struct{ Repo, Base, Head, Label string }

const maxCompares = 5

// reCompare requires a scheme + github.com as a real domain segment, so a
// look-alike host like "attacker-github.com" can't smuggle in an arbitrary
// repo target. (The host is otherwise discarded — fetches always go to
// api.github.com — but anchoring it is cheap defense in depth.)
// The separator is two-or-three dots: renovate uses the three-dot form for
// version bumps and the two-dot range form for action/image *digest* bumps
// (e.g. compare/<sha>..<sha>). CompareDiff always rebuilds it as three-dot.
var reCompare = regexp.MustCompile(`https?://(?:[\w-]+\.)*github\.com/([\w.\-]+)/([\w.\-]+)/compare/([\w.\-+/@]+?)\.\.\.?([\w.\-+/@]+)`)

// parseCompareURLs extracts GitHub compare links (owner/repo + base...head) from
// a renovate/dependabot PR body. Base/head are taken verbatim — the bot already
// resolved the repo's real tag form (incl. monorepos). Works for any ecosystem.
func parseCompareURLs(body string) []CompareRef {
	var out []CompareRef
	seen := map[string]bool{}
	for _, m := range reCompare.FindAllStringSubmatch(body, -1) {
		repo := m[1] + "/" + m[2]
		ref := CompareRef{Repo: repo, Base: m[3], Head: m[4],
			Label: fmt.Sprintf("%s %s..%s (PR body)", repo, m[3], m[4])}
		k := ref.Repo + "|" + ref.Base + "|" + ref.Head
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, ref)
	}
	return out
}

// compareTargets unifies upstream comparisons from the Go go.mod bumps and the
// PR-body compare links, deduped by repo|base|head and capped at maxCompares.
func compareTargets(diff []byte, body string) []CompareRef {
	var out []CompareRef
	seen := map[string]bool{}
	add := func(r CompareRef) {
		if r.Repo == "" {
			return
		}
		k := r.Repo + "|" + r.Base + "|" + r.Head
		if seen[k] {
			return
		}
		seen[k] = true
		out = append(out, r)
	}
	for _, b := range parseBumps(diff) {
		if repo, ok := moduleRepo(b.Module); ok {
			add(CompareRef{Repo: repo, Base: compareRef(b.From), Head: compareRef(b.To),
				Label: fmt.Sprintf("%s %s→%s", b.Module, b.From, b.To)})
		}
	}
	for _, c := range parseCompareURLs(body) {
		add(c)
	}
	if len(out) > maxCompares {
		out = out[:maxCompares]
	}
	return out
}
