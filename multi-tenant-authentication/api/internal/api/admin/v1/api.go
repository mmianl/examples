package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mmianl/api/v2/internal/api"
	"github.com/mmianl/api/v2/internal/config"
)

// AdminAPI is the API for the Admin endpoints
type AdminAPI struct {
	// router *httprouter.Router
}

// NewAPI Creates a new Admin API
func NewAPI(config *config.Config, router *httprouter.Router) api.Api {
	apiV1 := AdminAPI{
		// router: router,
	}

	// router.GET("/admin/v1/metrics", logging.LoggerMiddleware(apiV1.metricsHandler(), *config.Log))
	// router.GET("/admin/v1/ready", logging.LoggerMiddleware(apiV1.readyHandler(), *config.Log))
	// router.GET("/admin/v1/healthy", logging.LoggerMiddleware(apiV1.healthyHandler(), *config.Log))
	// router.GET("/v1/shutdown", logging.LoggerMiddleware(api.shutdownHandler(), *config.Log))

	return &apiV1
}
