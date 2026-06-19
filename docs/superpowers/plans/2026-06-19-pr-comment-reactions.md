# PR Comment Reactions + AI Prose — Implementation Plan (Plan 3)

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the remediation bot **react to review comments on the PRs it owns** — classify each new comment (request-change / question / nack / approve) with the self-hosted model, then adjust the bump to a maintainer-requested version, post a drafted reply, or close the PR — remembering which comments it already handled, and (optionally) drafting nicer PR prose with the model.

**Architecture:** A new comment-reaction layer in `internal/remediate`, split into: GitHub comment/close operations (`ghclient`, faked in tests), an AI **classifier** behind an interface (an OpenAI/LocalAI implementation using forced tool-calling, plus a fake), a **pure decision function** mapping a classification to an action, and a **pure-ish orchestrator** that ties them together over the ledger with `SeenComments` idempotency (react once per comment, never to the bot's own replies). The actual git re-bump (`Adjust`) and PR prose are deterministic-or-AI with deterministic fallback. Reactions run inside the existing `remediate` reconcile pass and stay **dry-run gated**.

**Tech Stack:** Go 1.22, `gh` CLI, the LocalAI OpenAI-compatible endpoint (forced tool calling, as in `internal/triage/openai.go`), `stretchr/testify`. Builds on Plan 1 (`ghclient`, `triage`) and Plan 2 (`remediate`, `state.Ledger`).

## Global Constraints

- Module `github.com/kairos-io/security`; binary `ksec`; Go 1.22.
- The bot reacts ONLY to comments on PRs **it owns** (ledger entries with `State=="open"` and `PRNumber>0`) and ONLY to comments authored by someone **other than** `kairos-security-bot` (never react to its own replies).
- **React once per comment:** every handled comment's id is appended to the ledger entry's `SeenComments`; a comment whose id is already there is skipped.
- **Classification is AI; the action mapping and the bump are deterministic.** The model returns an intent (+ an optional explicit version + a reply draft); the version change is still `go get <pkg>@<explicit-version>` verified by `go build ./...`.
- **AI is best-effort but a classification failure must not silently mis-act:** on classifier error, record `needs-human` in history and skip that comment (do not guess). Mirrors Plan 1/2 fail-safe behavior.
- `--dry-run` short-circuits every GitHub write (reply/close) and every git write (adjust/force-push) to a printed plan.
- Reactions are bounded by the existing run; no extra blast-radius beyond PRs already in the ledger.
- **Out of scope:** reacting to comments on human-authored PRs; multi-turn conversations (each comment is handled independently); changing PRs not in the ledger.
- **Operational prerequisite (not code):** live reactions need the same `KSEC_BOT_TOKEN` write scope as Plan 2 (`pull_requests:write` to comment/close, `contents:write` to force-push adjustments). Until then, runs are dry-run.

---

## File structure

```
internal/ghclient/ghclient.go        # + ReviewComment, ListPRComments, PostPRComment, ClosePR (modify)
internal/ghclient/fake.go            # + fakes for the above (modify)
internal/ghclient/comments_test.go   # (create)
internal/remediate/classify.go       # Classification, CommentClassifier interface, fake (create)
internal/remediate/openai_classify.go# OpenAI forced-tool-call classifier (create)
internal/remediate/openai_classify_test.go  # httptest parse test (create)
internal/remediate/react.go          # DecideReaction (pure) + ReactToComments orchestrator (create)
internal/remediate/react_test.go     # (create)
internal/remediate/git_executor.go   # + Adjust(entry, toVersion) re-bump+force-push (modify)
internal/remediate/prose.go          # AI PR title/body with deterministic fallback (create)
internal/remediate/prose_test.go     # (create)
cmd/ksec/main.go                     # wire reactions + prose into the remediate command (modify)
```

---

### Task 1: GitHub comment + close operations

**Files:**
- Modify: `internal/ghclient/ghclient.go`
- Modify: `internal/ghclient/fake.go`
- Test: `internal/ghclient/comments_test.go`

**Interfaces:**
- Consumes: the existing `GitHub` interface + `CLI`/`Fake`.
- Produces:
  - `type ReviewComment struct { ID, Author, Body string }`
  - three new `GitHub` interface methods: `ListPRComments(repo string, pr int) ([]ReviewComment, error)`, `PostPRComment(repo string, pr int, body string) error`, `ClosePR(repo string, pr int, comment string) error`.
  - `CLI` implementations (shell `gh`); `Fake` fields `PRComments map[string][]ReviewComment` (key `"<repo>#<pr>"`), `Posted []string` (record of `"<repo>#<pr>: <body>"`), `Closed []string` (record of `"<repo>#<pr>"`).

- [ ] **Step 1: Write the failing test (covers the Fake)**

Create `internal/ghclient/comments_test.go`:

```go
package ghclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakeCommentOps(t *testing.T) {
	f := NewFake()
	f.PRComments["kairos-io/immucore#412"] = []ReviewComment{
		{ID: "c1", Author: "maintainer", Body: "please pin to 0.36.0"},
	}
	got, err := f.ListPRComments("kairos-io/immucore", 412)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "c1", got[0].ID)

	require.NoError(t, f.PostPRComment("kairos-io/immucore", 412, "on it"))
	assert.Equal(t, []string{"kairos-io/immucore#412: on it"}, f.Posted)

	require.NoError(t, f.ClosePR("kairos-io/immucore", 412, "superseded"))
	assert.Equal(t, []string{"kairos-io/immucore#412"}, f.Closed)
}
```

- [ ] **Step 2: Run it — expect FAIL** (methods/type undefined). Run: `go test ./internal/ghclient/...`

- [ ] **Step 3: Implement on the interface + CLI + Fake**

In `internal/ghclient/ghclient.go`, add the type and extend the interface, then the `CLI` methods:

```go
type ReviewComment struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Body   string `json:"body"`
}
```

Add to the `GitHub` interface:

```go
	ListPRComments(repo string, pr int) ([]ReviewComment, error)
	PostPRComment(repo string, pr int, body string) error
	ClosePR(repo string, pr int, comment string) error
```

Add `CLI` methods:

```go
func (c *CLI) ListPRComments(repo string, pr int) ([]ReviewComment, error) {
	b, err := c.run("pr", "view", fmt.Sprint(pr), "-R", repo, "--json", "comments",
		"-q", "[.comments[] | {id: (.id|tostring), author: .author.login, body: .body}]")
	if err != nil {
		return nil, err
	}
	var out []ReviewComment
	return out, json.Unmarshal(b, &out)
}

func (c *CLI) PostPRComment(repo string, pr int, body string) error {
	_, err := c.run("pr", "comment", fmt.Sprint(pr), "-R", repo, "--body", body)
	return err
}

func (c *CLI) ClosePR(repo string, pr int, comment string) error {
	_, err := c.run("pr", "close", fmt.Sprint(pr), "-R", repo, "--comment", comment)
	return err
}
```

Add `Fake` fields (in `internal/ghclient/fake.go`) and methods. Extend `NewFake()` to initialize `PRComments`:

```go
// in Fake struct:
	PRComments map[string][]ReviewComment
	Posted     []string
	Closed     []string

// in NewFake():
		PRComments: map[string][]ReviewComment{},

func prKey(repo string, pr int) string { return fmt.Sprintf("%s#%d", repo, pr) }

func (f *Fake) ListPRComments(repo string, pr int) ([]ReviewComment, error) {
	return f.PRComments[prKey(repo, pr)], nil
}
func (f *Fake) PostPRComment(repo string, pr int, body string) error {
	f.Posted = append(f.Posted, prKey(repo, pr)+": "+body)
	return nil
}
func (f *Fake) ClosePR(repo string, pr int, comment string) error {
	f.Closed = append(f.Closed, prKey(repo, pr))
	return nil
}
```

(Add `"fmt"` to `fake.go` imports if not present.)

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/ghclient/...`

- [ ] **Step 5: Commit**

```bash
git add internal/ghclient/
git commit -m "feat(ghclient): PR comment list/post and close operations"
```

---

### Task 2: AI comment classifier (forced tool calling)

**Files:**
- Create: `internal/remediate/classify.go`
- Create: `internal/remediate/openai_classify.go`
- Test: `internal/remediate/openai_classify_test.go`

**Interfaces:**
- Consumes: `config.AIConfig` (for endpoint/model).
- Produces:
  - `type Classification struct { Intent, Version, Reply string }` where `Intent` is one of `"request-change"`, `"question"`, `"nack"`, `"approve"`, `"other"`.
  - `type CommentClassifier interface { Classify(prTitle, author, body string) (Classification, error) }`
  - `type FakeClassifier struct { Result Classification; Err error }` implementing it.
  - `type OpenAIClassifier struct { ... }` + `func NewOpenAIClassifier(cfg config.AIConfig) *OpenAIClassifier` — forces a `classify_comment` tool call (schema: intent enum, version, reply), parses the arguments. Mirror the structure of `internal/triage/openai.go` (chat-completions POST, forced `tool_choice`, tolerant parse).

- [ ] **Step 1: Write the classifier types + fake**

Create `internal/remediate/classify.go`:

```go
package remediate

type Classification struct {
	Intent  string // request-change | question | nack | approve | other
	Version string // explicit version requested, if any
	Reply   string // suggested reply text
}

type CommentClassifier interface {
	Classify(prTitle, author, body string) (Classification, error)
}

// FakeClassifier is a test double.
type FakeClassifier struct {
	Result Classification
	Err    error
}

func (f FakeClassifier) Classify(_, _, _ string) (Classification, error) {
	return f.Result, f.Err
}
```

- [ ] **Step 2: Write the failing test (httptest)**

Create `internal/remediate/openai_classify_test.go`:

```go
package remediate

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kairos-io/security/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAIClassifierForcesToolCall(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assert.Contains(t, string(body), "classify_comment")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"choices":[{"message":{"tool_calls":[{"function":{"name":"classify_comment",`+
			`"arguments":"{\"intent\":\"request-change\",\"version\":\"0.36.0\",\"reply\":\"Bumping to 0.36.0.\"}"}}]}}]}`)
	}))
	defer ts.Close()

	c := &OpenAIClassifier{endpoint: ts.URL, model: "m", httpc: ts.Client()}
	got, err := c.Classify("bump x", "maintainer", "please pin to 0.36.0")
	require.NoError(t, err)
	assert.Equal(t, "request-change", got.Intent)
	assert.Equal(t, "0.36.0", got.Version)
	assert.Equal(t, "Bumping to 0.36.0.", got.Reply)
}

func TestNewOpenAIClassifierReadsConfig(t *testing.T) {
	c := NewOpenAIClassifier(config.AIConfig{})
	require.NotNil(t, c)
}
```

- [ ] **Step 3: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 4: Implement the OpenAI classifier**

Create `internal/remediate/openai_classify.go`:

```go
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
```

(`extractJSON` is already defined in `internal/triage`... but this is package `remediate`. Add a local copy of `extractJSON` here OR — simpler — inline a minimal version. To avoid duplication, add a small unexported `extractJSON` in this file.)

```go
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}
```

- [ ] **Step 5: Run it — expect PASS.** Run: `go test ./internal/remediate/...`

- [ ] **Step 6: Commit**

```bash
git add internal/remediate/classify.go internal/remediate/openai_classify.go internal/remediate/openai_classify_test.go
git commit -m "feat(remediate): AI comment classifier via forced tool calling"
```

---

### Task 3: Pure reaction-decision logic

**Files:**
- Create: `internal/remediate/react.go` (first half — the pure `DecideReaction`)
- Test: `internal/remediate/react_test.go` (first half)

**Interfaces:**
- Consumes: `Classification`.
- Produces:
  - `type ReactionKind string` with `ReactReply`, `ReactAdjust`, `ReactClose`, `ReactNone`.
  - `type Reaction struct { Kind ReactionKind; ReplyBody, ToVersion string }`
  - `func DecideReaction(c Classification) Reaction` — `nack`→close (reply acks); `request-change`+non-empty `Version`→adjust to that version (reply included); `request-change` without a version OR `question`→reply; `approve`/`other`/anything else→none.

- [ ] **Step 1: Write the failing test**

Create `internal/remediate/react_test.go`:

```go
package remediate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecideReaction(t *testing.T) {
	cases := []struct {
		name string
		in   Classification
		want ReactionKind
	}{
		{"nack closes", Classification{Intent: "nack", Reply: "ok"}, ReactClose},
		{"request-change with version adjusts", Classification{Intent: "request-change", Version: "0.36.0", Reply: "ok"}, ReactAdjust},
		{"request-change without version replies", Classification{Intent: "request-change", Reply: "could you clarify?"}, ReactReply},
		{"question replies", Classification{Intent: "question", Reply: "it's automated"}, ReactReply},
		{"approve does nothing", Classification{Intent: "approve"}, ReactNone},
		{"other does nothing", Classification{Intent: "other"}, ReactNone},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, DecideReaction(tc.in).Kind)
		})
	}
	adj := DecideReaction(Classification{Intent: "request-change", Version: "0.36.0", Reply: "r"})
	assert.Equal(t, "0.36.0", adj.ToVersion)
	assert.Equal(t, "r", adj.ReplyBody)
}
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement `DecideReaction`**

Create `internal/remediate/react.go` with (the orchestrator is added in Task 4):

```go
package remediate

import (
	"fmt"

	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

type ReactionKind string

const (
	ReactReply  ReactionKind = "reply"
	ReactAdjust ReactionKind = "adjust"
	ReactClose  ReactionKind = "close"
	ReactNone   ReactionKind = "none"
)

type Reaction struct {
	Kind      ReactionKind
	ReplyBody string
	ToVersion string
}

func DecideReaction(c Classification) Reaction {
	switch c.Intent {
	case "nack":
		return Reaction{Kind: ReactClose, ReplyBody: c.Reply}
	case "request-change":
		if c.Version != "" {
			return Reaction{Kind: ReactAdjust, ToVersion: c.Version, ReplyBody: c.Reply}
		}
		return Reaction{Kind: ReactReply, ReplyBody: c.Reply}
	case "question":
		return Reaction{Kind: ReactReply, ReplyBody: c.Reply}
	default:
		return Reaction{Kind: ReactNone}
	}
}
```

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/remediate/...`

- [ ] **Step 5: Commit**

```bash
git add internal/remediate/react.go internal/remediate/react_test.go
git commit -m "feat(remediate): pure comment-reaction decision logic"
```

---

### Task 4: Reaction orchestrator (idempotent, dry-run aware)

**Files:**
- Modify: `internal/remediate/react.go` (add the orchestrator + `Adjuster` interface)
- Test: `internal/remediate/react_test.go` (add orchestrator tests)

**Interfaces:**
- Consumes: `ghclient.GitHub`, `CommentClassifier`, `Reaction`, `state.LedgerEntry`.
- Produces:
  - `type Adjuster interface { Adjust(entry state.LedgerEntry, toVersion, runID string) (state.LedgerEntry, error) }` (implemented by `GitExecutor` in Task 5; a fake here).
  - `func ReactToComments(entry *state.LedgerEntry, gh ghclient.GitHub, cls CommentClassifier, adj Adjuster, prTitle, runID string, dryRun bool) error` — for each comment NOT in `entry.SeenComments` and NOT authored by `kairos-security-bot`: classify (on classifier error, append a `needs-human` history event and skip — do not act); decide; execute (`ReactReply`→`gh.PostPRComment`; `ReactClose`→`gh.ClosePR` + set entry State `closed`; `ReactAdjust`→`adj.Adjust` then `gh.PostPRComment`; `ReactNone`→nothing); then append the comment id to `entry.SeenComments` and a history event. In `dryRun`, print intended actions and STILL record the comment id as seen (so a dry-run shows convergence) — but perform no writes. Only operates on entries with `State=="open"` and `PRNumber>0`.

- [ ] **Step 1: Write the failing test**

Add to `internal/remediate/react_test.go`:

```go
import (
	"github.com/kairos-io/security/internal/ghclient"
	"github.com/kairos-io/security/internal/state"
)

type fakeAdjuster struct{ called int }

func (f *fakeAdjuster) Adjust(e state.LedgerEntry, to, run string) (state.LedgerEntry, error) {
	f.called++
	e.Bump.To = to
	return e, nil
}

func TestReactToCommentsIdempotentAndActs(t *testing.T) {
	gh := ghclient.NewFake()
	gh.PRComments["kairos-io/immucore#412"] = []ghclient.ReviewComment{
		{ID: "c1", Author: "maintainer", Body: "pin to 0.36.0"},
		{ID: "self", Author: "kairos-security-bot", Body: "on it"}, // own reply: ignored
	}
	entry := &state.LedgerEntry{Key: "kairos-io/immucore|golang.org/x/net", Repo: "kairos-io/immucore",
		PRNumber: 412, State: "open", Bump: state.Bump{Package: "golang.org/x/net", To: "0.33.0"}}
	adj := &fakeAdjuster{}
	cls := FakeClassifier{Result: Classification{Intent: "request-change", Version: "0.36.0", Reply: "Bumping."}}

	require.NoError(t, ReactToComments(entry, gh, cls, adj, "bump", "run1", false))
	assert.Equal(t, 1, adj.called, "adjusted to requested version")
	assert.Equal(t, "0.36.0", entry.Bump.To)
	assert.Contains(t, gh.Posted, "kairos-io/immucore#412: Bumping.")
	assert.Contains(t, entry.SeenComments, "c1")
	assert.NotContains(t, entry.SeenComments, "self", "never react to its own comment")

	// Second run: c1 already seen -> no further action.
	adj.called = 0
	gh.Posted = nil
	require.NoError(t, ReactToComments(entry, gh, cls, adj, "bump", "run2", false))
	assert.Equal(t, 0, adj.called)
	assert.Empty(t, gh.Posted)
}

func TestReactToCommentsClassifierErrorSkips(t *testing.T) {
	gh := ghclient.NewFake()
	gh.PRComments["r#1"] = []ghclient.ReviewComment{{ID: "c1", Author: "m", Body: "?"}}
	entry := &state.LedgerEntry{Key: "r|p", Repo: "r", PRNumber: 1, State: "open"}
	require.NoError(t, ReactToComments(entry, gh, FakeClassifier{Err: assertErr()}, &fakeAdjuster{}, "t", "run1", false))
	assert.Empty(t, gh.Posted, "no action taken on classifier error")
	assert.NotContains(t, entry.SeenComments, "c1", "unhandled comment stays unseen for retry")
	assert.NotEmpty(t, entry.History)
}

func assertErr() error { return assert.AnError }
```

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement the orchestrator**

Append to `internal/remediate/react.go`:

```go
const botLogin = "kairos-security-bot"

type Adjuster interface {
	Adjust(entry state.LedgerEntry, toVersion, runID string) (state.LedgerEntry, error)
}

func seen(entry *state.LedgerEntry, id string) bool {
	for _, s := range entry.SeenComments {
		if s == id {
			return true
		}
	}
	return false
}

func ReactToComments(entry *state.LedgerEntry, gh ghclient.GitHub, cls CommentClassifier, adj Adjuster, prTitle, runID string, dryRun bool) error {
	if entry.State != "open" || entry.PRNumber == 0 {
		return nil
	}
	comments, err := gh.ListPRComments(entry.Repo, entry.PRNumber)
	if err != nil {
		return err
	}
	for _, cm := range comments {
		if cm.Author == botLogin || seen(entry, cm.ID) {
			continue
		}
		cl, err := cls.Classify(prTitle, cm.Author, cm.Body)
		if err != nil {
			// Do not guess; leave the comment unseen so a later run can retry.
			entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "needs-human", Detail: "classify failed: " + err.Error()})
			continue
		}
		r := DecideReaction(cl)
		action := string(r.Kind)
		if dryRun {
			fmt.Printf("[dry-run] would react to %s#%d comment %s: %s\n", entry.Repo, entry.PRNumber, cm.ID, action)
		} else {
			switch r.Kind {
			case ReactReply:
				if err := gh.PostPRComment(entry.Repo, entry.PRNumber, r.ReplyBody); err != nil {
					return err
				}
			case ReactAdjust:
				updated, err := adj.Adjust(*entry, r.ToVersion, runID)
				if err != nil {
					return err
				}
				*entry = updated
				if r.ReplyBody != "" {
					if err := gh.PostPRComment(entry.Repo, entry.PRNumber, r.ReplyBody); err != nil {
						return err
					}
				}
			case ReactClose:
				if err := gh.ClosePR(entry.Repo, entry.PRNumber, r.ReplyBody); err != nil {
					return err
				}
				entry.State = "closed"
			case ReactNone:
			}
		}
		entry.SeenComments = append(entry.SeenComments, cm.ID)
		entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "reacted", Detail: action + " to " + cm.ID})
	}
	return nil
}
```

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/remediate/...`

- [ ] **Step 5: Commit**

```bash
git add internal/remediate/react.go internal/remediate/react_test.go
git commit -m "feat(remediate): idempotent comment-reaction orchestrator"
```

---

### Task 5: `GitExecutor.Adjust` (re-bump + force-push)

**Files:**
- Modify: `internal/remediate/git_executor.go`

**Interfaces:**
- Consumes: `state.LedgerEntry`, the existing `GitExecutor` helpers (`cloneURL`, `g.run`, `BranchName` style).
- Produces: `func (g *GitExecutor) Adjust(entry state.LedgerEntry, toVersion, runID string) (state.LedgerEntry, error)` implementing `Adjuster`. Re-clone the repo, check out the existing `entry.Branch`, `go get <package>@<toVersion>`, `go mod tidy`, `go build ./...` (verify), commit, **force-push** the branch; update `entry.Bump.To`, append a history event. Dry-run prints and returns the entry with the new `Bump.To` but performs no writes. Build failure → record `build-failed` history, return entry with **nil** error (no force-push). Token never logged (uses `g.run`, already redacting).

- [ ] **Step 1: Implement**

Append to `internal/remediate/git_executor.go`:

```go
var _ Adjuster = (*GitExecutor)(nil)

func (g *GitExecutor) Adjust(entry state.LedgerEntry, toVersion, runID string) (state.LedgerEntry, error) {
	entry.LastActionRun = runID
	if g.DryRun {
		fmt.Printf("[dry-run] would adjust %s PR #%d: go get %s@%s, force-push %s\n",
			entry.Repo, entry.PRNumber, entry.Package, toVersion, entry.Branch)
		entry.Bump.To = toVersion
		return entry, nil
	}
	dir, err := os.MkdirTemp("", "ksec-adj-*")
	if err != nil {
		return entry, err
	}
	defer os.RemoveAll(dir)

	if _, err := g.run("", "git", "clone", g.cloneURL(entry.Repo), dir); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "git", "checkout", entry.Branch); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "go", "get", entry.Package+"@"+toVersion); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "go", "mod", "tidy"); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "go", "build", "./..."); err != nil {
		entry.State = "build-failed"
		entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "adjust-build-failed", Detail: err.Error()})
		return entry, nil
	}
	_, _ = g.run(dir, "git", "config", "user.name", "kairos-security-bot")
	_, _ = g.run(dir, "git", "config", "user.email", "bot@kairos.io")
	if _, err := g.run(dir, "git", "commit", "-am", "chore(security): adjust bump to "+toVersion); err != nil {
		return entry, err
	}
	if _, err := g.run(dir, "git", "push", "--force", "origin", entry.Branch); err != nil {
		return entry, err
	}
	entry.Bump.To = toVersion
	entry.History = append(entry.History, state.LedgerEvent{Run: runID, Action: "adjusted", Detail: "to " + toVersion})
	return entry, nil
}
```

(Note: the repo is cloned WITHOUT `--depth 1` here so the existing branch and its history are available to check out and amend; the default-branch clone in `Open` uses `--depth 1`, but `Adjust` needs the feature branch — clone full or add `--branch entry.Branch`. Use `git clone <url> dir` then `git checkout entry.Branch` as written; if the shallow optimization matters, clone with `--branch entry.Branch --depth 1` instead.)

- [ ] **Step 2: Build + vet + test**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: pass (no new unit test; `Adjust` is integration, exercised by the dry-run path and Task 6).

- [ ] **Step 3: Commit**

```bash
git add internal/remediate/git_executor.go
git commit -m "feat(remediate): GitExecutor.Adjust re-bump + force-push (dry-run aware)"
```

---

### Task 6: Wire reactions into the `remediate` command

**Files:**
- Modify: `cmd/ksec/main.go`

**Interfaces:**
- Consumes: `remediate.ReactToComments`, `remediate.NewOpenAIClassifier`, `remediate.GitExecutor` (as `Adjuster`), `ghclient.NewCLI`, `config.LoadAI`, the ledger.
- Produces: after the existing `remediate.Run(...)` updates the ledger, iterate the resulting ledger entries and call `ReactToComments` on each open PR entry, then save the ledger. Comment reactions honor the same `--dry-run`.

- [ ] **Step 1: Wire it**

In `newRemediateCmd`'s `RunE`, AFTER `out, results := remediate.Run(...)` and the result logging, BEFORE the dry-run save guard, add:

```go
			// React to review comments on PRs we own.
			aiCfg, _ := config.LoadAI("ai.yaml")
			classifier := remediate.NewOpenAIClassifier(aiCfg)
			gh := ghclient.NewCLI()
			for i := range out.Entries {
				e := &out.Entries[i]
				if e.State != "open" || e.PRNumber == 0 {
					continue
				}
				title := remediate.PRTitle(remediate.Intent{Package: e.Package, Bump: e.Bump})
				if err := remediate.ReactToComments(e, gh, classifier, ex, title, runID, gf.dryRun); err != nil {
					fmt.Fprintf(os.Stderr, "remediate: react %s: %v\n", e.Key, err)
				}
			}
