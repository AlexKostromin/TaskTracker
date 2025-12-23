package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPPort        string
	ShutdownTimeout time.Duration
}

func Load() *Config {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	shutdownTimeout := getEnvAsInt("SHUTDOWN_TIMEOUT_SEC", 30)

	return &Config{
		HTTPPort:        httpPort,
		ShutdownTimeout: time.Duration(shutdownTimeout) * time.Second,
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
