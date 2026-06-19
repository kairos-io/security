package remediate

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAIProseReturnsContent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), "golang.org/x/net")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"choices":[{"message":{"content":"  Fixes a high-severity flaw in x/net.  "}}]}`)
	}))
	defer ts.Close()

	p := &OpenAIProse{endpoint: ts.URL, model: "m", httpc: ts.Client()}
	got, err := p.DraftPRBody(sampleIntent())
	require.NoError(t, err)
	assert.Equal(t, "Fixes a high-severity flaw in x/net.", got, "content is trimmed")
}

func TestOpenAIProseEmptyEndpointErrors(t *testing.T) {
	p := NewOpenAIProse(config.AIConfig{})
	_, err := p.DraftPRBody(sampleIntent())
	require.Error(t, err)
}

func TestOpenAIProseHTTPErrorIsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "model not loaded", http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	p := &OpenAIProse{endpoint: ts.URL, model: "m", httpc: ts.Client()}
	_, err := p.DraftPRBody(sampleIntent())
	require.Error(t, err)
}
