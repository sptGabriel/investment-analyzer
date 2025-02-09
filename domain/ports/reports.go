package ports

import "time"

type ReportDTO struct {
	Timestamp         time.Time
	PortfolioValue    float64
	AccumulatedReturn float64
}
