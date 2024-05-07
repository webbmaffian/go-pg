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

	r.buf.Grow(200)

	if onConflict != nil {
		r.onConflict = onConflict[0]
	}

	return r
}

var _ RowInserter = (*rowInserter)(nil)

type rowInserter struct {
	columns    []string
	values     []any
	skipUpdate []bool
	db         conn
	table      TableSource
	onConflict ConflictAction
	buf        bytes.Buffer
}

func (r *rowInserter) Reset() {
	r.columns = r.columns[:0]
	r.values = r.values[:0]
	r.skipUpdate = r.skipUpdate[:0]
	r.buf.Reset()
}

func (r *rowInserter) Value(column string, value any, skipUpdate ...bool) RowInserter {
	r.columns = append(r.columns, column)
	r.values = append(r.values, value)

	if len(skipUpdate) > 0 {
		r.skipUpdate = append(r.skipUpdate, skipUpdate[0])
	} else {
		r.skipUpdate = append(r.skipUpdate, false)
	}

	return r
}

func (r *rowInserter) Exec(ctx context.Context) (err error) {
	r.buildQuery(&r.buf)
	queryString := b2s(r.buf.Bytes())
	_, err = r.db.Exec(ctx, queryString, r.values...)

	if err != nil {
		err = QueryError{
			err:   err.Error(),
			query: queryString,
			args:  r.values,
		}
	}

	r.Reset()

	return
}

func (r *rowInserter) ExecAndReturn(ctx context.Context, column string, bind any) (err error) {
	r.buildQuery(&r.buf)

	r.buf.WriteByte('\n')
	r.buf.WriteString("RETURNING ")
	r.buf.WriteByte('"')
	r.buf.WriteString(column)
	r.buf.WriteByte('"')

	queryString := b2s(r.buf.Bytes())
	rows, err := r.db.Query(ctx, queryString, r.values...)

	if err != nil {
		return QueryError{
			err:   err.Error(),
			query: queryString,
			args:  r.values,
		}
	}

	r.Reset()

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
		r.onConflict.encodeConflictHandler(b, r.columns, r.skipUpdate, &r.values)
	}
}
