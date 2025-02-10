package trades

import "context"

func (r Repository) Delete(ctx context.Context) error {
	const query = `DELETE FROM trades;`
	_, err := r.q.Exec(ctx, query)
	return err
}
