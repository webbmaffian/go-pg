package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Query(ctx context.Context, db *pgxpool.Pool, sql string, args ...any) (pgx.Rows, error) {
	return db.Query(ctx, sql, args...)
}

func QueryRow(ctx context.Context, db *pgxpool.Pool, sql string, args ...any) pgx.Row {
	return db.QueryRow(ctx, sql, args...)
}

func Exec(ctx context.Context, db *pgxpool.Pool, sql string, args ...any) (pgconn.CommandTag, error) {
	return db.Exec(ctx, sql, args...)
}
