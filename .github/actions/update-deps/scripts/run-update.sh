#!/usr/bin/env bash
# Update dependencies in the CWD repo: nib primary path, deterministic
# fallback, then a build+vet gate with one nib repair retry.
set -uo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib.sh disable=SC1091
. "$HERE/lib.sh"

LANGUAGE="${LANGUAGE:-go}"
MODEL="${MODEL:?MODEL is required}"
LOCALAI_URL="${LOCALAI_URL:-http://localhost:8080}"
NIB_AVAILABLE="${NIB_AVAILABLE:-1}"
NIB_PROMPT="${NIB_PROMPT:-}"       # optional: replaces the default per-language task
PR_BODY_FILE="${PR_BODY_FILE:-}"   # optional: where to write the generated PR body

# Validate the language via a helper that fails for unsupported values, then
# resolve the task: a custom NIB_PROMPT replaces the built-in task entirely.
verify="$(lang_verify_cmd "$LANGUAGE")" || { echo "unsupported language: $LANGUAGE" >&2; exit 2; }
fallback="$(lang_fallback_cmd "$LANGUAGE")"
task="$(resolve_nib_task "$LANGUAGE" "$NIB_PROMPT")"

# Run nib once with the given one-line task (matches internal/remediate/nib_agent.go).
run_nib() { # TASK
  MODEL="$MODEL" BASE_URL="$(nib_base_url "$LOCALAI_URL")" API_KEY="sk-localai" \
    bash -c 'printf "%s\n" "$1" | nib --cli --yolo' _ "$1"
  # nib exits non-zero on the stdin EOF after the turn; that is expected, so we
  # never trust its exit code — the verify gate below is the source of truth.
  return 0
}

# summarize_prose DIFF -> a short natural-language summary of the changes from
# the model (best-effort; prints nothing if the model is unavailable, python3 is
# missing, or the call fails). Kept bounded — only the manifest diff is sent.
summarize_prose() {
  [ "$NIB_AVAILABLE" = "1" ] || return 0
  command -v python3 >/dev/null 2>&1 || return 0
  local payload resp content
  payload="$(DIFF="$1" MODEL="$MODEL" python3 - <<'PY' 2>/dev/null
import json, os
diff = os.environ["DIFF"][:6000]
print(json.dumps({"model": os.environ["MODEL"], "max_tokens": 200, "temperature": 0.2,
  "messages": [{"role": "user", "content":
    "Summarize these dependency changes in 2-3 plain sentences for a pull request "
    "description. Call out any major version jumps or removed modules. Do not list "
    "every version.\n\n" + diff}]}))
PY
)" || return 0
  resp="$(curl -s -m 60 "$LOCALAI_URL/v1/chat/completions" -H 'Content-Type: application/json' -d "$payload" 2>/dev/null)" || return 0
  content="$(printf '%s' "$resp" | python3 -c 'import sys,json; print(json.load(sys.stdin)["choices"][0]["message"]["content"].strip())' 2>/dev/null)" || return 0
  [ -n "$content" ] && printf '## Summary\n\n%s\n\n' "$content"
  return 0
}

echo "== primary path: nib =="
if [ "$NIB_AVAILABLE" = "1" ]; then
  run_nib "$task"
else
  echo "nib/LocalAI unavailable — skipping primary path"
fi

# If nib produced no manifest change, fall back to the deterministic update so a
# PR still opens when the model is down or was a no-op. A custom prompt overrides
# the default "update everything" intent, so we do NOT bump-everything behind the
# caller's back — skip the fallback when a custom prompt is set.
if ! has_dep_changes "$LANGUAGE"; then
  if [ -n "$NIB_PROMPT" ]; then
    echo "== no change from nib and a custom prompt is set — skipping deterministic fallback =="
  else
    echo "== no change from nib — deterministic fallback: $fallback =="
    eval "$fallback" || true
  fi
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

# Generate the PR body (best-effort model prose + a deterministic dependency
# table from the manifest diff) so open-pr.sh can use it as the PR description.
if [ "$changed" = "true" ] && [ -n "$PR_BODY_FILE" ]; then
  manifest="$(lang_dep_paths "$LANGUAGE" | head -1)"   # go.mod
  mdiff="$(git diff -- "$manifest")"
  mkdir -p "$(dirname "$PR_BODY_FILE")"
  { summarize_prose "$mdiff"; printf '%s\n' "$mdiff" | summarize_go_mod_diff; } > "$PR_BODY_FILE"
  echo "wrote PR body to $PR_BODY_FILE"
fi
exit 0
