package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return db.Query(ctx, sql, args...)
}

func QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return db.QueryRow(ctx, sql, args...)
}

func Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return db.Exec(ctx, sql, args...)
}
