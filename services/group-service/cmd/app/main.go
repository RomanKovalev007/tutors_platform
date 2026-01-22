package main

import (
	"context"
	"fmt"
	"group_service/internal/app"
	"group_service/internal/config"
	"group_service/pkg/migrator"
)

func main() {
	ctx := context.Background()

	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("could not parse config: %w", err))
	}

	dsn := cfg.FormatConnectionString()

	err = migrator.RunMigrations(cfg.MigrationPath, dsn)
	if err != nil {
		panic(fmt.Errorf("failed to run migrations: %w", err))
	}

	app, err := app.NewApp(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to create app structure: %w", err))
	}

	app.MustRun(ctx, cfg.GRPCPort, cfg.GHTimeout)
}
