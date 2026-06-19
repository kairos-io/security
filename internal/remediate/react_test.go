package remediate

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
