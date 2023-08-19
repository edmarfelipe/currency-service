package logger

import (
	"context"
)

type requestIDCtxKeyType string

const requestIDCtxKey requestIDCtxKeyType = "request_id"

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey, requestID)
}

func GetRequestID(ctx context.Context) *string {
	requestID, ok := ctx.Value(requestIDCtxKey).(string)
	if !ok {
		return nil
	}
	return &requestID
}
