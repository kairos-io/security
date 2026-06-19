package discover

import (
	"regexp"
	"sort"
	"strings"
)

// Fixed component → repo slug map for kairos-init Makefile *_VERSION vars.
var makefileComponents = map[string]string{
	"AGENT":                       "kairos-io/kairos-agent",
	"IMMUCORE":                    "kairos-io/immucore",
	"KCRYPT_DISCOVERY_CHALLENGER": "kairos-io/kcrypt-discovery-challenger",
	"PROVIDER_KAIROS":             "kairos-io/provider-kairos",
	"EDGEVPN":                     "mudler/edgevpn",
}

var (
	reMakeVar = regexp.MustCompile(`(?m)^([A-Z_]+)_VERSION\??=`)
	reGoMod   = regexp.MustCompile(`(?m)^\s*(?:require\s+)?github\.com/(kairos-io|mudler|mauromorales)/([A-Za-z0-9._-]+)\s+v`)
)

// ParseDeps returns sorted unique owner/name slugs from the kairos-init
// Makefile and go.mod.
func ParseDeps(makefile, gomod []byte) []string {
	set := map[string]struct{}{}
	for _, m := range reMakeVar.FindAllStringSubmatch(string(makefile), -1) {
		if slug, ok := makefileComponents[m[1]]; ok {
			set[slug] = struct{}{}
		}
	}
	for _, m := range reGoMod.FindAllStringSubmatch(string(gomod), -1) {
		set[m[1]+"/"+strings.TrimSuffix(m[2], "/")] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