```

(`ex` is the `*remediate.GitExecutor` already constructed in this function; it satisfies `Adjuster`. `config` and `ghclient` are already imported in main.go.)

- [ ] **Step 2: Build + full test + smoke**

Run: `go build ./... && go vet ./... && go test ./...`
Expected: all pass. Smoke (dry-run, no network writes): with a `ledger.json` containing an `open` entry the reaction loop runs in dry-run and prints `[dry-run] would react ...` only if the PR has comments (the gh call is read-only). Confirm `go run ./cmd/ksec remediate --help` still works.

- [ ] **Step 3: Commit**

```bash
git add cmd/ksec/main.go
git commit -m "feat(remediate): wire comment reactions into the remediate phase"
```

---

### Task 7: AI-drafted PR prose (optional, deterministic fallback)

**Files:**
- Create: `internal/remediate/prose.go`
- Test: `internal/remediate/prose_test.go`
- Modify: `internal/remediate/git_executor.go` (use prose in `Open` when available)

**Interfaces:**
- Consumes: `Intent`, the existing deterministic `PRTitle`/`PRBody`.
- Produces:
  - `type ProseClient interface { DraftPRBody(in Intent) (string, error) }`
  - `func PRBodyWith(in Intent, prose ProseClient) string` — if `prose` is non-nil and returns a non-empty body without error, append the AI prose to the deterministic body (keeping the marker last); otherwise return the deterministic `PRBody(in)`. The marker MUST remain the last line either way.

- [ ] **Step 1: Write the failing test**

Create `internal/remediate/prose_test.go`:

```go
package remediate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeProse struct {
	body string
	err  error
}

