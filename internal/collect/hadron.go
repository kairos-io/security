package collect

import (
	"fmt"
	"os"
	"sort"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/state"
)

// osvAlpineBranch is the OSV.dev Alpine release branch ksec queries for
// component CVE matching. OSV.dev requires Alpine ecosystem strings to be
// release-qualified (e.g. "Alpine:v3.22"), not the bare "Alpine" family name —
// bump this periodically as Alpine branches age out of OSV's tracked set.
const osvAlpineBranch = "Alpine:v3.22"

// ComponentManifest collects CVEs for kairos-io/hadron's published
// component manifest: OSV.dev (Alpine ecosystem) first, falling back to an
// NVD CPE-match query for packages OSV has no advisory for.
type ComponentManifest struct {
	FetchManifest func() ([]byte, error)
	Components    map[string]config.HadronComponentEntry
	QueryOSV      OSVQueryFunc
	QueryNVD      NVDQueryFunc
}

func (ComponentManifest) Name() string { return "componentManifest" }

func (c ComponentManifest) Collect(repo state.Repo) ([]state.Finding, error) {
	hasManifest := false
	for _, a := range repo.Artifacts {
		if a.Type == "component-manifest" {
			hasManifest = true
		}
	}
	if !hasManifest {
		return nil, nil
	}

	raw, err := c.FetchManifest()
	if err != nil {
		return nil, err
	}
	components, err := ParseHadronManifest(raw)
	if err != nil {
		return nil, err
	}

	type hit struct {
		cveID, severity, fixed, title, url, source string
		details, affected                          string
	}

	out := map[string]state.Finding{}
	for _, comp := range components {
		entry, ok := c.Components[comp.Package]
		if !ok || entry.Skip {
			continue
		}

		var hits []hit
		if entry.OSV != nil && c.QueryOSV != nil {
			ecosystem := entry.OSV.Ecosystem
			if ecosystem == "Alpine" {
				ecosystem = osvAlpineBranch
			}
			results, err := QueryOSV(c.QueryOSV, ecosystem, entry.OSV.Package, comp.Version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "componentManifest: OSV query failed for %s@%s: %v\n", comp.Package, comp.Version, err)
			} else {
				for _, r := range results {
					hits = append(hits, hit{r.CVEID, r.Severity, r.FixedVersion, r.Title, r.URL, "osv", r.Details, r.AffectedRanges})
				}
			}
		}
		if len(hits) == 0 && entry.CPE != nil && c.QueryNVD != nil {
			results, err := QueryNVD(c.QueryNVD, entry.CPE.Vendor, entry.CPE.Product, comp.Version)
			if err != nil {
				fmt.Fprintf(os.Stderr, "componentManifest: NVD query failed for %s@%s: %v\n", comp.Package, comp.Version, err)
			} else {
				for _, r := range results {
					if r.VersionEndExcluding == "" {
						// QueryNVD emits one NVDResult per NVD "vulnerabilities[]"
						// entry regardless of whether any cpeMatch was actually
						// vulnerable==true for the queried CPE (NVD configurations
						// legitimately include vulnerable:false entries for
						// platform/AND conditions). VersionEndExcluding is only
						// populated from a vulnerable==true match, so an empty
						// value here means this CVE's configuration didn't
						// actually implicate the queried package — skip it rather
						// than surface a false positive.
						continue
					}
					hits = append(hits, hit{r.CVEID, r.Severity, r.VersionEndExcluding, r.Title, r.URL, "nvd", r.Details, r.AffectedRanges})
				}
			}
		}

		for _, h := range hits {
			f := state.Finding{
				ID:             FindingID(repo.Repo, "componentCVE", h.cveID, comp.Package),
				Repo:           repo.Repo,
				Type:           "componentCVE",
				CVEID:          h.cveID,
				Ecosystem:      "hadron",
				Package:        comp.Package,
				CurrentVersion: comp.Version,
				FixedVersion:   h.fixed,
				Severity:       h.severity,
				Source:         h.source,
				Title:          h.title,
				URL:            h.url,
				FirstSeen:      Today(),
				LastSeen:       Today(),
				Details:        h.details,
				AffectedRanges: h.affected,
			}
			out[f.ID] = f
		}
	}

	res := make([]state.Finding, 0, len(out))
	for _, f := range out {
		res = append(res, f)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ID < res[j].ID })
	return res, nil
}
