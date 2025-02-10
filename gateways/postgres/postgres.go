package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sptGabriel/investment-analyzer/extensions/migrate"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/migrations"
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

func New(addr string, minConn, maxConn int32) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, fmt.Errorf("parsing pgxpool config: %w", err)
	}

	config.MaxConns = maxConn
	config.MinConns = minConn

	pgxConn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("creating new pgxpool: %w", err)
	}

	if err = migrate.RunMigrations(addr, migrate.Migrations{
		Folder: ".",
		FS:     migrations.MigrationFS,
	}); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return pgxConn, nil
}
