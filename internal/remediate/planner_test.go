package remediate

import (
	"testing"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanOpensNewActionableTargetsDedupedAndCapped(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		// two CVEs in the same repo+package -> one target at the highest fixed version
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
		{ID: "b", Repo: "kairos-io/immucore", Type: "ghAlert", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.36.0", Severity: "critical"},
		// a different repo+package -> second target
		{ID: "c", Repo: "kairos-io/kairos", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/crypto", FixedVersion: "0.31.0", Severity: "high"},
		// not actionable: image CVE
		{ID: "d", Repo: "kairos-io/kairos", Type: "imageCVE", Package: "openssl", FixedVersion: "1.1.1w", Severity: "critical"},
		// not actionable: no fixed version
		{ID: "e", Repo: "kairos-io/kairos", Type: "sourceCVE", Ecosystem: "go", Package: "x/text", Severity: "low"},
	}}

	intents, deferred := Plan(c, state.Ledger{}, nil, nil, 1) // cap to 1 new PR
	require.Len(t, intents, 1)
	assert.Equal(t, 1, deferred)
	in := intents[0]
	assert.Equal(t, IntentOpen, in.Type)
	// highest severity target first: immucore/x/net (critical) at the highest fixed version
	assert.Equal(t, "kairos-io/immucore|golang.org/x/net", in.Key)
	assert.Equal(t, "0.36.0", in.Bump.To)
	assert.Equal(t, "critical", in.Severity)
}

func TestPlanReconcilesExistingLedgerEntries(t *testing.T) {
	c := state.Correlated{}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", Repo: "kairos-io/immucore", State: "open"},
	}}
	intents, _ := Plan(c, led, nil, nil, 10)
	require.Len(t, intents, 1)
	assert.Equal(t, IntentReconcile, intents[0].Type)
	require.NotNil(t, intents[0].Entry)
	assert.Equal(t, "open", intents[0].Entry.State)
}

func TestPlanSkipsTargetsAlreadyInLedger(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|golang.org/x/net", State: "open"},
	}}
	intents, _ := Plan(c, led, nil, nil, 10)
	// only the reconcile for the existing entry; no new open
	require.Len(t, intents, 1)
	assert.Equal(t, IntentReconcile, intents[0].Type)
}

// intentFor returns the first intent of the given type for a key, or nil.
func intentFor(intents []Intent, typ IntentType, key string) *Intent {
	for i := range intents {
		if intents[i].Type == typ && intents[i].Key == key {
			return &intents[i]
		}
	}
	return nil
}

// A non-live ledger state (planned, from a prior dry-run) must NOT permanently
// suppress the key: going live should re-open the PR while still reconciling.
func TestPlanReopensPlannedLedgerEntry(t *testing.T) {
	k := "kairos-io/immucore|golang.org/x/net"
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: k, Repo: "kairos-io/immucore", State: "planned", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	intents, _ := Plan(c, led, nil, nil, 10)
	require.NotNil(t, intentFor(intents, IntentReconcile, k), "expected reconcile for the existing entry")
	open := intentFor(intents, IntentOpen, k)
	require.NotNil(t, open, "expected re-open for the planned entry")
	assert.Equal(t, "0.33.0", open.Bump.To)
}

// A transient build-failed entry must retry (re-open) on a later run.
func TestPlanReopensBuildFailedLedgerEntry(t *testing.T) {
	k := "kairos-io/immucore|golang.org/x/net"
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: k, Repo: "kairos-io/immucore", State: "build-failed", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	intents, _ := Plan(c, led, nil, nil, 10)
	require.NotNil(t, intentFor(intents, IntentOpen, k), "expected re-open for the build-failed entry")
}

