# nib Dependency-Update Composite Action — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Ship a reusable composite action in `kairos-io/security` that starts LocalAI, has `nib` bump a repo's dependencies to latest (with a deterministic fallback), verifies the build, and opens/updates a CI-triggering PR.

**Architecture:** The action (`.github/actions/update-deps/action.yml`) is a thin orchestrator whose steps call three focused bash scripts (`start-localai.sh`, `run-update.sh`, `open-pr.sh`) that share pure helper functions in `lib.sh`. Pure logic is unit-tested with mocked `gh`/`curl`; the full nib run is exercised by a maintainer-run local e2e against a real LocalAI. Consumers add a ~15-line caller workflow.

**Tech Stack:** GitHub composite action (YAML), Bash, `nib`, LocalAI binary, `gh` CLI, `shellcheck` + `actionlint` for static gates. Go is the only implemented `language` in v1.

## Global Constraints

- Spec: `docs/superpowers/specs/2026-07-07-nib-dependency-update-action-design.md` — this plan implements it.
- Action path (exact): `.github/actions/update-deps/`.
- Default model: `gemma-4-e2b-it` (overridable via `model` input).
- `nib` is invoked exactly as `nib --cli --yolo` with env `MODEL` / `BASE_URL` (must carry `/v1` suffix) / `API_KEY=sk-localai`, task fed on stdin as a single line + trailing newline (matches `internal/remediate/nib_agent.go`).
- Verify gate is **`go build ./... && go vet ./...`** only — never run the target repo's test suite inside the action.
- Never open a PR when: (a) no dependency-manifest change, or (b) verify still fails after one nib repair retry (in that case the action **fails**).
- The built-in `GITHUB_TOKEN` must NOT be used to open the PR (it suppresses CI); the `token` input is required and expected to be a GitHub App token or PAT.
- All scripts must pass `shellcheck` with no warnings; `action.yml` and caller workflow must pass `actionlint`.
- Only `language: go` is implemented; unknown languages must fail fast with a clear message. No other ecosystem code (YAGNI).
- Commit style: end commit messages with the repo's `Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>` trailer.

---

## File Structure

- Create `.github/actions/update-deps/action.yml` — composite action: inputs → steps calling the scripts.
- Create `.github/actions/update-deps/scripts/lib.sh` — pure helpers (no side effects on source): base-URL normalization, per-language task/fallback/verify/paths, dep-change detection, open-PR lookup. Sourced by the other scripts and by unit tests.
- Create `.github/actions/update-deps/scripts/start-localai.sh` — skip if the endpoint already answers; else download the `local-ai` binary, start it, wait until a real completion succeeds.
- Create `.github/actions/update-deps/scripts/run-update.sh` — primary nib run, deterministic fallback, verify gate + one repair retry; emits `changed`/`verified` outputs.
- Create `.github/actions/update-deps/scripts/open-pr.sh` — commit, branch, dedupe against an existing open PR, push, `gh pr create`; honors dry-run.
- Create `.github/actions/update-deps/tests/assert.sh` — minimal assertion + PATH-mock helpers for the unit tests.
- Create `.github/actions/update-deps/tests/lib_test.sh` — unit tests for `lib.sh`.
- Create `.github/actions/update-deps/tests/open_pr_test.sh` — unit tests for `open-pr.sh` decision logic with mocked `gh`/`git`.
- Create `.github/actions/update-deps/tests/local-e2e.sh` — maintainer-run integration against a real LocalAI + throwaway worktree.
- Create `.github/actions/update-deps/README.md` — adoption docs (caller snippet + GitHub App setup).
- Create `.github/actions/update-deps/examples/caller-workflow.yml` — canonical consumer workflow.
- Create `.github/workflows/lint-actions.yml` — CI: `shellcheck` + `actionlint` + run the shell unit tests.

---

### Task 1: `lib.sh` pure helpers + unit tests

**Files:**
- Create: `.github/actions/update-deps/scripts/lib.sh`
- Create: `.github/actions/update-deps/tests/assert.sh`
- Test: `.github/actions/update-deps/tests/lib_test.sh`

**Interfaces:**
- Produces (sourced by later scripts and tests):
  - `nib_base_url ENDPOINT` → prints ENDPOINT with a single `/v1` suffix (idempotent), stdout.
  - `lang_nib_task LANG` → prints the single-line nib task; returns non-zero on unknown LANG.
  - `lang_fallback_cmd LANG` → prints the deterministic update command; non-zero on unknown LANG.
  - `lang_verify_cmd LANG` → prints the verify command; non-zero on unknown LANG.
  - `lang_dep_paths LANG` → prints newline-separated manifest paths; non-zero on unknown LANG.
  - `has_dep_changes LANG` → exit 0 if any manifest path differs from `HEAD` in the CWD git repo, else 1.
  - `open_pr_number HEAD` → prints the number of an open PR whose head branch is HEAD (via `gh`), empty if none.
