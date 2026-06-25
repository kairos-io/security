package review

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestCompareRef(t *testing.T) {
	assert.Equal(t, "v0.33.0", compareRef("0.33.0"))
	assert.Equal(t, "fab4fdf2f2f3", compareRef("0.0.0-20241017190036-fab4fdf2f2f3"))
	assert.Equal(t, "abcdef123456", compareRef("1.2.3-0.20240101000000-abcdef123456"))
	assert.Equal(t, "v2.0.0+incompatible", compareRef("2.0.0+incompatible")) // not a pseudo-version
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

func TestParseCompareURLs(t *testing.T) {
	body := "Release notes\n" +
		"[Compare Source](https://redirect.github.com/lucide-icons/lucide/compare/0.576.0...0.577.0)\n" +
		"Full Changelog: <https://github.com/lucide-icons/lucide/compare/0.468.0...0.577.0>\n" +
		"dup: https://github.com/lucide-icons/lucide/compare/0.576.0...0.577.0\n"
	got := parseCompareURLs(body)
	// deduped: two distinct compares
	assert.Len(t, got, 2)
	assert.Equal(t, CompareRef{Repo: "lucide-icons/lucide", Base: "0.576.0", Head: "0.577.0",
		Label: "lucide-icons/lucide 0.576.0..0.577.0 (PR body)"}, got[0])
	assert.Equal(t, "0.468.0", got[1].Base)
	assert.Equal(t, "0.577.0", got[1].Head)
}

func TestParseCompareURLsNone(t *testing.T) {
	assert.Empty(t, parseCompareURLs("no links here"))
}

func TestCompareTargetsUnifiesAndCaps(t *testing.T) {
	// a Go bump (via go.mod diff) + a body compare link → both, deduped
	diff := []byte("--- a/go.mod\n+++ b/go.mod\n" +
		"-\tgithub.com/foo/bar v1.2.0\n+\tgithub.com/foo/bar v1.3.0\n")
	body := "https://github.com/baz/qux/compare/v2.0.0...v2.1.0"
	got := compareTargets(diff, body)
	require.Len(t, got, 2)
	assert.Equal(t, "foo/bar", got[0].Repo) // Go path first
	assert.Equal(t, "v1.2.0", got[0].Base)  // compareRef adds v
	assert.Equal(t, "v1.3.0", got[0].Head)
	assert.Equal(t, "baz/qux", got[1].Repo) // body link
	assert.Equal(t, "v2.0.0", got[1].Base)  // verbatim from URL
}

func TestParseCompareURLsRejectsLookalikeHost(t *testing.T) {
	// a look-alike host must not be parsed as a github.com compare target
	assert.Empty(t, parseCompareURLs("https://attacker-github.com/evil/repo/compare/a...b"))
	// scheme-required: a bare host reference is ignored
	assert.Empty(t, parseCompareURLs("see github.com/foo/bar/compare/a...b for details"))
	// legitimate subdomains still match
	got := parseCompareURLs("https://redirect.github.com/o/r/compare/v1...v2")
	assert.Len(t, got, 1)
	assert.Equal(t, "o/r", got[0].Repo)
}

func TestParseCompareURLsTwoDot(t *testing.T) {
	// renovate uses a two-dot range for action/image digest bumps
	got := parseCompareURLs("bumps the action; see " +
		"https://redirect.github.com/docker/login-action/compare/4907a6ddec9925e35a0a9e82d7399ccc52663121..650006c6eb7dba73a995cc03b0b2d7f5ca915bee for details")
	require.Len(t, got, 1)
	assert.Equal(t, "docker/login-action", got[0].Repo)
	assert.Equal(t, "4907a6ddec9925e35a0a9e82d7399ccc52663121", got[0].Base)
	assert.Equal(t, "650006c6eb7dba73a995cc03b0b2d7f5ca915bee", got[0].Head)
}
