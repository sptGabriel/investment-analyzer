package assets

import (
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type report struct {
	Timestamp              time.Time
	PatrimonioTotal        float64
	RentabilidadeAcumulada float64
}

type customTime struct {
	time.Time
}

func (ct *customTime) UnmarshalCSV(csv string) (err error) {
	ct.Time, err = time.Parse("2006-01-02 15:04:05", csv)
	return err
}

type trade struct {
	Time     customTime `csv:"time"`
	Symbol   string     `csv:"symbol"`
	Side     string     `csv:"side"`
	Price    float64    `csv:"price"`
	Quantity int        `csv:"quantity"`
}

func readTrades(filename string) ([]trade, error) {
	file, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var trades []trade
	err = gocsv.UnmarshalFile(file, &trades)
	if err != nil {
		return nil, err
	}

	return trades, nil
}
