package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type tokenRepository struct {
	client *redis.Client
	prefix string
}

func NewRedisTokenRepository(client *redis.Client, prefix string) *tokenRepository {
	if prefix == "" {
		prefix = "auth:refresh_token:"
	}
	return &tokenRepository{client: client, prefix: prefix}
}

func (r *tokenRepository) key(token string) string {
	return fmt.Sprintf("%s%s", r.prefix, token)
}

func (r *tokenRepository) SaveToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	ttl := time.Until(expiresAt)
	return r.client.Set(ctx, r.key(token), userID, ttl).Err()
}

func (r *tokenRepository) GetToken(ctx context.Context, token string) (string, error) {
	val, err := r.client.Get(ctx, r.key(token)).Result()
	if err == redis.Nil {
		return "", ErrUserNotFound
	}
	return val, err
}

func (r *tokenRepository) DeleteToken(ctx context.Context, token string) error {
	err := r.client.Del(ctx, r.key(token)).Err()
	if err == redis.Nil {
		return ErrUserNotFound
	}
	return err
}

