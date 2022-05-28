package service

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tempConfigPath        = "/tmp/test-app.yaml"
	tempInvalidConfigPath = "/tmp/test-invalid-app.yaml"
)

func setupTestConfig(t *testing.T) {
	viper.SetConfigType("yaml")

	var configSample = []byte(`
env: dev
host: 127.0.0.1
port: 4747
log: /path/to/log/
db_driver: postgres
db_source: connection_string
`)

	if err := viper.ReadConfig(bytes.NewBuffer(configSample)); err != nil {
		require.NoError(t, err)
	}
	if err := viper.SafeWriteConfigAs(tempConfigPath); err != nil {
		require.NoError(t, err)
	}

	var invalidConfigSample = []byte(`
env: dev
host: 127.0.0.1
port:
  - 2323
log: /path/to/log/
db_drive: postgres
`)

	if err := viper.ReadConfig(bytes.NewBuffer(invalidConfigSample)); err != nil {
		require.NoError(t, err)
	}
	if err := viper.SafeWriteConfigAs(tempInvalidConfigPath); err != nil {
		require.NoError(t, err)
	}
}

func TestNewConfig(t *testing.T) {
	setupTestConfig(t)

	testCases := []struct {
		name       string
		sampleName string
		samplePath string
		wantErr    bool
	}{
		{
			name:       "Should pass when using valid and accessible filepath",
			sampleName: "test-app",
			samplePath: "/tmp/",
		},
		{
			name:       "Should fail when using invalid and or inaccessible filepath",
			sampleName: "test-example",
			samplePath: "/tmp/",
			wantErr:    true,
		},
		{
			name:       "Should fail when using invalid config structure even if the file is valid or accessible",
			sampleName: "test-invalid-app",
			samplePath: "/tmp/",
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewConfig(tc.sampleName, tc.samplePath)

			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, "dev", c.Env)
				assert.Equal(t, "127.0.0.1", c.Host)
				assert.Equal(t, "4747", c.PortNum)
				assert.Equal(t, "/path/to/log/", c.LogDir)
			case true:
				require.Error(t, err)
			}
		})
	}

	t.Cleanup(func() {
		os.Remove(tempConfigPath)
		os.Remove(tempInvalidConfigPath)
	})
}

func TestConfig_SanitizeEnv(t *testing.T) {
	testCases := []struct {
		name   string
		sample Config
		expect Config
	}{
		{
			name:   "Field `EnvIsProd` should be false if field `Env` is 'dev'",
			sample: Config{Env: "dev"},
			expect: Config{Env: "dev", EnvIsProd: false},
		}, {

			name:   "Field `EnvIsProd` should be true if field `Env` is 'prod'",
			sample: Config{Env: "prod"},
			expect: Config{Env: "prod", EnvIsProd: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.sample.SanitizeEnv()
			assert.Equal(t, tc.expect, tc.sample)
		})
	}
}

func TestConfig_SanitizeLog(t *testing.T) {
	testCases := []struct {
		name   string
		sample Config
		expect Config
	}{
		{
			name:   "Should has leading slash in LogDir",
			sample: Config{LogDir: "path/to/log/"},
			expect: Config{LogDir: "/path/to/log/"},
		},
		{
			name:   "Should has trailing slash in LogDir",
			sample: Config{LogDir: "/path/to/log"},
			expect: Config{LogDir: "/path/to/log/"},
		},
		{
			name:   "Should has leading and trailing slash in LogDir",
			sample: Config{LogDir: "path/to/log"},
			expect: Config{LogDir: "/path/to/log/"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.sample.SanitizeLog()
			assert.Equal(t, tc.expect, tc.sample)
		})
	}
}
