package review

import (
	"strconv"
	"strings"
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeGH struct {
	ghclient.GitHub
	prs      map[string][]ghclient.PullRequest
	diffs    map[string][]byte
	compares map[string][]byte // keyed "repo base..head"
	upserts  map[string]string // keyed "repo#pr" -> body (re-upsert overwrites: no spam)
	compErr  error             // if set, CompareDiff returns it
	approved []string          // "repo#pr"
	compared [][3]string       // records each (repo, base, head) CompareDiff was called with
}

func (f *fakeGH) ListOpenPRs(repo string) ([]ghclient.PullRequest, error) { return f.prs[repo], nil }
func (f *fakeGH) PRDiff(repo string, pr int) ([]byte, error) {
	return f.diffs[repo], nil
}
func (f *fakeGH) CompareDiff(repo, base, head string) ([]byte, error) {
	f.compared = append(f.compared, [3]string{repo, base, head})
	if f.compErr != nil {
		return nil, f.compErr
	}
	return f.compares[repo+" "+base+".."+head], nil
}
func (f *fakeGH) UpsertPRComment(repo string, pr int, marker, body string) error {
	if f.upserts == nil {
		f.upserts = map[string]string{}
	}
	f.upserts[repo+"#"+strconv.Itoa(pr)] = body
	return nil
}
func (f *fakeGH) ApprovePR(repo string, pr int, body string) error {
	f.approved = append(f.approved, repo+"#"+strconv.Itoa(pr))
	return nil
}

func TestRunAssessesBotPRsIdempotently(t *testing.T) {
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 2, Title: "bump x", Author: "app/dependabot", IsBot: true, HeadSHA: "sha2", URL: "u2", Body: "### Release Notes"},
			{Number: 3, Title: "human pr", Author: "alice", HeadSHA: "sha3"}, // not a bot -> skipped
		}},
		diffs: map[string][]byte{"o/r": []byte("go.mod bump")},
	}
	a := &FakeAssessor{Verdict: "good", Reasoning: "clean"}
	cfg := config.ReviewCfg{Enabled: true, AutoApprove: true, MaxPerRun: 20, Notify: []string{"@team"}}

	// First run: assesses the bot PR, upserts a comment (with cc), approves (good+autoApprove).
	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, out, 1)
	assert.Equal(t, "good", out[0].Verdict)
	assert.Equal(t, "sha2", out[0].HeadSHA)
	require.Len(t, gh.upserts, 1)
	assert.Contains(t, gh.upserts["o/r#2"], "good")
	assert.Contains(t, gh.upserts["o/r#2"], "@team")
	assert.Contains(t, gh.upserts["o/r#2"], reviewMarker)
	require.Len(t, gh.approved, 1)

	// Second run with the SAME head SHA in prev: carried forward, no new assessment.
	a2 := &panicAssessor{t} // would fail the test if called
	out2, _ := Run([]state.Repo{{Repo: "o/r"}}, gh, a2, cfg, out, "run2", false)
	require.Len(t, out2, 1)
	assert.Equal(t, "good", out2[0].Verdict)
	assert.Len(t, gh.upserts, 1)  // unchanged — idempotent
	assert.Len(t, gh.approved, 1) // unchanged — idempotent
}

// panicAssessor fails the test if Assess is ever invoked.
type panicAssessor struct{ t *testing.T }

func (p *panicAssessor) Assess(ghclient.PullRequest, string) (string, string, string, error) {
	p.t.Fatal("assessor must not be called for unchanged head SHA")
	return "", "", "", nil
}

