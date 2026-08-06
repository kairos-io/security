# Kairos Security Dashboard

_Updated 2026-08-06._

🌐 **[Live dashboard](https://kairos-io.github.io/security/)** — the published board with clickable links.

## 📋 This run

- **Scanned:** 25 repos · ⚠️ 1 errored
- **Findings:** 0 (0 critical / 0 high / 0 medium / 0 low / 0 unknown)
- **Informational (not counted):** 48
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
| libxml2 | 2.15.3 | 2.13.8 | high | [CVE-2025-32414](https://osv.dev/vulnerability/ALPINE-CVE-2025-32414) | already-fixed |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-28388](https://osv.dev/vulnerability/ALPINE-CVE-2026-28388) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
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
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2024-2511](https://osv.dev/vulnerability/ALPINE-CVE-2024-2511) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.6 | medium | [CVE-2026-22795](https://osv.dev/vulnerability/ALPINE-CVE-2026-22795) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.3.7 | high | [CVE-2026-31790](https://osv.dev/vulnerability/ALPINE-CVE-2026-31790) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-45445](https://osv.dev/vulnerability/ALPINE-CVE-2026-45445) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
| openssl-fips | 3.1.2 | 3.1.8 | medium | [CVE-2024-13176](https://osv.dev/vulnerability/ALPINE-CVE-2024-13176) | accepted-component: FIPS 140-3 validated module, pinned at 3.1.2; cannot bump without revalidation |
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
- [#695](https://github.com/kairos-io/AuroraBoot/pull/695) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - tailwindlabs/tailwindcss v4.3.2..v4.3.3 (PR body): compare v4.3.2...v4.3.3 ✓ 40000 bytes
    - context: 58866 bytes
- [#699](https://github.com/kairos-io/AuroraBoot/pull/699) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - typescript-eslint/typescript-eslint v8.65.0..v8.66.0 (PR body): compare v8.65.0...v8.66.0 ✓ 40000 bytes
    - typescript-eslint/typescript-eslint v8.64.0..v8.65.0 (PR body): compare v8.64.0...v8.65.0 ✓ 40000 bytes
    - context: 96848 bytes
- [#700](https://github.com/kairos-io/AuroraBoot/pull/700) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - vitejs/vite v8.1.5..v8.2.0 (PR body): compare v8.1.5...v8.2.0 ✓ 40000 bytes
    - context: 72989 bytes
- [#701](https://github.com/kairos-io/AuroraBoot/pull/701) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - eemeli/yaml v2.8.4..v2.9.0 (PR body): compare v2.8.4...v2.9.0 ✓ 11907 bytes
    - eemeli/yaml v2.8.3..v2.8.4 (PR body): compare v2.8.3...v2.8.4 ✓ 13617 bytes
    - context: 29655 bytes
- [#703](https://github.com/kairos-io/AuroraBoot/pull/703) — ✅ **good** — This is a routine version bump for a Helm chart to align with a newly released version (0.26.0). The description confirms this is intended to target the released image, and there are no security implications associated with updating to a stable, released version.
  ↳ This PR updates the `version` and `appVersion` in `deploy/helm/auroraboot/Chart.yaml` to `0.26.0`. This change ensures that Helm installations target the officially released chart artifact.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 816 bytes
- [#704](https://github.com/kairos-io/AuroraBoot/pull/704) — ✅ **good** — This pull request is a routine maintenance update to sync the Helm chart version to the latest released tag (0.26.1). Since the description confirms the artifact is published and the change is to a known, released version, it poses no security risk and is safe to auto-approve.
  ↳ This change updates the `version` and `appVersion` in `Chart.yaml` to `0.26.1`. This ensures that Helm installations targeting the main branch default to using the newly released chart artifact.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 817 bytes
- [#708](https://github.com/kairos-io/AuroraBoot/pull/708) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - helm/helm v3.21.2..v3.21.3 (PR body): compare v3.21.2...v3.21.3 ✓ 5490 bytes
    - context: 64643 bytes
- [#709](https://github.com/kairos-io/AuroraBoot/pull/709) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/kairos-io/kairos-sdk 0.24.0→0.25.1: compare v0.24.0...v0.25.1 ✓ 40000 bytes
    - github.com/go-logr/logr 1.4.3→1.4.4: compare v1.4.3...v1.4.4 ✓ 40000 bytes
    - context: 90757 bytes
- [#710](https://github.com/kairos-io/AuroraBoot/pull/710) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/stmcginnis/gofish 0.23.0→0.24.0: compare v0.23.0...v0.24.0 ✓ 33308 bytes
    - context: 37982 bytes
- [#712](https://github.com/kairos-io/AuroraBoot/pull/712) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/foxboron/sbctl 0.0.0-20250917190250-6b8ed8715652→0.0.0-20260802183653-a7168106e003: compare 6b8ed8715652...a7168106e003 ✓ 18927 bytes
    - context: 21798 bytes
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
    - k8s.io/api 0.24.0→0.36.3: compare v0.24.0...v0.36.3 ✓ 40000 bytes
    - context: 126088 bytes
- [#23](https://github.com/kairos-io/entangle-proxy/pull/23) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/checkout v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/checkout v6.0.3..v7.0.0 (PR body): compare v6.0.3...v7.0.0 ✓ 40000 bytes
    - context: 63254 bytes
- [#25](https://github.com/kairos-io/entangle-proxy/pull/25) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/go-logr/logr 1.4.3→1.4.4: compare v1.4.3...v1.4.4 ✓ 40000 bytes
    - context: 44091 bytes
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
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - context: 42046 bytes
**[kairos-io/go-ukify](https://github.com/kairos-io/go-ukify)**

- [#59](https://github.com/kairos-io/go-ukify/pull/59) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - securego/gosec v2.27.1..v2.28.0 (PR body): compare v2.27.1...v2.28.0 ✓ 40000 bytes
    - context: 44536 bytes
- [#60](https://github.com/kairos-io/go-ukify/pull/60) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - context: 42433 bytes
- [#61](https://github.com/kairos-io/go-ukify/pull/61) — ✅ **good** — The update is a standard version bump to a newer release (v1.6.8) of the dependency. The upstream changelog indicates this release addresses important structural changes, specifically restoring the module path and updating Go version policies, which is a necessary maintenance step. There are no immediate red flags suggesting a security risk.
  ↳ This PR updates the dependency `github.com/ThalesGroup/crypto11` from version v1.6.2 to v1.6.8. This update incorporates a repository migration to the `github.com/eclipse-keypont/crypto11` module path and includes fixes related to module path restoration and Go version policy.
    - github.com/ThalesGroup/crypto11 1.6.2→1.6.8: compare v1.6.2...v1.6.8 ✓ 2785 bytes
    - ThalesGroup/crypto11 v1.6.7..v1.6.8 (PR body): compare v1.6.7...v1.6.8 ✓ 2070 bytes
    - eclipse-keypont/crypto11 v1.6.5..v1.6.8 (PR body): compare v1.6.5...v1.6.8 ✓ 2070 bytes
    - ThalesGroup/crypto11 v1.6.6..v1.6.7 (PR body): compare v1.6.6...v1.6.7 ✓ 617 bytes
    - ThalesGroup/crypto11 v1.6.5..v1.6.6 (PR body): compare v1.6.5...v1.6.6 ✓ 936 bytes
    - context: 14589 bytes
**[kairos-io/kairos](https://github.com/kairos-io/kairos)**

- [#4229](https://github.com/kairos-io/kairos/pull/4229) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/login-action c99871dec2022cc055c062a10cc1a1310835ceb4..af1e73f918a031802d376d3c8bbc3fe56130a9b0 (PR body): compare c99871dec2022cc055c062a10cc1a1310835ceb4...af1e73f918a031802d376d3c8bbc3fe56130a9b0 ✓ 40000 bytes
    - context: 43062 bytes
- [#4234](https://github.com/kairos-io/kairos/pull/4234) — ✅ **good** — This is a routine minor version bump from a trusted dependency, and the changes are documented in the upstream release notes. The update is applied consistently across all relevant workflow files, suggesting a safe and necessary maintenance update.
  ↳ This PR updates the dependency `kairos-io/kairos-factory-action` from v1.1.3 to v1.2.0. This update incorporates changes from the upstream release, including updates to `actions/checkout` and `github/codeql-action`. The change is applied across multiple CI/CD workflow files.
    - kairos-io/kairos-factory-action v1.1.3..v1.2.0 (PR body): compare v1.1.3...v1.2.0 ✓ 7319 bytes
    - context: 15709 bytes
- [#4256](https://github.com/kairos-io/kairos/pull/4256) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/checkout v7.0.0..v7.0.1 (PR body): compare v7.0.0...v7.0.1 ✓ 40000 bytes
    - context: 48520 bytes
- [#4259](https://github.com/kairos-io/kairos/pull/4259) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - aws-actions/configure-aws-credentials 517a711dbcd0e402f90c77e7e2f81e849156e31d..e6de054238d6b7531b4efff3b6587d9aade6a06c (PR body): compare 517a711dbcd0e402f90c77e7e2f81e849156e31d...e6de054238d6b7531b4efff3b6587d9aade6a06c ✓ 40000 bytes
    - context: 42305 bytes
- [#4270](https://github.com/kairos-io/kairos/pull/4270) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/stale v11.0.0..v11.0.0 (PR body): compare v11.0.0...v11.0.0 failed/empty (no upstream diff)
    - actions/stale v10..v11.0.0 (PR body): compare v10...v11.0.0 ✓ 40000 bytes
    - actions/stale v10.4.0..v11.0.0 (PR body): compare v10.4.0...v11.0.0 ✓ 40000 bytes
    - context: 82981 bytes
- [#4272](https://github.com/kairos-io/kairos/pull/4272) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - kairos-io/kairos-init v0.16.1..v0.16.2 (PR body): compare v0.16.1...v0.16.2 ✓ 25265 bytes
    - context: 29768 bytes
- [#4273](https://github.com/kairos-io/kairos/pull/4273) — ✅ **good** — This pull request is a minor version bump for a specific Docker image tag. As this is a dependency update within a known project, and there are no indications of malicious code or breaking changes, it is safe to auto-approve.
  ↳ This PR updates the Docker image tag for `quay.io/kairos/auroraboot` from v0.25.2 to v0.26.1 across multiple GitHub Actions workflows. This is a standard dependency maintenance update.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 3582 bytes
- [#4277](https://github.com/kairos-io/kairos/pull/4277) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - azure/login 532459ea530d8321f2fb9bb10d1e0bcf23869a43..f5d393ae46f8fde4be8b75f32e3fc50e654ad0ca (PR body): compare 532459ea530d8321f2fb9bb10d1e0bcf23869a43...f5d393ae46f8fde4be8b75f32e3fc50e654ad0ca ✓ 40000 bytes
    - context: 42094 bytes
**[kairos-io/kairos-init](https://github.com/kairos-io/kairos-init)**

- [#422](https://github.com/kairos-io/kairos-init/pull/422) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/kairos-io/kairos-sdk 0.25.0→0.25.1: compare v0.25.0...v0.25.1 ✓ 37726 bytes
    - github.com/docker/cli 29.5.3+incompatible→29.6.2+incompatible: compare v29.5.3+incompatible...v29.6.2+incompatible failed/empty (no upstream diff)
    - github.com/go-logr/logr 1.4.3→1.4.4: compare v1.4.3...v1.4.4 ✓ 40000 bytes
    - context: 92488 bytes
**[kairos-io/kairos-installer](https://github.com/kairos-io/kairos-installer)**

- [#19](https://github.com/kairos-io/kairos-installer/pull/19) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/kairos-io/kairos-sdk 0.24.0→0.25.0: compare v0.24.0...v0.25.0 ✓ 26995 bytes
    - github.com/docker/cli 29.6.2+incompatible→29.7.1+incompatible: compare v29.6.2+incompatible...v29.7.1+incompatible failed/empty (no upstream diff)
    - github.com/docker/go-connections 0.8.0→0.8.1: compare v0.8.0...v0.8.1 ✓ 1617 bytes
    - github.com/google/go-containerregistry 0.21.7→0.21.8: compare v0.21.7...v0.21.8 ✓ 40000 bytes
    - context: 82356 bytes
**[kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator)**

- [#153](https://github.com/kairos-io/kairos-operator/pull/153) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - docker/login-action abd2ef45e78c5afb21d64d4ca52ee8550d9572c7..dbcb813823bdd20940b903addbd779551569679f (PR body): compare abd2ef45e78c5afb21d64d4ca52ee8550d9572c7...dbcb813823bdd20940b903addbd779551569679f ✓ 40000 bytes
    - context: 43269 bytes
- [#156](https://github.com/kairos-io/kairos-operator/pull/156) — ✅ **good** — This is a routine maintenance update to change a specific digest for a base image. It does not introduce new code, change dependencies, or alter the build logic in a way that introduces security risks. This type of digest update is safe and recommended for ensuring build reproducibility.
  ↳ This PR updates the specific SHA256 digest for the `docker.io/golang:1.26.5` base image across the Dockerfiles and Dockerfile.node-labeler. This change ensures the build uses a fixed, known version of the image, which is a standard maintenance practice.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1972 bytes
**[kairos-io/kairos-sdk](https://github.com/kairos-io/kairos-sdk)**

- [#827](https://github.com/kairos-io/kairos-sdk/pull/827) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - context: 83033 bytes
**[kairos-io/netboot](https://github.com/kairos-io/netboot)**

- [#46](https://github.com/kairos-io/netboot/pull/46) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.53.0→0.54.0: compare v0.53.0...v0.54.0 ✓ 40000 bytes
    - golang.org/x/net 0.56.0→0.57.0: compare v0.56.0...v0.57.0 ✓ 40000 bytes
    - context: 83965 bytes
- [#47](https://github.com/kairos-io/netboot/pull/47) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - context: 83203 bytes
- [#48](https://github.com/kairos-io/netboot/pull/48) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - golang.org/x/crypto 0.53.0→0.54.0: compare v0.53.0...v0.54.0 ✓ 40000 bytes
    - golang.org/x/sys 0.46.0→0.47.0: compare v0.46.0...v0.47.0 ✓ 33531 bytes
    - context: 76973 bytes
**[mauromorales/xpasswd](https://github.com/mauromorales/xpasswd)**

- [#53](https://github.com/mauromorales/xpasswd/pull/53) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - context: 42399 bytes
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
- [#1054](https://github.com/mudler/edgevpn/pull/1054) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - urfave/cli v3.10.0..v3.10.1 (PR body): compare v3.10.0...v3.10.1 ✓ 17319 bytes
    - urfave/cli v3.9.1..v3.10.0 (PR body): compare v3.9.1...v3.10.0 ✓ 40000 bytes
    - context: 100504 bytes
- [#1056](https://github.com/mudler/edgevpn/pull/1056) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/libp2p/go-libp2p-pubsub 0.16.0→0.17.0: compare v0.16.0...v0.17.0 ✓ 40000 bytes
    - context: 45662 bytes
- [#1057](https://github.com/mudler/edgevpn/pull/1057) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - actions/setup-go v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/setup-go v6..v7.0.0 (PR body): compare v6...v7.0.0 ✓ 40000 bytes
    - actions/setup-go v6.5.0..v7.0.0 (PR body): compare v6.5.0...v7.0.0 ✓ 40000 bytes
    - context: 84598 bytes
- [#1059](https://github.com/mudler/edgevpn/pull/1059) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/libp2p/go-libp2p-kad-dht 0.41.0→0.42.1: compare v0.41.0...v0.42.1 ✓ 40000 bytes
    - github.com/ipfs/boxo 0.39.0→0.41.0: compare v0.39.0...v0.41.0 ✓ 40000 bytes
    - context: 99296 bytes
- [#1060](https://github.com/mudler/edgevpn/pull/1060) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/labstack/echo/v5 5.3.0→5.3.1: compare v5.3.0...v5.3.1 ✓ 35382 bytes
    - golang.org/x/crypto 0.53.0→0.54.0: compare v0.53.0...v0.54.0 ✓ 40000 bytes
    - context: 86707 bytes
- [#1061](https://github.com/mudler/edgevpn/pull/1061) — ✅ **good** — This change is a routine dependency update, specifically updating the digest of `go-libp2p-pubsub`. There are no apparent security risks introduced by this version bump, and it aligns with standard dependency maintenance practices.
  ↳ This PR updates the dependency `github.com/mudler/go-libp2p-pubsub` by changing its digest from `205ded1` to `2a31b5e`. This is a routine maintenance update to ensure the project uses the latest version of this library.
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

- [#320](https://github.com/mudler/yip/pull/320) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/go-git/go-git/v5 5.19.1→5.19.2: compare v5.19.1...v5.19.2 ✓ 40000 bytes
    - context: 44928 bytes
- [#322](https://github.com/mudler/yip/pull/322) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/google/go-containerregistry 0.21.7→0.21.9: compare v0.21.7...v0.21.9 ✓ 40000 bytes
    - github.com/docker/cli 29.5.3+incompatible→29.6.2+incompatible: compare v29.5.3+incompatible...v29.6.2+incompatible failed/empty (no upstream diff)
    - github.com/klauspost/compress 1.18.6→1.19.1: compare v1.18.6...v1.19.1 ✓ 40000 bytes
    - context: 98197 bytes

