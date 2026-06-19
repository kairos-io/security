package remediate

import "github.com/kairos-io/security/internal/state"

// FakeExecutor is an in-memory Executor for tests.
type FakeExecutor struct {
	Opened     map[string]state.LedgerEntry
	Reconciled map[string]state.LedgerEntry
	OpenErr    map[string]error
}

func (f *FakeExecutor) Open(in Intent, run string) (state.LedgerEntry, error) {
	if err := f.OpenErr[in.Key]; err != nil {
		return state.LedgerEntry{}, err
	}
	e, ok := f.Opened[in.Key]
	if !ok {
		return state.LedgerEntry{Key: in.Key, Repo: in.Repo, Package: in.Package, State: "open", CreatedRun: run, LastActionRun: run}, nil
	}
	return e, nil
}

func (f *FakeExecutor) Reconcile(e state.LedgerEntry, run string) (state.LedgerEntry, error) {
	if r, ok := f.Reconciled[e.Key]; ok {
		return r, nil
	}
	return e, nil
}
