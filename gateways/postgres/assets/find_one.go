package assets

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
)

func (r Repository) FindOneBySymbol(
	ctx context.Context, symbol string,
) (entities.Asset, error) {
	const query = `SELECT id from assets where symbol = $1;`

	var (
		assetID string
	)

	if err := r.q.QueryRow(
		ctx, query, symbol,
	).Scan(
		&assetID,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Asset{}, fmt.Errorf("%w:asset not found", domain.ErrNotFound)
		}

		return entities.Asset{}, err
	}

	assetIDVO, err := vos.ParseID(assetID)
	if err != nil {
		return entities.Asset{}, err
	}

	p, err := entities.NewAsset(
		assetIDVO, symbol,
	)
	if err != nil {
		return entities.Asset{}, err
	}

	return p, nil
}
