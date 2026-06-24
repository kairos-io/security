package review

import (
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
