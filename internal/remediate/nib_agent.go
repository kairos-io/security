package remediate

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kairos-io/security/internal/config"
)

// NibAgent runs `nib --cli --yolo` in a repo clone to perform a code-edit task.
//
// nib reads its model from MODEL/BASE_URL/API_KEY in the environment (see nib's
// config.Load); NibAgent sets those from the AI config so repair works without
// relying on workflow-level env. nib's --cli is a line-at-a-time REPL that reads
// each prompt from a stdin line, so the task is collapsed to a single line with
// a trailing newline (see nibPrompt) — a multi-line prompt would otherwise be
// split into separate turns and a newline-less final line would be dropped.
//
// In --yolo mode every tool call is auto-approved without reading stdin, and nib
// exits non-zero on the stdin EOF after the turn; both are expected. A
// misconfigured or missing nib makes Repair return an error, which the caller
// treats as "no repair" (best-effort) and re-verifies the build regardless.
type NibAgent struct {
	cfg config.AIConfig
	run func(dir, task string) error
}

func NewNibAgent(cfg config.AIConfig) *NibAgent {
	return &NibAgent{cfg: cfg, run: func(dir, task string) error {
		cmd := exec.Command("nib", "--cli", "--yolo")
		cmd.Dir = dir
		cmd.Stdin = bytes.NewBufferString(nibPrompt(task))
		// nib reads MODEL/BASE_URL/API_KEY from env; appended last so they win
		// over anything inherited. BASE_URL must carry the /v1 suffix.
		cmd.Env = append(os.Environ(),
			"MODEL="+cfg.Nib.Model,
			"BASE_URL="+nibBaseURL(cfg.Nib.Endpoint),
			"API_KEY=sk-localai",
		)
		var errb bytes.Buffer
		cmd.Stderr = &errb
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("nib: %v: %s", err, errb.String())
		}
		return nil
	}}
}

// nibPrompt adapts a (possibly multi-line) task to nib's line-at-a-time CLI
// REPL: it collapses all whitespace to single spaces and terminates with a
// newline, so the whole task arrives as exactly one prompt.
func nibPrompt(task string) string {
	return strings.Join(strings.Fields(task), " ") + "\n"
}

// nibBaseURL normalizes an endpoint into nib's BASE_URL, which must include the
// /v1 suffix (e.g. http://localhost:8080 -> http://localhost:8080/v1).
func nibBaseURL(endpoint string) string {
	base := strings.TrimRight(endpoint, "/")
	if base == "" || strings.HasSuffix(base, "/v1") {
		return base
	}
	return base + "/v1"
}

func (a *NibAgent) Repair(dir, task string) error { return a.run(dir, task) }
