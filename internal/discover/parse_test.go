package discover

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDeps(t *testing.T) {
	makefile := []byte("AGENT_VERSION?=v1.2.3\nIMMUCORE_VERSION?=v0.5.0\nEDGEVPN_VERSION?=v0.30.0\n")
	gomod := []byte(`
module github.com/kairos-io/kairos-init
go 1.22
require (
	github.com/kairos-io/kairos-sdk v0.7.0
	github.com/mudler/yip v1.9.0
	github.com/mauromorales/xpasswd v0.3.0
	github.com/spf13/cobra v1.8.0
)
`)
	got := ParseDeps(makefile, gomod)
	assert.Equal(t, []string{
		"kairos-io/immucore",
		"kairos-io/kairos-agent",
		"kairos-io/kairos-sdk",
		"mauromorales/xpasswd",
		"mudler/edgevpn",
		"mudler/yip",
	}, got)
}
