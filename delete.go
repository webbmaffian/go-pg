package pg

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Delete(ctx context.Context, db *pgxpool.Pool, table TableSource, condition Condition) (err error) {
	var b strings.Builder
	b.Grow(64)
	args := make([]any, 0, 2)

	b.WriteString("DELETE FROM ")
	table.buildQuery(&b, nil)
	b.WriteByte('\n')
	b.WriteString("WHERE ")
	condition.run(&b, &args)

	_, err = db.Exec(ctx, b.String(), args...)

	if err != nil {
		err = QueryError{
			err:   err.Error(),
			query: b.String(),
			args:  args,
		}
	}

	return
}
