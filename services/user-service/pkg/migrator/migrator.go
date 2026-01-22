package migrator

import (
	"fmt"
	"os"

    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4"
)


func RunMigrations(migPath, connStr string) error {
	migrationsPath := migPath

	migration, err := migrate.New(
		migrationsPath,
		connStr,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}
	defer migration.Close()

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	return nil
}

func GetMigrationInfo(connStr string) (version uint, dirty bool, err error) {
	workDir, err := os.Getwd()
	if err != nil {
		return 0, false, err
	}

	migration, err := migrate.New(
		fmt.Sprintf("file://%s/migrations", workDir),
		connStr,
	)
	if err != nil {
		return 0, false, err
	}
	defer migration.Close()

	version, dirty, err = migration.Version()
	return version, dirty, err
}
