package pg

import (
	"context"
	"strings"

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
	encodeColumn(b *strings.Builder)
	encodeColumnIdentifier(b *strings.Builder)
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

type Condition interface {
	encodeCondition(b *strings.Builder, args *[]any)
}

type Queryable interface {
	encodeQuery(b *strings.Builder, args *[]any)
}

type conn interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
