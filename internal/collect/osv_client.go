package collect

import (
	"encoding/json"
	"strconv"
	"strings"

	ver "github.com/kairos-io/security/internal/version"
)

// OSVQueryFunc performs a single OSV.dev query and returns the raw JSON
// response body. Real implementations POST to https://api.osv.dev/v1/query;
// tests inject a fixture-returning fake.
type OSVQueryFunc func(ecosystem, pkg, version string) ([]byte, error)

type osvVuln struct {
	ID               string   `json:"id"`
	Aliases          []string `json:"aliases"`
	Upstream         []string `json:"upstream"`
	Summary          string   `json:"summary"`
	Details          string   `json:"details"`
	DatabaseSpecific struct {
		Severity string `json:"severity"`
	} `json:"database_specific"`
	Severity []struct {
		Type  string `json:"type"`
		Score string `json:"score"`
	} `json:"severity"`
	Affected []osvAffected `json:"affected"`
}

// osvAffected is the affected-package entry from an OSV record. Kept as a
// named type so it can be re-marshalled into the Finding's AffectedRanges
// field verbatim (feeds the applicability classifier).
type osvAffected struct {
	Package struct {
		Ecosystem string `json:"ecosystem,omitempty"`
		Name      string `json:"name,omitempty"`
	} `json:"package,omitempty"`
	Ranges []struct {
		Type   string          `json:"type,omitempty"`
		Events []osvRangeEvent `json:"events"`
	} `json:"ranges"`
	Versions []string `json:"versions,omitempty"`
}

// osvRangeEvent is one entry in a range's event timeline. The schema allows
// only one of the fields per object, and the array as a whole describes a
// sequence of vulnerable windows: see rangeIntervals.
type osvRangeEvent struct {
	Introduced string `json:"introduced,omitempty"`
	Fixed      string `json:"fixed,omitempty"`
}

type osvQueryResponse struct {
	Vulns []osvVuln `json:"vulns"`
}

// OSVResult is one CVE hit from an OSV.dev query, normalized to ksec's finding shape.
type OSVResult struct {
	CVEID        string
	Severity     string // critical|high|medium|low|unknown, via osvSeverity
	FixedVersion string
	Title        string
	URL          string
	// Details is the OSV `details` markdown (full CVE description).
	Details string
	// AffectedRanges is the OSV `affected[]` sub-object re-marshalled to a
	// compact JSON string so downstream (the applicability classifier) can feed
	// it to a model without a second fetch. Empty when marshalling fails or the
	// record has no affected entries.
	AffectedRanges string
}

// QueryOSV queries OSV.dev for (ecosystem, pkg, version) and normalizes every
// hit. For the Alpine ecosystem, OSV's "fixed" version carries Alpine's own
// package-revision suffix (e.g. "3.6.4-r0"); stripAlpineRevisionSuffix
// removes it so the result approximates hadron's upstream version numbering.
func QueryOSV(query OSVQueryFunc, ecosystem, pkg, version string) ([]OSVResult, error) {
	raw, err := query(ecosystem, pkg, version)
	if err != nil {
		return nil, err
	}
	var resp osvQueryResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	out := make([]OSVResult, 0, len(resp.Vulns))
	for _, v := range resp.Vulns {
		cve := v.ID
		for _, alias := range v.Aliases {
			if strings.HasPrefix(alias, "CVE-") {
				cve = alias
			}
		}
		// Alpine OSV-converted advisories carry the real CVE id in "upstream"
		// (and never populate "aliases"); honor it with the same last-match-wins
		// semantics as the alias loop above.
		for _, up := range v.Upstream {
			if strings.HasPrefix(up, "CVE-") {
				cve = up
			}
		}
		fixed, applicable := osvApplicableFix(v, ver.Compare, version)
		if !applicable {
			continue // version is below every introduced boundary — not vulnerable
		}
		affected := ""
		if len(v.Affected) > 0 {
			if b, mErr := json.Marshal(v.Affected); mErr == nil {
				affected = string(b)
			}
		}
		out = append(out, OSVResult{
			CVEID:          cve,
			Severity:       osvSeverity(v),
			FixedVersion:   fixed,
			Title:          v.Summary,
			URL:            "https://osv.dev/vulnerability/" + v.ID,
			Details:        v.Details,
			AffectedRanges: affected,
		})
	}
	return out, nil
}

