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

func CascadeBranchName(in Intent) string {
	return "ksec/cascade-" + slug(in.Package) + "-pseudo"
}

func CascadePRBody(in Intent) string {
	return fmt.Sprintf(`## Automated security cascade

This bumps **%s** to a pseudo-version of its latest default-branch commit, which
contains an unreleased security fix. Once a maintainer cuts a release tag for
that module, kairos-security will re-pin this PR to the tagged version.

- Module: `+"`%s`"+`
- Severity: %s
- Please **tag a release** of the upstream module so this can be pinned cleanly.

This PR was opened automatically by kairos-security. CI on this PR runs the tests.

%s`, in.Package, in.Package, in.Severity, PRMarker(in.Key))
}
