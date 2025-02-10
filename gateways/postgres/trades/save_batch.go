package trades

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
)

func (r Repository) SaveTradesBatch(ctx context.Context, trades []entities.Trade) error {
	if len(trades) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, trade := range trades {
		batch.Queue(
			`INSERT INTO trades(id, asset_id, side, price, quantity, timestamp) VALUES ($1, $2, $3, $4, $5, $6)`,
			trade.ID().Value(),
			trade.Asset().ID().Value(),
			trade.Side().Value(),
			trade.Price().Float64(),
			trade.Quantity(),
			trade.Time(),
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
