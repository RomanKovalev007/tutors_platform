package config

import (
	"fmt"
	"user-service/pkg/kafka"
	"user-service/pkg/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort    string `env:"GRPC_PORT" env-default:":50051"`
	MigrationPath string `env:"MIGRATION_PATH" env-default:":file://migrations"`

	postgres.PGConfig
	kafka.KafkaConfig
}


func ParseConfigFromEnv() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return &cfg, nil
}
