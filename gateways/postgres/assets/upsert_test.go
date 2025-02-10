package assets

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
	"github.com/sptGabriel/investment-analyzer/extensions/pgtest"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/testhelpers"
)

func TestUpsert(t *testing.T) {
	t.Parallel()

	type args struct {
		asset entities.Asset
	}

	tests := []struct {
		name    string
		args    args
		seed    func(t *testing.T, pool *pgxpool.Pool)
		wantErr error
	}{
		{
			name: "should upsert",
			seed: func(t *testing.T, pool *pgxpool.Pool) {
				assetID := vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf")
				saveAsset, err := entities.NewAsset(assetID, "bla")
				require.NoError(t, err)

				testhelpers.SaveAsset(t, pool, saveAsset)
			},
			args: args{
				asset: entities.MustAsset(vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"), "bla"),
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

			err := r.Upsert(testCtx, tt.args.asset)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
