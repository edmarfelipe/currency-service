package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
)

func SetDefault() {
	handler := tint.NewHandler(os.Stderr, &tint.Options{
		NoColor: isProd(),
	})
	slog.SetDefault(slog.New(NewRequestIDHandler(handler)))
}

func isProd() bool {
	return strings.ToLower(os.Getenv("ENV")) == "prd"
}
