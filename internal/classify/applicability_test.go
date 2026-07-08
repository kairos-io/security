package classify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kairos-io/security/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildApplicabilityPrompt_IncludesFields(t *testing.T) {
	f := state.Finding{
		CVEID: "CVE-2024-0001", Package: "openssl", Ecosystem: "hadron",
		CurrentVersion: "3.1.2", FixedVersion: "3.2.0", Severity: "high",
		Title: "TLS bug", Details: "Some CVE description here.",
		AffectedRanges: `[{"ranges":[{"events":[{"introduced":"0"},{"fixed":"3.2.0"}]}]}]`,
	}
	prompt := buildApplicabilityPrompt(f)
	for _, want := range []string{"CVE-2024-0001", "openssl", "hadron", "3.1.2", "3.2.0", "TLS bug", "CVE description", "introduced"} {
		assert.Contains(t, prompt, want, "prompt must include %q", want)
	}
}

func TestBuildApplicabilityPrompt_TruncatesLongDetails(t *testing.T) {
	f := state.Finding{CVEID: "CVE-1", Package: "p", Details: strings.Repeat("x", 5000)}
	prompt := buildApplicabilityPrompt(f)
	assert.Contains(t, prompt, "(truncated)")
	assert.Less(t, len(prompt), 5000+500, "prompt must be capped, not include full 5000-char details")
}

func TestMeetsThreshold(t *testing.T) {
	assert.True(t, meetsThreshold("high", "high"))
	assert.True(t, meetsThreshold("high", "medium"))
	assert.True(t, meetsThreshold("medium", "medium"))
	assert.False(t, meetsThreshold("medium", "high"))
	assert.False(t, meetsThreshold("low", "medium"))
	assert.False(t, meetsThreshold("", "high"))
}

func TestApplier_NilAndEmpty(t *testing.T) {
	var a *OpenAIApplier
	findings := []state.Finding{{ID: "1", CVEID: "CVE-1", Details: "x"}}
	// nil applier is a no-op
	got := a.Apply(findings)
	assert.Equal(t, findings, got)

	// empty-endpoint applier: NewOpenAIApplier returns nil
	assert.Nil(t, NewOpenAIApplier("", "m", 0, 0, ""))
}

func TestApplier_SkipsFindingsWithoutDetails(t *testing.T) {
	// server should never be hit because the finding has no details/ranges
	hit := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit++
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	a := NewOpenAIApplier(srv.URL, "m", 0, 0, "high")
	require.NotNil(t, a)
	findings := []state.Finding{{ID: "1", CVEID: "CVE-1"}}
	got := a.Apply(findings)
	assert.Equal(t, findings, got)
	assert.Equal(t, 0, hit)
}

func TestApplier_SkipsAlreadyInformational(t *testing.T) {
	// server should never be hit for already-informational findings
	hit := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit++
	}))
	defer srv.Close()

	a := NewOpenAIApplier(srv.URL, "m", 0, 0, "high")
	findings := []state.Finding{{ID: "1", CVEID: "CVE-1", Details: "x", Class: "informational", ClassReason: "already-fixed"}}
	got := a.Apply(findings)
	assert.Equal(t, "informational", got[0].Class)
	assert.Equal(t, "already-fixed", got[0].ClassReason, "existing reason must be preserved")
	assert.Nil(t, got[0].AIApplicability, "must not attach applicability to already-informational finding")
	assert.Equal(t, 0, hit)
}

// applierResp wraps the model reply the fake server returns.
type applierResp struct {
	Applicable bool   `json:"applicable"`
	Confidence string `json:"confidence"`
	Reasoning  string `json:"reasoning"`
}

func fakeServer(t *testing.T, verdict applierResp) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		// sanity: the request was a forced-tool-call
		assert.Contains(t, string(body), applicabilityToolName)
		args, _ := json.Marshal(verdict)
		reply := map[string]any{
			"choices": []map[string]any{{
				"message": map[string]any{
					"tool_calls": []map[string]any{{
						"function": map[string]string{
							"name":      applicabilityToolName,
							"arguments": string(args),
						},
					}},
				},
			}},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(reply)
	}))
}

