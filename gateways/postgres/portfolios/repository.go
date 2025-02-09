package portfolios

import "github.com/sptGabriel/investment-analyzer/gateways/postgres"

type Repository struct {
	q postgres.Querier
}

func New(q postgres.Querier) *Repository {
	return &Repository{
		q: q,
	}
}
