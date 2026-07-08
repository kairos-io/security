package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCVEPolicy(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cve-policy.yaml")
	os.WriteFile(path, []byte(
		"accepted-components:\n  openssl-fips:\n    reason: FIPS pinned\n"), 0o644)

	p, err := LoadCVEPolicy(path)
	if err != nil {
		t.Fatal(err)
	}
	if reason, ok := p.Accepted("openssl-fips"); !ok || reason != "FIPS pinned" {
		t.Fatalf("openssl-fips accepted lookup = (%q,%v)", reason, ok)
	}
	if _, ok := p.Accepted("openssl"); ok {
		t.Fatal("plain openssl must NOT be accepted")
	}

	// Missing file is valid (empty policy).
	empty, err := LoadCVEPolicy(filepath.Join(dir, "nope.yaml"))
	if err != nil {
		t.Fatalf("missing file should be valid: %v", err)
	}
	if _, ok := empty.Accepted("openssl-fips"); ok {
		t.Fatal("empty policy accepts nothing")
	}
}
