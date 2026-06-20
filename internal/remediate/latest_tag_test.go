package remediate

import "testing"

func TestLatestTag(t *testing.T) {
	// A `go list -m -versions` line: module path first, then ascending tags.
	// latestTag must skip the module path and return the highest vN.N.N token,
	// regardless of ordering.
	got := latestTag([]byte("github.com/kairos-io/kairos-sdk v0.0.1 v0.8.1 v0.10.0 v0.2.0"))
	if got != "v0.10.0" {
		t.Fatalf("latestTag picked %q, want v0.10.0", got)
	}

	// No tags published yet: `go list -m -versions` prints only the module path.
	if got := latestTag([]byte("github.com/kairos-io/kairos-sdk")); got != "" {
		t.Fatalf("latestTag with no tags returned %q, want \"\"", got)
	}

	// Empty output also yields "".
	if got := latestTag(nil); got != "" {
		t.Fatalf("latestTag(nil) returned %q, want \"\"", got)
	}
}
