package config

import (
	"context"
	"fmt"
	"os"

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