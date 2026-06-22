package collect

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassifyGovulncheck(t *testing.T) {
	// Success: no run error -> stdout passed through.
	out, err := ClassifyGovulncheck([]byte(`{"config":{}}`), nil, nil)
	require.NoError(t, err)
	assert.Equal(t, `{"config":{}}`, string(out))

	// Non-zero exit but stdout has an osv/finding record -> vulns found, normal.
	stdout := []byte(`{"config":{}}` + "\n" + `{"finding":{"osv":"GO-1"}}`)
	out, err = ClassifyGovulncheck(stdout, []byte("some progress"), errors.New("exit status 3"))
	require.NoError(t, err)
	assert.Equal(t, stdout, out)

	// Non-zero exit, only config/progress on stdout, build error on stderr -> real failure.
	_, err = ClassifyGovulncheck([]byte(`{"config":{}}`+"\n"+`{"progress":{}}`),
		[]byte("go: updates to go.mod needed; module requires go >= 1.26"), errors.New("exit status 1"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "go >= 1.26")
}
