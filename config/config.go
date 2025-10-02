package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application environment configuration (no connection logic here)
type Config struct {
	Port          string
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	DBUrl         string
	DBSSLMode     string
	MemcachedAddr string
	Version       string
	OtelEndpoint  string
	OtelService   string
	OtelEnv       string
	OtelSampler   string
}

// Load loads .env (if present) and builds a Config instance.
func Load() *Config {
	_ = godotenv.Load()
	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "postgres"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		MemcachedAddr: getEnv("MEMCACHED_ADDR", "127.0.0.1:11211"),
		Version:       "v1.0.0",
		OtelEndpoint:  getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318"),
		OtelService:   getEnv("OTEL_SERVICE_NAME", "practice-go-api"),
		OtelEnv:       getEnv("OTEL_ENVIRONMENT", "dev"),
		OtelSampler:   getEnv("OTEL_TRACES_SAMPLER", "parentbased_always_on"),
	}
	cfg.DBUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)
	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
