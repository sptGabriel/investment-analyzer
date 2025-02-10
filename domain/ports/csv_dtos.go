package ports

import "time"

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalCSV(csv string) (err error) {
	ct.Time, err = time.Parse("2006-01-02 15:04:05", csv)
	return err
}

type PriceCSVDTO struct {
	Time   CustomTime `csv:"time"`
	Price  float64    `csv:"price"`
	Symbol string
}

type TradeCSVDTO struct {
	Time     CustomTime `csv:"time"`
	Symbol   string     `csv:"symbol"`
	Side     string     `csv:"side"`
	Price    float64    `csv:"price"`
	Quantity int        `csv:"quantity"`
}
