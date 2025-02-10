package gbdb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

func (d Database) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

func (d Database) With(ctx context.Context, fn func(context.Context) error) error {
	return fn(withQuerier(ctx, d.pool))
}

func (d Database) Querier(ctx context.Context) Querier {
	return querierFromContext(ctx)
}

func (d Database) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return d.Querier(ctx).CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (d Database) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return d.Querier(ctx).Exec(ctx, sql, arguments...)
}

func (d Database) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return d.Querier(ctx).Query(ctx, sql, args...)
}

func (d Database) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return d.Querier(ctx).SendBatch(ctx, b)
}

func (d Database) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return d.Querier(ctx).QueryRow(ctx, sql, args...)
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{
		pool: pool,
	}
}
