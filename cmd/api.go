package main

import (
	"log/slog"

	"github.come/edmarfelipe/currency-service/internal"
	"github.come/edmarfelipe/currency-service/internal/httpserver"
	"github.come/edmarfelipe/currency-service/internal/logger"
)

func main() {
	logger.SetDefault()

	cfg, err := internal.LoadConfig()
	if err != nil {
		slog.Error("Error loading config: %v", err)
		return
	}

	ct, err := internal.NewContainer(cfg)
	if err != nil {
		slog.Error("Error creating container: %v", err)
		return
	}

	httpserver.New(ct).Start()
}
