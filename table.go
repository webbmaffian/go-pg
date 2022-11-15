package pg

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Table(db *pgxpool.Pool, name string) TableSource {
	return TableSource{db, name}
}

type TableSource struct {
	db   *pgxpool.Pool
	name string
}

func (t TableSource) buildQuery(b *strings.Builder, args *[]any) {
	writeIdentifier(b, t.name)
}

func (t TableSource) Select(ctx context.Context, dest any, q SelectQuery, options ...SelectOptions) error {
	q.From = t

	return Select(ctx, t.db, dest, q, options...)
}

func (t TableSource) Iterate(ctx context.Context, q SelectQuery, iterator func(values []any) error) error {
	q.From = t

	return Iterate(ctx, t.db, q, iterator)
}

func (t TableSource) IterateRaw(ctx context.Context, q SelectQuery, iterator func(values [][]byte) error) error {
	q.From = t

	return IterateRaw(ctx, t.db, q, iterator)
}

func (t TableSource) Insert(ctx context.Context, src any, onConflict ...OnConflictUpdate) error {
	return Insert(ctx, t.db, t.name, src, onConflict...)
}

func (t TableSource) Update(ctx context.Context, src any, condition Condition) error {
	return Update(ctx, t.db, t.name, src, condition)
}

func (t TableSource) Delete(ctx context.Context, condition Condition) error {
	return Delete(ctx, t.db, t.name, condition)
}
