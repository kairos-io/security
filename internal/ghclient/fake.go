package ghclient

import "fmt"

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

	PRComments map[string][]ReviewComment // key: "<repo>#<pr>"
	Posted     []string
	Closed     []string

	Statuses map[string]PRStatus // key: "<repo>#<pr>"
	Merged   []string
}

func NewFake() *Fake {
	return &Fake{
		OrgRepos:   map[string][]string{},
		Files:      map[string][]byte{},
		PRs:        map[string][]PullRequest{},
		Alerts:     map[string][]Alert{},
		Issues:     map[string]*FakeIssue{},
		PRComments: map[string][]ReviewComment{},
		Statuses:   map[string]PRStatus{},
	}
}

func prKey(repo string, pr int) string { return fmt.Sprintf("%s#%d", repo, pr) }

func (f *Fake) ListPRComments(repo string, pr int) ([]ReviewComment, error) {
	return f.PRComments[prKey(repo, pr)], nil
}
func (f *Fake) PostPRComment(repo string, pr int, body string) error {
	f.Posted = append(f.Posted, prKey(repo, pr)+": "+body)
	return nil
}
func (f *Fake) ClosePR(repo string, pr int, comment string) error {
	f.Closed = append(f.Closed, prKey(repo, pr))
	return nil
}
func (f *Fake) PRStatusOf(repo string, pr int) (PRStatus, error) { return f.Statuses[prKey(repo, pr)], nil }
func (f *Fake) MergePR(repo string, pr int, auto bool) error {
	k := prKey(repo, pr)
	if auto {
		k += " (auto)"
	}
	f.Merged = append(f.Merged, k)
	return nil
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
