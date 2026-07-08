# Kairos Security Dashboard

_Updated 2026-07-08._

🌐 **[Live dashboard](https://kairos-io.github.io/security/)** — the published board with clickable links.

## 📋 This run

- **Scanned:** 28 repos (1 skipped)
- **Findings:** 61 (3 critical / 29 high / 24 medium / 3 low / 2 unknown)
- **CVE-related PRs:** 0
- **Remediation:** 0 open · 0 superseded · 0 merged · 0 need-human
- **Why:** 61 finding(s); 0 PR(s) open.

> The immediate focus should be on the three critical vulnerabilities found in the openssl-fips package (F1, F2, F3). Additionally, high-severity issues in the core openssl and glib libraries (F4, F6, F7, F8) require urgent attention to mitigate potential system compromise.

## 🔥 Focus now

- [CVE-2026-31789](https://osv.dev/vulnerability/ALPINE-CVE-2026-31789) — Critical vulnerability in openssl-fips via CVE-2026-31789.
- [CVE-2024-5535](https://osv.dev/vulnerability/ALPINE-CVE-2024-5535) — Critical vulnerability in openssl-fips via CVE-2024-5535.
- [CVE-2026-34182](https://osv.dev/vulnerability/ALPINE-CVE-2026-34182) — Critical vulnerability in openssl-fips via CVE-2026-34182.
- [CVE-2023-4807](https://osv.dev/vulnerability/ALPINE-CVE-2023-4807) — High severity vulnerability in openssl via CVE-2023-4807.
- [CVE-2026-9076](https://osv.dev/vulnerability/ALPINE-CVE-2026-9076) — High severity vulnerability in openssl-fips via CVE-2026-9076.
- [CVE-2021-27219](https://osv.dev/vulnerability/ALPINE-CVE-2021-27219) — High severity vulnerability in glib via CVE-2021-27219.
- [CVE-2026-45447](https://osv.dev/vulnerability/ALPINE-CVE-2026-45447) — High severity vulnerability in openssl-fips via CVE-2026-45447.
- [CVE-2025-9230](https://osv.dev/vulnerability/ALPINE-CVE-2025-9230) — High severity vulnerability in openssl-fips via CVE-2025-9230.

## 🌊 Waterfall fronts

_None._

## 📦 Per-repo findings

| Repo | Critical | High | Medium | Total | Status |
|---|---|---|---|---|---|
| [kairos-io/hadron](https://github.com/kairos-io/hadron) | 3 | 29 | 24 | 56 | ok |
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
| [kairos-io/kairos-must-burn](https://github.com/kairos-io/kairos-must-burn) | 0 | 0 | 0 | 0 | skipped: not source-scannable |
| [kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kairos-sdk](https://github.com/kairos-io/kairos-sdk) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kcrypt](https://github.com/kairos-io/kcrypt) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/kcrypt-discovery-challenger](https://github.com/kairos-io/kcrypt-discovery-challenger) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/netboot](https://github.com/kairos-io/netboot) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/provider-kairos](https://github.com/kairos-io/provider-kairos) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/provider-kubernetes](https://github.com/kairos-io/provider-kubernetes) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/simple-mdns-server](https://github.com/kairos-io/simple-mdns-server) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [kairos-io/tpm-helpers](https://github.com/kairos-io/tpm-helpers) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mauromorales/xpasswd](https://github.com/mauromorales/xpasswd) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/edgevpn](https://github.com/mudler/edgevpn) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/entities](https://github.com/mudler/entities) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/go-pluggable](https://github.com/mudler/go-pluggable) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |
| [mudler/yip](https://github.com/mudler/yip) | 0 | 0 | 0 | 0 | clean (no crit/high/med) |

## 🧩 Hadron component CVEs

| Package | Current | Fixed | Severity | CVE |
|---|---|---|---|---|
| openssl-fips | 3.1.2 | 3.5.6 | critical | [CVE-2026-31789](https://osv.dev/vulnerability/ALPINE-CVE-2026-31789) |
| openssl-fips | 3.1.2 | 3.3.1 | critical | [CVE-2024-5535](https://osv.dev/vulnerability/ALPINE-CVE-2024-5535) |
| openssl-fips | 3.1.2 | 3.5.7 | critical | [CVE-2026-34182](https://osv.dev/vulnerability/ALPINE-CVE-2026-34182) |
| glib | 2.86.2 | 2.66.6 | high | [CVE-2021-27219](https://osv.dev/vulnerability/ALPINE-CVE-2021-27219) |
| libxml2 | 2.15.3 | 2.13.8 | high | [CVE-2025-32415](https://osv.dev/vulnerability/ALPINE-CVE-2025-32415) |
| libxml2 | 2.15.3 | 2.13.8 | high | [CVE-2025-32414](https://osv.dev/vulnerability/ALPINE-CVE-2025-32414) |
| openssl | 3.6.3 | 0 | high | [CVE-2022-2068](https://osv.dev/vulnerability/ALPINE-CVE-2022-2068) |
| openssl | 3.6.3 | 0 | high | [CVE-2023-4807](https://osv.dev/vulnerability/ALPINE-CVE-2023-4807) |
| openssl | 3.6.3 | 0 | high | [CVE-2022-1292](https://osv.dev/vulnerability/ALPINE-CVE-2022-1292) |
| openssl-fips | 3.1.2 | 0 | high | [CVE-2022-2068](https://osv.dev/vulnerability/ALPINE-CVE-2022-2068) |
| openssl-fips | 3.1.2 | 3.5.5 | high | [CVE-2025-69421](https://osv.dev/vulnerability/ALPINE-CVE-2025-69421) |
| openssl-fips | 3.1.2 | 3.3.0 | high | [CVE-2024-4741](https://osv.dev/vulnerability/ALPINE-CVE-2024-4741) |
| openssl-fips | 3.1.2 | 3.5.4 | high | [CVE-2025-9230](https://osv.dev/vulnerability/ALPINE-CVE-2025-9230) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-45445](https://osv.dev/vulnerability/ALPINE-CVE-2026-45445) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-7383](https://osv.dev/vulnerability/ALPINE-CVE-2026-7383) |
| openssl-fips | 3.1.2 | 3.5.6 | high | [CVE-2026-31790](https://osv.dev/vulnerability/ALPINE-CVE-2026-31790) |
| openssl-fips | 3.1.2 | 3.5.5 | high | [CVE-2025-69419](https://osv.dev/vulnerability/ALPINE-CVE-2025-69419) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-45447](https://osv.dev/vulnerability/ALPINE-CVE-2026-45447) |
| openssl-fips | 3.1.2 | 0 | high | [CVE-2022-1292](https://osv.dev/vulnerability/ALPINE-CVE-2022-1292) |
| openssl-fips | 3.1.2 | 0 | high | [CVE-2023-4807](https://osv.dev/vulnerability/ALPINE-CVE-2023-4807) |
| openssl-fips | 3.1.2 | 3.5.5 | high | [CVE-2025-69420](https://osv.dev/vulnerability/ALPINE-CVE-2025-69420) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-9076](https://osv.dev/vulnerability/ALPINE-CVE-2026-9076) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-34180](https://osv.dev/vulnerability/ALPINE-CVE-2026-34180) |
| openssl-fips | 3.1.2 | 3.5.6 | high | [CVE-2026-28390](https://osv.dev/vulnerability/ALPINE-CVE-2026-28390) |
| openssl-fips | 3.1.2 | 3.5.6 | high | [CVE-2026-28389](https://osv.dev/vulnerability/ALPINE-CVE-2026-28389) |
| openssl-fips | 3.1.2 | 3.5.6 | high | [CVE-2026-28388](https://osv.dev/vulnerability/ALPINE-CVE-2026-28388) |
| openssl-fips | 3.1.2 | 3.1.4 | high | [CVE-2023-5363](https://osv.dev/vulnerability/ALPINE-CVE-2023-5363) |
| openssl-fips | 3.1.2 | 3.3.2 | high | [CVE-2024-6119](https://osv.dev/vulnerability/ALPINE-CVE-2024-6119) |
| openssl-fips | 3.1.2 | 3.5.5 | high | [CVE-2025-15467](https://osv.dev/vulnerability/ALPINE-CVE-2025-15467) |
| openssl-fips | 3.1.2 | 3.5.6 | high | [CVE-2026-28387](https://osv.dev/vulnerability/ALPINE-CVE-2026-28387) |
| rsync | 3.4.4 | 0 | high | [CVE-2020-14387](https://osv.dev/vulnerability/ALPINE-CVE-2020-14387) |
| sqlite3 | 3.53.3 | 0 | high | [CVE-2022-35737](https://osv.dev/vulnerability/ALPINE-CVE-2022-35737) |
| busybox | 1.37.0 | 0 | medium | [CVE-2021-42376](https://osv.dev/vulnerability/ALPINE-CVE-2021-42376) |
| curl | 8.21.0 | 0 | medium | [CVE-2021-22897](https://osv.dev/vulnerability/ALPINE-CVE-2021-22897) |
| gcc | 15.3.0 | 13.2.1_git20231014 | medium | [CVE-2023-4039](https://osv.dev/vulnerability/ALPINE-CVE-2023-4039) |
| openssl | 3.6.3 | 0 | medium | [CVE-2023-0466](https://osv.dev/vulnerability/ALPINE-CVE-2023-0466) |
| openssl-fips | 3.1.2 | 0 | medium | [CVE-2023-0466](https://osv.dev/vulnerability/ALPINE-CVE-2023-0466) |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2024-0727](https://osv.dev/vulnerability/ALPINE-CVE-2024-0727) |
| openssl-fips | 3.1.2 | 3.5.4 | medium | [CVE-2025-9231](https://osv.dev/vulnerability/ALPINE-CVE-2025-9231) |
| openssl-fips | 3.1.2 | 3.3.3 | medium | [CVE-2024-12797](https://osv.dev/vulnerability/ALPINE-CVE-2024-12797) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2026-22796](https://osv.dev/vulnerability/ALPINE-CVE-2026-22796) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2025-68160](https://osv.dev/vulnerability/ALPINE-CVE-2025-68160) |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-6237](https://osv.dev/vulnerability/ALPINE-CVE-2023-6237) |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-6129](https://osv.dev/vulnerability/ALPINE-CVE-2023-6129) |
| openssl-fips | 3.1.2 | 3.5.1 | medium | [CVE-2025-4575](https://osv.dev/vulnerability/ALPINE-CVE-2025-4575) |
| openssl-fips | 3.1.2 | 3.5.4 | medium | [CVE-2025-9232](https://osv.dev/vulnerability/ALPINE-CVE-2025-9232) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2025-69418](https://osv.dev/vulnerability/ALPINE-CVE-2025-69418) |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-5678](https://osv.dev/vulnerability/ALPINE-CVE-2023-5678) |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-45446](https://osv.dev/vulnerability/ALPINE-CVE-2026-45446) |
| openssl-fips | 3.1.2 | 3.3.2 | medium | [CVE-2024-9143](https://osv.dev/vulnerability/ALPINE-CVE-2024-9143) |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-42766](https://osv.dev/vulnerability/ALPINE-CVE-2026-42766) |
| openssl-fips | 3.1.2 | 3.2.1 | medium | [CVE-2024-2511](https://osv.dev/vulnerability/ALPINE-CVE-2024-2511) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2026-22795](https://osv.dev/vulnerability/ALPINE-CVE-2026-22795) |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-42767](https://osv.dev/vulnerability/ALPINE-CVE-2026-42767) |
| openssl-fips | 3.1.2 | 3.3.0 | medium | [CVE-2024-4603](https://osv.dev/vulnerability/ALPINE-CVE-2024-4603) |
| openssl-fips | 3.1.2 | 3.3.2 | medium | [CVE-2024-13176](https://osv.dev/vulnerability/ALPINE-CVE-2024-13176) |
| busybox | 1.37.0 | 1.37.0 | low | [CVE-2024-58251](https://osv.dev/vulnerability/ALPINE-CVE-2024-58251) |
| busybox | 1.37.0 | 1.37.0 | low | [CVE-2025-46394](https://osv.dev/vulnerability/ALPINE-CVE-2025-46394) |
| openssl-fips | 3.1.2 | 3.5.7 | low | [CVE-2026-42770](https://osv.dev/vulnerability/ALPINE-CVE-2026-42770) |
| perl | 5.42.2 | 5.26.3 | unknown | [CVE-2018-18311](https://osv.dev/vulnerability/ALPINE-CVE-2018-18311) |
| perl | 5.42.2 | 5.26.3 | unknown | [CVE-2018-18312](https://osv.dev/vulnerability/ALPINE-CVE-2018-18312) |

## 📋 Open PRs

_None._

## 🤖 Bot PR ledger

| Repo | Bump | Kind | Source | State | PR |
|---|---|---|---|---|---|
| [mudler/edgevpn](https://github.com/mudler/edgevpn) | golang.org/x/net@0.33.0 | direct | ksec | error | — |

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
- [#591](https://github.com/kairos-io/AuroraBoot/pull/591) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - microsoft/TypeScript v6.0.2..v6.0.3 (PR body): compare v6.0.2...v6.0.3 ✓ 40000 bytes
    - microsoft/TypeScript v5.9.3..v6.0.2 (PR body): compare v5.9.3...v6.0.2 ✓ 40000 bytes
    - context: 85089 bytes
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
**[kairos-io/cluster-api-provider-kairos](https://github.com/kairos-io/cluster-api-provider-kairos)**

- [#38](https://github.com/kairos-io/cluster-api-provider-kairos/pull/38) — ✅ **good** — This pull request is a routine dependency update for golang.org/x/oauth2. Updating to a newer version is standard practice and generally safe, as it addresses potential minor issues or security patches without introducing significant risk.
**[kairos-io/entangle](https://github.com/kairos-io/entangle)**

- [#10](https://github.com/kairos-io/entangle/pull/10) — ✅ **good** — This pull request only updates several indirect dependencies to newer versions. These types of dependency bumps are routine maintenance and do not introduce new security risks. The changes appear safe to merge automatically.
**[kairos-io/entangle-proxy](https://github.com/kairos-io/entangle-proxy)**

- [#5](https://github.com/kairos-io/entangle-proxy/pull/5) — ✅ **good** — The changes are primarily feature additions (new matchers, plugin system) and internal refactoring aimed at improving robustness (cycle detection, iteration limits). There are no obvious security vulnerabilities introduced by these changes. The dependency update is a standard version bump, and the architectural changes in the SSH client appear to be performance/concurrency improvements. Therefore, this PR is safe to auto-approve.
  ↳ This PR updates the `gomega` dependency to v1.42.1, introduces new matchers (`BeASlice`, `BeAnArray`), and adds a plugin system for Claude Code skills. It also includes significant internal refactoring to improve robustness, such as cycle detection in formatting and iteration limits in PKCS#12 decryption. Finally, it implements request pipelining for the SSH agent client to improve concurrency.
    - github.com/onsi/gomega 1.40.0→1.42.1: compare v1.40.0...v1.42.1 ✓ 40000 bytes
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - context: 88194 bytes
- [#6](https://github.com/kairos-io/entangle-proxy/pull/6) — ✅ **good** — This pull request updates the dependency sigs.k8s.io/controller-runtime to version v0.24.1. This is a routine dependency update to a newer version, which is generally safe and necessary for maintaining security and compatibility.
- [#14](https://github.com/kairos-io/entangle-proxy/pull/14) — ✅ **good** — This pull request primarily updates several dependencies to newer versions, including core packages like `golang.org/x` and `google.golang.org/protobuf`. Updating dependencies is a crucial security practice to ensure that known vulnerabilities are patched. The changes appear to be dependency hygiene improvements and do not introduce any obvious security risks.
- [#18](https://github.com/kairos-io/entangle-proxy/pull/18) — ✅ **good** — This is a standard major version upgrade for a core dependency, which generally includes important security patches and feature improvements. The changes are comprehensive and align with maintaining up-to-date tooling. Therefore, it is safe to auto-approve.
  ↳ This PR updates the `docker/build-push-action` dependency from v2 to v7, which includes numerous underlying dependency updates. It also modifies the CI/CD workflow files to use the corresponding v7.0.0 and v4.1.0 versions of related Docker actions, ensuring compatibility with the new action.
    - docker/build-push-action v7.2.0..v7.3.0 (PR body): compare v7.2.0...v7.3.0 ✓ 40000 bytes
    - context: 83719 bytes
- [#20](https://github.com/kairos-io/entangle-proxy/pull/20) — ⚠️ **needs_human_verification** — Upgrading core Kubernetes dependencies is generally beneficial for security and feature updates. However, this is a major version jump (v0.24.0 to v0.36.2) across multiple critical libraries. A human review is required to ensure that no breaking changes have been introduced that could impact the application's functionality or security posture.
  ↳ This pull request updates three core Kubernetes dependencies (k8s.io/api, k8s.io/apimachinery, and k8s.io/client-go) from version v0.24.0 to v0.36.2. This involves significant version bumps across the board.
    - k8s.io/api 0.24.0→0.36.2: compare v0.24.0...v0.36.2 ✓ 40000 bytes
    - context: 126004 bytes
- [#23](https://github.com/kairos-io/entangle-proxy/pull/23) — ✅ **good** — This pull request updates the dependency for 'actions/checkout' from version v2 to v7 in two workflow files. This is a standard dependency update to a newer version, which is generally safe and beneficial for security and maintenance.
**[kairos-io/go-nodepair](https://github.com/kairos-io/go-nodepair)**

- [#27](https://github.com/kairos-io/go-nodepair/pull/27) — ✅ **good** — The changes involve updating several core dependencies across the project. The changelogs indicate that these updates include important security patches, such as restricting RSA key sizes in go-libp2p and fixing memory exhaustion attacks in quic-go. This is standard maintenance and security hygiene.
  ↳ This pull request updates several core dependencies, including go-libp2p, quic-go, golang.org/x/crypto, golang.org/x/image, golang.org/x/net, and google.golang.org/protobuf. The updates include critical security fixes, such as mitigating a DoS attack in go-libp2p and addressing memory exhaustion issues in quic-go.
- [#65](https://github.com/kairos-io/go-nodepair/pull/65) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - ipfs/go-log v2.9.1..v2.9.2 (PR body): compare v2.9.1...v2.9.2 ✓ 12110 bytes
    - ipfs/go-log v2.9.0..v2.9.1 (PR body): compare v2.9.0...v2.9.1 ✓ 2149 bytes
    - ipfs/go-log v2.8.2..v2.9.0 (PR body): compare v2.8.2...v2.9.0 ✓ 40000 bytes
    - context: 70615 bytes
**[kairos-io/hadron](https://github.com/kairos-io/hadron)**

- [#502](https://github.com/kairos-io/hadron/pull/502) — ✅ **good** — This change is necessary to resolve a build-time dependency issue where `libsystemd.so` requires `libucontext.so.1`. The implementation correctly adds a build stage for `libucontext` and injects the resulting library into the main build process, ensuring the build can complete successfully.
  ↳ The PR introduces a new build stage to download, compile, and install `libucontext`, which resolves a critical linker failure against `libsystemd.so` by making the required library available during the `fwupd` build.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 3403 bytes
**[kairos-io/kairos-must-burn](https://github.com/kairos-io/kairos-must-burn)**

- [#45](https://github.com/kairos-io/kairos-must-burn/pull/45) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - google/go-github v88.0.0..v89.0.0 (PR body): compare v88.0.0...v89.0.0 ✓ 40000 bytes
    - context: 54668 bytes
**[kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator)**

- [#136](https://github.com/kairos-io/kairos-operator/pull/136) — ✅ **good** — This is a patch update to a specific dependency. Updating dependencies to newer patch versions is a standard maintenance practice and is generally safe. There are no obvious security risks introduced by this minor version bump.
  ↳ This PR updates the Docker image tag for `docker.io/golang` from version `1.26.4` to `1.26.5` in the relevant Dockerfiles. This is a minor patch update to the dependency.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1872 bytes
**[kairos-io/kcrypt](https://github.com/kairos-io/kcrypt)**

- [#505](https://github.com/kairos-io/kcrypt/pull/505) — ✅ **good** — The change is a dependency update to a newer version of the Kairos SDK, which is generally safe. The accompanying code changes focus on internal refactoring and adding complex features like multipath device handling, which appear to be well-tested in the provided diff. This is safe to auto-approve.
  ↳ This PR updates the `github.com/kairos-io/kairos-sdk` dependency from v0.9.4 to v0.11.0. It also includes significant internal refactoring within the `ghw` package to introduce robust support for multipath devices and update internal logging and provider management structures.
    - github.com/kairos-io/kairos-sdk 0.9.4→0.11.0: compare v0.9.4...v0.11.0 ✓ 40000 bytes
    - context: 48550 bytes
- [#509](https://github.com/kairos-io/kcrypt/pull/509) — ✅ **good** — This is a necessary upgrade to a major dependency, which includes important security fixes mentioned in the release notes. Since the update is to a newer major version and addresses security concerns, it is safe to auto-approve.
  ↳ This PR updates the dependency github.com/docker/docker from version 27.5.1+incompatible to 28.0.0+incompatible. This upgrade incorporates new features and critical security fixes released in the Docker 28.0.0 release.
    - github.com/docker/docker 27.5.1+incompatible→28.0.0+incompatible: compare v27.5.1+incompatible...v28.0.0+incompatible failed: gh api: gh: Not Found (HTTP 404) (no upstream diff)
    - context: 46086 bytes
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
**[kairos-io/simple-mdns-server](https://github.com/kairos-io/simple-mdns-server)**

- [#4](https://github.com/kairos-io/simple-mdns-server/pull/4) — ✅ **good** — This is a standard dependency maintenance update performed by Dependabot. The changes involve updating core packages like `x/net` and `x/sys`, which are necessary for project health. Since this is an automated bump and no immediate security risks are apparent from the diffs, it is safe to auto-approve.
  ↳ This pull request updates two dependencies: `golang.org/x/net` to version 0.23.0 and `golang.org/x/sys` to version 0.18.0. The updates include significant changes to networking code, context handling, and low-level CPU/system interaction logic.
    - golang.org/x/net 0.0.0-20210410081132-afb366fc7cd1→0.23.0: compare afb366fc7cd1...v0.23.0 ✓ 40000 bytes
    - golang.org/x/sys 0.0.0-20210330210617-4fbd30eecc44→0.18.0: compare 4fbd30eecc44...v0.18.0 ✓ 40000 bytes
    - context: 85299 bytes
**[kairos-io/tpm-helpers](https://github.com/kairos-io/tpm-helpers)**

- [#6](https://github.com/kairos-io/tpm-helpers/pull/6) — ⚠️ **needs_human_verification** — review endpoint unreachable: Post "http://localhost:8080/v1/chat/completions": context deadline exceeded
    - golang.org/x/crypto 0.0.0-20220722155217-630584e8d5aa→0.17.0: compare 630584e8d5aa...v0.17.0 ✓ 40000 bytes
    - golang.org/x/net 0.0.0-20220722155237-a158d28d115b→0.10.0: compare a158d28d115b...v0.10.0 ✓ 40000 bytes
    - context: 90128 bytes
**[mudler/edgevpn](https://github.com/mudler/edgevpn)**

- [#804](https://github.com/mudler/edgevpn/pull/804) — ✅ **good** — This pull request appears to be a necessary and comprehensive migration to upgrade the `iplib` library to version 2.0.5. The changes reflect the breaking changes detailed in the v2 release notes, specifically the transition to `uint128` for IPv6 handling. The diffs suggest all necessary package imports and internal logic have been updated to accommodate this major version change.
  ↳ The PR updates the `github.com/c-robinson/iplib` dependency from v1.0.8 to v2.0.5. This migration involves adopting the new `uint128` type for IPv6 functions and updating related package imports and logic across the codebase.
    - c-robinson/iplib v2.0.4..v2.0.5 (PR body): compare v2.0.4...v2.0.5 ✓ 6378 bytes
    - c-robinson/iplib v2.0.3..v2.0.4 (PR body): compare v2.0.3...v2.0.4 ✓ 3273 bytes
    - c-robinson/iplib v2.0.2..v2.0.3 (PR body): compare v2.0.2...v2.0.3 ✓ 9999 bytes
    - c-robinson/iplib v2.0.1..v2.0.2 (PR body): compare v2.0.1...v2.0.2 ✓ 15662 bytes
    - c-robinson/iplib v2.0.0..v2.0.1 (PR body): compare v2.0.0...v2.0.1 ✓ 1844 bytes
    - context: 44540 bytes
- [#805](https://github.com/mudler/edgevpn/pull/805) — ✅ **good** — This is a standard dependency update to a newer version of the YAML library. Upgrading dependencies is crucial for security and accessing bug fixes and new features. The changes appear to be standard migration steps for moving from v2 to v3, making this change safe to auto-approve.
  ↳ This PR updates the `gopkg.in/yaml.v2` dependency from version v2.4.0 to v3.0.1. This upgrade incorporates changes from v3.0.0, including updates to the parser logic, documentation, and licensing information.
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
- [#914](https://github.com/mudler/edgevpn/pull/914) — ✅ **good** — This pull request updates a dependency to a version that directly addresses a reported security vulnerability (GHSA-vfp3-v2gw-7wfq). The change hardens the static file serving mechanism by disabling path unescaping by default, which mitigates the risk of encoded path separators bypassing route-level access controls. This is a positive security enhancement.
  ↳ This PR updates the dependency `github.com/labstack/echo/v4` to version v4.15.4, which includes a security fix. This fix prevents encoded path separators from bypassing route-level middleware when serving static files by disabling path unescaping by default.
    - github.com/labstack/echo/v4 4.15.2→4.15.4: compare v4.15.2...v4.15.4 ✓ 30288 bytes
    - github.com/mattn/go-colorable 0.1.14→0.1.15: compare v0.1.14...v0.1.15 ✓ 5234 bytes
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - context: 90621 bytes
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
- [#1001](https://github.com/mudler/edgevpn/pull/1001) — ⚠️ **needs_human_verification** — review endpoint unreachable: Post "http://localhost:8080/v1/chat/completions": context deadline exceeded
    - FortAwesome/Font-Awesome 7.0.0..7.0.1 (PR body): compare 7.0.0...7.0.1 ✓ 40000 bytes
    - FortAwesome/Font-Awesome 6.7.2..7.0.0 (PR body): compare 6.7.2...7.0.0 failed/empty (no upstream diff)
    - FortAwesome/Font-Awesome 6.7.1..6.7.2 (PR body): compare 6.7.1...6.7.2 ✓ 40000 bytes
    - context: 83524 bytes
- [#1006](https://github.com/mudler/edgevpn/pull/1006) — ✅ **good** — The upgrade is to a newer minor version (4.15.1) which includes security enhancements, such as the new CSRF middleware features detailed in the release notes. There are no immediate red flags or known critical vulnerabilities associated with this specific version jump. Therefore, this change is safe to auto-approve.
  ↳ This pull request updates the dependency `github.com/labstack/echo/v4` from version 4.13.3 to 4.15.1. This upgrade incorporates several enhancements, including improved CSRF protection features and minor internal fixes related to time comparison logic.
    - github.com/labstack/echo/v4 4.13.3→4.15.1: compare v4.13.3...v4.15.1 ✓ 40000 bytes
    - github.com/mattn/go-colorable 0.1.13→0.1.14: compare v0.1.13...v0.1.14 ✓ 6350 bytes
    - golang.org/x/time 0.12.0→0.14.0: compare v0.12.0...v0.14.0 ✓ 606 bytes
    - context: 76092 bytes
- [#1041](https://github.com/mudler/edgevpn/pull/1041) — ✅ **good** — This is a dependency upgrade to a newer version of a library, which is generally beneficial. Crucially, this update incorporates a security fix that mitigates a path traversal/ACL bypass vulnerability related to URL-encoded separators in static file serving. Therefore, this change is safe to auto-approve.
  ↳ The PR updates the `github.com/labstack/echo/v4` dependency to `v5.2.1`. This upgrade includes a security fix that changes the default behavior of static file serving to prevent encoded path separators from bypassing route-level access controls.
    - github.com/labstack/echo/v4 4.15.2→4.15.4: compare v4.15.2...v4.15.4 ✓ 30288 bytes
    - github.com/mattn/go-colorable 0.1.14→0.1.15: compare v0.1.14...v0.1.15 ✓ 5234 bytes
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - context: 108729 bytes
- [#1046](https://github.com/mudler/edgevpn/pull/1046) — ✅ **good** — This is a routine dependency update to a newer version of Autoprefixer. The changes primarily address compatibility fixes and minor logic adjustments, which are generally safe and beneficial for the project. There are no apparent security risks introduced by this version bump.
  ↳ The PR updates the `autoprefixer` dependency from v10.5.0 to v10.5.2. This update includes fixes for Firefox compatibility related to `-webkit-fill-available` and related adjustments in source code and test files.
    - postcss/autoprefixer 10.5.1..10.5.2 (PR body): compare 10.5.1...10.5.2 ✓ 2688 bytes
    - postcss/autoprefixer 10.5.0..10.5.1 (PR body): compare 10.5.0...10.5.1 ✓ 40000 bytes
    - context: 54014 bytes
- [#1054](https://github.com/mudler/edgevpn/pull/1054) — ⚠️ **needs_human_verification** — review endpoint returned HTTP 500
    - github.com/urfave/cli/v3 3.10.0→3.10.1: compare v3.10.0...v3.10.1 ✓ 17319 bytes
    - urfave/cli v3.9.1..v3.10.0 (PR body): compare v3.9.1...v3.10.0 ✓ 40000 bytes
    - context: 101293 bytes

