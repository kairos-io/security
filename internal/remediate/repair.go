package remediate

import "fmt"

// RepairTask is the prompt handed to the agent when `go build` fails after a
// dependency bump. It must fix the code to compile WITHOUT changing the
// dependency versions (the security bump is deterministic and must stand).
func RepairTask(buildErr string) string {
	return fmt.Sprintf(`A dependency security bump in this repository broke the build.
Make the minimal source-code changes needed so that `+"`go build ./...`"+` compiles again.
Do NOT change any dependency versions in go.mod/go.sum — only adapt the calling code
(e.g. to a changed API). The compiler error was:

%s`, buildErr)
}
