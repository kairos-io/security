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
	comments []string // "repo#pr: body"
	approved []string // "repo#pr"
}

func (f *fakeGH) ListOpenPRs(repo string) ([]ghclient.PullRequest, error) { return f.prs[repo], nil }
func (f *fakeGH) PRDiff(repo string, pr int) ([]byte, error) {
	return f.diffs[repo], nil
}
func (f *fakeGH) PostPRComment(repo string, pr int, body string) error {
	f.comments = append(f.comments, repo+"#"+strconv.Itoa(pr)+": "+body)
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

	// First run: assesses the bot PR, comments (with cc), approves (good+autoApprove).
	out, errs := Run([]state.Repo{{Repo: "o/r"}}, gh, a, cfg, nil, "run1", false)
	require.Empty(t, errs)
	require.Len(t, out, 1)
	assert.Equal(t, "good", out[0].Verdict)
	assert.Equal(t, "sha2", out[0].HeadSHA)
	require.Len(t, gh.comments, 1)
	assert.Contains(t, gh.comments[0], "good")
	assert.Contains(t, gh.comments[0], "@team")
	assert.Contains(t, gh.comments[0], reviewMarker)
	require.Len(t, gh.approved, 1)

	// Second run with the SAME head SHA in prev: carried forward, no new comment/approve.
	a2 := &panicAssessor{t} // would fail the test if called
	out2, _ := Run([]state.Repo{{Repo: "o/r"}}, gh, a2, cfg, out, "run2", false)
	require.Len(t, out2, 1)
	assert.Equal(t, "good", out2[0].Verdict)
	assert.Len(t, gh.comments, 1) // unchanged — idempotent
	assert.Len(t, gh.approved, 1) // unchanged — idempotent
}

// panicAssessor fails the test if Assess is ever invoked.
type panicAssessor struct{ t *testing.T }

func (p *panicAssessor) Assess([]byte, ghclient.PullRequest) (string, string, error) {
	p.t.Fatal("assessor must not be called for unchanged head SHA")
	return "", "", nil
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
	assert.Empty(t, gh.comments) // dry-run: zero writes
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
	require.Len(t, gh.comments, 1) // still comments
	assert.Empty(t, gh.approved)   // but does not approve a non-good verdict
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
	assert.Len(t, gh.comments, 2) // capped at MaxPerRun new assessments
	assert.Len(t, out, 2)
}
