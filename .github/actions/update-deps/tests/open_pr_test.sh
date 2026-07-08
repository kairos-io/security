#!/usr/bin/env bash
set -uo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tests/assert.sh disable=SC1091
. "$HERE/assert.sh"

SCRIPT="$HERE/../scripts/open-pr.sh"

# Mock dir: gh/git record their args to $CALLS; gh pr list returns $EXISTING.
MOCK="$(mktemp -d)"; CALLS="$MOCK/calls"; : > "$CALLS"
cat > "$MOCK/gh" <<'EOF'
#!/usr/bin/env bash
echo "gh $*" >> "$CALLS_FILE"
if [ "$1 $2" = "pr list" ]; then printf '%s' "${EXISTING:-}"; fi
EOF
cat > "$MOCK/git" <<'EOF'
#!/usr/bin/env bash
echo "git $*" >> "$CALLS_FILE"
EOF
chmod +x "$MOCK/gh" "$MOCK/git"

# Dry-run makes no gh create/push calls.
: > "$CALLS"
CALLS_FILE="$CALLS" DRY_RUN=true BRANCH=chore/update-deps BASE=main PR_TITLE=t PR_LABELS=deps TOKEN=x \
  PATH="$MOCK:$PATH" bash "$SCRIPT" >/dev/null 2>&1
assert_eq "0" "$(grep -c 'gh pr create' "$CALLS")" "dry-run does not create a PR"
assert_eq "0" "$(grep -c 'git push'    "$CALLS")" "dry-run does not push"
assert_eq "0" "$(grep -c 'git remote set-url' "$CALLS")" "dry-run does not re-auth the remote"

# Live, no existing PR -> pushes and creates.
: > "$CALLS"
CALLS_FILE="$CALLS" EXISTING="" DRY_RUN=false BRANCH=chore/update-deps BASE=main PR_TITLE=t PR_LABELS=deps TOKEN=x \
  GITHUB_REPOSITORY=owner/repo \
  PATH="$MOCK:$PATH" bash "$SCRIPT" >/dev/null 2>&1
assert_eq "1" "$(grep -c 'git push'    "$CALLS")" "live push happens"
assert_eq "1" "$(grep -c 'gh pr create' "$CALLS")" "live create happens when no PR exists"

# Live, existing PR #42 -> pushes but does NOT create a duplicate.
: > "$CALLS"
CALLS_FILE="$CALLS" EXISTING="42" DRY_RUN=false BRANCH=chore/update-deps BASE=main PR_TITLE=t PR_LABELS=deps TOKEN=x \
  GITHUB_REPOSITORY=owner/repo \
  PATH="$MOCK:$PATH" bash "$SCRIPT" >/dev/null 2>&1
assert_eq "1" "$(grep -c 'git push'     "$CALLS")" "live push happens with existing PR"
assert_eq "0" "$(grep -c 'gh pr create' "$CALLS")" "no duplicate PR when one is open"