- `assert.sh` produces: `assert_eq EXPECTED ACTUAL MSG`, `assert_ok MSG CMD...`, `assert_fail MSG CMD...`, `with_mock_path DIR CMD...` and a `_report` trap that exits non-zero if any assertion failed.

- [ ] **Step 1: Write the assertion helper**

Create `.github/actions/update-deps/tests/assert.sh`:

```bash
#!/usr/bin/env bash
# Minimal test helpers. Source this from *_test.sh files.
set -uo pipefail

_FAILS=0

assert_eq() { # EXPECTED ACTUAL MSG
  if [ "$1" = "$2" ]; then
    printf 'ok   - %s\n' "$3"
  else
    printf 'FAIL - %s\n       expected: %q\n       actual:   %q\n' "$3" "$1" "$2"
    _FAILS=$((_FAILS + 1))
  fi
}

assert_ok() { # MSG CMD...
  local msg="$1"; shift
  if "$@" >/dev/null 2>&1; then printf 'ok   - %s\n' "$msg"
  else printf 'FAIL - %s (command failed: %s)\n' "$msg" "$*"; _FAILS=$((_FAILS + 1)); fi
}

assert_fail() { # MSG CMD...
  local msg="$1"; shift
  if "$@" >/dev/null 2>&1; then printf 'FAIL - %s (command unexpectedly succeeded)\n' "$msg"; _FAILS=$((_FAILS + 1))
  else printf 'ok   - %s\n' "$msg"; fi
}

# Run CMD... with DIR prepended to PATH (for mock executables).
with_mock_path() { # DIR CMD...
  local dir="$1"; shift
  PATH="$dir:$PATH" "$@"
}

_report() { if [ "$_FAILS" -ne 0 ]; then printf '\n%d assertion(s) failed\n' "$_FAILS"; exit 1; fi; printf '\nall assertions passed\n'; }
trap _report EXIT
```

- [ ] **Step 2: Write the failing test for `lib.sh`**

Create `.github/actions/update-deps/tests/lib_test.sh`:

```bash
#!/usr/bin/env bash
set -uo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tests/assert.sh
. "$HERE/assert.sh"
# shellcheck source=scripts/lib.sh
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
TMP="$(mktemp -d)"; trap 'rm -rf "$TMP"' EXIT
(
  cd "$TMP"
  git init -q && git config user.email t@t && git config user.name t
  printf 'module x\n' > go.mod; printf '' > go.sum
  git add -A && git commit -qm init
)
assert_fail "no change -> has_dep_changes false" bash -c ". '$HERE/../scripts/lib.sh'; cd '$TMP'; has_dep_changes go"
printf 'module x\n// bump\n' > "$TMP/go.mod"
assert_ok   "modified go.mod -> has_dep_changes true" bash -c ". '$HERE/../scripts/lib.sh'; cd '$TMP'; has_dep_changes go"
```

- [ ] **Step 3: Run the test to verify it fails**

Run: `bash .github/actions/update-deps/tests/lib_test.sh`
Expected: FAIL — `lib.sh` does not exist yet (source error / functions not found).

- [ ] **Step 4: Implement `lib.sh`**

Create `.github/actions/update-deps/scripts/lib.sh`:

```bash
#!/usr/bin/env bash
# Pure helpers for the update-deps action. Sourcing this file has no side
# effects; each function writes only to stdout / its exit status.

# nib_base_url ENDPOINT -> ENDPOINT normalized to carry exactly one /v1 suffix.
nib_base_url() {
  local base="${1%/}"
  case "$base" in
    "" | */v1) printf '%s' "$base" ;;
    *)         printf '%s/v1' "$base" ;;
  esac
}

# lang_nib_task LANG -> single-line task string for nib (unknown LANG -> rc 1).
lang_nib_task() {
  case "$1" in
    go) printf '%s' 'Update all Go dependencies in this repository to their latest versions by running "go get -u ./..." and then "go mod tidy". After that, run "go build ./..." to confirm the project still compiles, and fix any compilation errors caused by the updates. Report a short summary of what changed when done.' ;;
    *) return 1 ;;
  esac
}

# lang_fallback_cmd LANG -> deterministic dependency-update command.
lang_fallback_cmd() {
  case "$1" in
    go) printf '%s' 'go get -u ./... && go mod tidy' ;;
    *) return 1 ;;
  esac
}

# lang_verify_cmd LANG -> verify command (build + vet; never tests).
lang_verify_cmd() {
  case "$1" in
    go) printf '%s' 'go build ./... && go vet ./...' ;;
    *) return 1 ;;
  esac
}

# lang_dep_paths LANG -> newline-separated manifest paths whose change signals an update.
lang_dep_paths() {
  case "$1" in
    go) printf '%s\n' go.mod go.sum ;;
    *) return 1 ;;
  esac
}

# has_dep_changes LANG -> exit 0 if any manifest path differs from HEAD in CWD.
has_dep_changes() {
  local lang="$1" path rc=1
  while IFS= read -r path; do
    if ! git diff --quiet -- "$path" 2>/dev/null; then rc=0; fi
  done < <(lang_dep_paths "$lang")
  return "$rc"
}

# open_pr_number HEAD -> number of an open PR whose head branch is HEAD (empty if none).
open_pr_number() {
  gh pr list --head "$1" --state open --json number --jq '.[0].number // empty' 2>/dev/null
}
```