func TestRunAssemblesContextAndUpsertsOnce(t *testing.T) {
	// A go.mod bump whose module resolves to a GitHub repo, with an upstream
	// compare diff available.
	diff := []byte("--- a/go.mod\n+++ b/go.mod\n" +
		"-	github.com/foo/bar v1.2.3\n" +
		"+	github.com/foo/bar v1.2.4\n")
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 2, Title: "bump bar", IsBot: true, HeadSHA: "sha2", URL: "u2", Body: "Bumps bar from 1.2.3 to 1.2.4 changelog"},
		}},
		diffs:    map[string][]byte{"o/r": diff},
		compares: map[string][]byte{"foo/bar v1.2.3..v1.2.4": []byte("upstream source change here")},
	}
	a := &FakeAssessor{Verdict: "good", Reasoning: "clean", ChangesSummary: "bumps bar, no behavior change"}
	cfg := config.ReviewCfg{Enabled: true, MaxPerRun: 20}

	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, out, 1)

	// Context carries the PR body and the labelled upstream diff.
	assert.Contains(t, a.GotContext, "Bumps bar from 1.2.3 to 1.2.4 changelog")
	assert.Contains(t, a.GotContext, "Upstream github.com/foo/bar 1.2.3→1.2.4")
	assert.Contains(t, a.GotContext, "upstream source change here")
	assert.Contains(t, a.GotContext, "PR diff:")

	// ChangesSummary is recorded and surfaced in the upserted comment.
	assert.Equal(t, "bumps bar, no behavior change", out[0].ChangesSummary)
	require.Len(t, gh.upserts, 1)
	assert.Contains(t, gh.upserts["o/r#2"], "**Dependency changes:** bumps bar, no behavior change")

	// Re-assessing the same PR (changed head) upserts in place — no spam.
	gh.prs["o/r"][0].HeadSHA = "sha3"
	_, _ = Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, out, "run2", false)
	assert.Len(t, gh.upserts, 1) // still one comment per PR (edited, not appended)
}

func TestRunCompareDiffErrorStillAssesses(t *testing.T) {
	diff := []byte("-	github.com/foo/bar v1.2.3\n+	github.com/foo/bar v1.2.4\n")
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 2, Title: "bump bar", IsBot: true, HeadSHA: "sha2", URL: "u2", Body: "changelog body"},
		}},
		diffs:   map[string][]byte{"o/r": diff},
		compErr: assert.AnError, // upstream compare fails -> degrade, but still assess
	}
	a := &FakeAssessor{Verdict: "needs_human_verification", Reasoning: "no upstream"}
	cfg := config.ReviewCfg{Enabled: true, MaxPerRun: 20}

	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, out, 1)
	assert.Equal(t, "needs_human_verification", out[0].Verdict)
	// Body still present; upstream diff label absent because compare failed.
	assert.Contains(t, a.GotContext, "changelog body")
	assert.NotContains(t, a.GotContext, "Upstream github.com/foo/bar")
	require.Len(t, gh.upserts, 1)
}

func TestRunDryRunWritesNothing(t *testing.T) {
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 2, Title: "bump x", IsBot: true, HeadSHA: "sha2", URL: "u2", Body: "### Release Notes"},
		}},
		diffs: map[string][]byte{"o/r": []byte("go.mod bump")},
	}
	a := &FakeAssessor{Verdict: "good", Reasoning: "clean"}
	cfg := config.ReviewCfg{Enabled: true, AutoApprove: true, MaxPerRun: 20}

	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", true)
	require.Empty(t, errs)
	require.Len(t, out, 1)
	assert.Equal(t, "good", out[0].Verdict)
	assert.Empty(t, gh.upserts) // dry-run: zero writes
	assert.Empty(t, gh.approved)
}

func TestRunAutoApproveOnlyOnGood(t *testing.T) {
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 2, Title: "bad bump", IsBot: true, HeadSHA: "sha2", URL: "u2"},
		}},
		diffs: map[string][]byte{"o/r": []byte("go.mod bump")},
	}
	a := &FakeAssessor{Verdict: "bad", Reasoning: "suspicious"}
	cfg := config.ReviewCfg{Enabled: true, AutoApprove: true, MaxPerRun: 20}

	_, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, gh.upserts, 1) // still comments
	assert.Empty(t, gh.approved)  // but does not approve a non-good verdict
}

func TestRunCapsNewAssessments(t *testing.T) {
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 1, Title: "a", IsBot: true, HeadSHA: "s1", URL: "u1"},
			{Number: 2, Title: "b", IsBot: true, HeadSHA: "s2", URL: "u2"},
			{Number: 3, Title: "c", IsBot: true, HeadSHA: "s3", URL: "u3"},
		}},
		diffs: map[string][]byte{"o/r": []byte("diff")},
	}
	a := &FakeAssessor{Verdict: "good", Reasoning: "ok"}
	cfg := config.ReviewCfg{Enabled: true, MaxPerRun: 2}

	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	assert.Len(t, gh.upserts, 2) // capped at MaxPerRun new assessments
	assert.Len(t, out, 2)
}

