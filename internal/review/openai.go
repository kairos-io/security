package review

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
	"github.com/kairos-io/security/internal/ghclient"
)

// OpenAIAssessor reviews bot-authored PRs against an OpenAI-compatible
// chat-completions endpoint (LocalAI). Like triage.OpenAIClient it forces a
// function/tool call whose parameters are a JSON schema, so the backend
// grammar-constrains the output to {verdict, reasoning}.
//
// It never returns a hard error: any failure (no endpoint, non-200, decode
// error, missing tool call, out-of-enum verdict) degrades to the safe
// "needs_human_verification" verdict so a human takes a look.
type OpenAIAssessor struct {
	endpoint    string // base URL, e.g. http://localhost:8080
	model       string
	maxTokens   int
	temperature float64
	httpc       *http.Client
}

func NewOpenAIAssessor(cfg config.AIConfig) *OpenAIAssessor {
	maxTokens := cfg.Nib.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 1024
	}
	return &OpenAIAssessor{
		endpoint:    strings.TrimRight(cfg.Nib.Endpoint, "/"),
		model:       cfg.Nib.Model,
		maxTokens:   maxTokens,
		temperature: cfg.Nib.Temperature,
		httpc:       &http.Client{Timeout: 5 * time.Minute},
	}
}

const (
	assessToolName  = "assess_pr"
	verdictNeeds    = "needs_human_verification"
	maxContextBytes = 60000
)

// assessToolParameters constrains the verdict to the three-value enum and
// requires a free-form reasoning string plus a short changes summary.
const assessToolParameters = `{
  "type": "object",
  "properties": {
    "verdict": {
      "type": "string",
      "enum": ["good", "bad", "needs_human_verification"],
      "description": "good = safe to auto-approve; bad = should not be merged; needs_human_verification = a human must review"
    },
    "reasoning": {
      "type": "string",
      "description": "One to three sentences explaining the verdict"
    },
    "changesSummary": {
      "type": "string",
      "description": "One to three sentences summarizing what the dependency change actually does, based on the changelog and upstream source diff"
    }
  },
  "required": ["verdict", "reasoning", "changesSummary"]
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

// assessArgs is the shape of the tool-call arguments (matches the schema above).
type assessArgs struct {
	Verdict        string `json:"verdict"`
	Reasoning      string `json:"reasoning"`
	ChangesSummary string `json:"changesSummary"`
}

func (a *OpenAIAssessor) Assess(pr ghclient.PullRequest, reviewContext string) (string, string, string, error) {
	if a.endpoint == "" {
		return verdictNeeds, "no AI endpoint configured; needs human verification", "", nil
	}

	truncated := reviewContext
	note := ""
	if len(truncated) > maxContextBytes {
		truncated = truncated[:maxContextBytes]
		note = "\n\n[context truncated]"
	}

	prompt := "You are a security reviewer for the Kairos project assessing an " +
		"automated dependency/bot pull request. Decide whether the change is safe " +
		"to auto-approve. The context below contains the PR description / changelog, " +
		"the upstream source diff for each dependency bump, and the PR's own diff. " +
		"Call the " + assessToolName + " function with your verdict, your reasoning, " +
		"and a one to three sentence changesSummary describing what the dependency " +
		"change actually does.\n\n" +
		"PR title: " + pr.Title + "\n\nContext:\n" + truncated + note

	reqBody, err := json.Marshal(chatRequest{
		Model:       a.model,
		Messages:    []chatMessage{{Role: "user", Content: prompt}},
		Temperature: a.temperature,
		MaxTokens:   a.maxTokens,
		Tools: []toolDef{{
			Type: "function",
			Function: toolFunctionDef{
				Name:        assessToolName,
				Description: "Report the review verdict for a bot pull request.",
				Parameters:  json.RawMessage(assessToolParameters),
			},
		}},
		// Force the model to call our function rather than reply with prose.
		ToolChoice: map[string]interface{}{
			"type":     "function",
			"function": map[string]string{"name": assessToolName},
		},
	})
	if err != nil {
		return verdictNeeds, "could not build review request: " + err.Error(), "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		a.endpoint+"/v1/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return verdictNeeds, "could not build review request: " + err.Error(), "", nil
	}
	req.Header.Set("Content-Type", "application/json")

	httpResp, err := a.httpc.Do(req)
	if err != nil {
		return verdictNeeds, fmt.Sprintf("review endpoint unreachable: %v", err), "", nil
	}
	defer httpResp.Body.Close()
	body, _ := io.ReadAll(httpResp.Body)
	if httpResp.StatusCode != http.StatusOK {
		return verdictNeeds, fmt.Sprintf("review endpoint returned HTTP %d", httpResp.StatusCode), "", nil
	}

	var cr chatResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return verdictNeeds, "could not decode review response", "", nil
	}
	if cr.Error != nil {
		return verdictNeeds, "review endpoint error: " + cr.Error.Message, "", nil
	}
	if len(cr.Choices) == 0 {
		return verdictNeeds, "review response had no choices", "", nil
	}

	msg := cr.Choices[0].Message
	var rawArgs string
	switch {
	case len(msg.ToolCalls) > 0 && msg.ToolCalls[0].Function.Arguments != "":
		rawArgs = msg.ToolCalls[0].Function.Arguments
	case msg.FunctionCall != nil && msg.FunctionCall.Arguments != "":
		rawArgs = msg.FunctionCall.Arguments
	default:
		return verdictNeeds, "model did not return a tool call", "", nil
	}

	var args assessArgs
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return verdictNeeds, "tool-call arguments were not valid JSON", "", nil
	}

	switch args.Verdict {
	case "good", "bad", verdictNeeds:
		return args.Verdict, args.Reasoning, args.ChangesSummary, nil
	default:
		return verdictNeeds, fmt.Sprintf("model returned an unrecognized verdict %q; needs human verification", args.Verdict), "", nil
	}
}
