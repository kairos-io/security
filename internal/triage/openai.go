package triage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/kairos-io/security/internal/config"
	"github.com/kairos-io/security/internal/state"
)

// OpenAIClient talks to an OpenAI-compatible chat-completions endpoint
// (LocalAI). It replaces the earlier nib-based client: nib is an interactive
// agent TUI, not a prompt->JSON client, so triage calls the model directly.
type OpenAIClient struct {
	endpoint    string // base URL, e.g. http://localhost:8080
	model       string
	maxTokens   int
	temperature float64
	httpc       *http.Client
}

func NewOpenAIClient(cfg config.AIConfig) *OpenAIClient {
	maxTokens := cfg.Nib.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 4096
	}
	return &OpenAIClient{
		endpoint:    strings.TrimRight(cfg.Nib.Endpoint, "/"),
		model:       cfg.Nib.Model,
		maxTokens:   maxTokens,
		temperature: cfg.Nib.Temperature,
		httpc:       &http.Client{Timeout: 5 * time.Minute},
	}
}

// aiResponse is the JSON contract we instruct the model to emit.
type aiResponse struct {
	Focus     []string          `json:"focus"`
	Summaries map[string]string `json:"summaries"`
	Narrative string            `json:"narrative"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// briefFinding is a compact, token-light view of a finding for the prompt, so
// a small model is not overwhelmed by the full finding records.
type briefFinding struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	CVE      string `json:"cve,omitempty"`
	Repo     string `json:"repo"`
	Package  string `json:"package,omitempty"`
}

// maxPromptFindings caps how many findings (highest severity first) are sent to
// the model, keeping the prompt within a small model's context window.
const maxPromptFindings = 50

func (c *OpenAIClient) Summarize(cor state.Correlated) ([]string, map[string]string, string, error) {
	if c.endpoint == "" {
		return nil, nil, "", fmt.Errorf("no AI endpoint configured")
	}

	findings := append([]state.Finding(nil), cor.Findings...)
	sort.SliceStable(findings, func(i, j int) bool {
		return sevRank[findings[i].Severity] > sevRank[findings[j].Severity]
	})
	if len(findings) > maxPromptFindings {
		findings = findings[:maxPromptFindings]
	}
	brief := struct {
		Findings  []briefFinding         `json:"findings"`
		Waterfall []state.WaterfallGroup `json:"waterfall"`
	}{Waterfall: cor.Waterfall}
	for _, f := range findings {
		brief.Findings = append(brief.Findings, briefFinding{
			ID: f.ID, Severity: f.Severity, CVE: f.CVEID, Repo: f.Repo, Package: f.Package,
		})
	}
	payload, err := json.Marshal(brief)
	if err != nil {
		return nil, nil, "", err
	}

	prompt := "You are a security triage assistant for the Kairos project. " +
		"Given this JSON of correlated security findings, return ONLY a single JSON object " +
		"(no prose, no markdown fences) with exactly these keys: " +
		`"focus" (array of the most urgent finding/waterfall "id" values, most urgent first, at most 20), ` +
		`"summaries" (object mapping each focus id to a one-line human summary), ` +
		`"narrative" (2-3 sentence overview of what to focus on). ` +
		"Only use id values that appear in the input. Findings:\n" + string(payload)

	reqBody, err := json.Marshal(chatRequest{
		Model:       c.model,
		Messages:    []chatMessage{{Role: "user", Content: prompt}},
		Temperature: c.temperature,
		MaxTokens:   c.maxTokens,
	})
	if err != nil {
		return nil, nil, "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.endpoint+"/v1/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpc.Do(req)
	if err != nil {
		return nil, nil, "", fmt.Errorf("call %s: %w", c.endpoint, err)
	}
	defer httpResp.Body.Close()
	body, _ := io.ReadAll(httpResp.Body)
	if httpResp.StatusCode != http.StatusOK {
		return nil, nil, "", fmt.Errorf("chat endpoint returned HTTP %d: %s", httpResp.StatusCode, snippet(body))
	}

	var cr chatResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return nil, nil, "", fmt.Errorf("decode chat response: %w (raw: %q)", err, snippet(body))
	}
	if cr.Error != nil {
		return nil, nil, "", fmt.Errorf("chat endpoint error: %s", cr.Error.Message)
	}
	if len(cr.Choices) == 0 {
		return nil, nil, "", fmt.Errorf("chat response had no choices (raw: %q)", snippet(body))
	}

	content := cr.Choices[0].Message.Content
	var resp aiResponse
	if err := json.Unmarshal([]byte(extractJSON(content)), &resp); err != nil {
		return nil, nil, "", fmt.Errorf("model did not return valid JSON: %w (content: %q)", err, snippet([]byte(content)))
	}
	return resp.Focus, resp.Summaries, resp.Narrative, nil
}

// extractJSON pulls the JSON object out of a model reply that may be wrapped in
// markdown fences or surrounded by prose.
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.Index(s, "```"); i >= 0 {
		rest := s[i+3:]
		rest = strings.TrimPrefix(rest, "json")
		rest = strings.TrimPrefix(rest, "JSON")
		if j := strings.Index(rest, "```"); j >= 0 {
			s = strings.TrimSpace(rest[:j])
		} else {
			s = strings.TrimSpace(rest)
		}
	}
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}

func snippet(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 300 {
		return s[:300] + "…"
	}
	return s
}
