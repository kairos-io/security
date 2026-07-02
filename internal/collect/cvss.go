package collect

import (
	"fmt"
	"math"
	"strings"
)

// CVSS v3.1 metric value tables, per the FIRST.org specification.
var (
	cvssAV = map[string]float64{"N": 0.85, "A": 0.62, "L": 0.55, "P": 0.2}
	cvssAC = map[string]float64{"L": 0.77, "H": 0.44}
	// Privileges Required depends on Scope.
	cvssPRUnchanged = map[string]float64{"N": 0.85, "L": 0.62, "H": 0.27}
	cvssPRChanged   = map[string]float64{"N": 0.85, "L": 0.68, "H": 0.5}
	cvssUI          = map[string]float64{"N": 0.85, "R": 0.62}
	cvssCIA         = map[string]float64{"H": 0.56, "L": 0.22, "N": 0}
)

// cvssV31BaseScore parses a CVSS v3.1 vector string (e.g.
// "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H") and computes its Base Score
// per the FIRST.org CVSS v3.1 specification. Returns an error if the vector
// is malformed or missing a required base metric.
func cvssV31BaseScore(vector string) (float64, error) {
	parts := strings.Split(strings.TrimSpace(vector), "/")
	if len(parts) == 0 || !strings.EqualFold(parts[0], "CVSS:3.1") {
		return 0, fmt.Errorf("cvss: not a CVSS:3.1 vector: %q", vector)
	}
	// Parse metric:value pairs, keyed by metric name. Unknown metrics
	// (temporal/environmental, e.g. E/RL/RC) are collected but ignored.
	metrics := map[string]string{}
	for _, p := range parts[1:] {
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, ":", 2)
		if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
			return 0, fmt.Errorf("cvss: malformed metric segment %q in %q", p, vector)
		}
		metrics[kv[0]] = kv[1]
	}

	scope, ok := metrics["S"]
	if !ok {
		return 0, fmt.Errorf("cvss: missing required metric S (Scope) in %q", vector)
	}
	if scope != "U" && scope != "C" {
		return 0, fmt.Errorf("cvss: invalid Scope value %q in %q", scope, vector)
	}

	prTable := cvssPRUnchanged
	if scope == "C" {
		prTable = cvssPRChanged
	}

	av, err := cvssLookup(metrics, "AV", cvssAV, vector)
	if err != nil {
		return 0, err
	}
	ac, err := cvssLookup(metrics, "AC", cvssAC, vector)
	if err != nil {
		return 0, err
	}
	pr, err := cvssLookup(metrics, "PR", prTable, vector)
	if err != nil {
		return 0, err
	}
	ui, err := cvssLookup(metrics, "UI", cvssUI, vector)
	if err != nil {
		return 0, err
	}
	c, err := cvssLookup(metrics, "C", cvssCIA, vector)
	if err != nil {
		return 0, err
	}
	i, err := cvssLookup(metrics, "I", cvssCIA, vector)
	if err != nil {
		return 0, err
	}
	a, err := cvssLookup(metrics, "A", cvssCIA, vector)
	if err != nil {
		return 0, err
	}

	iscBase := 1 - ((1 - c) * (1 - i) * (1 - a))

	var impact float64
	if scope == "U" {
		impact = 6.42 * iscBase
	} else {
		impact = 7.52*(iscBase-0.029) - 3.25*math.Pow(iscBase-0.02, 15)
	}

	exploitability := 8.22 * av * ac * pr * ui

	if impact <= 0 {
		return 0.0, nil
	}
	if scope == "U" {
		return cvssRoundup(math.Min(impact+exploitability, 10)), nil
	}
	return cvssRoundup(math.Min(1.08*(impact+exploitability), 10)), nil
}

// cvssLookup resolves one metric's numeric value from the metric table,
// erroring if the metric is absent or its value is not defined.
func cvssLookup(metrics map[string]string, name string, table map[string]float64, vector string) (float64, error) {
	raw, ok := metrics[name]
	if !ok {
		return 0, fmt.Errorf("cvss: missing required metric %s in %q", name, vector)
	}
	val, ok := table[raw]
	if !ok {
		return 0, fmt.Errorf("cvss: invalid value %q for metric %s in %q", raw, name, vector)
	}
	return val, nil
}

// cvssRoundup implements the official CVSS v3.1 rounding function, which avoids
// floating-point precision bugs at exact ".0" boundaries.
func cvssRoundup(input float64) float64 {
	intInput := int(math.Round(input * 100000))
	if intInput%10000 == 0 {
		return float64(intInput) / 100000.0
	}
	return (math.Floor(float64(intInput)/10000) + 1) / 10.0
}

// cvssSeverityLabel maps a CVSS v3.1 base score to ksec's severity enum
// (critical|high|medium|low), per the standard qualitative rating bands:
// 0.1-3.9 low, 4.0-6.9 medium, 7.0-8.9 high, 9.0-10.0 critical. A score of
// exactly 0.0 (no impact at all — essentially never occurs for a real
// published CVE) defensively maps to "low" rather than introducing a fifth
// enum value.
func cvssSeverityLabel(score float64) string {
	switch {
	case score >= 9.0:
		return "critical"
	case score >= 7.0:
		return "high"
	case score >= 4.0:
		return "medium"
	default:
		return "low"
	}
}
