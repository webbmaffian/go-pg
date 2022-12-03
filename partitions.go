package pg

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePartition(ctx context.Context, db *pgxpool.Pool, table pgx.Identifier, partition pgx.Identifier, value string) (err error) {
	var b strings.Builder
	b.Grow(64)
	args := make([]any, 0, 1)

	b.WriteString("CREATE TABLE ")
	b.WriteString(partition.Sanitize())
	b.WriteString(" PARTITION OF ")
	b.WriteString(table.Sanitize())
	b.WriteString(" FOR VALUES IN (")
	writeString(&b, value)
	b.WriteByte(')')

	_, err = db.Exec(ctx, b.String(), args...)

	return
}