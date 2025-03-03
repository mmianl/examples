package api

import (
	"context"
	"fmt"

	"net"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mmianl/api/v2/internal/config"
	log "github.com/sirupsen/logrus"
)

// HttpServer is the API server
type HttpServer struct {
	Router   *httprouter.Router
	shutdown <-chan os.Signal
	listen   string
}

// NewHttpServer creates a new API server
func NewHttpServer(config *config.Config, shutdown <-chan os.Signal) *HttpServer {
	router := httprouter.New()
	listen := fmt.Sprintf("%s:%d", config.Metrics.Host, config.Metrics.Port)

	NewObservabilityAPI(config, router)

	return &HttpServer{
		shutdown: shutdown,
		listen:   listen,
		Router:   router,
	}
}

func (s *HttpServer) Serve() error {
	ctx := context.Background()
	lc := &net.ListenConfig{}
	listenConfig, err := lc.Listen(ctx, "tcp", s.listen)
	if err != nil {
		log.Fatalf("Error listening on %s: %s", s.listen, err.Error())
	}

	httpServer := &http.Server{
		Addr:    s.listen,
		Handler: s.Router,
	}

	errCh := make(chan error)

	go func() {
		log.Infof("HTTP server listening on %s...", s.listen)
		errCh <- httpServer.Serve(listenConfig)
	}()

	select {
	case err := <-errCh:
		log.Fatalf("Error listening on %s: %s", s.listen, err.Error())
		return err
	case <-s.shutdown:
		log.Infof("Shutdown signal received.")
		deadline := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), deadline)
		defer cancel()
		log.Warnf("Shutting down in no more than %s.", deadline.String())
		return httpServer.Shutdown(ctx)
	}
}
