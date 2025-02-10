package testhelpers

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
)

func SavePrices(t *testing.T, db *pgxpool.Pool, prices ...entities.Price) {
	// language=PostgreSQL
	const query = `INSERT INTO prices(asset_id, value, timestamp) VALUES ($1,	$2,	$3)`

	for _, p := range prices {
		res, err := db.Exec(
			context.Background(),
			query,
			p.AssetID().Value(),
			p.Value().Float64(),
			p.AtTime(),
		)
		require.NoError(t, err, "on save asset")

		require.NotEqual(t, 0, res.RowsAffected(), "not affected rows")
	}
}
