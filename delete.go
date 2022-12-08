package pg

import (
	"context"
	"strings"
)

func Delete(ctx context.Context, db conn, table TableSource, condition Condition) (err error) {
	var b strings.Builder
	b.Grow(64)
	args := make([]any, 0, 2)

	b.WriteString("DELETE FROM ")
	table.encodeQuery(&b, nil)
	b.WriteByte('\n')
	b.WriteString("WHERE ")
	condition.encodeCondition(&b, &args)

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
