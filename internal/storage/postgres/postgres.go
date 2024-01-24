package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"net/http"

	"fangaoxs.com/go-chat/environment"
	"fangaoxs.com/go-chat/internal/storage"

	"github.com/golang-migrate/migrate/v4"
	dStub "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
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

//go:embed migrations/*.sql
var migrations embed.FS

func (p *postgres) migrate() error {
	instance, err := dStub.WithInstance(p.db, &dStub.Config{})
	if err != nil {
		return fmt.Errorf("unable to create migration driver instance: %w", err)
	}
	src, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return fmt.Errorf("failed to open migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("httpfs", src, "postgres", instance)
	if err != nil {
		return fmt.Errorf("unable to create migration instance: %w", err)
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil // do nothing, skip
	}

	return err
}
