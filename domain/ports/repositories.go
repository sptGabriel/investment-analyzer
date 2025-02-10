package ports

import (
	"context"
	"time"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
)

//go:generate moq -stub -pkg mocks -out mocks/repositories.go . PricesRepository TradesRepository PortfolioRepository

type FindOneByTimeInput struct {
	Date    time.Time
	AssetID entities.AssetID
}

type PricesRepository interface {
	FindOneByTime(
		context.Context, FindOneByTimeInput,
	) (entities.Price, error)
}

type TradesRepository interface {
	FindTradesByRange(
		ctx context.Context, start, end time.Time,
	) ([]entities.Trade, error)
}

type PortfolioRepository interface {
	FindOne(
		context.Context, entities.PortfolioID,
	) (entities.Portfolio, error)
}
