package state

// File name constants for each phase's output.
const (
	ReposFile      = "repos.json"
	FindingsFile   = "findings.json"
	CorrelatedFile = "correlated.json"
	TriageFile     = "triage.json"
)

type Artifact struct {
	Type    string `json:"type"`              // "image" | "go"
	Ref     string `json:"ref,omitempty"`     // image reference, when Type=="image"
	ModPath string `json:"modpath,omitempty"` // module path within repo, when Type=="go"
}

type Repo struct {
	Repo        string     `json:"repo"`        // "owner/name"
	Kind        string     `json:"kind"`        // "org" | "dep" | "external"
	Branch      string     `json:"branch"`
	Criticality string     `json:"criticality"` // "low" | "medium" | "high"
	Artifacts   []Artifact `json:"artifacts"`
}

type Finding struct {
	ID             string `json:"id"` // stable dedupe key
	Repo           string `json:"repo"`
	Type           string `json:"type"` // "pr" | "imageCVE" | "sourceCVE" | "ghAlert"
	CVEID          string `json:"cveID,omitempty"`
	GHSA           string `json:"ghsa,omitempty"`
	Ecosystem      string `json:"ecosystem,omitempty"`
	Package        string `json:"package,omitempty"`
	CurrentVersion string `json:"currentVersion,omitempty"`
	FixedVersion   string `json:"fixedVersion,omitempty"`
	Severity       string `json:"severity"` // critical|high|medium|low|unknown
	Source         string `json:"source"`   // tool/api that produced it
	Title          string `json:"title,omitempty"`
	URL            string `json:"url,omitempty"`
	FirstSeen      string `json:"firstSeen"` // YYYY-MM-DD
	LastSeen       string `json:"lastSeen"`  // YYYY-MM-DD
}

type CollectionError struct {
	Repo      string `json:"repo"`
	Collector string `json:"collector"`
	Message   string `json:"message"`
}

// Findings is the collect phase output: findings plus non-fatal errors.
type Findings struct {
	Findings []Finding         `json:"findings"`
	Errors   []CollectionError `json:"errors"`
}

type Bump struct {
	Package string `json:"package"`
	To      string `json:"to"`
}

type WaterfallGroup struct {
	ID            string   `json:"id"`
	RootCause     string   `json:"rootCause"`
	Ecosystem     string   `json:"ecosystem"`
	Severity      string   `json:"severity"`
	AffectedRepos []string `json:"affectedRepos"`
	SuggestedBump Bump     `json:"suggestedBump"`
}

type Correlated struct {
	Findings  []Finding        `json:"findings"`
	Waterfall []WaterfallGroup `json:"waterfall"`
}

type Triage struct {
	GeneratedAt string            `json:"generatedAt"`
	Model       string            `json:"model"`
	AIAvailable bool              `json:"aiAvailable"`
	Focus       []string          `json:"focus"`
	Summaries   map[string]string `json:"summaries"`
	Narrative   string            `json:"narrative"`
}
