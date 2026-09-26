# Kairos Security Dashboard

_Updated 2026-09-26._

🌐 **[Live dashboard](https://kairos-io.github.io/security/)** — the published board with clickable links.

## 📋 This run

- **Scanned:** 18 repos · ⚠️ 1 errored
- **Findings:** 0 (0 critical / 0 high / 0 medium / 0 low / 0 unknown)
- **Informational (not counted):** 52
- **CVE-related PRs:** 0
- **Remediation:** 0 open · 0 superseded · 0 merged · 0 need-human
- **Why:** No CVEs found, but 1 repo(s) could not be scanned — see collection errors.

## 🔥 Focus now

_Nothing flagged._

## 🌊 Waterfall fronts

_None._

## 📦 Per-repo findings

| Repo | Critical | High | Medium | Total | Status |
|---|---|---|---|---|---|
| [kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot) | 0 | 0 | 0 | 0 | ⚠️ errors |
| [kairos-io/cluster-api-provider-kairos](https://github.com/kairos-io/cluster-api-provider-kairos) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/entangle](https://github.com/kairos-io/entangle) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/entangle-proxy](https://github.com/kairos-io/entangle-proxy) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/go-nodepair](https://github.com/kairos-io/go-nodepair) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/go-ukify](https://github.com/kairos-io/go-ukify) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/hadron](https://github.com/kairos-io/hadron) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos](https://github.com/kairos-io/kairos) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-lab](https://github.com/kairos-io/kairos-lab) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/netboot](https://github.com/kairos-io/netboot) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/provider-kubernetes](https://github.com/kairos-io/provider-kubernetes) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/tpm-helpers](https://github.com/kairos-io/tpm-helpers) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mauromorales/xpasswd](https://github.com/mauromorales/xpasswd) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/edgevpn](https://github.com/mudler/edgevpn) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/entities](https://github.com/mudler/entities) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/go-pluggable](https://github.com/mudler/go-pluggable) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/yip](https://github.com/mudler/yip) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |

## Informational — not counted

These findings are separated from the counts above: CVEs we are already past, or components accepted as pinned risk.

| Package | Current | Fixed | Severity | CVE | Why |
|---|---|---|---|---|---|
| openssl-fips | 3.1.2 | 3.3.7 | critical | [CVE-2026-31789](https://osv.dev/vulnerability/ALPINE-CVE-2026-31789) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.6 | critical | [CVE-2024-5535](https://osv.dev/vulnerability/ALPINE-CVE-2024-5535) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-9076](https://osv.dev/vulnerability/ALPINE-CVE-2026-9076) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.5 | medium | [CVE-2025-9231](https://osv.dev/vulnerability/ALPINE-CVE-2025-9231) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| glib | 2.86.2 | 2.66.6 | high | [CVE-2021-27219](https://osv.dev/vulnerability/ALPINE-CVE-2021-27219) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2025-69418](https://osv.dev/vulnerability/ALPINE-CVE-2025-69418) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.1 | medium | [CVE-2025-4575](https://osv.dev/vulnerability/ALPINE-CVE-2025-4575) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.3 | medium | [CVE-2024-12797](https://osv.dev/vulnerability/ALPINE-CVE-2024-12797) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-45447](https://osv.dev/vulnerability/ALPINE-CVE-2026-45447) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | medium | [CVE-2026-45446](https://osv.dev/vulnerability/ALPINE-CVE-2026-45446) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.8 | high | [CVE-2025-9230](https://osv.dev/vulnerability/ALPINE-CVE-2025-9230) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.5 | medium | [CVE-2024-4603](https://osv.dev/vulnerability/ALPINE-CVE-2024-4603) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | critical | [CVE-2026-75803](https://osv.dev/vulnerability/ALPINE-CVE-2026-75803) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-7383](https://osv.dev/vulnerability/ALPINE-CVE-2026-7383) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2024-0727](https://osv.dev/vulnerability/ALPINE-CVE-2024-0727) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-5678](https://osv.dev/vulnerability/ALPINE-CVE-2023-5678) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.8 | medium | [CVE-2025-9232](https://osv.dev/vulnerability/ALPINE-CVE-2025-9232) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-69420](https://osv.dev/vulnerability/ALPINE-CVE-2025-69420) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-6129](https://osv.dev/vulnerability/ALPINE-CVE-2023-6129) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-6237](https://osv.dev/vulnerability/ALPINE-CVE-2023-6237) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | medium | [CVE-2026-42766](https://osv.dev/vulnerability/ALPINE-CVE-2026-42766) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | low | [CVE-2026-42770](https://osv.dev/vulnerability/ALPINE-CVE-2026-42770) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | high | [CVE-2023-5363](https://osv.dev/vulnerability/ALPINE-CVE-2023-5363) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.7 | high | [CVE-2024-6119](https://osv.dev/vulnerability/ALPINE-CVE-2024-6119) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| perl | 5.44.0 | 5.26.3 | unknown | [CVE-2018-18311](https://osv.dev/vulnerability/ALPINE-CVE-2018-18311) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.7 | critical | [CVE-2026-34182](https://osv.dev/vulnerability/ALPINE-CVE-2026-34182) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28387](https://osv.dev/vulnerability/ALPINE-CVE-2026-28387) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-15467](https://osv.dev/vulnerability/ALPINE-CVE-2025-15467) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | medium | [CVE-2026-63074](https://osv.dev/vulnerability/ALPINE-CVE-2026-63074) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-63072](https://osv.dev/vulnerability/ALPINE-CVE-2026-63072) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| libxml2 | 2.15.4 | 2.13.8 | high | [CVE-2025-32414](https://osv.dev/vulnerability/ALPINE-CVE-2025-32414) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-63076](https://osv.dev/vulnerability/ALPINE-CVE-2026-63076) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28388](https://osv.dev/vulnerability/ALPINE-CVE-2026-28388) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28389](https://osv.dev/vulnerability/ALPINE-CVE-2026-28389) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2025-68160](https://osv.dev/vulnerability/ALPINE-CVE-2025-68160) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.7 | medium | [CVE-2024-9143](https://osv.dev/vulnerability/ALPINE-CVE-2024-9143) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| libxml2 | 2.15.4 | 2.13.8 | high | [CVE-2025-32415](https://osv.dev/vulnerability/ALPINE-CVE-2025-32415) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28390](https://osv.dev/vulnerability/ALPINE-CVE-2026-28390) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2026-22796](https://osv.dev/vulnerability/ALPINE-CVE-2026-22796) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-34180](https://osv.dev/vulnerability/ALPINE-CVE-2026-34180) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-69421](https://osv.dev/vulnerability/ALPINE-CVE-2025-69421) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-69419](https://osv.dev/vulnerability/ALPINE-CVE-2025-69419) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-14456](https://osv.dev/vulnerability/ALPINE-CVE-2026-14456) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2024-2511](https://osv.dev/vulnerability/ALPINE-CVE-2024-2511) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2026-22795](https://osv.dev/vulnerability/ALPINE-CVE-2026-22795) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-31790](https://osv.dev/vulnerability/ALPINE-CVE-2026-31790) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-45445](https://osv.dev/vulnerability/ALPINE-CVE-2026-45445) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.8 | medium | [CVE-2024-13176](https://osv.dev/vulnerability/ALPINE-CVE-2024-13176) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-54874](https://osv.dev/vulnerability/ALPINE-CVE-2026-54874) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-42767](https://osv.dev/vulnerability/ALPINE-CVE-2026-42767) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.6 | high | [CVE-2024-4741](https://osv.dev/vulnerability/ALPINE-CVE-2024-4741) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| perl | 5.44.0 | 5.26.3 | unknown | [CVE-2018-18312](https://osv.dev/vulnerability/ALPINE-CVE-2018-18312) | already-fixed |

## ⚠️ 1 collection errors

- [kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot) / sourceCVE: govulncheck: exit status 1: govulncheck: loading packages: There are errors with the provided package patterns: -: # github.com/go-piv/piv-go/v2/piv # [pkg-config --cflags -- libpcsclite] Package libpcsclite was not found in the pkg-config search path. Perhaps you sho … (truncated)

## 📋 Open PRs

_None._

## 🤖 Bot PR ledger

_No bot PRs yet._

## 🔎 Bot-PR reviews

**[kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot)**

- [#674](https://github.com/kairos-io/AuroraBoot/pull/674) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - microsoft/TypeScript v6.0.3..2bd066d87f5bafd315be9f40889d0a60b9e58e0b (PR body): compare v6.0.3...2bd066d87f5bafd315be9f40889d0a60b9e58e0b failed/empty (no upstream diff)
    - microsoft/TypeScript v6.0.2..v6.0.3 (PR body): compare v6.0.2...v6.0.3 ✓ 40000 bytes
    - microsoft/TypeScript v5.9.3..v6.0.2 (PR body): compare v5.9.3...v6.0.2 failed/empty (no upstream diff)
    - context: 44246 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#813](https://github.com/kairos-io/AuroraBoot/pull/813) — ✅ **good** — The change is a routine dependency update to a minor version of the Go language image. This update does not introduce any new security vulnerabilities and is a standard maintenance practice. Therefore, it is safe to auto-approve.
  ↳ This PR updates the base Go Docker image tag from version 1.26 to 1.27 across all relevant Dockerfiles.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 3041 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#842](https://github.com/kairos-io/AuroraBoot/pull/842) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/containerd/containerd/v2 2.4.0→2.4.1: compare v2.4.0...v2.4.1 ✓ 40000 bytes
    - github.com/fxamacker/cbor/v2 2.9.1→2.9.4: compare v2.9.1...v2.9.4 ✓ 40000 bytes
    - context: 111556 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#843](https://github.com/kairos-io/AuroraBoot/pull/843) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/klauspost/compress 1.20.0→1.20.1: compare v1.20.0...v1.20.1 ✓ 40000 bytes
    - context: 66313 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#858](https://github.com/kairos-io/AuroraBoot/pull/858) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - k8s.io/api 0.37.0→0.37.1: compare v0.37.0...v0.37.1 ✓ 1313 bytes
    - k8s.io/apimachinery 0.37.0→0.37.1: compare v0.37.0...v0.37.1 ✓ 1582 bytes
    - k8s.io/client-go 0.37.0→0.37.1: compare v0.37.0...v0.37.1 ✓ 2874 bytes
    - context: 31318 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#861](https://github.com/kairos-io/AuroraBoot/pull/861) — ✅ **good** — This is a routine minor version update for a dependency. Minor version bumps are generally safe and necessary for receiving bug fixes and security patches. There are no obvious security risks indicated by the context.
  ↳ This PR updates the `LUET_VERSION` argument in the `Dockerfile` and `Dockerfile.riscv64` from `0.36.5` to `0.37.0` for the `quay.io/luet/base` package.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1610 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#863](https://github.com/kairos-io/AuroraBoot/pull/863) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/download-artifact v8.0.0..v8.0.1 (PR body): compare v8.0.0...v8.0.1 ✓ 40000 bytes
    - actions/download-artifact v8..v8.0.1 (PR body): compare v8...v8.0.1 failed/empty (no upstream diff)
    - actions/download-artifact v7.0.0..v8.0.0 (PR body): compare v7.0.0...v8.0.0 ✓ 40000 bytes
    - context: 85742 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#868](https://github.com/kairos-io/AuroraBoot/pull/868) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/containerd/containerd/v2 2.4.0→2.4.1: compare v2.4.0...v2.4.1 ✓ 40000 bytes
    - github.com/fxamacker/cbor/v2 2.9.1→2.9.4: compare v2.9.1...v2.9.4 ✓ 40000 bytes
    - context: 138379 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#870](https://github.com/kairos-io/AuroraBoot/pull/870) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/kairos-io/netboot 0.0.0-20260921125349-ba8ed34a440d→0.0.0-20260923193508-f2590943fb41: compare ba8ed34a440d...f2590943fb41 ✓ 1682 bytes
    - github.com/onsi/gomega 1.43.1→1.44.0: compare v1.43.1...v1.44.0 ✓ 31883 bytes
    - github.com/kairos-io/kairos/v4 4.3.1-0.20260921115150-170bfbb891c4→4.3.1-0.20260925140850-092013fed4d0: compare 170bfbb891c4...092013fed4d0 ✓ 40000 bytes
    - context: 104920 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#871](https://github.com/kairos-io/AuroraBoot/pull/871) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - lucide-icons/lucide 0.576.0..0.577.0 (PR body): compare 0.576.0...0.577.0 ✓ 40000 bytes
    - context: 125781 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#872](https://github.com/kairos-io/AuroraBoot/pull/872) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - azure/setup-helm v5.0.0..v5.0.1 (PR body): compare v5.0.0...v5.0.1 ✓ 40000 bytes
    - azure/setup-helm v4.3.1..v5.0.0 (PR body): compare v4.3.1...v5.0.0 ✓ 40000 bytes
    - context: 86218 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
**[kairos-io/entangle](https://github.com/kairos-io/entangle)**

- [#13](https://github.com/kairos-io/entangle/pull/13) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/emicklei/go-restful 2.9.5+incompatible→2.16.0+incompatible: compare v2.9.5+incompatible...v2.16.0+incompatible failed/empty (no upstream diff)
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - golang.org/x/net 0.55.0→0.56.0: compare v0.55.0...v0.56.0 ✓ 40000 bytes
    - context: 97666 bytes
**[kairos-io/entangle-proxy](https://github.com/kairos-io/entangle-proxy)**

- [#6](https://github.com/kairos-io/entangle-proxy/pull/6) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - sigs.k8s.io/controller-runtime 0.12.1→0.25.1: compare v0.12.1...v0.25.1 ✓ 40000 bytes
    - context: 98865 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#14](https://github.com/kairos-io/entangle-proxy/pull/14) — ✅ **good** — This pull request primarily updates several dependencies to newer versions, including core packages like `golang.org/x` and `google.golang.org/protobuf`. Updating dependencies is a crucial security practice to ensure that known vulnerabilities are patched. The changes appear to be dependency hygiene improvements and do not introduce any obvious security risks.
- [#18](https://github.com/kairos-io/entangle-proxy/pull/18) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/build-push-action v7.2.0..v7.3.0 (PR body): compare v7.2.0...v7.3.0 ✓ 40000 bytes
    - context: 82843 bytes
- [#20](https://github.com/kairos-io/entangle-proxy/pull/20) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - k8s.io/api 0.24.0→0.37.1: compare v0.24.0...v0.37.1 ✓ 40000 bytes
    - context: 129687 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#23](https://github.com/kairos-io/entangle-proxy/pull/23) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/checkout v7..v7.0.1 (PR body): compare v7...v7.0.1 failed/empty (no upstream diff)
    - actions/checkout v6.1.0..v7 (PR body): compare v6.1.0...v7 ✓ 40000 bytes
    - context: 63408 bytes
- [#25](https://github.com/kairos-io/entangle-proxy/pull/25) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/go-logr/logr 1.4.3→1.4.4: compare v1.4.3...v1.4.4 ✓ 40000 bytes
    - context: 44091 bytes
- [#26](https://github.com/kairos-io/entangle-proxy/pull/26) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/onsi/ginkgo/v2 2.32.0→2.33.0: compare v2.32.0...v2.33.0 ✓ 40000 bytes
    - onsi/ginkgo v2.32.2..v2.33.0 (PR body): compare v2.32.2...v2.33.0 ✓ 40000 bytes
    - context: 84296 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#27](https://github.com/kairos-io/entangle-proxy/pull/27) — ✅ **good** — This is a routine minor version update for the golang base image. Minor version bumps typically include bug fixes and security patches, making this change safe and necessary for maintaining a current build environment. There are no apparent security risks introduced by this update.
  ↳ This PR updates the base image for the Go build stage in the Dockerfile from golang:1.26 to golang:1.27.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1405 bytes
- [#29](https://github.com/kairos-io/entangle-proxy/pull/29) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/onsi/gomega 1.40.0→1.44.0: compare v1.40.0...v1.44.0 failed/empty (no upstream diff)
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - golang.org/x/net 0.55.0→0.56.0: compare v0.55.0...v0.56.0 ✓ 40000 bytes
    - context: 93241 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
**[kairos-io/go-nodepair](https://github.com/kairos-io/go-nodepair)**

- [#65](https://github.com/kairos-io/go-nodepair/pull/65) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - ipfs/go-log v2.9.1..v2.9.2 (PR body): compare v2.9.1...v2.9.2 ✓ 12110 bytes
    - ipfs/go-log v2.9.0..v2.9.1 (PR body): compare v2.9.0...v2.9.1 ✓ 2149 bytes
    - ipfs/go-log v2.8.2..v2.9.0 (PR body): compare v2.8.2...v2.9.0 ✓ 40000 bytes
    - context: 70615 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#66](https://github.com/kairos-io/go-nodepair/pull/66) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - context: 82644 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#69](https://github.com/kairos-io/go-nodepair/pull/69) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - google/osv-scanner-action v2.5.1..v2.6.0 (PR body): compare v2.5.1...v2.6.0 ✓ 9253 bytes
    - google/osv-scanner-action v2.5.0..v2.5.1 (PR body): compare v2.5.0...v2.5.1 ✓ 9140 bytes
    - google/osv-scanner-action v2.3.8..v2.5.0 (PR body): compare v2.3.8...v2.5.0 ✓ 12179 bytes
    - context: 36348 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#75](https://github.com/kairos-io/go-nodepair/pull/75) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/onsi/gomega 1.42.1→1.44.0: compare v1.42.1...v1.44.0 ✓ 40000 bytes
    - onsi/gomega v1.43.1..v1.44.0 (PR body): compare v1.43.1...v1.44.0 ✓ 31883 bytes
    - context: 80127 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
- [#81](https://github.com/kairos-io/go-nodepair/pull/81) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/onsi/ginkgo/v2 2.32.2→2.33.0: compare v2.32.2...v2.33.0 ✓ 40000 bytes
    - context: 43656 bytes
    - comment not posted: list comments: gh api: the `--slurp` option is not supported with `--jq` or `--template`

Usage:  gh api <endpoint> [flags]

Flags:
      --allow-escape-sequences   Allow printing terminal escape sequences
      --cache duration           Cache the response, e.g. "3600s", "60m", "1h"
  -F, --field key=value          Add a typed parameter in key=value format (use "@<path>" or "@-" to read value from file or stdin)
  -H, --header key:value         Add a HTTP request header in key:value format
      --hostname string          The GitHub hostname for the request (default "github.com")
  -i, --include                  Include HTTP response status line and headers in the output
      --input file               The file to use as body for the HTTP request (use "-" to read from standard input)
  -q, --jq string                Query to select values from the response using jq syntax
  -X, --method string            The HTTP method for the request (default "GET")
      --paginate                 Make additional HTTP requests to fetch all pages of results
  -p, --preview strings          Opt into GitHub API previews (names should omit '-preview')
  -f, --raw-field key=value      Add a string parameter in key=value format
      --silent                   Do not print the response body
      --slurp                    Use with "--paginate" to return an array of all pages of either JSON arrays or objects
  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
      --verbose                  Include full HTTP request and response in the output → retried next run
**[kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator)**

- [#172](https://github.com/kairos-io/kairos-operator/pull/172) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - k8s.io/api 0.36.3→0.37.0: compare v0.36.3...v0.37.0 ✓ 40000 bytes
    - k8s.io/apimachinery 0.36.3→0.37.0: compare v0.36.3...v0.37.0 ✓ 40000 bytes
    - context: 132062 bytes
**[kairos-io/netboot](https://github.com/kairos-io/netboot)**

- [#46](https://github.com/kairos-io/netboot/pull/46) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.53.0→0.57.0: compare v0.53.0...v0.57.0 ✓ 40000 bytes
    - golang.org/x/net 0.56.0→0.59.0: compare v0.56.0...v0.59.0 ✓ 40000 bytes
    - context: 83965 bytes
- [#47](https://github.com/kairos-io/netboot/pull/47) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - context: 83203 bytes
- [#48](https://github.com/kairos-io/netboot/pull/48) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.53.0→0.57.0: compare v0.53.0...v0.57.0 ✓ 40000 bytes
    - golang.org/x/net 0.56.0→0.58.0: compare v0.56.0...v0.58.0 ✓ 40000 bytes
    - context: 83976 bytes
- [#49](https://github.com/kairos-io/netboot/pull/49) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - google/osv-scanner-action v2.5.0..v2.5.1 (PR body): compare v2.5.0...v2.5.1 ✓ 9140 bytes
    - google/osv-scanner-action v2.3.8..v2.5.0 (PR body): compare v2.3.8...v2.5.0 ✓ 12179 bytes
    - context: 25806 bytes
**[kairos-io/tpm-helpers](https://github.com/kairos-io/tpm-helpers)**

- [#12](https://github.com/kairos-io/tpm-helpers/pull/12) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/gorilla/websocket 1.5.0→1.5.3: compare v1.5.0...v1.5.3 ✓ 7569 bytes
    - gorilla/websocket v1.5.1..v1.5.3 (PR body): compare v1.5.1...v1.5.3 ✓ 40000 bytes
    - context: 62623 bytes
**[mauromorales/xpasswd](https://github.com/mauromorales/xpasswd)**

- [#64](https://github.com/mauromorales/xpasswd/pull/64) — ✅ **good** — This is a standard dependency version bump to a newer release. The changelog indicates a new feature addition, and the diffs show the necessary code changes to support this new functionality. There are no obvious security risks or breaking changes indicated in the provided context.
  ↳ The PR updates the dependency `github.com/onsi/gomega` from version `v1.42.1` to `v1.43.0`. This update introduces a new feature: a gomock adaptor extension that allows Gomega matchers to be used as gomock argument matchers.
    - github.com/onsi/gomega 1.42.1→1.43.0: compare v1.42.1...v1.43.0 ✓ 3786 bytes
    - context: 6711 bytes
- [#65](https://github.com/mauromorales/xpasswd/pull/65) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang/go go1.27rc3..go1.27.1 (PR body): compare go1.27rc3...go1.27.1 ✓ 40000 bytes
    - context: 41548 bytes
- [#66](https://github.com/mauromorales/xpasswd/pull/66) — ⚠️ **needs_human_verification** — tool-call arguments were not valid JSON
    - github.com/onsi/ginkgo/v2 2.32.1→2.32.2: compare v2.32.1...v2.32.2 ✓ 18518 bytes
    - context: 21369 bytes
**[mudler/edgevpn](https://github.com/mudler/edgevpn)**

- [#804](https://github.com/mudler/edgevpn/pull/804) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - c-robinson/iplib v2.0.4..v2.0.5 (PR body): compare v2.0.4...v2.0.5 ✓ 6378 bytes
    - c-robinson/iplib v2.0.3..v2.0.4 (PR body): compare v2.0.3...v2.0.4 ✓ 3273 bytes
    - c-robinson/iplib v2.0.2..v2.0.3 (PR body): compare v2.0.2...v2.0.3 ✓ 9999 bytes
    - c-robinson/iplib v2.0.1..v2.0.2 (PR body): compare v2.0.1...v2.0.2 ✓ 15662 bytes
    - c-robinson/iplib v2.0.0..v2.0.1 (PR body): compare v2.0.0...v2.0.1 ✓ 1844 bytes
    - context: 44506 bytes
- [#905](https://github.com/mudler/edgevpn/pull/905) — ✅ **good** — This is a routine dependency bump for a tool used in the CI/CD workflow. The changelog indicates that version 2.4.0 includes various maintenance updates and fixes, suggesting this is a safe and necessary update. There are no immediate security red flags indicated by the context.
  ↳ This pull request updates the version of the `dependabot/fetch-metadata` dependency from 2.3.0 to 2.4.0. This upgrade incorporates various fixes, updates to actions, and improvements to the dependency fetching mechanism.
    - dependabot/fetch-metadata v2..v2.4.0 (PR body): compare v2...v2.4.0 failed/empty (no upstream diff)
    - dependabot/fetch-metadata v2.3.0..v2.4.0 (PR body): compare v2.3.0...v2.4.0 ✓ 40000 bytes
    - context: 49729 bytes
- [#923](https://github.com/mudler/edgevpn/pull/923) — ⚠️ **needs_human_verification** — review endpoint unreachable: Post "http://localhost:8080/v1/chat/completions": context deadline exceeded
    - github.com/miekg/dns 1.1.66→1.1.68: compare v1.1.66...v1.1.68 ✓ 40000 bytes
    - miekg/dns v1.1.64..v1.1.68 (PR body): compare v1.1.64...v1.1.68 ✓ 40000 bytes
    - context: 86107 bytes
- [#927](https://github.com/mudler/edgevpn/pull/927) — ✅ **good** — This is a standard dependency upgrade to a newer major version, which is generally a positive security and maintenance practice. The changes include necessary updates to workflows to use the new action version and Node.js version, as well as internal code refactoring to align with the v5 API. No security regressions are apparent.
  ↳ This pull request bumps the `actions/checkout` dependency from version 4 to 5.0.0 and updates related configurations across workflows and source code. It also updates the Node.js version used in workflows to 24.x and refactors the URL helper logic for improved handling of GitHub Enterprise Cloud and other hostnames.
    - actions/checkout v4..v5.0.0 (PR body): compare v4...v5.0.0 ✓ 11870 bytes
    - actions/checkout v4..v4.3.0 (PR body): compare v4...v4.3.0 failed/empty (no upstream diff)
    - actions/checkout v4.2.1..v4.2.2 (PR body): compare v4.2.1...v4.2.2 ✓ 9872 bytes
    - actions/checkout v4.2.0..v4.2.1 (PR body): compare v4.2.0...v4.2.1 ✓ 3510 bytes
    - actions/checkout v4..v5 (PR body): compare v4...v5 ✓ 40000 bytes
    - context: 84131 bytes
- [#942](https://github.com/mudler/edgevpn/pull/942) — ✅ **good** — This is a routine dependency update to a newer minor version of a well-known testing library. The changes primarily involve version bumps and internal code refactoring, which are typical for dependency maintenance. Since this is a standard update and the changes appear to be focused on compatibility and minor fixes, it is safe to auto-approve.
  ↳ This PR bumps github.com/onsi/gomega to version 1.38.2 and updates several related dependencies, including golang.org/x/net, google.golang.org/protobuf, and gopkg.in/yaml.v3. It also includes internal refactoring in gstruct to improve handling of unexported fields and updates to internal error handling.
    - github.com/onsi/gomega 1.37.0→1.38.2: compare v1.37.0...v1.38.2 ✓ 34194 bytes
    - github.com/Masterminds/semver/v3 3.3.1→3.4.0: compare v3.3.1...v3.4.0 ✓ 40000 bytes
    - context: 89137 bytes
- [#943](https://github.com/mudler/edgevpn/pull/943) — ✅ **good** — This is a routine dependency update for a well-known action. The changes are confined to updating the version number, and the upstream changes detailed in the changelog appear to be standard maintenance and minor feature updates, posing no immediate security risk.
  ↳ This pull request updates the `codecov/codecov-action` dependency from version 5.5.0 to 5.5.1. This version bump incorporates several underlying dependency updates for related actions, such as `actions/checkout` and `github/codeql-action`.
    - codecov/codecov-action v5.5.0..v5.5.1 (PR body): compare v5.5.0...v5.5.1 ✓ 10680 bytes
    - context: 21031 bytes
- [#951](https://github.com/mudler/edgevpn/pull/951) — ✅ **good** — This is a standard dependency bump for a widely used GitHub Action. The changes primarily involve updating the version number and migrating usage patterns in workflows, which is typical for dependency maintenance. The noted breaking change regarding Node v24.x support is documented, making the update safe to proceed with for automated approval.
  ↳ This PR bumps the dependency `actions/download-artifact` from version 5 to 6. It updates the dependency version in the configuration, modifies usage in workflow files to use the new version, and updates internal code imports. The release notes indicate a breaking change related to Node v24.x support.
    - actions/download-artifact v5..v6.0.0 (PR body): compare v5...v6.0.0 ✓ 40000 bytes
    - actions/download-artifact v5..v6 (PR body): compare v5...v6 ✓ 40000 bytes
    - context: 88447 bytes
- [#961](https://github.com/mudler/edgevpn/pull/961) — ✅ **good** — The changes are a dependency bump to a newer minor version of a well-maintained library. The diffs show internal refactoring, modernization of logging, and the addition of new features (GossipSub v1.3 support and peer extensions). There are no apparent security regressions or breaking API changes that would warrant manual review.
  ↳ This PR bumps `go-libp2p-pubsub` to version 0.15.0, which includes internal refactoring for logging (migrating to `log/slog`), the addition of support for GossipSub protocol version 1.3, and the implementation of a new Peer Extensions mechanism for testing. These changes are primarily internal improvements and feature additions.
    - github.com/libp2p/go-libp2p-pubsub 0.14.2→0.15.0: compare v0.14.2...v0.15.0 ✓ 40000 bytes
    - libp2p/go-libp2p-pubsub v0.14.3..v0.15.0 (PR body): compare v0.14.3...v0.15.0 ✓ 40000 bytes
    - context: 116011 bytes
- [#1006](https://github.com/mudler/edgevpn/pull/1006) — ✅ **good** — The upgrade is to a newer minor version (4.15.1) which includes security enhancements, such as the new CSRF middleware features detailed in the release notes. There are no immediate red flags or known critical vulnerabilities associated with this specific version jump. Therefore, this change is safe to auto-approve.
  ↳ This pull request updates the dependency `github.com/labstack/echo/v4` from version 4.13.3 to 4.15.1. This upgrade incorporates several enhancements, including improved CSRF protection features and minor internal fixes related to time comparison logic.
    - github.com/labstack/echo/v4 4.13.3→4.15.1: compare v4.13.3...v4.15.1 ✓ 40000 bytes
    - github.com/mattn/go-colorable 0.1.13→0.1.14: compare v0.1.13...v0.1.14 ✓ 6350 bytes
    - golang.org/x/time 0.12.0→0.14.0: compare v0.12.0...v0.14.0 ✓ 606 bytes
    - context: 76092 bytes
- [#1056](https://github.com/mudler/edgevpn/pull/1056) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/libp2p/go-libp2p-pubsub 0.16.0→0.17.0: compare v0.16.0...v0.17.0 ✓ 40000 bytes
    - context: 45663 bytes
- [#1057](https://github.com/mudler/edgevpn/pull/1057) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - context: 84432 bytes
- [#1059](https://github.com/mudler/edgevpn/pull/1059) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/libp2p/go-libp2p 0.48.0→0.49.0: compare v0.48.0...v0.49.0 ✓ 40000 bytes
    - github.com/libp2p/go-libp2p-kad-dht 0.41.0→0.42.2: compare v0.41.0...v0.42.2 ✓ 40000 bytes
    - context: 111567 bytes
- [#1061](https://github.com/mudler/edgevpn/pull/1061) — ✅ **good** — The change is a routine dependency update involving a digest change for a known library. There are no apparent security risks or changes to the application logic, making it safe to auto-approve.
  ↳ This PR updates the dependency `github.com/mudler/go-libp2p-pubsub` by replacing its old digest (`205ded1`) with a newer one (`2a31b5e`). This is a routine maintenance update to pull in the latest version of the library.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 2511 bytes
- [#1067](https://github.com/mudler/edgevpn/pull/1067) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6..v6.5.0 (PR body): compare v6...v6.5.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v6.4.0 (PR body): compare v6...v6.4.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v6.3.0 (PR body): compare v6...v6.3.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v7 (PR body): compare v6...v7 ✓ 40000 bytes
    - context: 89775 bytes
- [#1068](https://github.com/mudler/edgevpn/pull/1068) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/libp2p/go-libp2p-kad-dht 0.41.0→0.42.1: compare v0.41.0...v0.42.1 ✓ 40000 bytes
    - github.com/ipfs/boxo 0.39.0→0.41.0: compare v0.39.0...v0.41.0 ✓ 40000 bytes
    - context: 189140 bytes
- [#1069](https://github.com/mudler/edgevpn/pull/1069) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - nodejs/node v20.20.1..v20.20.2 (PR body): compare v20.20.1...v20.20.2 ✓ 40000 bytes
    - context: 65436 bytes
- [#1070](https://github.com/mudler/edgevpn/pull/1070) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-node v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - actions/setup-node v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - context: 95909 bytes
- [#1076](https://github.com/mudler/edgevpn/pull/1076) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - vitejs/vite v8.2.1..v8.2.2 (PR body): compare v8.2.1...v8.2.2 ✓ 40000 bytes
    - vitejs/vite v8.2.0..v8.2.1 (PR body): compare v8.2.0...v8.2.1 ✓ 40000 bytes
    - context: 111633 bytes
- [#1077](https://github.com/mudler/edgevpn/pull/1077) — ✅ **good** — The PR updates a dependency to a newer version that includes a bug fix, which is a positive change. The accompanying code modifications in `src/pure.js` and the new test file appear to improve the robustness of event handling and test for act warnings, making this a safe and beneficial update.
  ↳ This PR updates the `@testing-library/react` dependency to version `16.3.3`, which includes a bug fix for `act()` re-entrant behavior. It also introduces code modifications in `src/pure.js` to improve event wrapper handling and adds a new test file to verify the fix.
    - testing-library/react-testing-library v16.3.2..v16.3.3 (PR body): compare v16.3.2...v16.3.3 ✓ 4007 bytes
    - context: 7242 bytes
- [#1079](https://github.com/mudler/edgevpn/pull/1079) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - react/react v19.2.8..1d34f91dfde6bba84d08b683aaba164c7194dacb (PR body): compare v19.2.8...1d34f91dfde6bba84d08b683aaba164c7194dacb ✓ 40000 bytes
    - react/react v19.2.8..v19.3.0 (PR body): compare v19.2.8...v19.3.0 ✓ 40000 bytes
    - context: 87315 bytes
- [#1081](https://github.com/mudler/edgevpn/pull/1081) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - vitest-dev/vitest v4.1.10..v4.1.11 (PR body): compare v4.1.10...v4.1.11 ✓ 40000 bytes
    - context: 143656 bytes
**[mudler/entities](https://github.com/mudler/entities)**

- [#10](https://github.com/mudler/entities/pull/10) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/text 0.3.2→0.3.8: compare v0.3.2...v0.3.8 ✓ 40000 bytes
    - context: 7549113 bytes
- [#11](https://github.com/mudler/entities/pull/11) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/net 0.0.0-20191209160850-c0dbc17a3553→0.7.0: compare c0dbc17a3553...v0.7.0 ✓ 40000 bytes
    - context: 8673179 bytes
- [#12](https://github.com/mudler/entities/pull/12) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/sys 0.0.0-20200102141924-c96a22e43c9c→0.1.0: compare c96a22e43c9c...v0.1.0 ✓ 40000 bytes
    - context: 8083956 bytes
**[mudler/yip](https://github.com/mudler/yip)**

- [#322](https://github.com/mudler/yip/pull/322) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/google/go-containerregistry 0.21.7→0.22.1: compare v0.21.7...v0.22.1 ✓ 40000 bytes
    - golang.org/x/crypto 0.54.0→0.55.0: compare v0.54.0...v0.55.0 ✓ 40000 bytes
    - context: 108732 bytes
- [#323](https://github.com/mudler/yip/pull/323) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/go-git/go-git/v5 5.19.1→5.19.2: compare v5.19.1...v5.19.2 ✓ 40000 bytes
    - context: 52337 bytes
- [#324](https://github.com/mudler/yip/pull/324) — ✅ **good** — This is a patch-level update to a dependency, which typically introduces only bug fixes and minor improvements. The changelog explicitly mentions a fix, and the diffs show internal refactoring which is expected during dependency upgrades. There are no obvious security vulnerabilities introduced by this version bump.
  ↳ The PR updates the dependency `github.com/onsi/ginkgo/v2` from v2.32.0 to v2.32.1. This version bump includes a fix to defer `AfterAll` until repeated specs complete, along with several internal code and documentation updates related to test ordering and parallel execution logic.
    - github.com/onsi/ginkgo/v2 2.32.0→2.32.1: compare v2.32.0...v2.32.1 ✓ 12922 bytes
    - context: 15994 bytes
- [#325](https://github.com/mudler/yip/pull/325) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.54.0→0.56.0: compare v0.54.0...v0.56.0 ✓ 40000 bytes
    - golang.org/x/mod 0.37.0→0.38.0: compare v0.37.0...v0.38.0 ✓ 10336 bytes
    - golang.org/x/net 0.56.0→0.57.0: compare v0.56.0...v0.57.0 ✓ 40000 bytes
    - context: 97415 bytes
- [#327](https://github.com/mudler/yip/pull/327) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/sirupsen/logrus 1.9.4→1.10.2: compare v1.9.4...v1.10.2 ✓ 40000 bytes
    - sirupsen/logrus v1.10.1..v1.10.2 (PR body): compare v1.10.1...v1.10.2 ✓ 1852 bytes
    - context: 63402 bytes
- [#328](https://github.com/mudler/yip/pull/328) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/mauromorales/xpasswd 0.4.8→0.5.0: compare v0.4.8...v0.5.0 ✓ 15459 bytes
    - github.com/onsi/ginkgo/v2 2.32.0→2.32.1: compare v2.32.0...v2.32.1 ✓ 12922 bytes
    - mauromorales/xpasswd v0.4.9..v0.5.0 (PR body): compare v0.4.9...v0.5.0 ✓ 13436 bytes
    - mauromorales/xpasswd v0.4.8..v0.4.9 (PR body): compare v0.4.8...v0.4.9 ✓ 2023 bytes
    - context: 51133 bytes
- [#334](https://github.com/mudler/yip/pull/334) — ✅ **good** — This is a standard dependency update to a newer version of the library. The changelog confirms that the update includes a new feature (gomock adaptor extension), and the diffs show the corresponding code changes. This change is safe to auto-approve.
  ↳ This PR updates the dependency `github.com/onsi/gomega` from v1.42.1 to v1.43.0. This update introduces a new feature: a gomock adaptor extension that allows Gomega matchers to be used with gomock argument matchers.
    - github.com/onsi/gomega 1.42.1→1.43.0: compare v1.42.1...v1.43.0 ✓ 3786 bytes
    - context: 6898 bytes
- [#335](https://github.com/mudler/yip/pull/335) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/sys 0.47.0→0.48.0: compare v0.47.0...v0.48.0 ✓ 40000 bytes
    - context: 42714 bytes
- [#337](https://github.com/mudler/yip/pull/337) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/containerd/containerd/v2 2.3.3→2.3.5: compare v2.3.3...v2.3.5 ✓ 40000 bytes
    - context: 65326 bytes

