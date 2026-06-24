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

- [#409](https://github.com/kairos-io/AuroraBoot/pull/409) — ✅ **good** — This pull request only updates the digest for a dependency (`github.com/foxboron/sbctl`) and updates several indirect dependencies in go.mod and go.sum files. This is a routine dependency maintenance task, and there are no apparent security risks introduced by these changes.
- [#566](https://github.com/kairos-io/AuroraBoot/pull/566) — ✅ **good** — This pull request is a routine dependency update for 'lucide-react'. Updating dependencies is a standard maintenance task, and this change appears safe as it only updates the version number.
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

- [#47](https://github.com/kairos-io/go-nodepair/pull/47) — ✅ **good** — This pull request updates the version of the `github.com/onsi/gomega` dependency to v1.42.1 and also updates several related transitive dependencies (e.g., `golang.org/x/crypto`, `golang.org/x/mod`, etc.) to newer versions. This is a routine dependency maintenance task that improves security and stability. No immediate security concerns are identified from the diff.
- [#53](https://github.com/kairos-io/go-nodepair/pull/53) — ✅ **good** — This pull request updates a dependency, specifically the `google/osv-scanner-action`, to a newer version (v2.3.8). Updating dependencies is a standard maintenance practice and generally improves security and stability. There are no suspicious changes in the diff itself.
- [#55](https://github.com/kairos-io/go-nodepair/pull/55) — ✅ **good** — This pull request updates a dependency, github.com/lucasb-eyer/go-colorful, to version v1.4.0. This is a standard dependency update, and without further context indicating known vulnerabilities or breaking changes, it is considered safe to auto-approve.
- [#57](https://github.com/kairos-io/go-nodepair/pull/57) — ✅ **good** — This pull request only updates the version of the 'actions/setup-go' action from v5 to v6. This is a routine maintenance update for a standard dependency and does not introduce any new security risks or significant functional changes that require manual review.
- [#58](https://github.com/kairos-io/go-nodepair/pull/58) — ✅ **good** — This pull request updates the version of the github/codeql-action/upload-sarif action from v3 to v4. Updating dependencies, especially for security scanning tools, is a standard maintenance practice to ensure the latest features, bug fixes, and security patches are included. This change is safe to auto-approve.
- [#59](https://github.com/kairos-io/go-nodepair/pull/59) — ✅ **good** — This pull request only modifies the configuration file for the Renovate bot. The changes involve migrating to a recommended configuration and adjusting package matching rules. This is a standard configuration update and poses no security risk.
- [#62](https://github.com/kairos-io/go-nodepair/pull/62) — ✅ **good** — This pull request only updates the version of the 'actions/checkout' dependency from v4/v2 to v7 in two workflow files. This is a routine dependency update and does not introduce any new security concerns.

