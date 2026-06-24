package review

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

const (
	reviewMarker = "<!-- ksec:review -->"
	maxBumpDiff  = 40000 // per-bump upstream source diff cap
	maxContext   = 60000 // total assembled context cap before the PR diff
)

func Run(repos []state.Repo, gh ghclient.GitHub, a Assessor, cfg config.ReviewCfg, prev []state.PRReview, runID string, dryRun bool) ([]state.PRReview, []state.CollectionError) {
	prior := map[string]state.PRReview{}
	for _, r := range prev {
		prior[key(r.Repo, r.PR)] = r
	}
	var out []state.PRReview
	var errs []state.CollectionError
	assessed := 0
	for _, repo := range repos {
		prs, err := gh.ListOpenPRs(repo.Repo)
		if err != nil {
			errs = append(errs, state.CollectionError{Repo: repo.Repo, Collector: "review", Message: err.Error()})
			continue
		}
		for _, pr := range prs {
			if !pr.IsBot {
				continue
			}
			k := key(repo.Repo, pr.Number)
			// Idempotent: unchanged head -> carry the prior review forward.
			if p, ok := prior[k]; ok && p.HeadSHA == pr.HeadSHA {
				out = append(out, p)
				continue
			}
			if assessed >= cfg.MaxPerRun {
				// Over budget this run: keep any prior review so the dashboard
				// still shows it; otherwise skip until a future run.
				if p, ok := prior[k]; ok {
					out = append(out, p)
				}
				continue
			}
			assessed++
			// Assemble the assessment context: changelog (PR body) + upstream source
			// diffs for each bump + the PR's own diff.
			var ctx strings.Builder
			if strings.TrimSpace(pr.Body) != "" {
				ctx.WriteString("PR description / changelog:\n" + pr.Body + "\n\n")
			}
			diff, derr := gh.PRDiff(repo.Repo, pr.Number)
			if derr != nil {
				errs = append(errs, state.CollectionError{Repo: repo.Repo, Collector: "review", Message: derr.Error()})
				continue
			}
			var trace []string
			bumps := parseBumps(diff)
			if len(bumps) == 0 {
				trace = append(trace, "no go.mod dependency bumps parsed from the PR diff")
			}
			for _, b := range bumps {
				if ctx.Len() > maxContext {
					break
				}
				gr, ok := moduleRepo(b.Module)
				if !ok {
					trace = append(trace, fmt.Sprintf("%s %s→%s: module not resolvable to a GitHub repo (skipped)", b.Module, b.From, b.To))
					continue
				}
				baseRef, headRef := compareRef(b.From), compareRef(b.To)
				ud, uerr := gh.CompareDiff(gr, baseRef, headRef)
				if uerr != nil || len(ud) == 0 {
					trace = append(trace, fmt.Sprintf("%s %s→%s: compare %s...%s failed: %v (no upstream diff)", b.Module, b.From, b.To, baseRef, headRef, uerr))
					continue // degrade: no upstream source diff for this bump
				}
				if len(ud) > maxBumpDiff {
					ud = ud[:maxBumpDiff]
				}
				fmt.Fprintf(&ctx, "Upstream %s %s..%s:\n%s\n\n", b.Module, b.From, b.To, ud)
				trace = append(trace, fmt.Sprintf("%s %s→%s: compare %s...%s ✓ %d bytes", b.Module, b.From, b.To, baseRef, headRef, len(ud)))
			}
			ctx.WriteString("PR diff:\n" + string(diff))
			trace = append(trace, fmt.Sprintf("context: %d bytes", ctx.Len()))
			verdict, reasoning, summary, _ := a.Assess(pr, ctx.String()) // assessor never hard-errors (defaults needs_human)
			rv := state.PRReview{Repo: repo.Repo, PR: pr.Number, URL: pr.URL, HeadSHA: pr.HeadSHA,
				Verdict: verdict, Reasoning: reasoning, ChangesSummary: summary, Trace: trace, ReviewedRun: runID}
			out = append(out, rv)
			if dryRun {
				fmt.Printf("[dry-run] would comment on %s#%d: %s — %s\n", repo.Repo, pr.Number, verdict, summary)
				continue
			}
			_ = gh.UpsertPRComment(repo.Repo, pr.Number, reviewMarker, comment(rv, cfg.Notify))
			if cfg.AutoApprove && verdict == "good" {
				_ = gh.ApprovePR(repo.Repo, pr.Number, "kairos-security: automated review verdict good")
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Repo != out[j].Repo {
			return out[i].Repo < out[j].Repo
		}
		return out[i].PR < out[j].PR
	})
	return out, errs
}

func key(repo string, pr int) string { return fmt.Sprintf("%s#%d", repo, pr) }

func comment(r state.PRReview, notify []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "🔎 kairos-security review: **%s** — %s", r.Verdict, r.Reasoning)
	if r.ChangesSummary != "" {
		b.WriteString("\n\n**Dependency changes:** " + r.ChangesSummary)
	}
	if len(notify) > 0 {
		fmt.Fprintf(&b, "\n\ncc %s", strings.Join(notify, " "))
	}
	if len(r.Trace) > 0 {
		b.WriteString("\n\n<details><summary>review trace</summary>\n\n```\n")
		for _, line := range r.Trace {
			b.WriteString(line + "\n")
		}
		b.WriteString("```\n\n</details>")
	}
	b.WriteString("\n\n" + reviewMarker)
	return b.String()
}
