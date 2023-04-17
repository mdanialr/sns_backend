package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultHTTPErrorHandlerMsg(t *testing.T) {
	testCases := []struct {
		name   string
		sample int
		expect string
	}{
		{
			name:   "Given status code 404 should return Not Found",
			sample: 404,
			expect: "Not Found",
		},
		{
			name:   "Given status code 405 should return This method is not allowed here!",
			sample: 405,
			expect: "This method is not allowed here!",
		},
		{
			name:   "Given unmapped status code yet such as 401 should return Something was wrong!",
			sample: 401,
			expect: "Something was wrong!",
		},
		{
			name:   "Given unmapped status code yet such as 400 should return Something was wrong!",
			sample: 400,
			expect: "Something was wrong!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, defaultHTTPErrorHandlerMsg(tc.sample))
		})
	}
}
