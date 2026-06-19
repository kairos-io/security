package state

// ByKey returns a pointer to the entry with the given key so callers can mutate
// it in place, plus whether it was found.
func (l *Ledger) ByKey(key string) (*LedgerEntry, bool) {
	for i := range l.Entries {
		if l.Entries[i].Key == key {
			return &l.Entries[i], true
		}
	}
	return nil, false
}
