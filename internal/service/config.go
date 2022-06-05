package service

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config Stores all configuration of the application.
// The values are read by viper from a config file or env variables.
type Config struct {
	EnvIsProd bool
	Env       string     `mapstructure:"env"`
	Host      string     `mapstructure:"host"`
	PortNum   string     `mapstructure:"port"`
	LogDir    string     `mapstructure:"log"`
	UploadDir string     `mapstructure:"upload"`
	DB        DBPostgres `mapstructure:"db"`
}

// NewConfig create new Config instance by reading from file or env variable.
func NewConfig(name, path string) (conf *Config, err error) {
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("failed to read config: %s", err)
		return
	}
	if err = viper.Unmarshal(&conf); err != nil {
		err = fmt.Errorf("failed to unmarshalling config: %s", err)
		return
	}

	return
}

// SanitizeEnv set whether EnvIsProd should be true or false based on the value of Env.
// If Env is 'dev' then it should be false or If Env is 'prod' then it should be true.
func (c *Config) SanitizeEnv() {
	switch c.Env {
	case "dev":
		c.EnvIsProd = false
	case "prod":
		c.EnvIsProd = true
	}
}

// SanitizeDir make sure LogDir & UploadDir has leading and trailing slash, so it can be safely used in another place.
func (c *Config) SanitizeDir() {
	if !strings.HasPrefix(c.LogDir, "/") {
		c.LogDir = "/" + c.LogDir
	}
	if !strings.HasSuffix(c.LogDir, "/") {
		c.LogDir += "/"
	}
	if !strings.HasPrefix(c.UploadDir, "/") {
		c.UploadDir = "/" + c.UploadDir
	}
	if !strings.HasSuffix(c.UploadDir, "/") {
		c.UploadDir += "/"
	}
}
