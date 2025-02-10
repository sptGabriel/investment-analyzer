package entities

import (
	"fmt"
	"time"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
)

type TradeID = vos.ID

type Trade struct {
	id        TradeID
	asset     Asset
	side      vos.Side
	price     vos.Decimal
	quantity  int
	timeStamp time.Time
}

func (t Trade) ID() TradeID {
	return t.id
}

func (t Trade) Asset() Asset {
	return t.asset
}

func (t Trade) Side() vos.Side {
	return t.side
}

func (t Trade) Price() vos.Decimal {
	return t.price
}

func (t Trade) Quantity() int {
	return t.quantity
}

func (t Trade) Time() time.Time {
	return t.timeStamp
}

func NewTrade(
	id TradeID,
	asset Asset,
	side vos.Side,
	price vos.Decimal,
	quantity int,
	timeStamp time.Time,
) (Trade, error) {
	if id.IsZero() {
		return Trade{}, fmt.Errorf(
			"%w:invalid trade_id", domain.ErrMalformedParameters)
	}

	if side.IsZero() {
		return Trade{}, fmt.Errorf(
			"%w:invalid side", domain.ErrMalformedParameters)
	}

	if timeStamp.IsZero() {
		return Trade{}, fmt.Errorf(
			"%w:empty timestamp value", domain.ErrMalformedParameters)
	}

	if quantity == 0 {
		return Trade{}, fmt.Errorf(
			"%w: trade quantity must be non-zero", domain.ErrMalformedParameters)
	}

	if price.IsZero() {
		return Trade{}, fmt.Errorf(
			"%w: trade price must be non-zero", domain.ErrMalformedParameters)
	}

	return Trade{
		id:        id,
		asset:     asset,
		side:      side,
		price:     price,
		quantity:  quantity,
		timeStamp: timeStamp,
	}, nil
}

func MustTrade(
	id TradeID,
	asset Asset,
	side vos.Side,
	price vos.Decimal,
	quantity int,
	timeStamp time.Time,
) Trade {
	return Trade{
		id:        id,
		asset:     asset,
		side:      side,
		price:     price,
		quantity:  quantity,
		timeStamp: timeStamp,
	}
}
