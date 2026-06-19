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
