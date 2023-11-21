package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"sync"

	"fangaoxs.com/go-chat/internal/storage"
)

// 将nil赋值给*Tx类型的指针，并且将其断言为storage.Session
// 目的：在编译时检查*Tx是否实现了storage.Session接口
var _ storage.Session = (*Tx)(nil)

func NewTxContext(ctx context.Context, db *sql.DB) *Tx {
	return &Tx{
		ctx:       ctx,
		exec:      db,
		db:        db,
		savePoint: 0,
	}
}

type Tx struct {
	sync.Mutex
	ctx context.Context

	exec executor
	db   *sql.DB
	txi  *sql.Tx

	savePoint int
	next      *Tx
	resolved  bool
}

func (tx *Tx) Begin() (storage.Session, error) {
	tx.Lock()
	defer tx.Unlock()

	if tx.txi == nil {
		txi, err := tx.db.BeginTx(tx.ctx, nil)
		if err != nil {
			return nil, err
		}
		tx.txi = txi
		tx.exec = txi
		return tx, nil
	}

	tx.next = &Tx{
		ctx:       tx.ctx,
		exec:      tx.exec,
		db:        tx.db,
		txi:       tx.txi,
		savePoint: tx.savePoint + 1,
	}

	_, err := tx.txi.ExecContext(tx.ctx, "SAVEPOINT sp"+strconv.Itoa(tx.next.savePoint))
	if err != nil {
		return nil, err
	}

	return tx.next, nil
}

func (tx *Tx) Rollback() error {
	tx.Lock()
	defer tx.Unlock()

	if tx.savePoint > 0 {
		if tx.resolved {
			return nil
		}
		tx.resolved = true
		if _, err := tx.txi.ExecContext(tx.ctx, "ROLLBACK TO SAVEPOINT sp"+strconv.Itoa(tx.savePoint)); err != nil {
			return err
		}
		return nil
	}

	return tx.txi.Rollback()
}

func (tx *Tx) Commit() error {
	tx.Lock()
	defer tx.Unlock()

	if tx.savePoint > 0 {
		if tx.resolved {
			return nil
		}
		tx.resolved = true
		if _, err := tx.txi.ExecContext(tx.ctx, "RELEASE SAVEPOINT sp"+strconv.Itoa(tx.savePoint)); err != nil {
			return err
		}
		return nil
	}

	return tx.txi.Commit()
}

func (tx *Tx) Exec(query string, args ...any) (sql.Result, error) {
	return tx.exec.ExecContext(tx.ctx, query, args...)
}

func (tx *Tx) Query(query string, args ...any) (*sql.Rows, error) {
	return tx.exec.QueryContext(tx.ctx, query, args...)
}

func (tx *Tx) QueryRow(query string, args ...any) *sql.Row {
	return tx.exec.QueryRowContext(tx.ctx, query, args...)
}

type executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
