package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

// readyHandler Returns a httprouter.Handle function that always returns a 200 response for readiness checks
func (a *ObservabilityAPI) readyHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if _, err := w.Write([]byte("Ready")); err != nil {
			log.Errorf("Could not write response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// healthyHandler Returns a httprouter.Handle function that always returns a 200 response for health checks
func (a *ObservabilityAPI) healthyHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if _, err := w.Write([]byte("Healthy")); err != nil {
			log.Errorf("Could not write response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// metricsHandler Returns a httprouter.Handle function that handles metrics requests
func (a *ObservabilityAPI) metricsHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		promhttp.Handler().ServeHTTP(w, r)
	}
}

// // shutdownHandler Returns a httprouter.Handle function that shuts down the server
// func (a UserAPI) shutdownHandler() httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 		a.Router.shutdownServer()
// 	}
// }
