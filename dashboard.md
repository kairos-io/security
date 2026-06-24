# Kairos Security Dashboard

_Updated 2026-06-24._

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
| kairos-io/AuroraBoot | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/cluster-api-provider-kairos | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/entangle | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/entangle-proxy | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/go-nodepair | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/go-ukify | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/hadron | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/immucore | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-agent | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-init | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-installer | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-lab | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-must-burn | 0 | 0 | 0 | 0 | 0 | skipped: not source-scannable |
| kairos-io/kairos-operator | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kairos-sdk | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kcrypt | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/kcrypt-discovery-challenger | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/netboot | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/provider-kairos | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/provider-kubernetes | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/simple-mdns-server | 0 | 0 | 0 | 0 | 0 | clean |
| kairos-io/tpm-helpers | 0 | 0 | 0 | 0 | 0 | clean |
| mauromorales/xpasswd | 0 | 0 | 0 | 0 | 0 | clean |
| mudler/edgevpn | 0 | 0 | 0 | 0 | 0 | clean |
| mudler/entities | 0 | 0 | 0 | 0 | 0 | clean |
| mudler/go-pluggable | 0 | 0 | 0 | 0 | 0 | clean |
| mudler/yip | 0 | 0 | 0 | 0 | 0 | clean |

## 📋 Open PRs

_None._

## 🤖 Bot PR ledger

| Repo | Bump | Kind | Source | State | PR |
|---|---|---|---|---|---|
| mudler/edgevpn | golang.org/x/net@0.33.0 | direct | ksec | error | — |

## 🔎 Bot-PR reviews

**kairos-io/AuroraBoot**

