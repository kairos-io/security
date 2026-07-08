package correlate

import (
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeApplier lets tests assert what correlate hands to the classifier and
// attach an AIApplicability verdict without hitting the network.
type fakeApplier struct {
	got  []state.Finding
	mark map[string]string // finding.ID -> reason (attaches AIApplicability=not-applicable/high)
}

func (f *fakeApplier) Apply(in []state.Finding) []state.Finding {
	f.got = append([]state.Finding(nil), in...)
	out := make([]state.Finding, len(in))
	copy(out, in)
	for i := range out {
		if reason, ok := f.mark[out[i].ID]; ok {
			out[i].AIApplicability = &state.AIApplicability{
				Applicable: false,
				Confidence: "high",
				Reasoning:  reason,
			}
		}
	}
	return out
}

func TestCorrelateDedupesAndBuildsWaterfall(t *testing.T) {
	in := state.Findings{Findings: []state.Finding{
		// same CVE in immucore seen by two sources → dedupe to 1, severity high wins
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "unknown", FixedVersion: "0.33.0"},
		{ID: "b", Repo: "kairos-io/immucore", Type: "ghAlert", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high"},
		// same CVE/package in a second repo → waterfall group of 2 repos
		{ID: "c", Repo: "kairos-io/kairos-agent", Type: "sourceCVE", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high", FixedVersion: "0.33.0"},
	}}

	out := Run(in, config.CVEPolicy{}, nil)

	// dedupe: immucore CVE-2025-1 collapses to one finding, severity "high"
	count := 0
	for _, f := range out.Findings {
		if f.Repo == "kairos-io/immucore" && f.CVEID == "CVE-2025-1" {
			count++
			assert.Equal(t, "high", f.Severity)
			assert.Equal(t, "0.33.0", f.FixedVersion)
		}
	}
	assert.Equal(t, 1, count)

	require.Len(t, out.Waterfall, 1)
	g := out.Waterfall[0]
	assert.ElementsMatch(t, []string{"kairos-io/immucore", "kairos-io/kairos-agent"}, g.AffectedRepos)
	assert.Equal(t, "golang.org/x/net", g.SuggestedBump.Package)
	assert.Equal(t, "0.33.0", g.SuggestedBump.To)
	assert.Equal(t, "high", g.Severity)
}

func TestRun_WaterfallSkipsInformational(t *testing.T) {
	// The first finding is already-fixed (current 2.0.0 >= fixed 1.0.0), so
	// classify (now run inside Run) marks it informational and the waterfall
	// must not count it — leaving fewer than 2 repos, hence no front.
	in := state.Findings{Findings: []state.Finding{
		{ID: "1", Repo: "o/a", Ecosystem: "go", CVEID: "CVE-1", Package: "p", CurrentVersion: "2.0.0", FixedVersion: "1.0.0"},
		{ID: "2", Repo: "o/b", Ecosystem: "go", CVEID: "CVE-1", Package: "p"},
	}}
	out := Run(in, config.CVEPolicy{}, nil)
	if len(out.Waterfall) != 0 {
		t.Fatalf("informational finding must not count toward waterfall: %+v", out.Waterfall)
	}
}

// TestRun_ClassifyAfterDedupeIsOrderIndependent asserts that classification runs
// once on the merged finding, so a merge no longer lets an informational input
// hide an actionable one (and vice versa) depending on insertion order. Two
// findings share repo|cveID|package; the merged FixedVersion is "1.0.0" (first
// non-empty wins during dedupe here it comes from whichever has it). We assert
// the surviving finding's Class is identical regardless of input order.
func TestRun_ClassifyAfterDedupeIsOrderIndependent(t *testing.T) {
	// Merged current 1.5.0 vs fixed 2.0.0 -> not yet fixed -> actionable.
	a := state.Finding{ID: "x", Repo: "o/a", Ecosystem: "go", CVEID: "CVE-9", Package: "p", CurrentVersion: "1.5.0", FixedVersion: "2.0.0"}
	b := state.Finding{ID: "x", Repo: "o/a", Ecosystem: "go", CVEID: "CVE-9", Package: "p", CurrentVersion: "1.5.0"}

	classOf := func(in state.Findings) string {
		out := Run(in, config.CVEPolicy{}, nil)
		require.Len(t, out.Findings, 1)
		return out.Findings[0].Class
	}
	fwd := classOf(state.Findings{Findings: []state.Finding{a, b}})
	rev := classOf(state.Findings{Findings: []state.Finding{b, a}})
	assert.Equal(t, fwd, rev, "classification must be order-independent after dedupe")
	assert.Equal(t, "", fwd, "merged finding (current<fixed) stays actionable")
}

// TestRun_ApplierRunsAfterClassify verifies (a) the applier is invoked on the
// deduped, classified findings and (b) findings the applier flags with a
// not-applicable verdict stay COUNTED and actionable — the verdict is advisory
// metadata surfaced by renderers, not a hide-from-dashboard signal.
func TestRun_ApplierRunsAfterClassify(t *testing.T) {
	in := state.Findings{Findings: []state.Finding{
		{ID: "1", Repo: "o/a", Ecosystem: "go", CVEID: "CVE-1", Package: "p", CurrentVersion: "1.0.0", FixedVersion: "2.0.0"},
		{ID: "2", Repo: "o/b", Ecosystem: "go", CVEID: "CVE-1", Package: "p", CurrentVersion: "1.0.0", FixedVersion: "2.0.0"},
	}}
	fa := &fakeApplier{mark: map[string]string{"2": "does not affect linux"}}
	out := Run(in, config.CVEPolicy{}, fa)

	// applier saw the deduped, classified findings
	require.Len(t, fa.got, 2)
	for _, f := range fa.got {
		assert.Equal(t, "", f.Class, "applier must receive still-actionable findings, not pre-classified ones")
	}

	// finding 2 gets AIApplicability metadata but stays actionable
	byID := map[string]state.Finding{}
	for _, f := range out.Findings {
		byID[f.ID] = f
	}
	assert.Equal(t, "", byID["1"].Class)
	assert.Nil(t, byID["1"].AIApplicability)
	assert.Equal(t, "", byID["2"].Class, "AI-suspected findings must stay actionable")
	require.NotNil(t, byID["2"].AIApplicability)
	assert.False(t, byID["2"].AIApplicability.Applicable)
	assert.Contains(t, byID["2"].AIApplicability.Reasoning, "does not affect linux")

	// waterfall keeps counting both repos — AI suspicion does not silently
	// remove a finding from the counts.
	require.Len(t, out.Waterfall, 1, "AI-suspected findings must still count toward waterfall")
	assert.ElementsMatch(t, []string{"o/a", "o/b"}, out.Waterfall[0].AffectedRepos)
}

// TestRun_ApplierDoesNotOverrideAlreadyInformational makes sure existing
// classifications (accepted-component, already-fixed) are not touched by the
// applier.
func TestRun_ApplierDoesNotResurrect(t *testing.T) {
	in := state.Findings{Findings: []state.Finding{
		{ID: "1", Repo: "o/a", Ecosystem: "go", CVEID: "CVE-1", Package: "p", CurrentVersion: "2.0.0", FixedVersion: "1.0.0"},
	}}
	fa := &fakeApplier{}
	out := Run(in, config.CVEPolicy{}, fa)
	require.Len(t, out.Findings, 1)
	assert.Equal(t, "informational", out.Findings[0].Class)
	assert.Equal(t, "already-fixed", out.Findings[0].ClassReason)
}
