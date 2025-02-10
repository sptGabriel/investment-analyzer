package assets

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
	"github.com/sptGabriel/investment-analyzer/extensions/pgtest"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/testhelpers"
)

func TestFindOne(t *testing.T) {
	t.Parallel()

	type args struct {
		symbol string
	}

	tests := []struct {
		name    string
		args    args
		seed    func(t *testing.T, pool *pgxpool.Pool)
		want    entities.Asset
		wantErr error
	}{
		{
			name: "should get a asset by symbol",
			seed: func(t *testing.T, pool *pgxpool.Pool) {
				assetID := vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf")
				saveAsset, err := entities.NewAsset(assetID, "bla")
				require.NoError(t, err)

				testhelpers.SaveAsset(t, pool, saveAsset)
			},
			args: args{
				symbol: "bla",
			},
			want:    entities.MustAsset(vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"), "bla"),
			wantErr: nil,
		},
		{
			name: "should return err when not found",
			args: args{
				symbol: "bla",
			},
			want:    entities.Asset{},
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

			got, err := r.FindOneBySymbol(testCtx, tt.args.symbol)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
