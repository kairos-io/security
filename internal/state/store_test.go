package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	in := Findings{
		Findings: []Finding{{ID: "a", Repo: "kairos-io/x", Type: "sourceCVE", Severity: "high"}},
		Errors:   []CollectionError{{Repo: "kairos-io/y", Collector: "prs", Message: "boom"}},
	}
	require.NoError(t, Save(dir, FindingsFile, in))

	var out Findings
	require.NoError(t, Load(dir, FindingsFile, &out))
	assert.Equal(t, in, out)
}

func TestSaveIsStableIndentedJSON(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, Save(dir, RingFileForTest, map[string]int{"b": 2, "a": 1}))
	b, err := os.ReadFile(filepath.Join(dir, RingFileForTest))
	require.NoError(t, err)
	// keys sorted, 2-space indent, trailing newline
	assert.Equal(t, "{\n  \"a\": 1,\n  \"b\": 2\n}\n", string(b))
}

const RingFileForTest = "scratch.json"

func TestLoadMissingFileErrors(t *testing.T) {
	var out Findings
	assert.Error(t, Load(t.TempDir(), "nope.json", &out))
}
