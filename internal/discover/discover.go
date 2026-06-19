package discover

import (
	"sort"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

// Normalize is the per-run path: it applies defaults, drops excluded repos and
// sorts the curated repo list. It performs no I/O and never talks to GitHub.
func Normalize(cfg config.ReposConfig) []state.Repo {
	exclude := map[string]bool{}
	for _, ex := range cfg.Exclude {
		exclude[ex] = true
	}

	out := make([]state.Repo, 0, len(cfg.Repos))
	for _, r := range cfg.Repos {
		if exclude[r.Repo] {
			continue
		}
		if r.Kind == "" {
			if hasPrefix(r.Repo, "kairos-io/") {
				r.Kind = "org"
			} else {
				r.Kind = "external"
			}
		}
		if r.Branch == "" {
			r.Branch = "main"
		}
		if r.Criticality == "" {
			r.Criticality = "medium"
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Repo < out[j].Repo })
	return out
}

// SeedFromOrg enumerates a GitHub org and the kairos-init dependency graph to
// produce a candidate repo list. It is an opt-in, one-time seed helper used to
// bootstrap repos.yaml; the per-run pipeline does not call it.
func SeedFromOrg(gh ghclient.GitHub, org, initRepo, initRef string) ([]state.Repo, error) {
	orgRepos, err := gh.ListOrgRepos(org)
	if err != nil {
		return nil, err
	}
	orgSet := map[string]bool{}
	merged := map[string]state.Repo{}
	for _, r := range orgRepos {
		orgSet[r] = true
		merged[r] = state.Repo{Repo: r, Kind: "org"}
	}

	makefile, _ := gh.GetFile(initRepo, "Makefile", initRef)
	gomod, _ := gh.GetFile(initRepo, "go.mod", initRef)
	for _, slug := range ParseDeps(makefile, gomod) {
		if _, ok := merged[slug]; !ok {
			kind := "external"
			if orgSet[slug] {
				kind = "org"
			} else if hasPrefix(slug, "kairos-io/") {
				kind = "dep"
			}
			merged[slug] = state.Repo{Repo: slug, Kind: kind}
		}
	}

	out := make([]state.Repo, 0, len(merged))
	for _, r := range merged {
		if r.Branch == "" {
			r.Branch = "main"
		}
		if r.Criticality == "" {
			r.Criticality = "medium"
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Repo < out[j].Repo })
	return out, nil
}

func hasPrefix(s, p string) bool { return len(s) >= len(p) && s[:len(p)] == p }
