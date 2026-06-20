package render

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestDashboardMarkdownLedgerShowsSourceAndKind(t *testing.T) {
	in := Input{Ledger: state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", Repo: "kairos-io/immucore",
			Package: "golang.org/x/net", State: "open", PRNumber: 7,
			PRURL: "https://github.com/kairos-io/immucore/pull/7",
			Source: "dependabot", Kind: "direct", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "dependabot")
	assert.Contains(t, md, "direct")
}