func TestRunPseudoVersionComparesBySHA(t *testing.T) {
	// A pseudo-version bump: compare must use the embedded commit SHAs as refs.
	diff := []byte("--- a/go.mod\n+++ b/go.mod\n" +
		"-	github.com/foo/bar v0.0.0-20241017190036-fab4fdf2f2f3\n" +
		"+	github.com/foo/bar v0.0.0-20241101120000-d29549a44f29\n")
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 2, Title: "bump bar", IsBot: true, HeadSHA: "sha2", URL: "u2"},
		}},
		diffs:    map[string][]byte{"o/r": diff},
		compares: map[string][]byte{"foo/bar fab4fdf2f2f3..d29549a44f29": []byte("upstream pseudo source change")},
	}
	a := &FakeAssessor{Verdict: "good", Reasoning: "clean"}
	cfg := config.ReviewCfg{Enabled: true, MaxPerRun: 20}

	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, out, 1)

	// CompareDiff received the SHA refs, not the raw pseudo-version strings.
	require.Len(t, gh.compared, 1)
	assert.Equal(t, [3]string{"foo/bar", "fab4fdf2f2f3", "d29549a44f29"}, gh.compared[0])
	assert.Contains(t, a.GotContext, "upstream pseudo source change")

	// Trace records the successful compare with a byte count, plus the context line.
	require.NotEmpty(t, out[0].Trace)
	var sawCompare, sawContext bool
	for _, line := range out[0].Trace {
		if strings.Contains(line, "compare fab4fdf2f2f3...d29549a44f29 ✓") && strings.Contains(line, "bytes") {
			sawCompare = true
		}
		if strings.HasPrefix(line, "context: ") {
			sawContext = true
		}
	}
	assert.True(t, sawCompare, "trace should record the successful SHA compare: %v", out[0].Trace)
	assert.True(t, sawContext, "trace should record the context size: %v", out[0].Trace)

	// The upserted comment embeds the collapsible trace block before the marker.
	body := gh.upserts["o/r#2"]
	assert.Contains(t, body, "<details><summary>review trace</summary>")
	assert.Less(t, strings.Index(body, "<details>"), strings.Index(body, reviewMarker))
}

func TestRunTraceUnresolvableModule(t *testing.T) {
	// A vanity/unresolvable module path with no compare links in the body ->
	// compareTargets yields nothing, so no upstream comparison is attempted.
	diff := []byte("-	example.com/private/thing v1.0.0\n+	example.com/private/thing v1.1.0\n")
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 2, Title: "bump thing", IsBot: true, HeadSHA: "sha2", URL: "u2"},
		}},
		diffs: map[string][]byte{"o/r": diff},
	}
	a := &FakeAssessor{Verdict: "needs_human_verification", Reasoning: "no upstream"}
	cfg := config.ReviewCfg{Enabled: true, MaxPerRun: 20}

	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, out, 1)

	assert.Empty(t, gh.compared) // unresolvable -> never calls CompareDiff
	require.NotEmpty(t, out[0].Trace)
	var sawNoTargets bool
	for _, line := range out[0].Trace {
		if strings.Contains(line, "no upstream comparisons available") {
			sawNoTargets = true
		}
	}
	assert.True(t, sawNoTargets, "trace should record the no-targets case: %v", out[0].Trace)
}

func TestRunTraceNoBumps(t *testing.T) {
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 2, Title: "docs only", IsBot: true, HeadSHA: "sha2", URL: "u2"},
		}},
		diffs: map[string][]byte{"o/r": []byte("--- a/README.md\n+++ b/README.md\n-old\n+new\n")},
	}
	a := &FakeAssessor{Verdict: "good", Reasoning: "docs"}
	cfg := config.ReviewCfg{Enabled: true, MaxPerRun: 20}

	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, out, 1)

	require.NotEmpty(t, out[0].Trace)
	var sawNoTargets bool
	for _, line := range out[0].Trace {
		if strings.Contains(line, "no upstream comparisons available") {
			sawNoTargets = true
		}
	}
	assert.True(t, sawNoTargets, "trace should record the no-targets case: %v", out[0].Trace)
}

