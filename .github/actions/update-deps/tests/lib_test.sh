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

# resolve_nib_task: empty custom -> default per-language task; non-empty custom -> the custom prompt.
assert_eq "$(lang_nib_task go)" "$(resolve_nib_task go '')"          "resolve_nib_task falls back to default when custom empty"
assert_eq "bump only kairos-sdk" "$(resolve_nib_task go 'bump only kairos-sdk')" "resolve_nib_task uses the custom prompt"
# collapse_ws flattens newlines/tabs/runs to single spaces and trims.
assert_eq "one two three" "$(collapse_ws $'one   two\n\tthree')"     "collapse_ws squeezes whitespace"
assert_eq "line1 line2"   "$(resolve_nib_task go $'  line1\n  line2  ')" "resolve_nib_task flattens a multi-line custom prompt"

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

# localai_answers: mock curl to control the readyz HTTP code.
MOCK="$(mktemp -d)"
cat > "$MOCK/curl" <<'EOF'
#!/usr/bin/env bash
# emit the code requested via the sentinel env MOCK_CODE for -w '%{http_code}'
printf '%s' "${MOCK_CODE:-000}"
EOF
chmod +x "$MOCK/curl"
assert_ok   "localai_answers true on 200"  bash -c ". '$HERE/../scripts/lib.sh'; PATH=\"$MOCK:\$PATH\" MOCK_CODE=200 localai_answers http://x"
assert_fail "localai_answers false on 000" bash -c ". '$HERE/../scripts/lib.sh'; PATH=\"$MOCK:\$PATH\" MOCK_CODE=000 localai_answers http://x"
