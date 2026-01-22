package config

import (
	"fmt"
	postgres "group_service/pkg/db"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort        int           `env:"GRPC_PORT" env-default:"50051"`
	GHTimeout       time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" env-default:"15s"`
	UserServiceAddr string        `env:"USER_SERVICE_ADDRESS" env-default:"localhost:50051"`
	MigrationPath   string        `env:"MIGRATION_PATH" env-default:":file://migrations"`

	postgres.PostgresConfig
}

func ParseConfigFromEnv() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return &cfg, nil
}
