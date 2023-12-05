package pg

import (
	"context"
)

func Iterate(ctx context.Context, db conn, q SelectQuery, iterator func(values []any) error) (err error) {
	if q.Select == nil {
		return ErrNoColumns
	}

	err = q.run(ctx, db)

	if err != nil {
		return
	}

	defer q.Close()

	var values []any

	for q.Next() {
		if values, err = q.result.Values(); err != nil {
			return
		}

		if err = iterator(values); err != nil {
			return
		}
	}

	return
}

func IterateRaw(ctx context.Context, db conn, q SelectQuery, iterator func(values [][]byte) error) (err error) {
	if q.Select == nil {
		return ErrNoColumns
	}

	err = q.run(ctx, db)

	if err != nil {
		return
	}

	defer q.Close()

	for q.Next() {
		if err = iterator(q.result.RawValues()); err != nil {
			return
		}
	}

	return
}
