package reports

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/ports"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
	"github.com/sptGabriel/investment-analyzer/extensions/gbdb"
	"github.com/sptGabriel/investment-analyzer/extensions/gblib"
)

type GenerateReportInput struct {
	PortfolioID string
	StartDate   time.Time
	EndDate     time.Time
	Interval    time.Duration
}

type GenerateReportOutput struct {
	Reports []ports.ReportDTO
}

type generateReportUC struct {
	tradesRepository    ports.TradesRepository
	pricesRepository    ports.PricesRepository
	portfolioRepository ports.PortfolioRepository
}

func (uc generateReportUC) Execute(
	ctx context.Context, input GenerateReportInput,
) (GenerateReportOutput, error) {
	portfolioID, err := vos.ParseID(input.PortfolioID)
	if err != nil {
		return GenerateReportOutput{}, err
	}

	portfolio, err := uc.portfolioRepository.FindOne(
		ctx, portfolioID,
	)
	if err != nil {
		return GenerateReportOutput{}, err
	}

	output := GenerateReportOutput{
		Reports: []ports.ReportDTO{
			{
				Timestamp:         input.StartDate,
				PortfolioValue:    portfolio.InitialCash().Float64(),
				AccumulatedReturn: 0.0,
			},
		},
	}

	trades, err := uc.tradesRepository.FindTradesByRange(
		ctx, input.StartDate, input.EndDate,
	)
	if err != nil {
		return GenerateReportOutput{}, err
	}

	currentTime := input.StartDate
	for currentTime.Before(input.EndDate) {
		windowEnd := currentTime.Add(input.Interval)
		if windowEnd.After(input.EndDate) {
			windowEnd = input.EndDate
		}

		intervalTrades := uc.filterTrades(trades, currentTime, windowEnd)
		for _, trade := range intervalTrades {
			portfolio.ApplyTrade(trade)
		}

		currentValue := portfolio.Cash()
		for _, pos := range portfolio.Positions() {
			price, err := uc.pricesRepository.FindOneByTime(ctx, ports.FindOneByTimeInput{
				Date:    windowEnd,
				AssetID: pos.AssetID(),
			})
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					continue
				}

				return GenerateReportOutput{}, fmt.Errorf("%w:on find one by time", err)
			}

			positionValue := price.Value().Mul(vos.ParseToDecimalFromInt(pos.Quantity()))
			currentValue = currentValue.Add(positionValue)
		}

		accReturn := uc.calculateAccumulatedReturn(portfolio.InitialCash(), currentValue)

		output.Reports = append(output.Reports, ports.ReportDTO{
			Timestamp:         windowEnd,
			PortfolioValue:    currentValue.RoundBank(1),
			AccumulatedReturn: accReturn.RoundBank(5),
		})

		currentTime = currentTime.Add(input.Interval)
	}

	return output, nil
}

func (uc generateReportUC) calculateAccumulatedReturn(
	initialValue, currentValue vos.Decimal,
) vos.Decimal {
	if initialValue.IsZero() {
		return initialValue
	}

	return currentValue.Div(initialValue).Sub(vos.ParseToDecimal(1))
}

func (uc generateReportUC) filterTrades(
	trades []entities.Trade, start, end time.Time,
) []entities.Trade {
	filteredTrades := []entities.Trade{}
	for _, tr := range trades {
		if !tr.Time().Before(start) && tr.Time().Before(end) {
			filteredTrades = append(filteredTrades, tr)
		}
	}

	return filteredTrades
}

func NewGenerateReportUC(
	db *gbdb.Database,
	tradesRepository ports.TradesRepository,
	pricesRepository ports.PricesRepository,
	portfolioRepository ports.PortfolioRepository,
) gblib.UseCase[GenerateReportInput, GenerateReportOutput] {
	return gblib.New(
		generateReportUC{
			tradesRepository:    tradesRepository,
			pricesRepository:    pricesRepository,
			portfolioRepository: portfolioRepository,
		},
		gblib.WithDB(db),
	)
}
