package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToHumanize(t *testing.T) {
	testCases := []struct {
		name   string
		sample int64
		expect string
	}{
		{
			name:   "Given 263 should return 263B",
			sample: 263,
			expect: "263B",
		},
		{
			name:   "Given 1234 should return 1.23KB",
			sample: 1234,
			expect: "1.23KB",
		},
		{
			name:   "Given 1238 should return 1.24KB",
			sample: 1238,
			expect: "1.24KB",
		},
		{
			name:   "Given 1234567 should return 1.23MB",
			sample: 1234567,
			expect: "1.23MB",
		},
		{
			name:   "Given 1235667 should return 1.24MB",
			sample: 1235667,
			expect: "1.24MB",
		},
		{
			name:   "Given 123456787 should return 123.46MB",
			sample: 123456787,
			expect: "123.46MB",
		},
		{
			name:   "Given 1234567890 should return 1.23GB",
			sample: 1234567890,
			expect: "1.23GB",
		},
		{
			name:   "Given 1234567890123 should return 1234.57GB",
			sample: 1234567890123,
			expect: "1234.57GB",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, BytesToHumanize(tc.sample))
		})
	}
}
