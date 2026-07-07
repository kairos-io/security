#!/usr/bin/env bash
set -uo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tests/assert.sh disable=SC1091
. "$HERE/assert.sh"
# shellcheck source=scripts/lib.sh disable=SC1091
. "$HERE/../scripts/lib.sh"

# nib_base_url adds /v1 once and is idempotent.
assert_eq "http://localhost:8080/v1" "$(nib_base_url http://localhost:8080)"     "base url adds /v1"
assert_eq "http://localhost:8080/v1" "$(nib_base_url http://localhost:8080/)"    "base url trims trailing slash"
assert_eq "http://localhost:8080/v1" "$(nib_base_url http://localhost:8080/v1)"  "base url idempotent"

# Language config: go is known, others are rejected.
assert_ok   "go task known"        lang_nib_task go
assert_ok   "go fallback known"    lang_fallback_cmd go
assert_ok   "go verify known"      lang_verify_cmd go
assert_eq   "go build ./... && go vet ./..." "$(lang_verify_cmd go)" "go verify is build+vet"
assert_eq   $'go.mod\ngo.sum'      "$(lang_dep_paths go)"            "go dep paths are go.mod/go.sum"
assert_fail "python task rejected" lang_nib_task python

# has_dep_changes reflects git state in a throwaway repo.
TMP="$(mktemp -d)"; trap 'rm -rf "$TMP"; _report' EXIT
(
  cd "$TMP" || exit
  git init -q && git config user.email t@t && git config user.name t
  printf 'module x\n' > go.mod; printf '' > go.sum
  git add -A && git commit -qm init
)
assert_fail "no change -> has_dep_changes false" bash -c ". '$HERE/../scripts/lib.sh'; cd '$TMP'; has_dep_changes go"
printf 'module x\n// bump\n' > "$TMP/go.mod"
assert_ok   "modified go.mod -> has_dep_changes true" bash -c ". '$HERE/../scripts/lib.sh'; cd '$TMP'; has_dep_changes go"
