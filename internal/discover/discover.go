package discover

import (
	"sort"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

func Run(gh ghclient.GitHub, cfg config.ReposConfig, org, initRepo, initRef string) ([]state.Repo, error) {
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

	// Config additions / metadata overrides (matched by .Repo).
	for _, r := range cfg.Repos {
		existing := merged[r.Repo]
		existing.Repo = r.Repo
		if r.Kind != "" {
			existing.Kind = r.Kind
		}
		if r.Branch != "" {
			existing.Branch = r.Branch
		}
		if r.Criticality != "" {
			existing.Criticality = r.Criticality
		}
		if len(r.Artifacts) > 0 {
			existing.Artifacts = r.Artifacts
		}
		merged[r.Repo] = existing
	}

	for _, ex := range cfg.Exclude {
		delete(merged, ex)
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
