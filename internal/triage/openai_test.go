package triage

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleForClient = state.Correlated{
	Findings: []state.Finding{
		{ID: "a", Severity: "critical", CVEID: "CVE-1", Repo: "kairos-io/x", Package: "p"},
	},
}

// newTestClient builds an OpenAIClient pointed at a test server.
func newTestClient(ts *httptest.Server) *OpenAIClient {
	return &OpenAIClient{endpoint: ts.URL, model: "m", maxTokens: 256, httpc: ts.Client()}
}

func TestSummarizeForcesToolCallAndParsesArguments(t *testing.T) {
	var gotReq chatRequest
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(body, &gotReq))
		// The model is asked to call exactly our tool.
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"choices":[{"message":{"tool_calls":[{"function":{"name":"report_triage",`+
			`"arguments":"{\"focus\":[\"a\"],\"summaries\":[{\"id\":\"a\",\"summary\":\"crit\"}],\"narrative\":\"focus on a\"}"}}]}}]}`)
	}))
	defer ts.Close()

	focus, summaries, narrative, err := newTestClient(ts).Summarize(sampleForClient)
	require.NoError(t, err)
	assert.Equal(t, []string{"a"}, focus)
	assert.Equal(t, "crit", summaries["a"])
	assert.Equal(t, "focus on a", narrative)

	// The request forced our tool.
	require.Len(t, gotReq.Tools, 1)
	assert.Equal(t, "report_triage", gotReq.Tools[0].Function.Name)
	tc, _ := gotReq.ToolChoice.(map[string]interface{})
	assert.Equal(t, "function", tc["type"])
}

func TestSummarizeFallsBackToContentJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// No tool_calls; JSON wrapped in a markdown fence in the content.
		io.WriteString(w, "{\"choices\":[{\"message\":{\"content\":\"```json\\n{\\\"focus\\\":[\\\"a\\\"],"+
			"\\\"summaries\\\":[{\\\"id\\\":\\\"a\\\",\\\"summary\\\":\\\"s\\\"}],\\\"narrative\\\":\\\"n\\\"}\\n```\"}}]}")
	}))
	defer ts.Close()

	focus, summaries, narrative, err := newTestClient(ts).Summarize(sampleForClient)
	require.NoError(t, err)
	assert.Equal(t, []string{"a"}, focus)
	assert.Equal(t, "s", summaries["a"])
	assert.Equal(t, "n", narrative)
}

func TestSummarizeErrorsOnHTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "model not loaded", http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	_, _, _, err := newTestClient(ts).Summarize(sampleForClient)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 503")
}
