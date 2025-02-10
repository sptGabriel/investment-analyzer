package prices

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
)

func (r Repository) SavePricesBatch(ctx context.Context, prices []entities.Price) error {
	if len(prices) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, p := range prices {
		batch.Queue(
			`INSERT INTO prices(asset_id, value, timestamp) VALUES ($1,	$2,	$3)`,
			p.AssetID().Value(),
			p.Value().Float64(),
			p.AtTime(),
		)
	}

	br := r.q.SendBatch(ctx, batch)
	defer br.Close()

	_, err := br.Exec()
	if err != nil {
		return err
	}

	return nil
}
