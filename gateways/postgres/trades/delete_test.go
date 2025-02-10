package trades

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/sptGabriel/investment-analyzer/extensions/pgtest"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		seed    func(t *testing.T, pool *pgxpool.Pool)
		wantErr error
	}{
		{
			name:    "should delete",
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

			err := r.Delete(testCtx)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
