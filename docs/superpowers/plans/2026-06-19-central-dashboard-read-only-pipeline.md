# Central Dashboard — Read-Only Pipeline Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the `ksec` Go CLI that discovers all maintained repos, collects their security findings (open PRs, image CVEs, source/dep CVEs, GitHub alerts), correlates them into a waterfall view, triages them with a self-hosted AI, and renders a committed dashboard plus a tracking issue in `kairos-io/kairos` — with **no writes to other repos** (remediation is a separate Plan 2).

**Architecture:** A single Go CLI of phased subcommands (`discover → collect → correlate → triage → render`). Each phase reads/writes pretty-printed JSON state files under `state/`, making every phase independently runnable, testable, and auditable in git. Deterministic Go does all data work; a small self-hosted LocalAI model (driven by `nib --yolo`) is used only for the triage narrative/summaries, with a deterministic fallback when the model is unavailable.

**Tech Stack:** Go (1.22+), `cobra` (CLI), `gopkg.in/yaml.v3` (config), `stretchr/testify` (tests), `gh` CLI (GitHub API access), `trivy`/`govulncheck` (scanners), LocalAI + `nib` (AI). Orchestrated by GitHub Actions.

## Global Constraints

- Module path: `github.com/kairos-io/security`. Binary name: `ksec`.
- Go version floor: **1.22**.
- Every phase is invoked as `ksec <phase> --state-dir <dir> [--dry-run]`. `--state-dir` defaults to `./state`.
- State files are pretty-printed JSON (2-space indent), arrays sorted by a stable key, so git diffs are clean. Map keys are emitted sorted (Go default for `map[string]...`).
- A phase NEVER mutates a state file owned by an earlier phase; it only reads earlier files and writes its own.
- **Read-only on other repos.** The only GitHub *write* in this plan is upserting the single tracking issue in `kairos-io/kairos`. Dry-run turns that write into a printed plan.
- Per-repo / per-collector failures are isolated: record a `CollectionError`, never abort the run.
- AI is best-effort: on any failure, fall back to deterministic triage and set `aiAvailable: false`.
- Secrets (tokens) are never logged.
- Use `testify` (`require` for fatal preconditions, `assert` for value checks). Golden-file tests compare against committed `testdata/*.golden` files.

---

## File structure

```
go.mod
cmd/ksec/main.go                 # cobra root + global flags + subcommand wiring
internal/state/types.go          # all state structs
internal/state/store.go          # generic Load/Save (stable JSON)
internal/state/store_test.go
internal/config/config.go        # repos.yaml + ai.yaml loading
internal/config/config_test.go
internal/ghclient/ghclient.go    # GitHub interface + gh-CLI implementation
internal/ghclient/fake.go        # in-memory fake for tests
internal/discover/parse.go       # kairos-init Makefile/go.mod parsing (pure)
internal/discover/parse_test.go
internal/discover/discover.go    # org enum + parse + repos.yaml merge
internal/discover/discover_test.go
internal/collect/collect.go      # Collector interface + Run fan-out + error isolation
internal/collect/collect_test.go
internal/collect/source.go       # sourceCVE collector (govulncheck JSON)
internal/collect/source_test.go
internal/collect/image.go        # imageCVE collector (trivy JSON)
internal/collect/image_test.go
internal/collect/prs.go          # open security PRs collector
internal/collect/prs_test.go
internal/collect/alerts.go       # GitHub security alerts collector
internal/collect/alerts_test.go
internal/correlate/correlate.go  # pure dedupe + waterfall grouping
internal/correlate/correlate_test.go
internal/triage/triage.go        # AIClient interface + fallback + orchestration
internal/triage/nib.go           # nib --yolo AIClient implementation
internal/triage/triage_test.go
internal/render/render.go        # dashboard.json + dashboard.md
internal/render/render_test.go
internal/render/testdata/*.golden
internal/render/issue.go         # tracking-issue upsert
internal/render/issue_test.go
repos.yaml                       # hybrid repo overrides (sample committed)
ai.yaml                          # LocalAI + nib handles (sample committed)
.github/workflows/security-dashboard.yaml
```

---

### Task 1: Project scaffold + CLI skeleton

**Files:**
- Create: `go.mod`
- Create: `cmd/ksec/main.go`

**Interfaces:**
- Consumes: nothing.
- Produces: a `ksec` binary with global flags `--state-dir` (string, default `./state`) and `--dry-run` (bool), and five subcommands `discover`, `collect`, `correlate`, `triage`, `render`, each currently printing `"<phase>: not implemented"`. Later tasks replace each subcommand's `RunE`.

- [ ] **Step 1: Initialize the module and add cobra**

```bash
cd /home/mudler/_git/kairos-security
go mod init github.com/kairos-io/security
go get github.com/spf13/cobra@latest
```

- [ ] **Step 2: Write the CLI skeleton**

Create `cmd/ksec/main.go`:

```go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Flags shared by every phase.
type globalFlags struct {
	stateDir string
	dryRun   bool
}

func newRootCmd() *cobra.Command {
	gf := &globalFlags{}
	root := &cobra.Command{
		Use:   "ksec",
		Short: "Kairos central security dashboard engine",
	}
	root.PersistentFlags().StringVar(&gf.stateDir, "state-dir", "./state", "directory holding committed state JSON")
	root.PersistentFlags().BoolVar(&gf.dryRun, "dry-run", false, "print intended writes instead of performing them")

	for _, phase := range []string{"discover", "collect", "correlate", "triage", "render"} {
		p := phase
		root.AddCommand(&cobra.Command{
			Use:   p,
			Short: "run the " + p + " phase",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Printf("%s: not implemented\n", p)
				return nil
			},
		})
	}
	return root
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Build and smoke-test**

Run: `go build ./... && go run ./cmd/ksec discover --state-dir ./state`
Expected: prints `discover: not implemented`, exit 0. `go run ./cmd/ksec --help` lists all five subcommands.

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum cmd/ksec/main.go
git commit -m "feat: ksec CLI skeleton with phase subcommands"
```

---

### Task 2: State package — types + stable Load/Save

**Files:**
- Create: `internal/state/types.go`
- Create: `internal/state/store.go`
- Test: `internal/state/store_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces:
  - Types: `Artifact`, `Repo`, `Finding`, `CollectionError`, `Findings`, `Bump`, `WaterfallGroup`, `Correlated`, `Triage` (fields below).
  - `func Save[T any](dir, name string, v T) error` — marshals `v` as 2-space-indented JSON to `filepath.Join(dir, name)`, creating `dir` if needed, with a trailing newline.
  - `func Load[T any](dir, name string, v *T) error` — reads and unmarshals; returns the underlying error if the file is missing.
  - File name constants: `ReposFile="repos.json"`, `FindingsFile="findings.json"`, `CorrelatedFile="correlated.json"`, `TriageFile="triage.json"`.

- [ ] **Step 1: Write the types**

Create `internal/state/types.go`:

```go
package state

// File name constants for each phase's output.
const (
	ReposFile      = "repos.json"
	FindingsFile   = "findings.json"
	CorrelatedFile = "correlated.json"
	TriageFile     = "triage.json"
)

type Artifact struct {
	Type    string `json:"type"`              // "image" | "go"
	Ref     string `json:"ref,omitempty"`     // image reference, when Type=="image"
	ModPath string `json:"modpath,omitempty"` // module path within repo, when Type=="go"
}

type Repo struct {
	Repo        string     `json:"repo"`        // "owner/name"
	Kind        string     `json:"kind"`        // "org" | "dep" | "external"
	Branch      string     `json:"branch"`
	Criticality string     `json:"criticality"` // "low" | "medium" | "high"
	Artifacts   []Artifact `json:"artifacts"`
}

type Finding struct {
	ID             string   `json:"id"`              // stable dedupe key
	Repo           string   `json:"repo"`
	Type           string   `json:"type"`            // "pr" | "imageCVE" | "sourceCVE" | "ghAlert"
	CVEID          string   `json:"cveID,omitempty"`
	GHSA           string   `json:"ghsa,omitempty"`
	Ecosystem      string   `json:"ecosystem,omitempty"`
	Package        string   `json:"package,omitempty"`
	CurrentVersion string   `json:"currentVersion,omitempty"`
	FixedVersion   string   `json:"fixedVersion,omitempty"`
	Severity       string   `json:"severity"`        // critical|high|medium|low|unknown
	Source         string   `json:"source"`          // tool/api that produced it
	Title          string   `json:"title,omitempty"`
	URL            string   `json:"url,omitempty"`
	FirstSeen      string   `json:"firstSeen"`       // YYYY-MM-DD
	LastSeen       string   `json:"lastSeen"`        // YYYY-MM-DD
}

type CollectionError struct {
	Repo      string `json:"repo"`
	Collector string `json:"collector"`
	Message   string `json:"message"`
}

// Findings is the collect phase output: findings plus non-fatal errors.
type Findings struct {
	Findings []Finding         `json:"findings"`
	Errors   []CollectionError `json:"errors"`
}

type Bump struct {
	Package string `json:"package"`
	To      string `json:"to"`
}

type WaterfallGroup struct {
	ID            string   `json:"id"`
	RootCause     string   `json:"rootCause"`
	Ecosystem     string   `json:"ecosystem"`
	Severity      string   `json:"severity"`
	AffectedRepos []string `json:"affectedRepos"`
	SuggestedBump Bump     `json:"suggestedBump"`
}

type Correlated struct {
	Findings  []Finding        `json:"findings"`
	Waterfall []WaterfallGroup `json:"waterfall"`
}

type Triage struct {
	GeneratedAt string            `json:"generatedAt"`
	Model       string            `json:"model"`
	AIAvailable bool              `json:"aiAvailable"`
	Focus       []string          `json:"focus"`
	Summaries   map[string]string `json:"summaries"`
	Narrative   string            `json:"narrative"`
}
```

- [ ] **Step 2: Write the failing test**

Create `internal/state/store_test.go`:

```go
package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	in := Findings{
		Findings: []Finding{{ID: "a", Repo: "kairos-io/x", Type: "sourceCVE", Severity: "high"}},
		Errors:   []CollectionError{{Repo: "kairos-io/y", Collector: "prs", Message: "boom"}},
	}
	require.NoError(t, Save(dir, FindingsFile, in))

	var out Findings
	require.NoError(t, Load(dir, FindingsFile, &out))
	assert.Equal(t, in, out)
}

func TestSaveIsStableIndentedJSON(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, Save(dir, RingFileForTest, map[string]int{"b": 2, "a": 1}))
	b, err := os.ReadFile(filepath.Join(dir, RingFileForTest))
	require.NoError(t, err)
	// keys sorted, 2-space indent, trailing newline
	assert.Equal(t, "{\n  \"a\": 1,\n  \"b\": 2\n}\n", string(b))
}

const RingFileForTest = "scratch.json"

func TestLoadMissingFileErrors(t *testing.T) {
	var out Findings
	assert.Error(t, Load(t.TempDir(), "nope.json", &out))
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./internal/state/...`
Expected: FAIL — `Save`/`Load` undefined.

- [ ] **Step 4: Implement the store**

Create `internal/state/store.go`:

```go
package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Save[T any](dir, name string, v T) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(filepath.Join(dir, name), b, 0o644)
}

