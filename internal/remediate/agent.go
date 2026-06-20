package remediate

// Agent performs an in-repo code edit task (e.g. fix a build break) in the
// working directory dir. Returning nil means the agent ran; the caller MUST
// re-verify the build before trusting the result.
type Agent interface {
	Repair(dir, task string) error
}

// FakeAgent is a test double.
type FakeAgent struct {
	Calls []string
	Err   error
	Edit  func(dir string) // optional: simulate file edits
}

func (f *FakeAgent) Repair(dir, task string) error {
	f.Calls = append(f.Calls, task)
	if f.Edit != nil {
		f.Edit(dir)
	}
	return f.Err
}
