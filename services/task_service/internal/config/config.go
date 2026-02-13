package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"task_service/pkg/pool"
)

type Config struct {
	pool.PostgresCfg

	ServerPort       string `env:"TASK_PORT" env-default:"50051"`
	GroupServiceAddr string `env:"GROUP_SERV_ADDR" env-default:"localhost:50051"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
