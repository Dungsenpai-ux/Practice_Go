package db

import (
	"context"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Init establishes a pgx pool and runs migrations.
func Init(cfg *config.Config) (*pgxpool.Pool, error) {
	if err := runMigrations(cfg.DBUrl); err != nil {
		return nil, err
	}
	pcfg, err := pgxpool.ParseConfig(cfg.DBUrl)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), pcfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func runMigrations(dsn string) error {
	m, err := migrate.New("file:///app/service/migrations", dsn)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