func TestApplier_AttachesMetadataOnHighConfidenceNotApplicable(t *testing.T) {
	srv := fakeServer(t, applierResp{Applicable: false, Confidence: "high", Reasoning: "queried version below introduced boundary"})
	defer srv.Close()
	a := NewOpenAIApplier(srv.URL, "m", 0, 0, "high")
	a.Today = "2026-07-08"

	findings := []state.Finding{{ID: "1", CVEID: "CVE-1", Package: "p", Details: "x"}}
	got := a.Apply(findings)
	require.Len(t, got, 1)
	assert.Equal(t, "", got[0].Class, "finding must stay actionable (fail-visible)")
	require.NotNil(t, got[0].AIApplicability)
	assert.False(t, got[0].AIApplicability.Applicable)
	assert.Equal(t, "high", got[0].AIApplicability.Confidence)
	assert.Contains(t, got[0].AIApplicability.Reasoning, "queried version below introduced boundary")
	assert.Equal(t, "m", got[0].AIApplicability.Model)
	assert.Equal(t, "2026-07-08", got[0].AIApplicability.CheckedAt)
}

func TestApplier_NoAttachmentWhenApplicable(t *testing.T) {
	srv := fakeServer(t, applierResp{Applicable: true, Confidence: "high", Reasoning: "vulnerable range"})
	defer srv.Close()
	a := NewOpenAIApplier(srv.URL, "m", 0, 0, "high")

	findings := []state.Finding{{ID: "1", CVEID: "CVE-1", Details: "x"}}
	got := a.Apply(findings)
	assert.Equal(t, "", got[0].Class)
	assert.Nil(t, got[0].AIApplicability, "applicable=true must not attach warning metadata")
}

func TestApplier_LowConfidenceNoAttachment(t *testing.T) {
	// model says not-applicable but at low confidence — must stay unannotated
	srv := fakeServer(t, applierResp{Applicable: false, Confidence: "low", Reasoning: "unsure"})
	defer srv.Close()
	a := NewOpenAIApplier(srv.URL, "m", 0, 0, "high")

	findings := []state.Finding{{ID: "1", CVEID: "CVE-1", Details: "x"}}
	got := a.Apply(findings)
	assert.Equal(t, "", got[0].Class)
	assert.Nil(t, got[0].AIApplicability, "low-confidence verdict must not attach a warning")
}

func TestApplier_MediumThresholdAcceptsMedium(t *testing.T) {
	srv := fakeServer(t, applierResp{Applicable: false, Confidence: "medium", Reasoning: "some reason"})
	defer srv.Close()
	a := NewOpenAIApplier(srv.URL, "m", 0, 0, "medium")

	findings := []state.Finding{{ID: "1", CVEID: "CVE-1", Details: "x"}}
	got := a.Apply(findings)
	assert.Equal(t, "", got[0].Class)
	require.NotNil(t, got[0].AIApplicability)
	assert.Equal(t, "medium", got[0].AIApplicability.Confidence)
}

func TestApplier_HTTPErrorLeavesUnannotated(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	a := NewOpenAIApplier(srv.URL, "m", 0, 0, "high")

	findings := []state.Finding{{ID: "1", CVEID: "CVE-1", Details: "x"}}
	got := a.Apply(findings)
	assert.Equal(t, "", got[0].Class, "HTTP 500 must not touch Class")
	assert.Nil(t, got[0].AIApplicability, "HTTP 500 must not attach a bogus verdict")
}

// TestApplier_MemoizesAcrossFindings verifies the (cve,pkg,cur,fix) cache: two
// findings sharing the same tuple must hit the model exactly once per run.
func TestApplier_MemoizesAcrossFindings(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		args, _ := json.Marshal(applierResp{Applicable: false, Confidence: "high", Reasoning: "ranges exclude"})
		reply := map[string]any{
			"choices": []map[string]any{{
				"message": map[string]any{
					"tool_calls": []map[string]any{{
						"function": map[string]string{"name": applicabilityToolName, "arguments": string(args)},
					}},
				},
			}},
		}
		_ = json.NewEncoder(w).Encode(reply)
	}))
	defer srv.Close()
	a := NewOpenAIApplier(srv.URL, "m", 0, 0, "high")
	findings := []state.Finding{
		{ID: "1", Repo: "o/a", CVEID: "CVE-1", Package: "p", CurrentVersion: "1.0", FixedVersion: "2.0", Details: "x"},
		{ID: "2", Repo: "o/b", CVEID: "CVE-1", Package: "p", CurrentVersion: "1.0", FixedVersion: "2.0", Details: "x"},
		{ID: "3", Repo: "o/c", CVEID: "CVE-2", Package: "p", CurrentVersion: "1.0", FixedVersion: "2.0", Details: "x"},
	}
	got := a.Apply(findings)
	assert.Equal(t, 2, calls, "same (cve,pkg,cur,fix) tuple must call the model once; different CVEs get their own call")
	for _, f := range got {
		require.NotNil(t, f.AIApplicability)
		assert.False(t, f.AIApplicability.Applicable)
	}
}
