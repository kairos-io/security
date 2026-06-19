package collect

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"time"

	"github.com/kairos-io/security/internal/state"
)

// Collector gathers raw findings for a single repo.
type Collector interface {
	Name() string
	Collect(state.Repo) ([]state.Finding, error)
}

func defaultNow() string { return time.Now().UTC().Format("2006-01-02") }

// nowFn is overridable in tests.
var nowFn = defaultNow

func Today() string { return nowFn() }

func FindingID(repo, typ, cve, pkg string) string {
	sum := sha256.Sum256([]byte(repo + "|" + typ + "|" + cve + "|" + pkg))
	return hex.EncodeToString(sum[:])
}

func Run(repos []state.Repo, collectors []Collector, prev state.Findings) state.Findings {
	firstSeen := map[string]string{}
	for _, f := range prev.Findings {
		firstSeen[f.ID] = f.FirstSeen
	}

	var res state.Findings
	for _, repo := range repos {
		for _, col := range collectors {
			found, err := col.Collect(repo)
			if err != nil {
				res.Errors = append(res.Errors, state.CollectionError{
					Repo: repo.Repo, Collector: col.Name(), Message: err.Error(),
				})
				continue
			}
			for _, f := range found {
				if fs, ok := firstSeen[f.ID]; ok && fs != "" {
					f.FirstSeen = fs
				}
				res.Findings = append(res.Findings, f)
			}
		}
	}
	sort.Slice(res.Findings, func(i, j int) bool { return res.Findings[i].ID < res.Findings[j].ID })
	sort.Slice(res.Errors, func(i, j int) bool {
		if res.Errors[i].Repo != res.Errors[j].Repo {
			return res.Errors[i].Repo < res.Errors[j].Repo
		}
		return res.Errors[i].Collector < res.Errors[j].Collector
	})
	return res
}