- [#409](https://github.com/kairos-io/AuroraBoot/pull/409) — ✅ **good** — This pull request appears to be a standard automated dependency update, primarily updating a dependency digest for github.com/foxboron/sbctl. Dependency updates are routine maintenance and generally safe, provided the updated versions do not introduce known critical vulnerabilities.
  ↳ This PR updates the digest for github.com/foxboron/sbctl to a specific commit hash, and also updates several other dependencies like github.com/fatih/color, golang.org/x/sys, and adds/updates packages such as github.com/go-piv/piv-go/v2 and kernel.org/pub/linux/libs/security/libcap/psx.
- [#566](https://github.com/kairos-io/AuroraBoot/pull/566) — ✅ **good** — This is a routine dependency update for a well-known icon library (`lucide-react`). Since this change is driven by a bot and involves updating a dependency to a newer version, it is generally safe to auto-approve.
  ↳ This PR updates the dependency `lucide-react` from version ^0.468.0 to ^0.577.0. This is a standard dependency maintenance update, likely performed by a bot, to pull in the latest version of the library.
**kairos-io/cluster-api-provider-kairos**

- [#38](https://github.com/kairos-io/cluster-api-provider-kairos/pull/38) — ✅ **good** — This pull request is a routine dependency update for golang.org/x/oauth2. Updating to a newer version is standard practice and generally safe, as it addresses potential minor issues or security patches without introducing significant risk.
**kairos-io/entangle**

- [#10](https://github.com/kairos-io/entangle/pull/10) — ✅ **good** — This pull request only updates several indirect dependencies to newer versions. These types of dependency bumps are routine maintenance and do not introduce new security risks. The changes appear safe to merge automatically.
**kairos-io/entangle-proxy**

- [#5](https://github.com/kairos-io/entangle-proxy/pull/5) — ✅ **good** — This pull request only updates the version of the 'github.com/onsi/gomega' dependency to v1.42.0. This is a routine dependency update, which typically includes bug fixes and security patches, and poses a low risk. It is safe to auto-approve.
- [#6](https://github.com/kairos-io/entangle-proxy/pull/6) — ✅ **good** — This pull request updates the dependency sigs.k8s.io/controller-runtime to version v0.24.1. This is a routine dependency update to a newer version, which is generally safe and necessary for maintaining security and compatibility.
- [#10](https://github.com/kairos-io/entangle-proxy/pull/10) — ✅ **good** — This pull request updates several dependencies, including upgrading github.com/onsi/ginkgo to v2 and updating other related packages like go-logr, go-cmp, and protobuf. These are routine dependency maintenance updates, which generally improve security by incorporating patches for known vulnerabilities. The changes are reflected correctly in both go.mod and go.sum.
- [#14](https://github.com/kairos-io/entangle-proxy/pull/14) — ✅ **good** — This pull request primarily updates several dependencies to newer versions, including core packages like `golang.org/x` and `google.golang.org/protobuf`. Updating dependencies is a crucial security practice to ensure that known vulnerabilities are patched. The changes appear to be dependency hygiene improvements and do not introduce any obvious security risks.
- [#18](https://github.com/kairos-io/entangle-proxy/pull/18) — ✅ **good** — This pull request updates the version of the docker/build-push-action from v2 to v7. Updating dependencies to the latest stable version is a standard security and maintenance practice. This change is safe to auto-approve.
- [#19](https://github.com/kairos-io/entangle-proxy/pull/19) — ✅ **good** — This pull request updates the docker/login-action dependency from v1 to v4. Updating dependencies is a standard maintenance practice that generally improves security and stability by incorporating bug fixes and security patches from the maintainers. This change is safe to auto-approve.
- [#20](https://github.com/kairos-io/entangle-proxy/pull/20) — ✅ **good** — This pull request primarily updates several dependencies, including core Kubernetes libraries (k8s.io/api, k8s.io/client-go, k8s.io/apimachinery) and other related packages, to newer versions. This is a standard maintenance task aimed at applying security patches and leveraging recent features. There are no changes to the application source code itself, making this change safe to auto-approve.
- [#22](https://github.com/kairos-io/entangle-proxy/pull/22) — ✅ **good** — This is a dependency update for a logging library. Updating to a newer patch version (v1.4.3) is generally safe and often includes bug fixes or minor security patches. No immediate security risks are apparent from the change itself.
- [#23](https://github.com/kairos-io/entangle-proxy/pull/23) — ✅ **good** — This pull request updates the dependency for 'actions/checkout' from version v2 to v7 in two workflow files. This is a standard dependency update to a newer version, which is generally safe and beneficial for security and maintenance.
**kairos-io/go-nodepair**

- [#11](https://github.com/kairos-io/go-nodepair/pull/11) — ✅ **good** — This is a standard dependency update to a newer version of a well-known project. The changelog indicates that the upgrade includes fixes and new features, suggesting it is a safe and beneficial update. No immediate security risks are apparent from the provided context.
  ↳ This PR updates the dependency `github.com/ipfs/go-log` from version v1.0.5 to v2.9.2. This upgrade incorporates fixes and new features, including support for `slog.Group` in the zap bridge.
- [#27](https://github.com/kairos-io/go-nodepair/pull/27) — ✅ **good** — The changes involve updating several core dependencies across the project. The changelogs indicate that these updates include important security patches, such as restricting RSA key sizes in go-libp2p and fixing memory exhaustion attacks in quic-go. This is standard maintenance and security hygiene.
  ↳ This pull request updates several core dependencies, including go-libp2p, quic-go, golang.org/x/crypto, golang.org/x/image, golang.org/x/net, and google.golang.org/protobuf. The updates include critical security fixes, such as mitigating a DoS attack in go-libp2p and addressing memory exhaustion issues in quic-go.
- [#37](https://github.com/kairos-io/go-nodepair/pull/37) — ✅ **good** — This appears to be a routine dependency update. Since this is a digest bump for a known dependency, and no specific security issues are indicated in the context, the change is considered safe for auto-approval.
  ↳ This PR updates the dependency github.com/kbinani/screenshot from digest b87d318 to 089614a. This is a standard dependency update managed by Mend Renovate.
- [#44](https://github.com/kairos-io/go-nodepair/pull/44) — ✅ **good** — This is a minor version update for a dependency that appears to be a security scanning tool. Minor version bumps are typically safe and address maintenance or minor bug fixes without introducing significant breaking changes or new security risks.
  ↳ This PR updates the dependency `google/osv-scanner-action` from version v1.8.4 to v1.9.2 in the workflow configuration. This change ensures the project uses the latest version of the OSV scanner action.
- [#46](https://github.com/kairos-io/go-nodepair/pull/46) — ⚠️ **needs_human_verification** — This is a dependency update to a major version bump within the library. While dependency updates are generally safe, the new version introduces significant new features and internal API changes, such as the `--sleep-on-failure` mechanism and changes to suite control logic. A human review is recommended to ensure these changes do not introduce any unexpected breaking behavior in the application code.
  ↳ This PR updates the `github.com/onsi/ginkgo/v2` dependency from v2.29.0 to v2.32.0. This version introduces new features like RSpec-style documentation output (`-fd`) and a debugging aid called `--sleep-on-failure`. It also includes internal refactoring related to suite management and plugin support.
- [#47](https://github.com/kairos-io/go-nodepair/pull/47) — ✅ **good** — This pull request updates the version of the `github.com/onsi/gomega` dependency to v1.42.1 and also updates several related transitive dependencies (e.g., `golang.org/x/crypto`, `golang.org/x/mod`, etc.) to newer versions. This is a routine dependency maintenance task that improves security and stability. No immediate security concerns are identified from the diff.
- [#53](https://github.com/kairos-io/go-nodepair/pull/53) — ✅ **good** — This pull request updates a dependency, specifically the `google/osv-scanner-action`, to a newer version (v2.3.8). Updating dependencies is a standard maintenance practice and generally improves security and stability. There are no suspicious changes in the diff itself.
- [#55](https://github.com/kairos-io/go-nodepair/pull/55) — ✅ **good** — This pull request updates a dependency, github.com/lucasb-eyer/go-colorful, to version v1.4.0. This is a standard dependency update, and without further context indicating known vulnerabilities or breaking changes, it is considered safe to auto-approve.
- [#57](https://github.com/kairos-io/go-nodepair/pull/57) — ✅ **good** — This pull request only updates the version of the 'actions/setup-go' action from v5 to v6. This is a routine maintenance update for a standard dependency and does not introduce any new security risks or significant functional changes that require manual review.
- [#58](https://github.com/kairos-io/go-nodepair/pull/58) — ✅ **good** — This pull request updates the version of the github/codeql-action/upload-sarif action from v3 to v4. Updating dependencies, especially for security scanning tools, is a standard maintenance practice to ensure the latest features, bug fixes, and security patches are included. This change is safe to auto-approve.
- [#59](https://github.com/kairos-io/go-nodepair/pull/59) — ✅ **good** — This pull request only modifies the configuration file for the Renovate bot. The changes involve migrating to a recommended configuration and adjusting package matching rules. This is a standard configuration update and poses no security risk.
- [#62](https://github.com/kairos-io/go-nodepair/pull/62) — ✅ **good** — This pull request only updates the version of the 'actions/checkout' dependency from v4/v2 to v7 in two workflow files. This is a routine dependency update and does not introduce any new security concerns.
**kairos-io/go-ukify**

- [#31](https://github.com/kairos-io/go-ukify/pull/31) — ✅ **good** — This change is a routine dependency update to a specific digest for `github.com/foxboron/go-uefi`. Since this is an automated update, the risk is low, assuming the new digest is a standard patch or minor update. It is safe to auto-approve.
  ↳ This PR updates the dependency `github.com/foxboron/go-uefi` by changing its version digest from `fab4fdf` to `d29549a`. This is a standard dependency update performed by an automated tool.
- [#38](https://github.com/kairos-io/go-ukify/pull/38) — ✅ **good** — This is a routine dependency update for a security scanning tool. The version bump is minor, and the changelog suggests compatibility adjustments rather than breaking changes. Therefore, it is safe to auto-approve this update.
  ↳ This PR updates the `securego/gosec` action dependency from v2.22.5 to v2.27.1. The release notes indicate minor changes, including a downgrade of a Google library to manage Go version compatibility.
- [#39](https://github.com/kairos-io/go-ukify/pull/39) — ✅ **good** — This is a standard dependency bump to a newer version of a widely used testing library. The changes appear to be feature additions and internal maintenance/refactoring, which are generally safe. No immediate security vulnerabilities are apparent in the provided diffs.
  ↳ This PR updates the `github.com/onsi/gomega` dependency from v1.37.0 to v1.42.1. The changes introduce new features such as a Claude plugin, updated development container support, and internal refactoring across documentation and formatting logic.
- [#43](https://github.com/kairos-io/go-ukify/pull/43) — ✅ **good** — This is a routine dependency update for a widely used GitHub Action. Updating to a newer major version is generally safe and recommended to ensure the project benefits from bug fixes and security patches provided by the action maintainers.
  ↳ This PR updates the `actions/setup-go` action from version 5 to version 6 in the release and unit test workflows. This is a routine dependency update to incorporate the latest changes and potential security fixes from the action maintainers.
- [#46](https://github.com/kairos-io/go-ukify/pull/46) — ⚠️ **needs_human_verification** — review endpoint unreachable: Post "http://localhost:8080/v1/chat/completions": context deadline exceeded
- [#47](https://github.com/kairos-io/go-ukify/pull/47) — ✅ **good** — This pull request updates a dependency to a newer major version of a GitHub Action. Dependency updates, especially for security-related tools like CodeQL actions, are generally safe and recommended to incorporate bug fixes and security patches. The change is localized to updating the action version used for uploading SARIF files.
  ↳ This PR updates the `github/codeql-action/upload-sarif` action from version v3 to v4 in the workflow file. This is a standard dependency update to ensure the workflow uses the latest version of the action.
- [#48](https://github.com/kairos-io/go-ukify/pull/48) — ✅ **good** — This is a standard dependency update to a newer version of Ginkgo. The changes primarily involve adopting new features and updating project configuration and tooling to support the new version. Since this is a dependency bump and the changes appear to be related to feature adoption and configuration rather than core logic changes, it is safe to auto-approve.
  ↳ This PR updates the dependency `github.com/onsi/ginkgo/v2` from v2.23.4 to v2.32.0. The changes include updates to configuration files, GitHub Actions workflows, and internal code to integrate new Ginkgo features like RSpec-style documentation output and improved spec filtering. Additionally, several project dependencies in `Gemfile.lock` have been updated to align with the new version.
- [#50](https://github.com/kairos-io/go-ukify/pull/50) — ✅ **good** — This is a standard dependency update for a widely used library, and the release notes indicate bug fixes and enhancements. The accompanying code changes represent internal refactoring of the Viper library itself, which is typical for major version bumps. The dependency bump is safe to auto-approve.
  ↳ This PR updates the `github.com/spf13/viper` dependency to v1.21.0, which includes bug fixes and new features. It also introduces significant internal refactoring to the configuration file finding mechanism using `locafero` and updates the encoding layer interfaces.
- [#51](https://github.com/kairos-io/go-ukify/pull/51) — ✅ **good** — This is a standard dependency update to a newer version of a widely used library. While there are internal breaking changes within the library (like flag option renaming), the PR appears to implement the necessary migration steps, including updating `go.mod` and `go.sum`. The addition of a `SECURITY.md` file is a positive security enhancement. Therefore, this change is safe to auto-approve.
  ↳ This PR updates the `github.com/spf13/cobra` dependency from v1.9.1 to v1.10.2, which includes internal refactors, updated documentation, and dependency bumps for `pflag` and `go.yaml.in/yaml/v3`. The update involves renaming a flag parsing option and updating several internal code paths.
- [#53](https://github.com/kairos-io/go-ukify/pull/53) — ✅ **good** — This is a standard dependency update to a newer version of a well-known library. The changes involve feature additions (HMAC Session support) and bug fixes, which are generally beneficial. The extensive internal refactoring and added tests suggest a thorough update process. Therefore, this change is safe to auto-approve.
  ↳ This PR updates the `go-tpm` dependency from v0.9.5 to v0.9.8, introducing support for HMAC Sessions in ReadPublic and fixing a typo in the TPM hierarchy name. It also includes significant internal refactoring in marshalling and reflection logic, along with new tests for audit sessions and time retrieval.
**kairos-io/kairos**

- [#4104](https://github.com/kairos-io/kairos/pull/4104) — ⚠️ **needs_human_verification** — The PR title suggests an automation task related to dependency upgrades. While automation can be beneficial, security review requires inspecting the actual code changes to ensure no unintended side effects or vulnerabilities were introduced during the pipeline wiring.
  ↳ The PR aims to automate the process of fetching the latest release and running validation tests within the upgrade pipeline.
**kairos-io/kairos-operator**

- [#114](https://github.com/kairos-io/kairos-operator/pull/114) — ✅ **good** — This change is a routine dependency update for a standard GitHub Action. Updating the digest to a newer version is a maintenance task and does not introduce any new security risks. It is safe to auto-approve.
  ↳ This PR updates the digest for the `actions/checkout` dependency from `de0fac2` to `df4cb1c`. This is a routine dependency update for a standard GitHub Action and does not introduce any new security vulnerabilities.
- [#115](https://github.com/kairos-io/kairos-operator/pull/115) — ✅ **good** — This change is a routine update of a dependency's digest, which is a standard practice for ensuring build reproducibility or applying minor security patches associated with the base image. It does not introduce any new code, logic changes, or new dependencies, making it safe for automatic approval.
  ↳ This PR updates the Dockerfile and Dockerfile.node-labeler files to use a new digest (`478231bfd9677835606c249208483a3c43b31e941c1040c48747b111c7ab871c`) for the `docker.io/golang:1.26.4` image. This is a routine update of the image digest.

