// Package version compares OS-package version strings (dotted-numeric, e.g.
// "3.1.2", "2.86.2"). It is deliberately NOT a full SemVer implementation:
// these are Alpine/upstream package versions, and only ordering by numeric
// dotted segments is needed.
package version

import (
	"strconv"
	"strings"
)

// Compare returns -1 if a<b, 0 if equal, 1 if a>b. Versions are compared per
// dot-separated segment; within a segment, maximal runs of digits compare
// numerically and runs of non-digits compare bytewise ("natural order"), so
// "1.1.1n" < "1.1.1t", "3.1.4-r5" < "3.1.4-r6", "3.1.4-r9" < "3.1.4-r10", and
// "1.10" > "1.9". A missing trailing segment is treated as "0" (so "1.2" ==
// "1.2.0"). An empty string is the lowest value. This is deliberately NOT full
// SemVer — these are OS/upstream package versions.
func Compare(a, b string) int {
	as, bs := strings.Split(a, "."), strings.Split(b, ".")
	n := len(as)
	if len(bs) > n {
		n = len(bs)
	}
	for i := 0; i < n; i++ {
		sa, sb := "0", "0" // a MISSING trailing segment is a numeric 0
		if i < len(as) {
			sa = as[i]
		}
		if i < len(bs) {
			sb = bs[i]
		}
		if c := natCompare(sa, sb); c != 0 {
			return c
		}
	}
	return 0
}

// natCompare compares two segment strings in natural order: aligned maximal
// digit runs compare numerically, everything else bytewise. When one string is
// a prefix of the other, the longer one is greater.
func natCompare(a, b string) int {
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if isDigit(a[i]) && isDigit(b[j]) {
			ai, ni := takeDigits(a, i)
			bj, nj := takeDigits(b, j)
			if ai != bj {
				if ai < bj {
					return -1
				}
				return 1
			}
			i, j = ni, nj
			continue
		}
		if a[i] != b[j] {
			if a[i] < b[j] {
				return -1
			}
			return 1
		}
		i++
		j++
	}
	switch {
	case i < len(a):
		return 1 // a has leftover -> longer -> greater
	case j < len(b):
		return -1
	default:
		return 0
	}
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// takeDigits parses the maximal digit run starting at i and returns its value
// and the index just past it.
func takeDigits(s string, i int) (int, int) {
	j := i
	for j < len(s) && isDigit(s[j]) {
		j++
	}
	n, _ := strconv.Atoi(s[i:j])
	return n, j
}