- [ ] **Step 5: Run the test to verify it passes**

Run: `bash .github/actions/update-deps/tests/lib_test.sh`
Expected: PASS — `all assertions passed`.

- [ ] **Step 6: Shellcheck the new files**

Run: `shellcheck .github/actions/update-deps/scripts/lib.sh .github/actions/update-deps/tests/assert.sh .github/actions/update-deps/tests/lib_test.sh`
Expected: no output (clean).

- [ ] **Step 7: Commit**

```bash
git add .github/actions/update-deps/scripts/lib.sh .github/actions/update-deps/tests/assert.sh .github/actions/update-deps/tests/lib_test.sh
git commit -m "feat(update-deps): add pure lib.sh helpers + unit tests

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 2: `run-update.sh` — nib run, fallback, verify

**Files:**
- Create: `.github/actions/update-deps/scripts/run-update.sh`
- Create: `.github/actions/update-deps/tests/local-e2e.sh`

**Interfaces:**
- Consumes (from Task 1): `nib_base_url`, `lang_nib_task`, `lang_fallback_cmd`, `lang_verify_cmd`, `has_dep_changes`.
- Reads env: `LANGUAGE`, `MODEL`, `LOCALAI_URL`, `NIB_AVAILABLE` (`1`/`0`), `GITHUB_OUTPUT` (optional).
- Produces: on success writes `changed=<true|false>` to `$GITHUB_OUTPUT` (if set) and stdout; exit 0 when the tree is buildable (or unchanged), exit 1 when verify still fails after the repair retry.

- [ ] **Step 1: Write `run-update.sh`**

Create `.github/actions/update-deps/scripts/run-update.sh`:

```bash
#!/usr/bin/env bash
# Update dependencies in the CWD repo: nib primary path, deterministic
# fallback, then a build+vet gate with one nib repair retry.
set -uo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib.sh
. "$HERE/lib.sh"

LANGUAGE="${LANGUAGE:-go}"
MODEL="${MODEL:?MODEL is required}"
LOCALAI_URL="${LOCALAI_URL:-http://localhost:8080}"
NIB_AVAILABLE="${NIB_AVAILABLE:-1}"

task="$(lang_nib_task "$LANGUAGE")"   || { echo "unsupported language: $LANGUAGE" >&2; exit 2; }
fallback="$(lang_fallback_cmd "$LANGUAGE")"
verify="$(lang_verify_cmd "$LANGUAGE")"

# Run nib once with the given one-line task (matches internal/remediate/nib_agent.go).
run_nib() { # TASK
  MODEL="$MODEL" BASE_URL="$(nib_base_url "$LOCALAI_URL")" API_KEY="sk-localai" \
    bash -c 'printf "%s\n" "$1" | nib --cli --yolo' _ "$1"
  # nib exits non-zero on the stdin EOF after the turn; that is expected, so we
  # never trust its exit code — the verify gate below is the source of truth.
  return 0
}

echo "== primary path: nib =="
if [ "$NIB_AVAILABLE" = "1" ]; then
  run_nib "$task"
else
  echo "nib/LocalAI unavailable — skipping primary path"
fi

# If nib produced no manifest change, fall back to the deterministic update so a
# PR still opens when the model is down or was a no-op.
if ! has_dep_changes "$LANGUAGE"; then
  echo "== no change from nib — deterministic fallback: $fallback =="
  eval "$fallback" || true
fi

echo "== verify: $verify =="
if eval "$verify"; then
  echo "verify OK"
