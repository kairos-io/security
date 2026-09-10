# Kairos Security Dashboard

_Updated 2026-09-10._

🌐 **[Live dashboard](https://kairos-io.github.io/security/)** — the published board with clickable links.

## 📋 This run

- **Scanned:** 18 repos · ⚠️ 1 errored
- **Findings:** 37 (0 critical / 4 high / 0 medium / 0 low / 33 unknown)
- **Informational (not counted):** 57
- **CVE-related PRs:** 0
- **Remediation:** 0 open · 0 superseded · 0 merged · 0 need-human
- **Why:** 37 finding(s); 0 PR(s) open.

> Focus on the four high-severity findings affecting the 'expat' and 'rsync' packages, as these represent the most critical immediate risks.

## 🔥 Focus now

- [CVE-2026-76957](https://osv.dev/vulnerability/ALPINE-CVE-2026-76957) — High severity vulnerability in expat package (CVE-2026-76957).
- [CVE-2026-53789](https://osv.dev/vulnerability/ALPINE-CVE-2026-53789) — High severity vulnerability in rsync package (CVE-2026-53789).
- [CVE-2026-76956](https://osv.dev/vulnerability/ALPINE-CVE-2026-76956) — High severity vulnerability in expat package (CVE-2026-76956).
- [CVE-2026-70457](https://osv.dev/vulnerability/ALPINE-CVE-2026-70457) — High severity vulnerability in rsync package (CVE-2026-70457).

## 🌊 Waterfall fronts

_None._

## 📦 Per-repo findings

| Repo | Critical | High | Medium | Total | Status |
|---|---|---|---|---|---|
| [kairos-io/hadron](https://github.com/kairos-io/hadron) | 0 | 4 | 0 | 4 | ok |
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
| expat | 2.8.2 | 2.8.4 | high | [CVE-2026-76957](https://osv.dev/vulnerability/ALPINE-CVE-2026-76957) |
| expat | 2.8.2 | 2.8.4 | high | [CVE-2026-76956](https://osv.dev/vulnerability/ALPINE-CVE-2026-76956) |
| rsync | 3.4.4 | 3.5.0 | high | [CVE-2026-53789](https://osv.dev/vulnerability/ALPINE-CVE-2026-53789) |
| rsync | 3.4.4 | 3.5.0 | high | [CVE-2026-70457](https://osv.dev/vulnerability/ALPINE-CVE-2026-70457) |
| expat | 2.8.2 | 2.8.4 | unknown | [CVE-2026-76641](https://osv.dev/vulnerability/ALPINE-CVE-2026-76641) |
| expat | 2.8.2 | 2.8.4 | unknown | [CVE-2026-66046](https://osv.dev/vulnerability/ALPINE-CVE-2026-66046) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53788](https://osv.dev/vulnerability/ALPINE-CVE-2026-53788) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70453](https://osv.dev/vulnerability/ALPINE-CVE-2026-70453) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70452](https://osv.dev/vulnerability/ALPINE-CVE-2026-70452) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53783](https://osv.dev/vulnerability/ALPINE-CVE-2026-53783) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70459](https://osv.dev/vulnerability/ALPINE-CVE-2026-70459) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70458](https://osv.dev/vulnerability/ALPINE-CVE-2026-70458) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53801](https://osv.dev/vulnerability/ALPINE-CVE-2026-53801) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70460](https://osv.dev/vulnerability/ALPINE-CVE-2026-70460) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53794](https://osv.dev/vulnerability/ALPINE-CVE-2026-53794) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53784](https://osv.dev/vulnerability/ALPINE-CVE-2026-53784) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53800](https://osv.dev/vulnerability/ALPINE-CVE-2026-53800) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70461](https://osv.dev/vulnerability/ALPINE-CVE-2026-70461) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53786](https://osv.dev/vulnerability/ALPINE-CVE-2026-53786) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53790](https://osv.dev/vulnerability/ALPINE-CVE-2026-53790) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70454](https://osv.dev/vulnerability/ALPINE-CVE-2026-70454) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53796](https://osv.dev/vulnerability/ALPINE-CVE-2026-53796) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53792](https://osv.dev/vulnerability/ALPINE-CVE-2026-53792) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70455](https://osv.dev/vulnerability/ALPINE-CVE-2026-70455) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53803](https://osv.dev/vulnerability/ALPINE-CVE-2026-53803) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53791](https://osv.dev/vulnerability/ALPINE-CVE-2026-53791) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70464](https://osv.dev/vulnerability/ALPINE-CVE-2026-70464) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70462](https://osv.dev/vulnerability/ALPINE-CVE-2026-70462) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53797](https://osv.dev/vulnerability/ALPINE-CVE-2026-53797) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70463](https://osv.dev/vulnerability/ALPINE-CVE-2026-70463) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53795](https://osv.dev/vulnerability/ALPINE-CVE-2026-53795) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-70456](https://osv.dev/vulnerability/ALPINE-CVE-2026-70456) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53799](https://osv.dev/vulnerability/ALPINE-CVE-2026-53799) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53793](https://osv.dev/vulnerability/ALPINE-CVE-2026-53793) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53802](https://osv.dev/vulnerability/ALPINE-CVE-2026-53802) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53785](https://osv.dev/vulnerability/ALPINE-CVE-2026-53785) |
| rsync | 3.4.4 | 3.5.0 | unknown | [CVE-2026-53798](https://osv.dev/vulnerability/ALPINE-CVE-2026-53798) |

## Informational — not counted

These findings are separated from the counts above: CVEs we are already past, or components accepted as pinned risk.

| Package | Current | Fixed | Severity | CVE | Why |
|---|---|---|---|---|---|
| openssl-fips | 3.1.2 | 3.3.7 | critical | [CVE-2026-31789](https://osv.dev/vulnerability/ALPINE-CVE-2026-31789) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-18798](https://osv.dev/vulnerability/ALPINE-CVE-2026-18798) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
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
| gzip | 1.14 | 1.14 | high | [CVE-2026-41992](https://osv.dev/vulnerability/ALPINE-CVE-2026-41992) | already-fixed |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-14457](https://osv.dev/vulnerability/ALPINE-CVE-2026-14457) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.5 | medium | [CVE-2024-4603](https://osv.dev/vulnerability/ALPINE-CVE-2024-4603) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | critical | [CVE-2026-75803](https://osv.dev/vulnerability/ALPINE-CVE-2026-75803) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
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
| openssl-fips | 3.1.2 | 3.3.6 | high | [CVE-2025-69419](https://osv.dev/vulnerability/ALPINE-CVE-2025-69419) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-14456](https://osv.dev/vulnerability/ALPINE-CVE-2026-14456) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2024-2511](https://osv.dev/vulnerability/ALPINE-CVE-2024-2511) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2026-22795](https://osv.dev/vulnerability/ALPINE-CVE-2026-22795) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-31790](https://osv.dev/vulnerability/ALPINE-CVE-2026-31790) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-45445](https://osv.dev/vulnerability/ALPINE-CVE-2026-45445) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.8 | medium | [CVE-2024-13176](https://osv.dev/vulnerability/ALPINE-CVE-2024-13176) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | high | [CVE-2026-54874](https://osv.dev/vulnerability/ALPINE-CVE-2026-54874) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.8 | critical | [CVE-2026-63073](https://osv.dev/vulnerability/ALPINE-CVE-2026-63073) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
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
- [#740](https://github.com/kairos-io/AuroraBoot/pull/740) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/exp 0.0.0-20260824195058-e88cd73687aa→0.0.0-20260908205506-85c1c2202aba: compare e88cd73687aa...85c1c2202aba ✓ 9589 bytes
    - golang.org/x/mod 0.40.0→0.41.0: compare v0.40.0...v0.41.0 ✓ 1224 bytes
    - golang.org/x/net 0.58.0→0.59.0: compare v0.58.0...v0.59.0 ✓ 40000 bytes
    - golang.org/x/sys 0.47.0→0.48.0: compare v0.47.0...v0.48.0 ✓ 40000 bytes
    - context: 100705 bytes
- [#799](https://github.com/kairos-io/AuroraBoot/pull/799) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - sindresorhus/globals v17.11.0..v17.12.0 (PR body): compare v17.11.0...v17.12.0 ✓ 6645 bytes
    - sindresorhus/globals v17.10.0..8c599278a68a0a6ea17b0c12f976f2270f70f391 (PR body): compare v17.10.0...8c599278a68a0a6ea17b0c12f976f2270f70f391 ✓ 3619 bytes
    - sindresorhus/globals v17.10.0..v17.11.0 (PR body): compare v17.10.0...v17.11.0 ✓ 3619 bytes
    - sindresorhus/globals v17.9.0..7bed4af3730dcb5dbea4274b5264e6be4c4b8910 (PR body): compare v17.9.0...7bed4af3730dcb5dbea4274b5264e6be4c4b8910 ✓ 1014 bytes
    - sindresorhus/globals v17.9.0..v17.10.0 (PR body): compare v17.9.0...v17.10.0 ✓ 1014 bytes
    - context: 19956 bytes
- [#800](https://github.com/kairos-io/AuroraBoot/pull/800) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - helm/helm v3.21.4..v3.22.0-rc.1 (PR body): compare v3.21.4...v3.22.0-rc.1 ✓ 40000 bytes
    - context: 62189 bytes
**[kairos-io/entangle](https://github.com/kairos-io/entangle)**

- [#13](https://github.com/kairos-io/entangle/pull/13) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/emicklei/go-restful 2.9.5+incompatible→2.16.0+incompatible: compare v2.9.5+incompatible...v2.16.0+incompatible failed/empty (no upstream diff)
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - golang.org/x/net 0.55.0→0.56.0: compare v0.55.0...v0.56.0 ✓ 40000 bytes
    - context: 97666 bytes
**[kairos-io/entangle-proxy](https://github.com/kairos-io/entangle-proxy)**

- [#6](https://github.com/kairos-io/entangle-proxy/pull/6) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - sigs.k8s.io/controller-runtime 0.12.1→0.25.0: compare v0.12.1...v0.25.0 ✓ 40000 bytes
    - context: 98865 bytes
- [#14](https://github.com/kairos-io/entangle-proxy/pull/14) — ✅ **good** — This pull request primarily updates several dependencies to newer versions, including core packages like `golang.org/x` and `google.golang.org/protobuf`. Updating dependencies is a crucial security practice to ensure that known vulnerabilities are patched. The changes appear to be dependency hygiene improvements and do not introduce any obvious security risks.
- [#18](https://github.com/kairos-io/entangle-proxy/pull/18) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/build-push-action v7.2.0..v7.3.0 (PR body): compare v7.2.0...v7.3.0 ✓ 40000 bytes
    - context: 82843 bytes
- [#20](https://github.com/kairos-io/entangle-proxy/pull/20) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - k8s.io/api 0.24.0→0.37.0: compare v0.24.0...v0.37.0 ✓ 40000 bytes
    - context: 129685 bytes
- [#23](https://github.com/kairos-io/entangle-proxy/pull/23) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/checkout v7..v7.0.1 (PR body): compare v7...v7.0.1 failed/empty (no upstream diff)
    - actions/checkout v6.1.0..v7 (PR body): compare v6.1.0...v7 ✓ 40000 bytes
    - context: 63408 bytes
- [#25](https://github.com/kairos-io/entangle-proxy/pull/25) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/go-logr/logr 1.4.3→1.4.4: compare v1.4.3...v1.4.4 ✓ 40000 bytes
    - context: 44091 bytes
- [#26](https://github.com/kairos-io/entangle-proxy/pull/26) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/onsi/ginkgo/v2 2.32.0→2.32.2: compare v2.32.0...v2.32.2 ✓ 30208 bytes
    - onsi/ginkgo v2.32.1..v2.32.2 (PR body): compare v2.32.1...v2.32.2 ✓ 18518 bytes
    - onsi/ginkgo v2.32.0..v2.32.1 (PR body): compare v2.32.0...v2.32.1 ✓ 12922 bytes
    - context: 65036 bytes
- [#27](https://github.com/kairos-io/entangle-proxy/pull/27) — ✅ **good** — This is a routine minor version update for the golang base image. Minor version bumps typically include bug fixes and security patches, making this change safe and necessary for maintaining a current build environment. There are no apparent security risks introduced by this update.
  ↳ This PR updates the base image for the Go build stage in the Dockerfile from golang:1.26 to golang:1.27.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1405 bytes
- [#29](https://github.com/kairos-io/entangle-proxy/pull/29) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/onsi/gomega 1.40.0→1.43.0: compare v1.40.0...v1.43.0 failed/empty (no upstream diff)
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - golang.org/x/net 0.55.0→0.56.0: compare v0.55.0...v0.56.0 ✓ 40000 bytes
    - context: 88510 bytes
**[kairos-io/go-nodepair](https://github.com/kairos-io/go-nodepair)**

- [#27](https://github.com/kairos-io/go-nodepair/pull/27) — ✅ **good** — The changes involve updating several core dependencies across the project. The changelogs indicate that these updates include important security patches, such as restricting RSA key sizes in go-libp2p and fixing memory exhaustion attacks in quic-go. This is standard maintenance and security hygiene.
  ↳ This pull request updates several core dependencies, including go-libp2p, quic-go, golang.org/x/crypto, golang.org/x/image, golang.org/x/net, and google.golang.org/protobuf. The updates include critical security fixes, such as mitigating a DoS attack in go-libp2p and addressing memory exhaustion issues in quic-go.
- [#65](https://github.com/kairos-io/go-nodepair/pull/65) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - ipfs/go-log v2.9.1..v2.9.2 (PR body): compare v2.9.1...v2.9.2 ✓ 12110 bytes
    - ipfs/go-log v2.9.0..v2.9.1 (PR body): compare v2.9.0...v2.9.1 ✓ 2149 bytes
    - ipfs/go-log v2.8.2..v2.9.0 (PR body): compare v2.8.2...v2.9.0 ✓ 40000 bytes
    - context: 70615 bytes
- [#66](https://github.com/kairos-io/go-nodepair/pull/66) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - context: 82644 bytes
- [#69](https://github.com/kairos-io/go-nodepair/pull/69) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - google/osv-scanner-action v2.5.0..v2.5.1 (PR body): compare v2.5.0...v2.5.1 ✓ 9140 bytes
    - google/osv-scanner-action v2.3.8..v2.5.0 (PR body): compare v2.3.8...v2.5.0 ✓ 12179 bytes
    - context: 25809 bytes
- [#75](https://github.com/kairos-io/go-nodepair/pull/75) — ✅ **good** — The change is a standard version bump for a widely used testing library. The changelog indicates a new feature addition, which is a common and expected update for dependency maintenance. There are no immediate security red flags apparent from the provided context.
  ↳ This PR updates the dependency `github.com/onsi/gomega` from version v1.42.1 to v1.43.0. The update introduces a new feature: a gomock adaptor extension that allows Gomega matchers to be used with gomock argument matchers.
    - github.com/onsi/gomega 1.42.1→1.43.0: compare v1.42.1...v1.43.0 ✓ 3786 bytes
    - context: 6924 bytes
**[kairos-io/go-ukify](https://github.com/kairos-io/go-ukify)**

- [#59](https://github.com/kairos-io/go-ukify/pull/59) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - securego/gosec v2.28.0..v2.29.0 (PR body): compare v2.28.0...v2.29.0 ✓ 40000 bytes
    - securego/gosec v2.27.1..v2.28.0 (PR body): compare v2.27.1...v2.28.0 ✓ 40000 bytes
    - context: 87925 bytes
- [#60](https://github.com/kairos-io/go-ukify/pull/60) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - context: 83201 bytes
- [#62](https://github.com/kairos-io/go-ukify/pull/62) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/onsi/ginkgo/v2 2.32.0→2.32.2: compare v2.32.0...v2.32.2 ✓ 30208 bytes
    - onsi/ginkgo v2.32.1..v2.32.2 (PR body): compare v2.32.1...v2.32.2 ✓ 18518 bytes
    - onsi/ginkgo v2.32.0..v2.32.1 (PR body): compare v2.32.0...v2.32.1 ✓ 12922 bytes
    - context: 65061 bytes
- [#63](https://github.com/kairos-io/go-ukify/pull/63) — ✅ **good** — This is a standard dependency version bump from a trusted source. The changelog indicates this update adds a new feature (gomock adaptor extension) and does not introduce any apparent security risks. The changes are safe to merge.
  ↳ The PR updates the `github.com/onsi/gomega` dependency to version `v1.43.0`, which introduces a new gomock adaptor extension for using Gomega matchers with gomock.
    - github.com/onsi/gomega 1.42.1→1.43.0: compare v1.42.1...v1.43.0 ✓ 3786 bytes
    - context: 6845 bytes
**[kairos-io/hadron](https://github.com/kairos-io/hadron)**

- [#594](https://github.com/kairos-io/hadron/pull/594) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - react/react v19.2.8..1d34f91dfde6bba84d08b683aaba164c7194dacb (PR body): compare v19.2.8...1d34f91dfde6bba84d08b683aaba164c7194dacb ✓ 40000 bytes
    - react/react v19.2.8..v19.3.0 (PR body): compare v19.2.8...v19.3.0 ✓ 40000 bytes
    - context: 84686 bytes
**[kairos-io/kairos](https://github.com/kairos-io/kairos)**

- [#4464](https://github.com/kairos-io/kairos/pull/4464) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/mudler/edgevpn 0.35.3→0.35.4: compare v0.35.3...v0.35.4 ✓ 40000 bytes
    - context: 45429 bytes
- [#4472](https://github.com/kairos-io/kairos/pull/4472) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - charmbracelet/bubbles v2.2.0..v2.2.1 (PR body): compare v2.2.0...v2.2.1 ✓ 4095 bytes
    - charmbracelet/bubbles v2.1.1..v2.2.0 (PR body): compare v2.1.1...v2.2.0 ✓ 40000 bytes
    - context: 67972 bytes
- [#4473](https://github.com/kairos-io/kairos/pull/4473) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - charmbracelet/bubbletea v2.0.8..v2.0.9 (PR body): compare v2.0.8...v2.0.9 failed/empty (no upstream diff)
    - charmbracelet/bubbletea v2.0.7..v2.0.8 (PR body): compare v2.0.7...v2.0.8 ✓ 6758 bytes
    - charmbracelet/bubbletea v2.0.6..v2.0.7 (PR body): compare v2.0.6...v2.0.7 ✓ 19910 bytes
    - context: 69636 bytes
- [#4507](https://github.com/kairos-io/kairos/pull/4507) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/containerd/containerd/v2 2.3.4→2.3.5: compare v2.3.4...v2.3.5 ✓ 40000 bytes
    - context: 51951 bytes
- [#4509](https://github.com/kairos-io/kairos/pull/4509) — ✅ **good** — This PR is a necessary maintenance update to resolve build failures caused by breaking changes in the golang v1.27 release. Updating the base image ensures the build process remains functional and uses a more current version of the toolchain. This is a safe and required dependency update.
  ↳ The PR updates the `golang` Docker base image in the `examples/bundle/Dockerfile` from version 1.26 to 1.27 to resolve build errors related to internal API changes.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1643 bytes
- [#4554](https://github.com/kairos-io/kairos/pull/4554) — ✅ **good** — This is a routine dependency version bump for a standard library, `golang.org/x/mod`. The changes are confined to updating the version number and minor internal renaming/cleanup in the module files. There are no obvious security vulnerabilities introduced by this version upgrade, and the context suggests this is a standard maintenance task.
  ↳ This PR updates the `golang.org/x/mod` dependency from v0.40.0 to v0.41.0. It also includes minor internal refactoring in `modfile/rule.go` and `modfile/work.go` to rename functions and comments related to tool handling.
    - golang.org/x/mod 0.40.0→0.41.0: compare v0.40.0...v0.41.0 ✓ 1224 bytes
    - context: 4014 bytes
- [#4555](https://github.com/kairos-io/kairos/pull/4555) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/mod 0.40.0→0.41.0: compare v0.40.0...v0.41.0 ✓ 1224 bytes
    - golang.org/x/net 0.58.0→0.59.0: compare v0.58.0...v0.59.0 ✓ 40000 bytes
    - golang.org/x/sys 0.47.0→0.48.0: compare v0.47.0...v0.48.0 ✓ 40000 bytes
    - context: 89854 bytes
- [#4581](https://github.com/kairos-io/kairos/pull/4581) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github/codeql-action cdf488f595d80d6e07e03d4674febd5ab45fa938..b96794f015dfd88f77b49b1c93e0fa7110f94c63 (PR body): compare cdf488f595d80d6e07e03d4674febd5ab45fa938...b96794f015dfd88f77b49b1c93e0fa7110f94c63 ✓ 40000 bytes
    - context: 42954 bytes
- [#4582](https://github.com/kairos-io/kairos/pull/4582) — ✅ **good** — This pull request is a routine maintenance update that changes the digest of a base image. It does not introduce any new code, dependencies, or logic that would introduce a security vulnerability. Therefore, it is safe to auto-approve.
  ↳ The PR updates the Dockerfile to use a new digest (`829f6df217bcbae2b371026e81711d1a787c61b2967ad09d015063663ebafbf7`) for the ubuntu:22.04 base image.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1730 bytes
**[kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator)**

- [#153](https://github.com/kairos-io/kairos-operator/pull/153) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/login-action abd2ef45e78c5afb21d64d4ca52ee8550d9572c7..dbcb813823bdd20940b903addbd779551569679f (PR body): compare abd2ef45e78c5afb21d64d4ca52ee8550d9572c7...dbcb813823bdd20940b903addbd779551569679f ✓ 40000 bytes
    - context: 43836 bytes
- [#156](https://github.com/kairos-io/kairos-operator/pull/156) — ✅ **good** — This change is a routine maintenance update to change the digest of a base image. It does not introduce new dependencies, modify application code, or change security configurations. Therefore, it is safe to auto-approve.
  ↳ This PR updates the digest for the `docker.io/golang:1.26.5` base image from `3aff665` to `705e964` in the Dockerfile and Dockerfile.node-labeler files.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1972 bytes
- [#158](https://github.com/kairos-io/kairos-operator/pull/158) — ✅ **good** — This is a standard minor version update for a known dependency. The change only modifies the tag reference in a configuration file and does not introduce any new code or security risks. It is safe to auto-approve.
  ↳ This PR updates the tag reference for the `quay.io/kairos/operator` dependency from version `v0.1.2` to `v0.2.1` in the Kustomization configuration file.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1500 bytes
- [#162](https://github.com/kairos-io/kairos-operator/pull/162) — ✅ **good** — This is a routine minor version update for a stable dependency (`docker.io/golang`). The change is applied consistently across all relevant build configurations and is expected maintenance work, posing no security risks.
  ↳ This PR updates the `docker.io/golang` dependency from version `1.26.5` to `1.27.1` across the `devcontainer.json` image definition and the base images in the `Dockerfile` for both the builder and node-labeler stages.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 2622 bytes
- [#163](https://github.com/kairos-io/kairos-operator/pull/163) — ✅ **good** — This PR updates the busybox image reference to a specific digest, ensuring build reproducibility and stability by preventing unexpected changes from upstream tags. Pinning dependencies is a standard practice for maintaining a secure and reliable build environment.
  ↳ The PR pins the `defaultImage` for the `nodeops` configuration from the mutable tag `busybox:latest` to the specific SHA digest `sha256:dc2d74b28e4cf8984fa52af1f39bc7c3d9c73760b41a74d629f5d11b1ab28616`.
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
    - azure/setup-helm v4.3.1..v5 (PR body): compare v4.3.1...v5 ✓ 40000 bytes
    - context: 85807 bytes
- [#171](https://github.com/kairos-io/kairos-operator/pull/171) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/setup-buildx-action bb05f3f5519dd87d3ba754cc423b652a5edd6d2c..37fe631027851001ddb9b187196cc803df7f5f0e (PR body): compare bb05f3f5519dd87d3ba754cc423b652a5edd6d2c...37fe631027851001ddb9b187196cc803df7f5f0e ✓ 40000 bytes
    - context: 44259 bytes
- [#172](https://github.com/kairos-io/kairos-operator/pull/172) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - k8s.io/api 0.36.3→0.37.0: compare v0.36.3...v0.37.0 ✓ 40000 bytes
    - k8s.io/apimachinery 0.36.3→0.37.0: compare v0.36.3...v0.37.0 ✓ 40000 bytes
    - context: 132062 bytes
- [#179](https://github.com/kairos-io/kairos-operator/pull/179) — ✅ **good** — This is a routine version bump for existing components. It updates the images to a newer, presumably vetted version (v0.2.0) and does not introduce any new code or security-sensitive changes. Therefore, it is safe to auto-approve.
  ↳ This PR updates the Kairos operator and node-labeler image versions to v0.2.0 in both the chart definition and the kustomization file.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1468 bytes
- [#181](https://github.com/kairos-io/kairos-operator/pull/181) — ✅ **good** — This is a routine version bump for existing components. Since the PR is labeled as a 'chore' and the change involves updating versions of existing images, it is considered safe to auto-approve without further security concerns based on the provided diff.
  ↳ This PR updates the version of the Kairos operator and its node-labeler image from v0.1.3 to v0.2.1 in both the chart definition and the kustomization configuration.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1468 bytes
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

