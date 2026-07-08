package state

// File name constants for each phase's output.
const (
	ReposFile      = "repos.json"
	FindingsFile   = "findings.json"
	CorrelatedFile = "correlated.json"
	TriageFile     = "triage.json"
	LedgerFile     = "ledger.json"
	OpenPRsFile    = "openprs.json"
	ReviewsFile    = "reviews.json"
)

type Artifact struct {
	Type    string `json:"type"`              // "image" | "go" | "component-manifest"
	Ref     string `json:"ref,omitempty"`     // image reference, when Type=="image"
	ModPath string `json:"modpath,omitempty"` // module path within repo, when Type=="go"
}

// ScanConfig holds per-repo scan opt-outs. A nil Source means source scanning
// is enabled (the default); set it to false to skip source scans (e.g. for a
// GTK4/cgo app that govulncheck cannot build headless).
type ScanConfig struct {
	Source *bool `json:"source,omitempty" yaml:"source"`
}

type Repo struct {
	Repo        string     `json:"repo"` // "owner/name"
	Kind        string     `json:"kind"` // "org" | "dep" | "external"
	Branch      string     `json:"branch"`
	Criticality string     `json:"criticality"` // "low" | "medium" | "high"
	Artifacts   []Artifact `json:"artifacts"`
	Scan        ScanConfig `json:"scan,omitempty" yaml:"scan"`
}

// SourceScanEnabled reports whether source scanning is enabled for the repo.
// It defaults to true when unset.
func (r Repo) SourceScanEnabled() bool { return r.Scan.Source == nil || *r.Scan.Source }

type Finding struct {
	ID             string `json:"id"` // stable dedupe key
	Repo           string `json:"repo"`
	Type           string `json:"type"` // "imageCVE" | "sourceCVE" | "ghAlert"
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
	FirstSeen      string `json:"firstSeen"`             // YYYY-MM-DD
	LastSeen       string `json:"lastSeen"`              // YYYY-MM-DD
	Class          string `json:"class,omitempty"`       // "" == actionable; "informational" == separated + uncounted
	ClassReason    string `json:"classReason,omitempty"` // why it's informational
}

type TrackedPR struct {
	Repo   string `json:"repo"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"url"`
	Source string `json:"source"` // renovate|dependabot|ksec|human
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

type LedgerEvent struct {
	Run    string `json:"run"`
	Action string `json:"action"`
	Detail string `json:"detail,omitempty"`
}

type LedgerEntry struct {
	Key           string        `json:"key"` // "<repo>|<package>"
	Repo          string        `json:"repo"`
	Package       string        `json:"package"`
	Branch        string        `json:"branch"`
	PRNumber      int           `json:"prNumber,omitempty"`
	PRURL         string        `json:"prURL,omitempty"`
	State         string        `json:"state"` // planned|open|merged|closed|conflicted|build-failed|error
	Bump          Bump          `json:"bump"`
	Severity      string        `json:"severity,omitempty"`
	Source        string        `json:"source,omitempty"`  // ksec | dependabot | renovate | human
	Kind          string        `json:"kind,omitempty"`    // direct | cascade | toolchain
	Blocked       string        `json:"blocked,omitempty"` // human-readable reason progress is stuck
	NeedsHuman    bool          `json:"needsHuman,omitempty"`
	CascadeFrom   string        `json:"cascadeFrom,omitempty"` // upstream ledger key that triggered this cascade bump
	PinTarget     string        `json:"pinTarget,omitempty"`   // for a pseudo cascade: the tag to re-pin to ("" while still pseudo)
	Supersedes    string        `json:"supersedes,omitempty"`  // the foreign PR URL this entry replaced
	Pseudo        bool          `json:"pseudo,omitempty"`      // true while the bump points at a pseudo-version awaiting re-pin
	CreatedRun    string        `json:"createdRun"`
	LastActionRun string        `json:"lastActionRun"`
	SeenComments  []string      `json:"seenComments,omitempty"` // reserved for Plan 3
	History       []LedgerEvent `json:"history,omitempty"`
}

type Ledger struct {
	Entries []LedgerEntry `json:"entries"`
}

type PRReview struct {
	Repo           string   `json:"repo"`
	PR             int      `json:"pr"`
	URL            string   `json:"url,omitempty"`
	HeadSHA        string   `json:"headSHA"`
	Verdict        string   `json:"verdict"` // good | bad | needs_human_verification
	Reasoning      string   `json:"reasoning,omitempty"`
	ChangesSummary string   `json:"changesSummary,omitempty"`
	ReviewedRun    string   `json:"reviewedRun,omitempty"`
	Trace          []string `json:"trace,omitempty"`
}
