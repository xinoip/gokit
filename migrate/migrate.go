// Package migrate contains functions for managing database migrations.
package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Up runs the migrations in the [directory] found in [migrationsFS] using the
// supplied [context]. It does not take ownership of [conn].
func Up(ctx context.Context, conn *sql.DB, migrationsFS fs.FS, directory string) error {
	migrationRoot, err := fs.Sub(migrationsFS, directory)
	if err != nil {
		return fmt.Errorf("open migrations directory: %w", err)
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, conn, migrationRoot)
	if err != nil {
		return fmt.Errorf("create migration provider: %w", err)
	}

	_, err = provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

// UpURL is an opinionated helper that, opens a database connection, runs [Up]
// agains directory named 'migrations' and closes the connection before
// returning.
func UpURL(ctx context.Context, connURL string, migrationsFS fs.FS) (retErr error) {
	db, err := sql.Open("pgx", connURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		closeErr := db.Close()
		if closeErr != nil {
			retErr = errors.Join(retErr, fmt.Errorf("close database: %w", closeErr))
		}
	}()

	err = Up(ctx, db, migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	return nil
}
