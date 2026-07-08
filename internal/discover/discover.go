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

// FilterArchived drops repos that GitHub reports as archived. An archived repo
// is read-only: it accepts no PRs, patches, or automated fixes and usually
// signals discontinued software, so there is no value in tracking it.
//
// A per-repo lookup error is treated as non-fatal: the repo is kept (fail-safe,
// so a transient GitHub hiccup does not silently drop a live repo) and its slug
// is returned in failed for the caller to surface. dropped lists the archived
// repos that were removed.
func FilterArchived(gh ghclient.GitHub, repos []state.Repo) (kept []state.Repo, dropped, failed []string) {
	for _, r := range repos {
		archived, err := gh.RepoArchived(r.Repo)
		if err != nil {
			failed = append(failed, r.Repo)
			kept = append(kept, r)
			continue
		}
		if archived {
			dropped = append(dropped, r.Repo)
			continue
		}
		kept = append(kept, r)
	}
	return kept, dropped, failed
}

func hasPrefix(s, p string) bool { return len(s) >= len(p) && s[:len(p)] == p }
