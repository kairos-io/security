package review

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

const reviewMarker = "<!-- ksec:review -->"

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
			diff, derr := gh.PRDiff(repo.Repo, pr.Number)
			if derr != nil {
				errs = append(errs, state.CollectionError{Repo: repo.Repo, Collector: "review", Message: derr.Error()})
				continue
			}
			verdict, reasoning, _ := a.Assess(diff, pr) // assessor never hard-errors (defaults needs_human)
			rv := state.PRReview{Repo: repo.Repo, PR: pr.Number, URL: pr.URL, HeadSHA: pr.HeadSHA,
				Verdict: verdict, Reasoning: reasoning, ReviewedRun: runID}
			out = append(out, rv)
			body := comment(rv, cfg.Notify)
			if dryRun {
				fmt.Printf("[dry-run] would comment on %s#%d: %s\n", repo.Repo, pr.Number, verdict)
				continue
			}
			_ = gh.PostPRComment(repo.Repo, pr.Number, body)
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
	if len(notify) > 0 {
		fmt.Fprintf(&b, "\n\ncc %s", strings.Join(notify, " "))
	}
	b.WriteString("\n\n" + reviewMarker)
	return b.String()
}