// osvApplicableFix picks the fix version for the range that applies to the
// queried version v, and reports whether v is at/after some range's
// "introduced" boundary (i.e. potentially vulnerable). Alpine revision suffixes
// are stripped for comparison. When v sits inside a range (introduced <= v <
// fixed) that range's fix is returned; when v is past a range's fix, the
// largest such fix is returned (already-fixed, classified later). When v is
// below every introduced boundary, applicable=false and the vuln is dropped.
func osvApplicableFix(v osvVuln, cmp func(string, string) int, queried string) (string, bool) {
	q := stripAlpineRevisionSuffix(queried)
	parseable := looksNumeric(q)
	bestFix := ""
	applicable := false
	for _, a := range v.Affected {
		for _, rg := range a.Ranges {
			for _, iv := range rangeIntervals(rg.Events) {
				// Only trust the introduced-boundary drop when q is orderable. A
				// non-numeric version can't be compared meaningfully, so we fail
				// OPEN (keep the vuln visible) instead of dropping it below "0".
				if parseable && cmp(q, stripAlpineRevisionSuffix(iv.introduced)) < 0 {
					continue // below this interval's introduced boundary
				}
				applicable = true
				// Alpine OSV records use "0" as a placeholder meaning "no known
				// fix yet" — carrying that forward as FixedVersion tricks the
				// deterministic classifier (current >= "0" == true) into flagging
				// every such finding as already-fixed and quietly hiding it in
				// the informational section. Treat "0" the same as an empty
				// fixed: the vuln is applicable but we don't have a target
				// version to bump to.
				if iv.fixed == "" || iv.fixed == "0" {
					continue
				}
				// Inside the interval (q < fixed): this is the fix we want.
				if cmp(q, iv.fixed) < 0 {
					return iv.fixed, true
				}
				// Past the fix: remember the largest fix seen (already-fixed).
				if bestFix == "" || cmp(iv.fixed, bestFix) > 0 {
					bestFix = iv.fixed
				}
			}
		}
	}
	return bestFix, applicable
}

// osvInterval is one half-open vulnerable window [introduced, fixed) taken
// from a range's event timeline. An empty fixed means the window is still
// open (no known fix).
type osvInterval struct {
	introduced string
	fixed      string
}

// rangeIntervals turns one range's events into the windows they describe.
//
// The events array is a *timeline*, not a single pair: the OSV schema says
// each object "describes a single version that either introduces or fixes a
// vulnerability", and a range that was patched on several branches carries
// them all. Go's database encodes every backported fix this way, e.g.
// GO-2026-4340 is
//
//	[{introduced:0} {fixed:1.24.12} {introduced:1.25.0} {fixed:1.25.6}]
//
// meaning vulnerable in [0, 1.24.12) AND in [1.25.0, 1.25.6). Reading only
// the last introduced and the last fixed collapses that to [1.25.0, 1.25.6)
// and reports 1.24.5 as unaffected, which is a false negative on a real CVE.
//
// Each introduced opens a window and the next fixed closes it; an introduced
// that arrives while a window is open closes the previous one as unfixed. A
// range whose first event is a fix (the schema requires an introduced, but be
// forgiving) is treated as starting at "0", the value the schema reserves for
// "sorts before any other version".
//
// last_affected and limit events are not modelled here: the decoder does not
// carry them, so a window bounded only by one of those stays open and the
// vuln stays visible, which is the safe direction for a scanner.
func rangeIntervals(events []osvRangeEvent) []osvInterval {
	var out []osvInterval
	introduced := ""
	open := false
	for _, ev := range events {
		switch {
		case ev.Introduced != "":
			if open {
				out = append(out, osvInterval{introduced: introduced})
			}
			introduced = ev.Introduced
			open = true
		case ev.Fixed != "":
			if !open {
				introduced = "0"
			}
			out = append(out, osvInterval{introduced: introduced, fixed: stripAlpineRevisionSuffix(ev.Fixed)})
			open = false
		}
	}
	if open {
		out = append(out, osvInterval{introduced: introduced})
	}
	if len(out) == 0 {
		// A range with no usable events still describes the package as
		// affected; keep the old all-versions default rather than dropping it.
		out = append(out, osvInterval{introduced: "0"})
	}
	return out
}

// looksNumeric reports whether v begins with a digit (so version.Compare can
// order it meaningfully). Non-numeric versions fail OPEN — we keep the vuln
// visible rather than silently dropping it.
func looksNumeric(v string) bool {
	return v != "" && v[0] >= '0' && v[0] <= '9'
}

// osvSeverity derives a Finding's severity from an OSV vuln record, in
// precedence order:
//  1. An explicit database_specific.severity label (CRITICAL/HIGH/MODERATE/
//     MEDIUM/LOW) is trusted directly — it's a human/tooling-assigned rating.
//  2. Otherwise the first parseable CVSS_V3 entry in the top-level severity
//     array is computed into a base score and mapped to a band.
//  3. Otherwise "unknown" — being honest about missing data beats guessing.
func osvSeverity(v osvVuln) string {
	switch strings.ToUpper(strings.TrimSpace(v.DatabaseSpecific.Severity)) {
	case "CRITICAL":
		return "critical"
	case "HIGH":
		return "high"
	case "MODERATE", "MEDIUM":
		return "medium"
	case "LOW":
		return "low"
	}
	for _, s := range v.Severity {
		if s.Type != "CVSS_V3" {
			continue
		}
		score, err := cvssV31BaseScore(s.Score)
		if err != nil {
			continue // try any remaining entries
		}
		return cvssSeverityLabel(score)
	}
	return "unknown"
}

// stripAlpineRevisionSuffix strips a trailing "-rN" Alpine package-revision
// suffix (e.g. "3.6.4-r0" -> "3.6.4"), leaving other version strings
// untouched.
func stripAlpineRevisionSuffix(v string) string {
	i := strings.LastIndex(v, "-r")
	if i <= 0 {
		return v
	}
	if _, err := strconv.Atoi(v[i+2:]); err != nil {
		return v
	}
	return v[:i]
}
