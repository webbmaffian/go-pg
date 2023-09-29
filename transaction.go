package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrReleased = errors.New("connection released")

func Transaction(ctx context.Context, pool *pgxpool.Pool, cb func(context.Context, Tx) error) (err error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		return
	}

	defer conn.Release()

	tx := Tx{
		db: conn,
	}

	if _, err = tx.db.Exec(ctx, "BEGIN"); err != nil {
		return
	}

	if err = cb(ctx, tx); err != nil {
		conn.Exec(ctx, "ROLLBACK")
	} else {
		conn.Exec(ctx, "COMMIT")
	}

	return
}

func ReadonlyTransaction(ctx context.Context, pool *pgxpool.Pool, cb func(context.Context, Tx) error) (err error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		return
	}

	defer conn.Release()

	tx := Tx{
		db: conn,
	}

	if _, err = tx.db.Exec(ctx, "BEGIN READ ONLY"); err != nil {
		return
	}

	if err = cb(ctx, tx); err != nil {
		conn.Exec(ctx, "ROLLBACK")
	} else {
		conn.Exec(ctx, "COMMIT")
	}

	return
}

var _ conn = Tx{}

type Tx struct {
	db *pgxpool.Conn
}

func (tx Tx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return tx.db.Exec(ctx, sql, arguments...)
}

func (tx Tx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return tx.db.Query(ctx, sql, args...)
}

func (tx Tx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return tx.db.QueryRow(ctx, sql, args...)
}

func (tx Tx) Select(ctx context.Context, dest any, q SelectQuery, options ...SelectOptions) error {
	return Select(ctx, tx.db, dest, q, options...)
}

func (tx Tx) SelectTotal(ctx context.Context, dest *int, q SelectQuery) error {
	return SelectTotal(ctx, tx.db, dest, q)
}

func (tx Tx) SelectValue(ctx context.Context, dest any, q SelectQuery) error {
	return SelectValue(ctx, tx.db, dest, q)
}

func (tx Tx) Iterate(ctx context.Context, t TableSource, q SelectQuery, iterator func(values []any) error) error {
	return Iterate(ctx, tx.db, q, iterator)
}

func (tx Tx) IterateRaw(ctx context.Context, t TableSource, q SelectQuery, iterator func(values [][]byte) error) error {
	return IterateRaw(ctx, tx.db, q, iterator)
}

func (tx Tx) Insert(ctx context.Context, t TableSource, src any, onConflict ...ConflictAction) error {
	return Insert(ctx, tx.db, t, src, onConflict...)
}

func (tx Tx) InsertRow(t TableSource, onConflict ...ConflictAction) RowInserter {
	return InsertRow(tx.db, t, onConflict...)
}

func (tx Tx) Update(ctx context.Context, t TableSource, src any, condition Condition) error {
	return Update(ctx, tx.db, t, src, condition)
}

func (tx Tx) Delete(ctx context.Context, t TableSource, condition Condition) error {
	return Delete(ctx, tx.db, t, condition)
}

func (tx Tx) CopyFrom(ctx context.Context, t TableSource, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return CopyFrom(ctx, tx.db, t.identifier, columnNames, rowSrc)
}

func (tx Tx) Truncate(ctx context.Context, t TableSource) (err error) {
	return TruncateTable(ctx, tx.db, t.identifier)
}
