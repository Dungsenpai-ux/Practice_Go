package config

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

func Connect() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Lấy các giá trị từ biến môi trường
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	// Run migrations before establishing the pool to ensure schema is up to date
	if err := runMigrations(dsn); err != nil {
		return err
	}

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

func runMigrations(dsn string) error {
	m, err := migrate.New("file://service/migrations", dsn)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
