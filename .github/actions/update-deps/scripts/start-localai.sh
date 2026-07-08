#!/usr/bin/env bash
# Start LocalAI best-effort for driving nib. If LOCALAI_URL already answers, do
# nothing (the caller provided a server). Otherwise download the release binary
# and wait until a real chat completion succeeds.
set -uo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib.sh disable=SC1091
. "$HERE/lib.sh"

LOCALAI_URL="${LOCALAI_URL:-http://localhost:8080}"
MODEL="${MODEL:?MODEL is required}"
LOCALAI_VERSION="${LOCALAI_VERSION:-latest}"
STARTUP_TIMEOUT="${STARTUP_TIMEOUT:-1200}"   # seconds
CONTEXT_SIZE="${CONTEXT_SIZE:-}"             # optional: --context-size for the model
BIN_DIR="${BIN_DIR:-$PWD/bin}"
MODELS_PATH="${MODELS_PATH:-$PWD/models}"
LOCALAI_LOG="${LOCALAI_LOG:-localai.log}"
# LocalAI writes backends/ and data/ relative to its working directory; run it
# from here (a scratch dir outside the checkout) so it never litters the repo
# the action is updating. Defaults to the models dir's parent (runner temp).
LOCALAI_WORKDIR="${LOCALAI_WORKDIR:-$(dirname "$MODELS_PATH")}"

if localai_answers "$LOCALAI_URL"; then
  echo "LocalAI already answering at $LOCALAI_URL — skipping startup"
  echo "nib_available=1"; [ -n "${GITHUB_OUTPUT:-}" ] && echo "nib_available=1" >> "$GITHUB_OUTPUT"
  exit 0
fi

mkdir -p "$BIN_DIR" "$MODELS_PATH" "$LOCALAI_WORKDIR" "$(dirname "$LOCALAI_LOG")"
curl -sfL https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -o "$BIN_DIR/yq" && chmod +x "$BIN_DIR/yq"
ver="$LOCALAI_VERSION"
if [ "$ver" = "latest" ]; then
  # shellcheck disable=SC2086
  ver="$(curl -sfL ${GH_TOKEN:+-H "Authorization: Bearer $GH_TOKEN"} https://api.github.com/repos/mudler/LocalAI/releases/latest | "$BIN_DIR/yq" -p=json '.tag_name')"
fi
echo "LocalAI version: $ver"
url="https://github.com/mudler/LocalAI/releases/download/${ver}/local-ai-${ver}-linux-amd64"
code=$(curl -sL -w '%{http_code}' "$url" -o "$BIN_DIR/local-ai")
echo "download HTTP status: $code"
chmod +x "$BIN_DIR/local-ai" 2>/dev/null
run_args=(run "$MODEL" --address ":8080" --models-path "$MODELS_PATH")
[ -n "$CONTEXT_SIZE" ] && run_args+=(--context-size "$CONTEXT_SIZE")
( cd "$LOCALAI_WORKDIR" && exec "$BIN_DIR/local-ai" "${run_args[@]}" ) > "$LOCALAI_LOG" 2>&1 &
LAI_PID=$!

deadline=$(( SECONDS + STARTUP_TIMEOUT ))
ready=0
while [ "$SECONDS" -lt "$deadline" ]; do
  if ! kill -0 "$LAI_PID" 2>/dev/null; then echo "local-ai exited early:"; tail -n 100 "$LOCALAI_LOG"; break; fi
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
