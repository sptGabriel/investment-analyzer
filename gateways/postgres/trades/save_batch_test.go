package trades

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
	"github.com/sptGabriel/investment-analyzer/extensions/pgtest"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/testhelpers"
)

func TestSaveBatch(t *testing.T) {
	t.Parallel()

	type args struct {
		trades []entities.Trade
	}
	mockDate := time.Date(2021, 3, 1, 10, 20, 0, 0, time.UTC)

	tests := []struct {
		name    string
		args    args
		seed    func(t *testing.T, pool *pgxpool.Pool)
		wantErr bool
	}{
		{
			name: "should save batch",
			seed: func(t *testing.T, pool *pgxpool.Pool) {
				assetID := vos.MustID("73f65f2e-08d1-41b8-af5e-832ee3780c75")
				saveAsset, err := entities.NewAsset(assetID, "bla")
				require.NoError(t, err)

				testhelpers.SaveAsset(t, pool, saveAsset)
				testhelpers.SaveTrades(
					t, pool,
					entities.MustTrade(
						vos.MustID("2278efd5-aa20-4ea5-a5fd-bf66142b87ae"),
						saveAsset, vos.SideBuy,
						vos.ParseToDecimal(1000),
						100,
						mockDate,
					),
				)
			},
			args: args{
				trades: []entities.Trade{
					entities.MustTrade(
						vos.MustID("2278efd5-aa20-4ea5-a5fd-bf66142b87ae"),
						entities.MustAsset(vos.MustID("73f65f2e-08d1-41b8-af5e-832ee3780c75"), "bla"),
						vos.SideBuy,
						vos.ParseToDecimal(1000),
						100,
						mockDate,
					),
				},
			},
			wantErr: true,
		},
		{
			name: "should save batch",
			seed: func(t *testing.T, pool *pgxpool.Pool) {
				assetID := vos.MustID("73f65f2e-08d1-41b8-af5e-832ee3780c75")
				saveAsset, err := entities.NewAsset(assetID, "bla")
				require.NoError(t, err)

				testhelpers.SaveAsset(t, pool, saveAsset)
			},
			args: args{
				trades: []entities.Trade{
					entities.MustTrade(
						vos.MustID("2278efd5-aa20-4ea5-a5fd-bf66142b87ae"),
						entities.MustAsset(vos.MustID("73f65f2e-08d1-41b8-af5e-832ee3780c75"), "bla"),
						vos.SideBuy,
						vos.ParseToDecimal(1000),
						100,
						mockDate,
					),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pgPool := pgtest.NewDB(t, t.Name())
			r := New(pgPool)
			if tt.seed != nil {
				tt.seed(t, pgPool)
			}

			err := r.SaveTradesBatch(testCtx, tt.args.trades)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
