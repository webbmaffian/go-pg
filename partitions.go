package pg

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
)

var PartitionName = func(tableName pgx.Identifier, partition string) pgx.Identifier {
	l := len(tableName)
	partitionName := make(pgx.Identifier, l)
	copy(partitionName, tableName)
	partitionName[l-1] = strings.Join([]string{"part", partitionName[l-1], partition}, "_")

	// schema.table becomes schema.part_table_partitionkey
	return partitionName
}

func CreatePartition(ctx context.Context, db conn, table pgx.Identifier, partition pgx.Identifier, value string) (err error) {
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
