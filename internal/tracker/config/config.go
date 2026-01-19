package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPPort        string
	ShutdownTimeout time.Duration
	DBURI           string
	MigrationsDir   string
}

func Load() *Config {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	shutdownTimeout := getEnvAsInt("SHUTDOWN_TIMEOUT_SEC", 30)

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		// Совместимо с примером из postgres/README.md и docker-compose.
		dbURI = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		// Миграции лежат в папке postgres/migrations в корне репозитория.
		migrationsDir = "./postgres/migrations"
	}

	return &Config{
		HTTPPort:        httpPort,
		ShutdownTimeout: time.Duration(shutdownTimeout) * time.Second,
		DBURI:           dbURI,
		MigrationsDir:   migrationsDir,
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
