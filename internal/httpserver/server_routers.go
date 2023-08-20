package httpserver

import (
	"context"

	chimetrics "github.com/edmarfelipe/chi-prometheus"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.come/edmarfelipe/currency-service/usecase/convert"
)

// createRouter creates the router for the http server
func (s *server) createRouter(readyCtx context.Context) *chi.Mux {
	r := chi.NewRouter()
	r.Use(requestIDMiddleware)
	r.Use(loggerMiddleware)
	r.Use(chimetrics.NewMiddleware("currency-service", 5, 50, 100, 200, 2000))
	r.Route("/api", func(r chi.Router) {
		r.Get("/ready", probeHandler(readyCtx))
		r.Handle("/metrics", promhttp.Handler())
		r.Get("/convert/{currency}/{value}", adapt(convert.NewController(s.ct)))
	})
	return r
}
