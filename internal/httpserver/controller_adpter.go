package httpserver

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.come/edmarfelipe/currency-service/internal"
)

var (
	ErrInternalError = internal.APIError{Message: "Internal error", Status: http.StatusInternalServerError}
)

type Controller interface {
	Handler(w http.ResponseWriter, r *http.Request) error
}

func adapt(ctrl Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ctrl.Handler(w, r); err != nil {
			var e internal.APIError
			if errors.As(err, &e) {
				writeError(w, e)
				return
			}
			slog.ErrorContext(r.Context(), "Failed to handle request", "err", err)
			writeError(w, ErrInternalError)
		}
	}
}

func writeError(w http.ResponseWriter, apiError internal.APIError) {
	w.WriteHeader(apiError.Status)
	err := json.NewEncoder(w).Encode(apiError)
	if err != nil {
		slog.Error("Failed to write error response", "err", err)
	}
}
