package httpserver

import (
	"net/http"

	"github.com/google/uuid"
	"github.come/edmarfelipe/currency-service/internal/logger"
)

// RequestIDHeader is the name of the HTTP Header which contains the request id.
var RequestIDHeader = "X-Request-Id"

func requestIDMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		next.ServeHTTP(w, r.WithContext(logger.WithRequestID(ctx, requestID)))
	}
	return http.HandlerFunc(fn)
}
