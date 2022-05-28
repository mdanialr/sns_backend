package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	testCases := []struct {
		name   string
		sample uint
		expect int
	}{
		{
			name:   "Should pass and the string has 4 length",
			sample: 4, expect: 4,
		},
		{
			name:   "Should pass and the string has 55 length",
			sample: 55, expect: 55,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := RandomString(tc.sample)
			assert.Len(t, s, tc.expect)
		})
	}
}
