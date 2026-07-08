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
PR_BODY_FILE="${PR_BODY_FILE:-}"
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
# Stage only changes to TRACKED files (go.mod/go.sum + any files nib edited to
# fix the build). Never `git add -A`: LocalAI downloads backends/ and data/ into
# the checkout at runtime, and -A would commit that junk into the PR.
git add -u
git commit --signoff -m "$PR_TITLE" || { echo "nothing to commit"; exit 0; }

# Push with the token-authenticated remote. --force keeps an automation-owned
# reused branch in sync; the branch is intentionally force-updated.
if [ "$DRY_RUN" != "true" ]; then
  # Authenticate git (not just gh) with the provided token so pushes to a
  # reused PR branch trigger CI. Drop the checkout-persisted GITHUB_TOKEN
  # auth header so the token in the remote URL is the one used.
  git remote set-url origin "https://x-access-token:${TOKEN}@github.com/${GITHUB_REPOSITORY:?GITHUB_REPOSITORY is required}.git"
  git config --local --unset-all "http.https://github.com/.extraheader" 2>/dev/null || true
fi
run git push --force origin "HEAD:$BRANCH"

# PR body: use the generated summary file when present, else a static fallback.
body_flag=(--body "Automated dependency update opened by the update-deps action.")
if [ -n "$PR_BODY_FILE" ] && [ -s "$PR_BODY_FILE" ]; then
  body_flag=(--body-file "$PR_BODY_FILE")
fi

existing="$(open_pr_number "$BRANCH")"
if [ -n "$existing" ]; then
  echo "reusing open PR #$existing (branch force-updated)"
  # Refresh the existing PR's description with the new summary (best-effort).
  if [ "$DRY_RUN" != "true" ] && [ -n "$PR_BODY_FILE" ] && [ -s "$PR_BODY_FILE" ]; then
    gh pr edit "$BRANCH" --body-file "$PR_BODY_FILE" >/dev/null 2>&1 || true
  fi
  exit 0
fi

base_flag=()
[ -n "$BASE" ] && base_flag=(--base "$BASE")
run gh pr create --head "$BRANCH" "${base_flag[@]}" --title "$PR_TITLE" "${body_flag[@]}"

# Labels are best-effort: a label that doesn't exist in the repo (or missing
# Issues permission on the token) must not fail the run — the PR is what matters.
if [ "$DRY_RUN" != "true" ] && [ -n "$PR_LABELS" ]; then
  gh pr edit "$BRANCH" --add-label "$PR_LABELS" >/dev/null 2>&1 \
    || echo "note: could not apply labels '$PR_LABELS' (missing label or permission) — PR created without them"
fi
