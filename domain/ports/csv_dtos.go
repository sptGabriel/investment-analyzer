package ports

import "time"

type customTime struct {
	time.Time
}

func (ct *customTime) UnmarshalCSV(csv string) (err error) {
	ct.Time, err = time.Parse("2006-01-02 15:04:05", csv)
	return err
}

type PriceCSVDTO struct {
	Time   customTime `csv:"time"`
	Price  float64    `csv:"price"`
	Symbol string
}

type TradeCSVDTO struct {
	Time     customTime `csv:"time"`
	Symbol   string     `csv:"symbol"`
	Side     string     `csv:"side"`
	Price    float64    `csv:"price"`
	Quantity int        `csv:"quantity"`
}
