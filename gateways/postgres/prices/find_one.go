package prices

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/ports"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
)

func (r Repository) FindOneByTime(
	ctx context.Context, input ports.FindOneByTimeInput,
) (entities.Price, error) {
	const query = `select value from prices where asset_id = $1 and timestamp = $2;`

	var value float64
	if err := r.q.QueryRow(
		ctx, query, input.AssetID, input.Date,
	).Scan(
		&value,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Price{}, fmt.Errorf("%w:price not found", domain.ErrNotFound)
		}

		return entities.Price{}, err
	}

	p, err := entities.NewPrice(
		input.AssetID, vos.ParseToDecimal(value), input.Date.UTC(),
	)
	if err != nil {
		return entities.Price{}, err
	}

	return p, nil
}
