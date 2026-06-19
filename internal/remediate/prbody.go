package remediate

import (
	"fmt"
	"strings"
)

func PRMarker(key string) string { return "<!-- ksec:key=" + key + " -->" }

func slug(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteByte('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

func BranchName(in Intent) string {
	return "ksec/bump-" + slug(in.Package) + "-" + slug(in.Bump.To)
}

func PRTitle(in Intent) string {
	return fmt.Sprintf("chore(security): bump %s to %s", in.Package, in.Bump.To)
}

func PRBody(in Intent) string {
	return fmt.Sprintf(`## Automated security bump

Bumps **%s** to **%s** to address a %s-severity vulnerability detected by
[kairos-security](https://github.com/kairos-io/security).

- Package: `+"`%s`"+`
- Target version: `+"`%s`"+`
- Severity: %s

This PR was opened automatically. The change is a deterministic
`+"`go get %s@%s` + `go mod tidy`"+`; CI on this PR runs the repository's tests.

%s`, in.Package, in.Bump.To, in.Severity, in.Package, in.Bump.To, in.Severity,
		in.Package, in.Bump.To, PRMarker(in.Key))
}
