package storage

import (
	"context"
	"database/sql"
)

type Session interface {
	Begin() (Session, error)
	Rollback() error
	Commit() error
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

var sessionKey = "ctx:storage_session"

func WithContext(ctx context.Context, ses Session) context.Context {
	return context.WithValue(ctx, sessionKey, ses)
}

func FromContext(ctx context.Context) Session {
	v := ctx.Value(sessionKey)
	if v == nil {
		return nil
	}
	return v.(Session)
}
