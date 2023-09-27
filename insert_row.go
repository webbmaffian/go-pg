package pg

import (
	"bytes"
	"context"
	"strconv"
)

func InsertRow(db conn, table TableSource, onConflict ...ConflictAction) RowInserter {
	r := &rowInserter{
		db:    db,
		table: table,
	}

	if onConflict != nil {
		r.onConflict = onConflict[0]
	}

	return r
}

var _ RowInserter = (*rowInserter)(nil)

type rowInserter struct {
	columns    []string
	values     []any
	db         conn
	table      TableSource
	onConflict ConflictAction
}

func (r *rowInserter) Value(column string, value any) RowInserter {
	r.columns = append(r.columns, column)
	r.values = append(r.values, value)

	return r
}

func (r *rowInserter) Exec(ctx context.Context) (err error) {
	var b bytes.Buffer
	b.Grow(200)
	r.buildQuery(&b)
	queryString := b2s(b.Bytes())
	_, err = r.db.Exec(ctx, queryString, r.values...)

	if err != nil {
		err = QueryError{
			err:   err.Error(),
			query: queryString,
			args:  r.values,
		}
	}

	return
}

func (r *rowInserter) ExecAndReturn(ctx context.Context, column string, bind any) (err error) {
	var b bytes.Buffer
	b.Grow(200)
	r.buildQuery(&b)

	b.WriteByte('\n')
	b.WriteString("RETURNING ")
	b.WriteByte('"')
	b.WriteString(column)
	b.WriteByte('"')

	queryString := b2s(b.Bytes())
	rows, err := r.db.Query(ctx, queryString, r.values...)

	if err != nil {
		return QueryError{
			err:   err.Error(),
			query: queryString,
			args:  r.values,
		}
	}

	defer rows.Close()

	if rows.Next() {
		return rows.Scan(bind)
	}

	return
}

func (r *rowInserter) buildQuery(b *bytes.Buffer) {
	b.WriteString("INSERT INTO ")
	b.WriteString(r.table.identifier.Sanitize())
	b.WriteString(" (")

	for i := range r.columns {
		if i != 0 {
			b.WriteByte(',')
		}

		b.WriteByte('"')
		b.WriteString(r.columns[i])
		b.WriteByte('"')
	}

	b.WriteString(") VALUES(")

	for i := range r.values {
		if i != 0 {
			b.WriteByte(',')
		}

		b.WriteByte('$')
		b.Write(strconv.AppendInt(b.AvailableBuffer(), int64(i+1), 10))
	}

	b.WriteByte(')')

	if r.onConflict != nil {
		r.onConflict.encodeConflictHandler(b, r.columns, &r.values)
	}
}