func Load[T any](dir, name string, v *T) error {
	b, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./internal/state/...`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/state/
git commit -m "feat: state types and stable JSON store"
```

---

### Task 3: Config package — repos.yaml + ai.yaml

**Files:**
- Create: `internal/config/config.go`
- Test: `internal/config/config_test.go`
- Create: `repos.yaml` (sample)
- Create: `ai.yaml` (sample)

**Interfaces:**
- Consumes: `state.Repo` from Task 2.
- Produces:
  - `type ReposConfig struct { Repos []state.Repo; Exclude []string }` (yaml `repos`, `exclude`).
  - `type AIConfig struct { LocalAI LocalAICfg; Nib NibCfg }` with the fields shown below.
  - `func LoadRepos(path string) (ReposConfig, error)` — returns a zero-value config (no error) if the file does not exist.
  - `func LoadAI(path string) (AIConfig, error)` — same missing-file behavior; applies env overrides `LOCALAI_URL`, `LOCALAI_MODEL`.

- [ ] **Step 1: Write the failing test**

Create `internal/config/config_test.go`:

```go
package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadReposParsesAndDefaultsMissing(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "repos.yaml")
	require.NoError(t, os.WriteFile(p, []byte(`
repos:
  - repo: mudler/edgevpn
    kind: external
    branch: master
    criticality: high
    artifacts:
      - type: go
        modpath: .
exclude:
  - kairos-io/some-archive
`), 0o644))

	cfg, err := LoadRepos(p)
	require.NoError(t, err)
	require.Len(t, cfg.Repos, 1)
	assert.Equal(t, "mudler/edgevpn", cfg.Repos[0].Repo)
	assert.Equal(t, "external", cfg.Repos[0].Kind)
	assert.Equal(t, []string{"kairos-io/some-archive"}, cfg.Exclude)

	missing, err := LoadRepos(filepath.Join(dir, "nope.yaml"))
	require.NoError(t, err)
	assert.Empty(t, missing.Repos)
}

func TestLoadAIAppliesEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ai.yaml")
	require.NoError(t, os.WriteFile(p, []byte(`
localai:
  endpoint: http://localhost:8080
  model:
    name: base-model
nib:
  mode: yolo
`), 0o644))
	t.Setenv("LOCALAI_URL", "http://override:9000")
	t.Setenv("LOCALAI_MODEL", "override-model")

	cfg, err := LoadAI(p)
	require.NoError(t, err)
	assert.Equal(t, "http://override:9000", cfg.LocalAI.Endpoint)
	assert.Equal(t, "override-model", cfg.LocalAI.Model.Name)
	assert.Equal(t, "yolo", cfg.Nib.Mode)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/config/...`
Expected: FAIL — package/functions undefined.

- [ ] **Step 3: Implement config**

```bash
go get gopkg.in/yaml.v3
```

Create `internal/config/config.go`:

```go
package config

import (
	"errors"
	"io/fs"
	"os"

	"github.com/kairos-io/security/internal/state"
	"gopkg.in/yaml.v3"
)

type ReposConfig struct {
	Repos   []state.Repo `yaml:"repos"`
	Exclude []string     `yaml:"exclude"`
}

type ModelCfg struct {
	Name    string `yaml:"name"`
	Gallery string `yaml:"gallery"`
	Quant   string `yaml:"quant"`
}

type LocalAICfg struct {
	Version        string   `yaml:"version"`
	Endpoint       string   `yaml:"endpoint"`
	StartupTimeout string   `yaml:"startupTimeout"`
	Model          ModelCfg `yaml:"model"`
}

type NibCfg struct {
	Version     string  `yaml:"version"`
	Mode        string  `yaml:"mode"`
	Model       string  `yaml:"model"`
	Endpoint    string  `yaml:"endpoint"`
	MaxTokens   int     `yaml:"maxTokens"`
	Temperature float64 `yaml:"temperature"`
}

type AIConfig struct {
	LocalAI LocalAICfg `yaml:"localai"`
	Nib     NibCfg     `yaml:"nib"`
}

func readYAML[T any](path string, v *T) error {
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil // missing file → zero value, no error
	}
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

func LoadRepos(path string) (ReposConfig, error) {
	var cfg ReposConfig
	return cfg, readYAML(path, &cfg)
}

func LoadAI(path string) (AIConfig, error) {
	var cfg AIConfig
	if err := readYAML(path, &cfg); err != nil {
		return cfg, err
	}
	if v := os.Getenv("LOCALAI_URL"); v != "" {
		cfg.LocalAI.Endpoint = v
	}
	if v := os.Getenv("LOCALAI_MODEL"); v != "" {
		cfg.LocalAI.Model.Name = v
	}
	// nib defaults derive from localai so they cannot drift
	if cfg.Nib.Endpoint == "" {
		cfg.Nib.Endpoint = cfg.LocalAI.Endpoint
	}
	if cfg.Nib.Model == "" {
		cfg.Nib.Model = cfg.LocalAI.Model.Name
	}
	return cfg, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/config/...`
Expected: PASS.

- [ ] **Step 5: Add the committed sample config files**

Create `repos.yaml`:

```yaml
# Hybrid repo overrides. Auto-discovery (kairos-io org + kairos-init deps) is the
# base; this file ADDS external repos, EXCLUDES repos, and overrides metadata.
repos:
  - repo: mudler/edgevpn
    kind: external
    branch: master
    criticality: high
    artifacts:
      - type: go
        modpath: .
  - repo: mudler/yip
    kind: external
    branch: master
    criticality: medium
    artifacts:
      - type: go
        modpath: .
exclude: []
```

Create `ai.yaml`:

```yaml
localai:
  version: "latest"            # overridable: LOCALAI_VERSION
  endpoint: "http://localhost:8080"   # overridable: LOCALAI_URL
  startupTimeout: "5m"
  model:
    name: "small-instruct"     # overridable: LOCALAI_MODEL
    gallery: ""                # gallery entry / model URI to preload
    quant: ""
nib:
  version: "latest"            # overridable: NIB_VERSION
  mode: "yolo"
  model: ""                    # defaults to localai.model.name
  endpoint: ""                 # defaults to localai.endpoint
  maxTokens: 4096
  temperature: 0.2
```

- [ ] **Step 6: Commit**

```bash
git add internal/config/ repos.yaml ai.yaml
git commit -m "feat: repos.yaml and ai.yaml config loading"
```

---

### Task 4: GitHub client interface + gh-CLI implementation + fake

**Files:**
- Create: `internal/ghclient/ghclient.go`
- Create: `internal/ghclient/fake.go`
- Test: `internal/ghclient/fake_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces a `GitHub` interface used by discover/collect/render. Real impl shells out to `gh` (already available in CI and respects `GH_TOKEN`/`GITHUB_TOKEN`). The `Fake` is the test double every other package uses.

```go
type PullRequest struct{ Number int; Title, Author, URL string; Labels []string }
type Alert struct{ Number int; CVEID, GHSA, Package, Ecosystem, Severity, URL, FixedVersion string }

type GitHub interface {
	ListOrgRepos(org string) ([]string, error)               // "owner/name"
	GetFile(repo, path, ref string) ([]byte, error)          // raw file bytes
	ListOpenPRs(repo string) ([]PullRequest, error)
	ListDependabotAlerts(repo string) ([]Alert, error)
	UpsertIssue(repo, marker, title, body string, labels []string) (int, error)
}
```

- [ ] **Step 1: Write the failing test (covers the Fake)**

Create `internal/ghclient/fake_test.go`:

```go
package ghclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakeUpsertCreatesThenUpdates(t *testing.T) {
	f := NewFake()
	n, err := f.UpsertIssue("kairos-io/kairos", "<!-- ksec:dashboard -->", "Security", "body v1", []string{"security"})
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	n2, err := f.UpsertIssue("kairos-io/kairos", "<!-- ksec:dashboard -->", "Security", "body v2", []string{"security"})
	require.NoError(t, err)
	assert.Equal(t, 1, n2, "same marker reuses the issue")
	assert.Equal(t, "body v2", f.Issues["kairos-io/kairos"].Body)
}

func TestFakeListOrgRepos(t *testing.T) {
	f := NewFake()
	f.OrgRepos["kairos-io"] = []string{"kairos-io/kairos", "kairos-io/immucore"}
	got, err := f.ListOrgRepos("kairos-io")
	require.NoError(t, err)
	assert.Equal(t, []string{"kairos-io/kairos", "kairos-io/immucore"}, got)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/ghclient/...`
Expected: FAIL — undefined.

- [ ] **Step 3: Implement the interface + fake + gh-CLI client**

Create `internal/ghclient/ghclient.go`:

```go
package ghclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type PullRequest struct {
	Number int      `json:"number"`
	Title  string   `json:"title"`
	Author string   `json:"author"`
	URL    string   `json:"url"`
	Labels []string `json:"labels"`
}

type Alert struct {
	Number       int    `json:"number"`
	CVEID        string `json:"cveID"`
	GHSA         string `json:"ghsa"`
	Package      string `json:"package"`
	Ecosystem    string `json:"ecosystem"`
	Severity     string `json:"severity"`
	URL          string `json:"url"`
	FixedVersion string `json:"fixedVersion"`
}

type GitHub interface {
	ListOrgRepos(org string) ([]string, error)
	GetFile(repo, path, ref string) ([]byte, error)
	ListOpenPRs(repo string) ([]PullRequest, error)
	ListDependabotAlerts(repo string) ([]Alert, error)
	UpsertIssue(repo, marker, title, body string, labels []string) (int, error)
}

// CLI is the production GitHub client; it shells out to `gh`.
type CLI struct {
	run func(args ...string) ([]byte, error)
}

func NewCLI() *CLI {
	return &CLI{run: func(args ...string) ([]byte, error) {
		cmd := exec.Command("gh", args...)
		var out, errb bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &errb
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("gh %v: %v: %s", args, err, errb.String())
		}
		return out.Bytes(), nil
	}}
}

func (c *CLI) api(path string, jqOrFields ...string) ([]byte, error) {
	args := append([]string{"api", path}, jqOrFields...)
	return c.run(args...)
}

func (c *CLI) ListOrgRepos(org string) ([]string, error) {
	b, err := c.run("repo", "list", org, "--no-archived", "--limit", "1000", "--json", "nameWithOwner", "-q", ".[].nameWithOwner")
	if err != nil {
		return nil, err
	}
	return splitLines(b), nil
}

func (c *CLI) GetFile(repo, path, ref string) ([]byte, error) {
	// gh api returns the raw content with the proper Accept header.
	return c.run("api", fmt.Sprintf("repos/%s/contents/%s?ref=%s", repo, path, ref),
		"-H", "Accept: application/vnd.github.raw+json")
}

func (c *CLI) ListOpenPRs(repo string) ([]PullRequest, error) {
	b, err := c.run("pr", "list", "-R", repo, "--state", "open", "--limit", "200",
		"--json", "number,title,author,url,labels",
		"-q", "[.[] | {number, title, author: .author.login, url, labels: [.labels[].name]}]")
	if err != nil {
		return nil, err
	}
	var prs []PullRequest
	return prs, json.Unmarshal(b, &prs)
}

func (c *CLI) ListDependabotAlerts(repo string) ([]Alert, error) {
	b, err := c.api(fmt.Sprintf("repos/%s/dependabot/alerts?state=open&per_page=100", repo),
		"-q", "[.[] | {number, cveID: (.security_advisory.cve_id // \"\"), ghsa: .security_advisory.ghsa_id, package: .dependency.package.name, ecosystem: .dependency.package.ecosystem, severity: .security_advisory.severity, url: .html_url, fixedVersion: (.security_vulnerability.first_patched_version.identifier // \"\")}]")
	if err != nil {
		return nil, err
	}
	var alerts []Alert
	return alerts, json.Unmarshal(b, &alerts)
}

func (c *CLI) UpsertIssue(repo, marker, title, body string, labels []string) (int, error) {
	full := body + "\n\n" + marker
	// Find an existing issue containing the marker.
	listed, err := c.run("issue", "list", "-R", repo, "--state", "open", "--search", marker, "--limit", "1", "--json", "number", "-q", ".[].number")
	if err != nil {
		return 0, err
	}
	if lines := splitLines(listed); len(lines) > 0 {
		var n int
		fmt.Sscanf(lines[0], "%d", &n)
		_, err := c.run("issue", "edit", fmt.Sprint(n), "-R", repo, "--body", full)
		return n, err
	}
	args := []string{"issue", "create", "-R", repo, "--title", title, "--body", full}
	for _, l := range labels {
		args = append(args, "--label", l)
	}
	out, err := c.run(args...)
	if err != nil {
		return 0, err
	}
	return parseIssueNumberFromURL(out), nil
}

func splitLines(b []byte) []string {
	var out []string
	for _, line := range bytes.Split(bytes.TrimSpace(b), []byte("\n")) {
		if s := string(bytes.TrimSpace(line)); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func parseIssueNumberFromURL(b []byte) int {
	// `gh issue create` prints the new issue URL ending in /<number>.
	s := string(bytes.TrimSpace(b))
	var n int
	if i := bytes.LastIndexByte([]byte(s), '/'); i >= 0 {
		fmt.Sscanf(s[i+1:], "%d", &n)
	}
	return n
}
```

Create `internal/ghclient/fake.go`:

```go
package ghclient

type FakeIssue struct {
	Number int
	Title  string
	Body   string
	Labels []string
}

// Fake is an in-memory GitHub double for tests.
type Fake struct {
	OrgRepos map[string][]string
	Files    map[string][]byte // key: repo|path|ref
	PRs      map[string][]PullRequest
	Alerts   map[string][]Alert
	Issues   map[string]*FakeIssue // key: repo
	nextNum  int
}

func NewFake() *Fake {
	return &Fake{
		OrgRepos: map[string][]string{},
		Files:    map[string][]byte{},
		PRs:      map[string][]PullRequest{},
		Alerts:   map[string][]Alert{},
		Issues:   map[string]*FakeIssue{},
	}
}

func (f *Fake) ListOrgRepos(org string) ([]string, error) { return f.OrgRepos[org], nil }
func (f *Fake) GetFile(repo, path, ref string) ([]byte, error) {
	return f.Files[repo+"|"+path+"|"+ref], nil
}
func (f *Fake) ListOpenPRs(repo string) ([]PullRequest, error)        { return f.PRs[repo], nil }
func (f *Fake) ListDependabotAlerts(repo string) ([]Alert, error)     { return f.Alerts[repo], nil }

func (f *Fake) UpsertIssue(repo, marker, title, body string, labels []string) (int, error) {
	if iss, ok := f.Issues[repo]; ok {
		iss.Body, iss.Title, iss.Labels = body, title, labels
		return iss.Number, nil
	}
	f.nextNum++
	f.Issues[repo] = &FakeIssue{Number: f.nextNum, Title: title, Body: body, Labels: labels}
	return f.nextNum, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/ghclient/...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/ghclient/
git commit -m "feat: GitHub client interface, gh-CLI impl, and test fake"
```

---

### Task 5: `discover` phase — kairos-init parsing + org enum + merge

**Files:**
- Create: `internal/discover/parse.go`
- Test: `internal/discover/parse_test.go`
- Create: `internal/discover/discover.go`
- Test: `internal/discover/discover_test.go`
- Modify: `cmd/ksec/main.go` (wire the `discover` subcommand)

**Interfaces:**
- Consumes: `ghclient.GitHub`, `config.ReposConfig`, `state.Repo`.
- Produces:
  - `func ParseDeps(makefile, gomod []byte) []string` — returns sorted unique `owner/name` slugs found in the kairos-init Makefile (lines like `AGENT_VERSION?=...` mapped via a fixed component→slug table) and go.mod (require lines for `github.com/kairos-io/*`, `github.com/mudler/*`, `github.com/mauromorales/*`).
  - `func Run(gh ghclient.GitHub, cfg config.ReposConfig, org, initRepo, initRef string) ([]state.Repo, error)` — enumerate org repos, parse kairos-init deps, merge with config (apply excludes, add config repos, override metadata by `.Repo`), default missing `Branch="main"`, `Criticality="medium"`, `Kind` (org if in org list else dep/external). Returns repos sorted by `.Repo`.

- [ ] **Step 1: Write the failing parse test**

Create `internal/discover/parse_test.go`:

```go
package discover

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDeps(t *testing.T) {
	makefile := []byte("AGENT_VERSION?=v1.2.3\nIMMUCORE_VERSION?=v0.5.0\nEDGEVPN_VERSION?=v0.30.0\n")
	gomod := []byte(`
module github.com/kairos-io/kairos-init
go 1.22
require (
	github.com/kairos-io/kairos-sdk v0.7.0
	github.com/mudler/yip v1.9.0
	github.com/mauromorales/xpasswd v0.3.0
	github.com/spf13/cobra v1.8.0
)
`)
	got := ParseDeps(makefile, gomod)
	assert.Equal(t, []string{
		"kairos-io/immucore",
		"kairos-io/kairos-agent",
		"kairos-io/kairos-sdk",
		"mauromorales/xpasswd",
		"mudler/edgevpn",
		"mudler/yip",
	}, got)
}
```

- [ ] **Step 2: Run it; expect FAIL** — `ParseDeps` undefined. Run: `go test ./internal/discover/...`

- [ ] **Step 3: Implement the parser**

Create `internal/discover/parse.go`:

```go
package discover

import (
	"regexp"
	"sort"
	"strings"
)

// Fixed component → repo slug map for kairos-init Makefile *_VERSION vars.
var makefileComponents = map[string]string{
	"AGENT":                       "kairos-io/kairos-agent",
	"IMMUCORE":                    "kairos-io/immucore",
	"KCRYPT_DISCOVERY_CHALLENGER": "kairos-io/kcrypt-discovery-challenger",
	"PROVIDER_KAIROS":             "kairos-io/provider-kairos",
	"EDGEVPN":                     "mudler/edgevpn",
}

var (
	reMakeVar = regexp.MustCompile(`(?m)^([A-Z_]+)_VERSION\??=`)
	reGoMod   = regexp.MustCompile(`(?m)^\s*github\.com/(kairos-io|mudler|mauromorales)/([A-Za-z0-9._-]+)\s+v`)
)

// ParseDeps returns sorted unique owner/name slugs from the kairos-init
// Makefile and go.mod.
func ParseDeps(makefile, gomod []byte) []string {
	set := map[string]struct{}{}
	for _, m := range reMakeVar.FindAllStringSubmatch(string(makefile), -1) {
		if slug, ok := makefileComponents[m[1]]; ok {
			set[slug] = struct{}{}
		}
	}
	for _, m := range reGoMod.FindAllStringSubmatch(string(gomod), -1) {
		set[m[1]+"/"+strings.TrimSuffix(m[2], "/")] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
```

- [ ] **Step 4: Run it; expect PASS.** Run: `go test ./internal/discover/...`

- [ ] **Step 5: Write the failing discover test**

Create `internal/discover/discover_test.go`:

```go
package discover

import (
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunMergesOrgDepsAndConfig(t *testing.T) {
	gh := ghclient.NewFake()
	gh.OrgRepos["kairos-io"] = []string{"kairos-io/kairos", "kairos-io/immucore", "kairos-io/archived"}
	gh.Files["kairos-io/kairos-init|Makefile|main"] = []byte("AGENT_VERSION?=v1\n")
	gh.Files["kairos-io/kairos-init|go.mod|main"] = []byte("require github.com/mudler/yip v1.0.0\n")

	cfg := config.ReposConfig{
		Repos:   []state.Repo{{Repo: "mudler/edgevpn", Kind: "external", Branch: "master", Criticality: "high"}},
		Exclude: []string{"kairos-io/archived"},
	}

	repos, err := Run(gh, cfg, "kairos-io", "kairos-io/kairos-init", "main")
	require.NoError(t, err)

	names := map[string]state.Repo{}
	for _, r := range repos {
		names[r.Repo] = r
	}
	assert.Contains(t, names, "kairos-io/kairos")
	assert.Contains(t, names, "kairos-io/immucore")
	assert.Contains(t, names, "kairos-io/kairos-agent") // from Makefile
	assert.Contains(t, names, "mudler/yip")             // from go.mod
	assert.Contains(t, names, "mudler/edgevpn")         // from config
	assert.NotContains(t, names, "kairos-io/archived")  // excluded
	assert.Equal(t, "high", names["mudler/edgevpn"].Criticality)
	assert.Equal(t, "org", names["kairos-io/kairos"].Kind)
	// sorted output
	assert.True(t, repos[0].Repo < repos[len(repos)-1].Repo)
}
```

- [ ] **Step 6: Run it; expect FAIL** — `Run` undefined.

- [ ] **Step 7: Implement discover orchestration**

Create `internal/discover/discover.go`:

```go
package discover

import (
	"sort"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

func Run(gh ghclient.GitHub, cfg config.ReposConfig, org, initRepo, initRef string) ([]state.Repo, error) {
	orgRepos, err := gh.ListOrgRepos(org)
	if err != nil {
		return nil, err
	}
	orgSet := map[string]bool{}
	merged := map[string]state.Repo{}
	for _, r := range orgRepos {
		orgSet[r] = true
		merged[r] = state.Repo{Repo: r, Kind: "org"}
	}

	makefile, _ := gh.GetFile(initRepo, "Makefile", initRef)
	gomod, _ := gh.GetFile(initRepo, "go.mod", initRef)
	for _, slug := range ParseDeps(makefile, gomod) {
		if _, ok := merged[slug]; !ok {
			kind := "external"
			if orgSet[slug] {
				kind = "org"
			} else if hasPrefix(slug, "kairos-io/") {
				kind = "dep"
			}
			merged[slug] = state.Repo{Repo: slug, Kind: kind}
		}
	}

	// Config additions / metadata overrides (matched by .Repo).
	for _, r := range cfg.Repos {
		existing := merged[r.Repo]
		existing.Repo = r.Repo
		if r.Kind != "" {
			existing.Kind = r.Kind
		}
		if r.Branch != "" {
			existing.Branch = r.Branch
		}
		if r.Criticality != "" {
			existing.Criticality = r.Criticality
		}
		if len(r.Artifacts) > 0 {
			existing.Artifacts = r.Artifacts
		}
		merged[r.Repo] = existing
	}

	for _, ex := range cfg.Exclude {
		delete(merged, ex)
	}

	out := make([]state.Repo, 0, len(merged))
	for _, r := range merged {
		if r.Branch == "" {
			r.Branch = "main"
		}
		if r.Criticality == "" {
			r.Criticality = "medium"
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Repo < out[j].Repo })
	return out, nil
}

func hasPrefix(s, p string) bool { return len(s) >= len(p) && s[:len(p)] == p }
```

- [ ] **Step 8: Run it; expect PASS.** Run: `go test ./internal/discover/...`

- [ ] **Step 9: Wire the `discover` subcommand**

In `cmd/ksec/main.go`, replace the placeholder loop with explicit subcommands. Add this command builder and call it from `newRootCmd` (replace the `for _, phase := range ...` block with `root.AddCommand(newDiscoverCmd(gf))` plus the still-stubbed others). New function:

```go
func newDiscoverCmd(gf *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "discover",
		Short: "build the tracked-repo list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadRepos("repos.yaml")
			if err != nil {
				return err
			}
			repos, err := discover.Run(ghclient.NewCLI(), cfg, "kairos-io", "kairos-io/kairos-init", "main")
			if err != nil {
				return err
			}
			return state.Save(gf.stateDir, state.ReposFile, repos)
		},
	}
}
```

Add imports `github.com/kairos-io/security/internal/{config,discover,ghclient,state}`. Keep the other four phases as stub commands via a small `newStubCmd(name)` helper so the build stays green.

- [ ] **Step 10: Build + commit**

Run: `go build ./...`
Expected: success.

```bash
git add internal/discover/ cmd/ksec/main.go
git commit -m "feat: discover phase (hybrid repo discovery)"
```

---

### Task 6: `collect` — Collector interface + sourceCVE (govulncheck)

**Files:**
- Create: `internal/collect/collect.go`
- Create: `internal/collect/source.go`
- Test: `internal/collect/source_test.go`

**Interfaces:**
- Consumes: `state.Repo`, `state.Finding`.
- Produces:
  - `type Collector interface { Name() string; Collect(repo state.Repo) ([]state.Finding, error) }`
  - `func FindingID(repo, typ, cve, pkg string) string` — `sha256` hex of `repo|typ|cve|pkg`, used as the stable dedupe `ID` by every collector.
  - `func Today() string` — returns the injectable current date (`YYYY-MM-DD`); tests override `nowFn`.
  - `type SourceCVE struct { Runner func(repo state.Repo) ([]byte, error) }` implementing `Collector` (`Name()=="sourceCVE"`), parsing govulncheck JSON-lines into findings.

- [ ] **Step 1: Write the failing test**

Create `internal/collect/source_test.go`:

```go
package collect

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Two govulncheck JSON-line messages: one "osv" (advisory metadata) and one
// "finding" referencing it at module scope.
const govulnJSON = `
{"osv":{"id":"GO-2025-1234","aliases":["CVE-2025-1234"],"summary":"x/net flaw","affected":[{"package":{"name":"golang.org/x/net"},"ranges":[{"type":"SEMVER","events":[{"introduced":"0"},{"fixed":"0.33.0"}]}]}]}}
{"finding":{"osv":"GO-2025-1234","trace":[{"module":"golang.org/x/net","version":"v0.30.0"}]}}
`

func TestSourceCVEParse(t *testing.T) {
	nowFn = func() string { return "2026-06-19" }
	defer func() { nowFn = defaultNow }

	c := SourceCVE{Runner: func(state.Repo) ([]byte, error) { return []byte(govulnJSON), nil }}
	fs, err := c.Collect(state.Repo{Repo: "kairos-io/immucore"})
	require.NoError(t, err)
	require.Len(t, fs, 1)
	f := fs[0]
	assert.Equal(t, "sourceCVE", f.Type)
	assert.Equal(t, "CVE-2025-1234", f.CVEID)
	assert.Equal(t, "golang.org/x/net", f.Package)
	assert.Equal(t, "v0.30.0", f.CurrentVersion)
	assert.Equal(t, "0.33.0", f.FixedVersion)
	assert.Equal(t, "go", f.Ecosystem)
	assert.Equal(t, FindingID("kairos-io/immucore", "sourceCVE", "CVE-2025-1234", "golang.org/x/net"), f.ID)
	assert.Equal(t, "2026-06-19", f.FirstSeen)
}
```

- [ ] **Step 2: Run it; expect FAIL.** Run: `go test ./internal/collect/...`

- [ ] **Step 3: Implement the shared collector helpers + sourceCVE**

Create `internal/collect/collect.go`:

```go
package collect

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func defaultNow() string { return time.Now().UTC().Format("2006-01-02") }

// nowFn is overridable in tests.
var nowFn = defaultNow

func Today() string { return nowFn() }

func FindingID(repo, typ, cve, pkg string) string {
	sum := sha256.Sum256([]byte(repo + "|" + typ + "|" + cve + "|" + pkg))
	return hex.EncodeToString(sum[:])
}
```

Create `internal/collect/source.go`:

```go
package collect

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/kairos-io/security/internal/state"
)

type SourceCVE struct {
	Runner func(repo state.Repo) ([]byte, error)
}

func (SourceCVE) Name() string { return "sourceCVE" }

type govulnLine struct {
	OSV *struct {
		ID       string   `json:"id"`
		Aliases  []string `json:"aliases"`
		Summary  string   `json:"summary"`
		Affected []struct {
			Package struct {
				Name string `json:"name"`
			} `json:"package"`
			Ranges []struct {
				Events []struct {
					Fixed string `json:"fixed"`
				} `json:"events"`
			} `json:"ranges"`
		} `json:"affected"`
	} `json:"osv"`
	Finding *struct {
		OSV   string `json:"osv"`
		Trace []struct {
			Module  string `json:"module"`
			Version string `json:"version"`
		} `json:"trace"`
	} `json:"finding"`
}

func (c SourceCVE) Collect(repo state.Repo) ([]state.Finding, error) {
	raw, err := c.Runner(repo)
	if err != nil {
		return nil, err
	}
	type adv struct {
		cve, fixed, summary string
	}
	advisories := map[string]adv{}
	var findings []govulnLine

	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Buffer(make([]byte, 0, 1024*1024), 8*1024*1024)
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		var gl govulnLine
		if err := json.Unmarshal(line, &gl); err != nil {
			continue // tolerate non-JSON progress lines
		}
		if gl.OSV != nil {
			a := adv{summary: gl.OSV.Summary}
			for _, al := range gl.OSV.Aliases {
				if len(al) > 3 && al[:3] == "CVE" {
					a.cve = al
				}
			}
			for _, af := range gl.OSV.Affected {
				for _, rg := range af.Ranges {
					for _, ev := range rg.Events {
						if ev.Fixed != "" {
							a.fixed = ev.Fixed
						}
					}
				}
			}
			advisories[gl.OSV.ID] = a
		}
		if gl.Finding != nil {
			findings = append(findings, gl)
		}
	}

	out := map[string]state.Finding{}
	for _, gl := range findings {
		if len(gl.Finding.Trace) == 0 {
			continue
		}
		t := gl.Finding.Trace[0]
		a := advisories[gl.Finding.OSV]
		cve := a.cve
		if cve == "" {
			cve = gl.Finding.OSV
		}
		f := state.Finding{
			ID:             FindingID(repo.Repo, "sourceCVE", cve, t.Module),
			Repo:           repo.Repo,
			Type:           "sourceCVE",
			CVEID:          cve,
			Ecosystem:      "go",
			Package:        t.Module,
			CurrentVersion: t.Version,
			FixedVersion:   a.fixed,
			Severity:       "unknown",
			Source:         "govulncheck",
			Title:          a.summary,
			FirstSeen:      Today(),
			LastSeen:       Today(),
		}
		out[f.ID] = f // dedupe module-level findings
	}
	res := make([]state.Finding, 0, len(out))
	for _, f := range out {
		res = append(res, f)
	}
	return res, nil
}
```

- [ ] **Step 4: Run it; expect PASS.** Run: `go test ./internal/collect/...`

- [ ] **Step 5: Commit**

```bash
git add internal/collect/collect.go internal/collect/source.go internal/collect/source_test.go
git commit -m "feat: collect interface and sourceCVE (govulncheck) collector"
```

---

### Task 7: `collect` — imageCVE collector (trivy JSON)

**Files:**
- Create: `internal/collect/image.go`
- Test: `internal/collect/image_test.go`

**Interfaces:**
- Consumes: helpers from Task 6.
- Produces: `type ImageCVE struct { Runner func(ref string) ([]byte, error) }` implementing `Collector` (`Name()=="imageCVE"`); runs once per `Type=="image"` artifact, parses trivy `--format json` output.

- [ ] **Step 1: Write the failing test**

Create `internal/collect/image_test.go`:

```go
package collect

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const trivyJSON = `{"Results":[{"Target":"x","Vulnerabilities":[
{"VulnerabilityID":"CVE-2025-9999","PkgName":"openssl","InstalledVersion":"1.1.1","FixedVersion":"1.1.1w","Severity":"CRITICAL","PrimaryURL":"https://x/CVE-2025-9999","Title":"openssl flaw"}]}]}`

func TestImageCVEParse(t *testing.T) {
	nowFn = func() string { return "2026-06-19" }
	defer func() { nowFn = defaultNow }

	c := ImageCVE{Runner: func(string) ([]byte, error) { return []byte(trivyJSON), nil }}
	repo := state.Repo{Repo: "kairos-io/kairos", Artifacts: []state.Artifact{{Type: "image", Ref: "quay.io/kairos/x:latest"}}}
	fs, err := c.Collect(repo)
	require.NoError(t, err)
	require.Len(t, fs, 1)
	f := fs[0]
	assert.Equal(t, "imageCVE", f.Type)
	assert.Equal(t, "CVE-2025-9999", f.CVEID)
	assert.Equal(t, "openssl", f.Package)
	assert.Equal(t, "critical", f.Severity)
	assert.Equal(t, "1.1.1w", f.FixedVersion)
}
```

- [ ] **Step 2: Run it; expect FAIL.** Run: `go test ./internal/collect/...`

- [ ] **Step 3: Implement**

Create `internal/collect/image.go`:

```go
package collect

import (
	"encoding/json"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

type ImageCVE struct {
	Runner func(ref string) ([]byte, error)
}

func (ImageCVE) Name() string { return "imageCVE" }

type trivyReport struct {
	Results []struct {
		Vulnerabilities []struct {
			VulnerabilityID  string `json:"VulnerabilityID"`
			PkgName          string `json:"PkgName"`
			InstalledVersion string `json:"InstalledVersion"`
			FixedVersion     string `json:"FixedVersion"`
			Severity         string `json:"Severity"`
			PrimaryURL       string `json:"PrimaryURL"`
			Title            string `json:"Title"`
		} `json:"Vulnerabilities"`
	} `json:"Results"`
}

func (c ImageCVE) Collect(repo state.Repo) ([]state.Finding, error) {
	out := map[string]state.Finding{}
	for _, art := range repo.Artifacts {
		if art.Type != "image" {
			continue
		}
		raw, err := c.Runner(art.Ref)
		if err != nil {
			return nil, err
		}
		var rep trivyReport
		if err := json.Unmarshal(raw, &rep); err != nil {
			return nil, err
		}
		for _, res := range rep.Results {
			for _, v := range res.Vulnerabilities {
				f := state.Finding{
					ID:             FindingID(repo.Repo, "imageCVE", v.VulnerabilityID, v.PkgName),
					Repo:           repo.Repo,
					Type:           "imageCVE",
					CVEID:          v.VulnerabilityID,
					Package:        v.PkgName,
					CurrentVersion: v.InstalledVersion,
					FixedVersion:   v.FixedVersion,
					Severity:       strings.ToLower(v.Severity),
					Source:         "trivy",
					Title:          v.Title,
					URL:            v.PrimaryURL,
					FirstSeen:      Today(),
					LastSeen:       Today(),
				}
				out[f.ID] = f
			}
		}
	}
	res := make([]state.Finding, 0, len(out))
	for _, f := range out {
		res = append(res, f)
	}
	return res, nil
}
```

- [ ] **Step 4: Run it; expect PASS.** Run: `go test ./internal/collect/...`

- [ ] **Step 5: Commit**

```bash
git add internal/collect/image.go internal/collect/image_test.go
git commit -m "feat: imageCVE (trivy) collector"
```

---

### Task 8: `collect` — open security PRs collector

**Files:**
- Create: `internal/collect/prs.go`
- Test: `internal/collect/prs_test.go`

**Interfaces:**
- Consumes: `ghclient.GitHub`, helpers from Task 6.
- Produces: `type PRs struct { GH ghclient.GitHub }` implementing `Collector` (`Name()=="prs"`); emits a `Finding{Type:"pr"}` for each open PR authored by `renovate`/`dependabot`/`kairos-security-bot` OR labeled `security`/`dependencies`.

- [ ] **Step 1: Write the failing test**

Create `internal/collect/prs_test.go`:

```go
package collect

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPRsCollectorFiltersSecurityPRs(t *testing.T) {
	nowFn = func() string { return "2026-06-19" }
	defer func() { nowFn = defaultNow }

	gh := ghclient.NewFake()
	gh.PRs["kairos-io/immucore"] = []ghclient.PullRequest{
		{Number: 1, Title: "Bump x/net", Author: "renovate[bot]", URL: "u1"},
		{Number: 2, Title: "Feature", Author: "alice", Labels: []string{"enhancement"}},
		{Number: 3, Title: "Patch CVE", Author: "alice", Labels: []string{"security"}},
	}
	c := PRs{GH: gh}
	fs, err := c.Collect(state.Repo{Repo: "kairos-io/immucore"})
	require.NoError(t, err)
	require.Len(t, fs, 2) // PR 1 (renovate) + PR 3 (security label)
	assert.Equal(t, "pr", fs[0].Type)
}
```

- [ ] **Step 2: Run it; expect FAIL.** Run: `go test ./internal/collect/...`

- [ ] **Step 3: Implement**

Create `internal/collect/prs.go`:

```go
package collect

import (
	"fmt"
	"sort"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

type PRs struct {
	GH ghclient.GitHub
}

func (PRs) Name() string { return "prs" }

var botAuthors = map[string]bool{
	"renovate[bot]": true, "dependabot[bot]": true, "kairos-security-bot": true,
}
var secLabels = map[string]bool{"security": true, "dependencies": true}

func (c PRs) Collect(repo state.Repo) ([]state.Finding, error) {
	prs, err := c.GH.ListOpenPRs(repo.Repo)
	if err != nil {
		return nil, err
	}
	var out []state.Finding
	for _, pr := range prs {
		if !isSecurityPR(pr) {
			continue
		}
		out = append(out, state.Finding{
			ID:        FindingID(repo.Repo, "pr", fmt.Sprintf("#%d", pr.Number), ""),
			Repo:      repo.Repo,
			Type:      "pr",
			Severity:  "unknown",
			Source:    "github-pr",
			Title:     fmt.Sprintf("#%d %s (@%s)", pr.Number, pr.Title, pr.Author),
			URL:       pr.URL,
			FirstSeen: Today(),
			LastSeen:  Today(),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func isSecurityPR(pr ghclient.PullRequest) bool {
	if botAuthors[pr.Author] {
		return true
	}
	for _, l := range pr.Labels {
		if secLabels[l] {
			return true
		}
	}
	return false
}
```

- [ ] **Step 4: Run it; expect PASS.** Run: `go test ./internal/collect/...`

- [ ] **Step 5: Commit**

```bash
git add internal/collect/prs.go internal/collect/prs_test.go
git commit -m "feat: open security PRs collector"
```

---

### Task 9: `collect` — GitHub security alerts collector

**Files:**
- Create: `internal/collect/alerts.go`
- Test: `internal/collect/alerts_test.go`

**Interfaces:**
- Consumes: `ghclient.GitHub`, helpers from Task 6.
- Produces: `type GHAlerts struct { GH ghclient.GitHub }` implementing `Collector` (`Name()=="ghAlerts"`); maps each Dependabot alert to a `Finding{Type:"ghAlert"}`.

- [ ] **Step 1: Write the failing test**

Create `internal/collect/alerts_test.go`:

```go
package collect

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGHAlertsCollector(t *testing.T) {
	nowFn = func() string { return "2026-06-19" }
	defer func() { nowFn = defaultNow }

	gh := ghclient.NewFake()
	gh.Alerts["kairos-io/immucore"] = []ghclient.Alert{
		{Number: 7, CVEID: "CVE-2025-1234", GHSA: "GHSA-aaa", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high", URL: "u", FixedVersion: "0.33.0"},
	}
	c := GHAlerts{GH: gh}
	fs, err := c.Collect(state.Repo{Repo: "kairos-io/immucore"})
	require.NoError(t, err)
	require.Len(t, fs, 1)
	assert.Equal(t, "ghAlert", fs[0].Type)
	assert.Equal(t, "CVE-2025-1234", fs[0].CVEID)
	assert.Equal(t, "high", fs[0].Severity)
	assert.Equal(t, "0.33.0", fs[0].FixedVersion)
}
```

- [ ] **Step 2: Run it; expect FAIL.** Run: `go test ./internal/collect/...`

- [ ] **Step 3: Implement**

Create `internal/collect/alerts.go`:

```go
package collect

import (
	"strings"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

type GHAlerts struct {
	GH ghclient.GitHub
}

func (GHAlerts) Name() string { return "ghAlerts" }

func (c GHAlerts) Collect(repo state.Repo) ([]state.Finding, error) {
	alerts, err := c.GH.ListDependabotAlerts(repo.Repo)
	if err != nil {
		return nil, err
	}
	out := make([]state.Finding, 0, len(alerts))
	for _, a := range alerts {
		cve := a.CVEID
		if cve == "" {
			cve = a.GHSA
		}
		out = append(out, state.Finding{
			ID:           FindingID(repo.Repo, "ghAlert", cve, a.Package),
			Repo:         repo.Repo,
			Type:         "ghAlert",
			CVEID:        a.CVEID,
			GHSA:         a.GHSA,
			Ecosystem:    strings.ToLower(a.Ecosystem),
			Package:      a.Package,
			FixedVersion: a.FixedVersion,
			Severity:     strings.ToLower(a.Severity),
			Source:       "dependabot",
			URL:          a.URL,
			FirstSeen:    Today(),
			LastSeen:     Today(),
		})
	}
	return out, nil
}
```

- [ ] **Step 4: Run it; expect PASS.** Run: `go test ./internal/collect/...`

- [ ] **Step 5: Commit**

```bash
git add internal/collect/alerts.go internal/collect/alerts_test.go
git commit -m "feat: GitHub security alerts collector"
```

---

### Task 10: `collect` — fan-out orchestration with error isolation + aging

**Files:**
- Modify: `internal/collect/collect.go`
- Test: `internal/collect/collect_test.go`
- Modify: `cmd/ksec/main.go` (wire `collect`)

**Interfaces:**
- Consumes: `Collector`, `[]state.Repo`, the previous `state.Findings` (for aging).
- Produces: `func Run(repos []state.Repo, collectors []Collector, prev state.Findings) state.Findings` — runs every collector against every repo, isolates errors into `Findings.Errors`, preserves `FirstSeen` from `prev` for findings whose `ID` already existed (aging), sorts findings by `ID`.

- [ ] **Step 1: Write the failing test**

Create `internal/collect/collect_test.go`:

```go
package collect

import (
	"errors"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubCollector struct {
	name string
	out  []state.Finding
	err  error
}

func (s stubCollector) Name() string                                   { return s.name }
func (s stubCollector) Collect(state.Repo) ([]state.Finding, error)    { return s.out, s.err }

func TestRunIsolatesErrorsAndPreservesFirstSeen(t *testing.T) {
	nowFn = func() string { return "2026-06-19" }
	defer func() { nowFn = defaultNow }

	repos := []state.Repo{{Repo: "kairos-io/immucore"}}
	good := stubCollector{name: "good", out: []state.Finding{
		{ID: "x", Repo: "kairos-io/immucore", Type: "sourceCVE", FirstSeen: "2026-06-19", LastSeen: "2026-06-19"},
	}}
	bad := stubCollector{name: "bad", err: errors.New("rate limited")}

	prev := state.Findings{Findings: []state.Finding{{ID: "x", FirstSeen: "2026-06-01"}}}
	got := Run(repos, []Collector{good, bad}, prev)

	require.Len(t, got.Findings, 1)
	assert.Equal(t, "2026-06-01", got.Findings[0].FirstSeen, "aging preserved")
	assert.Equal(t, "2026-06-19", got.Findings[0].LastSeen)
	require.Len(t, got.Errors, 1)
	assert.Equal(t, "bad", got.Errors[0].Collector)
	assert.Contains(t, got.Errors[0].Message, "rate limited")
}
```

- [ ] **Step 2: Run it; expect FAIL** — `Run` undefined.

- [ ] **Step 3: Implement `Run` (append to `collect.go`)**

Add to `internal/collect/collect.go`:

```go
import (
	"sort"

	"github.com/kairos-io/security/internal/state"
)

func Run(repos []state.Repo, collectors []Collector, prev state.Findings) state.Findings {
	firstSeen := map[string]string{}
	for _, f := range prev.Findings {
		firstSeen[f.ID] = f.FirstSeen
	}

	var res state.Findings
	for _, repo := range repos {
		for _, col := range collectors {
			found, err := col.Collect(repo)
			if err != nil {
				res.Errors = append(res.Errors, state.CollectionError{
					Repo: repo.Repo, Collector: col.Name(), Message: err.Error(),
				})
				continue
			}
			for _, f := range found {
				if fs, ok := firstSeen[f.ID]; ok && fs != "" {
					f.FirstSeen = fs
				}
				res.Findings = append(res.Findings, f)
			}
		}
	}
	sort.Slice(res.Findings, func(i, j int) bool { return res.Findings[i].ID < res.Findings[j].ID })
	sort.Slice(res.Errors, func(i, j int) bool {
		if res.Errors[i].Repo != res.Errors[j].Repo {
			return res.Errors[i].Repo < res.Errors[j].Repo
		}
		return res.Errors[i].Collector < res.Errors[j].Collector
	})
	return res
}
```

> Note: the existing `collect.go` already declares `package collect` and imports `crypto/sha256`, `encoding/hex`, `time`. Merge the new imports into the existing import block rather than adding a second one.

- [ ] **Step 4: Run it; expect PASS.** Run: `go test ./internal/collect/...`

- [ ] **Step 5: Wire the `collect` subcommand**

Add to `cmd/ksec/main.go` (and register via `root.AddCommand(newCollectCmd(gf))`). The real scanners are wrapped here so the package stays testable:

```go
func newCollectCmd(gf *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "collect",
		Short: "gather raw findings per repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			var repos []state.Repo
			if err := state.Load(gf.stateDir, state.ReposFile, &repos); err != nil {
				return err
			}
			var prev state.Findings
			_ = state.Load(gf.stateDir, state.FindingsFile, &prev) // best-effort for aging

			gh := ghclient.NewCLI()
			collectors := []collect.Collector{
				collect.PRs{GH: gh},
				collect.GHAlerts{GH: gh},
				collect.ImageCVE{Runner: trivyRunner},
				collect.SourceCVE{Runner: govulncheckRunner},
			}
			out := collect.Run(repos, collectors, prev)
			return state.Save(gf.stateDir, state.FindingsFile, out)
		},
	}
}

func trivyRunner(ref string) ([]byte, error) {
	return exec.Command("trivy", "image", "--quiet", "--scanners", "vuln", "--format", "json", ref).Output()
}

// govulncheckRunner shallow-clones the repo to a temp dir and runs govulncheck.
func govulncheckRunner(r state.Repo) ([]byte, error) {
	dir, err := os.MkdirTemp("", "ksec-src-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	clone := exec.Command("git", "clone", "--depth", "1", "--branch", r.Branch,
		"https://github.com/"+r.Repo+".git", dir)
	if out, err := clone.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("clone: %v: %s", err, out)
	}
	cmd := exec.Command("govulncheck", "-json", "./...")
	cmd.Dir = dir
	return cmd.Output() // non-zero exit with findings still yields JSON on stdout
}
```

Add imports `os`, `os/exec`, `fmt`, and `github.com/kairos-io/security/internal/collect`.

- [ ] **Step 6: Build + commit**

Run: `go build ./...`

```bash
git add internal/collect/collect.go internal/collect/collect_test.go cmd/ksec/main.go
git commit -m "feat: collect fan-out with error isolation and CVE aging"
```

---

### Task 11: `correlate` phase — dedupe + waterfall grouping

**Files:**
- Create: `internal/correlate/correlate.go`
- Test: `internal/correlate/correlate_test.go`
- Modify: `cmd/ksec/main.go` (wire `correlate`)

**Interfaces:**
- Consumes: `state.Findings`, `state.Correlated`, `state.WaterfallGroup`, `state.Bump`.
- Produces: `func Run(in state.Findings) state.Correlated` — (1) merges findings that share the same `(Repo, CVEID, Package)` into one (keeping the highest severity and any non-empty `FixedVersion`); (2) builds waterfall groups for Go-ecosystem CVEs sharing the same `(CVEID, Package)` across ≥2 repos, with `SuggestedBump` = package@fixedVersion. Deterministic and pure.

- [ ] **Step 1: Write the failing test**

Create `internal/correlate/correlate_test.go`:

```go
package correlate

import (
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCorrelateDedupesAndBuildsWaterfall(t *testing.T) {
	in := state.Findings{Findings: []state.Finding{
		// same CVE in immucore seen by two sources → dedupe to 1, severity high wins
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "unknown", FixedVersion: "0.33.0"},
		{ID: "b", Repo: "kairos-io/immucore", Type: "ghAlert", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high"},
		// same CVE/package in a second repo → waterfall group of 2 repos
		{ID: "c", Repo: "kairos-io/kairos-agent", Type: "sourceCVE", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high", FixedVersion: "0.33.0"},
	}}

	out := Run(in)

	// dedupe: immucore CVE-2025-1 collapses to one finding, severity "high"
	count := 0
	for _, f := range out.Findings {
		if f.Repo == "kairos-io/immucore" && f.CVEID == "CVE-2025-1" {
			count++
			assert.Equal(t, "high", f.Severity)
			assert.Equal(t, "0.33.0", f.FixedVersion)
		}
	}
	assert.Equal(t, 1, count)

	require.Len(t, out.Waterfall, 1)
	g := out.Waterfall[0]
	assert.ElementsMatch(t, []string{"kairos-io/immucore", "kairos-io/kairos-agent"}, g.AffectedRepos)
	assert.Equal(t, "golang.org/x/net", g.SuggestedBump.Package)
	assert.Equal(t, "0.33.0", g.SuggestedBump.To)
	assert.Equal(t, "high", g.Severity)
}
```

- [ ] **Step 2: Run it; expect FAIL.** Run: `go test ./internal/correlate/...`

- [ ] **Step 3: Implement**

Create `internal/correlate/correlate.go`:

```go
package correlate

import (
	"fmt"
	"sort"

	"github.com/kairos-io/security/internal/state"
)

var sevRank = map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1, "unknown": 0, "": 0}

func worse(a, b string) string {
	if sevRank[a] >= sevRank[b] {
		return a
	}
	return b
}

func Run(in state.Findings) state.Correlated {
	// 1) dedupe by (repo, cveID, package); PR findings (no CVE) pass through.
	merged := map[string]state.Finding{}
	var order []string
	for _, f := range in.Findings {
		key := f.Repo + "|" + f.CVEID + "|" + f.Package
		if f.CVEID == "" {
			key = f.ID // PRs and CVE-less findings never merge
		}
		cur, ok := merged[key]
		if !ok {
			merged[key] = f
			order = append(order, key)
			continue
		}
		cur.Severity = worse(cur.Severity, f.Severity)
		if cur.FixedVersion == "" {
			cur.FixedVersion = f.FixedVersion
		}
		if cur.FirstSeen == "" || (f.FirstSeen != "" && f.FirstSeen < cur.FirstSeen) {
			cur.FirstSeen = f.FirstSeen
		}
		merged[key] = cur
	}

	findings := make([]state.Finding, 0, len(merged))
	for _, k := range order {
		findings = append(findings, merged[k])
	}
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })

	// 2) waterfall: group go-ecosystem CVEs by (cveID, package) across repos.
	type agg struct {
		repos    map[string]bool
		severity string
		fixed    string
	}
	groups := map[string]*agg{}
	for _, f := range findings {
		if f.Ecosystem != "go" || f.CVEID == "" || f.Package == "" {
			continue
		}
		gk := f.CVEID + "|" + f.Package
		g := groups[gk]
		if g == nil {
			g = &agg{repos: map[string]bool{}}
			groups[gk] = g
		}
		g.repos[f.Repo] = true
		g.severity = worse(g.severity, f.Severity)
		if g.fixed == "" {
			g.fixed = f.FixedVersion
		}
	}

	var waterfall []state.WaterfallGroup
	for gk, g := range groups {
		if len(g.repos) < 2 {
			continue
		}
		repos := make([]string, 0, len(g.repos))
		for r := range g.repos {
			repos = append(repos, r)
		}
		sort.Strings(repos)
		cve, pkg := splitKey(gk)
		waterfall = append(waterfall, state.WaterfallGroup{
			ID:            "go-" + cve + "-" + pkg,
			RootCause:     fmt.Sprintf("%s (%s)", pkg, cve),
			Ecosystem:     "go",
			Severity:      g.severity,
			AffectedRepos: repos,
			SuggestedBump: state.Bump{Package: pkg, To: g.fixed},
		})
	}
	sort.Slice(waterfall, func(i, j int) bool { return waterfall[i].ID < waterfall[j].ID })

	return state.Correlated{Findings: findings, Waterfall: waterfall}
}

func splitKey(k string) (cve, pkg string) {
	for i := 0; i < len(k); i++ {
		if k[i] == '|' {
			return k[:i], k[i+1:]
		}
	}
	return k, ""
}
```

- [ ] **Step 4: Run it; expect PASS.** Run: `go test ./internal/correlate/...`

- [ ] **Step 5: Wire `correlate` subcommand**

Add to `cmd/ksec/main.go`, register with `root.AddCommand(newCorrelateCmd(gf))`:

```go
func newCorrelateCmd(gf *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "correlate",
		Short: "dedupe findings and build the waterfall graph",
		RunE: func(cmd *cobra.Command, args []string) error {
			var in state.Findings
			if err := state.Load(gf.stateDir, state.FindingsFile, &in); err != nil {
				return err
			}
			return state.Save(gf.stateDir, state.CorrelatedFile, correlate.Run(in))
		},
	}
}
```

Add import `github.com/kairos-io/security/internal/correlate`.

- [ ] **Step 6: Build + commit**

```bash
go build ./...
git add internal/correlate/ cmd/ksec/main.go
git commit -m "feat: correlate phase (dedupe + waterfall grouping)"
```

---

### Task 12: `triage` phase — AIClient + deterministic fallback + nib impl

**Files:**
- Create: `internal/triage/triage.go`
- Create: `internal/triage/nib.go`
- Test: `internal/triage/triage_test.go`
- Modify: `cmd/ksec/main.go` (wire `triage`)

**Interfaces:**
- Consumes: `state.Correlated`, `state.Triage`, `config.AIConfig`.
- Produces:
  - `type AIClient interface { Summarize(c state.Correlated) (focus []string, summaries map[string]string, narrative string, err error) }`
  - `func Run(c state.Correlated, ai AIClient, model string) state.Triage` — calls `ai.Summarize`; on error, uses `deterministicFocus`/`templatedSummaries` and sets `AIAvailable=false`.
  - `func deterministicFocus(c state.Correlated) []string` — finding IDs + waterfall IDs ordered by severity desc, then by ID.
  - `type NibClient struct { cfg config.AIConfig; run func(prompt string) ([]byte, error) }` implementing `AIClient` (shells `nib --yolo`).

- [ ] **Step 1: Write the failing test**

Create `internal/triage/triage_test.go`:

```go
package triage

import (
	"errors"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubAI struct {
	focus     []string
	summaries map[string]string
	narrative string
	err       error
}

func (s stubAI) Summarize(state.Correlated) ([]string, map[string]string, string, error) {
	return s.focus, s.summaries, s.narrative, s.err
}

var sampleCorrelated = state.Correlated{
	Findings: []state.Finding{
		{ID: "low1", Severity: "low"},
		{ID: "crit1", Severity: "critical"},
		{ID: "high1", Severity: "high"},
	},
}

func TestRunUsesAIWhenAvailable(t *testing.T) {
	ai := stubAI{focus: []string{"crit1"}, summaries: map[string]string{"crit1": "bad"}, narrative: "n"}
	got := Run(sampleCorrelated, ai, "m")
	assert.True(t, got.AIAvailable)
	assert.Equal(t, []string{"crit1"}, got.Focus)
	assert.Equal(t, "n", got.Narrative)
}

func TestRunFallsBackOnAIError(t *testing.T) {
	got := Run(sampleCorrelated, stubAI{err: errors.New("model down")}, "m")
	assert.False(t, got.AIAvailable)
	// deterministic severity ordering: critical, high, low
	require.Equal(t, []string{"crit1", "high1", "low1"}, got.Focus)
	assert.NotEmpty(t, got.Summaries["crit1"])
}
```

- [ ] **Step 2: Run it; expect FAIL.** Run: `go test ./internal/triage/...`

- [ ] **Step 3: Implement triage core**

Create `internal/triage/triage.go`:

```go
package triage

import (
	"fmt"
	"sort"
	"time"

	"github.com/kairos-io/security/internal/state"
)

type AIClient interface {
	Summarize(c state.Correlated) (focus []string, summaries map[string]string, narrative string, err error)
}

var sevRank = map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1, "unknown": 0, "": 0}

var nowFn = func() string { return time.Now().UTC().Format("2006-01-02") }

func Run(c state.Correlated, ai AIClient, model string) state.Triage {
	t := state.Triage{GeneratedAt: nowFn(), Model: model, Summaries: map[string]string{}}
	if ai != nil {
		focus, summaries, narrative, err := ai.Summarize(c)
		if err == nil {
			t.AIAvailable = true
			t.Focus = focus
			t.Summaries = summaries
			t.Narrative = narrative
			return t
		}
	}
	t.AIAvailable = false
	t.Focus = deterministicFocus(c)
	t.Summaries = templatedSummaries(c)
	t.Narrative = fmt.Sprintf("AI unavailable this run. %d findings, %d waterfall groups, ordered by severity.",
		len(c.Findings), len(c.Waterfall))
	return t
}

func deterministicFocus(c state.Correlated) []string {
	type item struct {
		id  string
		sev string
	}
	var items []item
	for _, f := range c.Findings {
		items = append(items, item{f.ID, f.Severity})
	}
	for _, g := range c.Waterfall {
		items = append(items, item{g.ID, g.Severity})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if sevRank[items[i].sev] != sevRank[items[j].sev] {
			return sevRank[items[i].sev] > sevRank[items[j].sev]
		}
		return items[i].id < items[j].id
	})
	out := make([]string, 0, len(items))
	for _, it := range items {
		out = append(out, it.id)
	}
	return out
}

func templatedSummaries(c state.Correlated) map[string]string {
	out := map[string]string{}
	for _, f := range c.Findings {
		if sevRank[f.Severity] >= sevRank["high"] {
			out[f.ID] = fmt.Sprintf("%s %s in %s (%s)", f.Severity, f.CVEID, f.Repo, f.Package)
		}
	}
	for _, g := range c.Waterfall {
		out[g.ID] = fmt.Sprintf("%s affects %d repos via %s", g.Severity, len(g.AffectedRepos), g.RootCause)
	}
	return out
}
```

- [ ] **Step 4: Run it; expect PASS.** Run: `go test ./internal/triage/...`

- [ ] **Step 5: Implement the nib client (no new test; exercised by the E2E dry-run in Task 15)**

Create `internal/triage/nib.go`:

```go
package triage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/state"
)

type NibClient struct {
	cfg config.AIConfig
	run func(prompt string) ([]byte, error)
}

func NewNibClient(cfg config.AIConfig) *NibClient {
	return &NibClient{cfg: cfg, run: func(prompt string) ([]byte, error) {
		cmd := exec.Command("nib", "--"+cfg.Nib.Mode,
			"--model", cfg.Nib.Model, "--endpoint", cfg.Nib.Endpoint)
		cmd.Stdin = bytes.NewBufferString(prompt)
		var out, errb bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &errb
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("nib: %v: %s", err, errb.String())
		}
		return out.Bytes(), nil
	}}
}

// aiResponse is the JSON contract we instruct the model to emit.
type aiResponse struct {
	Focus     []string          `json:"focus"`
	Summaries map[string]string `json:"summaries"`
	Narrative string            `json:"narrative"`
}

func (n *NibClient) Summarize(c state.Correlated) ([]string, map[string]string, string, error) {
	payload, err := json.Marshal(c)
	if err != nil {
		return nil, nil, "", err
	}
	prompt := "You are a security triage assistant. Given this JSON of correlated security findings, " +
		"return ONLY a JSON object with keys: focus (array of finding/waterfall IDs ordered most-urgent first), " +
		"summaries (map of id to one-line human summary for high/critical items), and narrative (2-3 sentence " +
		"'what to focus on' overview). Do not invent IDs. Findings:\n" + string(payload)

	raw, err := n.run(prompt)
	if err != nil {
		return nil, nil, "", err
	}
	var resp aiResponse
	if err := json.Unmarshal(bytes.TrimSpace(raw), &resp); err != nil {
		return nil, nil, "", fmt.Errorf("parse model output: %w", err)
	}
	return resp.Focus, resp.Summaries, resp.Narrative, nil
}
```

- [ ] **Step 6: Wire `triage` subcommand**

Add to `cmd/ksec/main.go`, register with `root.AddCommand(newTriageCmd(gf))`:

```go
func newTriageCmd(gf *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "triage",
		Short: "prioritize findings and write the AI summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			var c state.Correlated
			if err := state.Load(gf.stateDir, state.CorrelatedFile, &c); err != nil {
				return err
			}
			aiCfg, err := config.LoadAI("ai.yaml")
			if err != nil {
				return err
			}
			out := triage.Run(c, triage.NewNibClient(aiCfg), aiCfg.LocalAI.Model.Name)
			return state.Save(gf.stateDir, state.TriageFile, out)
		},
	}
}
```

Add import `github.com/kairos-io/security/internal/triage`.

- [ ] **Step 7: Build + commit**

```bash
go build ./...
git add internal/triage/ cmd/ksec/main.go
git commit -m "feat: triage phase with nib AI client and deterministic fallback"
```

---

### Task 13: `render` — dashboard.json + dashboard.md (golden)

**Files:**
- Create: `internal/render/render.go`
- Test: `internal/render/render_test.go`
- Create: `internal/render/testdata/dashboard.md.golden`

**Interfaces:**
- Consumes: `state.Correlated`, `state.Triage`, `state.Findings` (for errors/aging).
- Produces:
  - `type Input struct { Correlated state.Correlated; Triage state.Triage; CollectErrors []state.CollectionError; RunURL string }`
  - `func DashboardJSON(in Input) ([]byte, error)` — stable indented JSON of the snapshot.
  - `func DashboardMarkdown(in Input) string` — renders Focus, Waterfall fronts, per-repo table, and a collection-errors note.

- [ ] **Step 1: Write the failing test**

Create `internal/render/render_test.go`:

```go
package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleInput() Input {
	return Input{
		Correlated: state.Correlated{
			Findings: []state.Finding{
				{ID: "crit1", Repo: "kairos-io/kairos", Type: "imageCVE", CVEID: "CVE-2025-9999", Package: "openssl", Severity: "critical", FirstSeen: "2026-06-01", LastSeen: "2026-06-19"},
			},
			Waterfall: []state.WaterfallGroup{
				{ID: "go-CVE-2025-1-golang.org/x/net", RootCause: "golang.org/x/net (CVE-2025-1)", Severity: "high", AffectedRepos: []string{"kairos-io/immucore", "kairos-io/kairos-agent"}, SuggestedBump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
			},
		},
		Triage: state.Triage{
			GeneratedAt: "2026-06-19", AIAvailable: true,
			Focus:     []string{"crit1", "go-CVE-2025-1-golang.org/x/net"},
			Summaries: map[string]string{"crit1": "Critical openssl CVE in kairos image"},
			Narrative: "Focus on the openssl critical first.",
		},
		CollectErrors: []state.CollectionError{{Repo: "kairos-io/x", Collector: "prs", Message: "rate limited"}},
		RunURL:        "https://github.com/kairos-io/security/actions/runs/1",
	}
}

func TestDashboardMarkdownGolden(t *testing.T) {
	got := DashboardMarkdown(sampleInput())
	golden := filepath.Join("testdata", "dashboard.md.golden")
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		require.NoError(t, os.WriteFile(golden, []byte(got), 0o644))
	}
	want, err := os.ReadFile(golden)
	require.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestDashboardJSONIsStable(t *testing.T) {
	a, err := DashboardJSON(sampleInput())
	require.NoError(t, err)
	b, err := DashboardJSON(sampleInput())
	require.NoError(t, err)
	assert.Equal(t, string(a), string(b))
}
```

- [ ] **Step 2: Run it; expect FAIL.** Run: `go test ./internal/render/...`

- [ ] **Step 3: Implement render**

Create `internal/render/render.go`:

```go
package render

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

type Input struct {
	Correlated    state.Correlated      `json:"correlated"`
	Triage        state.Triage          `json:"triage"`
	CollectErrors []state.CollectionError `json:"collectErrors"`
	RunURL        string                `json:"runURL"`
}

func DashboardJSON(in Input) ([]byte, error) {
	b, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

func DashboardMarkdown(in Input) string {
	var b strings.Builder
	b.WriteString("# Kairos Security Dashboard\n\n")
	fmt.Fprintf(&b, "_Updated %s", in.Triage.GeneratedAt)
	if !in.Triage.AIAvailable {
		b.WriteString(" — ⚠️ AI unavailable this run")
	}
	b.WriteString("._\n\n")

	if in.Triage.Narrative != "" {
		b.WriteString("> " + in.Triage.Narrative + "\n\n")
	}

	// Focus now
	b.WriteString("## 🔥 Focus now\n\n")
	if len(in.Triage.Focus) == 0 {
		b.WriteString("_Nothing flagged._\n\n")
	} else {
		for _, id := range in.Triage.Focus {
			if s, ok := in.Triage.Summaries[id]; ok {
				fmt.Fprintf(&b, "- **%s** — %s\n", id, s)
			} else {
				fmt.Fprintf(&b, "- **%s**\n", id)
			}
		}
		b.WriteString("\n")
	}

	// Waterfall fronts
	b.WriteString("## 🌊 Waterfall fronts\n\n")
	if len(in.Correlated.Waterfall) == 0 {
		b.WriteString("_None._\n\n")
	} else {
		b.WriteString("| Root cause | Severity | Bump | Affected repos |\n|---|---|---|---|\n")
		for _, g := range in.Correlated.Waterfall {
			fmt.Fprintf(&b, "| %s | %s | %s@%s | %s |\n",
				g.RootCause, g.Severity, g.SuggestedBump.Package, g.SuggestedBump.To,
				strings.Join(g.AffectedRepos, ", "))
		}
		b.WriteString("\n")
	}

	// Per-repo table
	b.WriteString("## 📦 Per-repo findings\n\n")
	b.WriteString("| Repo | Critical | High | Medium | Low | Total |\n|---|---|---|---|---|---|\n")
	for _, row := range perRepoRows(in.Correlated.Findings) {
		fmt.Fprintf(&b, "| %s | %d | %d | %d | %d | %d |\n",
			row.repo, row.crit, row.high, row.med, row.low, row.total)
	}
	b.WriteString("\n")

	// Collection errors
	if len(in.CollectErrors) > 0 {
		fmt.Fprintf(&b, "## ⚠️ %d collection errors\n\n", len(in.CollectErrors))
		for _, e := range in.CollectErrors {
			fmt.Fprintf(&b, "- `%s` / %s: %s\n", e.Repo, e.Collector, e.Message)
		}
		b.WriteString("\n")
	}

	if in.RunURL != "" {
		fmt.Fprintf(&b, "---\n[Run log](%s)\n", in.RunURL)
	}
	return b.String()
}

type repoRow struct {
	repo                     string
	crit, high, med, low, total int
}

func perRepoRows(findings []state.Finding) []repoRow {
	idx := map[string]*repoRow{}
	for _, f := range findings {
		r := idx[f.Repo]
		if r == nil {
			r = &repoRow{repo: f.Repo}
			idx[f.Repo] = r
		}
		r.total++
		switch f.Severity {
		case "critical":
			r.crit++
		case "high":
			r.high++
		case "medium":
			r.med++
		case "low":
			r.low++
		}
	}
	rows := make([]repoRow, 0, len(idx))
	for _, r := range idx {
		rows = append(rows, *r)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].crit != rows[j].crit {
			return rows[i].crit > rows[j].crit
		}
		if rows[i].high != rows[j].high {
			return rows[i].high > rows[j].high
		}
		return rows[i].repo < rows[j].repo
	})
	return rows
}
```

- [ ] **Step 4: Generate the golden file, then run the test**

Run: `UPDATE_GOLDEN=1 go test ./internal/render/... && go test ./internal/render/...`
Expected: first run writes `testdata/dashboard.md.golden`; second run PASSES. Open the golden file and eyeball it for correctness (headers, focus list, waterfall row, per-repo row, errors note).

- [ ] **Step 5: Commit**

```bash
git add internal/render/render.go internal/render/render_test.go internal/render/testdata/
git commit -m "feat: render dashboard markdown + json"
```

---

### Task 14: `render` — tracking-issue upsert + phase wiring

**Files:**
- Create: `internal/render/issue.go`
- Test: `internal/render/issue_test.go`
- Modify: `cmd/ksec/main.go` (wire `render`)

**Interfaces:**
- Consumes: `ghclient.GitHub`, `DashboardMarkdown` from Task 13.
- Produces: `func UpsertTrackingIssue(gh ghclient.GitHub, repo, body string, dryRun bool) (int, error)` — when `dryRun`, prints the intended action and returns `0,nil`; otherwise calls `gh.UpsertIssue` with the fixed marker `<!-- ksec:dashboard -->`, title `Kairos Security Dashboard`, labels `security`,`kairos-security-bot`.

- [ ] **Step 1: Write the failing test**

Create `internal/render/issue_test.go`:

```go
package render

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpsertTrackingIssueWrites(t *testing.T) {
	gh := ghclient.NewFake()
	n, err := UpsertTrackingIssue(gh, "kairos-io/kairos", "body", false)
	require.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, "body", gh.Issues["kairos-io/kairos"].Body)
	assert.Contains(t, gh.Issues["kairos-io/kairos"].Labels, "kairos-security-bot")
}

func TestUpsertTrackingIssueDryRunSkips(t *testing.T) {
	gh := ghclient.NewFake()
	n, err := UpsertTrackingIssue(gh, "kairos-io/kairos", "body", true)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Empty(t, gh.Issues)
}
```

- [ ] **Step 2: Run it; expect FAIL.** Run: `go test ./internal/render/...`

- [ ] **Step 3: Implement**

Create `internal/render/issue.go`:

```go
package render

import (
	"fmt"

	"github.com/kairos-io/security/internal/ghclient"
)

const (
	IssueMarker = "<!-- ksec:dashboard -->"
	IssueTitle  = "Kairos Security Dashboard"
)

var IssueLabels = []string{"security", "kairos-security-bot"}

func UpsertTrackingIssue(gh ghclient.GitHub, repo, body string, dryRun bool) (int, error) {
	if dryRun {
		fmt.Printf("[dry-run] would upsert tracking issue in %s (%d bytes)\n", repo, len(body))
		return 0, nil
	}
	return gh.UpsertIssue(repo, IssueMarker, IssueTitle, body, IssueLabels)
}
```

- [ ] **Step 4: Run it; expect PASS.** Run: `go test ./internal/render/...`

- [ ] **Step 5: Wire the `render` subcommand**

Add to `cmd/ksec/main.go`, register with `root.AddCommand(newRenderCmd(gf))`:

```go
func newRenderCmd(gf *globalFlags) *cobra.Command {
	var trackingRepo string
	cmd := &cobra.Command{
		Use:   "render",
		Short: "write dashboard files and upsert the tracking issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			var c state.Correlated
			if err := state.Load(gf.stateDir, state.CorrelatedFile, &c); err != nil {
				return err
			}
			var tr state.Triage
			if err := state.Load(gf.stateDir, state.TriageFile, &tr); err != nil {
				return err
			}
			var findings state.Findings
			_ = state.Load(gf.stateDir, state.FindingsFile, &findings)

			in := render.Input{
				Correlated:    c,
				Triage:        tr,
				CollectErrors: findings.Errors,
				RunURL:        os.Getenv("KSEC_RUN_URL"),
			}
			md := render.DashboardMarkdown(in)
			j, err := render.DashboardJSON(in)
			if err != nil {
				return err
			}
			if err := os.WriteFile("dashboard.md", []byte(md), 0o644); err != nil {
				return err
			}
			if err := os.WriteFile("dashboard.json", j, 0o644); err != nil {
				return err
			}
			_, err = render.UpsertTrackingIssue(ghclient.NewCLI(), trackingRepo, md, gf.dryRun)
			return err
		},
	}
	cmd.Flags().StringVar(&trackingRepo, "tracking-repo", "kairos-io/kairos", "repo to upsert the tracking issue into")
	return cmd
}
```

Add import `github.com/kairos-io/security/internal/render`.

- [ ] **Step 6: Build + full test suite + commit**

Run: `go build ./... && go test ./...`
Expected: all PASS.

```bash
git add internal/render/issue.go internal/render/issue_test.go cmd/ksec/main.go
git commit -m "feat: render tracking-issue upsert and wire render phase"
```

---

### Task 15: GHA workflow + LocalAI/nib setup + end-to-end dry-run

**Files:**
- Create: `.github/workflows/security-dashboard.yaml`
- Create: `internal/e2e/pipeline_test.go`

**Interfaces:**
- Consumes: the entire CLI built in Tasks 1–14.
- Produces: a scheduled workflow that runs the five phases live by default (dry-run on dispatch/forks) and an end-to-end test that runs every phase in-process against fakes/stubs and asserts the surfaces are produced without any network or repo writes.

- [ ] **Step 1: Write the end-to-end dry-run test**

Create `internal/e2e/pipeline_test.go`:

```go
package e2e

import (
	"testing"

	"github.com/kairos-io/security/internal/correlate"
	"github.com/kairos-io/security/internal/render"
	"github.com/kairos-io/security/internal/state"
	"github.com/kairos-io/security/internal/triage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// failingAI forces the deterministic fallback so the test needs no model.
type failingAI struct{}

func (failingAI) Summarize(state.Correlated) ([]string, map[string]string, string, error) {
	return nil, nil, "", assert.AnError
}

func TestPipelineCorrelateTriageRenderProducesSurfaces(t *testing.T) {
	findings := state.Findings{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high", FixedVersion: "0.33.0", FirstSeen: "2026-06-01", LastSeen: "2026-06-19"},
		{ID: "b", Repo: "kairos-io/kairos-agent", Type: "sourceCVE", CVEID: "CVE-2025-1", Package: "golang.org/x/net", Ecosystem: "go", Severity: "high", FixedVersion: "0.33.0", FirstSeen: "2026-06-10", LastSeen: "2026-06-19"},
	}}

	c := correlate.Run(findings)
	require.Len(t, c.Waterfall, 1, "two repos sharing a go CVE form a waterfall front")

	tr := triage.Run(c, failingAI{}, "test-model")
	assert.False(t, tr.AIAvailable)
	assert.NotEmpty(t, tr.Focus)

	md := render.DashboardMarkdown(render.Input{Correlated: c, Triage: tr})
	assert.Contains(t, md, "Waterfall fronts")
	assert.Contains(t, md, "golang.org/x/net@0.33.0")

	j, err := render.DashboardJSON(render.Input{Correlated: c, Triage: tr})
	require.NoError(t, err)
	assert.Contains(t, string(j), "CVE-2025-1")
}
```

- [ ] **Step 2: Run it; expect PASS** (all referenced code already exists).

Run: `go test ./internal/e2e/...`

- [ ] **Step 3: Write the workflow**

Create `.github/workflows/security-dashboard.yaml`:

```yaml
name: Security Dashboard

on:
  workflow_dispatch:
    inputs:
      dry_run:
        description: "Print intended writes instead of performing them"
        type: boolean
        default: false
  schedule:
    - cron: "0 6 * * *"

concurrency:
  group: security-dashboard
  cancel-in-progress: false

permissions:
  contents: write   # commit state + dashboards back to this repo

jobs:
  dashboard:
    runs-on: ubuntu-latest
    env:
      # Scheduled runs are LIVE by default; dispatch can request dry-run; forks force it.
      DRYRUN: ${{ (github.event.inputs.dry_run == 'true' || github.event.pull_request.head.repo.fork) && '--dry-run' || '' }}
      KSEC_RUN_URL: ${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}
      # Cross-repo reads + the kairos tracking-issue write use the scoped bot token.
      GH_TOKEN: ${{ secrets.KSEC_BOT_TOKEN }}
      LOCALAI_URL: http://localhost:8080
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: Install scanners
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin

      - name: Start LocalAI (small model)
        run: |
          MODEL="${LOCALAI_MODEL:-$(yq '.localai.model.name' ai.yaml)}"
          docker run -d --name localai -p 8080:8080 \
            -e MODELS_PATH=/models \
            quay.io/go-skynet/local-ai:${LOCALAI_VERSION:-latest} \
            "$MODEL"
          for i in $(seq 1 60); do
            curl -sf "$LOCALAI_URL/readyz" && break || sleep 5
          done

      - name: Install nib
        run: go install github.com/mudler/nib@${NIB_VERSION:-latest}

      - name: Build ksec
        run: go build -o /usr/local/bin/ksec ./cmd/ksec

      - name: Run pipeline
        run: |
          ksec discover  --state-dir state
          ksec collect   --state-dir state
          ksec correlate --state-dir state
          ksec triage    --state-dir state
          ksec render    --state-dir state $DRYRUN

      - name: Commit state + dashboards
        if: ${{ env.DRYRUN == '' }}
        run: |
          git config user.name "kairos-security-bot"
          git config user.email "bot@kairos.io"
          git add state/ dashboard.md dashboard.json
          if git diff --cached --quiet; then
            echo "no changes"
          else
            git commit -m "chore: update security dashboard [skip ci]"
            git push
          fi
```

> Notes for the implementer: `KSEC_BOT_TOKEN` is the scoped GitHub App / fine-grained PAT secret (read on all tracked repos + issues:write on `kairos-io/kairos`). `yq` is preinstalled on `ubuntu-latest`. The exact LocalAI model/run invocation may need adjusting to the chosen small model (a §16 open item from the spec) — keep the `ai.yaml` handles as the single source of truth.

- [ ] **Step 4: Validate the workflow YAML locally**

Run: `python -c "import yaml,sys; yaml.safe_load(open('.github/workflows/security-dashboard.yaml'))" && echo OK`
Expected: `OK`.

- [ ] **Step 5: Final full build + test + commit**

Run: `go build ./... && go test ./...`
Expected: all PASS.

```bash
git add .github/workflows/security-dashboard.yaml internal/e2e/
git commit -m "feat: security-dashboard workflow and end-to-end dry-run test"
```

---

## Self-review

**Spec coverage** (against `2026-06-19-...-design.md`):

- §5 phases discover/collect/correlate/triage/render → Tasks 5, 6–10, 11, 12, 13–14. ✓
- §6.1 hybrid discovery (org enum + kairos-init parse + repos.yaml) → Task 5. ✓
- §6.2 four collectors (prs, imageCVE, sourceCVE, ghAlerts) + per-repo error isolation → Tasks 6–10. ✓
- §6.3 correlate dedupe + waterfall → Task 11. ✓
- §6.4 triage AI + deterministic fallback (`aiAvailable`) → Task 12. ✓
- §6.6 render dashboard + tracking issue → Tasks 13–14. ✓
- §7 state files (repos/findings/correlated/triage) → Tasks 2, 5, 10, 11, 12. Note: `findings.json` is the `Findings{findings,errors}` struct (refines the spec's bare array to carry collection errors; ledger.json is Plan 2). ✓
- §9 dashboard.md sections (focus, waterfall, per-repo) + single tracking issue by marker → Tasks 13–14. ✓
- §10 error handling: per-collector isolation (Task 10), AI best-effort fallback (Task 12), secrets not logged (gh/nib clients never echo tokens). ✓
- §11 GHA orchestration, live-by-default dry-run, LocalAI service, idempotent commit → Task 15. ✓
- §12 AI config handles (model select/preload, nib→LocalAI wiring, pin/install) → Task 3 (`ai.yaml`) + Task 15 (workflow consumes them). ✓
- §13 config files repos.yaml/ai.yaml → Task 3. ✓
- §14 testing: pure-phase golden/table tests, collector fixtures, AI mock, E2E dry-run → Tasks 2–14 tests + Task 15 E2E. ✓

**Out of scope (correctly deferred to Plan 2):** `remediate`, `ledger.json`, PR creation/rebasing, comment reactions, blast-radius cap. These appear in the spec (§6.5, §8) but are explicitly Plan 2. The read-only pipeline here is independently shippable.

**Placeholder scan:** No `TODO`/`TBD`/"add error handling"/"similar to Task N" left; every code step shows complete code.

**Type consistency:** `state.Finding`/`Findings`/`Correlated`/`Triage`/`WaterfallGroup`/`Bump` defined in Task 2 are used unchanged in Tasks 5–14. `ghclient.GitHub` (Task 4) consumed by Tasks 5, 8, 9, 14. `collect.Collector`, `FindingID`, `Today`/`nowFn` (Task 6) reused by Tasks 7–10. `triage.AIClient` (Task 12) consumed by Task 15's E2E. `render.Input` (Task 13) consumed by Task 14 and Task 15. Subcommand builders `newDiscoverCmd`/`newCollectCmd`/`newCorrelateCmd`/`newTriageCmd`/`newRenderCmd` all registered on the root in `cmd/ksec/main.go`.

---

## Open items carried from the spec (resolve during execution / Plan 2)

- Exact scoped-token mechanism (GitHub App install token vs fine-grained PAT) and minimal per-repo permission set for `KSEC_BOT_TOKEN`.
- Concrete small LocalAI model id + gallery reference + runner resource footprint; adjust the `Start LocalAI` step accordingly.
- Retirement of legacy `scan.yaml`/`autobump.yaml`/`automerge.yaml` once this reaches parity (the `imageCVE` collector subsumes the framework-image scan via `repos.yaml` artifacts).
