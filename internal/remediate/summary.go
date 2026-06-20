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
	"github.com/kairos-io/security/internal/state"
)

// SummarizeLedger asks the LocalAI OpenAI-compatible chat endpoint for a short
// cross-repo coordination narrative over the committed ledger. It is
// best-effort: any error lets the caller fall back to rendering without the
// coordination section. The endpoint, model and temperature come from the nib
// AI config.
func SummarizeLedger(cfg config.AIConfig, led state.Ledger) (string, error) {
	endpoint := strings.TrimRight(cfg.Nib.Endpoint, "/")
	if endpoint == "" {
		return "", fmt.Errorf("no AI endpoint configured")
	}

	prompt := "You are a security remediation coordinator. Given this status of automated " +
		"security PRs across repositories, write 2-4 plain sentences summarizing what is open, " +
		"what is cascading, what is blocked or needs a human, and what to prioritize. " +
		"Output only the sentences. Status:\n" + ledgerView(led)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":       cfg.Nib.Model,
		"temperature": cfg.Nib.Temperature,
		"max_tokens":  256,
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/v1/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	httpc := &http.Client{Timeout: 2 * time.Minute}
	resp, err := httpc.Do(req)
	if err != nil {
		return "", fmt.Errorf("call %s: %w", endpoint, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("summary HTTP %d", resp.StatusCode)
	}
	var cr struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &cr); err != nil {
		return "", fmt.Errorf("decode summary response: %w", err)
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("summary: no choices")
	}
	return strings.TrimSpace(cr.Choices[0].Message.Content), nil
}

// ledgerView renders the ledger as one compact line per entry for the prompt.
func ledgerView(led state.Ledger) string {
	var b strings.Builder
	for _, e := range led.Entries {
		fmt.Fprintf(&b, "%s %s@%s kind=%s source=%s state=%s",
			e.Repo, e.Package, e.Bump.To, e.Kind, e.Source, e.State)
		if e.Pseudo {
			b.WriteString(" (pseudo)")
		}
		if e.NeedsHuman {
			b.WriteString(" NEEDS-HUMAN")
		}
		if e.Blocked != "" {
			fmt.Fprintf(&b, " blocked=%s", e.Blocked)
		}
		b.WriteString("\n")
	}
	return b.String()
}
