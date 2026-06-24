package review

import "github.com/kairos-io/security/internal/ghclient"

// Assessor judges a bot-authored pull request from an assembled review context
// (changelog + upstream source diffs + the PR diff), returning one of the
// verdicts good | bad | needs_human_verification along with reasoning and a
// short summary of the dependency change.
type Assessor interface {
	Assess(pr ghclient.PullRequest, context string) (verdict, reasoning, changesSummary string, err error)
}

// FakeAssessor is a test double returning canned values. It records the
// assembled context it was given so tests can assert on it.
type FakeAssessor struct {
	Verdict, Reasoning, ChangesSummary string
	Err                                error
	GotContext                         string // records the assembled context for assertions
}

func (f *FakeAssessor) Assess(_ ghclient.PullRequest, context string) (string, string, string, error) {
	f.GotContext = context
	return f.Verdict, f.Reasoning, f.ChangesSummary, f.Err
}
