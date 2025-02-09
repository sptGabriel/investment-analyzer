package entities

import (
	"fmt"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
)

var (
	ErrInvalidAssetID = fmt.Errorf(
		"%w:invalid asset_id", domain.ErrMalformedParameters)
)

type AssetID = vos.ID

type Asset struct {
	id     AssetID
	symbol string
}

func (a Asset) ID() AssetID {
	return a.id
}

func (a Asset) Symbol() string {
	return a.symbol
}

func NewAsset(
	id AssetID,
	symbol string,
) (Asset, error) {
	if symbol == "" {
		return Asset{}, fmt.Errorf(
			"%w:empty symbol", domain.ErrMalformedParameters)
	}

	if id.IsZero() {
		return Asset{}, ErrInvalidAssetID
	}

	return Asset{
		id:     id,
		symbol: symbol,
	}, nil
}

func MustAsset(
	id AssetID,
	symbol string,
) Asset {
	return Asset{
		id:     id,
		symbol: symbol,
	}
}
