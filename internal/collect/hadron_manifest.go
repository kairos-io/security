package collect

import (
	"encoding/json"
	"sort"
)

// HadronComponent is one package pin from hadron's published component
// manifest (https://hadron-linux.io/components/main.json).
type HadronComponent struct {
	Group   string
	Package string
	Version string
}

type hadronManifestDoc struct {
	Ref    string                       `json:"ref"`
	Commit string                       `json:"commit"`
	Groups map[string]map[string]string `json:"groups"`
}

// ParseHadronManifest flattens the manifest's grouped package->version maps
// into a sorted (group, then package) slice for deterministic iteration.
func ParseHadronManifest(raw []byte) ([]HadronComponent, error) {
	var m hadronManifestDoc
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	out := make([]HadronComponent, 0, len(m.Groups))
	for group, pkgs := range m.Groups {
		for pkg, version := range pkgs {
			out = append(out, HadronComponent{Group: group, Package: pkg, Version: version})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Group != out[j].Group {
			return out[i].Group < out[j].Group
		}
		return out[i].Package < out[j].Package
	})
	return out, nil
}
