package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/storage"

	"github.com/golang-migrate/migrate/v4"
	dStub "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // 导入 PostgresSQL 驱动
)

var _ storage.Storage = (*postgres)(nil)

func New(env environment.Env) (storage.Storage, error) {
	// TODO: 解析env.DSN

	db, err := sql.Open("postgres", env.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres with %s, %w", env.DSN, err)
	}

	pg := &postgres{db: db}
	if err = pg.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	return pg, nil
}

type postgres struct {
	db *sql.DB
}

func (p *postgres) Close() error {
	if p.db != nil {
		if err := p.db.Close(); err != nil {
			return fmt.Errorf("failed to close db: %w", err)
		}
	}

	return nil
}

func (p *postgres) NewSession(ctx context.Context) (storage.Session, error) {
	if ses := storage.FromContext(ctx); ses != nil {
		return ses, nil
	}

	return NewTxContext(ctx, p.db), nil
}

func (p *postgres) migrate() error {
	instance, err := dStub.WithInstance(p.db, &dStub.Config{})
	if err != nil {
		return fmt.Errorf("unable to create migration driver instance: %w", err)
	}

	// migrations目录的路径
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	u, _ := url.Parse(dir)
	u.Scheme = "file"
	u = u.JoinPath("migrations")
	m, err := migrate.NewWithDatabaseInstance(u.String(), "postgres", instance)
	if err != nil {
		return fmt.Errorf("unable to create migration instance: %w", err)
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil // do nothing, skip
	}

	return err
}
