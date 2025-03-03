package logging

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type contextKey string

const loggerKey contextKey = "logger"

// LoggerMiddleware Returns a middleware that adds a logger with a requestId field in the context
func LoggerMiddleware(next httprouter.Handle, config LogConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := uuid.New()
		ctx := r.Context()
		logger := GetRequestLogger(id.String())
		ctx = context.WithValue(ctx, loggerKey, logger)
		r = r.WithContext(ctx)

		if config.LogRequests {
			logger.WithFields(log.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"remote": r.RemoteAddr,
			}).Info("Request received")
		}

		next(w, r, p)
	}
}
