package trades

import (
	"context"
	"time"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
)

func (r Repository) FindTradesBeforeDate(
	ctx context.Context, dateTime time.Time,
) ([]entities.Trade, error) {
	const query = `
		SELECT 
			t.id, 
			t.asset_id, 
			a.symbol, 
			t.side, 
			t.price, 
			t.quantity, 
			t.timestamp
		FROM trades t
		JOIN assets a ON t.asset_id = a.id
		WHERE t.timestamp < $1;
	`

	rows, err := r.q.Query(ctx, query, dateTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trades := make([]entities.Trade, 0)
	for rows.Next() {
		var (
			id, assetID, symbol, side string
			price                     float64
			quantity                  int
			timestamp                 time.Time
		)

		if err := rows.Scan(
			&id, &assetID, &symbol, &side, &price, &quantity, &timestamp,
		); err != nil {
			return nil, err
		}

		assetIDVo, err := vos.ParseID(assetID)
		if err != nil {
			return nil, err
		}

		asset, err := entities.NewAsset(assetIDVo, symbol)
		if err != nil {
			return nil, err
		}

		tradeID, err := vos.ParseID(id)
		if err != nil {
			return nil, err
		}

		sideVO, err := vos.ParseSide(side)
		if err != nil {
			return nil, err
		}

		trade, err := entities.NewTrade(
			tradeID,
			asset,
			sideVO,
			vos.ParseToDecimal(price),
			quantity,
			timestamp.UTC(),
		)
		if err != nil {
			return nil, err
		}

		trades = append(trades, trade)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trades, nil
}