else
  echo "== verify failed — one nib repair retry =="
  if [ "$NIB_AVAILABLE" = "1" ]; then
    run_nib "Fix the build after the dependency update. Run \"$verify\" and resolve every compilation or vet error it reports. Do not change application logic beyond what is needed to compile."
  fi
  if eval "$verify"; then
    echo "verify OK after repair"
  else
    echo "verify STILL FAILING after repair — refusing to open a PR" >&2
    exit 1
  fi
fi

if has_dep_changes "$LANGUAGE"; then changed=true; else changed=false; fi
echo "changed=$changed"
[ -n "${GITHUB_OUTPUT:-}" ] && echo "changed=$changed" >> "$GITHUB_OUTPUT"
exit 0
```

- [ ] **Step 2: Shellcheck it**

Run: `shellcheck .github/actions/update-deps/scripts/run-update.sh`
Expected: no output (clean).

- [ ] **Step 3: Write the maintainer-run local e2e harness**

Create `.github/actions/update-deps/tests/local-e2e.sh`:

```bash
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
```

- [ ] **Step 4: Run the local e2e against kairos-installer (maintainer, LocalAI up)**

Run:
```bash
LOCALAI_URL=http://localhost:8080 MODEL=gemma-4-e2b-it-qat-q4_0 \
  bash .github/actions/update-deps/tests/local-e2e.sh ~/_git/kairos-installer
```
Expected: `baseline build OK` … `PASS: updated and builds`.
(If no LocalAI is available in the execution environment, mark this step done by noting the earlier manual validation on edgevpn + kairos-installer recorded in the spec, and rely on CI/static gates.)

- [ ] **Step 5: Shellcheck the e2e**

Run: `shellcheck .github/actions/update-deps/tests/local-e2e.sh`
Expected: no output (clean).

- [ ] **Step 6: Commit**

```bash
git add .github/actions/update-deps/scripts/run-update.sh .github/actions/update-deps/tests/local-e2e.sh
git commit -m "feat(update-deps): add nib run + fallback + verify script

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 3: `start-localai.sh` — download + start, skip if already up

**Files:**
- Create: `.github/actions/update-deps/scripts/start-localai.sh`
- Test: extend `.github/actions/update-deps/tests/lib_test.sh` with the skip-decision function.

**Interfaces:**
- Consumes: nothing from other tasks (self-contained, uses `curl`/`yq`).
- Adds to `lib.sh`: `localai_answers URL` → exit 0 if `URL/readyz` returns HTTP 200 (used to decide skip).
- Reads env: `LOCALAI_URL`, `MODEL`, `LOCALAI_VERSION`, `STARTUP_TIMEOUT`, `GH_TOKEN`, `BIN_DIR`, `MODELS_PATH`. Produces: a running `local-ai` (or a clear failure) and prints readiness.

- [ ] **Step 1: Write the failing test for `localai_answers`**

Append to `.github/actions/update-deps/tests/lib_test.sh` (before nothing special; order-independent):

```bash
# localai_answers: mock curl to control the readyz HTTP code.
MOCK="$(mktemp -d)"
cat > "$MOCK/curl" <<'EOF'
#!/usr/bin/env bash
# emit the code requested via the sentinel env MOCK_CODE for -w '%{http_code}'
printf '%s' "${MOCK_CODE:-000}"
EOF
chmod +x "$MOCK/curl"
assert_ok   "localai_answers true on 200"  bash -c "MOCK_CODE=200 . '$HERE/../scripts/lib.sh'; PATH='$MOCK:\$PATH' localai_answers http://x"
assert_fail "localai_answers false on 000" bash -c "MOCK_CODE=000 . '$HERE/../scripts/lib.sh'; PATH='$MOCK:\$PATH' localai_answers http://x"
```

- [ ] **Step 2: Run to verify it fails**

Run: `bash .github/actions/update-deps/tests/lib_test.sh`
Expected: FAIL — `localai_answers` not defined.

- [ ] **Step 3: Add `localai_answers` to `lib.sh`**

Append to `.github/actions/update-deps/scripts/lib.sh`:

```bash
# localai_answers URL -> exit 0 if URL/readyz returns HTTP 200.
localai_answers() {
  local code
  code="$(curl -s -m 5 -o /dev/null -w '%{http_code}' "${1%/}/readyz" 2>/dev/null)"
  [ "$code" = "200" ]
}
```

- [ ] **Step 4: Run to verify it passes**

Run: `bash .github/actions/update-deps/tests/lib_test.sh`
Expected: PASS — `all assertions passed`.

- [ ] **Step 5: Write `start-localai.sh`** (adapted verbatim from `.github/workflows/security-dashboard.yaml` lines 85–131, parameterized and with a skip path)

Create `.github/actions/update-deps/scripts/start-localai.sh`:

