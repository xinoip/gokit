// Package migrate contains functions for managing database migrations.
package migrate

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Up runs the migrations in the "migrations" directory in given
// [migrationsFS], using connection [conn].
func Up(conn *sql.DB, migrationsFS fs.FS) error {
	goose.SetBaseFS(migrationsFS)
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	err = goose.Up(conn, "migrations")
	if err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}

// UpURL opens the connection to the database using [connURL] and runs [Up].
func UpURL(connURL string, migrationsFS fs.FS) error {
	db, err := sql.Open("pgx", connURL)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}

	err = Up(db, migrationsFS)
	if err != nil {
		return fmt.Errorf("Up: %w", err)
	}

	err = db.Close()
	if err != nil {
		return fmt.Errorf("db.Close: %w", err)
	}

	return nil
}
