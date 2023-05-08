package pg

import (
	"context"
)

func SelectValue(ctx context.Context, db conn, dest any, q SelectQuery) (err error) {
	if err = q.run(ctx, db); err != nil {
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
