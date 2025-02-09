package reports

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/ports"
	"github.com/sptGabriel/investment-analyzer/domain/ports/mocks"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	assetAID := vos.MustID("84b3cf08-db0c-46aa-9a3c-cb97c98337ec")
	assetBID := vos.MustID("b94b8afb-dc8a-450f-8ebc-83cabe3b3c3a")

	type priceKey struct {
		assetID entities.AssetID
		date    time.Time
	}

	var pricesMap = map[priceKey]entities.Price{
		{assetID: assetAID, date: time.Date(2021, 3, 1, 10, 10, 0, 0, time.UTC)}: entities.MustPrice(
			assetAID,
			vos.ParseToDecimal(23.92),
			time.Date(2021, 3, 1, 10, 10, 0, 0, time.UTC),
		),
		{assetID: assetAID, date: time.Date(2021, 3, 1, 10, 20, 0, 0, time.UTC)}: entities.MustPrice(
			assetAID,
			vos.ParseToDecimal(23.71),
			time.Date(2021, 3, 1, 10, 20, 0, 0, time.UTC),
		),
		{assetID: assetBID, date: time.Date(2021, 3, 1, 10, 10, 0, 0, time.UTC)}: entities.MustPrice(
			assetBID,
			vos.ParseToDecimal(23.82),
			time.Date(2021, 3, 1, 10, 10, 0, 0, time.UTC),
		),
		{assetID: assetBID, date: time.Date(2021, 3, 1, 10, 20, 0, 0, time.UTC)}: entities.MustPrice(
			assetBID,
			vos.ParseToDecimal(23.82),
			time.Date(2021, 3, 1, 10, 20, 0, 0, time.UTC),
		),
	}

	type args struct {
		ctx   context.Context
		input GenerateReportInput
	}

	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) generateReportUC
		want    []ports.ReportDTO
		wantErr error
	}{
		{
			name: "should return a report from March 1, 2021 10:00:00 to 10:20:00",
			args: args{
				ctx: context.TODO(),
				input: GenerateReportInput{
					PortfolioID: "02f168da-4ee6-4876-9821-5be6764e6eb1",
					StartDate:   time.Date(2021, 3, 1, 10, 0, 0, 0, time.UTC),
					EndDate:     time.Date(2021, 3, 1, 10, 20, 0, 0, time.UTC),
					Interval:    10 * time.Minute,
				},
			},
			setup: func(t *testing.T) generateReportUC {
				return generateReportUC{
					tradesRepository: &mocks.TradesRepositoryMock{
						FindTradesByRangeFunc: func(ctx context.Context, start, end time.Time) ([]entities.Trade, error) {
							return []entities.Trade{
								entities.MustTrade(
									vos.MustID(uuid.NewString()),
									entities.MustAsset(assetAID, "A"),
									vos.SideBuy,
									vos.ParseToDecimal(23.86),
									200,
									time.Date(2021, 3, 1, 10, 4, 23, 0, time.UTC),
								),
								entities.MustTrade(
									vos.MustID(uuid.NewString()),
									entities.MustAsset(assetBID, "B"),
									vos.SideBuy,
									vos.ParseToDecimal(23.79),
									100,
									time.Date(2021, 3, 1, 10, 5, 55, 0, time.UTC),
								),
								entities.MustTrade(
									vos.MustID(uuid.NewString()),
									entities.MustAsset(assetAID, "A"),
									vos.SideBuy,
									vos.ParseToDecimal(23.89),
									300,
									time.Date(2021, 3, 1, 10, 9, 20, 0, time.UTC),
								),
							}, nil
						},
					},
					pricesRepository: &mocks.PricesRepositoryMock{
						FindOneByTimeFunc: func(_ context.Context, input ports.FindOneByTimeInput) (entities.Price, error) {
							if price, ok := pricesMap[priceKey{
								assetID: input.AssetID,
								date:    input.Date,
							}]; ok {
								return price, nil
							}
							return entities.Price{}, errors.New("price not found")
						},
					},
					portfolioRepository: &mocks.PortfolioRepositoryMock{
						FindOneFunc: func(
							_ context.Context, id entities.PortfolioID,
						) (entities.Portfolio, error) {
							return entities.NewPortfolio(
								id,
								vos.ParseToDecimal(100000.0),
								vos.ParseToDecimal(100000.0),
								map[string]entities.Position{},
							)
						},
					},
				}
			},
			want: []ports.ReportDTO{
				{
					Timestamp:         time.Date(2021, 3, 1, 10, 0, 0, 0, time.UTC),
					PortfolioValue:    100000,
					AccumulatedReturn: 0.0,
				},
				{
					Timestamp:         time.Date(2021, 3, 1, 10, 10, 0, 0, time.UTC),
					PortfolioValue:    100024,
					AccumulatedReturn: 0.00024,
				},
				{
					Timestamp:         time.Date(2021, 3, 1, 10, 20, 0, 0, time.UTC),
					PortfolioValue:    99919,
					AccumulatedReturn: -0.00081,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := tt.setup(t).Execute(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, output.Reports, tt.want)
		})
	}
}
