package pg

import (
	"context"
	"strings"
)

func Exists(ctx context.Context, db conn, table TableSource, condition Condition) (exists bool, err error) {
	var b strings.Builder
	b.Grow(64)
	args := make([]any, 0, 2)

	b.WriteString("SELECT EXISTS (SELECT 1 FROM ")
	table.encodeQuery(&b, nil)
	b.WriteByte('\n')
	b.WriteString("WHERE ")
	condition.encodeCondition(&b, &args)
	b.WriteByte(')')

	row := db.QueryRow(ctx, b.String(), args...)
	err = row.Scan(&exists)
	return
}
