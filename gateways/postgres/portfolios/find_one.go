package portfolios

import (
	"context"
	"database/sql"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
	"github.com/sptGabriel/investment-analyzer/domain/vos"
)

func (r Repository) FindOne(
	ctx context.Context, portfolioID entities.PortfolioID,
) (entities.Portfolio, error) {
	const query = `
		SELECT 
			p.initial_cash, 
			p.cash, 
			pos.asset_id, 
			pos.quantity
		FROM portfolios p
		LEFT JOIN positions pos ON p.id = pos.portfolio_id
		WHERE p.id = $1;`

	var (
		initialCash, cash float64
	)

	rows, err := r.q.Query(ctx, query, portfolioID.Value())
	if err != nil {
		return entities.Portfolio{}, err
	}
	defer rows.Close()

	positionsMap := make(map[string]entities.Position)
	for rows.Next() {
		var assetID sql.NullString
		var quantity sql.NullInt64

		if err := rows.Scan(&initialCash, &cash, &assetID, &quantity); err != nil {
			return entities.Portfolio{}, err
		}

		if !assetID.Valid {
			continue
		}

		assetIDVO, err := vos.ParseID(assetID.String)
		if err != nil {
			return entities.Portfolio{}, err
		}

		p, err := entities.NewPosition(assetIDVO, int(quantity.Int64))
		if err != nil {
			return entities.Portfolio{}, err
		}

		positionsMap[assetID.String] = p
	}

	return entities.NewPortfolio(
		portfolioID, vos.ParseToDecimal(initialCash),
		vos.ParseToDecimal(cash), positionsMap,
	)
}
