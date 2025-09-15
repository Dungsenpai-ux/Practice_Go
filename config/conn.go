package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() error {
	dsn := "postgres://postgres:123456@localhost:5432/movies_db?sslmode=disable" // Cập nhật mật khẩu
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}
	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}
	return Pool.Ping(context.Background())
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
