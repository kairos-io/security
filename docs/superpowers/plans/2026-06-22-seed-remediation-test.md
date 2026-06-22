# Seed Remediation Test Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `ksec remediate --seed <owner/repo>=<package>@<version>` (repeatable) to inject a synthetic finding before planning, so the full plan → fork → PR path can be exercised on demand (dry-run prints the would-be PR; live opens one).

**Architecture:** A pure `ParseSeed` turns the spec into a synthetic `sourceCVE` finding satisfying `actionable()`; `newRemediateCmd` appends parsed seeds to the loaded correlated findings before `Plan`. Everything downstream (Plan cap, fork rules, dry-run) is unchanged. Builds on `internal/remediate`, `internal/state`, `cmd/ksec`.

**Tech Stack:** Go 1.22, `stretchr/testify`, cobra, existing remediate engine.

## Global Constraints

- Module `github.com/kairos-io/security`; Go 1.22.
- A seed is a planning-time injection only — never persisted to state.
- Seeds flow through the real engine: `--dry-run` (no writes), `--max-prs` (cap), and fork rules (external → fork PR) all apply unchanged.
- The synthetic finding must satisfy `actionable()`: `Type:"sourceCVE"`, `Ecosystem:"go"`, non-empty non-`stdlib` `Package`, non-empty `FixedVersion`.

---

## File structure

```
internal/remediate/seed.go        # ParseSeed (create)
internal/remediate/seed_test.go   # (create)
internal/remediate/planner_test.go # assert a seed finding plans an intent (modify)
cmd/ksec/main.go                  # --seed flag; append to c.Findings before Plan (modify)
```

---

### Task 1: ParseSeed + planner coverage

**Files:** Create `internal/remediate/seed.go`, `internal/remediate/seed_test.go`; modify `internal/remediate/planner_test.go`.

**Interfaces:** `func ParseSeed(spec string) (state.Finding, error)`.

- [ ] **Step 1: Write the failing tests**

Create `internal/remediate/seed_test.go`:

```go
package remediate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSeed(t *testing.T) {
	f, err := ParseSeed("mudler/edgevpn=golang.org/x/net@0.33.0")
	require.NoError(t, err)
	assert.Equal(t, "mudler/edgevpn", f.Repo)
	assert.Equal(t, "golang.org/x/net", f.Package) // slashes preserved
	assert.Equal(t, "0.33.0", f.FixedVersion)
	assert.Equal(t, "sourceCVE", f.Type)
	assert.Equal(t, "go", f.Ecosystem)
	assert.Equal(t, "high", f.Severity)
	assert.True(t, actionable(f)) // the planner will turn it into a target
}

func TestParseSeedErrors(t *testing.T) {
	for _, bad := range []string{"", "no-eq", "repo=pkgNoVersion", "=pkg@1", "r=@1", "r=pkg@"} {
		_, err := ParseSeed(bad)
		assert.Error(t, err, "spec %q should error", bad)
	}
}
```

Add to `planner_test.go` (a seed finding produces an actionable intent):

```go
func TestPlanActionsSeedFinding(t *testing.T) {
	f, err := ParseSeed("mudler/edgevpn=golang.org/x/net@0.33.0")
	require.NoError(t, err)
	intents, _ := Plan(state.Correlated{Findings: []state.Finding{f}}, state.Ledger{}, nil, nil, 10)
	var got *Intent
	for i := range intents {
		if intents[i].Key == "mudler/edgevpn|golang.org/x/net" {
			got = &intents[i]
		}
	}
	require.NotNil(t, got, "seed finding must produce an intent")
	assert.Equal(t, "0.33.0", got.Bump.To)
}
```

(Confirm `require` is imported in `planner_test.go`; add if missing.)

- [ ] **Step 2: Run red.** `go test ./internal/remediate/...`

- [ ] **Step 3: Implement** — `internal/remediate/seed.go`:

