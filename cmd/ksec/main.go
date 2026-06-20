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
	"github.com/kairos-io/security/internal/remediate"
	"github.com/kairos-io/security/internal/render"
	"github.com/kairos-io/security/internal/state"
	"github.com/kairos-io/security/internal/triage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
	root.AddCommand(newRemediateCmd(gf))
	root.AddCommand(newRenderCmd(gf))
	return root
}

func newDiscoverCmd(gf *globalFlags) *cobra.Command {
	var seedFrom string
	cmd := &cobra.Command{
		Use:   "discover",
		Short: "build the tracked-repo list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if seedFrom != "" {
				repos, err := discover.SeedFromOrg(ghclient.NewCLI(), seedFrom, seedFrom+"/kairos-init", "main")
				if err != nil {
					return err
				}
				out, err := yaml.Marshal(struct {
					Repos []state.Repo `yaml:"repos"`
				}{repos})
				if err != nil {
					return err
				}
				_, err = cmd.OutOrStdout().Write(out)
				return err
			}

			cfg, err := config.LoadRepos("repos.yaml")
			if err != nil {
				return err
			}
			repos := discover.Normalize(cfg)
			return state.Save(gf.stateDir, state.ReposFile, repos)
		},
	}
	cmd.Flags().StringVar(&seedFrom, "seed-from", "",
		"enumerate the given GitHub org once and print a curated repos.yaml block to stdout (does not write state)")
	return cmd
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
// Non-Go repos (no root go.mod) are skipped without error so they neither
// produce a finding nor a collection error.
func govulncheckRunner(r state.Repo) ([]byte, error) {
	// Only scan repos that have a root go.mod; otherwise govulncheck just
	// fails ("exit status 1") on docs/helm/.github repos.
	if err := exec.Command("gh", "api", "repos/"+r.Repo+"/contents/go.mod").Run(); err != nil {
		return nil, nil // not a Go repo (or inaccessible): skip
	}

	dir, err := os.MkdirTemp("", "ksec-src-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	// Clone the default branch (no --branch): repos differ between main/master.
	// Authenticate with the token so private repos and rate limits work.
	url := "https://github.com/" + r.Repo + ".git"
	if token := os.Getenv("GH_TOKEN"); token != "" {
		url = "https://x-access-token:" + token + "@github.com/" + r.Repo + ".git"
	}
	clone := exec.Command("git", "clone", "--depth", "1", url, dir)
	if out, err := clone.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("clone: %v: %s", err, out)
	}

	cmd := exec.Command("govulncheck", "-json", "./...")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil && len(out) == 0 {
		// No JSON at all: a real failure. A non-zero exit with output is
		// normal — govulncheck exits non-zero when it finds vulnerabilities.
		return nil, err
	}
	return out, nil
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
	var requireAI bool
	cmd := &cobra.Command{
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
			fmt.Fprintf(os.Stderr, "triage: requesting AI summary — model=%q endpoint=%q\n",
				aiCfg.Nib.Model, aiCfg.Nib.Endpoint)
			out, aiErr := triage.Run(c, triage.NewOpenAIClient(aiCfg), aiCfg.LocalAI.Model.Name)
			if aiErr != nil {
				fmt.Fprintf(os.Stderr, "triage: ❌ AI UNAVAILABLE: %v\n", aiErr)
				if requireAI {
					return fmt.Errorf("AI triage required (--require-ai) but failed: %w", aiErr)
				}
				fmt.Fprintln(os.Stderr, "triage: falling back to deterministic prioritization (use --require-ai to fail instead)")
			} else {
				fmt.Fprintf(os.Stderr, "triage: ✅ AI summary OK — model=%q focus=%d narrative=%dB\n",
					out.Model, len(out.Focus), len(out.Narrative))
			}
			return state.Save(gf.stateDir, state.TriageFile, out)
		},
	}
	cmd.Flags().BoolVar(&requireAI, "require-ai", false,
		"fail the phase if the AI summary cannot be produced (instead of falling back)")
	return cmd
}