```bash
#!/usr/bin/env bash
# Start LocalAI best-effort for driving nib. If LOCALAI_URL already answers, do
# nothing (the caller provided a server). Otherwise download the release binary
# and wait until a real chat completion succeeds.
set -uo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib.sh
. "$HERE/lib.sh"

LOCALAI_URL="${LOCALAI_URL:-http://localhost:8080}"
MODEL="${MODEL:?MODEL is required}"
LOCALAI_VERSION="${LOCALAI_VERSION:-latest}"
STARTUP_TIMEOUT="${STARTUP_TIMEOUT:-1200}"   # seconds
BIN_DIR="${BIN_DIR:-$PWD/bin}"
MODELS_PATH="${MODELS_PATH:-$PWD/models}"

if localai_answers "$LOCALAI_URL"; then
  echo "LocalAI already answering at $LOCALAI_URL — skipping startup"
  echo "nib_available=1"; [ -n "${GITHUB_OUTPUT:-}" ] && echo "nib_available=1" >> "$GITHUB_OUTPUT"
  exit 0
fi

mkdir -p "$BIN_DIR" "$MODELS_PATH"
curl -sfL https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -o "$BIN_DIR/yq" && chmod +x "$BIN_DIR/yq"
ver="$LOCALAI_VERSION"
if [ "$ver" = "latest" ]; then
  ver="$(curl -sfL ${GH_TOKEN:+-H "Authorization: Bearer $GH_TOKEN"} https://api.github.com/repos/mudler/LocalAI/releases/latest | "$BIN_DIR/yq" -p=json '.tag_name')"
fi
echo "LocalAI version: $ver"
url="https://github.com/mudler/LocalAI/releases/download/${ver}/local-ai-${ver}-linux-amd64"
code=$(curl -sL -w '%{http_code}' "$url" -o "$BIN_DIR/local-ai")
echo "download HTTP status: $code"
chmod +x "$BIN_DIR/local-ai" 2>/dev/null
"$BIN_DIR/local-ai" run "$MODEL" --address ":8080" --models-path "$MODELS_PATH" > localai.log 2>&1 &
LAI_PID=$!

deadline=$(( SECONDS + STARTUP_TIMEOUT ))
ready=0
while [ "$SECONDS" -lt "$deadline" ]; do
  if ! kill -0 "$LAI_PID" 2>/dev/null; then echo "local-ai exited early:"; tail -n 100 localai.log; break; fi
  resp=$(curl -s -m 180 "$LOCALAI_URL/v1/chat/completions" -H 'Content-Type: application/json' \
    -d "{\"model\":\"$MODEL\",\"messages\":[{\"role\":\"user\",\"content\":\"ping\"}],\"max_tokens\":1}")
  if printf '%s' "$resp" | grep -q '"choices"'; then ready=1; break; fi
  echo "model not ready yet (~$(( (deadline - SECONDS) / 60 ))m left)"
  sleep 15
done

if [ "$ready" = 1 ]; then
  echo "LocalAI model '$MODEL' is loaded."
  echo "nib_available=1"; [ -n "${GITHUB_OUTPUT:-}" ] && echo "nib_available=1" >> "$GITHUB_OUTPUT"
else
  echo "LocalAI not ready — nib primary path will be skipped, fallback will run."
  echo "nib_available=0"; [ -n "${GITHUB_OUTPUT:-}" ] && echo "nib_available=0" >> "$GITHUB_OUTPUT"
fi
exit 0
```

- [ ] **Step 6: Shellcheck both files**

Run: `shellcheck .github/actions/update-deps/scripts/start-localai.sh .github/actions/update-deps/scripts/lib.sh`
Expected: no output (clean). If shellcheck flags the `${GH_TOKEN:+-H "..."}` word-splitting, silence it with a targeted `# shellcheck disable=SC2086` line above that `curl` (intentional).

- [ ] **Step 7: Commit**

```bash
git add .github/actions/update-deps/scripts/start-localai.sh .github/actions/update-deps/scripts/lib.sh .github/actions/update-deps/tests/lib_test.sh
git commit -m "feat(update-deps): add LocalAI start/skip script

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 4: `open-pr.sh` — commit, dedupe, push, PR

**Files:**
- Create: `.github/actions/update-deps/scripts/open-pr.sh`
- Test: `.github/actions/update-deps/tests/open_pr_test.sh`

**Interfaces:**
- Consumes (from Task 1): `open_pr_number`.
- Reads env: `TOKEN`, `BRANCH`, `BASE`, `PR_TITLE`, `PR_LABELS`, `DRY_RUN` (`true`/`false`), `BOT_NAME`, `BOT_EMAIL`.
- Behavior: stages the manifest paths, commits on `BRANCH`, force-pushes, then reuses an existing open PR (via `open_pr_number`) or creates one. In dry-run, prints the intended git/gh commands and makes no network calls.

- [ ] **Step 1: Write the failing test with mocked `gh`/`git`**

Create `.github/actions/update-deps/tests/open_pr_test.sh`:

```bash
#!/usr/bin/env bash
set -uo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=tests/assert.sh
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

