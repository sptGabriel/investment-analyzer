package gbdb

import "context"

type ctxKey struct{}

var ctxConn = ctxKey{}

func withQuerier(ctx context.Context, q Querier) context.Context {
	return context.WithValue(ctx, ctxKey{}, q)
}

func querierFromContext(ctx context.Context) Querier {
	return ctx.Value(ctxConn).(Querier)
}
