package entities

import (
	"fmt"
	"time"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
)

type ReportID = vos.ID

type Report struct {
	id                ReportID
	totalEquity       vos.Decimal
	accumulatedReturn float64
	timestamp         time.Time
}

func (r Report) ID() ReportID {
	return r.id
}

func (r Report) TotalEquity() vos.Decimal {
	return r.totalEquity
}

func (r Report) AccumulatedReturn() float64 {
	return r.accumulatedReturn
}

func (r Report) OccuredAt() time.Time {
	return r.timestamp
}

func NewReport(
	id ReportID,
	totalEquity vos.Decimal,
	aAccumulatedReturn float64,
	timeStamp time.Time,
) (Report, error) {
	if id.IsZero() {
		return Report{}, fmt.Errorf(
			"%w:invalid report_id", domain.ErrMalformedParameters)
	}

	if timeStamp.IsZero() {
		return Report{}, fmt.Errorf(
			"%w:empty timestamp value", domain.ErrMalformedParameters)
	}

	return Report{
		id:                id,
		totalEquity:       totalEquity,
		accumulatedReturn: aAccumulatedReturn,
		timestamp:         timeStamp,
	}, nil
}