# Live, no existing PR -> pushes and creates.
: > "$CALLS"
CALLS_FILE="$CALLS" EXISTING="" DRY_RUN=false BRANCH=chore/update-deps BASE=main PR_TITLE=t PR_LABELS=deps TOKEN=x \
  PATH="$MOCK:$PATH" bash "$SCRIPT" >/dev/null 2>&1
assert_eq "1" "$(grep -c 'git push'    "$CALLS")" "live push happens"
assert_eq "1" "$(grep -c 'gh pr create' "$CALLS")" "live create happens when no PR exists"

# Live, existing PR #42 -> pushes but does NOT create a duplicate.
: > "$CALLS"
CALLS_FILE="$CALLS" EXISTING="42" DRY_RUN=false BRANCH=chore/update-deps BASE=main PR_TITLE=t PR_LABELS=deps TOKEN=x \
  PATH="$MOCK:$PATH" bash "$SCRIPT" >/dev/null 2>&1
assert_eq "1" "$(grep -c 'git push'     "$CALLS")" "live push happens with existing PR"
assert_eq "0" "$(grep -c 'gh pr create' "$CALLS")" "no duplicate PR when one is open"
```

- [ ] **Step 2: Run to verify it fails**

Run: `bash .github/actions/update-deps/tests/open_pr_test.sh`
Expected: FAIL — `open-pr.sh` does not exist.

- [ ] **Step 3: Implement `open-pr.sh`**

Create `.github/actions/update-deps/scripts/open-pr.sh`:

```bash
#!/usr/bin/env bash
# Commit the dependency update, push the branch, and open or update the PR.
# In DRY_RUN mode, print the intended git/gh commands and make no changes.
set -uo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib.sh
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
```

- [ ] **Step 4: Run to verify it passes**

Run: `bash .github/actions/update-deps/tests/open_pr_test.sh`
Expected: PASS — `all assertions passed`.

- [ ] **Step 5: Shellcheck**

Run: `shellcheck .github/actions/update-deps/scripts/open-pr.sh .github/actions/update-deps/tests/open_pr_test.sh`
Expected: no output (clean).

- [ ] **Step 6: Commit**

```bash
git add .github/actions/update-deps/scripts/open-pr.sh .github/actions/update-deps/tests/open_pr_test.sh
git commit -m "feat(update-deps): add commit/dedupe/push/PR script

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 5: `action.yml` — wire inputs to steps

**Files:**
- Create: `.github/actions/update-deps/action.yml`

**Interfaces:**
- Consumes: all three scripts + `nib_available`/`changed` outputs.
- Produces: the composite action consumed by caller workflows (Task 6).

- [ ] **Step 1: Write `action.yml`**

Create `.github/actions/update-deps/action.yml`:

