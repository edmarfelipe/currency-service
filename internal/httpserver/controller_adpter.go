package httpserver

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.come/edmarfelipe/currency-service/internal/xhttp"
)

var (
	errInternalError = xhttp.NewAPIError("Internal error", http.StatusInternalServerError)
)

func Adapt(ctrl xhttp.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ctrl.Handler(w, r); err != nil {
			if e, ok := err.(*xhttp.APIError); ok {
				writeError(w, e)
				return
			}
			slog.ErrorContext(r.Context(), "Failed to handle request", "err", err)
			writeError(w, errInternalError)
		}
	}
}

func writeError(w http.ResponseWriter, apiError *xhttp.APIError) {
	w.WriteHeader(apiError.Status())
	err := json.NewEncoder(w).Encode(apiError)
	if err != nil {
		slog.Error("Failed to write error response", "err", err)
	}
}
