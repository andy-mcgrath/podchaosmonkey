package config

import (
	"io"
	"os"
	"time"

	"github.com/spf13/viper"
)

type (
	IConfig interface {
		GetEnvironment() string
		GetLogLevel() string
		GetLogOutput() io.Writer
		GetNamespace() string
		GetPodFilter() string
		GetConnectionTimeout() time.Duration
		GetKillTimeDelay() time.Duration
	}

	Config struct {
		Environment    string    `mapstructure:"ENVIRONMENT"`
		Namespace      string    `mapstructure:"NAMESPACE"`
		PodFilter      string    `mapstructure:"POD_FILTER"`
		RequestTimeout int64     `mapstructure:"DEFAULT_TIMEOUT"`
		KillTimeDelay  int64     `mapstructure:"KILL_TIME_DELAY"`
		LogLevel       string    `mapstructure:"LOG_LEVEL"`
		LogOutput      io.Writer `mapstructure:"LOG_OUTPUT"`
	}
)

func (c *Config) GetEnvironment() string  { return c.Environment }
func (c *Config) GetLogLevel() string     { return c.LogLevel }
func (c *Config) GetLogOutput() io.Writer { return c.LogOutput }
func (c *Config) GetNamespace() string    { return c.Namespace }
func (c *Config) GetPodFilter() string    { return c.PodFilter }
func (c *Config) GetKillTimeDelay() time.Duration {
	return time.Duration(c.KillTimeDelay) * time.Second
}
func (c *Config) GetConnectionTimeout() time.Duration {
	return time.Duration(c.RequestTimeout) * time.Second
}

func Init() (*Config, error) {
	cfg := &Config{}

	// Set configuration default values
	viper.SetDefault("ENVIRONMENT", "DEV")
	viper.SetDefault("KILL_TIME_DELAY", 120)
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_OUTPUT", os.Stdout)
	viper.SetDefault("NAMESPACE", "chaos")
	viper.SetDefault("POD_FILTER", ".*")
	viper.SetDefault("DEFAULT_TIMEOUT", 5)

	// Get configuration values from Environment variables and override from configuration file
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Select the correct configuration file `development.env` or `production.env`
	if viper.GetString("ENVIRONMENT") == "DEV" {
		viper.SetConfigName("development")
	} else {
		viper.SetConfigName("production")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
