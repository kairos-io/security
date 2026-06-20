package remediate

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/kairos-io/security/internal/config"
)

// NibAgent runs `nib --cli --yolo` in a repo clone to perform a code-edit task.
// nib reads its model/endpoint from the environment/config provided by the
// workflow; a misconfigured or missing nib makes Repair return an error, which
// the caller treats as "no repair" (best-effort).
type NibAgent struct {
	cfg config.AIConfig
	run func(dir, task string) error
}

func NewNibAgent(cfg config.AIConfig) *NibAgent {
	return &NibAgent{cfg: cfg, run: func(dir, task string) error {
		cmd := exec.Command("nib", "--cli", "--yolo")
		cmd.Dir = dir
		cmd.Stdin = bytes.NewBufferString(task)
		var errb bytes.Buffer
		cmd.Stderr = &errb
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("nib: %v: %s", err, errb.String())
		}
		return nil
	}}
}

func (a *NibAgent) Repair(dir, task string) error { return a.run(dir, task) }
