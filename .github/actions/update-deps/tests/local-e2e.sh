#!/usr/bin/env bash
# Maintainer-run integration test. Requires: a running LocalAI at LOCALAI_URL
# with MODEL loaded, `nib` and `go` on PATH, and a Go repo path to test against.
#
#   LOCALAI_URL=http://localhost:8080 MODEL=gemma-4-e2b-it-qat-q4_0 \
#     bash tests/local-e2e.sh /path/to/go/repo
#
# It runs the update in an isolated worktree and asserts the manifest changed
# and the tree still builds. It NEVER pushes or opens a PR.
set -euo pipefail

REPO="${1:?usage: local-e2e.sh /path/to/go/repo}"
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WT="$(mktemp -d)/wt"
BR="local-e2e/$$"

cleanup() { git -C "$REPO" worktree remove --force "$WT" 2>/dev/null || true; git -C "$REPO" branch -D "$BR" 2>/dev/null || true; }
trap cleanup EXIT

git -C "$REPO" worktree add -b "$BR" "$WT" >/dev/null
( cd "$WT" && go build ./... ) && echo "baseline build OK"

LANGUAGE=go MODEL="${MODEL:?}" LOCALAI_URL="${LOCALAI_URL:-http://localhost:8080}" \
  bash -c "cd '$WT' && bash '$HERE/../scripts/run-update.sh'"

git -C "$WT" diff --quiet -- go.mod go.sum && { echo "FAIL: no dependency change produced"; exit 1; }
( cd "$WT" && go build ./... ) && echo "PASS: updated and builds"
