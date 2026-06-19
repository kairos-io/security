package main

import (
	"fmt"
	"os"

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

	for _, phase := range []string{"discover", "collect", "correlate", "triage", "render"} {
		p := phase
		root.AddCommand(&cobra.Command{
			Use:   p,
			Short: "run the " + p + " phase",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Printf("%s: not implemented\n", p)
				return nil
			},
		})
	}
	return root
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
