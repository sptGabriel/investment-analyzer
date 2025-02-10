package testhelpers

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
)

func SaveAsset(t *testing.T, db *pgxpool.Pool, assets ...entities.Asset) {
	// language=PostgreSQL
	const query = `INSERT INTO assets(id, symbol) VALUES ($1, $2)`

	for _, a := range assets {
		res, err := db.Exec(
			context.Background(),
			query,
			a.ID().Value(),
			a.Symbol(),
		)
		require.NoError(t, err, "on save asset")

		require.NotEqual(t, 0, res.RowsAffected(), "not affected rows")
	}
}
