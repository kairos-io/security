package triage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/state"
)

type NibClient struct {
	cfg config.AIConfig
	run func(prompt string) ([]byte, error)
}

func NewNibClient(cfg config.AIConfig) *NibClient {
	return &NibClient{cfg: cfg, run: func(prompt string) ([]byte, error) {
		cmd := exec.Command("nib", "--"+cfg.Nib.Mode,
			"--model", cfg.Nib.Model, "--endpoint", cfg.Nib.Endpoint)
		cmd.Stdin = bytes.NewBufferString(prompt)
		var out, errb bytes.Buffer
		cmd.Stdout, cmd.Stderr = &out, &errb
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("nib: %v: %s", err, errb.String())
		}
		return out.Bytes(), nil
	}}
}

// aiResponse is the JSON contract we instruct the model to emit.
type aiResponse struct {
	Focus     []string          `json:"focus"`
	Summaries map[string]string `json:"summaries"`
	Narrative string            `json:"narrative"`
}

func (n *NibClient) Summarize(c state.Correlated) ([]string, map[string]string, string, error) {
	payload, err := json.Marshal(c)
	if err != nil {
		return nil, nil, "", err
	}
	prompt := "You are a security triage assistant. Given this JSON of correlated security findings, " +
		"return ONLY a JSON object with keys: focus (array of finding/waterfall IDs ordered most-urgent first), " +
		"summaries (map of id to one-line human summary for high/critical items), and narrative (2-3 sentence " +
		"'what to focus on' overview). Do not invent IDs. Findings:\n" + string(payload)

	raw, err := n.run(prompt)
	if err != nil {
		return nil, nil, "", err
	}
	var resp aiResponse
	if err := json.Unmarshal(bytes.TrimSpace(raw), &resp); err != nil {
		snippet := string(bytes.TrimSpace(raw))
		if len(snippet) > 300 {
			snippet = snippet[:300] + "…"
		}
		return nil, nil, "", fmt.Errorf("parse model output: %w (raw: %q)", err, snippet)
	}
	return resp.Focus, resp.Summaries, resp.Narrative, nil
}
