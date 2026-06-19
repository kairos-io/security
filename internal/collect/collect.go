package collect

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

func defaultNow() string { return time.Now().UTC().Format("2006-01-02") }

// nowFn is overridable in tests.
var nowFn = defaultNow

func Today() string { return nowFn() }

func FindingID(repo, typ, cve, pkg string) string {
	sum := sha256.Sum256([]byte(repo + "|" + typ + "|" + cve + "|" + pkg))
	return hex.EncodeToString(sum[:])
}
