package assets

import (
	_ "embed"
	"sort"

	"github.com/gocarina/gocsv"

	"github.com/sptGabriel/investment-analyzer/domain/ports"
)

//go:embed march_2021_pricesA.csv
var pricesAFile []byte

//go:embed march_2021_pricesB.csv
var pricesBFile []byte

func readPrices() ([]ports.PriceCSVDTO, error) {
	var pricesA []ports.PriceCSVDTO
	if err := gocsv.UnmarshalBytes(pricesAFile, &pricesA); err != nil {
		return nil, err
	}

	for i := range pricesA {
		pricesA[i].Symbol = "A"
	}

	var pricesB []ports.PriceCSVDTO
	if err := gocsv.UnmarshalBytes(pricesBFile, &pricesB); err != nil {
		return nil, err
	}

	for i := range pricesB {
		pricesB[i].Symbol = "B"
	}

	allPrices := append(pricesA, pricesB...)

	sort.Slice(allPrices, func(i, j int) bool {
		return allPrices[i].Time.Before(allPrices[j].Time.Time)
	})

	return allPrices, nil
}