// An open entry has a live PR maintained via reconcile: no new open.
func TestPlanSkipsOpenLedgerEntryButReconciles(t *testing.T) {
	k := "kairos-io/immucore|golang.org/x/net"
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	led := state.Ledger{Entries: []state.LedgerEntry{
		{Key: k, Repo: "kairos-io/immucore", State: "open", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	intents, _ := Plan(c, led, nil, nil, 10)
	require.NotNil(t, intentFor(intents, IntentReconcile, k))
	assert.Nil(t, intentFor(intents, IntentOpen, k), "open entry must not be re-opened")
}

// A merged entry re-opens only when a NEWER fixed version is later required.
func TestPlanReopensMergedOnlyForHigherVersion(t *testing.T) {
	k := "kairos-io/immucore|golang.org/x/net"

	// merged at 0.33.0, finding needs 0.36.0 -> re-open (re-bump).
	cHigher := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.36.0", Severity: "high"},
	}}
	ledLow := state.Ledger{Entries: []state.LedgerEntry{
		{Key: k, Repo: "kairos-io/immucore", State: "merged", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	intents, _ := Plan(cHigher, ledLow, nil, nil, 10)
	open := intentFor(intents, IntentOpen, k)
	require.NotNil(t, open, "merged at lower version must re-open for the newer fix")
	assert.Equal(t, "0.36.0", open.Bump.To)

	// merged at 0.36.0, finding needs 0.33.0 -> skip (already addressed).
	cLower := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	ledHigh := state.Ledger{Entries: []state.LedgerEntry{
		{Key: k, Repo: "kairos-io/immucore", State: "merged", Bump: state.Bump{Package: "golang.org/x/net", To: "0.36.0"}},
	}}
	intents2, _ := Plan(cLower, ledHigh, nil, nil, 10)
	assert.Nil(t, intentFor(intents2, IntentOpen, k), "merged at >= version must not re-open")
}

func TestPlanAdoptsExistingExternalPR(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go",
			Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	prs := map[string][]ghclient.PullRequest{
		"kairos-io/immucore": {{Number: 7, Title: "Bump golang.org/x/net to 0.33.0", Author: "renovate[bot]", URL: "u7"}},
	}
	intents, _ := Plan(c, state.Ledger{}, prs, nil, 10)
	require.Len(t, intents, 1)
	assert.Equal(t, IntentAdopt, intents[0].Type)
	assert.Equal(t, 7, intents[0].PRNumber)
	assert.Equal(t, "renovate", intents[0].Source)
}

func TestPlanOpensWhenNoExternalPR(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "a", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go",
			Package: "golang.org/x/net", FixedVersion: "0.33.0", Severity: "high"},
	}}
	intents, _ := Plan(c, state.Ledger{}, nil, nil, 10)
	require.Len(t, intents, 1)
	assert.Equal(t, IntentOpen, intents[0].Type)
}

func TestPlanCascadesMergedFirstPartyFix(t *testing.T) {
	repos := []state.Repo{
		{Repo: "kairos-io/kairos-sdk", Branch: "main"},
		{Repo: "kairos-io/immucore", Branch: "master"},
	}
	gomod := map[string][]byte{
		"kairos-io/kairos-sdk": []byte("module github.com/kairos-io/kairos-sdk\n"),
		"kairos-io/immucore":   []byte("module github.com/kairos-io/immucore\nrequire github.com/kairos-io/kairos-sdk v0.7.0\n"),
	}
	g := BuildGraph(repos, gomod)
	// A merged fix in the sdk repo -> cascade a pseudo bump into immucore.
	ledger := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/kairos-sdk|golang.org/x/net", Repo: "kairos-io/kairos-sdk", State: "merged",
			Kind: "direct", Severity: "high", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
	}}
	intents, _ := Plan(state.Correlated{}, ledger, nil, g, 10)

	var cas *Intent
	for i := range intents {
		if intents[i].Type == IntentCascade {
			cas = &intents[i]
		}
	}
	require.NotNil(t, cas, "expected a cascade intent")
	assert.Equal(t, "kairos-io/immucore", cas.Repo)
	assert.Equal(t, "github.com/kairos-io/kairos-sdk", cas.Package)
	assert.Equal(t, "main", cas.Ref) // sdk's default branch for the pseudo go get
	assert.Equal(t, "kairos-io/kairos-sdk|golang.org/x/net", cas.CascadeFrom)
}

// A maintainer who closes a cascade PR must not be fought: a closed cascade
// ledger entry for the consumer suppresses re-creation.
func TestPlanSkipsClosedCascade(t *testing.T) {
	repos := []state.Repo{
		{Repo: "kairos-io/kairos-sdk", Branch: "main"},
		{Repo: "kairos-io/immucore", Branch: "master"},
	}
	gomod := map[string][]byte{
		"kairos-io/kairos-sdk": []byte("module github.com/kairos-io/kairos-sdk\n"),
		"kairos-io/immucore":   []byte("module github.com/kairos-io/immucore\nrequire github.com/kairos-io/kairos-sdk v0.7.0\n"),
	}
	g := BuildGraph(repos, gomod)
	ck := "kairos-io/immucore|github.com/kairos-io/kairos-sdk"
	ledger := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/kairos-sdk|golang.org/x/net", Repo: "kairos-io/kairos-sdk", State: "merged",
			Kind: "direct", Severity: "high", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
		{Key: ck, Repo: "kairos-io/immucore", Package: "github.com/kairos-io/kairos-sdk",
			Kind: "cascade", State: "closed"},
	}}
	intents, _ := Plan(state.Correlated{}, ledger, nil, g, 10)
	assert.Nil(t, intentFor(intents, IntentCascade, ck),
		"a human-closed cascade must not be re-created")
}

