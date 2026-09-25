package remediate

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// shimEnv puts stub `git`, `go` and `gh` commands ahead of the real ones on
// PATH so a whole GitExecutor entry point can run end to end with no network.
// `git` is the real binary with any https URL rewritten to a local bare repo,
// so clone, commit and push all really happen and the pushed tree can be read
// back. `go` and `gh` are simulations.
type shimEnv struct {
	realGit string
	origin  string
	root    string
}

func newShimEnv(t *testing.T) *shimEnv {
	t.Helper()
	realGit, err := exec.LookPath("git")
	require.NoError(t, err)

	root := t.TempDir()
	sh := &shimEnv{realGit: realGit, origin: filepath.Join(root, "origin.git"), root: root}

	sh.git(t, "", "init", "--bare", "-b", "main", sh.origin)

	seed := filepath.Join(root, "seed")
	require.NoError(t, os.MkdirAll(seed, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(seed, "go.mod"), []byte("module m\n\ngo 1.22\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(seed, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o644))
	sh.git(t, seed, "init", "-b", "main")
	sh.git(t, seed, "config", "user.email", "seed@example.com")
	sh.git(t, seed, "config", "user.name", "seed")
	sh.git(t, seed, "add", "-A")
	sh.git(t, seed, "commit", "-m", "init")
	sh.git(t, seed, "remote", "add", "origin", sh.origin)
	sh.git(t, seed, "push", "-q", "origin", "main")

	bin := filepath.Join(root, "bin")
	require.NoError(t, os.MkdirAll(bin, 0o755))
	write := func(name, body string) {
		require.NoError(t, os.WriteFile(filepath.Join(bin, name), []byte(body), 0o755))
	}
	write("git", "#!/bin/bash\n"+
		"a=()\n"+
		"for x in \"$@\"; do if [[ \"$x\" == https://* ]]; then x=\"$KSEC_SHIM_ORIGIN\"; fi; a+=(\"$x\"); done\n"+
		"exec "+realGit+" \"${a[@]}\"\n")
	// `go get` and `go mod edit` must leave a tracked change behind, or the
	// no-op guards in Adjust/Repin/Toolchain short-circuit before committing.
	write("go", "#!/bin/bash\n"+
		"if [ \"$1 $2\" = 'list -m' ]; then echo 'example.com/x v1.0.0 v1.2.3'; exit 0; fi\n"+
		"case \"$1\" in\n"+
		"  get) printf '\\nrequire example.com/x v1.2.3\\n' >> go.mod ;;\n"+
		"  mod) if [ \"$2\" = edit ]; then printf '\\n// toolchain bumped\\n' >> go.mod; fi ;;\n"+
		"  build) if [ -f \"$KSEC_SHIM_BUILDFAIL\" ]; then rm -f \"$KSEC_SHIM_BUILDFAIL\";"+
		" echo 'undefined: Foo' >&2; exit 1; fi ;;\n"+
		"esac\nexit 0\n")
	write("gh", "#!/bin/bash\n"+
		"if [ \"$1 $2\" = 'pr create' ]; then echo 'https://github.com/o/r/pull/7'; fi\nexit 0\n")

	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("KSEC_SHIM_ORIGIN", sh.origin)
	t.Setenv("KSEC_SHIM_BUILDFAIL", filepath.Join(root, "buildfail"))
	return sh
}

func (s *shimEnv) git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command(s.realGit, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, out)
	return string(out)
}

// failNextBuild makes the next `go build ./...` fail once, which is what drives
// verifyOrRepair into the repair agent.
func (s *shimEnv) failNextBuild(t *testing.T) {
	t.Helper()
	require.NoError(t, os.WriteFile(os.Getenv("KSEC_SHIM_BUILDFAIL"), nil, 0o644))
}

// seedBranch publishes a ksec branch on origin, which Adjust and Repin require
// because they check out an existing branch rather than creating one.
func (s *shimEnv) seedBranch(t *testing.T, branch string) {
	t.Helper()
	seed := filepath.Join(s.root, "seed")
	s.git(t, seed, "checkout", "-q", "-b", branch)
	s.git(t, seed, "push", "-q", "origin", branch)
	s.git(t, seed, "checkout", "-q", "main")
}

// fileInBranch returns the contents of path on branch in the bare origin.
func (s *shimEnv) fileInBranch(t *testing.T, branch, path string) (string, error) {
	t.Helper()
	out, err := exec.Command(s.realGit, "-C", s.origin, "show", branch+":"+path).CombinedOutput()
	return string(out), err
}

// agentCreating returns an agent that writes name into the clone, the way the
// nib agent adds a compatibility shim when it adapts code to a changed API.
func agentCreating(name string) *FakeAgent {
	return &FakeAgent{Edit: func(dir string) {
		_ = os.WriteFile(filepath.Join(dir, name), []byte("package main\n\n// repaired\n"), 0o644)
	}}
}

// Every path that commits must push the tree it verified. `git commit -am`
// stages tracked files only, so a file the repair agent created is dropped from
// the commit and then destroyed with the temp clone: the branch on GitHub holds
// the dependency bump without the change that made it compile.
func TestCommitPathsPushTheTreeTheyVerified(t *testing.T) {
	bump := state.Bump{Package: "example.com/x", To: "v1.2.3"}

	cases := []struct {
		name string
		// branch is where the pushed tree is read back from; wantInMod is the
		// text the path's own edit leaves in go.mod, so a subtest cannot pass
		// on a tree where nothing but the agent's file changed.
		branch    string
		wantInMod string
		seed      bool
		call      func(g *GitExecutor) error
	}{
		{
			name:      "Open",
			branch:    "ksec/bump-example-com-x-v1-2-3",
			wantInMod: "example.com/x",
			call: func(g *GitExecutor) error {
				_, err := g.Open(Intent{Key: "k", Repo: "o/r", Package: "example.com/x", Bump: bump}, "run")
				return err
			},
		},
		{
			name:      "Cascade",
			branch:    "ksec/cascade-example-com-x-pseudo",
			wantInMod: "example.com/x",
			call: func(g *GitExecutor) error {
				_, err := g.Cascade(Intent{Key: "k", Repo: "o/r", Package: "example.com/x", Ref: "v0.0.0-pseudo"}, "run")
				return err
			},
		},
		{
			name:      "Supersede",
			branch:    "ksec/bump-example-com-x-v1-2-3",
			wantInMod: "example.com/x",
			call: func(g *GitExecutor) error {
				_, err := g.Supersede(Intent{Key: "k", Repo: "o/r", Package: "example.com/x", Bump: bump}, "run")
				return err
			},
		},
		{
			name:      "Toolchain",
			branch:    "ksec/toolchain-1-24-0",
			wantInMod: "toolchain bumped",
			call: func(g *GitExecutor) error {
				_, err := g.Toolchain(Intent{Key: "k", Repo: "o/r", ToolchainVersion: "1.24.0"}, "run")
				return err
			},
		},
		{
			name:      "Adjust",
			branch:    "ksec/bump-example-com-x-v1-0-0",
			wantInMod: "example.com/x",
			seed:      true,
			call: func(g *GitExecutor) error {
				_, err := g.Adjust(state.LedgerEntry{
					Repo: "o/r", Package: "example.com/x", PRNumber: 1,
					Branch: "ksec/bump-example-com-x-v1-0-0",
				}, "v1.2.3", "run")
				return err
			},
		},
		{
			name:      "Repin",
			branch:    "ksec/cascade-example-com-x-pseudo",
			wantInMod: "example.com/x",
			seed:      true,
			call: func(g *GitExecutor) error {
				_, err := g.Repin(state.LedgerEntry{
					State: "open", Repo: "o/r", Package: "example.com/x",
					Branch: "ksec/cascade-example-com-x-pseudo",
				}, "run")
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sh := newShimEnv(t)
			if tc.seed {
				sh.seedBranch(t, tc.branch)
			}
			sh.failNextBuild(t)

			g := &GitExecutor{Agent: agentCreating("compat.go")}
			require.NoError(t, tc.call(g))

			// The bump itself always lands; that is not what regresses.
			modAtOrigin, err := sh.fileInBranch(t, tc.branch, "go.mod")
			require.NoError(t, err, modAtOrigin)
			assert.Contains(t, modAtOrigin, tc.wantInMod)

			got, err := sh.fileInBranch(t, tc.branch, "compat.go")
			require.NoError(t, err,
				"the repair agent's new file never reached the pushed branch: %s", got)
			assert.Contains(t, got, "repaired")
		})
	}
}

// ResolveConflict already staged everything before continuing a rebase; keep
// that covered so the shared helper cannot regress it.
func TestCommitAllStagesNewAndModifiedAndDeleted(t *testing.T) {
	sh := newShimEnv(t)
	dir := filepath.Join(sh.root, "work")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	sh.git(t, "", "clone", "-q", sh.origin, dir)
	sh.git(t, dir, "config", "user.email", "bot@kairos.io")
	sh.git(t, dir, "config", "user.name", "bot")

	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module m\n\ngo 1.23\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "compat.go"), []byte("package main\n"), 0o644))
	require.NoError(t, os.Remove(filepath.Join(dir, "main.go")))

	g := &GitExecutor{}
	require.NoError(t, g.commitAll(dir, "chore: everything"))

	out := sh.git(t, dir, "show", "--name-status", "--format=", "HEAD")
	assert.Contains(t, out, "A\tcompat.go")
	assert.Contains(t, out, "M\tgo.mod")
	assert.Contains(t, out, "D\tmain.go")
	assert.Empty(t, sh.git(t, dir, "status", "--porcelain"))
}
