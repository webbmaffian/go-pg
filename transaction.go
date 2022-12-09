package pg

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrReleased = errors.New("connection released")

func Transaction(ctx context.Context, pool *pgxpool.Pool) (tx *Tx) {
	tx = new(Tx)
	tx.ctx, tx.cancel = context.WithCancel(ctx)
	tx.db, tx.err = pool.Acquire(tx.ctx)

	if tx.err != nil {
		return
	}

	_, tx.err = tx.db.Exec(tx.ctx, "BEGIN")

	tx.maybeRollback()

	go func() {
		<-tx.ctx.Done()
		tx.Rollback()
	}()

	return
}

type Tx struct {
	ctx    context.Context
	cancel context.CancelFunc
	db     *pgxpool.Conn
	err    error
	mu     sync.Mutex
}

func (tx *Tx) Select(t TableSource, dest any, q SelectQuery, options ...SelectOptions) error {
	if tx.err != nil {
		return tx.err
	}

	q.From = t
	tx.err = Select(tx.ctx, t.db, dest, q, options...)
	return tx.maybeRollback()
}

func (tx *Tx) SelectTotal(t TableSource, dest *int, q SelectQuery) error {
	if tx.err != nil {
		return tx.err
	}

	q.From = t
	tx.err = SelectTotal(tx.ctx, t.db, dest, q)
	return tx.maybeRollback()
}

func (tx *Tx) Iterate(t TableSource, q SelectQuery, iterator func(values []any) error) error {
	if tx.err != nil {
		return tx.err
	}

	q.From = t
	tx.err = Iterate(tx.ctx, t.db, q, iterator)
	return tx.maybeRollback()
}

func (tx *Tx) IterateRaw(t TableSource, q SelectQuery, iterator func(values [][]byte) error) error {
	if tx.err != nil {
		return tx.err
	}

	q.From = t
	tx.err = IterateRaw(tx.ctx, t.db, q, iterator)
	return tx.maybeRollback()
}

func (tx *Tx) Insert(t TableSource, src any, onConflict ...ConflictAction) error {
	if tx.err != nil {
		return tx.err
	}

	tx.err = Insert(tx.ctx, t.db, t, src, onConflict...)
	return tx.maybeRollback()
}

func (tx *Tx) Update(t TableSource, src any, condition Condition) error {
	if tx.err != nil {
		return tx.err
	}

	tx.err = Update(tx.ctx, t.db, t, src, condition)
	return tx.maybeRollback()
}

func (tx *Tx) Delete(t TableSource, condition Condition) error {
	if tx.err != nil {
		return tx.err
	}

	tx.err = Delete(tx.ctx, t.db, t, condition)
	return tx.maybeRollback()
}

func (tx *Tx) Truncate(t TableSource) (err error) {
	if tx.err != nil {
		return tx.err
	}

	tx.err = TruncateTable(tx.ctx, t.db, t.identifier)
	return tx.maybeRollback()
}

func (tx *Tx) maybeRollback() error {
	if tx.err != nil {
		tx.Rollback()
	}

	return tx.err
}

func (tx *Tx) Commit() (err error) {
	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.err != nil {
		return tx.err
	}

	_, tx.err = tx.db.Exec(tx.ctx, "COMMIT")
	tx.release()

	if tx.err != nil {
		return tx.err
	}

	tx.err = ErrReleased
	return
}

func (tx *Tx) Rollback() {
	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.db == nil {
		return
	}

	tx.db.Exec(tx.ctx, "ROLLBACK")
	tx.release()
	tx.err = ErrReleased
}

func (tx *Tx) release() {
	if tx.db != nil {
		tx.db.Release()
		tx.db = nil
	}

	tx.cancel()
}
