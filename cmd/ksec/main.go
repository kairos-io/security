package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kairos-io/security/internal/collect"
	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/discover"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
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

	root.AddCommand(newDiscoverCmd(gf))
	root.AddCommand(newCollectCmd(gf))
	root.AddCommand(newStubCmd("correlate"))
	root.AddCommand(newStubCmd("triage"))
	root.AddCommand(newStubCmd("render"))
	return root
}

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

func newStubCmd(name string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "run the " + name + " phase",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%s: not implemented\n", name)
			return nil
		},
	}
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
