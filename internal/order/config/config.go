package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPPort        string
	ShutdownTimeout time.Duration
	// DatabaseURL string
	// RedisURL    string
}

func Load() *Config {
	// Загрузка из env переменных или файла
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080" // значение по умолчанию
	}

	// Таймаут для graceful shutdown
	shutdownTimeout := getEnvAsInt("SHUTDOWN_TIMEOUT_SEC", 30)

	return &Config{
		HTTPPort:        httpPort,
		ShutdownTimeout: time.Duration(shutdownTimeout) * time.Second,
	}
}

// Вспомогательная функция для чтения env переменных как int
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
