# Kairos Security Dashboard

_Updated 2026-08-30._

🌐 **[Live dashboard](https://kairos-io.github.io/security/)** — the published board with clickable links.

## 📋 This run

- **Scanned:** 19 repos · ⚠️ 1 errored
- **Findings:** 3 (0 critical / 1 high / 2 medium / 0 low / 0 unknown)
- **Informational (not counted):** 58
- **CVE-related PRs:** 2 (2 human)
- **Remediation:** 0 open · 0 superseded · 0 merged · 0 need-human
- **Why:** 3 finding(s); 0 PR(s) open.

> The most urgent finding is F1, which has a high severity. It is recommended to prioritize remediation for F1, followed by F2 and F3, which are medium severity issues affecting the libkcapi package.

## 🔥 Focus now

- [CVE-2026-71226](https://osv.dev/vulnerability/ALPINE-CVE-2026-71226) — High severity vulnerability (CVE-2026-71226) in libkcapi.
- [CVE-2026-71227](https://osv.dev/vulnerability/ALPINE-CVE-2026-71227) — Medium severity vulnerability (CVE-2026-71227) in libkcapi.
- [CVE-2026-71225](https://osv.dev/vulnerability/ALPINE-CVE-2026-71225) — Medium severity vulnerability (CVE-2026-71225) in libkcapi.

## 🌊 Waterfall fronts

_None._

## 📦 Per-repo findings

| Repo | Critical | High | Medium | Total | Status |
|---|---|---|---|---|---|
| [kairos-io/hadron](https://github.com/kairos-io/hadron) | 0 | 1 | 2 | 3 | ok |
| [kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot) | 0 | 0 | 0 | 0 | ⚠️ errors |
| [kairos-io/cluster-api-provider-kairos](https://github.com/kairos-io/cluster-api-provider-kairos) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/entangle](https://github.com/kairos-io/entangle) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/entangle-proxy](https://github.com/kairos-io/entangle-proxy) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/go-nodepair](https://github.com/kairos-io/go-nodepair) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/go-ukify](https://github.com/kairos-io/go-ukify) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos](https://github.com/kairos-io/kairos) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-lab](https://github.com/kairos-io/kairos-lab) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/netboot](https://github.com/kairos-io/netboot) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/provider-kairos](https://github.com/kairos-io/provider-kairos) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/provider-kubernetes](https://github.com/kairos-io/provider-kubernetes) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/tpm-helpers](https://github.com/kairos-io/tpm-helpers) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mauromorales/xpasswd](https://github.com/mauromorales/xpasswd) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/edgevpn](https://github.com/mudler/edgevpn) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/entities](https://github.com/mudler/entities) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/go-pluggable](https://github.com/mudler/go-pluggable) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/yip](https://github.com/mudler/yip) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |

## 🧩 Hadron component CVEs

| Package | Current | Fixed | Severity | CVE |
|---|---|---|---|---|
| libkcapi | 1.5.0 | 1.5.1 | high | [CVE-2026-71226](https://osv.dev/vulnerability/ALPINE-CVE-2026-71226) |
| libkcapi | 1.5.0 | 1.5.1 | medium | [CVE-2026-71227](https://osv.dev/vulnerability/ALPINE-CVE-2026-71227) |
| libkcapi | 1.5.0 | 1.5.1 | medium | [CVE-2026-71225](https://osv.dev/vulnerability/ALPINE-CVE-2026-71225) |

## Informational — not counted

These findings are separated from the counts above: CVEs we are already past, or components accepted as pinned risk.

| Package | Current | Fixed | Severity | CVE | Why |
|---|---|---|---|---|---|
| openssl-fips | 3.1.2 | 3.3.7 | critical | [CVE-2026-31789](https://osv.dev/vulnerability/ALPINE-CVE-2026-31789) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-18798](https://osv.dev/vulnerability/ALPINE-CVE-2026-18798) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| busybox | 1.37.0 | 1.37.0 | low | [CVE-2025-46394](https://osv.dev/vulnerability/ALPINE-CVE-2025-46394) | already-fixed |
| openssl-fips | 3.1.2 | 3.1.6 | critical | [CVE-2024-5535](https://osv.dev/vulnerability/ALPINE-CVE-2024-5535) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-9076](https://osv.dev/vulnerability/ALPINE-CVE-2026-9076) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.5 | medium | [CVE-2025-9231](https://osv.dev/vulnerability/ALPINE-CVE-2025-9231) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| glib | 2.86.2 | 2.66.6 | high | [CVE-2021-27219](https://osv.dev/vulnerability/ALPINE-CVE-2021-27219) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2025-69418](https://osv.dev/vulnerability/ALPINE-CVE-2025-69418) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.1 | medium | [CVE-2025-4575](https://osv.dev/vulnerability/ALPINE-CVE-2025-4575) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.3 | medium | [CVE-2024-12797](https://osv.dev/vulnerability/ALPINE-CVE-2024-12797) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-45447](https://osv.dev/vulnerability/ALPINE-CVE-2026-45447) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-45446](https://osv.dev/vulnerability/ALPINE-CVE-2026-45446) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.8 | high | [CVE-2025-9230](https://osv.dev/vulnerability/ALPINE-CVE-2025-9230) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-14457](https://osv.dev/vulnerability/ALPINE-CVE-2026-14457) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.5 | medium | [CVE-2024-4603](https://osv.dev/vulnerability/ALPINE-CVE-2024-4603) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | unknown | [CVE-2026-75803](https://osv.dev/vulnerability/ALPINE-CVE-2026-75803) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-7383](https://osv.dev/vulnerability/ALPINE-CVE-2026-7383) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2024-0727](https://osv.dev/vulnerability/ALPINE-CVE-2024-0727) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-5678](https://osv.dev/vulnerability/ALPINE-CVE-2023-5678) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.8 | medium | [CVE-2025-9232](https://osv.dev/vulnerability/ALPINE-CVE-2025-9232) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-69420](https://osv.dev/vulnerability/ALPINE-CVE-2025-69420) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-6129](https://osv.dev/vulnerability/ALPINE-CVE-2023-6129) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-6237](https://osv.dev/vulnerability/ALPINE-CVE-2023-6237) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-42766](https://osv.dev/vulnerability/ALPINE-CVE-2026-42766) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | low | [CVE-2026-42770](https://osv.dev/vulnerability/ALPINE-CVE-2026-42770) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | high | [CVE-2023-5363](https://osv.dev/vulnerability/ALPINE-CVE-2023-5363) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.7 | high | [CVE-2024-6119](https://osv.dev/vulnerability/ALPINE-CVE-2024-6119) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| perl | 5.44.0 | 5.26.3 | unknown | [CVE-2018-18311](https://osv.dev/vulnerability/ALPINE-CVE-2018-18311) | already-fixed |
| openssl-fips | 3.1.2 | 3.5.7 | critical | [CVE-2026-34182](https://osv.dev/vulnerability/ALPINE-CVE-2026-34182) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28387](https://osv.dev/vulnerability/ALPINE-CVE-2026-28387) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-15467](https://osv.dev/vulnerability/ALPINE-CVE-2025-15467) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | medium | [CVE-2026-63074](https://osv.dev/vulnerability/ALPINE-CVE-2026-63074) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-63072](https://osv.dev/vulnerability/ALPINE-CVE-2026-63072) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| libxml2 | 2.15.3 | 2.13.8 | high | [CVE-2025-32414](https://osv.dev/vulnerability/ALPINE-CVE-2025-32414) | already-fixed |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-63076](https://osv.dev/vulnerability/ALPINE-CVE-2026-63076) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28388](https://osv.dev/vulnerability/ALPINE-CVE-2026-28388) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-63075](https://osv.dev/vulnerability/ALPINE-CVE-2026-63075) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28389](https://osv.dev/vulnerability/ALPINE-CVE-2026-28389) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2025-68160](https://osv.dev/vulnerability/ALPINE-CVE-2025-68160) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.7 | medium | [CVE-2024-9143](https://osv.dev/vulnerability/ALPINE-CVE-2024-9143) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| libxml2 | 2.15.3 | 2.13.8 | high | [CVE-2025-32415](https://osv.dev/vulnerability/ALPINE-CVE-2025-32415) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28390](https://osv.dev/vulnerability/ALPINE-CVE-2026-28390) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2026-22796](https://osv.dev/vulnerability/ALPINE-CVE-2026-22796) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-34180](https://osv.dev/vulnerability/ALPINE-CVE-2026-34180) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-69421](https://osv.dev/vulnerability/ALPINE-CVE-2025-69421) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| busybox | 1.37.0 | 1.37.0 | low | [CVE-2024-58251](https://osv.dev/vulnerability/ALPINE-CVE-2024-58251) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-69419](https://osv.dev/vulnerability/ALPINE-CVE-2025-69419) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-14456](https://osv.dev/vulnerability/ALPINE-CVE-2026-14456) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2024-2511](https://osv.dev/vulnerability/ALPINE-CVE-2024-2511) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2026-22795](https://osv.dev/vulnerability/ALPINE-CVE-2026-22795) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-31790](https://osv.dev/vulnerability/ALPINE-CVE-2026-31790) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-45445](https://osv.dev/vulnerability/ALPINE-CVE-2026-45445) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.8 | medium | [CVE-2024-13176](https://osv.dev/vulnerability/ALPINE-CVE-2024-13176) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-54874](https://osv.dev/vulnerability/ALPINE-CVE-2026-54874) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | unknown | [CVE-2026-63073](https://osv.dev/vulnerability/ALPINE-CVE-2026-63073) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-42767](https://osv.dev/vulnerability/ALPINE-CVE-2026-42767) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.6 | high | [CVE-2024-4741](https://osv.dev/vulnerability/ALPINE-CVE-2024-4741) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| perl | 5.44.0 | 5.26.3 | unknown | [CVE-2018-18312](https://osv.dev/vulnerability/ALPINE-CVE-2018-18312) | already-fixed |

## ⚠️ 1 collection errors

- [kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot) / sourceCVE: govulncheck: exit status 1: govulncheck: loading packages: There are errors with the provided package patterns: -: # github.com/go-piv/piv-go/v2/piv # [pkg-config --cflags -- libpcsclite] Package libpcsclite was not found in the pkg-config search path. Perhaps you sho … (truncated)

## 📋 Open PRs

**[kairos-io/hadron](https://github.com/kairos-io/hadron)**

- [#559 Automatic version bumps for libffi, libkcapi](https://github.com/kairos-io/hadron/pull/559) — human — tracked
- [#560 Automatic version bumps for busybox](https://github.com/kairos-io/hadron/pull/560) — human — tracked

## 🤖 Bot PR ledger

_No bot PRs yet._

## 🔎 Bot-PR reviews

**[kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot)**

- [#674](https://github.com/kairos-io/AuroraBoot/pull/674) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - microsoft/TypeScript v6.0.3..2bd066d87f5bafd315be9f40889d0a60b9e58e0b (PR body): compare v6.0.3...2bd066d87f5bafd315be9f40889d0a60b9e58e0b failed/empty (no upstream diff)
    - microsoft/TypeScript v6.0.2..v6.0.3 (PR body): compare v6.0.2...v6.0.3 ✓ 40000 bytes
    - microsoft/TypeScript v5.9.3..v6.0.2 (PR body): compare v5.9.3...v6.0.2 failed/empty (no upstream diff)
    - context: 44246 bytes
- [#699](https://github.com/kairos-io/AuroraBoot/pull/699) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - typescript-eslint/typescript-eslint v8.67.0..v8.68.0 (PR body): compare v8.67.0...v8.68.0 ✓ 40000 bytes
    - typescript-eslint/typescript-eslint v8.66.0..v8.67.0 (PR body): compare v8.66.0...v8.67.0 ✓ 40000 bytes
    - context: 98494 bytes
- [#700](https://github.com/kairos-io/AuroraBoot/pull/700) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - vitejs/vite v8.2.1..v8.2.2 (PR body): compare v8.2.1...v8.2.2 ✓ 40000 bytes
    - vitejs/vite v8.2.0..v8.2.1 (PR body): compare v8.2.0...v8.2.1 ✓ 40000 bytes
    - context: 123673 bytes
- [#712](https://github.com/kairos-io/AuroraBoot/pull/712) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/foxboron/sbctl 0.0.0-20250917190250-6b8ed8715652→0.0.0-20260802183653-a7168106e003: compare 6b8ed8715652...a7168106e003 ✓ 18927 bytes
    - context: 21798 bytes
- [#723](https://github.com/kairos-io/AuroraBoot/pull/723) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/go-git/go-git/v5 5.19.1→5.19.2: compare v5.19.1...v5.19.2 ✓ 40000 bytes
    - context: 49866 bytes
- [#740](https://github.com/kairos-io/AuroraBoot/pull/740) — ✅ **good** — This is a routine digest update for a known dependency. Since the change is a maintenance task and the new digest is from the upstream source, it is safe to auto-approve.
  ↳ The PR updates the digest for the `golang.org/x/exp` dependency from `c1d0aac` to `e88cd73`. This change is part of a routine dependency maintenance task.
    - golang.org/x/exp 0.0.0-20260813180055-c1d0aacb2297→0.0.0-20260824195058-e88cd73687aa: compare c1d0aacb2297...e88cd73687aa ✓ 7904 bytes
    - context: 10661 bytes
- [#744](https://github.com/kairos-io/AuroraBoot/pull/744) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - vitejs/vite-plugin-react 39b31735bf79c2dd380eedaba7ed849256f92a29..04cac5020e349f452d76c5a4f6d788ad4b38930a (PR body): compare 39b31735bf79c2dd380eedaba7ed849256f92a29...04cac5020e349f452d76c5a4f6d788ad4b38930a ✓ 40000 bytes
    - vitejs/vite-plugin-react 68c0cb8796ce18bd049c3d05c5210eaf0617eac0..39b31735bf79c2dd380eedaba7ed849256f92a29 (PR body): compare 68c0cb8796ce18bd049c3d05c5210eaf0617eac0...39b31735bf79c2dd380eedaba7ed849256f92a29 ✓ 40000 bytes
    - context: 85269 bytes
- [#747](https://github.com/kairos-io/AuroraBoot/pull/747) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - cypress-io/cypress v15.21.0..v15.21.1 (PR body): compare v15.21.0...v15.21.1 ✓ 32157 bytes
    - cypress-io/cypress v15.20.1..v15.21.0 (PR body): compare v15.20.1...v15.21.0 failed/empty (no upstream diff)
    - cypress-io/cypress v15.20.0..v15.20.1 (PR body): compare v15.20.0...v15.20.1 ✓ 40000 bytes
    - context: 79466 bytes
- [#748](https://github.com/kairos-io/AuroraBoot/pull/748) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - eslint/eslint v10.9.0..5c8c2417b9ff462f2dc4e54a062c59135b45b845 (PR body): compare v10.9.0...5c8c2417b9ff462f2dc4e54a062c59135b45b845 ✓ 7871 bytes
    - eslint/eslint v10.9.0..v10.9.1 (PR body): compare v10.9.0...v10.9.1 ✓ 7871 bytes
    - eslint/eslint v10.8.1..c27bc926e496985eb7911c09eb60914b2e4b5d0f (PR body): compare v10.8.1...c27bc926e496985eb7911c09eb60914b2e4b5d0f ✓ 40000 bytes
    - eslint/eslint v10.8.1..v10.9.0 (PR body): compare v10.8.1...v10.9.0 ✓ 40000 bytes
    - context: 98852 bytes
- [#755](https://github.com/kairos-io/AuroraBoot/pull/755) — ✅ **good** — This is a routine maintenance task to sync the chart definition version with the latest stable release (v0.27.0). The description confirms that the corresponding artifact is already published, and there are no apparent security risks introduced by this version bump.
  ↳ The PR updates the `version` and `appVersion` fields in `deploy/helm/auroraboot/Chart.yaml` to `0.27.0` to ensure Helm deployments target the latest chart release.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 819 bytes
- [#757](https://github.com/kairos-io/AuroraBoot/pull/757) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/kairos-io/kairos/v4 4.2.1-0.20260825070803-b5e8809b4aee→4.2.1-0.20260828082842-a3619b06ea86: compare b5e8809b4aee...a3619b06ea86 ✓ 40000 bytes
    - context: 43022 bytes
- [#761](https://github.com/kairos-io/AuroraBoot/pull/761) — ✅ **good** — The PR updates a dependency to a newer version (`16.3.3`) which includes a bug fix for act() re-entrant behavior. This is a safe and beneficial update that improves the stability of the testing library.
  ↳ This PR updates the `@testing-library/react` dependency from version `16.3.2` to `16.3.3`. This version bump includes a bug fix addressing act() re-entrant issues when dispatching events, improving the stability of the testing library.
    - testing-library/react-testing-library v16.3.2..v16.3.3 (PR body): compare v16.3.2...v16.3.3 ✓ 4007 bytes
    - context: 7111 bytes
**[kairos-io/cluster-api-provider-kairos](https://github.com/kairos-io/cluster-api-provider-kairos)**

- [#38](https://github.com/kairos-io/cluster-api-provider-kairos/pull/38) — ✅ **good** — This pull request is a routine dependency update for golang.org/x/oauth2. Updating to a newer version is standard practice and generally safe, as it addresses potential minor issues or security patches without introducing significant risk.
**[kairos-io/entangle](https://github.com/kairos-io/entangle)**

- [#13](https://github.com/kairos-io/entangle/pull/13) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/emicklei/go-restful 2.9.5+incompatible→2.16.0+incompatible: compare v2.9.5+incompatible...v2.16.0+incompatible failed/empty (no upstream diff)
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - golang.org/x/net 0.55.0→0.56.0: compare v0.55.0...v0.56.0 ✓ 40000 bytes
    - context: 97666 bytes
**[kairos-io/entangle-proxy](https://github.com/kairos-io/entangle-proxy)**

- [#5](https://github.com/kairos-io/entangle-proxy/pull/5) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/onsi/gomega 1.40.0→1.43.0: compare v1.40.0...v1.43.0 failed/empty (no upstream diff)
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - golang.org/x/net 0.55.0→0.56.0: compare v0.55.0...v0.56.0 ✓ 40000 bytes
    - context: 88494 bytes
- [#6](https://github.com/kairos-io/entangle-proxy/pull/6) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - sigs.k8s.io/controller-runtime 0.12.1→0.24.1: compare v0.12.1...v0.24.1 ✓ 40000 bytes
    - context: 98801 bytes
- [#14](https://github.com/kairos-io/entangle-proxy/pull/14) — ✅ **good** — This pull request primarily updates several dependencies to newer versions, including core packages like `golang.org/x` and `google.golang.org/protobuf`. Updating dependencies is a crucial security practice to ensure that known vulnerabilities are patched. The changes appear to be dependency hygiene improvements and do not introduce any obvious security risks.
- [#18](https://github.com/kairos-io/entangle-proxy/pull/18) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/build-push-action v7.2.0..v7.3.0 (PR body): compare v7.2.0...v7.3.0 ✓ 40000 bytes
    - context: 83719 bytes
- [#20](https://github.com/kairos-io/entangle-proxy/pull/20) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - k8s.io/api 0.24.0→0.37.0: compare v0.24.0...v0.37.0 ✓ 40000 bytes
    - context: 129685 bytes
- [#23](https://github.com/kairos-io/entangle-proxy/pull/23) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/checkout v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/checkout v6.0.3..v7.0.0 (PR body): compare v6.0.3...v7.0.0 ✓ 40000 bytes
    - context: 63254 bytes
- [#25](https://github.com/kairos-io/entangle-proxy/pull/25) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/go-logr/logr 1.4.3→1.4.4: compare v1.4.3...v1.4.4 ✓ 40000 bytes
    - context: 44091 bytes
- [#26](https://github.com/kairos-io/entangle-proxy/pull/26) — ✅ **good** — This is a routine version bump for a dependency from the official upstream repository. The changelog indicates a fix for spec lifecycle management, and there are no indications of any security vulnerabilities or malicious code introduced by this update. Therefore, it is safe to auto-approve.
  ↳ This PR updates the dependency `github.com/onsi/ginkgo/v2` from version v2.32.0 to v2.32.1. The update includes a fix to defer `AfterAll` until repeated specs complete, improving spec lifecycle management.
    - github.com/onsi/ginkgo/v2 2.32.0→2.32.1: compare v2.32.0...v2.32.1 ✓ 12922 bytes
    - context: 16130 bytes
- [#27](https://github.com/kairos-io/entangle-proxy/pull/27) — ✅ **good** — This is a minor version bump for the Go base image. Updating the base image is a standard maintenance task and generally improves security and stability by incorporating bug fixes and minor security patches. There are no known security risks associated with this specific version upgrade.
  ↳ The PR updates the base image for the Go build stage in the Dockerfile from golang:1.26 to golang:1.27.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1405 bytes
**[kairos-io/go-nodepair](https://github.com/kairos-io/go-nodepair)**

- [#27](https://github.com/kairos-io/go-nodepair/pull/27) — ✅ **good** — The changes involve updating several core dependencies across the project. The changelogs indicate that these updates include important security patches, such as restricting RSA key sizes in go-libp2p and fixing memory exhaustion attacks in quic-go. This is standard maintenance and security hygiene.
  ↳ This pull request updates several core dependencies, including go-libp2p, quic-go, golang.org/x/crypto, golang.org/x/image, golang.org/x/net, and google.golang.org/protobuf. The updates include critical security fixes, such as mitigating a DoS attack in go-libp2p and addressing memory exhaustion issues in quic-go.
- [#65](https://github.com/kairos-io/go-nodepair/pull/65) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - ipfs/go-log v2.9.1..v2.9.2 (PR body): compare v2.9.1...v2.9.2 ✓ 12110 bytes
    - ipfs/go-log v2.9.0..v2.9.1 (PR body): compare v2.9.0...v2.9.1 ✓ 2149 bytes
    - ipfs/go-log v2.8.2..v2.9.0 (PR body): compare v2.8.2...v2.9.0 ✓ 40000 bytes
    - context: 70615 bytes
- [#66](https://github.com/kairos-io/go-nodepair/pull/66) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - context: 82814 bytes
- [#69](https://github.com/kairos-io/go-nodepair/pull/69) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - google/osv-scanner-action v2.5.0..v2.5.1 (PR body): compare v2.5.0...v2.5.1 ✓ 9140 bytes
    - google/osv-scanner-action v2.3.8..v2.5.0 (PR body): compare v2.3.8...v2.5.0 ✓ 12179 bytes
    - context: 25809 bytes
- [#72](https://github.com/kairos-io/go-nodepair/pull/72) — ✅ **good** — This is a routine dependency update for a testing framework. The changes are a standard version bump and the addition of a new, non-breaking feature (gomock adaptor), which is safe to auto-approve.
  ↳ This PR updates the `github.com/onsi/gomega` dependency from v1.42.1 to v1.43.0. This version includes a new feature: a gomock adaptor extension that allows Gomega matchers to be used with gomock argument matchers.
    - github.com/onsi/gomega 1.42.1→1.43.0: compare v1.42.1...v1.43.0 ✓ 3786 bytes
    - context: 6924 bytes
**[kairos-io/go-ukify](https://github.com/kairos-io/go-ukify)**

- [#59](https://github.com/kairos-io/go-ukify/pull/59) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - securego/gosec v2.28.0..v2.29.0 (PR body): compare v2.28.0...v2.29.0 ✓ 40000 bytes
    - securego/gosec v2.27.1..v2.28.0 (PR body): compare v2.27.1...v2.28.0 ✓ 40000 bytes
    - context: 87925 bytes
- [#60](https://github.com/kairos-io/go-ukify/pull/60) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - context: 42433 bytes
- [#62](https://github.com/kairos-io/go-ukify/pull/62) — ✅ **good** — This is a minor version bump for a widely used library. The changes appear to be internal refactoring, bug fixes, and documentation updates related to spec lifecycle management. There are no apparent security risks introduced by this version update.
  ↳ This PR updates the dependency `github.com/onsi/ginkgo/v2` from v2.32.0 to v2.32.1. The changes include a fix deferring `AfterAll` until repeated specs complete, updates to internal logic for spec attempt handling, and corresponding documentation and test updates.
    - github.com/onsi/ginkgo/v2 2.32.0→2.32.1: compare v2.32.0...v2.32.1 ✓ 12922 bytes
    - context: 16020 bytes
- [#63](https://github.com/kairos-io/go-ukify/pull/63) — ✅ **good** — This is a routine version bump for a well-known dependency, `github.com/onsi/gomega`. The release notes indicate a new feature (gomock adaptor extension), and the diffs show standard version updates in `go.mod` and `go.sum`. There are no immediate security concerns apparent from the context.
  ↳ The PR updates the `github.com/onsi/gomega` dependency from v1.42.1 to v1.43.0, adding a gomock adaptor extension to allow Gomega matchers to be used with gomock argument matchers.
    - github.com/onsi/gomega 1.42.1→1.43.0: compare v1.42.1...v1.43.0 ✓ 3786 bytes
    - context: 6845 bytes
**[kairos-io/hadron](https://github.com/kairos-io/hadron)**

- [#557](https://github.com/kairos-io/hadron/pull/557) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - updatecli/updatecli-action v3.5.0..v3.6.0 (PR body): compare v3.5.0...v3.6.0 ✓ 40000 bytes
    - updatecli/updatecli-action v3.4.0..v3.5.0 (PR body): compare v3.4.0...v3.5.0 ✓ 40000 bytes
    - context: 86163 bytes
- [#561](https://github.com/kairos-io/hadron/pull/561) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/checkout v7..v7.0.1 (PR body): compare v7...v7.0.1 failed/empty (no upstream diff)
    - actions/checkout v7..v7 (PR body): compare v7...v7 failed/empty (no upstream diff)
    - actions/checkout v6.1.0..v7 (PR body): compare v6.1.0...v7 ✓ 40000 bytes
    - actions/checkout v6.0.3..v6.1.0 (PR body): compare v6.0.3...v6.1.0 ✓ 32809 bytes
    - context: 82977 bytes
- [#562](https://github.com/kairos-io/hadron/pull/562) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/setup-qemu-action v4.1.0..v4.2.0 (PR body): compare v4.1.0...v4.2.0 ✓ 40000 bytes
    - docker/setup-qemu-action v4..v4.1.0 (PR body): compare v4...v4.1.0 failed/empty (no upstream diff)
    - docker/setup-qemu-action v4.0.0..v4.1.0 (PR body): compare v4.0.0...v4.1.0 failed/empty (no upstream diff)
    - docker/setup-qemu-action v4..v4 (PR body): compare v4...v4 failed/empty (no upstream diff)
    - docker/setup-qemu-action v3.7.0..v4.0.0 (PR body): compare v3.7.0...v4.0.0 failed/empty (no upstream diff)
    - context: 45783 bytes
**[kairos-io/kairos](https://github.com/kairos-io/kairos)**

- [#4342](https://github.com/kairos-io/kairos/pull/4342) — ✅ **good** — The change is a version bump for a widely used and fundamental action (`actions/checkout`) to a patch release. This is a routine maintenance update, and there are no apparent security risks associated with this specific version change.
  ↳ This PR updates the version of the `actions/checkout` action from v7.0.0 to v7.0.1 in the workflow files. This is a minor patch update to the action.
    - actions/checkout v7..v7.0.1 (PR body): compare v7...v7.0.1 failed/empty (no upstream diff)
    - context: 4313 bytes
- [#4346](https://github.com/kairos-io/kairos/pull/4346) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 39237 bytes
- [#4347](https://github.com/kairos-io/kairos/pull/4347) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/setup-buildx-action bb05f3f5519dd87d3ba754cc423b652a5edd6d2c..37fe631027851001ddb9b187196cc803df7f5f0e (PR body): compare bb05f3f5519dd87d3ba754cc423b652a5edd6d2c...37fe631027851001ddb9b187196cc803df7f5f0e ✓ 40000 bytes
    - context: 43866 bytes
- [#4348](https://github.com/kairos-io/kairos/pull/4348) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github/codeql-action v4.37.6..v4.37.7 (PR body): compare v4.37.6...v4.37.7 ✓ 40000 bytes
    - context: 42447 bytes
- [#4349](https://github.com/kairos-io/kairos/pull/4349) — ✅ **good** — The change is a minor version bump for a security scanning action. The changelog indicates that the updates are bug fixes and feature additions, which are generally considered safe. The PR is generated by an automated tool (Mend Renovate) and the diffs confirm only version bumps are being applied across multiple files. No malicious code or breaking changes are apparent from the context.
  ↳ This PR updates the `google/osv-scanner-action` and related actions from version v2.5.0 to v2.5.1. The update includes bug fixes and new support for environment variables, such as `OSV_SCANNER_LOCAL_DB_CACHE_DIRECTORY`, and improves local vulnerability matching.
    - google/osv-scanner-action v2.5.0..v2.5.1 (PR body): compare v2.5.0...v2.5.1 ✓ 9140 bytes
    - context: 12069 bytes
- [#4352](https://github.com/kairos-io/kairos/pull/4352) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - k8s.io/api 0.36.3→0.36.4: compare v0.36.3...v0.36.4 ✓ 3078 bytes
    - k8s.io/apimachinery 0.36.3→0.36.4: compare v0.36.3...v0.36.4 ✓ 3427 bytes
    - k8s.io/client-go 0.36.3→0.36.4: compare v0.36.3...v0.36.4 ✓ 7769 bytes
    - k8s.io/mount-utils 0.36.3→0.36.4: compare v0.36.3...v0.36.4 ✓ 1432 bytes
    - context: 22199 bytes
- [#4370](https://github.com/kairos-io/kairos/pull/4370) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/google/go-attestation 0.6.1→0.6.4: compare v0.6.1...v0.6.4 ✓ 40000 bytes
    - google/go-attestation v0.6.3..v0.6.4 (PR body): compare v0.6.3...v0.6.4 ✓ 17615 bytes
    - context: 63581 bytes
- [#4371](https://github.com/kairos-io/kairos/pull/4371) — ✅ **good** — This change is a routine dependency update for the golang language across several build components. Updating to a newer, stable minor version is a standard maintenance practice and does not introduce any apparent security risks. Therefore, it is safe to auto-approve.
  ↳ This PR updates the base Docker image tag for the golang dependency across multiple build stages in several Dockerfiles (e.g., from 1.26.6-bookworm to 1.27.0-bookworm). This is a routine maintenance update to incorporate minor version improvements and security patches.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 3298 bytes
- [#4373](https://github.com/kairos-io/kairos/pull/4373) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/sirupsen/logrus 1.9.4→1.10.1: compare v1.9.4...v1.10.1 ✓ 40000 bytes
    - sirupsen/logrus v1.10.0..v1.10.1 (PR body): compare v1.10.0...v1.10.1 ✓ 40000 bytes
    - context: 99693 bytes
- [#4406](https://github.com/kairos-io/kairos/pull/4406) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - azure/login f5d393ae46f8fde4be8b75f32e3fc50e654ad0ca..7ddb5af1ef8758cf1353cf3b42f940aee27ba21c (PR body): compare f5d393ae46f8fde4be8b75f32e3fc50e654ad0ca...7ddb5af1ef8758cf1353cf3b42f940aee27ba21c ✓ 40000 bytes
    - context: 42227 bytes
**[kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator)**

- [#153](https://github.com/kairos-io/kairos-operator/pull/153) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/login-action abd2ef45e78c5afb21d64d4ca52ee8550d9572c7..dbcb813823bdd20940b903addbd779551569679f (PR body): compare abd2ef45e78c5afb21d64d4ca52ee8550d9572c7...dbcb813823bdd20940b903addbd779551569679f ✓ 40000 bytes
    - context: 43836 bytes
- [#156](https://github.com/kairos-io/kairos-operator/pull/156) — ✅ **good** — The pull request is a standard dependency update, specifically changing the digest for the `docker.io/golang:1.26.5` image. This is a routine maintenance task and does not introduce any new security risks or breaking changes. The change is safe to auto-approve.
  ↳ This PR updates the Docker image digest for the `docker.io/golang:1.26.5` dependency from an older SHA to a newer one. This is a standard maintenance update to ensure the build uses the latest artifact for the specified version.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1972 bytes
- [#158](https://github.com/kairos-io/kairos-operator/pull/158) — ✅ **good** — The PR is a routine dependency update to a minor version of a well-known operator. Updating to v0.1.3 is a standard maintenance task and does not introduce obvious security risks. It is safe to auto-approve.
  ↳ This PR updates the Docker tag for the `quay.io/kairos/operator` image in the Kustomization file from version v0.1.2 to v0.1.3. This is a routine dependency update to a newer, presumably stable version.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1500 bytes
- [#162](https://github.com/kairos-io/kairos-operator/pull/162) — ✅ **good** — This is a routine dependency update to a minor version, which is standard maintenance. There are no obvious security risks associated with this version bump, and the change is clearly documented in the changelog.
  ↳ This PR updates the `docker.io/golang` dependency from version `1.26.5` to `1.27.0` across the project's development container configuration and Dockerfiles.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 2622 bytes
- [#163](https://github.com/kairos-io/kairos-operator/pull/163) — ✅ **good** — This change is a dependency pinning update. Pinning a dependency to a specific digest (SHA) is a security best practice that ensures the build uses a known, immutable version of the image, mitigating risks associated with mutable tags.
  ↳ The PR updates the `defaultImage` for the `nodeops` configuration in `values.yaml` to pin the `busybox` image to a specific SHA digest. This change enforces image immutability, which is a security best practice to prevent supply chain attacks via mutable tags.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1570 bytes
- [#164](https://github.com/kairos-io/kairos-operator/pull/164) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/checkout 9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0..3d3c42e5aac5ba805825da76410c181273ba90b1 (PR body): compare 9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0...3d3c42e5aac5ba805825da76410c181273ba90b1 ✓ 40000 bytes
    - context: 41932 bytes
- [#165](https://github.com/kairos-io/kairos-operator/pull/165) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - azure/setup-helm b9e51907a09c216f16ebe8536097933489208112..1a275c3b69536ee54be43f2070a358922e12c8d4 (PR body): compare b9e51907a09c216f16ebe8536097933489208112...1a275c3b69536ee54be43f2070a358922e12c8d4 ✓ 40000 bytes
    - context: 41976 bytes
- [#166](https://github.com/kairos-io/kairos-operator/pull/166) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - azure/setup-helm v5.0.0..v5.0.1 (PR body): compare v5.0.0...v5.0.1 ✓ 40000 bytes
    - azure/setup-helm v5..v5.0.1 (PR body): compare v5...v5.0.1 failed/empty (no upstream diff)
    - azure/setup-helm v5..v5 (PR body): compare v5...v5 failed/empty (no upstream diff)
    - azure/setup-helm v4.3.1..v5.0.0 (PR body): compare v4.3.1...v5.0.0 ✓ 40000 bytes
    - context: 85977 bytes
- [#171](https://github.com/kairos-io/kairos-operator/pull/171) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/setup-buildx-action bb05f3f5519dd87d3ba754cc423b652a5edd6d2c..37fe631027851001ddb9b187196cc803df7f5f0e (PR body): compare bb05f3f5519dd87d3ba754cc423b652a5edd6d2c...37fe631027851001ddb9b187196cc803df7f5f0e ✓ 40000 bytes
    - context: 44259 bytes
- [#172](https://github.com/kairos-io/kairos-operator/pull/172) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - k8s.io/api 0.36.3→0.37.0: compare v0.36.3...v0.37.0 ✓ 40000 bytes
    - k8s.io/apimachinery 0.36.3→0.37.0: compare v0.36.3...v0.37.0 ✓ 40000 bytes
    - context: 101770 bytes
**[kairos-io/netboot](https://github.com/kairos-io/netboot)**

- [#46](https://github.com/kairos-io/netboot/pull/46) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.53.0→0.55.0: compare v0.53.0...v0.55.0 ✓ 40000 bytes
    - golang.org/x/net 0.56.0→0.58.0: compare v0.56.0...v0.58.0 ✓ 40000 bytes
    - context: 83965 bytes
- [#47](https://github.com/kairos-io/netboot/pull/47) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - context: 83203 bytes
- [#48](https://github.com/kairos-io/netboot/pull/48) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.53.0→0.55.0: compare v0.53.0...v0.55.0 ✓ 40000 bytes
    - golang.org/x/net 0.56.0→0.57.0: compare v0.56.0...v0.57.0 ✓ 40000 bytes
    - context: 83976 bytes
- [#49](https://github.com/kairos-io/netboot/pull/49) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - google/osv-scanner-action v2.5.0..v2.5.1 (PR body): compare v2.5.0...v2.5.1 ✓ 9140 bytes
    - google/osv-scanner-action v2.3.8..v2.5.0 (PR body): compare v2.3.8...v2.5.0 ✓ 12179 bytes
    - context: 25025 bytes
**[kairos-io/provider-kairos](https://github.com/kairos-io/provider-kairos)**

- [#965](https://github.com/kairos-io/provider-kairos/pull/965) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/mudler/edgevpn 0.35.3→0.35.4: compare v0.35.3...v0.35.4 ✓ 40000 bytes
    - github.com/labstack/echo/v4 4.15.3→4.15.4: compare v4.15.3...v4.15.4 ✓ 34118 bytes
    - context: 82037 bytes
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
**[mudler/edgevpn](https://github.com/mudler/edgevpn)**

- [#804](https://github.com/mudler/edgevpn/pull/804) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - c-robinson/iplib v2.0.4..v2.0.5 (PR body): compare v2.0.4...v2.0.5 ✓ 6378 bytes
    - c-robinson/iplib v2.0.3..v2.0.4 (PR body): compare v2.0.3...v2.0.4 ✓ 3273 bytes
    - c-robinson/iplib v2.0.2..v2.0.3 (PR body): compare v2.0.2...v2.0.3 ✓ 9999 bytes
    - c-robinson/iplib v2.0.1..v2.0.2 (PR body): compare v2.0.1...v2.0.2 ✓ 15662 bytes
    - c-robinson/iplib v2.0.0..v2.0.1 (PR body): compare v2.0.0...v2.0.1 ✓ 1844 bytes
    - context: 44543 bytes
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
    - context: 45662 bytes
- [#1057](https://github.com/mudler/edgevpn/pull/1057) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - context: 84598 bytes
- [#1059](https://github.com/mudler/edgevpn/pull/1059) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/libp2p/go-libp2p 0.48.0→0.49.0: compare v0.48.0...v0.49.0 ✓ 40000 bytes
    - github.com/libp2p/go-libp2p-kad-dht 0.41.0→0.42.2: compare v0.41.0...v0.42.2 ✓ 40000 bytes
    - context: 117355 bytes
- [#1060](https://github.com/mudler/edgevpn/pull/1060) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/labstack/echo/v5 5.3.0→5.3.1: compare v5.3.0...v5.3.1 ✓ 35382 bytes
    - golang.org/x/crypto 0.53.0→0.54.0: compare v0.53.0...v0.54.0 ✓ 40000 bytes
    - context: 86707 bytes
- [#1061](https://github.com/mudler/edgevpn/pull/1061) — ✅ **good** — The change is a routine dependency digest update. There are no obvious security implications, and this type of maintenance is necessary to keep dependencies current. It is safe to auto-approve.
  ↳ This PR updates the dependency `github.com/mudler/go-libp2p-pubsub` by replacing the old digest (`205ded1`) with a newer one (`2a31b5e`). This is a routine maintenance update to ensure the project uses the latest version of the library.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 2511 bytes
- [#1062](https://github.com/mudler/edgevpn/pull/1062) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/libp2p/go-libp2p 0.48.0→0.49.0: compare v0.48.0...v0.49.0 ✓ 40000 bytes
    - github.com/ipfs/go-cid 0.6.1→0.6.2: compare v0.6.1...v0.6.2 ✓ 8550 bytes
    - github.com/ipfs/go-datastore 0.9.1→0.9.2: compare v0.9.1...v0.9.2 ✓ 16069 bytes
    - context: 93503 bytes
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
    - actions/setup-node v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/setup-node v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - actions/setup-node v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - context: 96419 bytes
- [#1076](https://github.com/mudler/edgevpn/pull/1076) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - vitejs/vite v8.2.1..v8.2.2 (PR body): compare v8.2.1...v8.2.2 ✓ 40000 bytes
    - vitejs/vite v8.2.0..v8.2.1 (PR body): compare v8.2.0...v8.2.1 ✓ 40000 bytes
    - context: 111628 bytes
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
    - github.com/google/go-containerregistry 0.21.7→0.22.0: compare v0.21.7...v0.22.0 ✓ 40000 bytes
    - golang.org/x/crypto 0.54.0→0.55.0: compare v0.54.0...v0.55.0 ✓ 40000 bytes
    - context: 103597 bytes
- [#323](https://github.com/mudler/yip/pull/323) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/go-git/go-git/v5 5.19.1→5.19.2: compare v5.19.1...v5.19.2 ✓ 40000 bytes
    - context: 52469 bytes
- [#324](https://github.com/mudler/yip/pull/324) — ✅ **good** — This PR updates a dependency to a patch release (v2.32.1) which contains a fix for spec timeout behavior. This is a standard maintenance update and does not introduce breaking changes or security vulnerabilities.
  ↳ Updates the `github.com/onsi/ginkgo/v2` dependency from version v2.32.0 to v2.32.1, which includes a fix to defer `AfterAll` until repeated specs complete.
    - github.com/onsi/ginkgo/v2 2.32.0→2.32.1: compare v2.32.0...v2.32.1 ✓ 12922 bytes
    - context: 16126 bytes
- [#325](https://github.com/mudler/yip/pull/325) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.54.0→0.55.0: compare v0.54.0...v0.55.0 ✓ 40000 bytes
    - golang.org/x/mod 0.37.0→0.38.0: compare v0.37.0...v0.38.0 ✓ 10336 bytes
    - golang.org/x/net 0.56.0→0.57.0: compare v0.56.0...v0.57.0 ✓ 40000 bytes
    - context: 97547 bytes
- [#326](https://github.com/mudler/yip/pull/326) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/containerd/containerd/v2 2.3.3→2.3.4: compare v2.3.3...v2.3.4 ✓ 40000 bytes
    - context: 53007 bytes
- [#327](https://github.com/mudler/yip/pull/327) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/sirupsen/logrus 1.9.4→1.10.2: compare v1.9.4...v1.10.2 ✓ 40000 bytes
    - sirupsen/logrus v1.10.1..v1.10.2 (PR body): compare v1.10.1...v1.10.2 ✓ 1852 bytes
    - context: 63402 bytes
- [#328](https://github.com/mudler/yip/pull/328) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/mauromorales/xpasswd 0.4.8→0.5.0: compare v0.4.8...v0.5.0 ✓ 15459 bytes
    - github.com/onsi/ginkgo/v2 2.32.0→2.32.1: compare v2.32.0...v2.32.1 ✓ 12922 bytes
    - mauromorales/xpasswd v0.4.9..v0.5.0 (PR body): compare v0.4.9...v0.5.0 ✓ 13436 bytes
    - mauromorales/xpasswd v0.4.8..v0.4.9 (PR body): compare v0.4.8...v0.4.9 ✓ 2023 bytes
    - context: 51001 bytes
- [#331](https://github.com/mudler/yip/pull/331) — ✅ **good** — This is a routine dependency update to a newer version. The changelog indicates a new feature has been added, and there are no apparent breaking changes mentioned. The update is safe and beneficial.
  ↳ The PR updates the `github.com/onsi/gomega` dependency from v1.42.1 to v1.43.0. This version bump introduces a new feature: a gomock adaptor extension, allowing Gomega matchers to be used as gomock argument matchers.
    - github.com/onsi/gomega 1.42.1→1.43.0: compare v1.42.1...v1.43.0 ✓ 3786 bytes
    - context: 7030 bytes

