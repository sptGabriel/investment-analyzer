package entities

import (
	"errors"
	"fmt"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
)

var (
	ErrPositionNotFound = fmt.Errorf("%w:position not found in portfolio", domain.ErrFailedDependency)
)

type PortfolioID = vos.ID

type Portfolio struct {
	id          PortfolioID
	initialCash vos.Decimal
	cash        vos.Decimal
	positions   map[string]Position
}

func (p Portfolio) ID() PortfolioID {
	return p.id
}

func (p Portfolio) InitialCash() vos.Decimal {
	return p.initialCash
}

func (p Portfolio) Cash() vos.Decimal {
	return p.cash
}

func (p Portfolio) Positions() map[string]Position {
	return p.positions
}

func (p *Portfolio) ApplyTrade(t Trade) error {
	amount := t.Price().Mul(vos.ParseToDecimalFromInt(t.Quantity()))

	switch t.Side() {
	case vos.SideBuy:
		p.cash = p.cash.Sub(amount)

		pos, exists := p.positions[t.Asset().ID().Value()]
		if !exists {
			p.positions[t.Asset().ID().Value()] = Position{
				assetID:  t.Asset().ID(),
				quantity: t.Quantity(),
			}

			return nil
		}

		pos.quantity += t.Quantity()
		p.positions[t.Asset().ID().Value()] = pos

		return nil
	case vos.SideSell:
		pos, exists := p.positions[t.Asset().ID().Value()]
		if !exists {
			return ErrPositionNotFound
		}

		p.cash = p.cash.Add(amount)
		pos.quantity -= t.Quantity()
		p.positions[t.Asset().ID().Value()] = pos

		return nil
	default:
		return errors.New("invalid trade side")
	}
}

func NewPortfolio(
	id PortfolioID,
	initialCash vos.Decimal,
	cash vos.Decimal,
	positions map[string]Position,
) (Portfolio, error) {
	if id.IsZero() {
		return Portfolio{}, fmt.Errorf(
			"%w:invalid portfolio_id", domain.ErrMalformedParameters)
	}

	return Portfolio{
		id:          id,
		initialCash: initialCash,
		cash:        cash,
		positions:   positions,
	}, nil
}
