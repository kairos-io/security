package version

import "testing"

func TestCompare(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"3.1.2", "3.1.4", -1},
		{"3.5.6", "3.1.2", 1},
		{"2.86.2", "2.66.6", 1}, // glib: current newer than "fixed"
		{"2.15.3", "2.13.8", 1}, // libxml2: current newer than "fixed"
		{"3.1.2", "3.3.1", -1},  // openssl: current older than fix
		{"1.2.3", "1.2.3", 0},
		{"1.2", "1.2.0", 0}, // missing trailing segment == 0
		{"1.2.0", "1.2", 0},
		{"1.10.0", "1.9.0", 1}, // numeric, not lexical
		{"", "0.0.1", -1},      // empty is lowest
		{"1.0.0", "", 1},
		// suffix-only differences must NOT collapse to equal (C1 regression)
		{"1.1.1n", "1.1.1t", -1},
		{"1.1.1t", "1.1.1n", 1},
		{"3.1.4-r5", "3.1.4-r6", -1},
		{"3.1.4-r9", "3.1.4-r10", -1},
		{"3.1.4-r10", "3.1.4-r9", 1},
		{"1.1.1w", "1.1.1w", 0},
	}
	for _, c := range cases {
		if got := Compare(c.a, c.b); got != c.want {
			t.Errorf("Compare(%q,%q)=%d want %d", c.a, c.b, got, c.want)
		}
	}
}

// TestCompareGoVersionPrefix pins the ordering of the version strings
// govulncheck writes into a sourceCVE finding against the bare version the
// advisory carries. A leading letter compares above every digit in natural
// order, so before this was handled every reachable Go CVE ordered as
// "current >= fixed" and classify.Apply filed it as already-fixed.
func TestCompareGoVersionPrefix(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"v0.30.0", "0.33.0", -1}, // the repo's own govulncheck fixture
		{"v0.33.0", "0.33.0", 0},
		{"v0.34.0", "0.33.0", 1},
		{"go1.25.1", "1.25.3", -1}, // stdlib, as govulncheck reports it
		{"go1.25.3", "1.25.3", 0},
		{"go1.26.0", "1.25.3", 1},
		{"v1.2.3", "v1.2.4", -1}, // both sides prefixed
		{"v", "0.1", 1},          // a bare "v" is not a prefix, it stays a letter
		{"go", "0.1", 1},         // nor is a bare "go"
	}
	for _, c := range cases {
		if got := Compare(c.a, c.b); got != c.want {
			t.Errorf("Compare(%q,%q)=%d want %d", c.a, c.b, got, c.want)
		}
	}
}

// TestComparePreRelease pins a pre-release below the release it leads to, and
// an Alpine package revision above it. Go advisories write the "introduced"
// boundary of a whole minor line as "1.26.0-0", so ordering 1.26.0 below that
// boundary dropped the range and reported a vulnerable version as unaffected.
func TestComparePreRelease(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.26.0", "1.26.0-0", 1}, // GO-2026-4601's introduced boundary
		{"1.26.0-0", "1.26.0", -1},
		{"1.26.0-0", "1.26.1", -1},
		{"1.26.0", "1.26.0-rc.1", 1},
		{"1.26.0-rc.1", "1.26.0-rc.2", -1},
		{"1.26.0-rc.2", "1.26.0-rc.10", -1},
		{"1.26.0-beta", "1.26.0-rc", -1},
		{"1.25.9", "1.26.0-0", -1}, // still below the whole line
		// an Alpine revision is a rebuild of the release, so it stays above it
		{"3.1.4-r5", "3.1.4", 1},
		{"3.1.4", "3.1.4-r5", -1},
		{"3.1.4-r5", "3.1.4-rc.1", 1}, // revision beats pre-release
		{"1.26.0-rc.1-r0", "1.26.0-rc.1", 1},
		{"1.26.0-rc.1-r0", "1.26.0", -1},
	}
	for _, c := range cases {
		if got := Compare(c.a, c.b); got != c.want {
			t.Errorf("Compare(%q,%q)=%d want %d", c.a, c.b, got, c.want)
		}
	}
}

// TestCompareIsAntisymmetric guards the three-part split: every pair must
// order the same way read backwards, or a sort over these versions is
// undefined.
func TestCompareIsAntisymmetric(t *testing.T) {
	vs := []string{
		"", "0", "1.2", "1.2.0", "1.10", "1.9", "1.1.1n", "1.1.1t",
		"3.1.4", "3.1.4-r5", "3.1.4-r10", "1.26.0-0", "1.26.0",
		"1.26.0-rc.1", "1.26.1", "v0.30.0", "0.33.0", "go1.25.1", "1.25.3",
	}
	for _, a := range vs {
		for _, b := range vs {
			if got, rev := Compare(a, b), Compare(b, a); got != -rev {
				t.Errorf("Compare(%q,%q)=%d but Compare(%q,%q)=%d", a, b, got, b, a, rev)
			}
		}
	}
}
