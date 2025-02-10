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

func TestFindMany(t *testing.T) {
	t.Parallel()

	type args struct {
		start, end time.Time
	}

	mockDate := time.Date(2021, 3, 1, 10, 11, 0, 0, time.UTC)

	startDate := time.Date(2021, 3, 1, 10, 10, 0, 0, time.UTC)
	endDate := time.Date(2021, 3, 1, 10, 20, 0, 0, time.UTC)

	tests := []struct {
		name    string
		args    args
		seed    func(t *testing.T, pool *pgxpool.Pool)
		want    []entities.Trade
		wantErr error
	}{
		{
			name: "should get trades in range",
			seed: func(t *testing.T, pool *pgxpool.Pool) {
				assetID := vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf")
				saveAsset, err := entities.NewAsset(assetID, "bla")
				require.NoError(t, err)

				testhelpers.SaveAsset(t, pool, saveAsset)
				testhelpers.SaveTrades(
					t, pool,
					entities.MustTrade(
						vos.MustID("2278efd5-aa20-4ea5-a5fd-bf66142b87ae"),
						saveAsset, vos.SideBuy,
						vos.ParseToDecimal(999),
						99,
						mockDate,
					),
					entities.MustTrade(
						vos.MustID("68d1df59-004c-4326-917d-792845b032bc"),
						saveAsset, vos.SideBuy,
						vos.ParseToDecimal(1000),
						100,
						mockDate,
					),
				)
			},
			args: args{
				start: startDate,
				end:   endDate,
			},
			want: []entities.Trade{
				entities.MustTrade(
					vos.MustID("68d1df59-004c-4326-917d-792845b032bc"),
					entities.MustAsset(vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"), "bla"), vos.SideBuy,
					vos.ParseToDecimal(1000),
					100,
					mockDate,
				),
				entities.MustTrade(
					vos.MustID("2278efd5-aa20-4ea5-a5fd-bf66142b87ae"),
					entities.MustAsset(vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"), "bla"), vos.SideBuy,
					vos.ParseToDecimal(999),
					99,
					mockDate,
				),
			},
			wantErr: nil,
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

			got, err := r.FindTradesByRange(testCtx, tt.args.start, tt.args.end)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
