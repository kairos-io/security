#!/usr/bin/env bash
# Commit the dependency update, push the branch, and open or update the PR.
# In DRY_RUN mode, print the intended git/gh commands and make no changes.
set -uo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib.sh disable=SC1091
. "$HERE/lib.sh"

BRANCH="${BRANCH:-chore/update-deps}"
BASE="${BASE:-}"
PR_TITLE="${PR_TITLE:-chore(deps): update dependencies}"
PR_LABELS="${PR_LABELS:-dependencies}"
DRY_RUN="${DRY_RUN:-false}"
BOT_NAME="${BOT_NAME:-kairos-deps-bot}"
BOT_EMAIL="${BOT_EMAIL:-bot@kairos.io}"
TOKEN="${TOKEN:?TOKEN is required}"
export GH_TOKEN="$TOKEN"

run() { # echo + execute unless dry-run
  if [ "$DRY_RUN" = "true" ]; then echo "DRYRUN: $*"; else "$@"; fi
}

git config user.name  "$BOT_NAME"
git config user.email "$BOT_EMAIL"
git checkout -B "$BRANCH"
git add -A
git commit -m "$PR_TITLE" || { echo "nothing to commit"; exit 0; }

# Push with the token-authenticated remote. --force-with-lease keeps a reused
# branch in sync without clobbering unrelated pushes.
run git push --force-with-lease origin "HEAD:$BRANCH"

existing="$(open_pr_number "$BRANCH")"
if [ -n "$existing" ]; then
  echo "reusing open PR #$existing (branch force-updated)"
  exit 0
fi

base_flag=()
[ -n "$BASE" ] && base_flag=(--base "$BASE")
run gh pr create --head "$BRANCH" "${base_flag[@]}" --title "$PR_TITLE" \
  --label "$PR_LABELS" --body "Automated dependency update opened by the update-deps action (nib-driven)."
