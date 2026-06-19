package remediate

type Classification struct {
	Intent  string // request-change | question | nack | approve | other
	Version string // explicit version requested, if any
	Reply   string // suggested reply text
}

type CommentClassifier interface {
	Classify(prTitle, author, body string) (Classification, error)
}

// FakeClassifier is a test double.
type FakeClassifier struct {
	Result Classification
	Err    error
}

func (f FakeClassifier) Classify(_, _, _ string) (Classification, error) {
	return f.Result, f.Err
}
