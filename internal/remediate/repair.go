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

// ConflictTask is the prompt handed to the agent when a rebase of an owned PR
// branch onto the base branch produces git merge conflicts. The dependency
// version change from the branch must be preserved.
func ConflictTask() string {
	return "This branch has git merge conflicts after rebasing onto the base branch. " +
		"Resolve every conflict marker so the code is correct and `go build ./...` compiles, " +
		"keeping the dependency-version change from this branch."
}
