package assets

import (
	"os"

	"github.com/gocarina/gocsv"
)

type price struct {
	Time  customTime `csv:"time"`
	Price float64    `csv:"price"`
}

func readPrices(filename string) ([]price, error) {
	file, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var prices []price
	err = gocsv.UnmarshalFile(file, &prices)
	if err != nil {
		return nil, err
	}

	// // Ordenar preços por timestamp
	// sort.Slice(prices, func(i, j int) bool {
	// 	return prices[i].Time.Before(prices[j].Time.Time)
	// })

	return prices, nil
}
