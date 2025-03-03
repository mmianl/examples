package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mmianl/api/v2/internal/config"
	"github.com/mmianl/api/v2/pkg/logging"
)

// Api is the interface for the Api
type Api interface {
	// Router() *httprouter.Router
}

// ObservabilityAPI is the API for the Observability endpoints
type ObservabilityAPI struct {
	// router *httprouter.Router
}

// NewObservabilityAPI Creates a new Observability API
func NewObservabilityAPI(config *config.Config, router *httprouter.Router) Api {
	obsAPI := ObservabilityAPI{
		// router: router,
	}

	router.GET("/metrics", logging.LoggerMiddleware(obsAPI.metricsHandler(), *config.Log))
	router.GET("/ready", logging.LoggerMiddleware(obsAPI.readyHandler(), *config.Log))
	router.GET("/healthy", logging.LoggerMiddleware(obsAPI.healthyHandler(), *config.Log))
	// router.GET("/v1/shutdown", logging.LoggerMiddleware(api.shutdownHandler(), *config.Log))

	return &obsAPI
}
