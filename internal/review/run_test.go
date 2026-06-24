package review

import (
	"strconv"
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
}

func (f *fakeGH) ListOpenPRs(repo string) ([]ghclient.PullRequest, error) { return f.prs[repo], nil }
func (f *fakeGH) PRDiff(repo string, pr int) ([]byte, error) {
	return f.diffs[repo], nil
}
func (f *fakeGH) CompareDiff(repo, base, head string) ([]byte, error) {
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
			{Number: 2, Title: "bump x", Author: "app/dependabot", IsBot: true, HeadSHA: "sha2", URL: "u2"},
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
	assert.Contains(t, a.GotContext, "Upstream github.com/foo/bar 1.2.3..1.2.4:")
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
			{Number: 2, Title: "bump x", IsBot: true, HeadSHA: "sha2", URL: "u2"},
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
