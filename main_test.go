package main

import (
	"os"
	"testing"

	"github.com/mdanialr/sns_backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	testCases := []struct {
		name    string
		sample  *service.Config
		wantErr bool
	}{
		{
			name:   "Should pass using valid and accessible LogDir in config instance",
			sample: &service.Config{LogDir: "/tmp/", EnvIsProd: false},
		},
		{
			name:   "Should pass and has ProxyHeader as `X-Real-Ip` and DisableStartupMessage as `true`",
			sample: &service.Config{LogDir: "/tmp/", EnvIsProd: true},
		},
		{
			name:    "Should pass and has ProxyHeader as `X-Real-Ip` and DisableStartupMessage as `true`",
			sample:  &service.Config{LogDir: "/fake/path/log/"},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app, err := setup(tc.sample)
			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				conf := app.Config()
				if tc.sample.EnvIsProd {
					assert.Equal(t, "X-Real-Ip", conf.ProxyHeader)
					assert.Equal(t, true, conf.DisableStartupMessage)
				}
			case true:
				require.Error(t, err)
			}
		})
	}

	t.Cleanup(func() {
		os.Remove("/tmp/sns_backend-internal-log")
	})
}
