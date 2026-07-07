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
//
// To get reliable structured output we never ask the model for free-form JSON;
// instead we force a function/tool call whose parameters are a JSON schema, so
// the backend grammar-constrains the output to that shape.
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

const triageToolName = "report_triage"

// triageToolParameters is the JSON schema the model's tool-call arguments must
// satisfy. summaries is an array of {id, summary} (not an arbitrary-key object)
// because fixed-shape objects are far more reliable under grammar constraints.
const triageToolParameters = `{
  "type": "object",
  "properties": {
    "focus": {
      "type": "array",
      "items": {"type": "string"},
      "description": "Most urgent finding/waterfall ids, most urgent first, at most 20"
    },
    "summaries": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "id": {"type": "string"},
          "summary": {"type": "string"}
        },
        "required": ["id", "summary"]
      },
      "description": "One-line human summary per focus id"
    },
    "narrative": {
      "type": "string",
      "description": "2-3 sentence overview of what to focus on"
    }
  },
  "required": ["focus", "summaries", "narrative"]
}`

type toolFunctionDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type toolDef struct {
	Type     string          `json:"type"`
	Function toolFunctionDef `json:"function"`
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
	Tools       []toolDef     `json:"tools,omitempty"`
	ToolChoice  interface{}   `json:"tool_choice,omitempty"`
}

type toolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content   string `json:"content"`
			ToolCalls []struct {
				Function toolCallFunction `json:"function"`
			} `json:"tool_calls"`
			// Older OpenAI-compatible shape.
			FunctionCall *toolCallFunction `json:"function_call"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// toolArgs is the shape of the tool-call arguments (matches the schema above).
type toolArgs struct {
	Focus     []string `json:"focus"`
	Summaries []struct {
		ID      string `json:"id"`
		Summary string `json:"summary"`
	} `json:"summaries"`
	Narrative string `json:"narrative"`
}

// briefFinding is a compact, token-light view of a finding for the prompt, so
// a small model is not overwhelmed by the full finding records. The id is a
// short alias (F1, F2, …), never the 64-char fingerprint — see Summarize.
type briefFinding struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	CVE      string `json:"cve,omitempty"`
	Repo     string `json:"repo"`
	Package  string `json:"package,omitempty"`
}

// briefWaterfall is the prompt-facing view of a waterfall group, keyed by a
// short alias (W1, W2, …) for the same reason as briefFinding.
type briefWaterfall struct {
	ID            string `json:"id"`
	Severity      string `json:"severity"`
	RootCause     string `json:"rootCause,omitempty"`
	Ecosystem     string `json:"ecosystem,omitempty"`
	AffectedRepos int    `json:"affectedRepos,omitempty"`
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
	// Give every finding/group a short alias (F1, W1, …). The model echoes these
	// aliases in its tool call instead of the 64-char SHA-256 fingerprints; a
	// small model forced to reproduce the fingerprints overran max_tokens and
	// truncated the tool-call JSON mid-array. We translate aliases back to the
	// real ids after parsing, dropping any the model hallucinates.
	alias := make(map[string]string, len(findings)+len(cor.Waterfall))
	brief := struct {
		Findings  []briefFinding   `json:"findings"`
		Waterfall []briefWaterfall `json:"waterfall,omitempty"`
	}{}
	for i, f := range findings {
		id := fmt.Sprintf("F%d", i+1)
		alias[id] = f.ID
		brief.Findings = append(brief.Findings, briefFinding{
			ID: id, Severity: f.Severity, CVE: f.CVEID, Repo: f.Repo, Package: f.Package,
		})
	}
	for i, g := range cor.Waterfall {
		id := fmt.Sprintf("W%d", i+1)
		alias[id] = g.ID
		brief.Waterfall = append(brief.Waterfall, briefWaterfall{
			ID: id, Severity: g.Severity, RootCause: g.RootCause,
			Ecosystem: g.Ecosystem, AffectedRepos: len(g.AffectedRepos),
		})
	}
	payload, err := json.Marshal(brief)
	if err != nil {
		return nil, nil, "", err
	}

	prompt := "You are a security triage assistant for the Kairos project. " +
		"Analyse these correlated security findings and call the " + triageToolName +
		" function to report the most urgent items to focus on. " +
		"Only use the short id values (like F1 or W2) that appear in the input. Findings:\n" + string(payload)

	reqBody, err := json.Marshal(chatRequest{
		Model:       c.model,
		Messages:    []chatMessage{{Role: "user", Content: prompt}},
		Temperature: c.temperature,
		MaxTokens:   c.maxTokens,
		Tools: []toolDef{{
			Type: "function",
			Function: toolFunctionDef{
				Name:        triageToolName,
				Description: "Report prioritized security triage results.",
				Parameters:  json.RawMessage(triageToolParameters),
			},
		}},
		// Force the model to call our function rather than reply with prose.
		ToolChoice: map[string]interface{}{
			"type":     "function",
			"function": map[string]string{"name": triageToolName},
		},
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

	// Prefer the forced tool call; fall back to function_call, then to parsing
	// JSON out of the message content (for backends that ignore tools).
	msg := cr.Choices[0].Message
	var rawArgs string
	switch {
	case len(msg.ToolCalls) > 0 && msg.ToolCalls[0].Function.Arguments != "":
		rawArgs = msg.ToolCalls[0].Function.Arguments
	case msg.FunctionCall != nil && msg.FunctionCall.Arguments != "":
		rawArgs = msg.FunctionCall.Arguments
	case strings.TrimSpace(msg.Content) != "":
		rawArgs = extractJSON(msg.Content)
	default:
		return nil, nil, "", fmt.Errorf("model returned neither a tool call nor content (raw: %q)", snippet(body))
	}

	var args toolArgs
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return nil, nil, "", fmt.Errorf("tool-call arguments were not valid JSON: %w (args: %q)", err, snippet([]byte(rawArgs)))
	}

	// Translate aliases back to real ids, dropping any the model invented.
	focus := make([]string, 0, len(args.Focus))
	for _, a := range args.Focus {
		if real, ok := alias[a]; ok {
			focus = append(focus, real)
		}
	}
	summaries := make(map[string]string, len(args.Summaries))
	for _, s := range args.Summaries {
		if real, ok := alias[s.ID]; ok {
			summaries[real] = s.Summary
		}
	}
	return focus, summaries, args.Narrative, nil
}

// extractJSON pulls the JSON object out of a model reply that may be wrapped in
// markdown fences or surrounded by prose (content-fallback path only).
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
