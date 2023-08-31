package pg

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
)

type SelectQuery struct {
	Select  Columnar
	From    Queryable
	Join    Join
	Where   Condition
	GroupBy Columnar
	Having  Condition
	OrderBy OrderByColumnar
	Limit   int
	Offset  int

	result pgx.Rows
	error  error
}

type SelectOptions struct {
	BeforeMarshal func(data *map[string]any) error
	AfterMarshal  func(data *map[string]any) error
}

func (q SelectQuery) IsZero() bool {
	return q.From.IsZero()
}

func (q *SelectQuery) Error() error {
	return q.error
}

func (q *SelectQuery) String() string {
	var b strings.Builder
	q.encodeQuery(&b, &[]any{})
	return b.String()
}

func (q *SelectQuery) encodeQuery(b *strings.Builder, args *[]any) {
	b.Grow(300)

	if q.Select != nil && !q.Select.IsZero() {
		b.WriteString("SELECT ")
		q.Select.encodeColumn(b)
		b.WriteByte('\n')
	} else {
		b.WriteString("SELECT *")
		b.WriteByte('\n')
	}

	if q.From != nil && !q.From.IsZero() {
		b.WriteString("FROM ")
		q.From.encodeQuery(b, args)
		b.WriteByte('\n')
	}

	if q.Join != nil && !q.Join.IsZero() {
		q.Join.encodeJoin(b, args)
	}

	if q.Where != nil && !q.Where.IsZero() {
		b.WriteString("WHERE ")
		q.Where.encodeCondition(b, args)
		b.WriteByte('\n')
	}

	if q.GroupBy != nil && !q.GroupBy.IsZero() {
		b.WriteString("GROUP BY ")
		q.GroupBy.encodeColumnIdentifier(b)
		b.WriteByte('\n')
	}

	if q.Having != nil && !q.Having.IsZero() {
		b.WriteString("HAVING ")
		q.Having.encodeCondition(b, args)
		b.WriteByte('\n')
	}

	if q.OrderBy != nil && !q.OrderBy.IsZero() {
		b.WriteString("ORDER BY ")
		q.OrderBy.encodeOrderBy(b)
		b.WriteByte('\n')
	}

	if q.Limit > 0 {
		b.WriteString("LIMIT ")
		writeInt(b, q.Limit)
		b.WriteByte('\n')
	}

	if q.Offset > 0 {
		b.WriteString("OFFSET ")
		writeInt(b, q.Offset)
		b.WriteByte('\n')
	}
}

func (q *SelectQuery) run(ctx context.Context, db conn) (err error) {
	var b strings.Builder
	args := make([]any, 0, 5)
	q.encodeQuery(&b, &args)

	q.result, err = db.Query(ctx, q.String(), args...)

	if err != nil {
		err = QueryError{
			err:   err.Error(),
			query: q.String(),
			args:  args,
		}
	}

	return
}

func (q *SelectQuery) Run(ctx context.Context, db conn) error {
	return q.run(ctx, db)
}

func (q *SelectQuery) Next() bool {
	if q.result == nil {
		return false
	}

	n := q.result.Next()

	if !n {
		q.result = nil
	}

	return n
}

func (q *SelectQuery) Scan(dest ...any) error {
	if q.result == nil {
		return errors.New("Result is closed")
	}

	return q.result.Scan(dest...)
}

func (q *SelectQuery) Close() {
	if q.result != nil {
		q.result.Close()
		q.result = nil
	}
}
