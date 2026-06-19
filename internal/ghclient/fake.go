package ghclient

type FakeIssue struct {
	Number int
	Title  string
	Body   string
	Labels []string
}

// Fake is an in-memory GitHub double for tests.
type Fake struct {
	OrgRepos map[string][]string
	Files    map[string][]byte // key: repo|path|ref
	PRs      map[string][]PullRequest
	Alerts   map[string][]Alert
	Issues   map[string]*FakeIssue // key: repo
	nextNum  int
}

func NewFake() *Fake {
	return &Fake{
		OrgRepos: map[string][]string{},
		Files:    map[string][]byte{},
		PRs:      map[string][]PullRequest{},
		Alerts:   map[string][]Alert{},
		Issues:   map[string]*FakeIssue{},
	}
}

func (f *Fake) ListOrgRepos(org string) ([]string, error) { return f.OrgRepos[org], nil }
func (f *Fake) GetFile(repo, path, ref string) ([]byte, error) {
	return f.Files[repo+"|"+path+"|"+ref], nil
}
func (f *Fake) ListOpenPRs(repo string) ([]PullRequest, error)    { return f.PRs[repo], nil }
func (f *Fake) ListDependabotAlerts(repo string) ([]Alert, error) { return f.Alerts[repo], nil }

func (f *Fake) UpsertIssue(repo, marker, title, body string, labels []string) (int, error) {
	if iss, ok := f.Issues[repo]; ok {
		iss.Body, iss.Title, iss.Labels = body, title, labels
		return iss.Number, nil
	}
	f.nextNum++
	f.Issues[repo] = &FakeIssue{Number: f.nextNum, Title: title, Body: body, Labels: labels}
	return f.nextNum, nil
}
