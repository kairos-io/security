package ghclient

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListDependabotAlertsTreats403AsNoAlerts(t *testing.T) {
	c := &CLI{run: func(args ...string) ([]byte, error) {
		return nil, errors.New("gh api: HTTP 403: Dependabot alerts are disabled")
	}}
	alerts, err := c.ListDependabotAlerts("kairos-io/kairos")
	require.NoError(t, err)
	assert.Nil(t, alerts)
}

func TestListDependabotAlertsPropagatesOtherErrors(t *testing.T) {
	c := &CLI{run: func(args ...string) ([]byte, error) {
		return nil, errors.New("gh api: HTTP 500: server error")
	}}
	_, err := c.ListDependabotAlerts("kairos-io/kairos")
	require.Error(t, err)
}

func TestListDependabotAlertsParsesJSON(t *testing.T) {
	c := &CLI{run: func(args ...string) ([]byte, error) {
		return []byte(`[{"number":1,"cveID":"CVE-2024-0001","package":"foo","severity":"high"}]`), nil
	}}
	alerts, err := c.ListDependabotAlerts("kairos-io/kairos")
	require.NoError(t, err)
	require.Len(t, alerts, 1)
	assert.Equal(t, "CVE-2024-0001", alerts[0].CVEID)
	assert.Equal(t, "high", alerts[0].Severity)
}

// labelRun builds a `run` fake around a repo whose labels are `have`. When
// canCreate is false, `label create` fails the way it does when the token
// cannot write labels in the tracking repo; when true, it succeeds and the
// label becomes available. `issue create` rejects an unknown label exactly as
// gh does. Every call is recorded.
func labelRun(calls *[][]string, have []string, canCreate bool) func(...string) ([]byte, error) {
	owned := append([]string(nil), have...)
	return func(args ...string) ([]byte, error) {
		*calls = append(*calls, args)
		switch args[0] {
		case "label":
			switch args[1] {
			case "list":
				return []byte(strings.Join(owned, "\n")), nil
			case "create":
				if !canCreate {
					return nil, errors.New("gh label: HTTP 403: Resource not accessible by integration")
				}
				owned = append(owned, args[2])
				return nil, nil
			}
		case "issue":
			switch args[1] {
			case "list":
				return nil, nil // no existing issue: take the create path
			case "create":
				for i, a := range args {
					if a != "--label" {
						continue
					}
					if !slices.Contains(owned, args[i+1]) {
						// What gh really does with an unknown label.
						return nil, errors.New("gh issue: could not add label: '" + args[i+1] + "' not found")
					}
				}
				return []byte("https://github.com/o/r/issues/7\n"), nil
			}
		}
		return nil, errors.New("unexpected call: gh " + strings.Join(args, " "))
	}
}

func labelArgsOf(calls [][]string, verb string) []string {
	var out []string
	for _, c := range calls {
		if c[0] == "label" && c[1] == verb {
			out = append(out, c[2])
		}
	}
	return out
}

// Where the token can create labels, every wanted label must reach the issue,
// and each create must happen before the issue is created.
func TestUpsertIssueCreatesMissingLabelsBeforeCreate(t *testing.T) {
	var calls [][]string
	c := &CLI{run: labelRun(&calls, nil, true)}

	n, err := c.UpsertIssue("o/r", "<!-- marker -->", "Title", "body",
		[]string{"security", "kairos-security-bot"})
	require.NoError(t, err)
	assert.Equal(t, 7, n)

	assert.Equal(t, []string{"security", "kairos-security-bot"}, labelArgsOf(calls, "create"))

	createIdx, lastLabelIdx := -1, -1
	for i, call := range calls {
		if call[0] == "label" && call[1] == "create" {
			lastLabelIdx = i
		}
		if call[0] == "issue" && call[1] == "create" {
			createIdx = i
			assert.Contains(t, call, "security")
			assert.Contains(t, call, "kairos-security-bot")
		}
	}
	require.NotEqual(t, -1, createIdx, "issue create must be invoked")
	assert.Less(t, lastLabelIdx, createIdx, "labels must be created before the issue")
}

// The tracking repo (kairos-io/kairos) has neither `security` nor
// `kairos-security-bot`, and the bot cannot create them there. `gh issue
// create` fails outright on an unknown label, so a label we cannot have must
// be dropped: an unlabelled tracking issue beats no tracking issue.
func TestUpsertIssueDropsLabelsTheRepoCannotHave(t *testing.T) {
	var calls [][]string
	c := &CLI{run: labelRun(&calls, []string{"area/security", "bug"}, false)}

	n, err := c.UpsertIssue("o/r", "<!-- marker -->", "Title", "body",
		[]string{"security", "kairos-security-bot"})
	require.NoError(t, err)
	assert.Equal(t, 7, n)

	for _, call := range calls {
		if call[0] == "issue" && call[1] == "create" {
			assert.NotContains(t, call, "security")
			assert.NotContains(t, call, "kairos-security-bot")
		}
	}
}

// A label the repo already owns must not be re-created. `gh label create
// --force` rewrites an existing label's colour and description, so running it
// over kairos-io/kairos' own `area/security` would repaint it.
func TestUpsertIssueDoesNotRecreateALabelTheRepoOwns(t *testing.T) {
	var calls [][]string
	c := &CLI{run: labelRun(&calls, []string{"area/security"}, false)}

	_, err := c.UpsertIssue("o/r", "<!-- marker -->", "Title", "body",
		[]string{"area/security", "kairos-security-bot"})
	require.NoError(t, err)

	assert.NotContains(t, labelArgsOf(calls, "create"), "area/security")
	assert.Contains(t, labelArgsOf(calls, "create"), "kairos-security-bot",
		"a missing label is still worth one best-effort create")
	for _, call := range calls {
		if call[0] == "label" && call[1] == "create" {
			assert.NotContains(t, call, "--force",
				"--force would repaint a label the repo owns")
		}
	}
}

// The existing-issue lookup filters by the last label. A label that did not
// survive must not be used as the filter, or the lookup returns nothing and
// the upsert creates a duplicate instead of editing.
func TestUpsertIssueLookupFiltersOnlyBySurvivingLabels(t *testing.T) {
	var calls [][]string
	c := &CLI{run: labelRun(&calls, []string{"area/security"}, false)}

	_, err := c.UpsertIssue("o/r", "<!-- marker -->", "Title", "body",
		[]string{"area/security", "kairos-security-bot"})
	require.NoError(t, err)

	for _, call := range calls {
		if call[0] == "issue" && call[1] == "list" {
			assert.NotContains(t, call, "kairos-security-bot")
		}
	}
}
