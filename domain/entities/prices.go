package entities

import (
	"fmt"
	"time"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
)

type Price struct {
	assetID   AssetID
	value     vos.Decimal
	timestamp time.Time
}

func (p Price) AssetID() AssetID {
	return p.assetID
}

func (p Price) Value() vos.Decimal {
	return p.value
}

func (p Price) AtTime() time.Time {
	return p.timestamp
}

func NewPrice(
	assetID AssetID,
	value vos.Decimal,
	timestamp time.Time,
) (Price, error) {
	if assetID.IsZero() {
		return Price{}, ErrInvalidAssetID
	}

	if value.IsZero() {
		return Price{}, fmt.Errorf(
			"%w:empty price value", domain.ErrMalformedParameters)
	}

	if timestamp.IsZero() {
		return Price{}, fmt.Errorf(
			"%w:empty timestamp value", domain.ErrMalformedParameters)
	}

	return Price{
		value:     value,
		assetID:   assetID,
		timestamp: timestamp,
	}, nil
}

func MustPrice(
	assetID AssetID,
	value vos.Decimal,
	timestamp time.Time,
) Price {
	return Price{
		assetID:   assetID,
		value:     value,
		timestamp: timestamp,
	}
}
