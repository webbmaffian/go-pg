package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SelectTotal(ctx context.Context, db *pgxpool.Pool, dest *int, q SelectQuery) (err error) {
	q.GroupBy = nil
	q.Limit = 0
	q.Offset = 0
	q.OrderBy = nil
	q.Select = Count("count")

	err = q.run(ctx, db)

	if err != nil {
		return
	}

	defer q.Close()

	if q.Next() {
		err = q.Scan(dest)

		if err != nil {
			return
		}
	}

	return
}
