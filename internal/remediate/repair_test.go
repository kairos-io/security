package remediate

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepairTask(t *testing.T) {
	task := RepairTask("./foo.go:10: undefined: Bar")
	assert.Contains(t, task, "undefined: Bar")
	assert.Contains(t, strings.ToLower(task), "go build")
	assert.Contains(t, strings.ToLower(task), "compile")
}

func TestConflictTask(t *testing.T) {
	assert.Contains(t, strings.ToLower(ConflictTask()), "conflict")
}