// Two merged fixes in the SAME upstream repo must cascade a shared consumer at
// most once per run.
func TestPlanCascadeDedupsPerRun(t *testing.T) {
	repos := []state.Repo{
		{Repo: "kairos-io/kairos-sdk", Branch: "main"},
		{Repo: "kairos-io/immucore", Branch: "master"},
	}
	gomod := map[string][]byte{
		"kairos-io/kairos-sdk": []byte("module github.com/kairos-io/kairos-sdk\n"),
		"kairos-io/immucore":   []byte("module github.com/kairos-io/immucore\nrequire github.com/kairos-io/kairos-sdk v0.7.0\n"),
	}
	g := BuildGraph(repos, gomod)
	ck := "kairos-io/immucore|github.com/kairos-io/kairos-sdk"
	ledger := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/kairos-sdk|golang.org/x/net", Repo: "kairos-io/kairos-sdk", State: "merged",
			Kind: "direct", Severity: "high", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}},
		{Key: "kairos-io/kairos-sdk|golang.org/x/crypto", Repo: "kairos-io/kairos-sdk", State: "merged",
			Kind: "direct", Severity: "high", Bump: state.Bump{Package: "golang.org/x/crypto", To: "0.31.0"}},
	}}
	intents, _ := Plan(state.Correlated{}, ledger, nil, g, 10)
	n := 0
	for _, in := range intents {
		if in.Type == IntentCascade && in.Key == ck {
			n++
		}
	}
	assert.Equal(t, 1, n, "consumer must be cascaded at most once per run")
}

func TestPlanRepinsPseudoCascade(t *testing.T) {
	ledger := state.Ledger{Entries: []state.LedgerEntry{
		{Key: "kairos-io/immucore|github.com/kairos-io/kairos-sdk", Repo: "kairos-io/immucore",
			Package: "github.com/kairos-io/kairos-sdk", State: "open", Kind: "cascade", Pseudo: true},
	}}
	intents, _ := Plan(state.Correlated{}, ledger, nil, nil, 10)
	var found bool
	for _, in := range intents {
		if in.Type == IntentRepin && in.Key == "kairos-io/immucore|github.com/kairos-io/kairos-sdk" {
			found = true
		}
	}
	assert.True(t, found, "expected a repin intent for the pseudo cascade entry")
}

func TestPlanToolchainForStdlib(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "s", Repo: "kairos-io/immucore", Type: "sourceCVE", Ecosystem: "go",
			Package: "stdlib", FixedVersion: "go1.22.5", Severity: "high"},
	}}
	intents, _ := Plan(c, state.Ledger{}, nil, nil, 10)
	var tc *Intent
	for i := range intents {
		if intents[i].Type == IntentToolchain {
			tc = &intents[i]
		}
	}
	require.NotNil(t, tc)
	assert.Equal(t, "kairos-io/immucore", tc.Repo)
	assert.Equal(t, "1.22.5", tc.ToolchainVersion) // leading "go" stripped
	assert.Equal(t, "kairos-io/immucore|go-toolchain", tc.Key)
}

func TestPlanSupersedesConflictedAdoptedPR(t *testing.T) {
	c := state.Correlated{Findings: []state.Finding{
		{ID: "f1", Repo: "o/r", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net",
			FixedVersion: "0.33.0", Severity: "high"},
	}}
	ledger := state.Ledger{Entries: []state.LedgerEntry{{
		Key: "o/r|golang.org/x/net", Repo: "o/r", Package: "golang.org/x/net", State: "open",
		Source: "bot", Blocked: "upstream-conflict", PRNumber: 38, PRURL: "https://github.com/o/r/pull/38",
		Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"},
	}}}
	intents, _ := Plan(c, ledger, nil, nil, 10)
	var sup *Intent
	for i := range intents {
		if intents[i].Type == IntentSupersede {
			sup = &intents[i]
		}
	}
	require.NotNil(t, sup)
	assert.Equal(t, "o/r|golang.org/x/net", sup.Key)
	assert.Equal(t, 38, sup.PRNumber)
	assert.Equal(t, "https://github.com/o/r/pull/38", sup.PRURL)
}

func TestPlanTerminalKsecEntryNotReadopted(t *testing.T) {
	// A ksec entry that build-failed (needs human) must NOT be re-adopted or
	// re-opened even though a matching external bot PR exists for the target.
	c := state.Correlated{Findings: []state.Finding{
		{ID: "f1", Repo: "o/r", Type: "sourceCVE", Ecosystem: "go", Package: "golang.org/x/net",
			FixedVersion: "0.33.0", Severity: "high"},
	}}
	prs := map[string][]ghclient.PullRequest{"o/r": {
		{Number: 9, Title: "bump golang.org/x/net to 0.33.0", Author: "dependabot[bot]", URL: "u9"},
	}}
	ledger := state.Ledger{Entries: []state.LedgerEntry{{
		Key: "o/r|golang.org/x/net", Repo: "o/r", Package: "golang.org/x/net", State: "build-failed",
		Source: "ksec", NeedsHuman: true, Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"},
	}}}
	intents, _ := Plan(c, ledger, prs, nil, 10)
	for _, in := range intents {
		if in.Key != "o/r|golang.org/x/net" {
			continue
		}
		// A bare reconcile is expected for every ledger entry; what must NOT
		// happen is a re-adopt / re-open / re-supersede of the terminal entry.
		if in.Type == IntentAdopt || in.Type == IntentOpen || in.Type == IntentSupersede {
			t.Fatalf("terminal ksec entry must not be re-actioned, got %s", in.Type)
		}
	}
}

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
