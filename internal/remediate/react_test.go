package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecideReaction(t *testing.T) {
	cases := []struct {
		name string
		in   Classification
		want ReactionKind
	}{
		{"nack closes", Classification{Intent: "nack", Reply: "ok"}, ReactClose},
		{"request-change with version adjusts", Classification{Intent: "request-change", Version: "0.36.0", Reply: "ok"}, ReactAdjust},
		{"request-change without version replies", Classification{Intent: "request-change", Reply: "could you clarify?"}, ReactReply},
		{"question replies", Classification{Intent: "question", Reply: "it's automated"}, ReactReply},
		{"approve does nothing", Classification{Intent: "approve"}, ReactNone},
		{"other does nothing", Classification{Intent: "other"}, ReactNone},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, DecideReaction(tc.in).Kind)
		})
	}
	adj := DecideReaction(Classification{Intent: "request-change", Version: "0.36.0", Reply: "r"})
	assert.Equal(t, "0.36.0", adj.ToVersion)
	assert.Equal(t, "r", adj.ReplyBody)
}

type fakeAdjuster struct{ called int }

func (f *fakeAdjuster) Adjust(e state.LedgerEntry, to, run string) (state.LedgerEntry, error) {
	f.called++
	e.Bump.To = to
	return e, nil
}

func TestReactToCommentsIdempotentAndActs(t *testing.T) {
	gh := ghclient.NewFake()
	gh.PRComments["kairos-io/immucore#412"] = []ghclient.ReviewComment{
		{ID: "c1", Author: "maintainer", Body: "pin to 0.36.0"},
		{ID: "self", Author: "kairos-security-bot", Body: "on it"}, // own reply: ignored
	}
	entry := &state.LedgerEntry{Key: "kairos-io/immucore|golang.org/x/net", Repo: "kairos-io/immucore",
		PRNumber: 412, State: "open", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}}
	adj := &fakeAdjuster{}
	cls := FakeClassifier{Result: Classification{Intent: "request-change", Version: "0.36.0", Reply: "Bumping."}}

	require.NoError(t, ReactToComments(entry, gh, cls, adj, "bump", "run1", false))
	assert.Equal(t, 1, adj.called, "adjusted to requested version")
	assert.Equal(t, "0.36.0", entry.Bump.To)
	assert.Contains(t, gh.Posted, "kairos-io/immucore#412: Bumping.")
	assert.Contains(t, entry.SeenComments, "c1")
	assert.NotContains(t, entry.SeenComments, "self", "never react to its own comment")

	// Second run: c1 already seen -> no further action.
	adj.called = 0
	gh.Posted = nil
	require.NoError(t, ReactToComments(entry, gh, cls, adj, "bump", "run2", false))
	assert.Equal(t, 0, adj.called)
	assert.Empty(t, gh.Posted)
}

func TestReactToCommentsClassifierErrorSkips(t *testing.T) {
	gh := ghclient.NewFake()
	gh.PRComments["r#1"] = []ghclient.ReviewComment{{ID: "c1", Author: "m", Body: "?"}}
	entry := &state.LedgerEntry{Key: "r|p", Repo: "r", PRNumber: 1, State: "open"}
	require.NoError(t, ReactToComments(entry, gh, FakeClassifier{Err: assertErr()}, &fakeAdjuster{}, "t", "run1", false))
	assert.Empty(t, gh.Posted, "no action taken on classifier error")
	assert.NotContains(t, entry.SeenComments, "c1", "unhandled comment stays unseen for retry")
	assert.NotEmpty(t, entry.History)
}

func assertErr() error { return assert.AnError }
