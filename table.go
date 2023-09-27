package pg

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Table(db *pgxpool.Pool, name ...string) TableSource {
	return TableSource{
		db:         db,
		identifier: pgx.Identifier(name),
	}
}

type TableSource struct {
	db                 *pgxpool.Pool
	identifier         pgx.Identifier
	originalIdentifier pgx.Identifier
}

func (t TableSource) IsZero() bool {
	return t.identifier == nil
}

func (t TableSource) encodeQuery(b ByteStringWriter, args *[]any) {
	if t.originalIdentifier != nil {
		b.WriteString(t.originalIdentifier.Sanitize())
		b.WriteString(" AS ")
	}

	b.WriteString(t.identifier.Sanitize())
}

func (t TableSource) Alias(alias string) TableSource {
	aliased := t
	aliased.identifier = pgx.Identifier{alias}

	if t.originalIdentifier == nil {
		aliased.originalIdentifier = t.identifier
	}

	return aliased
}

func (t TableSource) Select(ctx context.Context, dest any, q SelectQuery, options ...SelectOptions) error {
	q.From = t

	return Select(ctx, t.db, dest, q, options...)
}

func (t TableSource) SelectTotal(ctx context.Context, dest *int, q SelectQuery) error {
	q.From = t

	return SelectTotal(ctx, t.db, dest, q)
}

func (t TableSource) SelectValue(ctx context.Context, dest any, q SelectQuery) error {
	q.From = t

	return SelectValue(ctx, t.db, dest, q)
}

func (t TableSource) Iterate(ctx context.Context, q SelectQuery, iterator func(values []any) error) error {
	q.From = t

	return Iterate(ctx, t.db, q, iterator)
}

func (t TableSource) IterateRaw(ctx context.Context, q SelectQuery, iterator func(values [][]byte) error) error {
	q.From = t

	return IterateRaw(ctx, t.db, q, iterator)
}

func (t TableSource) Query(ctx context.Context, q *SelectQuery) (err error) {
	q.From = t

	return q.run(ctx, t.db)
}

func (t TableSource) Insert(ctx context.Context, src any, onConflict ...ConflictAction) error {
	return Insert(ctx, t.db, t, src, onConflict...)
}

func (t TableSource) InsertRow(ctx context.Context, onConflict ...ConflictAction) RowInserter {
	return InsertRow(t.db, t, onConflict...)
}

func (t TableSource) Update(ctx context.Context, src any, condition Condition) error {
	return Update(ctx, t.db, t, src, condition)
}

func (t TableSource) Delete(ctx context.Context, condition Condition) error {
	return Delete(ctx, t.db, t, condition)
}

// Presumes that there is a scheme with the same name as the table.
func (t TableSource) Partition(partition string) TableSource {
	return TableSource{
		db:         t.db,
		identifier: PartitionName(t.identifier, partition),
	}
}

func (t TableSource) Truncate(ctx context.Context) (err error) {
	return TruncateTable(ctx, t.db, t.identifier)
}

func (t TableSource) Drop(ctx context.Context) (err error) {
	return DropTable(ctx, t.db, t.identifier)
}

// Presumes that there is a scheme with the same name as the table, and that the partition value is a string.
func (t TableSource) CreatePartition(ctx context.Context, partition string) (err error) {
	return CreatePartition(ctx, t.db, t.identifier, PartitionName(t.identifier, partition), partition)
}

func (t TableSource) CopyFrom(ctx context.Context, columnsNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return t.db.CopyFrom(ctx, t.identifier, columnsNames, rowSrc)
}

func (t *TableSource) Column(path ...string) AliasedColumnar {
	return column{
		path:  path,
		table: t,
	}
}

func (t *TableSource) JsonColumn(path ...string) AliasedColumnar {
	return jsonColumn{
		path:  path,
		table: t,
	}
}

func TruncateTable(ctx context.Context, db *pgxpool.Pool, table pgx.Identifier) (err error) {
	var b strings.Builder
	b.Grow(64)

	b.WriteString("TRUNCATE TABLE ")
	b.WriteString(table.Sanitize())

	_, err = db.Exec(ctx, b.String())

	return
}

func DropTable(ctx context.Context, db *pgxpool.Pool, table pgx.Identifier) (err error) {
	var b strings.Builder
	b.Grow(64)

	b.WriteString("DROP TABLE ")
	b.WriteString(table.Sanitize())

	_, err = db.Exec(ctx, b.String())

	return
}
