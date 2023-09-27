package pg

import (
	"context"
	"io"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type MutationType uint8

const (
	Inserting MutationType = 1
	Updating  MutationType = 2
)

type BeforeMutation interface {
	BeforeMutation(ctx context.Context, mutationType MutationType) error
}

type AfterMutation interface {
	AfterMutation(ctx context.Context, mutationType MutationType)
}

type IsZeroer interface {
	IsZero() bool
}

type Columnar interface {
	IsZeroer
	encodeColumn(b ByteStringWriter)
	encodeColumnIdentifier(b ByteStringWriter)
	has(col string) bool
}

type AliasedColumnar interface {
	Columnar
	Alias(string) AliasedColumnar
}

type MultiColumnar interface {
	Columnar
	Append(...Columnar) MultiColumnar
}

type OrderByColumnar interface {
	IsZeroer
	encodeOrderBy(b ByteStringWriter)
}

type Condition interface {
	IsZeroer
	encodeCondition(b ByteStringWriter, args *[]any)
}

type Queryable interface {
	IsZeroer
	encodeQuery(b ByteStringWriter, args *[]any)
}

type conn interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type RawData interface {
	Columnar
	Condition
	Queryable
	Alias(alias string) RawData
	Column(path ...string) AliasedColumnar
}

type ByteStringWriter interface {
	io.Writer
	io.ByteWriter
	io.StringWriter
	Grow(n int)
}

type RowInserter interface {
	Value(column string, value any) RowInserter
	Exec(ctx context.Context) (err error)
	ExecAndReturn(ctx context.Context, column string, bind any) (err error)
}
