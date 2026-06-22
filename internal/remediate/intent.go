package remediate

import "github.com/kairos-io/security/internal/state"

type IntentType string

const (
	IntentOpen      IntentType = "open"
	IntentReconcile IntentType = "reconcile"
	IntentAdopt     IntentType = "adopt"
	IntentCascade   IntentType = "cascade"
	IntentRepin     IntentType = "repin"
	IntentToolchain IntentType = "toolchain"
	IntentSupersede IntentType = "supersede"
)

type Intent struct {
	Type             IntentType
	Key              string
	Repo             string
	Package          string
	Severity         string
	Bump             state.Bump
	Entry            *state.LedgerEntry // set for IntentReconcile
	PRNumber         int
	PRURL            string
	Source           string // dependabot | renovate | human (for IntentAdopt)
	Ref              string // module's default branch for the pseudo `go get` (IntentCascade)
	CascadeFrom      string // upstream ledger key that triggered this cascade (IntentCascade)
	ToolchainVersion string // target Go toolchain version, "go" prefix stripped (IntentToolchain)
}

type Result struct {
	Key    string
	Action string
	State  string
	Detail string
}
