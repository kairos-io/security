package render

import (
	"strings"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestDashboardMarkdownLedgerSection(t *testing.T) {
	in := Input{Ledger: state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", Repo: "kairos-io/immucore",
			Package: "golang.org/x/net", State: "open", PRNumber: 412,
			PRURL: "https://github.com/kairos-io/immucore/pull/412",
			Bump:  state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}}
	md := DashboardMarkdown(in)
	assert.Contains(t, md, "Bot PR ledger")
	assert.Contains(t, md, "kairos-io/immucore")
	assert.Contains(t, md, "golang.org/x/net")
	assert.Contains(t, md, "0.33.0")
	assert.Contains(t, md, "open")
	assert.Contains(t, md, "412")
}

func TestDashboardMarkdownLedgerEmpty(t *testing.T) {
	md := DashboardMarkdown(Input{})
	assert.Contains(t, md, "No bot PRs yet")
	assert.True(t, strings.Contains(md, "Bot PR ledger"))
}
