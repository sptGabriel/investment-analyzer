package testhelpers

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
)

func SaveTrades(t *testing.T, db *pgxpool.Pool, trades ...entities.Trade) {
	// language=PostgreSQL
	const query = `
		INSERT INTO trades(id, asset_id, side, price, quantity, timestamp) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	for _, trade := range trades {
		res, err := db.Exec(
			context.Background(),
			query,
			trade.ID().Value(),
			trade.Asset().ID().Value(),
			trade.Side().Value(),
			trade.Price().Float64(),
			trade.Quantity(),
			trade.Time(),
		)
		require.NoError(t, err, "on save trade")

		require.NotEqual(t, 0, res.RowsAffected(), "not affected rows")
	}
}
