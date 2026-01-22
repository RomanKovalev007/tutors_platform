package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string `env:"REDIS_CACHE_HOST" env-default:"redis-cache"`
	Port     string `env:"REDIS_CACHE_PORT" env-default:"6379"`
	DB       int    `env:"REDIS_CACHE_DB" env-default:"0"`
	Password string `env:"REDIS_CACHE_PASSWORD" env-default:""`
}

type Cache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func NewCache(cfg Config, prefix string, ttl time.Duration) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Cache{
		client: client,
		prefix: prefix,
		ttl:    ttl,
	}, nil
}

func (c *Cache) key(id string) string {
	return fmt.Sprintf("%s:%s", c.prefix, id)
}

func (c *Cache) Set(ctx context.Context, id string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, c.key(id), data, c.ttl).Err()
}

func (c *Cache) Get(ctx context.Context, id string, dest interface{}) error {
	data, err := c.client.Get(ctx, c.key(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

func (c *Cache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, c.key(id)).Err()
}

func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, c.key(pattern), 100).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}

type CacheStats struct {
	Hits   int64
	Misses int64
}

func (c *Cache) Stats() *CacheStats {
	return &CacheStats{}
}

var ErrCacheMiss = fmt.Errorf("cache miss")
