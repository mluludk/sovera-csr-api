package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations executes DDL migration files from db/migrations directory into PostgreSQL.
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	upFilePath := filepath.Join(migrationsDir, "000001_init_schema.up.sql")
	sqlBytes, err := os.ReadFile(upFilePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", upFilePath, err)
	}

	_, err = pool.Exec(ctx, string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute migration DDL: %w", err)
	}

	return nil
}
