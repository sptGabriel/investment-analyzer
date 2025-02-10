package prices

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/ports"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
	"github.com/sptGabriel/investment-analyzer/extensions/pgtest"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/testhelpers"
)

func TestFindOne(t *testing.T) {
	t.Parallel()

	type args struct {
		input ports.FindOneByTimeInput
	}

	mockDate := time.Date(2021, 3, 1, 10, 20, 0, 0, time.UTC)

	tests := []struct {
		name    string
		args    args
		seed    func(t *testing.T, pool *pgxpool.Pool)
		want    entities.Price
		wantErr error
	}{
		{
			name: "should get price",
			seed: func(t *testing.T, pool *pgxpool.Pool) {
				assetID := vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf")
				saveAsset, err := entities.NewAsset(assetID, "bla")
				require.NoError(t, err)

				testhelpers.SaveAsset(t, pool, saveAsset)

				testhelpers.SavePrices(
					t, pool,
					entities.MustPrice(
						vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"),
						vos.ParseToDecimal(1000),
						mockDate,
					))
			},
			args: args{
				input: ports.FindOneByTimeInput{
					Date:    mockDate,
					AssetID: vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"),
				},
			},
			want:    entities.MustPrice(vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"), vos.ParseToDecimal(1000), mockDate),
			wantErr: nil,
		},
		{
			name: "should return err when not found",
			args: args{
				input: ports.FindOneByTimeInput{},
			},
			want:    entities.Price{},
			wantErr: domain.ErrNotFound,
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

			got, err := r.FindOneByTime(testCtx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
