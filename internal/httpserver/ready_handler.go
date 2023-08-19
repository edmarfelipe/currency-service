package httpserver

import (
	"context"
	"net/http"
)

// readyHandler returns a 200 if the server is ready to receive requests
func readyHandler(readyCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if readyCtx.Err() != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