func newRemediateCmd(gf *globalFlags) *cobra.Command {
	var maxPRs int
	var aiProse bool
	var automerge bool
	var repair bool
	cmd := &cobra.Command{
		Use:   "remediate",
		Short: "open and maintain dependency-bump PRs for actionable findings",
		RunE: func(cmd *cobra.Command, args []string) error {
			var c state.Correlated
			if err := state.Load(gf.stateDir, state.CorrelatedFile, &c); err != nil {
				return err
			}
			var ledger state.Ledger
			_ = state.Load(gf.stateDir, state.LedgerFile, &ledger) // best-effort: empty on first run

			runID := os.Getenv("KSEC_RUN_URL")
			if runID == "" {
				runID = "local"
			}
			gh := ghclient.NewCLI()
			// Collect open PRs per tracked repo so the planner can adopt existing
			// dependabot/renovate/human PRs instead of duplicating them.
			prsByRepo := map[string][]ghclient.PullRequest{}
			var repos []state.Repo
			if err := state.Load(gf.stateDir, state.ReposFile, &repos); err == nil {
				for _, r := range repos {
					if prs, err := gh.ListOpenPRs(r.Repo); err == nil {
						prsByRepo[r.Repo] = prs
					}
				}
			}
			// Build the dependency graph (module<->repo, consumers) by
			// fetching each tracked repo's go.mod. Best-effort: a repo whose
			// go.mod can't be read is simply absent from the graph.
			gomodByRepo := map[string][]byte{}
			for _, r := range repos {
				if b, err := gh.GetFile(r.Repo, "go.mod", r.Branch); err == nil {
					gomodByRepo[r.Repo] = b
				}
			}
			graph := remediate.BuildGraph(repos, gomodByRepo)
			intents, deferred := remediate.Plan(c, ledger, prsByRepo, graph, maxPRs)
			if deferred > 0 {
				fmt.Fprintf(os.Stderr, "remediate: %d new bumps deferred by --max-prs=%d\n", deferred, maxPRs)
			}
			// Load AI config up front: it enables AI-drafted PR prose (used when
			// PRs are opened in remediate.Run below) and the comment reactions.
			aiCfg, _ := config.LoadAI("ai.yaml")
			ex := &remediate.GitExecutor{Token: os.Getenv("GH_TOKEN"), DryRun: gf.dryRun, GH: gh, Automerge: automerge}
			if aiProse && aiCfg.Nib.Endpoint != "" {
				ex.Prose = remediate.NewOpenAIProse(aiCfg)
			}
			// The nib agent repairs build breaks / conflicts on owned PRs. It is
			// only wired when --repair is set and an AI endpoint is configured;
			// it is used solely on build/verify failures, degrading otherwise to
			// the deterministic build-failed/needsHuman behavior.
			if repair && aiCfg.Nib.Endpoint != "" {
				ex.Agent = remediate.NewNibAgent(aiCfg)
			}
			out, results := remediate.Run(intents, ex, ledger, runID)
			for _, r := range results {
				fmt.Fprintf(os.Stderr, "remediate: %s %s -> %s %s\n", r.Action, r.Key, r.State, r.Detail)
			}
			// React to review comments on PRs we own. Only run the reaction
			// loop when an AI endpoint is configured: without one the classifier
			// errors on every comment, appending a needs-human event forever and
			// churning the ledger across runs.
			if aiCfg.Nib.Endpoint != "" {
				classifier := remediate.NewOpenAIClassifier(aiCfg)
				for i := range out.Entries {
					e := &out.Entries[i]
					if e.State != "open" || e.PRNumber == 0 {
						continue
					}
					title := remediate.PRTitle(remediate.Intent{Package: e.Package, Bump: e.Bump})
					if err := remediate.ReactToComments(e, gh, classifier, ex, title, runID, gf.dryRun); err != nil {
						fmt.Fprintf(os.Stderr, "remediate: react %s: %v\n", e.Key, err)
					}
				}
			} else {
				fmt.Fprintln(os.Stderr, "remediate: comment reactions disabled (no AI endpoint configured)")
			}
			if gf.dryRun {
				fmt.Fprintln(os.Stderr, "remediate: dry-run — ledger not persisted")
				return nil
			}
			return state.Save(gf.stateDir, state.LedgerFile, out)
		},
	}
	cmd.Flags().IntVar(&maxPRs, "max-prs", 10, "maximum NEW PRs to open per run (blast-radius guard)")
	cmd.Flags().BoolVar(&aiProse, "ai-pr-prose", true, "use the AI model to draft PR descriptions (falls back to deterministic text)")
	cmd.Flags().BoolVar(&automerge, "automerge", false, "merge addressing PRs (ours/dependabot/renovate) when green and unblocked")
	cmd.Flags().BoolVar(&repair, "repair", true, "use the nib agent to repair build breaks / conflicts")
	return cmd
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
			var repos []state.Repo
			_ = state.Load(gf.stateDir, state.ReposFile, &repos) // best-effort: show all tracked repos
			var ledger state.Ledger
			_ = state.Load(gf.stateDir, state.LedgerFile, &ledger) // best-effort

			// Best-effort cross-repo coordination narrative for the dashboard.
			// On any AI failure the summary stays empty and the section is omitted.
			summary := ""
			if aiCfg, err := config.LoadAI("ai.yaml"); err == nil {
				if s, err := remediate.SummarizeLedger(aiCfg, ledger); err == nil {
					summary = s
				}
			}

			in := render.Input{
				Correlated:          c,
				Triage:              tr,
				Repos:               repos,
				Ledger:              ledger,
				CollectErrors:       findings.Errors,
				RunURL:              os.Getenv("KSEC_RUN_URL"),
				CoordinationSummary: summary,
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
			// The dashboard files and site/index.html are already written, so a
			// flaky issue API call must not fail the pipeline; warn and move on.
			if _, err := render.UpsertTrackingIssue(ghclient.NewCLI(), trackingRepo, issueBody, gf.dryRun); err != nil {
				fmt.Fprintf(os.Stderr, "warning: tracking issue upsert failed: %v\n", err)
			}
			return nil
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
