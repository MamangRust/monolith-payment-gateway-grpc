package middleware

import (
	"context"
	"fmt"
	"time"

	randomstring "github.com/MamangRust/monolith-payment-gateway-pkg/random_string"
)

type ctxKey string

const requestIDKey ctxKey = "requestId"

func GenerateRequestID() string {
	rand, _ := randomstring.GenerateRandomString(8)
	return fmt.Sprintf("req_%d_%s", time.Now().UnixNano(), rand)
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

// context.go

const startTimeKey ctxKey = "startTime"

func WithStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, startTimeKey, t)
}

func StartTime(ctx context.Context) time.Time {
	if v, ok := ctx.Value(startTimeKey).(time.Time); ok {
		return v
	}
	return time.Time{}
}

const (
	methodKey    ctxKey = "method"
	operationKey ctxKey = "operation"
)

func WithMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, methodKey, method)
}

func Method(ctx context.Context) string {
	if v, ok := ctx.Value(methodKey).(string); ok {
		return v
	}
	return ""
}

func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, operationKey, operation)
}

func Operation(ctx context.Context) string {
	if v, ok := ctx.Value(operationKey).(string); ok {
		return v
	}
	return ""
}
