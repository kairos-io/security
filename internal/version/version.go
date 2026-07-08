// Package version compares OS-package version strings (dotted-numeric, e.g.
// "3.1.2", "2.86.2"). It is deliberately NOT a full SemVer implementation:
// these are Alpine/upstream package versions, and only ordering by numeric
// dotted segments is needed.
package version

import (
	"strconv"
	"strings"
)

// Compare returns -1 if a<b, 0 if equal, 1 if a>b. An empty string is the
// lowest possible value. Each dot-separated segment is compared by its leading
// integer; a segment whose leading token is not an integer sorts below one that
// is. Missing trailing segments are treated as 0.
func Compare(a, b string) int {
	as, bs := strings.Split(a, "."), strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		av, aok := segInt(as, i)
		bv, bok := segInt(bs, i)
		switch {
		case aok && !bok:
			return 1
		case !aok && bok:
			return -1
		case av < bv:
			return -1
		case av > bv:
			return 1
		}
	}
	return 0
}

// segInt returns the leading integer of segment i and whether that value is
// "numeric-comparable". A missing trailing segment is a numeric 0. A present
// but non-numeric segment is not numeric-comparable (sorts lowest).
func segInt(seg []string, i int) (int, bool) {
	if i >= len(seg) {
		return 0, true // trailing missing == 0
	}
	s := seg[i]
	j := 0
	for j < len(s) && s[j] >= '0' && s[j] <= '9' {
		j++
	}
	if j == 0 {
		return 0, false
	}
	n, _ := strconv.Atoi(s[:j])
	return n, true
}
