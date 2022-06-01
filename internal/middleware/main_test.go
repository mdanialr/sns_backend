package middleware

import (
	"os"
	"testing"
)

// JsonResponse standard response for all handler in this api pkg.
type JsonResponse struct {
	Msg string `json:"message"`
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
