package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTPPort     int           `env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"30s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"30s"`

	GHTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT" env-default:"15s"`

	AuthGRPC  string `env:"AUTH_GRPC" env-default:"auth-go:50051"`
	GroupGRPC string `env:"GROUP_GRPC" env-default:"group-go:50051"`
	UserGRPC  string `env:"USER_GRPC" env-default:"user-go:50051"`
	TasksGRPC string `env:"TASKS_GRPC" env-default:"tasks-go:50051"`
}

func ParseConfigFromEnv() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return &cfg, nil
}
