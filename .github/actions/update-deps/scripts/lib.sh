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

# collapse_ws STRING -> STRING with every whitespace run (incl. newlines/tabs)
# squeezed to a single space and no leading/trailing space. nib's --cli reads
# one prompt per stdin line, so a multi-line custom prompt must be flattened.
collapse_ws() {
  printf '%s' "$1" | tr '[:space:]' ' ' | sed -e 's/  */ /g' -e 's/^ //' -e 's/ $//'
}

# resolve_nib_task LANG CUSTOM -> the task string handed to nib: the CUSTOM
# prompt (flattened) when non-empty, otherwise the built-in per-language task.
resolve_nib_task() {
  if [ -n "$2" ]; then collapse_ws "$2"; else lang_nib_task "$1"; fi
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

# summarize_go_mod_diff -> reads a `git diff` of go.mod on stdin and prints a
# grouped markdown summary (direct deps updated/added/removed with from->to
# versions, plus a count of indirect changes). Empty output if nothing parsed.
summarize_go_mod_diff() {
  local records bd ad rd ind
  records="$(awk '
    /^[-+]/ {
      sign = substr($0, 1, 1); line = substr($0, 2)
      sub(/^[ \t]+/, "", line); sub(/^require[ \t]+/, "", line)
      if (split(line, f, /[ \t]+/) < 2) next
      p = f[1]; v = f[2]
      if (p !~ /\./) next                       # skip non-module lines (go directive, etc.)
      ind = (line ~ /\/\/[ \t]*indirect/)
      if (sign == "-") { o[p] = v; oi[p] = ind } else { nw[p] = v; ni[p] = ind }
    }
    END {
      for (p in nw) {
        if (p in o) { if (o[p] != nw[p]) print (ni[p] ? "BI" : "BD") "\t" p "\t" o[p] " -> " nw[p] }
        else print (ni[p] ? "AI" : "AD") "\t" p "\t" nw[p]
      }
      for (p in o) if (!(p in nw)) print (oi[p] ? "RI" : "RD") "\t" p
    }
  ')"
  [ -z "$records" ] && return 0
  bd="$(printf '%s\n' "$records" | awk -F'\t' '$1=="BD"{printf "- `%s` %s\n",$2,$3}' | sort)"
  ad="$(printf '%s\n' "$records" | awk -F'\t' '$1=="AD"{printf "- `%s` %s\n",$2,$3}' | sort)"
  rd="$(printf '%s\n' "$records" | awk -F'\t' '$1=="RD"{printf "- `%s`\n",$2}'        | sort)"
  ind="$(printf '%s\n' "$records" | grep -cE '^(BI|AI|RI)' || true)"

  echo "## Dependency updates"
  echo
  [ -n "$bd" ] && { echo "**Updated:**"; echo "$bd"; echo; }
  [ -n "$ad" ] && { echo "**Added:**";   echo "$ad"; echo; }
  [ -n "$rd" ] && { echo "**Removed:**"; echo "$rd"; echo; }
  if [ "${ind:-0}" -eq 1 ]; then
    echo "_1 indirect dependency also changed._"
  elif [ "${ind:-0}" -gt 1 ]; then
    echo "_${ind} indirect dependencies also changed._"
  fi
  return 0
}

# has_dep_changes LANG -> exit 0 if any manifest path has unstaged changes in CWD.
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

# localai_answers URL -> exit 0 if URL/readyz returns HTTP 200.
localai_answers() {
  local code
  code="$(curl -s -m 5 -o /dev/null -w '%{http_code}' "${1%/}/readyz" 2>/dev/null)"
  [ "$code" = "200" ]
}
