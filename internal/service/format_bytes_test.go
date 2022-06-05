package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatBytesToHumanString(t *testing.T) {
	testCases := []struct {
		name   string
		sample int64
		expect string
	}{
		{
			name:   "Given 512 should return 512B",
			sample: 512,
			expect: "512B",
		},
		{
			name:   "Given 1023 should return 1023B",
			sample: 1023,
			expect: "1023B",
		},
		{
			name:   "Given 1024 should return 1.0KB",
			sample: 1024,
			expect: "1.0KB",
		},
		{
			name:   "Given 2048 should return 2.0KB",
			sample: 2048,
			expect: "2.0KB",
		},
		{
			name:   "Given 1048576 should return 1.0MB",
			sample: 1048576,
			expect: "1.0MB",
		},
		{
			name:   "Given 10485760 should return 10.0MB",
			sample: 10485760,
			expect: "10.0MB",
		},
		{
			name:   "Given 1073741824 should return 1.0GB",
			sample: 1073741824,
			expect: "1.0GB",
		},
		{
			name:   "Given 10737418240 should return 10.0GB",
			sample: 10737418240,
			expect: "10.0GB",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := FormatBytesToHumanString(tc.sample)
			assert.Equal(t, tc.expect, out)
		})
	}
}
