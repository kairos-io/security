package ghclient

import (
	"errors"
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

func TestUpsertIssueCreatesLabelsBeforeCreate(t *testing.T) {
	var calls [][]string
	c := &CLI{run: func(args ...string) ([]byte, error) {
		calls = append(calls, args)
		switch args[0] {
		case "issue":
			if args[1] == "list" {
				return nil, nil // no existing issue: take the create path
			}
			if args[1] == "create" {
				return []byte("https://github.com/o/r/issues/7\n"), nil
			}
		case "label":
			// Simulate "already exists": label create must be tolerated.
			return nil, errors.New("could not add label: already exists")
		}
		return nil, nil
	}}

	n, err := c.UpsertIssue("o/r", "<!-- marker -->", "Title", "body", []string{"security", "kairos-security-bot"})
	require.NoError(t, err)
	assert.Equal(t, 7, n)

	// Every label must be created, and all label-create calls must precede
	// the issue-create call.
	var labelCreated []string
	createIdx, lastLabelIdx := -1, -1
	for i, c := range calls {
		if c[0] == "label" && c[1] == "create" {
			labelCreated = append(labelCreated, c[2])
			lastLabelIdx = i
		}
		if c[0] == "issue" && c[1] == "create" {
			createIdx = i
		}
	}
	assert.Equal(t, []string{"security", "kairos-security-bot"}, labelCreated)
	require.NotEqual(t, -1, createIdx, "issue create must be invoked")
	assert.Less(t, lastLabelIdx, createIdx, "labels must be created before the issue")
}
