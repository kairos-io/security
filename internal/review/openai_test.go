package review

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// toolCallResponse builds an OpenAI-compatible chat response carrying a single
// forced tool call whose arguments are the given verdict/reasoning JSON.
func toolCallResponse(t *testing.T, verdict, reasoning string) []byte {
	t.Helper()
	args, err := json.Marshal(map[string]string{
		"verdict":        verdict,
		"reasoning":      reasoning,
		"changesSummary": "summary of " + verdict,
	})
	require.NoError(t, err)
	resp := map[string]any{
		"choices": []any{map[string]any{
			"message": map[string]any{
				"tool_calls": []any{map[string]any{
					"function": map[string]any{
						"name":      assessToolName,
						"arguments": string(args),
					},
				}},
			},
		}},
	}
	b, err := json.Marshal(resp)
	require.NoError(t, err)
	return b
}

func TestOpenAIAssessorReturnsVerdict(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(toolCallResponse(t, "bad", "dependency looks malicious"))
	}))
	defer srv.Close()

	a := NewOpenAIAssessor(config.AIConfig{Nib: config.NibCfg{Endpoint: srv.URL, Model: "m"}})
	verdict, reasoning, summary, err := a.Assess(ghclient.PullRequest{Title: "bump"}, "review context")
	require.NoError(t, err)
	assert.Equal(t, "bad", verdict)
	assert.Equal(t, "dependency looks malicious", reasoning)
	assert.Equal(t, "summary of bad", summary)
}

func TestOpenAIAssessorEmptyEndpointDegrades(t *testing.T) {
	a := NewOpenAIAssessor(config.AIConfig{})
	verdict, reasoning, _, err := a.Assess(ghclient.PullRequest{Title: "x"}, "ctx")
	require.NoError(t, err) // never a hard error
	assert.Equal(t, "needs_human_verification", verdict)
	assert.NotEmpty(t, reasoning)
}

func TestOpenAIAssessorNon200Degrades(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	a := NewOpenAIAssessor(config.AIConfig{Nib: config.NibCfg{Endpoint: srv.URL, Model: "m"}})
	verdict, _, _, err := a.Assess(ghclient.PullRequest{Title: "x"}, "ctx")
	require.NoError(t, err)
	assert.Equal(t, "needs_human_verification", verdict)
}

func TestOpenAIAssessorNoToolCallDegrades(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"hi"}}]}`))
	}))
	defer srv.Close()

	a := NewOpenAIAssessor(config.AIConfig{Nib: config.NibCfg{Endpoint: srv.URL, Model: "m"}})
	verdict, _, _, err := a.Assess(ghclient.PullRequest{Title: "x"}, "ctx")
	require.NoError(t, err)
	assert.Equal(t, "needs_human_verification", verdict)
}

func TestOpenAIAssessorOutOfEnumDegrades(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(toolCallResponse(t, "maybe", "unsure"))
	}))
	defer srv.Close()

	a := NewOpenAIAssessor(config.AIConfig{Nib: config.NibCfg{Endpoint: srv.URL, Model: "m"}})
	verdict, _, _, err := a.Assess(ghclient.PullRequest{Title: "x"}, "ctx")
	require.NoError(t, err)
	assert.Equal(t, "needs_human_verification", verdict)
}
