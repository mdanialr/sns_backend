package middleware

import (
	"log"
	"os"
	"testing"
)

// JsonResponse standard response for all handler in this api pkg.
type JsonResponse struct {
	Msg string `json:"message"`
}

const fakeFilePath = "/tmp/fake-test-file.txt"

func TestMain(m *testing.M) {
	// prepare fake file for Send test
	err := os.WriteFile(fakeFilePath, []byte(`hello universe`), 0664)
	if err != nil {
		log.Fatalln("failed preparing fake file for testing Send middleware:", err)
	}

	os.Exit(m.Run())
}
