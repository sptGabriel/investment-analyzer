package prices

import "context"

func (r Repository) Delete(ctx context.Context) error {
	const query = `DELETE FROM prices;`
	_, err := r.q.Exec(ctx, query)
	return err
}
