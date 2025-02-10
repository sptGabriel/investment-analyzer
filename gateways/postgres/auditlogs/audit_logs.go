package auditlogs

import (
	"context"

	"github.com/sptGabriel/investment-analyzer/domain/entities"
)

func (r Repository) Save(
	ctx context.Context, log entities.AuditLog,
) error {
	const query = `
	INSERT INTO audit_logs (timestamp, ip, port, params, result, error_message)
	VALUES ($1, $2, $3, $4, $5, $6);
`

	_, err := r.q.Exec(
		ctx, query,
		log.Timestamp, log.IP, log.Port, log.Params, log.Result, log.ErrorMessage,
	)
	if err != nil {
		return err
	}

	return nil
}
