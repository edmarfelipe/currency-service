package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.come/edmarfelipe/currency-service/internal"
)

type Server struct {
	server http.Server
	ct     *internal.Container
}

// New creates a new http Server
func New(ct *internal.Container) *Server {
	return &Server{
		ct: ct,
		server: http.Server{
			ReadTimeout:       1 * time.Second,
			WriteTimeout:      1 * time.Second,
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		},
	}
}

// Start starts the http Server
func (s *Server) Start() error {
	runCtx, runCancel := context.WithCancel(context.Background())
	readyCtx, readyCancel := context.WithCancel(context.Background())

	s.server.Addr = s.ct.Config.ServerAddr
	s.server.Handler = s.createRouter(readyCtx)

	go s.waitForShutdown(runCancel, readyCancel)

	slog.Info("Starting the http Server on " + s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-runCtx.Done()
	return nil
}

// TestServer returns a httptest.Server
func (s *Server) TestServer(t *testing.T) *httptest.Server {
	serv := httptest.NewServer(s.createRouter(context.Background()))
	t.Cleanup(func() {
		t.Helper()
		serv.Close()
	})
	return serv
}

// waitForShutdown waits for a SIGTERM, SIGINT or SIGQUIT signal and then shuts down the Server
func (s *Server) waitForShutdown(runCancel context.CancelFunc, readyCancel context.CancelFunc) {
	trap := make(chan os.Signal, 1)

	signal.Notify(trap, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-trap
	slog.Info("Shutting down the Server")
	readyCancel()
	<-time.After(2 * time.Second)

	err := s.server.Shutdown(context.Background())
	if err != nil {
		slog.Error("Error shutting down the Server", "err", err)
	}

	s.ct.Shutdown()
	runCancel()
}
