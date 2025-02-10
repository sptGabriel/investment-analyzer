package assets

import (
	"context"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
)

func (r Repository) Upsert(
	ctx context.Context, a entities.Asset,
) error {
	const query = `
		INSERT INTO assets(id, symbol)
		VALUES ($1,	$2) ON CONFLICT (id) DO UPDATE
		SET
		id = $1,
		symbol = $2
		WHERE
			assets.id      	IS DISTINCT FROM $1 OR
			assets.symbol	IS DISTINCT FROM $2`

	_, err := r.q.Exec(
		ctx, query, a.ID().Value(), a.Symbol(),
	)
	if err != nil {
		return err
	}

	return nil
}
