package review

import "github.com/kairos-io/security/internal/ghclient"

// Assessor judges a bot-authored pull request from its diff, returning one of
// the verdicts good | bad | needs_human_verification along with reasoning.
type Assessor interface {
	Assess(diff []byte, pr ghclient.PullRequest) (verdict, reasoning string, err error)
}

// FakeAssessor is a test double returning canned values.
type FakeAssessor struct {
	Verdict, Reasoning string
	Err                error
}

func (f *FakeAssessor) Assess([]byte, ghclient.PullRequest) (string, string, error) {
	return f.Verdict, f.Reasoning, f.Err
}
