package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

// Gather sets the config with values read from the passed config file
func Gather[T any](configPath string, defaultConfig *T) (*T, error) {
	config := defaultConfig

	// Try to gather config from file but don't fail if it doesn't exist
	gatherFromFile(configPath, config)

	if err := gatherFromEnv(config); err != nil {
		return nil, err
	}

	return config, nil
}

func gatherFromFile[T any](configPath string, config *T) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(config); err != nil {
		return err
	}

	return nil
}

func gatherFromEnv[T any](config *T) error {
	if err := envconfig.Process("", config); err != nil {
		return err
	}

	return nil
}
