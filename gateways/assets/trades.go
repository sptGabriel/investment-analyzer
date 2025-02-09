package assets

import (
	_ "embed"

	"github.com/gocarina/gocsv"

	"github.com/sptGabriel/investment-analyzer/domain/ports"
)

//go:embed march_2021_trades.csv
var tradesFile []byte

func readTrades() ([]ports.TradeCSVDTO, error) {
	var trades []ports.TradeCSVDTO

	if err := gocsv.UnmarshalBytes(tradesFile, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}

// func readTrades(filename string) ([]ports.TradeCSVDTO, error) {
// 	file, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	var trades []ports.TradeCSVDTO
// 	err = gocsv.UnmarshalFile(file, &trades)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return trades, nil
// }
