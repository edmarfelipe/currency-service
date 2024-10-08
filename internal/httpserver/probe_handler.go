package httpserver

import (
	"context"
	"net/http"
)

// probeHandler returns a 200 if the Server is ready to receive requests
func probeHandler(readyCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if readyCtx.Err() != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
