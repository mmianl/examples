package logging

import (
	log "github.com/sirupsen/logrus"
)

type LogConfig struct {
	Level       string `yaml:"level" envconfig:"LOG_LEVEL" required:"false"`          // possible values: trace, debug, info, warn, error, fatal, panic
	Format      string `yaml:"format" envconfig:"LOG_FORMAT" required:"false"`        // possible values: logfmt, json
	LogRequests bool   `yaml:"logRequests" envconfig:"LOG_REQUESTS" required:"false"` // possible values: true, false
}

var DefaultLogConfig = &LogConfig{
	Level:       "info",
	Format:      "logfmt",
	LogRequests: true,
}

// InitLogger initializes the logger with the given configuration
func InitLogger(config *LogConfig) {
	level, err := log.ParseLevel(config.Level)
	if err != nil {
		level = log.InfoLevel
		log.Warnf("Invalid log level '%s', defaulting to 'info'", config.Level)
	}
	log.SetLevel(level)

	switch config.Format {
	case "logfmt":
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
		log.Warnf("Invalid log format '%s', defaulting to 'logfmt'", config.Format)
	}
}

// GetRequestLogger Returns a logger with requestId set
func GetRequestLogger(requestId string) *log.Entry {
	return log.WithField("requestId", requestId)
}
