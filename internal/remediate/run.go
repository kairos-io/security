package remediate

import (
	"sort"

	"github.com/kairos-io/security/internal/state"
)

type Executor interface {
	Open(in Intent, run string) (state.LedgerEntry, error)
	Reconcile(e state.LedgerEntry, run string) (state.LedgerEntry, error)
	Adopt(in Intent, run string) (state.LedgerEntry, error)
	Cascade(in Intent, run string) (state.LedgerEntry, error)
	Repin(e state.LedgerEntry, run string) (state.LedgerEntry, error)
	Toolchain(in Intent, run string) (state.LedgerEntry, error)
}

func Run(intents []Intent, ex Executor, ledger state.Ledger, run string) (state.Ledger, []Result) {
	// Index existing entries by key for in-place replacement.
	byKey := map[string]state.LedgerEntry{}
	for _, e := range ledger.Entries {
		byKey[e.Key] = e
	}
	var results []Result

	for _, in := range intents {
		switch in.Type {
		case IntentOpen:
			entry, err := ex.Open(in, run)
			if err != nil {
				rec := state.LedgerEntry{
					Key: in.Key, Repo: in.Repo, Package: in.Package, State: "error",
					Bump: in.Bump, Severity: in.Severity, CreatedRun: run, LastActionRun: run,
					History: []state.LedgerEvent{{Run: run, Action: "open-failed", Detail: err.Error()}},
				}
				byKey[in.Key] = rec
				results = append(results, Result{Key: in.Key, Action: "open", State: "error", Detail: err.Error()})
				continue
			}
			byKey[entry.Key] = entry
			results = append(results, Result{Key: entry.Key, Action: "open", State: entry.State})
		case IntentReconcile:
			prior := *in.Entry
			entry, err := ex.Reconcile(prior, run)
			if err != nil {
				prior.LastActionRun = run
				prior.History = append(prior.History, state.LedgerEvent{Run: run, Action: "reconcile-failed", Detail: err.Error()})
				byKey[prior.Key] = prior
				results = append(results, Result{Key: prior.Key, Action: "reconcile", State: "error", Detail: err.Error()})
				continue
			}
			byKey[entry.Key] = entry
			results = append(results, Result{Key: entry.Key, Action: "reconcile", State: entry.State})
		case IntentAdopt:
			entry, err := ex.Adopt(in, run)
			if err != nil {
				rec := state.LedgerEntry{
					Key: in.Key, Repo: in.Repo, Package: in.Package, State: "error",
					Source: in.Source, Kind: "direct", Bump: in.Bump, Severity: in.Severity,
					PRNumber: in.PRNumber, PRURL: in.PRURL, CreatedRun: run, LastActionRun: run,
					History: []state.LedgerEvent{{Run: run, Action: "adopt-failed", Detail: err.Error()}},
				}
				byKey[in.Key] = rec
				results = append(results, Result{Key: in.Key, Action: "adopt", State: "error", Detail: err.Error()})
				continue
			}
			byKey[entry.Key] = entry
			results = append(results, Result{Key: entry.Key, Action: "adopt", State: entry.State})
		case IntentCascade:
			entry, err := ex.Cascade(in, run)
			if err != nil {
				rec := state.LedgerEntry{
					Key: in.Key, Repo: in.Repo, Package: in.Package, State: "error",
					Kind: "cascade", CascadeFrom: in.CascadeFrom, Pseudo: true, Severity: in.Severity,
					CreatedRun: run, LastActionRun: run,
					History: []state.LedgerEvent{{Run: run, Action: "cascade-failed", Detail: err.Error()}},
				}
				byKey[in.Key] = rec
				results = append(results, Result{Key: in.Key, Action: "cascade", State: "error", Detail: err.Error()})
				continue
			}
			byKey[entry.Key] = entry
			results = append(results, Result{Key: entry.Key, Action: "cascade", State: entry.State})
		case IntentToolchain:
			entry, err := ex.Toolchain(in, run)
			if err != nil {
				rec := state.LedgerEntry{
					Key: in.Key, Repo: in.Repo, Package: "go-toolchain", State: "error",
					Kind: "toolchain", Severity: in.Severity, CreatedRun: run, LastActionRun: run,
					Bump:    state.Bump{Package: "go", To: in.ToolchainVersion},
					History: []state.LedgerEvent{{Run: run, Action: "toolchain-failed", Detail: err.Error()}},
				}
				byKey[in.Key] = rec
				results = append(results, Result{Key: in.Key, Action: "toolchain", State: "error", Detail: err.Error()})
				continue
			}
			byKey[entry.Key] = entry
			results = append(results, Result{Key: entry.Key, Action: "toolchain", State: entry.State})
		case IntentRepin:
			// Repin operates on the live (possibly reconciled) entry, not the
			// stale pre-run snapshot: the planner emits both an IntentReconcile
			// and an IntentRepin for each pseudo cascade entry, and reconcile
			// runs first and updates byKey[key].
			prior := *in.Entry
			if live, ok := byKey[in.Key]; ok {
				prior = live
			}
			entry, err := ex.Repin(prior, run)
			if err != nil {
				prior.LastActionRun = run
				prior.History = append(prior.History, state.LedgerEvent{Run: run, Action: "repin-failed", Detail: err.Error()})
				byKey[prior.Key] = prior
				results = append(results, Result{Key: prior.Key, Action: "repin", State: "error", Detail: err.Error()})
				continue
			}
			byKey[entry.Key] = entry
			results = append(results, Result{Key: entry.Key, Action: "repin", State: entry.State})
		}
	}

	out := state.Ledger{Entries: make([]state.LedgerEntry, 0, len(byKey))}
	for _, e := range byKey {
		out.Entries = append(out.Entries, e)
	}
	sort.Slice(out.Entries, func(i, j int) bool { return out.Entries[i].Key < out.Entries[j].Key })
	return out, results
}
