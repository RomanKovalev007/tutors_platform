package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PGConfig struct {
	PGHost     string `env:"POSTGRES_HOST" env-default:"db"`
	PGPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	PGUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	PGPassword string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	PGName     string `env:"POSTGRES_DB" env-default:"user_db"`
	DSN        string  
}

func (c *PGConfig) FormatConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PGUser, c.PGPassword, c.PGHost, c.PGPort, c.PGName)
}

type PostgresDB struct {
	DB  *sql.DB
	DSN string
}

func NewPostgresDB(cfg PGConfig) (*PostgresDB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PGHost, cfg.PGPort, cfg.PGUser, cfg.PGPassword, cfg.PGName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres db: %v", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping postgrses db: %v", err)
	}

	// db.SetMaxOpenConns(25)
	// db.SetMaxIdleConns(5)
	// db.SetConnMaxLifetime(5 * 60)

	return &PostgresDB{
		DB:             db,
		DSN:            cfg.DSN,
	}, nil
}