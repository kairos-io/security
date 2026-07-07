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
