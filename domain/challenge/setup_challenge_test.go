package challenge

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/ports"
)

func TestSetupChallenge(t *testing.T) {
	t.Parallel()

	errBla := errors.New("bla")

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) setupChallengeUC
		wantErr error
	}{
		{
			name: "should return err when call IsCsvImported",
			args: args{},
			setup: func(t *testing.T) setupChallengeUC {
				return setupChallengeUC{
					settingsRepository: &settingsRepositoryMock{
						IsCsvImportedFunc: func(contextMoqParam context.Context) (bool, error) {
							return false, errBla
						},
					},
				}
			},
			wantErr: errBla,
		},
		{
			name: "should return nill when has already imported",
			args: args{},
			setup: func(t *testing.T) setupChallengeUC {
				return setupChallengeUC{
					settingsRepository: &settingsRepositoryMock{
						IsCsvImportedFunc: func(contextMoqParam context.Context) (bool, error) {
							return true, nil
						},
					},
				}
			},
			wantErr: nil,
		},
		{
			name: "should setup challenge",
			args: args{},
			setup: func(t *testing.T) setupChallengeUC {
				return setupChallengeUC{
					portfolioRepository: &portfolioRepositoryMock{},
					settingsRepository: &settingsRepositoryMock{
						IsCsvImportedFunc: func(contextMoqParam context.Context) (bool, error) {
							return false, nil
						},
					},
					csvReader: &csvServicesMock{
						GetPricesFunc: func() []ports.PriceCSVDTO {
							return []ports.PriceCSVDTO{}
						},
						GetTradesFunc: func() []ports.TradeCSVDTO {
							return []ports.TradeCSVDTO{}
						},
					},
					assetsRepository: &assetsRepositoryMock{
						UpsertFunc: func(context.Context, entities.Asset) error {
							return nil
						},
					},
					pricesRepository: &pricesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
					},
					tradesRepository: &tradesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
					},
				}
			},
			wantErr: nil,
		},
		{
			name: "should return err when upsert asset a",
			args: args{},
			setup: func(t *testing.T) setupChallengeUC {
				return setupChallengeUC{
					settingsRepository: &settingsRepositoryMock{
						IsCsvImportedFunc: func(context.Context) (bool, error) {
							return false, nil
						},
					},
					assetsRepository: &assetsRepositoryMock{
						UpsertFunc: func(context.Context, entities.Asset) error {
							return errBla
						},
					},
				}
			},
			wantErr: errBla,
		},
		{
			name: "should return err when upsert asset b",
			args: args{},
			setup: func(t *testing.T) setupChallengeUC {
				return setupChallengeUC{
					settingsRepository: &settingsRepositoryMock{
						IsCsvImportedFunc: func(context.Context) (bool, error) {
							return false, nil
						},
					},
					assetsRepository: &assetsRepositoryMock{
						UpsertFunc: func(_ context.Context, a entities.Asset) error {
							if a.Symbol() == "A" {
								return nil
							}

							return errBla
						},
					},
				}
			},
			wantErr: errBla,
		},
		{
			name: "should return err when save prices batch",
			args: args{},
			setup: func(t *testing.T) setupChallengeUC {
				return setupChallengeUC{
					portfolioRepository: &portfolioRepositoryMock{},
					settingsRepository: &settingsRepositoryMock{
						IsCsvImportedFunc: func(contextMoqParam context.Context) (bool, error) {
							return false, nil
						},
					},
					csvReader: &csvServicesMock{
						GetPricesFunc: func() []ports.PriceCSVDTO {
							return []ports.PriceCSVDTO{
								{
									Symbol: "B",
									Time: ports.CustomTime{
										Time: time.Now(),
									},
								},
								{
									Symbol: "A",
									Time: ports.CustomTime{
										Time: time.Now(),
									},
								},
								{
									Symbol: "C",
								},
							}
						},
						GetTradesFunc: func() []ports.TradeCSVDTO {
							return []ports.TradeCSVDTO{}
						},
					},
					assetsRepository: &assetsRepositoryMock{
						UpsertFunc: func(context.Context, entities.Asset) error {
							return nil
						},
					},
					pricesRepository: &pricesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
						SavePricesBatchFunc: func(contextMoqParam context.Context, prices []entities.Price) error {
							return errBla
						},
					},
					tradesRepository: &tradesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
					},
				}
			},
			wantErr: errBla,
		},
		{
			name: "should return err when save trades batch",
			args: args{},
			setup: func(t *testing.T) setupChallengeUC {
				return setupChallengeUC{
					portfolioRepository: &portfolioRepositoryMock{},
					settingsRepository: &settingsRepositoryMock{
						IsCsvImportedFunc: func(contextMoqParam context.Context) (bool, error) {
							return false, nil
						},
					},
					csvReader: &csvServicesMock{
						GetPricesFunc: func() []ports.PriceCSVDTO {
							return []ports.PriceCSVDTO{
								{
									Symbol: "B",
									Time: ports.CustomTime{
										Time: time.Now(),
									},
								},
								{
									Symbol: "A",
									Time: ports.CustomTime{
										Time: time.Now(),
									},
								},
								{
									Symbol: "C",
								},
							}
						},
						GetTradesFunc: func() []ports.TradeCSVDTO {
							return []ports.TradeCSVDTO{
								{
									Time: ports.CustomTime{
										Time: time.Now(),
									},
									Symbol:   "B",
									Side:     "BUY",
									Price:    100,
									Quantity: 100,
								},
								{
									Time: ports.CustomTime{
										Time: time.Now(),
									},
									Symbol:   "A",
									Side:     "BUY",
									Price:    100,
									Quantity: 100,
								},
								{
									Time: ports.CustomTime{
										Time: time.Now(),
									},
									Symbol:   "C",
									Side:     "BUY",
									Price:    100,
									Quantity: 100,
								},
							}
						},
					},
					assetsRepository: &assetsRepositoryMock{
						UpsertFunc: func(context.Context, entities.Asset) error {
							return nil
						},
					},
					pricesRepository: &pricesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
						SavePricesBatchFunc: func(contextMoqParam context.Context, prices []entities.Price) error {
							return nil
						},
					},
					tradesRepository: &tradesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
						SaveTradesBatchFunc: func(context.Context, []entities.Trade) error {
							return errBla
						},
					},
				}
			},
			wantErr: errBla,
		},
		{
			name: "should return err when parse side",
			args: args{},
			setup: func(t *testing.T) setupChallengeUC {
				return setupChallengeUC{
					portfolioRepository: &portfolioRepositoryMock{},
					settingsRepository: &settingsRepositoryMock{
						IsCsvImportedFunc: func(contextMoqParam context.Context) (bool, error) {
							return false, nil
						},
					},
					csvReader: &csvServicesMock{
						GetPricesFunc: func() []ports.PriceCSVDTO {
							return []ports.PriceCSVDTO{
								{
									Symbol: "B",
									Time: ports.CustomTime{
										Time: time.Now(),
									},
								},
								{
									Symbol: "A",
									Time: ports.CustomTime{
										Time: time.Now(),
									},
								},
								{
									Symbol: "C",
								},
							}
						},
						GetTradesFunc: func() []ports.TradeCSVDTO {
							return []ports.TradeCSVDTO{
								{
									Time: ports.CustomTime{
										Time: time.Now(),
									},
									Symbol:   "B",
									Side:     "BUYX",
									Price:    100,
									Quantity: 100,
								},
								{
									Time: ports.CustomTime{
										Time: time.Now(),
									},
									Symbol:   "A",
									Side:     "BUY",
									Price:    100,
									Quantity: 100,
								},
								{
									Time: ports.CustomTime{
										Time: time.Now(),
									},
									Symbol:   "C",
									Side:     "BUY",
									Price:    100,
									Quantity: 100,
								},
							}
						},
					},
					assetsRepository: &assetsRepositoryMock{
						UpsertFunc: func(context.Context, entities.Asset) error {
							return nil
						},
					},
					pricesRepository: &pricesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
						SavePricesBatchFunc: func(contextMoqParam context.Context, prices []entities.Price) error {
							return nil
						},
					},
					tradesRepository: &tradesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
						SaveTradesBatchFunc: func(context.Context, []entities.Trade) error {
							return nil
						},
					},
				}
			},
			wantErr: domain.ErrFailedDependency,
		},
		{
			name: "should return err when upsert portfolio",
			args: args{},
			setup: func(t *testing.T) setupChallengeUC {
				return setupChallengeUC{
					portfolioRepository: &portfolioRepositoryMock{
						UpsertFunc: func(contextMoqParam context.Context, portfolio entities.Portfolio) error {
							return errBla
						},
					},
					settingsRepository: &settingsRepositoryMock{
						IsCsvImportedFunc: func(contextMoqParam context.Context) (bool, error) {
							return false, nil
						},
					},
					csvReader: &csvServicesMock{
						GetPricesFunc: func() []ports.PriceCSVDTO {
							return []ports.PriceCSVDTO{}
						},
						GetTradesFunc: func() []ports.TradeCSVDTO {
							return []ports.TradeCSVDTO{}
						},
					},
					assetsRepository: &assetsRepositoryMock{
						UpsertFunc: func(context.Context, entities.Asset) error {
							return nil
						},
					},
					pricesRepository: &pricesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
					},
					tradesRepository: &tradesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
					},
				}
			},
			wantErr: errBla,
		},
		{
			name: "should return err when set csv imported",
			args: args{},
			setup: func(t *testing.T) setupChallengeUC {
				return setupChallengeUC{
					portfolioRepository: &portfolioRepositoryMock{
						UpsertFunc: func(contextMoqParam context.Context, portfolio entities.Portfolio) error {
							return nil
						},
					},
					settingsRepository: &settingsRepositoryMock{
						IsCsvImportedFunc: func(contextMoqParam context.Context) (bool, error) {
							return false, nil
						},
						SetCsvImportedFunc: func(contextMoqParam context.Context) error {
							return errBla
						},
					},
					csvReader: &csvServicesMock{
						GetPricesFunc: func() []ports.PriceCSVDTO {
							return []ports.PriceCSVDTO{}
						},
						GetTradesFunc: func() []ports.TradeCSVDTO {
							return []ports.TradeCSVDTO{}
						},
					},
					assetsRepository: &assetsRepositoryMock{
						UpsertFunc: func(context.Context, entities.Asset) error {
							return nil
						},
					},
					pricesRepository: &pricesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
					},
					tradesRepository: &tradesRepositoryMock{
						DeleteFunc: func(contextMoqParam context.Context) error {
							return nil
						},
					},
				}
			},
			wantErr: errBla,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.setup(t).Execute(tt.args.ctx, SetupChallengeInput{})
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
