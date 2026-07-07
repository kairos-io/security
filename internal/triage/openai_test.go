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
			`"arguments":"{\"focus\":[\"F1\"],\"summaries\":[{\"id\":\"F1\",\"summary\":\"crit\"}],\"narrative\":\"focus on a\"}"}}]}}]}`)
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
		io.WriteString(w, "{\"choices\":[{\"message\":{\"content\":\"```json\\n{\\\"focus\\\":[\\\"F1\\\"],"+
			"\\\"summaries\\\":[{\\\"id\\\":\\\"F1\\\",\\\"summary\\\":\\\"s\\\"}],\\\"narrative\\\":\\\"n\\\"}\\n```\"}}]}")
	}))
	defer ts.Close()

	focus, summaries, narrative, err := newTestClient(ts).Summarize(sampleForClient)
	require.NoError(t, err)
	assert.Equal(t, []string{"a"}, focus)
	assert.Equal(t, "s", summaries["a"])
	assert.Equal(t, "n", narrative)
}

// TestSummarizeUsesShortAliasesNotFingerprints guards the CI regression where a
// small model, forced to echo 64-char SHA-256 fingerprint ids into the focus
// array, overran max_tokens and truncated the tool-call JSON. The prompt must
// carry short aliases (F1, W1, …) and the model refers to those; we translate
// aliases back to the real ids on return.
func TestSummarizeUsesShortAliasesNotFingerprints(t *testing.T) {
	const fp = "02a6ba6130f153efe9eee50cc18c14ea5c25821be870c44cae32faedec50627e"
	cor := state.Correlated{
		Findings: []state.Finding{
			{ID: fp, Severity: "critical", CVEID: "CVE-1", Repo: "kairos-io/x", Package: "p"},
		},
		Waterfall: []state.WaterfallGroup{
			{ID: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				Severity: "high", RootCause: "openssl", Ecosystem: "alpine"},
		},
	}

	var promptContent string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req chatRequest
		body, _ := io.ReadAll(r.Body)
		require.NoError(t, json.Unmarshal(body, &req))
		promptContent = req.Messages[0].Content
		w.Header().Set("Content-Type", "application/json")
		// The model answers with the short aliases, never the fingerprints.
		io.WriteString(w, `{"choices":[{"message":{"tool_calls":[{"function":{"name":"report_triage",`+
			`"arguments":"{\"focus\":[\"F1\",\"W1\"],\"summaries\":[{\"id\":\"F1\",\"summary\":\"crit\"},`+
			`{\"id\":\"Z9\",\"summary\":\"ghost\"}],\"narrative\":\"n\"}"}}]}}]}`)
	}))
	defer ts.Close()

	focus, summaries, _, err := newTestClient(ts).Summarize(cor)
	require.NoError(t, err)

	// Prompt carries aliases, not the raw fingerprints.
	assert.NotContains(t, promptContent, fp, "prompt must not send 64-char fingerprints the model has to echo back")
	assert.Contains(t, promptContent, "F1")

	// Aliases are translated back to the real ids; hallucinated ids are dropped.
	assert.Equal(t, []string{fp, cor.Waterfall[0].ID}, focus)
	assert.Equal(t, "crit", summaries[fp])
	assert.NotContains(t, summaries, "Z9")
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