```yaml
name: Update dependencies
description: Have nib bump a repo's dependencies to latest and open a CI-triggering PR.
inputs:
  token:            { description: "CI-triggering token (GitHub App token recommended; PAT fallback).", required: true }
  language:         { description: "Ecosystem. Only 'go' is implemented.", required: false, default: "go" }
  model:            { description: "LocalAI model driving nib.", required: false, default: "gemma-4-e2b-it" }
  localai-url:      { description: "If it already answers, startup is skipped.", required: false, default: "http://localhost:8080" }
  localai-version:  { description: "LocalAI release to download.", required: false, default: "latest" }
  nib-version:      { description: "nib version to install.", required: false, default: "latest" }
  go-version:       { description: "Go toolchain.", required: false, default: "stable" }
  startup-timeout:  { description: "Max seconds to wait for model load.", required: false, default: "1200" }
  branch:           { description: "PR head branch.", required: false, default: "chore/update-deps" }
  base:             { description: "PR target branch (default: repo default).", required: false, default: "" }
  pr-title:         { description: "PR title / commit subject.", required: false, default: "chore(deps): update dependencies" }
  pr-labels:        { description: "Comma-separated PR labels.", required: false, default: "dependencies" }
  dry-run:          { description: "Do everything except push/open PR.", required: false, default: "false" }
runs:
  using: composite
  steps:
    - uses: actions/setup-go@v5
      with: { go-version: "${{ inputs.go-version }}" }
    - name: Install nib
      shell: bash
      run: go install github.com/mudler/nib@${{ inputs.nib-version }}
    - name: Start LocalAI
      id: localai
      shell: bash
      env:
        LOCALAI_URL: ${{ inputs.localai-url }}
        MODEL: ${{ inputs.model }}
        LOCALAI_VERSION: ${{ inputs.localai-version }}
        STARTUP_TIMEOUT: ${{ inputs.startup-timeout }}
        GH_TOKEN: ${{ inputs.token }}
      run: bash "${{ github.action_path }}/scripts/start-localai.sh"
    - name: Update dependencies
      id: update
      shell: bash
      env:
        LANGUAGE: ${{ inputs.language }}
        MODEL: ${{ inputs.model }}
        LOCALAI_URL: ${{ inputs.localai-url }}
        NIB_AVAILABLE: ${{ steps.localai.outputs.nib_available }}
      run: bash "${{ github.action_path }}/scripts/run-update.sh"
    - name: Open pull request
      if: ${{ steps.update.outputs.changed == 'true' }}
      shell: bash
      env:
        TOKEN: ${{ inputs.token }}
        BRANCH: ${{ inputs.branch }}
        BASE: ${{ inputs.base }}
        PR_TITLE: ${{ inputs.pr-title }}
        PR_LABELS: ${{ inputs.pr-labels }}
        DRY_RUN: ${{ inputs.dry-run }}
      run: bash "${{ github.action_path }}/scripts/open-pr.sh"
```

- [ ] **Step 2: Validate with actionlint**

Run: `actionlint .github/actions/update-deps/action.yml` (install first if missing: `go install github.com/rhysd/actionlint/cmd/actionlint@latest`)
Expected: no output (clean).

- [ ] **Step 3: Commit**

```bash
git add .github/actions/update-deps/action.yml
git commit -m "feat(update-deps): wire composite action

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 6: Docs, example caller, and adoption for kairos-installer

**Files:**
- Create: `.github/actions/update-deps/README.md`
- Create: `.github/actions/update-deps/examples/caller-workflow.yml`

**Interfaces:**
- Consumes: the finished action (Task 5). Produces: the canonical consumer workflow that `kairos-io/kairos-installer` (and others) copy in.

- [ ] **Step 1: Write the example caller workflow**

Create `.github/actions/update-deps/examples/caller-workflow.yml`:

```yaml
# Copy to <target-repo>/.github/workflows/update-deps.yml
name: Update dependencies
on:
  schedule: [{ cron: "0 5 * * 1" }]   # weekly, Monday 05:00 UTC
  workflow_dispatch:
    inputs:
      dry-run: { description: "Preview without opening a PR", type: boolean, default: false }
concurrency:
  group: update-deps
  cancel-in-progress: false
jobs:
  update:
    runs-on: oracle-vm-16cpu-64gb-x86-64   # needs RAM for LocalAI + model
    steps:
      - uses: actions/checkout@v4
        with: { fetch-depth: 0 }
      - uses: actions/create-github-app-token@v2
        id: app-token
        with:
          app-id: ${{ secrets.DEPS_BOT_APP_ID }}
          private-key: ${{ secrets.DEPS_BOT_APP_KEY }}
      - uses: kairos-io/security/.github/actions/update-deps@main
        with:
          language: go
          token: ${{ steps.app-token.outputs.token }}
          dry-run: ${{ github.event.inputs.dry-run || 'false' }}
```

- [ ] **Step 2: Write the README**

Create `.github/actions/update-deps/README.md`:

```markdown
# update-deps action

Have `nib` (driven by a self-hosted LocalAI) bump a repository's dependencies to
their latest versions, verify the build, and open a CI-triggering pull request.

## Usage

Copy [`examples/caller-workflow.yml`](examples/caller-workflow.yml) to
`.github/workflows/update-deps.yml` in the target repo. Minimal call:

    - uses: kairos-io/security/.github/actions/update-deps@main
      with:
        language: go
        token: ${{ steps.app-token.outputs.token }}

## Inputs

See [`action.yml`](action.yml). Key ones: `token` (required), `language`
(only `go` today), `model` (default `gemma-4-e2b-it`), `branch`, `base`,
`dry-run`.

## Token: use a GitHub App (not the built-in token)

The built-in `GITHUB_TOKEN` opens a PR but its checks never run (GitHub
suppresses workflow runs triggered by that token). Use a **GitHub App** token so
CI triggers, with no personal PAT:

1. Create an org GitHub App (`kairos-deps-bot`) with **Contents: write** and
   **Pull requests: write**.
