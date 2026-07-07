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
