package gbdb

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Transactioner interface {
	With(ctx context.Context, fn func(context.Context) error) error
}

type transactioner struct {
	db *Database
}

func (t *transactioner) With(ctx context.Context, f func(context.Context) error) error {
	q := querierFromContext(ctx)

	_, hasTx := q.(pgx.Tx)
	if hasTx {
		return f(ctx)
	}

	nTx, err := t.db.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if rErr := recover(); rErr != nil {
			_ = nTx.Rollback(ctx)
			panic(err)
		}
	}()

	if err = f(withQuerier(ctx, nTx)); err != nil {
		if rollBackErr := nTx.Rollback(ctx); rollBackErr != nil {
			return fmt.Errorf("%w:%s", err, rollBackErr)
		}

		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	if err := nTx.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func NewTransactioner(db *Database) *transactioner {
	return &transactioner{db: db}
}