```go
package remediate

import (
	"fmt"
	"strings"

	"github.com/kairos-io/security/internal/state"
)

// ParseSeed turns "owner/repo=package@version" into a synthetic sourceCVE
// finding, so an operator can exercise the remediation pipeline on demand:
//
//	ksec remediate --seed mudler/edgevpn=golang.org/x/net@0.33.0 --dry-run
//
// The fields are exactly those actionable() requires, so the planner turns the
// seed into a real bump intent. Seeds are planning-time only (never persisted).
func ParseSeed(spec string) (state.Finding, error) {
	repo, rest, ok := strings.Cut(spec, "=")
	if !ok || repo == "" {
		return state.Finding{}, fmt.Errorf("seed %q: want owner/repo=package@version", spec)
	}
	pkg, version, ok := strings.Cut(rest, "@")
	if !ok || pkg == "" || version == "" {
		return state.Finding{}, fmt.Errorf("seed %q: want owner/repo=package@version", spec)
	}
	return state.Finding{
		ID:           "seed:" + repo + "|" + pkg,
		Repo:         repo,
		Type:         "sourceCVE",
		Ecosystem:    "go",
		Package:      pkg,
		FixedVersion: version,
		Severity:     "high",
		Source:       "seed",
		CVEID:        "SEED",
		Title:        "synthetic seed finding (remediation test)",
	}, nil
}
```

- [ ] **Step 4: Run green + build + commit**

Run: `go test ./internal/remediate/... && go build ./...`
```bash
git add internal/remediate/seed.go internal/remediate/seed_test.go internal/remediate/planner_test.go
git commit -m "feat(remediate): ParseSeed — synthetic finding for pipeline testing"
```

---

### Task 2: Wire the --seed flag

**Files:** Modify `cmd/ksec/main.go`.

**Interfaces:** `newRemediateCmd` gains a repeatable `--seed` string flag; parsed seeds are appended to the loaded correlated findings before `Plan`.

- [ ] **Step 1: Add the flag + injection**

In `newRemediateCmd`, declare the flag (near the other flag vars):

```go
	var seeds []string
```
and register it (near the other `cmd.Flags()` calls):
```go
	cmd.Flags().StringArrayVar(&seeds, "seed", nil,
		"inject a synthetic finding to test remediation: owner/repo=package@version (repeatable)")
```

In the `RunE`, immediately after the correlated state is loaded (`state.Load(gf.stateDir, state.CorrelatedFile, &c)`), inject the seeds before `Plan`:

```go
			for _, s := range seeds {
				f, err := remediate.ParseSeed(s)
				if err != nil {
					return err
				}
				c.Findings = append(c.Findings, f)
			}
			if len(seeds) > 0 {
				fmt.Fprintf(os.Stderr, "remediate: injected %d seed finding(s) for testing\n", len(seeds))
			}
```

(`fmt`/`os`/`remediate` are already imported.)

- [ ] **Step 2: Build + vet + gofmt + test + smoke**

Run: `go build ./... && go vet ./... && test -z "$(gofmt -l .)" && go test ./...`
Smoke:
```bash
go run ./cmd/ksec remediate --help        # shows --seed
```
Confirm `--seed` appears in help with the repeatable usage.

- [ ] **Step 3: Commit**

```bash
git add cmd/ksec/main.go
git commit -m "feat(remediate): --seed flag injects synthetic findings before planning"
```

---

## Self-review

**Spec coverage:**
- `--seed owner/repo=package@version` repeatable → Task 2. ✓
- `ParseSeed` synthetic sourceCVE finding satisfying `actionable()` → Task 1. ✓
- Append before `Plan`; honor `--dry-run`/`--max-prs`/fork (no special-casing) → Task 2 (injection only). ✓
- Malformed spec → clear error → Task 1. ✓
- Planning-time only (never persisted) → Task 2 (appends to in-memory `c.Findings`, not saved). ✓

**Placeholder scan:** none — full code for the pure parser (Task 1) and the exact flag/injection lines (Task 2).

**Type consistency:** `ParseSeed` returns `state.Finding` with the fields `actionable()`/target-build read (`Repo`/`Package`/`FixedVersion`/`Severity`/`Type`/`Ecosystem`); the command appends to `c.Findings` ([]state.Finding) before the existing `Plan(c, …)` call.

---

## Operational notes

- Safe test: `ksec remediate --seed mudler/edgevpn=golang.org/x/net@0.33.0 --dry-run --state-dir state` → prints the planned fork PR, zero writes.
- Live demo: drop `--dry-run` → opens one real PR (on the bot's fork for an external repo), capped by `--max-prs`.
- Seed a **tracked** repo (in `repos.yaml`) so fork-detection and PR-matching apply; seeding an untracked repo push-directs and may fail.
