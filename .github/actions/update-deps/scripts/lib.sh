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
