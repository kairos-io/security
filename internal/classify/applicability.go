package classify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kairos-io/security/internal/state"
)

// Applier annotates findings with an AI applicability verdict. It never sets
// Class — the finding stays counted and actionable. The verdict is advisory
// metadata surfaced by renderers (a warning + reasoning popup) so an operator
// still sees the finding but knows the model suspects it may not affect us.
type Applier interface {
	// Apply returns findings with Finding.AIApplicability populated for the
	// ones the model reached a confident non-applicable verdict on.
	// Confidence-threshold and error handling are implementation details —
	// see OpenAIApplier for the concrete contract.
	Apply(findings []state.Finding) []state.Finding
}

// OpenAIApplier is an OpenAI-compatible (LocalAI) tool-call classifier that
// judges CVE applicability from the upstream Details + AffectedRanges the
// collector already recorded. It never makes a second HTTP fetch of its own —
// if the finding has no Details, it is skipped (no verdict attached).
//
// Verdicts are memoized per (cveID, package, currentVersion, fixedVersion) so
// the same CVE hitting multiple repos only calls the model once per run.
type OpenAIApplier struct {
	Endpoint            string // base URL, e.g. http://localhost:8080
	Model               string
	Temperature         float64
	MaxTokens           int
	ConfidenceThreshold string // "high" (default) | "medium" — min confidence to record a not-applicable verdict
	Today               string // YYYY-MM-DD stamped onto verdicts; empty = time.Now().UTC()
	HTTPClient          *http.Client
}

// NewOpenAIApplier returns a classifier with sane defaults. Missing endpoint
// yields a nil-safe classifier whose Apply is a no-op (fail-open).
func NewOpenAIApplier(endpoint, model string, temperature float64, maxTokens int, threshold string) *OpenAIApplier {
	if endpoint == "" {
		return nil
	}
	if maxTokens <= 0 {
		maxTokens = 1024
	}
	if threshold == "" {
		threshold = "high"
	}
	return &OpenAIApplier{
		Endpoint:            strings.TrimRight(endpoint, "/"),
		Model:               model,
		Temperature:         temperature,
		MaxTokens:           maxTokens,
		ConfidenceThreshold: threshold,
		HTTPClient:          &http.Client{Timeout: 3 * time.Minute},
	}
}

// Apply attaches an AIApplicability record to each finding the classifier
// reaches a confident non-applicable verdict on. Findings without upstream
// Details/AffectedRanges are skipped (nothing to reason over). Any transport
// or decode error is logged and the finding is left unannotated. Findings
// already marked informational (accepted-component or already-fixed) are
// skipped — no point spending model tokens on them.
func (a *OpenAIApplier) Apply(findings []state.Finding) []state.Finding {
	if a == nil || a.Endpoint == "" {
		return findings
	}
	out := make([]state.Finding, len(findings))
	copy(out, findings)

	stamp := a.Today
	if stamp == "" {
		stamp = time.Now().UTC().Format("2006-01-02")
	}

	// Memoize verdicts per (cve, pkg, current, fixed) — a CVE hitting N repos
	// gets one model call, not N.
	type memoKey struct{ cve, pkg, cur, fix string }
	type memoVal struct {
		args applicabilityArgs
		err  error
	}
	cache := map[memoKey]memoVal{}

	for i := range out {
		f := &out[i]
		if f.Class == "informational" {
			continue
		}
		if strings.TrimSpace(f.Details) == "" && strings.TrimSpace(f.AffectedRanges) == "" {
			continue
		}
		k := memoKey{f.CVEID, f.Package, f.CurrentVersion, f.FixedVersion}
		v, seen := cache[k]
		if !seen {
			args, err := a.classify(*f)
			v = memoVal{args: args, err: err}
			cache[k] = v
		}
		if v.err != nil {
			fmt.Printf("classify: applicability query failed for %s (%s): %v\n", f.CVEID, f.Package, v.err)
			continue
		}
		// Only surface a warning for a confident NOT-applicable verdict —
		// low-confidence noise or an "applicable" verdict would just clutter
		// the dashboard.
		if v.args.Applicable || !meetsThreshold(v.args.Confidence, a.ConfidenceThreshold) {
			continue
		}
		reason := strings.TrimSpace(v.args.Reasoning)
		if reason == "" {
			reason = "model reports finding does not apply to the queried version"
		}
		f.AIApplicability = &state.AIApplicability{
			Applicable: false,
			Confidence: strings.ToLower(strings.TrimSpace(v.args.Confidence)),
			Reasoning:  reason,
			Model:      a.Model,
			CheckedAt:  stamp,
		}
	}
	return out
}

// meetsThreshold reports whether got clears the bar. Order: high > medium > low.
func meetsThreshold(got, want string) bool {
	rank := map[string]int{"high": 3, "medium": 2, "low": 1, "": 0}
	return rank[strings.ToLower(got)] >= rank[strings.ToLower(want)]
}

// applicabilityToolName is the forced function-call name the model must invoke.
const applicabilityToolName = "classify_applicability"

