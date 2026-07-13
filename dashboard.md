# Kairos Security Dashboard

_Updated 2026-07-13._

🌐 **[Live dashboard](https://kairos-io.github.io/security/)** — the published board with clickable links.

## 📋 This run

- **Scanned:** 25 repos
- **Findings:** 6 (0 critical / 4 high / 2 medium / 0 low / 0 unknown)
- **Informational (not counted):** 55
- **CVE-related PRs:** 0
- **Remediation:** 0 open · 0 superseded · 0 merged · 0 need-human
- **Why:** 6 finding(s); 0 PR(s) open.

> There are multiple high-severity vulnerabilities identified across the project's dependencies, primarily impacting the 'openssl' and 'sqlite3' packages. Immediate attention is required to address these critical security risks.

## 🔥 Focus now

- [CVE-2023-4807](https://osv.dev/vulnerability/ALPINE-CVE-2023-4807) — High severity vulnerability in the openssl package.
- [CVE-2022-35737](https://osv.dev/vulnerability/ALPINE-CVE-2022-35737) ⚠️ — High severity vulnerability in the sqlite3 package.
- [CVE-2022-2068](https://osv.dev/vulnerability/ALPINE-CVE-2022-2068) — High severity vulnerability in the openssl package.
- [CVE-2022-1292](https://osv.dev/vulnerability/ALPINE-CVE-2022-1292) — High severity vulnerability in the openssl package.

## 🌊 Waterfall fronts

_None._

## 📦 Per-repo findings

| Repo | Critical | High | Medium | Total | Status |
|---|---|---|---|---|---|
| [kairos-io/hadron](https://github.com/kairos-io/hadron) | 0 | 4 | 2 | 6 | ok |
| [kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/cluster-api-provider-kairos](https://github.com/kairos-io/cluster-api-provider-kairos) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/entangle](https://github.com/kairos-io/entangle) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/entangle-proxy](https://github.com/kairos-io/entangle-proxy) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/go-nodepair](https://github.com/kairos-io/go-nodepair) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/go-ukify](https://github.com/kairos-io/go-ukify) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/immucore](https://github.com/kairos-io/immucore) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos](https://github.com/kairos-io/kairos) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-agent](https://github.com/kairos-io/kairos-agent) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-init](https://github.com/kairos-io/kairos-init) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-installer](https://github.com/kairos-io/kairos-installer) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-lab](https://github.com/kairos-io/kairos-lab) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-sdk](https://github.com/kairos-io/kairos-sdk) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kcrypt-discovery-challenger](https://github.com/kairos-io/kcrypt-discovery-challenger) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
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
| openssl | 3.6.3 | — | high | [CVE-2023-4807](https://osv.dev/vulnerability/ALPINE-CVE-2023-4807) |
| openssl | 3.6.3 | — | high | [CVE-2022-2068](https://osv.dev/vulnerability/ALPINE-CVE-2022-2068) |
| openssl | 3.6.3 | — | high | [CVE-2022-1292](https://osv.dev/vulnerability/ALPINE-CVE-2022-1292) |
| sqlite3 | 3.53.3 | — | high | [CVE-2022-35737](https://osv.dev/vulnerability/ALPINE-CVE-2022-35737) ⚠️ |
| busybox | 1.37.0 | — | medium | [CVE-2021-42376](https://osv.dev/vulnerability/ALPINE-CVE-2021-42376) |
| curl | 8.21.0 | — | medium | [CVE-2021-22897](https://osv.dev/vulnerability/ALPINE-CVE-2021-22897) ⚠️ |

## ⚠️ 2 finding(s) possibly not applicable (AI)

> These findings are still counted and listed above. The AI applicability check thinks they may not affect us — verify the reasoning below and, if you agree, silence via `cve-policy.yaml`.

<details>
<summary>⚠️ [CVE-2022-35737](https://osv.dev/vulnerability/ALPINE-CVE-2022-35737) — [kairos-io/hadron](https://github.com/kairos-io/hadron) (sqlite3 / confidence: high)</summary>

**Reason:** The vulnerability affects versions up to 3.39.x before 3.39.2, and the queried version 3.53.3 is outside this range.

- CVE: `CVE-2022-35737`
- Current: `3.53.3`
- Fixed: `—`
- Checked by: `gemma-4-e2b-it` on 2026-07-13

</details>

<details>
<summary>⚠️ [CVE-2021-22897](https://osv.dev/vulnerability/ALPINE-CVE-2021-22897) — [kairos-io/hadron](https://github.com/kairos-io/hadron) (curl / confidence: high)</summary>

**Reason:** The CVE only affects curl versions 7.61.0 through 7.76.1, and the queried version 8.21.0 is outside this range.

- CVE: `CVE-2021-22897`
- Current: `8.21.0`
- Fixed: `—`
- Checked by: `gemma-4-e2b-it` on 2026-07-13

</details>

## Informational — not counted

These findings are separated from the counts above: CVEs we are already past, or components accepted as pinned risk.

| Package | Current | Fixed | Severity | CVE | Why |
|---|---|---|---|---|---|
| openssl-fips | 3.1.2 | 3.3.7 | critical | [CVE-2026-31789](https://osv.dev/vulnerability/ALPINE-CVE-2026-31789) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
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
| openssl-fips | 3.1.2 | 3.1.5 | medium | [CVE-2024-4603](https://osv.dev/vulnerability/ALPINE-CVE-2024-4603) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-7383](https://osv.dev/vulnerability/ALPINE-CVE-2026-7383) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2024-0727](https://osv.dev/vulnerability/ALPINE-CVE-2024-0727) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-5678](https://osv.dev/vulnerability/ALPINE-CVE-2023-5678) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.8 | medium | [CVE-2025-9232](https://osv.dev/vulnerability/ALPINE-CVE-2025-9232) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | — | high | [CVE-2023-4807](https://osv.dev/vulnerability/ALPINE-CVE-2023-4807) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-69420](https://osv.dev/vulnerability/ALPINE-CVE-2025-69420) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-6129](https://osv.dev/vulnerability/ALPINE-CVE-2023-6129) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-6237](https://osv.dev/vulnerability/ALPINE-CVE-2023-6237) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-42766](https://osv.dev/vulnerability/ALPINE-CVE-2026-42766) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | low | [CVE-2026-42770](https://osv.dev/vulnerability/ALPINE-CVE-2026-42770) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | high | [CVE-2023-5363](https://osv.dev/vulnerability/ALPINE-CVE-2023-5363) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.7 | high | [CVE-2024-6119](https://osv.dev/vulnerability/ALPINE-CVE-2024-6119) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| perl | 5.42.2 | 5.26.3 | unknown | [CVE-2018-18311](https://osv.dev/vulnerability/ALPINE-CVE-2018-18311) | already-fixed |
| openssl-fips | 3.1.2 | 3.5.7 | critical | [CVE-2026-34182](https://osv.dev/vulnerability/ALPINE-CVE-2026-34182) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28387](https://osv.dev/vulnerability/ALPINE-CVE-2026-28387) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-15467](https://osv.dev/vulnerability/ALPINE-CVE-2025-15467) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| libxml2 | 2.15.3 | 2.13.8 | high | [CVE-2025-32414](https://osv.dev/vulnerability/ALPINE-CVE-2025-32414) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28388](https://osv.dev/vulnerability/ALPINE-CVE-2026-28388) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28389](https://osv.dev/vulnerability/ALPINE-CVE-2026-28389) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2025-68160](https://osv.dev/vulnerability/ALPINE-CVE-2025-68160) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.7 | medium | [CVE-2024-9143](https://osv.dev/vulnerability/ALPINE-CVE-2024-9143) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl | 3.6.3 | 3.0.8 | medium | [CVE-2023-0466](https://osv.dev/vulnerability/ALPINE-CVE-2023-0466) | already-fixed |
| libxml2 | 2.15.3 | 2.13.8 | high | [CVE-2025-32415](https://osv.dev/vulnerability/ALPINE-CVE-2025-32415) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28390](https://osv.dev/vulnerability/ALPINE-CVE-2026-28390) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2026-22796](https://osv.dev/vulnerability/ALPINE-CVE-2026-22796) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-34180](https://osv.dev/vulnerability/ALPINE-CVE-2026-34180) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.0.8 | medium | [CVE-2023-0466](https://osv.dev/vulnerability/ALPINE-CVE-2023-0466) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | — | high | [CVE-2022-1292](https://osv.dev/vulnerability/ALPINE-CVE-2022-1292) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | — | high | [CVE-2022-2068](https://osv.dev/vulnerability/ALPINE-CVE-2022-2068) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-69421](https://osv.dev/vulnerability/ALPINE-CVE-2025-69421) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| busybox | 1.37.0 | 1.37.0 | low | [CVE-2024-58251](https://osv.dev/vulnerability/ALPINE-CVE-2024-58251) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-69419](https://osv.dev/vulnerability/ALPINE-CVE-2025-69419) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2024-2511](https://osv.dev/vulnerability/ALPINE-CVE-2024-2511) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2026-22795](https://osv.dev/vulnerability/ALPINE-CVE-2026-22795) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-31790](https://osv.dev/vulnerability/ALPINE-CVE-2026-31790) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-45445](https://osv.dev/vulnerability/ALPINE-CVE-2026-45445) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.8 | medium | [CVE-2024-13176](https://osv.dev/vulnerability/ALPINE-CVE-2024-13176) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-42767](https://osv.dev/vulnerability/ALPINE-CVE-2026-42767) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.6 | high | [CVE-2024-4741](https://osv.dev/vulnerability/ALPINE-CVE-2024-4741) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| gcc | 15.3.0 | 13.2.1_git20231014 | medium | [CVE-2023-4039](https://osv.dev/vulnerability/ALPINE-CVE-2023-4039) | already-fixed |
| rsync | 3.4.4 | 3.2.4 | high | [CVE-2020-14387](https://osv.dev/vulnerability/ALPINE-CVE-2020-14387) | already-fixed |
| perl | 5.42.2 | 5.26.3 | unknown | [CVE-2018-18312](https://osv.dev/vulnerability/ALPINE-CVE-2018-18312) | already-fixed |

## 📋 Open PRs

_None._

## 🤖 Bot PR ledger

_No bot PRs yet._

## 🔎 Bot-PR reviews

**[kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot)**

- [#409](https://github.com/kairos-io/AuroraBoot/pull/409) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/foxboron/sbctl 0.0.0-20240526163235-64e649b31c8e→0.0.0-20260316200809-1b913e78d38c: compare 64e649b31c8e...1b913e78d38c ✓ 40000 bytes
    - github.com/fatih/color 1.15.0→1.17.0: compare v1.15.0...v1.17.0 ✓ 9976 bytes
    - context: 58329 bytes
- [#590](https://github.com/kairos-io/AuroraBoot/pull/590) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - sindresorhus/globals v17.6.0..a19670cc86c1218e915657c55ea02ba3e7623834 (PR body): compare v17.6.0...a19670cc86c1218e915657c55ea02ba3e7623834 ✓ 11637 bytes
    - sindresorhus/globals v17.6.0..v17.7.0 (PR body): compare v17.6.0...v17.7.0 ✓ 11637 bytes
    - sindresorhus/globals v17.5.0..v17.6.0 (PR body): compare v17.5.0...v17.6.0 ✓ 3099 bytes
    - sindresorhus/globals v17.4.0..v17.5.0 (PR body): compare v17.4.0...v17.5.0 ✓ 5103 bytes
    - sindresorhus/globals v17.3.0..v17.4.0 (PR body): compare v17.3.0...v17.4.0 ✓ 4284 bytes
    - context: 45798 bytes
- [#594](https://github.com/kairos-io/AuroraBoot/pull/594) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - facebook/react eslint-plugin-react-hooks@7.1.0..eslint-plugin-react-hooks@7.1.1 (PR body): compare eslint-plugin-react-hooks@7.1.0...eslint-plugin-react-hooks@7.1.1 ✓ 24066 bytes
    - facebook/react 408b38ef7304faf022d2a37110c57efce12c6bad..eslint-plugin-react-hooks@7.1.0 (PR body): compare 408b38ef7304faf022d2a37110c57efce12c6bad...eslint-plugin-react-hooks@7.1.0 ✓ 40000 bytes
    - context: 100027 bytes
- [#599](https://github.com/kairos-io/AuroraBoot/pull/599) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - eslint/eslint v10.0.0..v10.0.1 (PR body): compare v10.0.0...v10.0.1 ✓ 40000 bytes
    - context: 77814 bytes
- [#626](https://github.com/kairos-io/AuroraBoot/pull/626) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - typescript-eslint/typescript-eslint v8.62.1..v8.63.0 (PR body): compare v8.62.1...v8.63.0 ✓ 40000 bytes
    - context: 53882 bytes
- [#627](https://github.com/kairos-io/AuroraBoot/pull/627) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - eemeli/yaml v2.8.4..v2.9.0 (PR body): compare v2.8.4...v2.9.0 ✓ 11907 bytes
    - eemeli/yaml v2.8.3..v2.8.4 (PR body): compare v2.8.3...v2.8.4 ✓ 13617 bytes
    - context: 29659 bytes
- [#629](https://github.com/kairos-io/AuroraBoot/pull/629) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/stmcginnis/gofish 0.22.0→0.23.0: compare v0.22.0...v0.23.0 ✓ 25082 bytes
    - context: 29808 bytes
- [#630](https://github.com/kairos-io/AuroraBoot/pull/630) — ✅ **good** — This is a routine dependency update for a widely used package, `@types/node`, to a minor version. Such updates are generally low-risk and are performed by an automated tool. Therefore, it is safe to auto-approve.
  ↳ This PR updates the `@types/node` dependency from version 25.9.4 to 25.9.5. This is a minor version bump, which typically includes bug fixes or minor updates for the Node.js type definitions.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 2275 bytes
- [#631](https://github.com/kairos-io/AuroraBoot/pull/631) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/mod 0.37.0→0.38.0: compare v0.37.0...v0.38.0 ✓ 10336 bytes
    - golang.org/x/tools 0.46.0→0.47.0: compare v0.46.0...v0.47.0 ✓ 40000 bytes
    - context: 54560 bytes
- [#632](https://github.com/kairos-io/AuroraBoot/pull/632) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/sys 0.46.0→0.47.0: compare v0.46.0...v0.47.0 ✓ 33531 bytes
    - context: 36312 bytes
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
    - github.com/onsi/gomega 1.40.0→1.42.1: compare v1.40.0...v1.42.1 ✓ 40000 bytes
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - context: 88243 bytes
- [#6](https://github.com/kairos-io/entangle-proxy/pull/6) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - sigs.k8s.io/controller-runtime 0.12.1→0.24.1: compare v0.12.1...v0.24.1 ✓ 40000 bytes
    - context: 98801 bytes
- [#14](https://github.com/kairos-io/entangle-proxy/pull/14) — ✅ **good** — This pull request primarily updates several dependencies to newer versions, including core packages like `golang.org/x` and `google.golang.org/protobuf`. Updating dependencies is a crucial security practice to ensure that known vulnerabilities are patched. The changes appear to be dependency hygiene improvements and do not introduce any obvious security risks.
- [#18](https://github.com/kairos-io/entangle-proxy/pull/18) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/build-push-action v7.2.0..v7.3.0 (PR body): compare v7.2.0...v7.3.0 ✓ 40000 bytes
    - context: 83719 bytes
- [#20](https://github.com/kairos-io/entangle-proxy/pull/20) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - k8s.io/api 0.24.0→0.36.2: compare v0.24.0...v0.36.2 ✓ 40000 bytes
    - context: 126088 bytes
- [#23](https://github.com/kairos-io/entangle-proxy/pull/23) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/checkout v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/checkout v6.0.3..v7.0.0 (PR body): compare v6.0.3...v7.0.0 ✓ 40000 bytes
    - context: 63254 bytes
**[kairos-io/go-nodepair](https://github.com/kairos-io/go-nodepair)**

- [#27](https://github.com/kairos-io/go-nodepair/pull/27) — ✅ **good** — The changes involve updating several core dependencies across the project. The changelogs indicate that these updates include important security patches, such as restricting RSA key sizes in go-libp2p and fixing memory exhaustion attacks in quic-go. This is standard maintenance and security hygiene.
  ↳ This pull request updates several core dependencies, including go-libp2p, quic-go, golang.org/x/crypto, golang.org/x/image, golang.org/x/net, and google.golang.org/protobuf. The updates include critical security fixes, such as mitigating a DoS attack in go-libp2p and addressing memory exhaustion issues in quic-go.
- [#65](https://github.com/kairos-io/go-nodepair/pull/65) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - ipfs/go-log v2.9.1..v2.9.2 (PR body): compare v2.9.1...v2.9.2 ✓ 12110 bytes
    - ipfs/go-log v2.9.0..v2.9.1 (PR body): compare v2.9.0...v2.9.1 ✓ 2149 bytes
    - ipfs/go-log v2.8.2..v2.9.0 (PR body): compare v2.8.2...v2.9.0 ✓ 40000 bytes
    - context: 70615 bytes
**[kairos-io/hadron](https://github.com/kairos-io/hadron)**

- [#506](https://github.com/kairos-io/hadron/pull/506) — ⚠️ **needs_human_verification** — This change involves a significant architectural shift by moving from an in-tree DRBD module to an out-of-tree DRBD 9 module, which is experimental. Although the PR description outlines the necessary steps, it explicitly flags potential risks regarding kernel compatibility layers (e.g., coccinelle/spatch) that must be verified in CI. Therefore, a human review is required to confirm the build stability and ensure all dependencies are correctly handled before merging.
  ↳ This pull request implements a major architectural change by dropping the in-tree DRBD 8 module and replacing it with an out-of-tree DRBD 9 module and its utilities. It modifies kernel configurations to disable DRBD 8 and updates the Dockerfile to download, build, and integrate the new DRBD 9 components.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 10912 bytes
- [#512](https://github.com/kairos-io/hadron/pull/512) — ✅ **good** — The changes directly address the stated requirement to enable kernel BTF support, which is necessary for CO-RE eBPF consumers. The implementation involves necessary kernel configuration modifications, integration of the `pahole` toolchain via `dwarves`, and the addition of verification steps in the CI workflows to ensure correctness. This is a necessary feature enablement and the associated verification makes it safe to auto-approve.
  ↳ This PR enables kernel BTF support by modifying kernel configuration files across multiple architectures to enable `CONFIG_DEBUG_INFO_BTF=y`. It also updates the build toolchain by adding `pahole` support via a new `dwarves-download` stage in the Dockerfile, and introduces CI checks to verify that the required BTF configuration is correctly set in all kernel variants.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 13446 bytes
**[kairos-io/immucore](https://github.com/kairos-io/immucore)**

- [#591](https://github.com/kairos-io/immucore/pull/591) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/jaypipes/ghw 0.24.0→0.25.0: compare v0.24.0...v0.25.0 ✓ 40000 bytes
    - context: 44957 bytes
**[kairos-io/kairos](https://github.com/kairos-io/kairos)**

- [#4229](https://github.com/kairos-io/kairos/pull/4229) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/login-action c99871dec2022cc055c062a10cc1a1310835ceb4..af1e73f918a031802d376d3c8bbc3fe56130a9b0 (PR body): compare c99871dec2022cc055c062a10cc1a1310835ceb4...af1e73f918a031802d376d3c8bbc3fe56130a9b0 ✓ 40000 bytes
    - context: 43062 bytes
- [#4234](https://github.com/kairos-io/kairos/pull/4234) — ✅ **good** — This is a routine minor version bump from a trusted dependency, and the changes are documented in the upstream release notes. The update is applied consistently across all relevant workflow files, suggesting a safe and necessary maintenance update.
  ↳ This PR updates the dependency `kairos-io/kairos-factory-action` from v1.1.3 to v1.2.0. This update incorporates changes from the upstream release, including updates to `actions/checkout` and `github/codeql-action`. The change is applied across multiple CI/CD workflow files.
    - kairos-io/kairos-factory-action v1.1.3..v1.2.0 (PR body): compare v1.1.3...v1.2.0 ✓ 7319 bytes
    - context: 15709 bytes
- [#4236](https://github.com/kairos-io/kairos/pull/4236) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/stale v10.3.0..v10.4.0 (PR body): compare v10.3.0...v10.4.0 ✓ 40000 bytes
    - context: 43055 bytes
**[kairos-io/kairos-init](https://github.com/kairos-io/kairos-init)**

- [#397](https://github.com/kairos-io/kairos-init/pull/397) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - kairos-io/hadron v0.5.0..v0.5.1 (PR body): compare v0.5.0...v0.5.1 ✓ 40000 bytes
    - context: 48773 bytes
**[kairos-io/kairos-installer](https://github.com/kairos-io/kairos-installer)**

- [#14](https://github.com/kairos-io/kairos-installer/pull/14) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/jaypipes/ghw 0.24.0→0.25.0: compare v0.24.0...v0.25.0 ✓ 40000 bytes
    - github.com/containerd/containerd/v2 2.3.2→2.3.3: compare v2.3.2...v2.3.3 ✓ 40000 bytes
    - context: 87431 bytes
**[kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator)**

- [#136](https://github.com/kairos-io/kairos-operator/pull/136) — ✅ **good** — This is a routine patch update for a standard toolchain dependency. Updating to a newer patch version is generally safe and beneficial for receiving bug fixes and minor security patches. There are no apparent security risks introduced by this change.
  ↳ This PR updates the base image tag for `docker.io/golang` from version `1.26.4` to `1.26.5` in the relevant Dockerfiles. This is a routine dependency update for a toolchain component.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1872 bytes
- [#138](https://github.com/kairos-io/kairos-operator/pull/138) — ✅ **good** — This change is a routine digest update for a base image, which is a standard maintenance practice. It does not introduce any new code, functional changes, or security vulnerabilities that would require manual review. Therefore, it is safe to auto-approve.
  ↳ This PR updates the digest of the `gcr.io/distroless/static:nonroot` Docker base image from an older SHA to a newer one. This is a routine maintenance update to ensure the build uses the latest, verified version of the base image.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1727 bytes
**[kairos-io/kcrypt-discovery-challenger](https://github.com/kairos-io/kcrypt-discovery-challenger)**

- [#41](https://github.com/kairos-io/kcrypt-discovery-challenger/pull/41) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - k8s.io/api 0.27.2→0.36.0: compare v0.27.2...v0.36.0 ✓ 40000 bytes
    - context: 123952 bytes
- [#190](https://github.com/kairos-io/kcrypt-discovery-challenger/pull/190) — ✅ **good** — Updating core infrastructure dependencies like Kubernetes components to the latest stable version is a crucial security and stability practice. This change incorporates bug fixes and security patches from the upstream, making the project more resilient. Therefore, it is safe to auto-approve.
  ↳ This PR updates the core Kubernetes dependencies, k8s.io/api, k8s.io/apimachinery, and k8s.io/client-go, to version v0.36.2. This brings the project up to a recent, patched version of the Kubernetes ecosystem components.
    - k8s.io/apimachinery 0.27.4→0.27.2: compare v0.27.4...v0.27.2 failed: <nil> (no upstream diff)
    - github.com/emicklei/go-restful/v3 3.10.1→3.13.0: compare v3.10.1...v3.13.0 ✓ 40000 bytes
    - context: 131955 bytes
- [#240](https://github.com/kairos-io/kcrypt-discovery-challenger/pull/240) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/google/go-attestation 0.5.1→0.6.1: compare v0.5.1...v0.6.1 ✓ 40000 bytes
    - github.com/kairos-io/tpm-helpers 0.0.0-20260608091616-8a4ccb53d8f7→0.0.0-20260702080541-9b3e057e2f32: compare 8a4ccb53d8f7...9b3e057e2f32 ✓ 11771 bytes
    - github.com/google/go-tpm-tools 0.4.4→0.4.7: compare v0.4.4...v0.4.7 ✓ 40000 bytes
    - context: 97188 bytes
- [#241](https://github.com/kairos-io/kcrypt-discovery-challenger/pull/241) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/google/go-attestation 0.5.1→0.6.1: compare v0.5.1...v0.6.1 ✓ 40000 bytes
    - github.com/kairos-io/kairos-sdk 0.23.1→0.23.2: compare v0.23.1...v0.23.2 ✓ 26235 bytes
    - context: 77901 bytes
- [#243](https://github.com/kairos-io/kcrypt-discovery-challenger/pull/243) — ✅ **good** — This is a minor patch update to a dependency. Patch updates are generally low-risk and are safe to auto-approve, as they typically address bug fixes or minor non-security related improvements.
  ↳ This PR updates the version of the `quay.io/kairos/kairos-init` dependency from v0.15.1 to v0.15.2 by changing the build argument in `Dockerfile.kairos-image`.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1501 bytes
- [#244](https://github.com/kairos-io/kcrypt-discovery-challenger/pull/244) — ✅ **good** — This change is a routine patch update to the Go language version. Updating to a newer patch version is generally safe and often includes important bug fixes or security patches. Since this is an automated dependency update, it is considered safe to auto-approve.
  ↳ This PR updates the base image for the build process in the Dockerfile and Dockerfile.kairos-image from `golang:1.26.4` to `golang:1.26.5` and its corresponding bookworm tags. This is a standard dependency bump for the Go language version.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 2473 bytes
**[kairos-io/netboot](https://github.com/kairos-io/netboot)**

- [#45](https://github.com/kairos-io/netboot/pull/45) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.53.0→0.54.0: compare v0.53.0...v0.54.0 ✓ 40000 bytes
    - golang.org/x/sys 0.46.0→0.47.0: compare v0.46.0...v0.47.0 ✓ 33531 bytes
    - context: 76977 bytes
- [#46](https://github.com/kairos-io/netboot/pull/46) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.53.0→0.54.0: compare v0.53.0...v0.54.0 ✓ 40000 bytes
    - golang.org/x/net 0.56.0→0.57.0: compare v0.56.0...v0.57.0 ✓ 40000 bytes
    - context: 83965 bytes
**[mudler/edgevpn](https://github.com/mudler/edgevpn)**

- [#804](https://github.com/mudler/edgevpn/pull/804) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - c-robinson/iplib v2.0.4..v2.0.5 (PR body): compare v2.0.4...v2.0.5 ✓ 6378 bytes
    - c-robinson/iplib v2.0.3..v2.0.4 (PR body): compare v2.0.3...v2.0.4 ✓ 3273 bytes
    - c-robinson/iplib v2.0.2..v2.0.3 (PR body): compare v2.0.2...v2.0.3 ✓ 9999 bytes
    - c-robinson/iplib v2.0.1..v2.0.2 (PR body): compare v2.0.1...v2.0.2 ✓ 15662 bytes
    - c-robinson/iplib v2.0.0..v2.0.1 (PR body): compare v2.0.0...v2.0.1 ✓ 1844 bytes
    - context: 44543 bytes
- [#805](https://github.com/mudler/edgevpn/pull/805) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - go-yaml/yaml v3.0.0..v3.0.1 (PR body): compare v3.0.0...v3.0.1 ✓ 2202 bytes
    - go-yaml/yaml v2.4.0..v3.0.0 (PR body): compare v2.4.0...v3.0.0 ✓ 40000 bytes
    - context: 44433 bytes
- [#808](https://github.com/mudler/edgevpn/pull/808) — ⚠️ **needs_human_verification** — review endpoint unreachable: Post "http://localhost:8080/v1/chat/completions": context deadline exceeded
    - FortAwesome/Font-Awesome 6.7.1..6.7.2 (PR body): compare 6.7.1...6.7.2 ✓ 40000 bytes
    - FortAwesome/Font-Awesome 6.7.0..6.7.1 (PR body): compare 6.7.0...6.7.1 ✓ 40000 bytes
    - context: 83016 bytes
- [#905](https://github.com/mudler/edgevpn/pull/905) — ✅ **good** — This is a routine dependency bump for a tool used in the CI/CD workflow. The changelog indicates that version 2.4.0 includes various maintenance updates and fixes, suggesting this is a safe and necessary update. There are no immediate security red flags indicated by the context.
  ↳ This pull request updates the version of the `dependabot/fetch-metadata` dependency from 2.3.0 to 2.4.0. This upgrade incorporates various fixes, updates to actions, and improvements to the dependency fetching mechanism.
    - dependabot/fetch-metadata v2..v2.4.0 (PR body): compare v2...v2.4.0 failed/empty (no upstream diff)
    - dependabot/fetch-metadata v2.3.0..v2.4.0 (PR body): compare v2.3.0...v2.4.0 ✓ 40000 bytes
    - context: 49729 bytes
- [#914](https://github.com/mudler/edgevpn/pull/914) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/labstack/echo/v4 4.15.2→4.15.4: compare v4.15.2...v4.15.4 ✓ 30288 bytes
    - github.com/mattn/go-colorable 0.1.14→0.1.15: compare v0.1.14...v0.1.15 ✓ 5234 bytes
    - labstack/echo v4.15.3..v4.15.4 (PR body): compare v4.15.3...v4.15.4 ✓ 34118 bytes
    - context: 77780 bytes
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
- [#939](https://github.com/mudler/edgevpn/pull/939) — ✅ **good** — The change is a dependency bump for a widely used action, which is a standard maintenance task. The changelog indicates breaking changes, specifically a Node.js runtime upgrade, which requires verification in the CI pipeline. Assuming the project's CI passes successfully after this update, the change is safe to auto-approve.
  ↳ This PR upgrades the `actions/setup-go` dependency from version 5 to 6. This involves updating references in various GitHub Actions workflow files to use the new action version and incorporates breaking changes from v6.0.0, such as upgrading the Node.js runtime to node 24.x in affected workflows.
    - actions/setup-go v5..v6.0.0 (PR body): compare v5...v6.0.0 ✓ 40000 bytes
    - actions/setup-go v5..v5.5.0 (PR body): compare v5...v5.5.0 failed/empty (no upstream diff)
    - actions/setup-go v5..v6 (PR body): compare v5...v6 ✓ 40000 bytes
    - context: 92926 bytes
- [#942](https://github.com/mudler/edgevpn/pull/942) — ✅ **good** — This is a routine dependency update to a newer minor version of a well-known testing library. The changes primarily involve version bumps and internal code refactoring, which are typical for dependency maintenance. Since this is a standard update and the changes appear to be focused on compatibility and minor fixes, it is safe to auto-approve.
  ↳ This PR bumps github.com/onsi/gomega to version 1.38.2 and updates several related dependencies, including golang.org/x/net, google.golang.org/protobuf, and gopkg.in/yaml.v3. It also includes internal refactoring in gstruct to improve handling of unexported fields and updates to internal error handling.
    - github.com/onsi/gomega 1.37.0→1.38.2: compare v1.37.0...v1.38.2 ✓ 34194 bytes
    - github.com/Masterminds/semver/v3 3.3.1→3.4.0: compare v3.3.1...v3.4.0 ✓ 40000 bytes
    - context: 89137 bytes
- [#943](https://github.com/mudler/edgevpn/pull/943) — ✅ **good** — This is a routine dependency update for a well-known action. The changes are confined to updating the version number, and the upstream changes detailed in the changelog appear to be standard maintenance and minor feature updates, posing no immediate security risk.
  ↳ This pull request updates the `codecov/codecov-action` dependency from version 5.5.0 to 5.5.1. This version bump incorporates several underlying dependency updates for related actions, such as `actions/checkout` and `github/codeql-action`.
    - codecov/codecov-action v5.5.0..v5.5.1 (PR body): compare v5.5.0...v5.5.1 ✓ 10680 bytes
    - context: 21031 bytes
- [#946](https://github.com/mudler/edgevpn/pull/946) — ⚠️ **needs_human_verification** — review endpoint unreachable: Post "http://localhost:8080/v1/chat/completions": context deadline exceeded
    - github.com/libp2p/go-libp2p-kad-dht 0.36.0→0.39.0: compare v0.36.0...v0.39.0 ✓ 40000 bytes
    - golang.org/x/sys 0.41.0→0.42.0: compare v0.41.0...v0.42.0 ✓ 28932 bytes
    - context: 212268 bytes
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
- [#1001](https://github.com/mudler/edgevpn/pull/1001) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - FortAwesome/Font-Awesome 7.2.0..7.3.0 (PR body): compare 7.2.0...7.3.0 ✓ 40000 bytes
    - FortAwesome/Font-Awesome 7.1.0..7.2.0 (PR body): compare 7.1.0...7.2.0 ✓ 40000 bytes
    - context: 84231 bytes
- [#1006](https://github.com/mudler/edgevpn/pull/1006) — ✅ **good** — The upgrade is to a newer minor version (4.15.1) which includes security enhancements, such as the new CSRF middleware features detailed in the release notes. There are no immediate red flags or known critical vulnerabilities associated with this specific version jump. Therefore, this change is safe to auto-approve.
  ↳ This pull request updates the dependency `github.com/labstack/echo/v4` from version 4.13.3 to 4.15.1. This upgrade incorporates several enhancements, including improved CSRF protection features and minor internal fixes related to time comparison logic.
    - github.com/labstack/echo/v4 4.13.3→4.15.1: compare v4.13.3...v4.15.1 ✓ 40000 bytes
    - github.com/mattn/go-colorable 0.1.13→0.1.14: compare v0.1.13...v0.1.14 ✓ 6350 bytes
    - golang.org/x/time 0.12.0→0.14.0: compare v0.12.0...v0.14.0 ✓ 606 bytes
    - context: 76092 bytes
- [#1041](https://github.com/mudler/edgevpn/pull/1041) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/labstack/echo/v4 4.15.2→4.15.4: compare v4.15.2...v4.15.4 ✓ 30288 bytes
    - github.com/labstack/echo/v5 5.2.1→5.3.0: compare v5.2.1...v5.3.0 ✓ 40000 bytes
    - context: 99797 bytes
- [#1046](https://github.com/mudler/edgevpn/pull/1046) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - postcss/autoprefixer 10.5.1..10.5.2 (PR body): compare 10.5.1...10.5.2 ✓ 2688 bytes
    - postcss/autoprefixer 10.5.0..10.5.1 (PR body): compare 10.5.0...10.5.1 ✓ 40000 bytes
    - context: 54014 bytes
- [#1054](https://github.com/mudler/edgevpn/pull/1054) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - urfave/cli v3.10.0..v3.10.1 (PR body): compare v3.10.0...v3.10.1 ✓ 17319 bytes
    - urfave/cli v3.9.1..v3.10.0 (PR body): compare v3.9.1...v3.10.0 ✓ 40000 bytes
    - context: 100504 bytes
- [#1055](https://github.com/mudler/edgevpn/pull/1055) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/sys 0.46.0→0.47.0: compare v0.46.0...v0.47.0 ✓ 33531 bytes
    - context: 36260 bytes
- [#1056](https://github.com/mudler/edgevpn/pull/1056) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/libp2p/go-libp2p-pubsub 0.16.0→0.17.0: compare v0.16.0...v0.17.0 ✓ 40000 bytes
    - context: 45662 bytes
**[mudler/yip](https://github.com/mudler/yip)**

- [#310](https://github.com/mudler/yip/pull/310) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/sys 0.46.0→0.47.0: compare v0.46.0...v0.47.0 ✓ 33531 bytes
    - context: 36249 bytes
- [#311](https://github.com/mudler/yip/pull/311) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.53.0→0.54.0: compare v0.53.0...v0.54.0 ✓ 40000 bytes
    - golang.org/x/sys 0.46.0→0.47.0: compare v0.46.0...v0.47.0 ✓ 33531 bytes
    - context: 80404 bytes
- [#312](https://github.com/mudler/yip/pull/312) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/containerd/containerd/v2 2.3.2→2.3.3: compare v2.3.2...v2.3.3 ✓ 40000 bytes
    - context: 49610 bytes

