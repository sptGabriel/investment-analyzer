package portfolios

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
	"github.com/sptGabriel/investment-analyzer/extensions/pgtest"
)

func TestUpsert(t *testing.T) {
	t.Parallel()

	type args struct {
		portfolio entities.Portfolio
	}

	tests := []struct {
		name    string
		args    args
		seed    func(t *testing.T, pool *pgxpool.Pool)
		wantErr error
	}{
		{
			name: "should upsert",
			args: args{
				portfolio: entities.MustPortfilio(
					vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"),
					vos.ParseToDecimal(100), vos.ParseToDecimal(100),
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

			err := r.Upsert(testCtx, tt.args.portfolio)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
