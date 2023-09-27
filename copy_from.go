package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CopyFrom(ctx context.Context, db *pgxpool.Conn, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return db.Conn().CopyFrom(ctx, tableName, columnNames, rowSrc)
}
