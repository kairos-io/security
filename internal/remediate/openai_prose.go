package remediate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kairos-io/security/internal/config"
)

// OpenAIProse drafts a short, human PR-body paragraph via the LocalAI
// OpenAI-compatible chat endpoint. Unlike the classifier it does NOT force a
// tool call — the output is free-form prose, returned as plain text. It is
// best-effort: any error makes PRBodyWith fall back to the deterministic body.
type OpenAIProse struct {
	endpoint    string
	model       string
	temperature float64
	httpc       *http.Client
}

var _ ProseClient = (*OpenAIProse)(nil)

func NewOpenAIProse(cfg config.AIConfig) *OpenAIProse {
	return &OpenAIProse{
		endpoint:    strings.TrimRight(cfg.Nib.Endpoint, "/"),
		model:       cfg.Nib.Model,
		temperature: cfg.Nib.Temperature,
		httpc:       &http.Client{Timeout: 2 * time.Minute},
	}
}

func (p *OpenAIProse) DraftPRBody(in Intent) (string, error) {
	if p.endpoint == "" {
		return "", fmt.Errorf("no AI endpoint configured")
	}
	prompt := fmt.Sprintf("Write 1-2 concise, factual sentences for a pull-request description explaining why "+
		"bumping the Go dependency %q to version %s matters (severity: %s). "+
		"Output only the sentences — no preamble, no markdown headings, no code fences.",
		in.Package, in.Bump.To, in.Severity)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":       p.model,
		"temperature": p.temperature,
		"max_tokens":  200,
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint+"/v1/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.httpc.Do(req)
	if err != nil {
		return "", fmt.Errorf("call %s: %w", p.endpoint, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("prose HTTP %d", resp.StatusCode)
	}
	var cr struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &cr); err != nil {
		return "", fmt.Errorf("decode prose response: %w", err)
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("prose: no choices")
	}
	return strings.TrimSpace(cr.Choices[0].Message.Content), nil
}
