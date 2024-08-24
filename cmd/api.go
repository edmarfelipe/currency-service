package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.come/edmarfelipe/currency-service/internal"
	"github.come/edmarfelipe/currency-service/internal/httpserver"
	"github.come/edmarfelipe/currency-service/internal/logger"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to start", "err", err)
		os.Exit(1)
	}
}

func run() error {
	logger.SetDefault()

	cfg, err := internal.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading configs: %w", err)
	}

	ct, err := internal.NewContainer(cfg)
	if err != nil {
		return fmt.Errorf("error creating container: %w", err)
	}

	err = httpserver.New(ct).Start()
	if err != nil {
		return fmt.Errorf("error starting the http server: %w", err)
	}

	return nil
}
