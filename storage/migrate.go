package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

const migrationsTableName = "migration"

// Migrate runs migrations from migrationsDir.
func Migrate(
	ctx context.Context,
	pool *pgxpool.Pool,
	migrationsDir string,
) (version int32, err error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to get connection from pool: %w", err)
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), migrationsTableName)
	if err != nil {
		return 0, fmt.Errorf("unable to create migrator: %w", err)
	}
	if err = migrator.LoadMigrations(migrationsDir); err != nil {
		return 0, fmt.Errorf("unable to load migrations: %w", err)
	}
	if err = migrator.Migrate(ctx); err != nil {
		return 0, fmt.Errorf("unable to migrate: %w", err)
	}
	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to get  migration version: %w", err)
	}
	return ver, nil
}