func TestRunFetchesUpstreamFromBodyCompareLink(t *testing.T) {
	// A bot PR with NO go.mod bump but a renovate body carrying a GitHub compare
	// link (e.g. an npm bump). The upstream diff must be fetched using the body's
	// verbatim base/head refs.
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 7, Title: "chore(deps): update lucide", IsBot: true, HeadSHA: "sha7", URL: "u7",
				Body: "Updates lucide-react. See https://github.com/lucide-icons/lucide/compare/0.576.0...0.577.0 for details."},
		}},
		diffs:    map[string][]byte{"o/r": []byte("--- a/package.json\n+++ b/package.json\n-lucide\n+lucide\n")},
		compares: map[string][]byte{"lucide-icons/lucide 0.576.0..0.577.0": []byte("upstream lucide source change")},
	}
	a := &FakeAssessor{Verdict: "good", Reasoning: "clean"}
	cfg := config.ReviewCfg{Enabled: true, MaxPerRun: 20}

	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, out, 1)

	// CompareDiff was called with the verbatim refs from the body link.
	require.Len(t, gh.compared, 1)
	assert.Equal(t, [3]string{"lucide-icons/lucide", "0.576.0", "0.577.0"}, gh.compared[0])

	// The fetched upstream diff is labelled and embedded in the assessor context.
	assert.Contains(t, a.GotContext, "Upstream lucide-icons/lucide")
	assert.Contains(t, a.GotContext, "upstream lucide source change")

	// Trace records the successful compare with a ✓ and byte count.
	require.NotEmpty(t, out[0].Trace)
	var sawCompare bool
	for _, line := range out[0].Trace {
		if strings.Contains(line, "lucide-icons/lucide") &&
			strings.Contains(line, "compare 0.576.0...0.577.0 ✓") && strings.Contains(line, "bytes") {
			sawCompare = true
		}
	}
	assert.True(t, sawCompare, "trace should record the successful body-link compare: %v", out[0].Trace)
}

func TestRunForcesNeedsHumanWhenNoContext(t *testing.T) {
	// No changelog (empty body) and the upstream compare fails -> 0 bytes of
	// evidence. Even if the model says "good", force needs_human_verification.
	diff := []byte("-	github.com/foo/bar v1.2.3\n+	github.com/foo/bar v1.2.4\n")
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 2, Title: "bump bar", IsBot: true, HeadSHA: "sha2", URL: "u2", Body: ""},
		}},
		diffs:   map[string][]byte{"o/r": diff},
		compErr: assert.AnError, // no upstream diff fetched
	}
	a := &FakeAssessor{Verdict: "good", Reasoning: "looks fine"} // model would pass it
	cfg := config.ReviewCfg{Enabled: true, MaxPerRun: 20}

	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, out, 1)
	assert.Equal(t, "needs_human_verification", out[0].Verdict) // overridden
	found := false
	for _, tl := range out[0].Trace {
		if strings.Contains(tl, "insufficient context") {
			found = true
		}
	}
	assert.True(t, found, "trace should record the forced override")
}

func TestRunKeepsVerdictWhenChangelogPresent(t *testing.T) {
	// Empty upstream diff but a non-empty changelog body -> the model has
	// something to go on, so its verdict stands.
	gh := &fakeGH{
		prs: map[string][]ghclient.PullRequest{"o/r": {
			{Number: 3, Title: "bump x", IsBot: true, HeadSHA: "sha3", URL: "u3", Body: "### Release Notes\nfixed a bug"},
		}},
		diffs: map[string][]byte{"o/r": []byte("docs only\n")},
	}
	a := &FakeAssessor{Verdict: "good", Reasoning: "clean"}
	cfg := config.ReviewCfg{Enabled: true, MaxPerRun: 20}
	out, _ := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Len(t, out, 1)
	assert.Equal(t, "good", out[0].Verdict) // not forced (changelog present)
}
