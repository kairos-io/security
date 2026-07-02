package collect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCVSSV31BaseScore(t *testing.T) {
	tests := []struct {
		name   string
		vector string
		want   float64
	}{
		{
			// Widely-cited unauthenticated critical RCE vector.
			name:   "critical 9.8 scope unchanged",
			vector: "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			want:   9.8,
		},
		{
			// ALPINE-CVE-2025-66199's real vector: the regression case.
			name:   "medium 5.9 scope unchanged (regression anchor)",
			vector: "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H",
			want:   5.9,
		},
		{
			name:   "low 2.0 scope unchanged",
			vector: "CVSS:3.1/AV:N/AC:H/PR:H/UI:R/S:U/C:L/I:N/A:N",
			want:   2.0,
		},
		{
			name:   "critical 10.0 scope changed",
			vector: "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:C/C:H/I:H/A:H",
			want:   10.0,
		},
		{
			name:   "high 7.8 scope unchanged local",
			vector: "CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:H/I:H/A:H",
			want:   7.8,
		},
		{
			// Trailing temporal/environmental metrics must be ignored.
			name:   "ignores trailing non-base metrics",
			vector: "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H/E:P/RL:O/RC:C",
			want:   5.9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cvssV31BaseScore(tt.vector)
			require.NoError(t, err)
			assert.InDelta(t, tt.want, got, 1e-9)
		})
	}
}

func TestCVSSV31BaseScoreErrors(t *testing.T) {
	tests := []struct {
		name   string
		vector string
	}{
		{"malformed garbage", "not-a-cvss-vector"},
		{"missing required metric A", "CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:N/I:N"},
		{"bad metric value", "CVSS:3.1/AV:X/AC:H/PR:N/UI:N/S:U/C:N/I:N/A:H"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cvssV31BaseScore(tt.vector)
			assert.Error(t, err)
		})
	}
}

func TestCVSSSeverityLabel(t *testing.T) {
	tests := []struct {
		score float64
		want  string
	}{
		{0.0, "low"},
		{3.9, "low"},
		{4.0, "medium"},
		{6.9, "medium"},
		{7.0, "high"},
		{8.9, "high"},
		{9.0, "critical"},
		{10.0, "critical"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, cvssSeverityLabel(tt.score), "score %v", tt.score)
	}
}
