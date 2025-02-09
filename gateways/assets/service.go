package assets

import (
	"github.com/sptGabriel/investment-analyzer/domain/ports"
)

func (s assetSvc) GetPrices() []ports.PriceCSVDTO {
	return s.prices
}

func (s assetSvc) GetTrades() []ports.TradeCSVDTO {
	return s.trades
}

func New() (assetSvc, error) {
	trades, err := readTrades()
	if err != nil {
		return assetSvc{}, err
	}

	prices, err := readPrices()
	if err != nil {
		return assetSvc{}, err
	}

	return assetSvc{
		trades: trades,
		prices: prices,
	}, nil
}

type assetSvc struct {
	prices []ports.PriceCSVDTO
	trades []ports.TradeCSVDTO
}
