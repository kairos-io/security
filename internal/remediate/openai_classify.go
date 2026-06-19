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

type OpenAIClassifier struct {
	endpoint    string
	model       string
	temperature float64
	httpc       *http.Client
}

func NewOpenAIClassifier(cfg config.AIConfig) *OpenAIClassifier {
	return &OpenAIClassifier{
		endpoint:    strings.TrimRight(cfg.Nib.Endpoint, "/"),
		model:       cfg.Nib.Model,
		temperature: cfg.Nib.Temperature,
		httpc:       &http.Client{Timeout: 3 * time.Minute},
	}
}

const classifyToolName = "classify_comment"

const classifyToolParameters = `{
  "type": "object",
  "properties": {
    "intent": {"type": "string", "enum": ["request-change","question","nack","approve","other"]},
    "version": {"type": "string", "description": "explicit version requested, else empty"},
    "reply": {"type": "string", "description": "a short, polite reply to post"}
  },
  "required": ["intent","reply"]
}`

func (c *OpenAIClassifier) Classify(prTitle, author, body string) (Classification, error) {
	if c.endpoint == "" {
		return Classification{}, fmt.Errorf("no AI endpoint configured")
	}
	prompt := fmt.Sprintf("A maintainer (%s) left this review comment on the automated dependency-bump PR "+
		"titled %q:\n\n%s\n\nClassify the intent and draft a short reply. If they ask for a specific version, "+
		"put it in `version`. Call the %s function.", author, prTitle, body, classifyToolName)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":       c.model,
		"temperature": c.temperature,
		"messages":    []map[string]string{{"role": "user", "content": prompt}},
		"tools": []map[string]interface{}{{
			"type": "function",
			"function": map[string]interface{}{
				"name":        classifyToolName,
				"description": "Report the classification of a PR review comment.",
				"parameters":  json.RawMessage(classifyToolParameters),
			},
		}},
		"tool_choice": map[string]interface{}{"type": "function", "function": map[string]string{"name": classifyToolName}},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/v1/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return Classification{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpc.Do(req)
	if err != nil {
		return Classification{}, fmt.Errorf("call %s: %w", c.endpoint, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return Classification{}, fmt.Errorf("classify HTTP %d", resp.StatusCode)
	}

	var cr struct {
		Choices []struct {
			Message struct {
				Content   string `json:"content"`
				ToolCalls []struct {
					Function struct {
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &cr); err != nil {
		return Classification{}, fmt.Errorf("decode classify response: %w", err)
	}
	if len(cr.Choices) == 0 {
		return Classification{}, fmt.Errorf("classify: no choices")
	}
	args := ""
	if tc := cr.Choices[0].Message.ToolCalls; len(tc) > 0 {
		args = tc[0].Function.Arguments
	} else {
		args = extractJSON(cr.Choices[0].Message.Content)
	}
	var out Classification
	if err := json.Unmarshal([]byte(args), &out); err != nil {
		return Classification{}, fmt.Errorf("classify args not valid JSON: %w", err)
	}
	return out, nil
}

func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}