func (f fakeProse) DraftPRBody(Intent) (string, error) { return f.body, f.err }

func TestPRBodyWithFallsBackOnError(t *testing.T) {
	in := sampleIntent()
	out := PRBodyWith(in, fakeProse{err: assertErrP()})
	assert.Equal(t, PRBody(in), out, "error -> deterministic body")
}

func TestPRBodyWithAppendsAIAndKeepsMarkerLast(t *testing.T) {
	in := sampleIntent()
	out := PRBodyWith(in, fakeProse{body: "AI: this is safe."})
	assert.Contains(t, out, "AI: this is safe.")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(out), PRMarker(in.Key)))
}

func assertErrP() error { return assert.AnError }
```

(`sampleIntent()` already exists in `prbody_test.go` within this package.)

- [ ] **Step 2: Run it — expect FAIL.** Run: `go test ./internal/remediate/...`

- [ ] **Step 3: Implement**

Create `internal/remediate/prose.go`:

```go
package remediate

import "strings"

type ProseClient interface {
	DraftPRBody(in Intent) (string, error)
}

// PRBodyWith returns the deterministic body, optionally enriched with an AI
// paragraph inserted before the trailing marker. On any AI error or empty
// output it returns the plain deterministic body. The marker stays last.
func PRBodyWith(in Intent, prose ProseClient) string {
	body := PRBody(in)
	if prose == nil {
		return body
	}
	extra, err := prose.DraftPRBody(in)
	if err != nil || strings.TrimSpace(extra) == "" {
		return body
	}
	marker := PRMarker(in.Key)
	withoutMarker := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(body), marker))
	return withoutMarker + "\n\n" + strings.TrimSpace(extra) + "\n\n" + marker
}
```

In `git_executor.go` `Open`, change the PR creation to use prose when a `ProseClient` is configured. Add a `Prose ProseClient` field to `GitExecutor` (nil by default → deterministic), and in `Open` replace `PRBody(in)` in the `gh pr create` call with `PRBodyWith(in, g.Prose)`. (Leave `Prose` nil unless wired; this keeps Plan 3 shippable with deterministic bodies and AI prose strictly opt-in.)

- [ ] **Step 4: Run it — expect PASS.** Run: `go test ./internal/remediate/...`

- [ ] **Step 5: Commit**

```bash
git add internal/remediate/prose.go internal/remediate/prose_test.go internal/remediate/git_executor.go
git commit -m "feat(remediate): optional AI PR prose with deterministic fallback"
```

---

## Self-review

**Spec coverage** (design §8 comment-reaction lifecycle):
- New review comments (not in SeenComments) → classify → act → record → Tasks 1, 2, 4. ✓
- Classify intent {request-change | question | nack | approve} → Task 2 (AI, forced tool call) + `other`. ✓
- Adjust bump to explicit version & push → Task 5 (`Adjust`, deterministic `go get @ver` + force-push, verify-before-push). ✓
- Post drafted reply → Task 4 (`ReactReply` → `PostPRComment`). ✓
- Close + record if nack → Task 4 (`ReactClose` → `ClosePR`, State `closed`). ✓
- Only touch own PRs; never react to own comments → Task 4 (`State=="open"`+`PRNumber>0`, skip `botLogin`). ✓
- React once (idempotent) → Task 4 (`SeenComments`). ✓
- Bump deterministic; only reaction is AI → Tasks 2 (AI classify), 5 (deterministic bump). ✓
- AI best-effort, never mis-act → Task 4 (classify error → `needs-human`, skip, leave unseen for retry). ✓
- Dry-run short-circuits writes → Tasks 4, 5. ✓
- AI PR prose (deferred from Plan 2) → Task 7, opt-in, deterministic fallback. ✓

**Placeholder scan:** none — complete code in every step.

**Type consistency:** `ReviewComment` + the three `GitHub` methods (Task 1) used by Task 4. `Classification`/`CommentClassifier`/`FakeClassifier` (Task 2) used by Tasks 4, 6. `OpenAIClassifier` (Task 2) constructed in Task 6. `Reaction`/`ReactionKind`/`DecideReaction` (Task 3) used by Task 4. `Adjuster` (Task 4) implemented by `GitExecutor.Adjust` (Task 5) and passed as `ex` in Task 6. `ProseClient`/`PRBodyWith` (Task 7) used in `Open`. `PRTitle`/`PRMarker`/`PRBody`/`Intent`/`sampleIntent` reused from Plan 2.

**Note on `extractJSON` duplication:** Plan 1's `internal/triage` has an `extractJSON`; Task 2 adds a smaller copy in `internal/remediate` (different package). Acceptable (no cross-package export); a future refactor could lift it to a shared `internal/ai` package — out of scope here.

---

## Operational notes (resolve before live reactions)

- **Token:** same `KSEC_BOT_TOKEN` write scope as Plan 2 (`pull_requests:write` for comment/close, `contents:write` for force-push). Reactions stay dry-run until then.
- **Model quality:** comment classification leans on the small model. `intent` is a constrained enum (reliable under grammar), but a wrong `version` from the model would cause an unwanted re-bump — `Adjust` still verifies `go build` before force-pushing, so a bad version that breaks the build is recorded `build-failed`, not pushed. Consider requiring a maintainer label (e.g. `ksec-approve-adjust`) before honoring `request-change` adjustments as a future safety gate.
- **No multi-turn memory:** each comment is handled once, independently. Threaded back-and-forth is out of scope.
