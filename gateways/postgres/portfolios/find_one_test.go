package portfolios

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/sptGabriel/investment-analyzer/domain"
	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
	"github.com/sptGabriel/investment-analyzer/extensions/pgtest"
	"github.com/sptGabriel/investment-analyzer/gateways/postgres/testhelpers"
)

func TestFindOne(t *testing.T) {
	t.Parallel()

	type args struct {
		portfolioID entities.PortfolioID
	}

	tests := []struct {
		name    string
		args    args
		seed    func(t *testing.T, pool *pgxpool.Pool)
		want    entities.Portfolio
		wantErr error
	}{
		{
			name: "should find portfolio",
			seed: func(t *testing.T, pool *pgxpool.Pool) {
				portfolio := entities.MustPortfilio(
					vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"),
					vos.ParseToDecimal(100), vos.ParseToDecimal(100),
				)

				testhelpers.SavePortfolio(t, pool, portfolio)
			},
			args: args{
				portfolioID: vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"),
			},
			want:    entities.MustPortfilio(vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"), vos.ParseToDecimal(100), vos.ParseToDecimal(100)),
			wantErr: nil,
		},
		{
			name: "should return err when not found",
			args: args{
				portfolioID: vos.MustID("db4b1faf-292a-44a3-9eb0-6ef99d463fdf"),
			},
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

			got, err := r.FindOne(testCtx, tt.args.portfolioID)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
