package pool

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCfg struct {
	Host     string `env:"POSTGRES_HOST"     env-default:"postgres-task"`
	Port     string `env:"POSTGRES_PORT"     env-default:"5432"`
	User     string `env:"POSTGRES_USER"     env-default:"task_postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"task_postgres"`
	DBName   string `env:"POSTGRES_DB"       env-default:"task_postgres"`
}

func CreateNewPool(cfg PostgresCfg) (*pgxpool.Pool, error) {
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.User, cfg.Password, addr, cfg.DBName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %v", err)
	}

	return pool, nil
}
