package remediate

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummarizeLedgerReturnsContent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), "kairos-io/immucore")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"choices":[{"message":{"content":"  One PR is open and cascading.  "}}]}`)
	}))
	defer ts.Close()

	cfg := config.AIConfig{}
	cfg.Nib.Endpoint = ts.URL
	cfg.Nib.Model = "m"
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Repo: "kairos-io/immucore", Package: "golang.org/x/net", State: "open", Kind: "direct", Source: "ksec", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	got, err := SummarizeLedger(cfg, led)
	require.NoError(t, err)
	assert.Equal(t, "One PR is open and cascading.", got, "content is trimmed")
}

func TestSummarizeLedgerEmptyEndpointErrors(t *testing.T) {
	_, err := SummarizeLedger(config.AIConfig{}, state.Ledger{})
	require.Error(t, err)
}
