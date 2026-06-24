package review

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBumps(t *testing.T) {
	diff := []byte("--- a/go.mod\n+++ b/go.mod\n" +
		"-\tgolang.org/x/net v0.30.0\n+\tgolang.org/x/net v0.33.0\n" +
		"-\tgithub.com/foo/bar v1.2.0 // indirect\n+\tgithub.com/foo/bar v1.3.0 // indirect\n")
	bumps := parseBumps(diff)
	assert.Equal(t, []DepBump{
		{Module: "golang.org/x/net", From: "0.30.0", To: "0.33.0"},
		{Module: "github.com/foo/bar", From: "1.2.0", To: "1.3.0"},
	}, bumps)
}

func TestModuleRepo(t *testing.T) {
	cases := map[string]struct {
		repo string
		ok   bool
	}{
		"github.com/foo/bar":       {"foo/bar", true},
		"github.com/foo/bar/v2":    {"foo/bar", true},
		"github.com/foo/bar/sub":   {"foo/bar", true},
		"golang.org/x/net":         {"golang/net", true},
		"k8s.io/api":               {"kubernetes/api", true},
		"sigs.k8s.io/yaml":         {"kubernetes-sigs/yaml", true},
		"example.com/vanity/thing": {"", false},
	}
	for mod, want := range cases {
		got, ok := moduleRepo(mod)
		assert.Equal(t, want.ok, ok, mod)
		assert.Equal(t, want.repo, got, mod)
	}
}
