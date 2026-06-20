package remediate

import (
	"regexp"
	"sort"

	"github.com/kairos-io/security/internal/state"
)

type DepGraph struct {
	moduleOf  map[string]string   // repo -> module import path
	repoOf    map[string]string   // module import path -> repo
	consumers map[string][]string // module import path -> repos requiring it
	branchOf  map[string]string   // repo -> default branch
}

var (
	reModuleLine  = regexp.MustCompile(`(?m)^module\s+(\S+)`)
	reRequireLine = regexp.MustCompile(`(?m)^\s*(?:require\s+)?(github\.com/\S+)\s+v\S+`)
)

// BuildGraph parses each tracked repo's go.mod to map module<->repo and to find,
// for each first-party module, the tracked repos that require it.
func BuildGraph(repos []state.Repo, gomodByRepo map[string][]byte) *DepGraph {
	g := &DepGraph{
		moduleOf:  map[string]string{},
		repoOf:    map[string]string{},
		consumers: map[string][]string{},
		branchOf:  map[string]string{},
	}
	for _, r := range repos {
		b := r.Branch
		if b == "" {
			b = "main"
		}
		g.branchOf[r.Repo] = b
		if m := reModuleLine.FindSubmatch(gomodByRepo[r.Repo]); m != nil {
			mod := string(m[1])
			g.moduleOf[r.Repo] = mod
			g.repoOf[mod] = r.Repo
		}
	}
	// Second pass: requires, keeping only first-party modules (those we map to a repo).
	for repo, mod := range g.moduleOf {
		_ = mod
		for _, m := range reRequireLine.FindAllSubmatch(gomodByRepo[repo], -1) {
			req := string(m[1])
			if _, ok := g.repoOf[req]; ok {
				g.consumers[req] = append(g.consumers[req], repo)
			}
		}
	}
	for k := range g.consumers {
		sort.Strings(g.consumers[k])
	}
	return g
}

func (g *DepGraph) ModuleOf(repo string) string { return g.moduleOf[repo] }
func (g *DepGraph) RepoOf(module string) (string, bool) {
	r, ok := g.repoOf[module]
	return r, ok
}
func (g *DepGraph) Consumers(module string) []string { return g.consumers[module] }
func (g *DepGraph) BranchOf(repo string) string {
	if b, ok := g.branchOf[repo]; ok {
		return b
	}
	return "main"
}
