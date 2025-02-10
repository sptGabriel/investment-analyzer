package vos

import "github.com/shopspring/decimal"

type Decimal struct {
	value decimal.Decimal
}

func (c Decimal) Add(number Decimal) Decimal {
	return Decimal{value: c.value.Add(number.value).RoundBank(8)}
}

func (c Decimal) Sub(number Decimal) Decimal {
	return Decimal{value: c.value.Sub(number.value).RoundBank(8)}
}

func (c Decimal) Mul(number Decimal) Decimal {
	return Decimal{value: c.value.Mul(number.value).RoundBank(8)}
}

func (c Decimal) Div(number Decimal) Decimal {
	return Decimal{value: c.value.Div(number.value).RoundBank(8)}
}

func (c Decimal) Float64() float64 {
	return c.value.InexactFloat64()
}

func (c Decimal) RoundUP(value int) float64 {
	return c.value.RoundUp(int32(value)).InexactFloat64()
}

func (c Decimal) IsZero() bool {
	return c.value.IsZero()
}

func NewDecimal() Decimal {
	return Decimal{}
}

func ParseToDecimal(value float64) Decimal {
	parsedValue := decimal.NewFromFloat(value)
	return Decimal{parsedValue}
}

func ParseToDecimalFromInt(value int) Decimal {
	parsedValue := decimal.NewFromInt(int64(value))
	return Decimal{parsedValue}
}
