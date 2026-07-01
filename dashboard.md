# Kairos Security Dashboard

_Updated 2026-07-01._

🌐 **[Live dashboard](https://kairos-io.github.io/security/)** — the published board with clickable links.

## 📋 This run

- **Scanned:** 28 repos (1 skipped)
- **Findings:** 0 (0 critical / 0 high / 0 medium / 0 low / 0 unknown)
- **CVE-related PRs:** 0
- **Remediation:** 0 open · 0 superseded · 0 merged · 0 need-human
- **Why:** No CVEs found across 28 repos — nothing to remediate.

> No security findings to triage this run.

## 🔥 Focus now

_Nothing flagged._

## 🌊 Waterfall fronts

_None._

## 📦 Per-repo findings

| Repo | Critical | High | Medium | Low | Total | Status |
|---|---|---|---|---|---|---|
| [kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/cluster-api-provider-kairos](https://github.com/kairos-io/cluster-api-provider-kairos) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/entangle](https://github.com/kairos-io/entangle) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/entangle-proxy](https://github.com/kairos-io/entangle-proxy) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/go-nodepair](https://github.com/kairos-io/go-nodepair) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/go-ukify](https://github.com/kairos-io/go-ukify) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/hadron](https://github.com/kairos-io/hadron) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/immucore](https://github.com/kairos-io/immucore) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/kairos](https://github.com/kairos-io/kairos) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/kairos-agent](https://github.com/kairos-io/kairos-agent) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/kairos-init](https://github.com/kairos-io/kairos-init) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/kairos-installer](https://github.com/kairos-io/kairos-installer) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/kairos-lab](https://github.com/kairos-io/kairos-lab) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/kairos-must-burn](https://github.com/kairos-io/kairos-must-burn) | 0 | 0 | 0 | 0 | 0 | skipped: not source-scannable |
| [kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/kairos-sdk](https://github.com/kairos-io/kairos-sdk) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/kcrypt](https://github.com/kairos-io/kcrypt) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/kcrypt-discovery-challenger](https://github.com/kairos-io/kcrypt-discovery-challenger) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/netboot](https://github.com/kairos-io/netboot) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/provider-kairos](https://github.com/kairos-io/provider-kairos) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/provider-kubernetes](https://github.com/kairos-io/provider-kubernetes) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/simple-mdns-server](https://github.com/kairos-io/simple-mdns-server) | 0 | 0 | 0 | 0 | 0 | clean |
| [kairos-io/tpm-helpers](https://github.com/kairos-io/tpm-helpers) | 0 | 0 | 0 | 0 | 0 | clean |
| [mauromorales/xpasswd](https://github.com/mauromorales/xpasswd) | 0 | 0 | 0 | 0 | 0 | clean |
| [mudler/edgevpn](https://github.com/mudler/edgevpn) | 0 | 0 | 0 | 0 | 0 | clean |
| [mudler/entities](https://github.com/mudler/entities) | 0 | 0 | 0 | 0 | 0 | clean |
| [mudler/go-pluggable](https://github.com/mudler/go-pluggable) | 0 | 0 | 0 | 0 | 0 | clean |
| [mudler/yip](https://github.com/mudler/yip) | 0 | 0 | 0 | 0 | 0 | clean |

## 📋 Open PRs

_None._

## 🤖 Bot PR ledger

| Repo | Bump | Kind | Source | State | PR |
|---|---|---|---|---|---|
| [mudler/edgevpn](https://github.com/mudler/edgevpn) | golang.org/x/net@0.33.0 | direct | ksec | error | — |

## 🔎 Bot-PR reviews

**[kairos-io/AuroraBoot](https://github.com/kairos-io/AuroraBoot)**

- [#409](https://github.com/kairos-io/AuroraBoot/pull/409) — ⚠️ **needs_human_verification** — While updating a dependency digest is a standard maintenance task, this PR involves significant internal code changes, including the introduction of new backend implementations for key management and numerous dependency upgrades. A human review is necessary to ensure the refactoring is correct, secure, and aligns with project standards before merging.
  ↳ This PR updates the digest for the dependency `github.com/foxboron/sbctl` to a new version. Additionally, it includes extensive internal refactoring, such as adding new key backends (File, TPM, Yubikey), updating build/CI workflows, and upgrading several other dependencies.
    - github.com/foxboron/sbctl 0.0.0-20240526163235-64e649b31c8e→0.0.0-20260316200809-1b913e78d38c: compare 64e649b31c8e...1b913e78d38c ✓ 40000 bytes
    - github.com/fatih/color 1.15.0→1.17.0: compare v1.15.0...v1.17.0 ✓ 9976 bytes
    - context: 58356 bytes
- [#566](https://github.com/kairos-io/AuroraBoot/pull/566) — ✅ **good** — This is a standard dependency update managed by Renovate, moving to a newer version of a trusted library. The changes align with the upstream release notes, specifically adding a new icon. There are no obvious security vulnerabilities introduced by this version bump.
  ↳ This PR updates the `lucide-react` dependency from version ^0.468.0 to ^0.577.0. This update includes a new icon, `ellipse`, and minor configuration adjustments in the CI workflow and `.gitignore` file.
    - lucide-icons/lucide 0.576.0..0.577.0 (PR body): compare 0.576.0...0.577.0 ✓ 40000 bytes
    - context: 99754 bytes
- [#588](https://github.com/kairos-io/AuroraBoot/pull/588) — ✅ **good** — This pull request is a routine dependency update generated by an automated tool (Mend Renovate) to bump the debian version from 12 to 13. This type of maintenance update is generally safe and necessary for keeping dependencies up-to-date.
  ↳ This PR updates the version of the debian dependency tag from v12 to v13 in the configuration file. This is a routine maintenance update to keep the project using a newer base image version.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1596 bytes
- [#590](https://github.com/kairos-io/AuroraBoot/pull/590) — ✅ **good** — This is a routine dependency update for a utility package. The change is a version bump for `globals`, which is a common maintenance task. While the `globals` package has had breaking changes in version 17.0.0, this PR updates to a version within that series, which is generally safe for automated approval unless specific breaking changes are known to impact the project's functionality.
  ↳ This PR updates the dependency `globals` from version 15.x.x to version 17.7.0. This change updates the package version in `package.json` and `package-lock.json` to align with the new major version.
    - sindresorhus/globals v17.6.0..a19670cc86c1218e915657c55ea02ba3e7623834 (PR body): compare v17.6.0...a19670cc86c1218e915657c55ea02ba3e7623834 ✓ 11637 bytes
    - sindresorhus/globals v17.6.0..v17.7.0 (PR body): compare v17.6.0...v17.7.0 ✓ 11637 bytes
    - sindresorhus/globals v17.5.0..v17.6.0 (PR body): compare v17.5.0...v17.6.0 ✓ 3099 bytes
    - sindresorhus/globals v17.4.0..v17.5.0 (PR body): compare v17.4.0...v17.5.0 ✓ 5103 bytes
    - sindresorhus/globals v17.3.0..v17.4.0 (PR body): compare v17.3.0...v17.4.0 ✓ 4284 bytes
    - context: 45798 bytes
- [#591](https://github.com/kairos-io/AuroraBoot/pull/591) — ✅ **good** — This is a standard dependency update to a newer patch version of a widely used library (TypeScript). The changes are confined to updating version numbers in lock files and internal source files. This is a routine maintenance task and poses no immediate security risk.
  ↳ This PR updates the `typescript` dependency from version 6.0.2 to 6.0.3 in `package.json` and `package-lock.json`. It also includes corresponding internal updates within the TypeScript source files to align with the new version. This is a routine dependency maintenance update.
    - microsoft/TypeScript v6.0.2..v6.0.3 (PR body): compare v6.0.2...v6.0.3 ✓ 40000 bytes
    - microsoft/TypeScript v5.9.3..v6.0.2 (PR body): compare v5.9.3...v6.0.2 ✓ 40000 bytes
    - context: 85089 bytes
- [#594](https://github.com/kairos-io/AuroraBoot/pull/594) — ✅ **good** — This is a routine dependency update to a major version of a well-known package. The changes primarily involve version bumping, updating documentation, and implementing deprecation handling for a removed rule, which is standard maintenance. The changes appear safe and necessary for project health.
  ↳ This PR updates the `eslint-plugin-react-hooks` dependency from v5 to v7.1.1, incorporating fixes and handling the deprecation of the `component-hook-factories` rule. It also includes related updates to test files and build cache keys to align with the new major version.
    - facebook/react eslint-plugin-react-hooks@7.1.0..eslint-plugin-react-hooks@7.1.1 (PR body): compare eslint-plugin-react-hooks@7.1.0...eslint-plugin-react-hooks@7.1.1 ✓ 24066 bytes
    - facebook/react 408b38ef7304faf022d2a37110c57efce12c6bad..eslint-plugin-react-hooks@7.1.0 (PR body): compare 408b38ef7304faf022d2a37110c57efce12c6bad...eslint-plugin-react-hooks@7.1.0 ✓ 40000 bytes
    - context: 100048 bytes
- [#599](https://github.com/kairos-io/AuroraBoot/pull/599) — ✅ **good** — This is a dependency update to a newer major version of ESLint. The release notes confirm that this version includes important bug fixes and security vulnerability patches, making this change safe to auto-approve.
  ↳ This PR updates the ESLint and related packages to version 10.0.1. This upgrade includes several bug fixes and security updates, such as updating `minimatch` to address known vulnerabilities. It also includes necessary configuration and documentation updates for the new version.
    - eslint/eslint v10.0.0..v10.0.1 (PR body): compare v10.0.0...v10.0.1 ✓ 40000 bytes
    - context: 77824 bytes
- [#606](https://github.com/kairos-io/AuroraBoot/pull/606) — ✅ **good** — This is a routine dependency update to a patch version of a widely used library. The change is low-risk as it is a minor version bump and addresses a specific bug fix documented in the release notes. It is safe to auto-approve.
  ↳ This PR updates the `react-router-dom` dependency from version 7.18.0 to 7.18.1. This is a patch release that includes a fix for CSRF check logic.
    - remix-run/react-router react-router-dom@7.18.0..react-router-dom@7.18.1 (PR body): compare react-router-dom@7.18.0...react-router-dom@7.18.1 ✓ 40000 bytes
    - context: 43533 bytes
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
- [#18](https://github.com/kairos-io/entangle-proxy/pull/18) — ✅ **good** — This pull request updates the version of the docker/build-push-action from v2 to v7. Updating dependencies to the latest stable version is a standard security and maintenance practice. This change is safe to auto-approve.
- [#19](https://github.com/kairos-io/entangle-proxy/pull/19) — ✅ **good** — This pull request updates the docker/login-action dependency from v1 to v4. Updating dependencies is a standard maintenance practice that generally improves security and stability by incorporating bug fixes and security patches from the maintainers. This change is safe to auto-approve.
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
**[kairos-io/kairos](https://github.com/kairos-io/kairos)**

- [#4104](https://github.com/kairos-io/kairos/pull/4104) — ⚠️ **needs_human_verification** — The PR title suggests an automation task related to dependency upgrades. While automation can be beneficial, security review requires inspecting the actual code changes to ensure no unintended side effects or vulnerabilities were introduced during the pipeline wiring.
  ↳ The PR aims to automate the process of fetching the latest release and running validation tests within the upgrade pipeline.
**[kairos-io/kairos-agent](https://github.com/kairos-io/kairos-agent)**

- [#1282](https://github.com/kairos-io/kairos-agent/pull/1282) — ✅ **good** — The change is a refactoring of the initialization logic to improve configuration flexibility by allowing users to specify provider paths dynamically. The implementation includes error handling to gracefully fall back to default paths if the configuration parsing fails, which enhances robustness. This change does not introduce any obvious security vulnerabilities.
  ↳ This PR refactors the initialization sequence to allow provider directories to be dynamically configured via a `providers.paths` setting in the collector configuration. It introduces YAML parsing to read these paths and uses them to initialize the bus manager, falling back to default paths if parsing fails.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 5365 bytes
**[kairos-io/kairos-operator](https://github.com/kairos-io/kairos-operator)**

- [#88](https://github.com/kairos-io/kairos-operator/pull/88) — ✅ **good** — This change is a routine version bump as described in the PR description, updating images from a previous beta tag to the next one. There are no immediate security concerns apparent from the diff itself, suggesting this is a safe maintenance update.
  ↳ This PR updates the container images for the Kairos operator and node-labeler to version v0.1.0-beta5.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 886 bytes
- [#89](https://github.com/kairos-io/kairos-operator/pull/89) — ✅ **good** — This change is a routine dependency digest update for a widely used Docker action. Updating to a newer digest is standard maintenance and does not introduce any immediate security risks or obvious breaking changes. Therefore, it is safe to auto-approve.
  ↳ This PR updates the digest of the `docker/build-push-action` dependency across multiple workflow files from an older version to a newer one. This is a routine maintenance update for a standard Docker action.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 2152 bytes
- [#91](https://github.com/kairos-io/kairos-operator/pull/91) — ✅ **good** — This is a routine dependency digest update for an external tool. There are no apparent security risks introduced by updating a package digest, and this change is necessary for maintenance.
  ↳ This PR updates the digest for the `azure/setup-kubectl` dependency from an older version to a newer one. This change is applied in the workflow file to ensure the build environment uses the latest specified version of the tool.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 1690 bytes
- [#92](https://github.com/kairos-io/kairos-operator/pull/92) — ✅ **good** — This is a routine dependency digest update for a base image. Updating base image digests is standard practice for security hygiene and ensuring the use of the latest available image, which is safe to auto-approve.
  ↳ This PR updates the base image for the Dockerfile from an older Alpine digest (`2510918`) to a newer one (`fd791d7`). This ensures the build uses the latest, verified image digest for Alpine 3.23.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 1604 bytes
- [#93](https://github.com/kairos-io/kairos-operator/pull/93) — ✅ **good** — The changes are dependency updates for core Kubernetes libraries, which is routine maintenance. The diffs show that the updates primarily involve adding optional fields to Protobuf messages, which is generally a safe, additive change. Since this is a core project dependency, it is likely safe to auto-approve.
  ↳ This PR updates three core Kubernetes dependencies (k8s.io/api, k8s.io/apimachinery, and k8s.io/client-go) to version v0.36.2. The changes primarily involve adding optional fields to Protobuf definitions within admission webhooks, which is an additive change. This update aligns the project with the latest Kubernetes API versions.
    - k8s.io/api 0.35.3→0.36.2: compare v0.35.3...v0.36.2 ✓ 40000 bytes
    - k8s.io/apimachinery 0.35.3→0.36.2: compare v0.35.3...v0.36.2 ✓ 40000 bytes
    - context: 97288 bytes
- [#96](https://github.com/kairos-io/kairos-operator/pull/96) — ✅ **good** — This pull request is a routine version bump for existing images as described in the changelog. Since this is a dependency update and no specific security vulnerabilities are indicated in the context, it is safe to auto-approve.
  ↳ This PR updates the Docker images for the Kairos operator and the operator-node-labeler to version v0.1.0-beta6.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 886 bytes
- [#98](https://github.com/kairos-io/kairos-operator/pull/98) — ✅ **good** — This is a minor version bump for a widely used testing framework. The changes primarily introduce new features, documentation, and minor dependency updates. There are no obvious security vulnerabilities introduced by this update, making it safe to auto-approve.
  ↳ This PR updates the `github.com/onsi/ginkgo/v2` dependency to version v2.32.0. The changes introduce new features such as a Claude plugin, new devcontainer configurations, and updates to documentation and internal code structure. It also includes minor dependency updates in the Gemfile.lock.
    - github.com/onsi/ginkgo/v2 2.28.1→2.32.0: compare v2.28.1...v2.32.0 ✓ 40000 bytes
    - github.com/onsi/gomega 1.39.1→1.40.0: compare v1.39.1...v1.40.0 ✓ 40000 bytes
    - context: 87633 bytes
- [#100](https://github.com/kairos-io/kairos-operator/pull/100) — ✅ **good** — This is a dependency update to a minor version, which is generally safe. The changes appear to be feature additions (Claude plugin, devcontainer setup) and internal code cleanup/refactoring, which are positive for the project. The release notes suggest a strategy to minimize dependency bloat, indicating thoughtful maintenance.
  ↳ This PR updates the dependency `github.com/onsi/gomega` from v1.39.1 to v1.42.1, introducing new features like a Claude plugin and updated development container configurations. It also includes significant internal refactoring in the `format` package and updates to documentation across the repository.
    - github.com/onsi/gomega 1.39.1→1.42.1: compare v1.39.1...v1.42.1 ✓ 40000 bytes
    - golang.org/x/mod 0.35.0→0.36.0: compare v0.35.0...v0.36.0 ✓ 16991 bytes
    - context: 65356 bytes
- [#102](https://github.com/kairos-io/kairos-operator/pull/102) — ✅ **good** — This change is a simple update of the digest for a base image. Updating the digest is a standard practice for ensuring the use of a specific, verified version of an image, which is a routine maintenance task and does not introduce new security risks.
  ↳ This PR updates the digest of the `gcr.io/distroless/static:nonroot` Docker base image to a new, specified digest. This is a routine maintenance update for a base image dependency.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 1727 bytes
- [#112](https://github.com/kairos-io/kairos-operator/pull/112) — ✅ **good** — This is a routine dependency update bumping the image version within the beta track. Since this is a version bump and not a major architectural change, it is safe to auto-approve.
  ↳ This PR updates the image tags for the Kairos operator and node-labeler from v0.1.0-beta4 to v0.1.0-beta7 within the Kustomize configuration.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 886 bytes
- [#117](https://github.com/kairos-io/kairos-operator/pull/117) — ✅ **good** — This is a routine dependency update for a standard GitHub action. Since this change is driven by an automated tool (Mend Renovate) and involves updating a common setup utility, it is generally safe to auto-approve.
  ↳ This PR updates the digest for the `actions/setup-go` action across multiple workflow files from an older version to a newer one.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 3067 bytes
- [#118](https://github.com/kairos-io/kairos-operator/pull/118) — ✅ **good** — This pull request is a routine dependency digest update for the `docker/login-action`. It does not introduce any new code or apparent security vulnerabilities, making it safe to auto-approve.
  ↳ This PR updates the digest of the `docker/login-action` dependency from `4907a6d` to `650006c`. This ensures the project is using the latest version of the action.
    - no go.mod dependency bumps parsed from the PR diff
    - context: 3058 bytes
- [#119](https://github.com/kairos-io/kairos-operator/pull/119) — ✅ **good** — This change is a simple update of a Docker image digest for an existing, specified version of golang. This is a standard maintenance task to ensure the build uses the latest artifact for that tag, which poses a very low security risk.
  ↳ This PR updates the specific Docker digest used for the `docker.io/golang:1.26.4` base image in the `Dockerfile` and `Dockerfile.node-labeler`. This is a routine update to ensure the build uses the latest, verified artifact for that version.
    - no upstream comparisons available (no go.mod bumps or compare links in the PR body)
    - context: 1875 bytes
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
- [#1050](https://github.com/mudler/edgevpn/pull/1050) — ✅ **good** — This is a dependency update to a newer major version of a widely used library. The changes reflect a necessary migration to the v7 API, which appears to be handled correctly by the bot, including updates to documentation and test cases. Since this is a dependency bump, it is safe to auto-approve.
  ↳ This PR updates the dependency `github.com/cenkalti/backoff` from version v4 to v7. This migration involves significant API changes, including updates to error handling, the `Retry` function signature, and documentation across the library.
    - cenkalti/backoff v6.0.1..v7.0.0 (PR body): compare v6.0.1...v7.0.0 ✓ 21611 bytes
    - cenkalti/backoff v6.0.0..v6.0.1 (PR body): compare v6.0.0...v6.0.1 ✓ 2512 bytes
    - cenkalti/backoff v5.0.3..v6.0.0 (PR body): compare v5.0.3...v6.0.0 ✓ 25480 bytes
    - cenkalti/backoff v5.0.2..v5.0.3 (PR body): compare v5.0.2...v5.0.3 ✓ 3314 bytes
    - cenkalti/backoff v5.0.1..v5.0.2 (PR body): compare v5.0.1...v5.0.2 ✓ 2019 bytes
    - context: 59288 bytes
**[mudler/yip](https://github.com/mudler/yip)**

- [#270](https://github.com/mudler/yip/pull/270) — ✅ **good** — This is a routine maintenance update for a core GitHub Action. The changes reflect standard version bumps within the action's dependencies, which is typical for keeping the action secure and compatible. There are no immediate red flags indicating a critical security vulnerability introduced by this version change.
  ↳ This PR updates the `actions/setup-go` action from version v5 to v6. This major version upgrade includes several internal dependency updates across the action's ecosystem, such as upgrading `actions/cache` to v5 and updating various underlying npm packages.
    - actions/setup-go v6.4.0..v6.5.0 (PR body): compare v6.4.0...v6.5.0 ✓ 40000 bytes
    - actions/setup-go v6.3.0..v6.4.0 (PR body): compare v6.3.0...v6.4.0 ✓ 40000 bytes
    - context: 91927 bytes
- [#294](https://github.com/mudler/yip/pull/294) — ✅ **good** — This is a minor version bump for a dependency. The associated changes appear to be performance optimizations and new tests, which do not introduce any apparent security risks. Therefore, it is safe to auto-approve.
  ↳ This PR updates the dependency gopkg.in/ini.v1 from v1.67.2 to v1.67.3. The update includes an optimization for string handling in the `Key.Strings` method and the addition of a new benchmark test for large arrays.
    - go-ini/ini v1.67.2..v1.67.3 (PR body): compare v1.67.2...v1.67.3 ✓ 29840 bytes
    - context: 33138 bytes
- [#299](https://github.com/mudler/yip/pull/299) — ⚠️ **needs_human_verification** — review endpoint unreachable: Post "http://localhost:8080/v1/chat/completions": context deadline exceeded
    - github.com/google/go-containerregistry 0.21.6→0.21.7: compare v0.21.6...v0.21.7 ✓ 40000 bytes
    - github.com/docker/cli 29.4.3+incompatible→29.5.3+incompatible: compare v29.4.3+incompatible...v29.5.3+incompatible failed/empty (no upstream diff)
    - golang.org/x/mod 0.36.0→0.37.0: compare v0.36.0...v0.37.0 ✓ 15670 bytes
    - context: 70397 bytes
- [#300](https://github.com/mudler/yip/pull/300) — ✅ **good** — This is a necessary dependency update to a newer major version of a core action. The upstream changes explicitly mention security improvements, specifically regarding fork PR handling, which mitigates potential "pwn request" vulnerabilities. Therefore, this change is safe to auto-approve.
  ↳ This PR updates the `actions/checkout` dependency from v6 to v7. This major version upgrade includes important security fixes, such as safer handling of fork pull requests, and updates to several related actions and dependencies.
    - actions/checkout v7.0.0..v7.0.0 (PR body): compare v7.0.0...v7.0.0 failed/empty (no upstream diff)
    - actions/checkout v6.0.3..v7.0.0 (PR body): compare v6.0.3...v7.0.0 ✓ 40000 bytes
    - context: 46073 bytes
- [#304](https://github.com/mudler/yip/pull/304) — ✅ **good** — The dependency update is to a major version (v1 to v2) and the release notes explicitly mention multiple security updates (CVEs). This type of dependency bump, especially when driven by security fixes, is generally safe to auto-approve.
  ↳ This PR updates the `github.com/containerd/containerd` dependency from v1.7.33 to v2.3.2. This upgrade includes several security patches for containerd, addressing multiple CVEs listed in the release notes.
    - github.com/Microsoft/go-winio 0.6.2→0.6.3-0.20251027160822-ad3df93bed29: compare v0.6.2...ad3df93bed29 ✓ 28979 bytes
    - context: 100866 bytes
- [#305](https://github.com/mudler/yip/pull/305) — ✅ **good** — This is a routine dependency update to a minor version of a widely used library (`gomega`) and several related Go packages. The changes do not introduce any obvious security risks or breaking API changes that would require manual review.
  ↳ Update the `github.com/onsi/gomega` dependency to v1.42.1 and update several other Go modules (`golang.org/x/net`, `golang.org/x/text`, `gopkg.in/check.v1`) to their latest versions. This PR also adds new documentation and plugin files for the updated Gomega version.
    - github.com/onsi/gomega 1.41.0→1.42.1: compare v1.41.0...v1.42.1 ✓ 40000 bytes
    - golang.org/x/net 0.55.0→0.56.0: compare v0.55.0...v0.56.0 ✓ 40000 bytes
    - context: 84517 bytes

