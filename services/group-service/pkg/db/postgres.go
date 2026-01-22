package postgres

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DbName   string `env:"POSTGRES_DB"`
}

func (c *PostgresConfig) FormatConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username, c.Password, c.Host, c.Port, c.DbName)
}

type Database struct {
	Pool *pgxpool.Pool
}

func New(config PostgresConfig) (*Database, error) {
	dataSource := fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		config.Username,
		url.QueryEscape(config.Password),
		config.Host,
		config.Port,
		config.DbName)

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &Database{
		Pool: pool,
	}, nil
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
