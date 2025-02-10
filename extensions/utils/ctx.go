package utils

import "context"

type RequestData struct {
	IP        string
	Port      string
	RequestID string
}

type requestDataCtxKey struct{}

func WithRequestData(ctx context.Context, data RequestData) context.Context {
	return context.WithValue(ctx, requestDataCtxKey{}, data)
}

func RequestDataFromContext(ctx context.Context) (RequestData, bool) {
	data, ok := ctx.Value(requestDataCtxKey{}).(RequestData)
	return data, ok
}
