// Package version compares OS-package version strings (dotted-numeric, e.g.
// "3.1.2", "2.86.2"). It is deliberately NOT a full SemVer implementation:
// these are Alpine/upstream package versions, and only ordering by numeric
// dotted segments is needed.
package version

import (
	"strconv"
	"strings"
)

// Compare returns -1 if a<b, 0 if equal, 1 if a>b.
//
// A version is read as three parts: a dotted-numeric release core, an optional
// pre-release suffix, and an optional Alpine package revision.
//
//   - The core is compared per dot-separated segment; within a segment,
//     maximal runs of digits compare numerically and runs of non-digits compare
//     bytewise ("natural order"), so "1.1.1n" < "1.1.1t" and "1.10" > "1.9". A
//     missing trailing segment is treated as "0", so "1.2" == "1.2.0".
//   - A leading "v" or "go" before a digit is a prefix, not part of the
//     version: Go module versions ("v0.30.0") and Go toolchain versions
//     ("go1.25.1") order against the bare numbers the advisories carry.
//   - A pre-release suffix sorts BELOW the bare release, as SemVer says, so
//     "1.26.0-0" < "1.26.0" < "1.26.1". Go advisories put "-0" on an
//     "introduced" boundary to mean "from the first 1.26.0 pre-release on".
//   - An Alpine package revision ("-r5") sorts ABOVE the bare release, because
//     it is a rebuild of it: "3.1.4" < "3.1.4-r5" < "3.1.4-r10".
//
// An empty string is the lowest value. This is deliberately NOT full SemVer —
// these are OS/upstream package versions.
func Compare(a, b string) int {
	ac, apre, arev := parse(a)
	bc, bpre, brev := parse(b)

	if c := compareDotted(ac, bc); c != 0 {
		return c
	}
	// Same release core: a pre-release is below the release it leads to.
	switch {
	case apre != "" && bpre == "":
		return -1
	case apre == "" && bpre != "":
		return 1
	case apre != "" && bpre != "":
		if c := compareDotted(apre, bpre); c != 0 {
			return c
		}
	}
	// Same release: an Alpine revision is above the release it rebuilds.
	switch {
	case arev != "" && brev == "":
		return 1
	case arev == "" && brev != "":
		return -1
	case arev != "" && brev != "":
		return natCompare(arev, brev)
	}
	return 0
}

// parse splits a version into its release core, pre-release suffix and Alpine
// package revision. The revision is taken off first so that a version carrying
// both ("1.26.0-rc.1-r0") keeps them apart. The returned revision keeps its "r"
// so natCompare orders "r9" below "r10".
func parse(s string) (core, pre, rev string) {
	s = stripBuildPrefix(s)
	s, rev = splitAlpineRevision(s)
	if i := strings.IndexByte(s, '-'); i > 0 {
		return s[:i], s[i+1:], rev
	}
	return s, "", rev
}

// stripBuildPrefix removes the "v" of a Go module version and the "go" of a Go
// toolchain version. Both are written by govulncheck ("v0.30.0", "go1.25.1")
// while the advisory's fixed version is bare ("0.33.0"), and a leading letter
// compares ABOVE every digit in natCompare — which would order the vulnerable
// version above its own fix and hide the finding as already-fixed.
func stripBuildPrefix(s string) string {
	switch {
	case len(s) > 1 && s[0] == 'v' && isDigit(s[1]):
		return s[1:]
	case len(s) > 2 && s[0] == 'g' && s[1] == 'o' && isDigit(s[2]):
		return s[2:]
	}
	return s
}

// splitAlpineRevision splits a trailing Alpine package-revision suffix
// ("3.6.4-r0" -> "3.6.4", "r0"), leaving other version strings untouched.
func splitAlpineRevision(s string) (string, string) {
	i := strings.LastIndex(s, "-r")
	if i <= 0 {
		return s, ""
	}
	if _, err := strconv.Atoi(s[i+2:]); err != nil {
		return s, ""
	}
	return s[:i], s[i+1:]
}

// compareDotted compares two dot-separated strings segment by segment, in
// natural order. A missing trailing segment is a numeric 0.
func compareDotted(a, b string) int {
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
