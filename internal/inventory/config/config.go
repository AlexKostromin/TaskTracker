package config

import (
	"time"
)

type Config struct {
	GRPCPort        string        `mapstructure:"GRPC_PORT"`
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
}

func Load() *Config {
	return &Config{
		GRPCPort:        "50053",
		ShutdownTimeout: 30 * time.Second,
	}
}
