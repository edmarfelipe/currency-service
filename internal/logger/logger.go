package logger

import (
	"context"
	"log/slog"
)

type loggerWithRequestID struct {
	slog.Handler
}

func NewRequestIDHandler(handler slog.Handler) slog.Handler {
	return loggerWithRequestID{
		handler,
	}
}

func (l loggerWithRequestID) Handle(ctx context.Context, record slog.Record) error {
	requestID := GetRequestID(ctx)
	if requestID != nil {
		record.Add("requestId", *requestID)
	}
	return l.Handler.Handle(ctx, record)
}
