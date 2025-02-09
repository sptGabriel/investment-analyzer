package portfolios

import (
	"context"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
)

func (r Repository) Upsert(
	ctx context.Context, p entities.Portfolio,
) error {
	const query = `
		INSERT INTO portfolios(id, initial_cash, cash)
		VALUES ($1,	$2, $3) ON CONFLICT (id) DO UPDATE
		SET
		id = $1,
		initial_cash = $2,
		cash = $3
		WHERE
			portfolios.id      		IS DISTINCT FROM $1 OR
			portfolios.initial_cash	IS DISTINCT FROM $2 OR
			portfolios.cash			IS DISTINCT FROM $3`

	_, err := r.q.Exec(
		ctx, query, p.ID().Value(), p.InitialCash().Float64(), p.Cash().Float64(),
	)
	if err != nil {
		return err
	}

	return nil
}
