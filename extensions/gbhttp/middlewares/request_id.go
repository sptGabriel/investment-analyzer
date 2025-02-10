package middlewares

import "context"

type requestIDCtxKey struct{}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey{}, requestID)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDCtxKey{}).(string)
	return requestID, ok
}
