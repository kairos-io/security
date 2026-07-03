# Kairos Security Dashboard

_Updated 2026-07-03._

🌐 **[Live dashboard](https://kairos-io.github.io/security/)** — the published board with clickable links.

## 📋 This run

- **Scanned:** 28 repos (1 skipped)
- **Findings:** 54 (3 critical / 21 high / 26 medium / 4 low / 0 unknown)
- **CVE-related PRs:** 0
- **Remediation:** 0 open · 0 superseded · 0 merged · 0 need-human
- **Why:** 54 finding(s); 0 PR(s) open.

> The triage is focused on critical severity vulnerabilities found in the openssl-fips package within the kairos-io/hadron repository. Immediate attention is required to mitigate these critical risks before they can be exploited.

## 🔥 Focus now

- [CVE-2026-31789](https://osv.dev/vulnerability/ALPINE-CVE-2026-31789) — Critical vulnerability CVE-2026-31789 in openssl-fips.
- [CVE-2024-5535](https://osv.dev/vulnerability/ALPINE-CVE-2024-5535) — Critical vulnerability CVE-2024-5535 in openssl-fips.
- [CVE-2026-34182](https://osv.dev/vulnerability/ALPINE-CVE-2026-34182) — Critical vulnerability CVE-2026-34182 in openssl-fips.

## 🌊 Waterfall fronts

_None._

## 📦 Per-repo findings

| Repo | Critical | High | Medium | Total | Status |
|---|---|---|---|---|---|
| [kairos-io/hadron](https://github.com/kairos-io/hadron) | 3 | 21 | 26 | 50 | ok |
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
| openssl-fips | 3.1.2 | 3.5.7 | critical | [CVE-2026-34182](https://osv.dev/vulnerability/ALPINE-CVE-2026-34182) |
| openssl-fips | 3.1.2 | 3.3.1 | critical | [CVE-2024-5535](https://osv.dev/vulnerability/ALPINE-CVE-2024-5535) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-34183](https://osv.dev/vulnerability/ALPINE-CVE-2026-34183) |
| openssl-fips | 3.1.2 | 3.5.5 | high | [CVE-2025-69421](https://osv.dev/vulnerability/ALPINE-CVE-2025-69421) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-45445](https://osv.dev/vulnerability/ALPINE-CVE-2026-45445) |
| openssl-fips | 3.1.2 | 3.1.4 | high | [CVE-2023-5363](https://osv.dev/vulnerability/ALPINE-CVE-2023-5363) |
| openssl-fips | 3.1.2 | 3.5.6 | high | [CVE-2026-31790](https://osv.dev/vulnerability/ALPINE-CVE-2026-31790) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-45447](https://osv.dev/vulnerability/ALPINE-CVE-2026-45447) |
| openssl-fips | 3.1.2 | 3.5.5 | high | [CVE-2025-69419](https://osv.dev/vulnerability/ALPINE-CVE-2025-69419) |
| openssl-fips | 3.1.2 | 3.5.4 | high | [CVE-2025-9230](https://osv.dev/vulnerability/ALPINE-CVE-2025-9230) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-34181](https://osv.dev/vulnerability/ALPINE-CVE-2026-34181) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-9076](https://osv.dev/vulnerability/ALPINE-CVE-2026-9076) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-7383](https://osv.dev/vulnerability/ALPINE-CVE-2026-7383) |
| openssl-fips | 3.1.2 | 3.3.0 | high | [CVE-2024-4741](https://osv.dev/vulnerability/ALPINE-CVE-2024-4741) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-34180](https://osv.dev/vulnerability/ALPINE-CVE-2026-34180) |
| openssl-fips | 3.1.2 | 3.5.6 | high | [CVE-2026-28390](https://osv.dev/vulnerability/ALPINE-CVE-2026-28390) |
| openssl-fips | 3.1.2 | 3.5.6 | high | [CVE-2026-28389](https://osv.dev/vulnerability/ALPINE-CVE-2026-28389) |
| openssl-fips | 3.1.2 | 3.5.6 | high | [CVE-2026-28388](https://osv.dev/vulnerability/ALPINE-CVE-2026-28388) |
| openssl-fips | 3.1.2 | 3.5.5 | high | [CVE-2025-69420](https://osv.dev/vulnerability/ALPINE-CVE-2025-69420) |
| openssl-fips | 3.1.2 | 3.5.5 | high | [CVE-2025-15467](https://osv.dev/vulnerability/ALPINE-CVE-2025-15467) |
| openssl-fips | 3.1.2 | 3.5.6 | high | [CVE-2026-28387](https://osv.dev/vulnerability/ALPINE-CVE-2026-28387) |
| openssl-fips | 3.1.2 | 3.3.2 | high | [CVE-2024-6119](https://osv.dev/vulnerability/ALPINE-CVE-2024-6119) |
| openssl-fips | 3.1.2 | 3.5.7 | high | [CVE-2026-42764](https://osv.dev/vulnerability/ALPINE-CVE-2026-42764) |
| gcc | 15.3.0 | 13.2.1_git20231014 | medium | [CVE-2023-4039](https://osv.dev/vulnerability/ALPINE-CVE-2023-4039) |
| openssl-fips | 3.1.2 | 3.3.0 | medium | [CVE-2024-4603](https://osv.dev/vulnerability/ALPINE-CVE-2024-4603) |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2024-0727](https://osv.dev/vulnerability/ALPINE-CVE-2024-0727) |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-42766](https://osv.dev/vulnerability/ALPINE-CVE-2026-42766) |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-6237](https://osv.dev/vulnerability/ALPINE-CVE-2023-6237) |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-6129](https://osv.dev/vulnerability/ALPINE-CVE-2023-6129) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2025-15469](https://osv.dev/vulnerability/ALPINE-CVE-2025-15469) |
| openssl-fips | 3.1.2 | 3.5.4 | medium | [CVE-2025-9232](https://osv.dev/vulnerability/ALPINE-CVE-2025-9232) |
| openssl-fips | 3.1.2 | 3.5.6 | medium | [CVE-2026-2673](https://osv.dev/vulnerability/ALPINE-CVE-2026-2673) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2025-15468](https://osv.dev/vulnerability/ALPINE-CVE-2025-15468) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2025-68160](https://osv.dev/vulnerability/ALPINE-CVE-2025-68160) |
| openssl-fips | 3.1.2 | 3.3.2 | medium | [CVE-2024-9143](https://osv.dev/vulnerability/ALPINE-CVE-2024-9143) |
| openssl-fips | 3.1.2 | 3.1.4 | medium | [CVE-2023-5678](https://osv.dev/vulnerability/ALPINE-CVE-2023-5678) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2025-11187](https://osv.dev/vulnerability/ALPINE-CVE-2025-11187) |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-42769](https://osv.dev/vulnerability/ALPINE-CVE-2026-42769) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2026-22796](https://osv.dev/vulnerability/ALPINE-CVE-2026-22796) |
| openssl-fips | 3.1.2 | 3.5.4 | medium | [CVE-2025-9231](https://osv.dev/vulnerability/ALPINE-CVE-2025-9231) |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-42767](https://osv.dev/vulnerability/ALPINE-CVE-2026-42767) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2025-66199](https://osv.dev/vulnerability/ALPINE-CVE-2025-66199) |
| openssl-fips | 3.1.2 | 3.5.7 | medium | [CVE-2026-45446](https://osv.dev/vulnerability/ALPINE-CVE-2026-45446) |
| openssl-fips | 3.1.2 | 3.2.1 | medium | [CVE-2024-2511](https://osv.dev/vulnerability/ALPINE-CVE-2024-2511) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2026-22795](https://osv.dev/vulnerability/ALPINE-CVE-2026-22795) |
| openssl-fips | 3.1.2 | 3.3.3 | medium | [CVE-2024-12797](https://osv.dev/vulnerability/ALPINE-CVE-2024-12797) |
| openssl-fips | 3.1.2 | 3.5.1 | medium | [CVE-2025-4575](https://osv.dev/vulnerability/ALPINE-CVE-2025-4575) |
| openssl-fips | 3.1.2 | 3.5.5 | medium | [CVE-2025-69418](https://osv.dev/vulnerability/ALPINE-CVE-2025-69418) |
| openssl-fips | 3.1.2 | 3.3.2 | medium | [CVE-2024-13176](https://osv.dev/vulnerability/ALPINE-CVE-2024-13176) |
| busybox | 1.37.0 | 1.37.0 | low | [CVE-2024-58251](https://osv.dev/vulnerability/ALPINE-CVE-2024-58251) |
| busybox | 1.37.0 | 1.37.0 | low | [CVE-2025-46394](https://osv.dev/vulnerability/ALPINE-CVE-2025-46394) |
| openssl-fips | 3.1.2 | 3.5.7 | low | [CVE-2026-42768](https://osv.dev/vulnerability/ALPINE-CVE-2026-42768) |
| openssl-fips | 3.1.2 | 3.5.7 | low | [CVE-2026-42770](https://osv.dev/vulnerability/ALPINE-CVE-2026-42770) |

## 📋 Open PRs

_None._

## 🤖 Bot PR ledger

| Repo | Bump | Kind | Source | State | PR |
|---|---|---|---|---|---|
| [mudler/edgevpn](https://github.com/mudler/edgevpn) | golang.org/x/net@0.33.0 | direct | ksec | error | — |

## 🔎 Bot-PR reviews

**[kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot)**

- [#409](https://github.com/kairos-io/AuroraBoot/pull/409) — ⚠️ **needs_human_verification** — Although this PR primarily updates a dependency digest, the diff shows extensive new code implementing a complex key hierarchy and backend support for PKI, TPM, and YubiKey. This level of architectural change requires a thorough manual security and functional review before merging.
  ↳ This PR updates the digest for github.com/foxboron/sbctl and introduces significant new functionality for cryptographic key management, including support for PKI, TPM, and YubiKey backends. It also includes updates to CI/CD workflows and various internal code files.
    - github.com/foxboron/sbctl 0.0.0-20240526163235-64e649b31c8e→0.0.0-20260316200809-1b913e78d38c: compare 64e649b31c8e...1b913e78d38c ✓ 40000 bytes
    - github.com/fatih/color 1.15.0→1.17.0: compare v1.15.0...v1.17.0 ✓ 9976 bytes
    - context: 58356 bytes
- [#566](https://github.com/kairos-io/AuroraBoot/pull/566) — ✅ **good** — This is a standard dependency update for a widely used icon library, and the version bump appears to follow a normal release progression. The changes introduced are related to configuration cleanup and adding a new icon, which do not present any immediate security risks. Therefore, it is safe to auto-approve.
  ↳ This PR updates the `lucide-react` dependency from version ^0.468.0 to ^0.577.0. The changes include a configuration update in `.github/workflows/ci.yml`, a `.gitignore` file cleanup, and the addition of a new `ellipse` icon definition.
    - lucide-icons/lucide 0.576.0..0.577.0 (PR body): compare 0.576.0...0.577.0 ✓ 40000 bytes
    - context: 99754 bytes
- [#588](https://github.com/kairos-io/AuroraBoot/pull/588) — ✅ **good** — This is a routine dependency update managed by an automated tool (Mend Renovate) to update a base OS version. As long as the new version does not introduce critical regressions or known security issues, this change is safe to merge automatically.
  ↳ This PR updates the version tag for the 'debian' dependency from v12 to v13 within the codebase. This is a routine dependency maintenance update.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1596 bytes
- [#590](https://github.com/kairos-io/AuroraBoot/pull/590) — ✅ **good** — This is a routine dependency update for a well-known package. The change is a version bump from v15 to v17, which is a standard maintenance task. Since the PR is generated by a bot and the change is a dependency update, it is considered safe to auto-approve.
  ↳ This PR updates the `globals` dependency from version 15.x.x to 17.7.0. This is a routine dependency update to pull in newer features and bug fixes from the library.
    - sindresorhus/globals v17.6.0..a19670cc86c1218e915657c55ea02ba3e7623834 (PR body): compare v17.6.0...a19670cc86c1218e915657c55ea02ba3e7623834 ✓ 11637 bytes
    - sindresorhus/globals v17.6.0..v17.7.0 (PR body): compare v17.6.0...v17.7.0 ✓ 11637 bytes
    - sindresorhus/globals v17.5.0..v17.6.0 (PR body): compare v17.5.0...v17.6.0 ✓ 3099 bytes
    - sindresorhus/globals v17.4.0..v17.5.0 (PR body): compare v17.4.0...v17.5.0 ✓ 5103 bytes
    - sindresorhus/globals v17.3.0..v17.4.0 (PR body): compare v17.3.0...v17.4.0 ✓ 4284 bytes
    - context: 45798 bytes
- [#591](https://github.com/kairos-io/AuroraBoot/pull/591) — ✅ **good** — This is a standard dependency update to a minor patch release of a major library (TypeScript). The changes are purely mechanical version bumps and necessary internal code adjustments to accommodate the new version. There are no apparent security risks introduced by this update.
  ↳ This PR updates the TypeScript dependency from version 6.0.2 to 6.0.3. The changes involve updating version strings in `package.json` and `package-lock.json`, along with minor internal code adjustments in TypeScript source files to align with the new version's API.
    - microsoft/TypeScript v6.0.2..v6.0.3 (PR body): compare v6.0.2...v6.0.3 ✓ 40000 bytes
    - microsoft/TypeScript v5.9.3..v6.0.2 (PR body): compare v5.9.3...v6.0.2 ✓ 40000 bytes
    - context: 85089 bytes
- [#594](https://github.com/kairos-io/AuroraBoot/pull/594) — ✅ **good** — This is a necessary major version upgrade for a core dependency. The changes appear to be forward-compatible, addressing deprecations by adding backward compatibility rules and implementing necessary internal refactors in React's reconciliation and testing logic. The comprehensive set of updates suggests a healthy and necessary evolution of the project's tooling.
  ↳ This PR updates `eslint-plugin-react-hooks` from v5 to v7.1.1, which introduces support for ESLint v10 and includes internal refactors to handle deprecations and new React features. The changes also include updates to build configurations and internal testing logic across React packages.
    - facebook/react eslint-plugin-react-hooks@7.1.0..eslint-plugin-react-hooks@7.1.1 (PR body): compare eslint-plugin-react-hooks@7.1.0...eslint-plugin-react-hooks@7.1.1 ✓ 24066 bytes
    - facebook/react 408b38ef7304faf022d2a37110c57efce12c6bad..eslint-plugin-react-hooks@7.1.0 (PR body): compare 408b38ef7304faf022d2a37110c57efce12c6bad...eslint-plugin-react-hooks@7.1.0 ✓ 40000 bytes
    - context: 100048 bytes
- [#599](https://github.com/kairos-io/AuroraBoot/pull/599) — ✅ **good** — This is a standard dependency update for a core tool (ESLint) driven by an automated bot. The changelog indicates that the upgrade includes bug fixes and security vulnerability patches, making this change safe to merge. The changes are well-tracked by Renovate and appear to be a routine maintenance task.
  ↳ This PR updates the core ESLint packages, `@eslint/js` and `eslint`, to version 10.0.1. This upgrade includes several bug fixes and dependency updates, such as updating `minimatch` to address security vulnerabilities. It also includes corresponding updates to configuration files, documentation, and internal code to align with the new version.
    - eslint/eslint v10.0.0..v10.0.1 (PR body): compare v10.0.0...v10.0.1 ✓ 40000 bytes
    - context: 77824 bytes
- [#612](https://github.com/kairos-io/AuroraBoot/pull/612) — ✅ **good** — This is a routine dependency update to a newer minor version of Vite, driven by the upstream project. The changes primarily consist of bug fixes and dependency bumps, which are standard maintenance activities and do not introduce any apparent security risks.
  ↳ This PR updates the `vite` dependency from v8.1.2 to v8.1.3, which includes several bug fixes and updates the `es-module-lexer` dependency. This is a routine maintenance update for the build tool.
    - vitejs/vite v8.1.2..v8.1.3 (PR body): compare v8.1.2...v8.1.3 ✓ 21170 bytes
    - context: 24551 bytes
- [#613](https://github.com/kairos-io/AuroraBoot/pull/613) — ✅ **good** — This is a standard dependency bump to a newer minor version of a well-known library. The upstream release notes indicate feature additions and internal improvements rather than critical security fixes, making this update safe to merge automatically.
  ↳ Update github.com/klauspost/compress from v1.18.6 to v1.19.0, incorporating new features like improved zstd encoding, added inflate checkpoints, and internal logic adjustments to support these changes.
    - github.com/klauspost/compress 1.18.6→1.19.0: compare v1.18.6...v1.19.0 ✓ 40000 bytes
    - klauspost/compress v1.18.7..v1.19.0 (PR body): compare v1.18.7...v1.19.0 ✓ 40000 bytes
    - context: 86385 bytes
- [#614](https://github.com/kairos-io/AuroraBoot/pull/614) — ✅ **good** — This pull request is a routine dependency update, bumping the `fedora` package version from 44 to 45. Since this is an automated dependency update, it is generally safe to approve, assuming the upstream version is stable and the change is intended for maintenance.
  ↳ This PR updates the version of the `fedora` dependency from v44 to v45 in both the Dockerfile and a TypeScript configuration file.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1894 bytes
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
- [#19](https://github.com/kairos-io/entangle-proxy/pull/19) — ✅ **good** — This pull request is a standard dependency update to a major version of a widely used action. The changelog indicates routine maintenance and dependency bumps, which typically include security patches. Since this is an automated PR, and the changes appear to be standard version upgrades, it is safe to auto-approve.
  ↳ This PR updates the `docker/login-action` dependency from v1 to v4, which includes several underlying dependency bumps like AWS SDKs and sigstore. Additionally, the PR updates various GitHub Actions and related workflow files to newer versions, ensuring the CI/CD setup remains current.
    - docker/login-action v4.2.0..v4.3.0 (PR body): compare v4.2.0...v4.3.0 ✓ 40000 bytes
    - context: 64282 bytes
- [#20](https://github.com/kairos-io/entangle-proxy/pull/20) — ⚠️ **needs_human_verification** — Upgrading core Kubernetes dependencies is generally beneficial for security and feature updates. However, this is a major version jump (v0.24.0 to v0.36.2) across multiple critical libraries. A human review is required to ensure that no breaking changes have been introduced that could impact the application's functionality or security posture.
  ↳ This pull request updates three core Kubernetes dependencies (k8s.io/api, k8s.io/apimachinery, and k8s.io/client-go) from version v0.24.0 to v0.36.2. This involves significant version bumps across the board.
    - k8s.io/api 0.24.0→0.36.2: compare v0.24.0...v0.36.2 ✓ 40000 bytes
    - context: 126004 bytes
- [#23](https://github.com/kairos-io/entangle-proxy/pull/23) — ✅ **good** — This pull request updates the dependency for 'actions/checkout' from version v2 to v7 in two workflow files. This is a standard dependency update to a newer version, which is generally safe and beneficial for security and maintenance.
- [#24](https://github.com/kairos-io/entangle-proxy/pull/24) — ✅ **good** — This is a standard dependency update to a newer major version of a library. The release notes for v2.32.0 indicate new features and maintenance updates, and there are no immediate security concerns suggested by this change. Therefore, it is safe to auto-approve.
  ↳ This PR updates the dependency `github.com/onsi/ginkgo` from version v1.16.5 to v2.32.0. This upgrade brings new features such as `-fd` for RSpec-style documentation output and a `--sleep-on-failure` flag to pause a failed spec for debugging purposes.
    - onsi/ginkgo v2.31.0..v2.32.0 (PR body): compare v2.31.0...v2.32.0 ✓ 35617 bytes
    - context: 94069 bytes
**[kairos-io/go-nodepair](https://github.com/kairos-io/go-nodepair)**

- [#27](https://github.com/kairos-io/go-nodepair/pull/27) — ✅ **good** — The changes involve updating several core dependencies across the project. The changelogs indicate that these updates include important security patches, such as restricting RSA key sizes in go-libp2p and fixing memory exhaustion attacks in quic-go. This is standard maintenance and security hygiene.
  ↳ This pull request updates several core dependencies, including go-libp2p, quic-go, golang.org/x/crypto, golang.org/x/image, golang.org/x/net, and google.golang.org/protobuf. The updates include critical security fixes, such as mitigating a DoS attack in go-libp2p and addressing memory exhaustion issues in quic-go.
- [#37](https://github.com/kairos-io/go-nodepair/pull/37) — ✅ **good** — This appears to be a standard dependency update initiated by Renovate. The changes involve updating the dependency digest and corresponding internal code refactoring, which is typical when upgrading a library. Since this is an automated PR, we trust the upstream maintainer, and the changes do not immediately suggest a critical security risk.
  ↳ This PR updates the dependency github.com/kbinani/screenshot to a new digest, which includes significant internal refactoring within the package to support new platform abstractions like Wayland and D-Bus. This involves renaming files, removing internal utility functions, and introducing new platform-specific capture logic.
    - github.com/kbinani/screenshot 0.0.0-20230812210009-b87d31814237→0.0.0-20250624051815-089614a94018: compare b87d31814237...089614a94018 ✓ 23617 bytes
    - github.com/gen2brain/shm 0.0.0-20230802011745-f2460f5984f7→0.1.0: compare f2460f5984f7...v0.1.0 ✓ 1878 bytes
    - github.com/jezek/xgb 1.1.0→1.1.1: compare v1.1.0...v1.1.1 ✓ 431 bytes
    - context: 32368 bytes
- [#46](https://github.com/kairos-io/go-nodepair/pull/46) — ⚠️ **needs_human_verification** — This is a major version bump for a core testing framework, which involves significant internal refactoring and new features. While the changes appear to be feature additions and internal improvements, a human review is necessary to ensure compatibility with the Kairos project's existing code, verify that the new features do not introduce regressions, and confirm that no security vulnerabilities were inadvertently introduced in the updated library.
  ↳ This PR updates the dependency github.com/onsi/ginkgo/v2 from v2.29.0 to v2.32.0. The update introduces new features like the `--sleep-on-failure` debugging flag, enhancements to suite management via `globals.Reset()`, and integration for Claude Code skills. It also includes significant internal refactoring across the library's core logic and documentation.
    - github.com/onsi/ginkgo/v2 2.29.0→2.32.0: compare v2.29.0...v2.32.0 ✓ 40000 bytes
    - onsi/ginkgo v2.31.0..v2.32.0 (PR body): compare v2.31.0...v2.32.0 ✓ 35617 bytes
    - context: 80051 bytes
- [#57](https://github.com/kairos-io/go-nodepair/pull/57) — ✅ **good** — This pull request only updates the version of the 'actions/setup-go' action from v5 to v6. This is a routine maintenance update for a standard dependency and does not introduce any new security risks or significant functional changes that require manual review.
- [#59](https://github.com/kairos-io/go-nodepair/pull/59) — ✅ **good** — This pull request only modifies the configuration file for the Renovate bot. The changes involve migrating to a recommended configuration and adjusting package matching rules. This is a standard configuration update and poses no security risk.
- [#63](https://github.com/kairos-io/go-nodepair/pull/63) — ✅ **good** — The changes introduce significant new features, specifically robust integration with the standard Go `log/slog` package and fine-grained, per-subsystem log level control. The extensive test coverage provided in `setup_test.go` suggests the new logic is sound and safe. This is a necessary and beneficial update for the project's logging infrastructure.
  ↳ This PR updates `github.com/ipfs/go-log` to v2.9.2, introducing comprehensive support for Go's `log/slog` package. Key changes include implementing a `slog` bridge to route logs through Zap, adding subsystem-aware level control via atomic levels, and improving handling for `slog.Group` attributes.
    - ipfs/go-log v2.9.1..v2.9.2 (PR body): compare v2.9.1...v2.9.2 ✓ 12110 bytes
    - ipfs/go-log v2.9.0..v2.9.1 (PR body): compare v2.9.0...v2.9.1 ✓ 2149 bytes
    - ipfs/go-log v2.8.2..v2.9.0 (PR body): compare v2.8.2...v2.9.0 ✓ 40000 bytes
    - context: 70615 bytes
**[kairos-io/go-ukify](https://github.com/kairos-io/go-ukify)**

- [#58](https://github.com/kairos-io/go-ukify/pull/58) — ✅ **good** — This is a standard dependency update to a minor version, and the upstream release notes document the changes, including adjustments to the Go version policy. Since this is an automated bot PR, and the changes appear to be documented, it is safe to auto-approve.
  ↳ This PR updates the dependency `github.com/ThalesGroup/crypto11` from version v1.6.1 to v1.6.2. This update includes changes to the Go version directives in `go.mod` and updates a sub-dependency, `github.com/miekg/pkcs11`, to a newer version.
    - github.com/ThalesGroup/crypto11 1.6.1→1.6.2: compare v1.6.1...v1.6.2 ✓ 1201 bytes
    - thales-transfer/crypto11 v1.6.0..v1.6.2 (PR body): compare v1.6.0...v1.6.2 ✓ 2393 bytes
    - context: 6940 bytes
**[kairos-io/kairos](https://github.com/kairos-io/kairos)**

- [#4104](https://github.com/kairos-io/kairos/pull/4104) — ⚠️ **needs_human_verification** — The PR title suggests an automation task related to dependency upgrades. While automation can be beneficial, security review requires inspecting the actual code changes to ensure no unintended side effects or vulnerabilities were introduced during the pipeline wiring.
  ↳ The PR aims to automate the process of fetching the latest release and running validation tests within the upgrade pipeline.
- [#4209](https://github.com/kairos-io/kairos/pull/4209) — ✅ **good** — This is a routine dependency update to a minor version of a well-maintained action. There are no apparent security risks introduced by this update, and it aligns with standard maintenance practices.
  ↳ This PR updates the `docker/build-push-action` dependency across multiple CI/CD workflows to version v7.3.0. This is a routine minor version bump for a widely used action.
    - docker/build-push-action v7.2.0..v7.3.0 (PR body): compare v7.2.0...v7.3.0 ✓ 40000 bytes
    - context: 42761 bytes
- [#4210](https://github.com/kairos-io/kairos/pull/4210) — ✅ **good** — This PR is a standard dependency update for a specific action, which is a routine maintenance task. Since it is a dependency bump and the changes involve updating to newer versions of related tools, it is considered safe to auto-approve.
  ↳ This pull request updates the digest for the `docker/login-action` dependency from `650006c` to `c99871d`. Additionally, it updates several other related GitHub Actions and Docker tooling dependencies to their latest versions.
    - docker/login-action 650006c6eb7dba73a995cc03b0b2d7f5ca915bee..c99871dec2022cc055c062a10cc1a1310835ceb4 (PR body): compare 650006c6eb7dba73a995cc03b0b2d7f5ca915bee...c99871dec2022cc055c062a10cc1a1310835ceb4 ✓ 40000 bytes
    - context: 43062 bytes
- [#4211](https://github.com/kairos-io/kairos/pull/4211) — ✅ **good** — This is a routine dependency update to a newer digest, which typically includes security patches and bug fixes. The change is localized to updating the action reference in CI/CD workflows and poses no immediate security risk.
  ↳ This PR updates the digest for the `docker/setup-buildx-action` dependency across multiple GitHub Actions workflows from an older version to a newer one.
    - docker/setup-buildx-action d7f5e7f509e45cec5c76c4d5afdd7de93d0b3df5..bb05f3f5519dd87d3ba754cc423b652a5edd6d2c (PR body): compare d7f5e7f509e45cec5c76c4d5afdd7de93d0b3df5...bb05f3f5519dd87d3ba754cc423b652a5edd6d2c ✓ 40000 bytes
    - context: 43864 bytes
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

- [#41](https://github.com/kairos-io/kcrypt-discovery-challenger/pull/41) — ✅ **good** — This is a dependency update for a core Kubernetes controller library. The changelog indicates bug fixes and necessary dependency alignment across the Kubernetes ecosystem, which is standard maintenance practice. The changes appear safe and necessary for project health.
  ↳ This PR updates the `sigs.k8s.io/controller-runtime` dependency from v0.15.0 to v0.24.1. This update includes several bug fixes and dependency upgrades across the Kubernetes ecosystem, including major bumps in core `k8s.io` packages.
    - k8s.io/api 0.27.2→0.36.0: compare v0.27.2...v0.36.0 ✓ 40000 bytes
    - context: 123948 bytes
- [#190](https://github.com/kairos-io/kcrypt-discovery-challenger/pull/190) — ✅ **good** — Updating core infrastructure dependencies like Kubernetes components to the latest stable version is a crucial security and stability practice. This change incorporates bug fixes and security patches from the upstream, making the project more resilient. Therefore, it is safe to auto-approve.
  ↳ This PR updates the core Kubernetes dependencies, k8s.io/api, k8s.io/apimachinery, and k8s.io/client-go, to version v0.36.2. This brings the project up to a recent, patched version of the Kubernetes ecosystem components.
    - k8s.io/apimachinery 0.27.4→0.27.2: compare v0.27.4...v0.27.2 failed: <nil> (no upstream diff)
    - github.com/emicklei/go-restful/v3 3.10.1→3.13.0: compare v3.10.1...v3.13.0 ✓ 40000 bytes
    - context: 131955 bytes
- [#240](https://github.com/kairos-io/kcrypt-discovery-challenger/pull/240) — ⚠️ **needs_human_verification** — review endpoint unreachable: Post "http://localhost:8080/v1/chat/completions": context deadline exceeded
    - github.com/google/go-attestation 0.5.1→0.6.1: compare v0.5.1...v0.6.1 ✓ 40000 bytes
    - github.com/kairos-io/tpm-helpers 0.0.0-20260608091616-8a4ccb53d8f7→0.0.0-20260702080541-9b3e057e2f32: compare 8a4ccb53d8f7...9b3e057e2f32 ✓ 11771 bytes
    - github.com/google/go-tpm-tools 0.4.4→0.4.7: compare v0.4.4...v0.4.7 ✓ 40000 bytes
    - context: 97184 bytes
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
**[mauromorales/xpasswd](https://github.com/mauromorales/xpasswd)**

- [#47](https://github.com/mauromorales/xpasswd/pull/47) — ✅ **good** — The pull request is a routine dependency update to a newer version of Ginkgo. The changes documented in the changelog point to new features and documentation improvements, which are beneficial for the project. There are no apparent security risks introduced by this version bump. Therefore, it is safe to auto-approve.
  ↳ This PR updates the dependency `github.com/onsi/ginkgo/v2` from v2.28.3 to v2.32.0. The changes include new features like RSpec-style documentation output, updates to the plugin marketplace, and significant refactoring and documentation improvements within the Ginkgo project itself. This is a standard dependency upgrade.
    - github.com/onsi/ginkgo/v2 2.28.3→2.32.0: compare v2.28.3...v2.32.0 ✓ 40000 bytes
    - context: 44543 bytes
- [#48](https://github.com/mauromorales/xpasswd/pull/48) — ✅ **good** — This is a dependency update to a well-known testing library, Gomega. The changes involve adding new matchers and internal code refactoring for better debugging, which are generally positive for maintainability. There are no obvious security vulnerabilities introduced by this update.
  ↳ This PR updates the dependency `github.com/onsi/gomega` from v1.40.0 to v1.42.1. It introduces new matchers (`BeASlice`, `BeAnArray`) and significantly refactors the internal formatting functions to handle cyclic references for improved debugging output. The plugin version is also updated to match the new library version.
    - github.com/onsi/gomega 1.40.0→1.42.1: compare v1.40.0...v1.42.1 ✓ 40000 bytes
    - golang.org/x/mod 0.35.0→0.36.0: compare v0.35.0...v0.36.0 ✓ 16991 bytes
    - golang.org/x/net 0.53.0→0.56.0: compare v0.53.0...v0.56.0 ✓ 40000 bytes
    - context: 105286 bytes
- [#49](https://github.com/mauromorales/xpasswd/pull/49) — ✅ **good** — This change is a minor patch update to the Go toolchain, moving from 1.26.3 to 1.26.4. Toolchain updates are generally low-risk maintenance tasks. Therefore, it is safe to auto-approve.
  ↳ This PR updates the Go toolchain directive in go.mod from version 1.26.3 to 1.26.4. This is a routine patch update to the Go toolchain.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 1472 bytes
- [#50](https://github.com/mauromorales/xpasswd/pull/50) — ✅ **good** — This is a standard dependency update for a widely used GitHub Action. The upgrade is to a newer major version (v7), which typically includes security patches and improvements. Since this is an automated dependency update, it is safe to auto-approve.
  ↳ This PR updates the version of the `actions/checkout` dependency from v6 to v7 across the workflow files. This upgrade incorporates updates from the official release notes, which include dependency bumps and updates to related actions like `actions/core`.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 5388 bytes
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
- [#1009](https://github.com/mudler/edgevpn/pull/1009) — ✅ **good** — This is a routine dependency update for a documentation library. The change is a simple commit hash update, which is typical for dependency management. Without further context on the specific changes in the new version, this appears safe to auto-approve.
  ↳ This pull request updates the dependency `docs/themes/docsy` from commit `bbf68d4` to `01c827e`. This change reflects a standard version bump for the documentation library, likely preparing for a release.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 4207 bytes
- [#1041](https://github.com/mudler/edgevpn/pull/1041) — ✅ **good** — This is a dependency upgrade to a newer version of a library, which is generally beneficial. Crucially, this update incorporates a security fix that mitigates a path traversal/ACL bypass vulnerability related to URL-encoded separators in static file serving. Therefore, this change is safe to auto-approve.
  ↳ The PR updates the `github.com/labstack/echo/v4` dependency to `v5.2.1`. This upgrade includes a security fix that changes the default behavior of static file serving to prevent encoded path separators from bypassing route-level access controls.
    - github.com/labstack/echo/v4 4.15.2→4.15.4: compare v4.15.2...v4.15.4 ✓ 30288 bytes
    - github.com/mattn/go-colorable 0.1.14→0.1.15: compare v0.1.14...v0.1.15 ✓ 5234 bytes
    - golang.org/x/crypto 0.52.0→0.53.0: compare v0.52.0...v0.53.0 ✓ 40000 bytes
    - context: 108729 bytes
- [#1045](https://github.com/mudler/edgevpn/pull/1045) — ⚠️ **needs_human_verification** — Upgrading to a major version introduces potential breaking changes that require thorough testing to ensure the project remains functional and secure. A human review is necessary to validate that all functionality is preserved after this substantial refactoring.
  ↳ This PR updates the dependency `github.com/urfave/cli/v2` from version v2.27.7 to v3.10.1. This major version upgrade includes significant refactoring across command parsing, flag handling, completion logic, and documentation.
    - github.com/urfave/cli/v3 3.10.0→3.10.1: compare v3.10.0...v3.10.1 ✓ 17319 bytes
    - urfave/cli v3.9.1..v3.10.0 (PR body): compare v3.9.1...v3.10.0 ✓ 40000 bytes
    - context: 100044 bytes
- [#1046](https://github.com/mudler/edgevpn/pull/1046) — ✅ **good** — This is a routine dependency update to a newer version of Autoprefixer. The changes primarily address compatibility fixes and minor logic adjustments, which are generally safe and beneficial for the project. There are no apparent security risks introduced by this version bump.
  ↳ The PR updates the `autoprefixer` dependency from v10.5.0 to v10.5.2. This update includes fixes for Firefox compatibility related to `-webkit-fill-available` and related adjustments in source code and test files.
    - postcss/autoprefixer 10.5.1..10.5.2 (PR body): compare 10.5.1...10.5.2 ✓ 2688 bytes
    - postcss/autoprefixer 10.5.0..10.5.1 (PR body): compare 10.5.0...10.5.1 ✓ 40000 bytes
    - context: 54014 bytes
- [#1051](https://github.com/mudler/edgevpn/pull/1051) — ✅ **good** — The changes are internal fixes and structural improvements within the dependency itself, specifically addressing keyspace region logic and adding performance-critical features like key counting. Since this is an update to a newer version and the changes appear to be bug fixes and optimizations from the upstream project, it is safe to auto-approve.
  ↳ This PR updates the `go-libp2p-kad-dht` dependency to v0.41.0. The changes include internal refactoring of the keyspace region calculation logic, performance enhancements via the introduction of `CountKeysUpTo` for efficient key counting, and various bug fixes within the provider and keystore implementations.
    - github.com/libp2p/go-libp2p-kad-dht 0.40.0→0.41.0: compare v0.40.0...v0.41.0 ✓ 40000 bytes
    - github.com/quic-go/quic-go 0.59.0→0.59.1: compare v0.59.0...v0.59.1 ✓ 7704 bytes
    - context: 53854 bytes
- [#1052](https://github.com/mudler/edgevpn/pull/1052) — ✅ **good** — This is a dependency update to a major version of a library. The context shows that the upstream changes are comprehensive, including updated documentation, API changes, and corresponding test updates, suggesting a clean migration. Since this is a dependency bump and no application code changes are introduced, it is safe to auto-approve.
  ↳ This PR updates the dependency `github.com/cenkalti/backoff` from version v4 to v7. This migration involves updating import paths, modifying the `RetryAfter` function signature, and updating error handling logic within the package. The upstream changes appear to be a standard, documented migration path.
    - cenkalti/backoff v6.0.1..v7.0.0 (PR body): compare v6.0.1...v7.0.0 ✓ 21611 bytes
    - cenkalti/backoff v6.0.0..v6.0.1 (PR body): compare v6.0.0...v6.0.1 ✓ 2512 bytes
    - cenkalti/backoff v5.0.3..v6.0.0 (PR body): compare v5.0.3...v6.0.0 ✓ 25480 bytes
    - cenkalti/backoff v5.0.2..v5.0.3 (PR body): compare v5.0.2...v5.0.3 ✓ 3314 bytes
    - cenkalti/backoff v5.0.1..v5.0.2 (PR body): compare v5.0.1...v5.0.2 ✓ 2019 bytes
    - context: 59288 bytes

