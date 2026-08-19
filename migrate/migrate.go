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

// UpParams of the [Up] function.
type UpParams struct {
	Conn         *sql.DB
	MigrationsFS fs.FS
	Directory    string
	Dialect      goose.Dialect
}

// Up runs the migrations in the [directory] found in [migrationsFS] using the
// supplied [context]. It does not take ownership of [conn].
func Up(ctx context.Context, p UpParams) error {
	migrationRoot, err := fs.Sub(p.MigrationsFS, p.Directory)
	if err != nil {
		return fmt.Errorf("open migrations directory: %w", err)
	}

	provider, err := goose.NewProvider(p.Dialect, p.Conn, migrationRoot)
	if err != nil {
		return fmt.Errorf("create migration provider: %w", err)
	}

	_, err = provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

// UpURLParams of the [UpURLPostgres] function.
type UpURLParams struct {
	ConnURL      string
	MigrationsFS fs.FS
}

// UpURLPostgres is an opinionated helper that, opens a database connection,
// runs [Up] against directory named 'migrations' and closes the connection
// before returning. Only Postgres is supported.
func UpURLPostgres(ctx context.Context, p UpURLParams) (retErr error) {
	db, err := sql.Open("pgx", p.ConnURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		closeErr := db.Close()
		if closeErr != nil {
			retErr = errors.Join(retErr, fmt.Errorf("close database: %w", closeErr))
		}
	}()

	err = Up(ctx, UpParams{
		Conn:         db,
		MigrationsFS: p.MigrationsFS,
		Dialect:      goose.DialectPostgres,
		Directory:    "migrations",
	})
	if err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	return nil
}
