package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mmianl/api/v2/internal/api"
	"github.com/mmianl/api/v2/internal/config"
)

// UserAPI is the API for the User endpoints
type UserAPI struct {
	// router *httprouter.Router
}

// NewAPI Creates a new User API
func NewAPI(config *config.Config, router *httprouter.Router) api.Api {
	userV1 := UserAPI{
		// router: router,
	}

	// router.GET("/user/v1/metrics", logging.LoggerMiddleware(userV1.metricsHandler(), *config.Log))
	// router.GET("/user/v1/ready", logging.LoggerMiddleware(userV1.readyHandler(), *config.Log))
	// router.GET("/user/v1/healthy", logging.LoggerMiddleware(userV1.healthyHandler(), *config.Log))
	// router.GET("/shutdown", logging.LoggerMiddleware(userV1.shutdownHandler(), *config.Log))

	return &userV1
}
