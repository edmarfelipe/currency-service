package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"time"

	chimetrics "github.com/edmarfelipe/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.come/edmarfelipe/currency-service/internal"
	"github.come/edmarfelipe/currency-service/usecase/convert"
)

// Server is the interface for the http server
type Server interface {
	// Start starts the http server
	Start()
	// TestServer returns a httptest.Server
	TestServer() *httptest.Server
}

type server struct {
	server http.Server
	ct     *internal.Container
}

// New creates a new http server
func New(ct *internal.Container) Server {
	return &server{
		server: http.Server{
			ReadTimeout:       1 * time.Second,
			WriteTimeout:      1 * time.Second,
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		},
		ct: ct,
	}
}

// Start starts the http server
func (s *server) Start() {
	runCtx, runCancel := context.WithCancel(context.Background())
	readyCtx, readyCancel := context.WithCancel(context.Background())

	s.server.Handler = s.createRouter(readyCtx)
	s.server.Addr = s.ct.Config.ServerAddr

	go s.waitForShutdown(runCancel, readyCancel)

	slog.Info("Starting the http server on " + s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Failed to start the http server", "err", err)
	}
	slog.Info("Http server stopped")
	<-runCtx.Done()
}

// createRouter creates the router for the http server
func (s *server) createRouter(readyCtx context.Context) *chi.Mux {
	r := chi.NewRouter()
	r.Use(requestIDMiddleware)
	r.Use(loggerMiddleware)
	r.Use(chimetrics.NewMiddleware("currency-service", 5, 50, 100, 200, 2000))
	r.Route("/api", func(r chi.Router) {
		r.Get("/ready", readyHandler(readyCtx))
		r.Handle("/metrics", promhttp.Handler())
		r.Get("/convert/{currency}/{value}", Adapt(convert.NewController(s.ct)))
	})
	return r
}

// TestServer returns a httptest.Server
func (s *server) TestServer() *httptest.Server {
	return httptest.NewServer(s.createRouter(context.Background()))
}

// waitForShutdown waits for a SIGTERM, SIGINT or SIGQUIT signal and then shuts down the server
func (s *server) waitForShutdown(runCancel context.CancelFunc, readyCancel context.CancelFunc) {
	trap := make(chan os.Signal, 1)

	signal.Notify(trap, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-trap
	slog.Info("Shutting down the server")
	readyCancel()
	<-time.After(2 * time.Second)

	err := s.server.Shutdown(context.Background())
	if err != nil {
		slog.Error("Error shutting down the server", "err", err)
	}

	s.ct.Shutdown()
	runCancel()
}
