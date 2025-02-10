package scheme

import (
	"time"

	"github.com/sptGabriel/investment-analyzer/domain/reports"
)

type GenerateReportResponse struct {
	Timestamp         time.Time `json:"timestamp"`
	TotalEquity       float64   `json:"total_equity"`
	AccumulatedReturn float64   `json:"accumulated_return"`
}

func BuildGenerateReportResponse(output reports.GenerateReportOutput) []GenerateReportResponse {
	response := make([]GenerateReportResponse, 0, len(output.Reports))
	for _, r := range output.Reports {
		response = append(response, GenerateReportResponse{
			Timestamp:         r.Timestamp,
			TotalEquity:       r.PortfolioValue,
			AccumulatedReturn: r.AccumulatedReturn,
		})
	}

	return response
}
