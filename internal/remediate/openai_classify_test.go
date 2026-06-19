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

func TestOpenAIClassifierForcesToolCall(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), "classify_comment")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"choices":[{"message":{"tool_calls":[{"function":{"name":"classify_comment",`+
			`"arguments":"{\"intent\":\"request-change\",\"version\":\"0.36.0\",\"reply\":\"Bumping to 0.36.0.\"}"}}]}}]}`)
	}))
	defer ts.Close()

	c := &OpenAIClassifier{endpoint: ts.URL, model: "m", httpc: ts.Client()}
	got, err := c.Classify("bump x", "maintainer", "please pin to 0.36.0")
	require.NoError(t, err)
	assert.Equal(t, "request-change", got.Intent)
	assert.Equal(t, "0.36.0", got.Version)
	assert.Equal(t, "Bumping to 0.36.0.", got.Reply)
}

func TestNewOpenAIClassifierReadsConfig(t *testing.T) {
	c := NewOpenAIClassifier(config.AIConfig{})
	require.NotNil(t, c)
}
