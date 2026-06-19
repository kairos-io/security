package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kairos-io/security/internal/collect"
	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/correlate"
	"github.com/kairos-io/security/internal/discover"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/render"
	"github.com/kairos-io/security/internal/state"
	"github.com/kairos-io/security/internal/triage"
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
	root.AddCommand(newCorrelateCmd(gf))
	root.AddCommand(newTriageCmd(gf))
	root.AddCommand(newRenderCmd(gf))
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

func newCorrelateCmd(gf *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "correlate",
		Short: "dedupe findings and build the waterfall graph",
		RunE: func(cmd *cobra.Command, args []string) error {
			var in state.Findings
			if err := state.Load(gf.stateDir, state.FindingsFile, &in); err != nil {
				return err
			}
			return state.Save(gf.stateDir, state.CorrelatedFile, correlate.Run(in))
		},
	}
}

func newTriageCmd(gf *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "triage",
		Short: "prioritize findings and write the AI summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			var c state.Correlated
			if err := state.Load(gf.stateDir, state.CorrelatedFile, &c); err != nil {
				return err
			}
			aiCfg, err := config.LoadAI("ai.yaml")
			if err != nil {
				return err
			}
			out := triage.Run(c, triage.NewNibClient(aiCfg), aiCfg.LocalAI.Model.Name)
			return state.Save(gf.stateDir, state.TriageFile, out)
		},
	}
}

func newRenderCmd(gf *globalFlags) *cobra.Command {
	var trackingRepo string
	cmd := &cobra.Command{
		Use:   "render",
		Short: "write dashboard files and upsert the tracking issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			var c state.Correlated
			if err := state.Load(gf.stateDir, state.CorrelatedFile, &c); err != nil {
				return err
			}
			var tr state.Triage
			if err := state.Load(gf.stateDir, state.TriageFile, &tr); err != nil {
				return err
			}
			var findings state.Findings
			_ = state.Load(gf.stateDir, state.FindingsFile, &findings)

			in := render.Input{
				Correlated:    c,
				Triage:        tr,
				CollectErrors: findings.Errors,
				RunURL:        os.Getenv("KSEC_RUN_URL"),
			}
			// Committed artifacts must be deterministic across runs of the
			// same data, so render them with the volatile RunURL stripped.
			// This lets the workflow's no-op commit guard succeed when only
			// the run id changed.
			fileIn := in
			fileIn.RunURL = ""
			md := render.DashboardMarkdown(fileIn)
			j, err := render.DashboardJSON(fileIn)
			if err != nil {
				return err
			}
			if err := os.WriteFile("dashboard.md", []byte(md), 0o644); err != nil {
				return err
			}
			if err := os.WriteFile("dashboard.json", j, 0o644); err != nil {
				return err
			}
			// The HTML dashboard is published to GitHub Pages, not committed,
			// so it can carry the per-run RunURL footer.
			if err := os.MkdirAll("site", 0o755); err != nil {
				return err
			}
			if err := os.WriteFile("site/index.html", []byte(render.DashboardHTML(in)), 0o644); err != nil {
				return err
			}
			// The tracking issue body is not committed, so keep the run-log
			// footer (RunURL) there for traceability.
			issueBody := render.DashboardMarkdown(in)
			_, err = render.UpsertTrackingIssue(ghclient.NewCLI(), trackingRepo, issueBody, gf.dryRun)
			return err
		},
	}
	cmd.Flags().StringVar(&trackingRepo, "tracking-repo", "kairos-io/kairos", "repo to upsert the tracking issue into")
	return cmd
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
