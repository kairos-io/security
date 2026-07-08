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
	}
	for _, c := range cases {
		if got := Compare(c.a, c.b); got != c.want {
			t.Errorf("Compare(%q,%q)=%d want %d", c.a, c.b, got, c.want)
		}
	}
}