const applicabilityToolParameters = `{
  "type": "object",
  "properties": {
    "applicable": {
      "type": "boolean",
      "description": "true if the CVE actually affects the queried package at the queried version; false only when the affected ranges clearly exclude it"
    },
    "confidence": {
      "type": "string",
      "enum": ["low", "medium", "high"],
      "description": "how sure the reasoning is; low forces a fail-open (finding stays actionable)"
    },
    "reasoning": {
      "type": "string",
      "description": "one short sentence explaining why (cite version ranges when possible)"
    }
  },
  "required": ["applicable", "confidence", "reasoning"]
}`

type applicabilityArgs struct {
	Applicable bool   `json:"applicable"`
	Confidence string `json:"confidence"`
	Reasoning  string `json:"reasoning"`
}

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
	ToolChoice  any   `json:"tool_choice,omitempty"`
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
			FunctionCall *toolCallFunction `json:"function_call"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (a *OpenAIApplier) classify(f state.Finding) (applicabilityArgs, error) {
	prompt := buildApplicabilityPrompt(f)
	body, err := json.Marshal(chatRequest{
		Model:       a.Model,
		Messages:    []chatMessage{{Role: "user", Content: prompt}},
		Temperature: a.Temperature,
		MaxTokens:   a.MaxTokens,
		Tools: []toolDef{{
			Type: "function",
			Function: toolFunctionDef{
				Name:        applicabilityToolName,
				Description: "Report whether a CVE applies to the queried package version.",
				Parameters:  json.RawMessage(applicabilityToolParameters),
			},
		}},
		ToolChoice: map[string]any{
			"type":     "function",
			"function": map[string]string{"name": applicabilityToolName},
		},
	})
	if err != nil {
		return applicabilityArgs{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.Endpoint+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return applicabilityArgs{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	httpc := a.HTTPClient
	if httpc == nil {
		httpc = &http.Client{Timeout: 3 * time.Minute}
	}
	resp, err := httpc.Do(req)
	if err != nil {
		return applicabilityArgs{}, fmt.Errorf("call %s: %w", a.Endpoint, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return applicabilityArgs{}, fmt.Errorf("chat endpoint returned HTTP %d: %s", resp.StatusCode, snippet(raw))
	}

	var cr chatResponse
	if err := json.Unmarshal(raw, &cr); err != nil {
		return applicabilityArgs{}, fmt.Errorf("decode chat response: %w (raw: %q)", err, snippet(raw))
	}
	if cr.Error != nil {
		return applicabilityArgs{}, fmt.Errorf("chat endpoint error: %s", cr.Error.Message)
	}
	if len(cr.Choices) == 0 {
		return applicabilityArgs{}, fmt.Errorf("chat response had no choices (raw: %q)", snippet(raw))
	}
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
		return applicabilityArgs{}, fmt.Errorf("model returned neither a tool call nor content (raw: %q)", snippet(raw))
	}
	var args applicabilityArgs
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return applicabilityArgs{}, fmt.Errorf("tool-call arguments were not valid JSON: %w (args: %q)", err, snippet([]byte(rawArgs)))
	}
	return args, nil
}

// buildApplicabilityPrompt renders a compact human-readable prompt from a
// finding. Kept in a separate function so tests can exercise it without a
// live model.
func buildApplicabilityPrompt(f state.Finding) string {
	var b strings.Builder
	b.WriteString("You are a security analyst deciding whether a CVE actually affects a specific package version.\n\n")
	b.WriteString("Rules:\n")
	b.WriteString("- Answer applicable=false ONLY when the upstream affected-version data clearly excludes the queried version.\n")
	b.WriteString("- If you are unsure, answer applicable=true with confidence=low (leaves the finding actionable).\n")
	b.WriteString("- Cite the range or condition in `reasoning` (one short sentence).\n\n")
	fmt.Fprintf(&b, "CVE: %s\n", nonEmpty(f.CVEID, "(none)"))
	fmt.Fprintf(&b, "Package: %s (ecosystem: %s)\n", nonEmpty(f.Package, "(none)"), nonEmpty(f.Ecosystem, "(none)"))
	fmt.Fprintf(&b, "Queried version: %s\n", nonEmpty(f.CurrentVersion, "(none)"))
	fmt.Fprintf(&b, "Reported fixed version: %s\n", nonEmpty(f.FixedVersion, "(none)"))
	fmt.Fprintf(&b, "Severity: %s\n", nonEmpty(f.Severity, "unknown"))
	if f.Title != "" {
		fmt.Fprintf(&b, "Title: %s\n", f.Title)
	}
	if f.Details != "" {
		b.WriteString("\n--- CVE details ---\n")
		b.WriteString(truncate(f.Details, 4000))
		b.WriteString("\n--- end details ---\n")
	}
	if f.AffectedRanges != "" {
		b.WriteString("\n--- Upstream affected ranges (JSON) ---\n")
		b.WriteString(truncate(f.AffectedRanges, 4000))
		b.WriteString("\n--- end affected ---\n")
	}
	b.WriteString("\nCall classify_applicability with your verdict.\n")
	return b.String()
}

func nonEmpty(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "… (truncated)"
}

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
