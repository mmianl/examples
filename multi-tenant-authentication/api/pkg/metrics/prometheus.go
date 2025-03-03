package metrics

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsConfig struct {
	Port uint16 `yaml:"port" envconfig:"METRICS_PORT" required:"false"`
	Host string `yaml:"host" envconfig:"METRICS_HOST" required:"false"`
}

var DefaultMetricsConfig = &MetricsConfig{
	Host: "localhost",
	Port: 8080,
}

var (
	VersionGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "build_info",
			Help: "Metric with a constant '1' value labeled by version and goversion from which ga service was built.",
		},
		[]string{"service", "version", "commit", "goversion"},
	)
	StartTimeGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "start_time_seconds",
			Help: "Start time of the process since unix epoch in seconds.",
		},
		[]string{"service"},
	)
)

// Setup Sets up the exported metrics
func Setup(service, version, commit string) {
	VersionGauge.WithLabelValues(service, version, commit, runtime.Version()).Set(1)
	now := time.Now()
	StartTimeGauge.WithLabelValues(service).Set(float64(now.Unix()))
}
