package prices

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
		prices []entities.Price
	}

	tests := []struct {
		name    string
		args    args
		seed    func(t *testing.T, pool *pgxpool.Pool)
		wantErr error
	}{
		{
			name: "should save batch",
			seed: func(t *testing.T, pool *pgxpool.Pool) {
				assetID := vos.MustID("73f65f2e-08d1-41b8-af5e-832ee3780c75")
				saveAsset, err := entities.NewAsset(assetID, "bla")
				require.NoError(t, err)

				testhelpers.SaveAsset(t, pool, saveAsset)
			},
			args: args{
				prices: []entities.Price{
					entities.MustPrice(
						vos.MustID("73f65f2e-08d1-41b8-af5e-832ee3780c75"),
						vos.ParseToDecimal(1000),
						time.Now()),
				},
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

			err := r.SavePricesBatch(testCtx, tt.args.prices)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
