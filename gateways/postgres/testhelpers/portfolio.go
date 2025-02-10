package testhelpers

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
)

func SavePortfolio(t *testing.T, db *pgxpool.Pool, assets ...entities.Portfolio) {
	// language=PostgreSQL
	const query = `INSERT INTO portfolios(id, initial_cash, cash) VALUES ($1, $2, $3)`

	for _, a := range assets {
		res, err := db.Exec(
			context.Background(),
			query,
			a.ID().Value(),
			a.InitialCash().Float64(),
			a.Cash().Float64(),
		)
		require.NoError(t, err, "on save portfolio")

		require.NotEqual(t, 0, res.RowsAffected(), "not affected rows")
	}
}
