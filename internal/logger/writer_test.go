package logger

import (
	"os"
	"testing"

	"github.com/mdanialr/sns_backend/internal/service"
	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	testCases := []struct {
		name    string
		sample  service.Config
		wantErr bool
	}{
		{
			name:   "Should pass using valid and accessible LogDir",
			sample: service.Config{LogDir: "/tmp/"},
		},
		{
			name:    "Should fail using invalid and or inaccessible LogDir",
			sample:  service.Config{LogDir: "/fake/path/"},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.sample.SanitizeDir()

			switch tt.wantErr {
			case false:
				require.NoError(t, InitLogger(&tt.sample))
			case true:
				require.Error(t, InitLogger(&tt.sample))
			}
		})
	}

	t.Cleanup(func() {
		os.Remove("/tmp/sns_backend-internal-log")
	})
}
