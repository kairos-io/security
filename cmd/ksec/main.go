package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/kairos-io/security/internal/collect"
	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/correlate"
	"github.com/kairos-io/security/internal/discover"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/remediate"
	"github.com/kairos-io/security/internal/render"
	"github.com/kairos-io/security/internal/review"
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
	root.AddCommand(newReviewCmd(gf))
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
				collect.GHAlerts{GH: gh},
				collect.ImageCVE{Runner: trivyRunner},
				collect.SourceCVE{Runner: govulncheckRunner},
			}
			out := collect.Run(repos, collectors, prev)
			if err := state.Save(gf.stateDir, state.FindingsFile, out); err != nil {
				return err
			}
			prs, prErrs := collect.OpenPRs(repos, gh, out.Findings)
			out.Errors = append(out.Errors, prErrs...)
			if len(prErrs) > 0 {
				_ = state.Save(gf.stateDir, state.FindingsFile, out) // include PR-list errors
			}
			fmt.Fprintf(os.Stderr, "collect: %d repos · %d findings · %d errors · %d PRs tied to CVEs\n",
				len(repos), len(out.Findings), len(out.Errors), len(prs))
			return state.Save(gf.stateDir, state.OpenPRsFile, prs)
		},
	}
}

func trivyRunner(ref string) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "image-scan: trivy %s\n", ref)
	out, err := exec.Command("trivy", "image", "--quiet", "--scanners", "vuln", "--format", "json", ref).Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "image-scan: %s → error: %v\n", ref, err)
	} else {
		fmt.Fprintf(os.Stderr, "image-scan: %s → ok\n", ref)
	}
	return out, err
}

// govulncheckRunner shallow-clones the repo to a temp dir and runs govulncheck.
// Non-Go repos (no root go.mod) are skipped without error so they neither
// produce a finding nor a collection error.
func govulncheckRunner(r state.Repo) ([]byte, error) {
	if !r.SourceScanEnabled() {
		return nil, nil // explicitly opted out of source scanning
	}

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
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, runErr := cmd.Output()
	return collect.ClassifyGovulncheck(out, stderr.Bytes(), runErr)
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
			fmt.Fprintf(os.Stderr, "triage: %d findings, focus=%d\n", len(c.Findings), len(out.Focus))
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
	var seeds []string
	cmd := &cobra.Command{
		Use:   "remediate",
		Short: "open and maintain dependency-bump PRs for actionable findings",
		RunE: func(cmd *cobra.Command, args []string) error {
			var c state.Correlated
			if err := state.Load(gf.stateDir, state.CorrelatedFile, &c); err != nil {
				return err
			}
			for _, s := range seeds {
				f, err := remediate.ParseSeed(s)
				if err != nil {
					return err
				}
				c.Findings = append(c.Findings, f)
			}
			if len(seeds) > 0 {
				fmt.Fprintf(os.Stderr, "remediate: injected %d seed finding(s) for testing\n", len(seeds))
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
			// Fork external repos (kind: external) the bot can't push to
			// directly; org repos push direct. ForkOwner is the PAT's own
			// login, with KSEC_FORK_OWNER as override/fallback.
			ex.ShouldFork = remediate.ForkByKind(repos)
			forkOwner := os.Getenv("KSEC_FORK_OWNER")
			if out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output(); err == nil {
				if login := strings.TrimSpace(string(out)); login != "" {
					forkOwner = login
				}
			}
			ex.ForkOwner = forkOwner
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
			fmt.Fprintf(os.Stderr, "remediate: %d intents → %s\n", len(intents), actionCounts(results))
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
	cmd.Flags().StringArrayVar(&seeds, "seed", nil,
		"inject a synthetic finding to test remediation: owner/repo=package@version (repeatable)")
	return cmd
}

// actionCounts tallies results by their Action into a stable-ordered summary
// string like "open=2 adopt=1". Ordering is fixed (not map iteration order)
// so the log line is deterministic across runs.
func actionCounts(results []remediate.Result) string {
	counts := map[string]int{}
	for _, r := range results {
		counts[r.Action]++
	}
	order := []string{"open", "adopt", "supersede", "cascade", "toolchain", "repin", "reconcile", "needs-human"}
	seen := map[string]bool{}
	var parts []string
	for _, a := range order {
		if n := counts[a]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s=%d", a, n))
			seen[a] = true
		}
	}
	// Append any actions not in the known order, sorted for stability.
	var extra []string
	for a := range counts {
		if !seen[a] {
			extra = append(extra, a)
		}
	}
	sort.Strings(extra)
	for _, a := range extra {
		parts = append(parts, fmt.Sprintf("%s=%d", a, counts[a]))
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, " ")
}

func newReviewCmd(gf *globalFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "review",
		Short: "AI-assess open bot PRs and post a verdict",
		RunE: func(cmd *cobra.Command, args []string) error {
			aiCfg, err := config.LoadAI("ai.yaml")
			if err != nil {
				return err
			}
			if !aiCfg.Review.Enabled || aiCfg.Nib.Endpoint == "" {
				fmt.Fprintln(os.Stderr, "review: disabled or no AI endpoint — skipping")
				return nil
			}
			var repos []state.Repo
			if err := state.Load(gf.stateDir, state.ReposFile, &repos); err != nil {
				return err
			}
			var prev []state.PRReview
			_ = state.Load(gf.stateDir, state.ReviewsFile, &prev) // best-effort
			gh := ghclient.NewCLI()
			reviews, errs := review.Run(repos, gh, review.NewOpenAIAssessor(aiCfg), aiCfg.Review, prev, collect.Today(), gf.dryRun)
			counts := map[string]int{}
			for _, r := range reviews {
				counts[r.Verdict]++
			}
			fmt.Fprintf(os.Stderr, "review: %d reviews (good=%d bad=%d needs-human=%d) · %d errors\n",
				len(reviews), counts["good"], counts["bad"], counts["needs_human_verification"], len(errs))
			return state.Save(gf.stateDir, state.ReviewsFile, reviews)
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
			var repos []state.Repo
			_ = state.Load(gf.stateDir, state.ReposFile, &repos) // best-effort: show all tracked repos
			var ledger state.Ledger
			_ = state.Load(gf.stateDir, state.LedgerFile, &ledger) // best-effort
			var openPRs []state.TrackedPR
			_ = state.Load(gf.stateDir, state.OpenPRsFile, &openPRs) // best-effort

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
				OpenPRs:             openPRs,
			}
			// Committed artifacts must be deterministic across runs of the
			// same data, so render them with the volatile RunURL stripped.
			// This lets the workflow's no-op commit guard succeed when only
			// the run id changed.
			fileIn := in
			fileIn.RunURL = ""
			fileIn.CoordinationSummary = ""
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
			issueNum, err := render.UpsertTrackingIssue(ghclient.NewCLI(), trackingRepo, issueBody, gf.dryRun)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: tracking issue upsert failed: %v\n", err)
				fmt.Fprintf(os.Stderr, "render: dashboard + site + issue (upsert failed)\n")
			} else if issueNum > 0 {
				fmt.Fprintf(os.Stderr, "render: dashboard + site + issue #%d\n", issueNum)
			} else {
				fmt.Fprintf(os.Stderr, "render: dashboard + site + issue\n")
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
