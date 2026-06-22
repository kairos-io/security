package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSourceScanEnabled(t *testing.T) {
	assert.True(t, Repo{}.SourceScanEnabled()) // default
	f := false
	assert.False(t, Repo{Scan: ScanConfig{Source: &f}}.SourceScanEnabled())
	tr := true
	assert.True(t, Repo{Scan: ScanConfig{Source: &tr}}.SourceScanEnabled())
}
