package main

import (
	"fmt"
	"os"

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
	root.AddCommand(newStubCmd("collect"))
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