2. Generate a private key; note the App ID.
3. Install the App on the target repos.
4. Store `DEPS_BOT_APP_ID` and `DEPS_BOT_APP_KEY` as org secrets.
5. Mint a token per run with `actions/create-github-app-token@v2` and pass it as
   `token` (see the example workflow).

## Behavior

- nib is the primary engine; if LocalAI can't load, a deterministic
  `go get -u ./... && go mod tidy` fallback runs so a PR still opens.
- Verify gate is `go build ./... && go vet ./...`. The repo's own CI runs tests
  on the PR.
- No PR is opened when there is no dependency change; the action fails (no PR)
  when the build can't be made to pass.
- An already-open PR on the same branch is force-updated instead of duplicated.
```

- [ ] **Step 3: Commit**

```bash
git add .github/actions/update-deps/README.md .github/actions/update-deps/examples/caller-workflow.yml
git commit -m "docs(update-deps): adoption README + example caller workflow

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

- [ ] **Step 4: Adopt in kairos-installer (separate repo — manual)**

In `~/_git/kairos-installer`, on a new branch, add `.github/workflows/update-deps.yml` with the contents of `examples/caller-workflow.yml`. Confirm the `kairos-deps-bot` App is installed on the `kairos-io` org and `DEPS_BOT_APP_ID`/`DEPS_BOT_APP_KEY` secrets exist. Open that as its own PR in kairos-installer (out of scope for this repo's branch). Record completion here once the caller workflow PR is opened.

---

### Task 7: CI lint gate for the action

**Files:**
- Create: `.github/workflows/lint-actions.yml`

**Interfaces:**
- Consumes: all scripts, `action.yml`, and the `*_test.sh` files. Produces: a CI job that keeps them green.

- [ ] **Step 1: Write the lint workflow**

Create `.github/workflows/lint-actions.yml`:

```yaml
name: Lint actions
on:
  pull_request:
    paths: [".github/actions/update-deps/**", ".github/workflows/lint-actions.yml"]
  workflow_dispatch:
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: shellcheck
        run: shellcheck .github/actions/update-deps/scripts/*.sh .github/actions/update-deps/tests/*.sh
      - name: actionlint
        run: |
          go install github.com/rhysd/actionlint/cmd/actionlint@latest
          "$(go env GOPATH)/bin/actionlint" .github/actions/update-deps/action.yml .github/actions/update-deps/examples/caller-workflow.yml
      - name: unit tests
        run: |
          bash .github/actions/update-deps/tests/lib_test.sh
          bash .github/actions/update-deps/tests/open_pr_test.sh
```

- [ ] **Step 2: Run the unit tests locally one more time (full suite)**

Run:
```bash
bash .github/actions/update-deps/tests/lib_test.sh && bash .github/actions/update-deps/tests/open_pr_test.sh
```
Expected: both print `all assertions passed`.

- [ ] **Step 3: actionlint the new workflow**

Run: `actionlint .github/workflows/lint-actions.yml`
Expected: no output (clean).

- [ ] **Step 4: Commit**

```bash
git add .github/workflows/lint-actions.yml
git commit -m "ci(update-deps): lint scripts + action + run unit tests

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Self-Review

**Spec coverage:**
- Composite action packaging → Task 5 + example (Task 6). ✅
- LocalAI start (best-effort, skip if provided) → Task 3. ✅
- nib primary path (exact invocation/env) → Task 2 (`run_nib`). ✅
- Deterministic fallback (decided: included) → Task 2. ✅
- Verify gate build+vet, one repair retry, fail-no-PR → Task 2. ✅
- No-diff → no PR → Task 2 (`changed`) gates the PR step in Task 5. ✅
- PR push + dedupe existing + dry-run → Task 4. ✅
- Token model / App-token docs → Task 6 README. ✅
- Inputs table → Task 5 `action.yml`. ✅
- Extensibility via `language` switch, go-only, unknown fails fast → Task 1 (`lang_*` return 1) + Task 2 (`exit 2`). ✅
- First adopter kairos-installer → Task 6 Step 4. ✅
- Static gates (shellcheck/actionlint) → every task + Task 7 CI. ✅

**Placeholder scan:** No TBD/TODO; every code step shows complete content. The only manual/environment-dependent step (Task 2 Step 4 local e2e) has an explicit fallback instruction. ✅

**Type consistency:** Function names are stable across tasks — `nib_base_url`, `lang_nib_task`, `lang_fallback_cmd`, `lang_verify_cmd`, `lang_dep_paths`, `has_dep_changes`, `open_pr_number`, `localai_answers`. Env contract between scripts (`nib_available` output → `NIB_AVAILABLE`, `changed` output → PR-step `if`) matches between Tasks 2/3/5. ✅
