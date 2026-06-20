package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedgerByKey(t *testing.T) {
	l := Ledger{Entries: []LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", State: "open"},
		{Key: "kairos-io/kairos|stdlib", State: "merged"},
	}}
	got, ok := l.ByKey("kairos-io/kairos|stdlib")
	require.True(t, ok)
	assert.Equal(t, "merged", got.State)

	got.State = "closed" // pointer write must mutate the slice
	again, _ := l.ByKey("kairos-io/kairos|stdlib")
	assert.Equal(t, "closed", again.State)

	_, ok = l.ByKey("nope")
	assert.False(t, ok)
}

func TestLedgerRoundTrip(t *testing.T) {
	dir := t.TempDir()
	in := Ledger{Entries: []LedgerEntry{{Key: "a|b", Repo: "a", Package: "b", State: "open",
		Source: "dependabot", Kind: "direct", Blocked: "checks failing", NeedsHuman: true,
		CascadeFrom: "kairos-io/kairos-sdk|golang.org/x/net", PinTarget: "v0.8.1", Pseudo: true,
		Bump: Bump{Package: "b", To: "1.2.3"}, History: []LedgerEvent{{Run: "2026-06-20", Action: "opened"}}}}}
	require.NoError(t, Save(dir, LedgerFile, in))
	var out Ledger
	require.NoError(t, Load(dir, LedgerFile, &out))
	assert.Equal(t, in, out)
}
