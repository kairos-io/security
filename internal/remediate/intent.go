package remediate

import "github.com/kairos-io/security/internal/state"

type IntentType string

const (
	IntentOpen      IntentType = "open"
	IntentReconcile IntentType = "reconcile"
)

type Intent struct {
	Type     IntentType
	Key      string
	Repo     string
	Package  string
	Severity string
	Bump     state.Bump
	Entry    *state.LedgerEntry // set for IntentReconcile
}

type Result struct {
	Key    string
	Action string
	State  string
	Detail string
}
