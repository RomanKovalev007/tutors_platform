package config

import (
	"auth_service/pkg/kafka"
	"auth_service/pkg/postgres"
	"auth_service/pkg/redis"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPCPort    string `env:"GRPC_PORT" env-default:":50051"`
	MigrationPath string `env:"MIGRATION_PATH" env-default:":file:///migrations"`

	postgres.PGConfig
	redis.RedisConfig
	kafka.KafkaConfig
	TokenConfig
}

type TokenConfig struct {
	AccessTTL  int32 `env:"ACCESS_TTL_M" env-default:"5"` 
	RefreshTTL int32 `env:"REFRESH_TTL_H" env-default:"336"`
	Secret     []byte     `env:"SECRET,required"`   
}

func ParseConfigFromEnv() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return &cfg, nil
}
