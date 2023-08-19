package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		defer func() {
			slog.InfoContext(r.Context(),
				"Request received",
				slog.String("path", r.URL.String()),
				slog.String("method", r.Method),
				slog.Int("status", ww.Status()),
				slog.Duration("duration", time.Since(start)),
			)
		}()
		next.ServeHTTP(ww, r)
	})
}
