package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	RedisPort     string `env:"REDIS_PORT" env-default:"6379"`
	RedisHost     string `env:"REDIS_HOST" env-default:"localhost"`
	RedisDB       int    `env:"REDIS_DB" env-default:"0"`
	RedisPassword string `env:"REDIS_PASSWORD" env-default:""`
}

func NewRedisDB(cfg RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %v", err)
	}

	return rdb, nil
}
