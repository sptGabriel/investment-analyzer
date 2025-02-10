package gblib

import (
	"context"

	"github.com/sptGabriel/investment-analyzer/extensions/gbdb"
)

func WithDB(database *gbdb.Database) Interceptor {
	return func(ctx context.Context, input interface{}, handler InterceptorFunc) (interface{}, error) {
		var result interface{}
		var err error

		database.With(ctx, func(ctx context.Context) error {
			result, err = handler(ctx, input)
			return err
		})
		return result, err
	}
}

func WithTx(tx gbdb.Transactioner) Interceptor {
	return func(ctx context.Context, input interface{}, handler InterceptorFunc) (interface{}, error) {
		var (
			result interface{}
			err    error
		)

		tx.With(ctx, func(txCtx context.Context) error {
			result, err = handler(txCtx, input)
			return err
		})

		return result, err
	}
}
