package collect

import (
	"bytes"
	"fmt"
	"strings"
)

// ClassifyGovulncheck decides whether a govulncheck run that exited non-zero
// failed for real (build/load error) or merely found vulnerabilities. In -json
// mode govulncheck emits config/progress objects on stdout before analysis, so
// a non-empty stdout does NOT mean it succeeded — only the presence of an
// "osv"/"finding" record does. A non-zero exit with no such record is a real
// failure and must surface (it is how a Go-toolchain mismatch was silently
// reported as zero vulnerabilities).
func ClassifyGovulncheck(stdout, stderr []byte, runErr error) ([]byte, error) {
	if runErr == nil {
		return stdout, nil
	}
	if bytes.Contains(stdout, []byte(`"osv"`)) || bytes.Contains(stdout, []byte(`"finding"`)) {
		return stdout, nil // vulnerabilities found: non-zero exit is expected
	}
	return nil, fmt.Errorf("govulncheck: %v: %s", runErr, truncErr(stderr, 240))
}

// truncErr renders tool stderr as a one-line, length-capped summary so a
// build-failure (e.g. a cgo app missing system libs) doesn't flood the
// dashboard or the committed findings.json with a multi-KB wall.
func truncErr(b []byte, max int) string {
	s := strings.Join(strings.Fields(string(b)), " ") // collapse all whitespace/newlines
	if len(s) <= max {
		return s
	}
	return s[:max] + " … (truncated)"
}
