package settings

import "context"

func (r Repository) IsCsvImported(ctx context.Context) (bool, error) {
	const query = `
		SELECT value FROM system_config WHERE key = 'csv_imported';`

	var value string

	if err := r.q.QueryRow(ctx, query).Scan(&value); err != nil {
		return false, err
	}

	return value == "true", nil
}

func (r Repository) SetCsvImported(ctx context.Context) error {
	const query = `
		UPDATE system_config
		SET value = 'true'
		WHERE key = 'csv_imported';
	`
	_, err := r.q.Exec(ctx, query)
	return err
}
