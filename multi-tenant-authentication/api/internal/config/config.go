package config

import (
	config "github.com/mmianl/api/v2/pkg/config"
	"github.com/mmianl/api/v2/pkg/logging"
	"github.com/mmianl/api/v2/pkg/metrics"
)

var global *Config
var defaultConfig = &Config{
	Log:     logging.DefaultLogConfig,
	Metrics: metrics.DefaultMetricsConfig,
}

type Config struct {
	Log     *logging.LogConfig     `yaml:"log"`
	Metrics *metrics.MetricsConfig `yaml:"metrics"`
}

// Global returns the globalConfig
func Global() *Config {
	return global
}

// Gather loads the configuration into the global variable
func Gather(configPath string) error {
	cfg, err := config.Gather(configPath, defaultConfig)
	if err != nil {
		return err
	}
	global = cfg
	return nil
}
