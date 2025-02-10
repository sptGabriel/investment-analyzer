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
	PortfolioID string // should be ignored on challenge
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
	portfolio, err := uc.rebuildPortfolioSnapshot(ctx, input.PortfolioID, input.StartDate)
	if err != nil {
		return GenerateReportOutput{}, err
	}

	initialValue, err := uc.calculatePortfolioValue(ctx, portfolio, input.StartDate)
	if err != nil {
		return GenerateReportOutput{}, err
	}

	output := GenerateReportOutput{
		Reports: []ports.ReportDTO{
			{
				Timestamp:         input.StartDate,
				PortfolioValue:    initialValue.RoundUP(1),
				AccumulatedReturn: 0.0,
			},
		},
	}

	reports, err := uc.processTradesInInterval(ctx, input, portfolio, initialValue)
	if err != nil {
		return GenerateReportOutput{}, err
	}

	output.Reports = append(output.Reports, reports...)

	return output, nil
}

func (uc generateReportUC) calculateAccumulatedReturn(
	initialValue, currentValue vos.Decimal,
) float64 {
	if initialValue.IsZero() {
		return 0
	}

	accReturn := currentValue.Div(
		initialValue).Sub(vos.ParseToDecimal(1))

	return accReturn.RoundUP(6)
}

func (uc generateReportUC) rebuildPortfolioSnapshot(
	ctx context.Context, portfolioID string, startDate time.Time,
) (entities.Portfolio, error) {
	id, err := vos.ParseID(portfolioID)
	if err != nil {
		return entities.Portfolio{}, err
	}

	portfolio, err := uc.portfolioRepository.FindOne(ctx, id)
	if err != nil {
		return entities.Portfolio{}, err
	}

	snapShot, err := uc.tradesRepository.FindTradesBeforeDate(ctx, startDate)
	if err != nil {
		return entities.Portfolio{}, err
	}

	for _, trade := range snapShot {
		if err := portfolio.ApplyTrade(trade); err != nil {
			return entities.Portfolio{}, fmt.Errorf("on apply snapShot: %w", err)
		}
	}

	return portfolio, nil
}

func (uc generateReportUC) calculatePortfolioValue(
	ctx context.Context, portfolio entities.Portfolio, at time.Time,
) (vos.Decimal, error) {
	total := portfolio.Cash()

	for _, pos := range portfolio.Positions() {
		price, err := uc.pricesRepository.FindOneByTime(ctx, ports.FindOneByTimeInput{
			Date:    at,
			AssetID: pos.AssetID(),
		})
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				continue
			}
			return vos.Decimal{}, fmt.Errorf("%w: on find one by time", err)
		}

		positionValue := price.Value().Mul(vos.ParseToDecimalFromInt(pos.Quantity()))
		total = total.Add(positionValue)
	}

	return total, nil
}

func (uc generateReportUC) processTradesInInterval(
	ctx context.Context, input GenerateReportInput, portfolio entities.Portfolio, initialValue vos.Decimal,
) ([]ports.ReportDTO, error) {
	reports := make([]ports.ReportDTO, 0)

	rangeTrades, err := uc.tradesRepository.FindTradesByRange(ctx, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	currentTime := input.StartDate
	for currentTime.Before(input.EndDate) {
		windowEnd := currentTime.Add(input.Interval)
		if windowEnd.After(input.EndDate) {
			windowEnd = input.EndDate
		}

		intervalTrades := uc.filterTrades(rangeTrades, currentTime, windowEnd)
		for _, trade := range intervalTrades {
			portfolio.ApplyTrade(trade)
		}

		currentValue, err := uc.calculatePortfolioValue(ctx, portfolio, windowEnd)
		if err != nil {
			return nil, err
		}

		reports = append(reports, ports.ReportDTO{
			Timestamp:         windowEnd,
			PortfolioValue:    currentValue.Float64(),
			AccumulatedReturn: uc.calculateAccumulatedReturn(initialValue, currentValue),
		})

		currentTime = currentTime.Add(input.Interval)
	}

	return reports, nil
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
