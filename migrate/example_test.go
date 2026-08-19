package migrate_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing/fstest"

	"github.com/pressly/goose/v3"
	"github.com/xinoip/gokit/migrate"

	_ "modernc.org/sqlite"
)

//nolint:gochecknoglobals // For example.
var migrationsFS = fstest.MapFS{
	"migrations/001_create_users.sql": {
		Data: []byte(`
-- +goose Up
CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
`),
	},
}

func ExampleUp() {
	ctx := context.Background()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	err = migrate.Up(ctx, migrate.UpParams{
		Conn:         db,
		MigrationsFS: migrationsFS,
		Dialect:      goose.DialectSQLite3,
		Directory:    "migrations",
	})
	if err != nil {
		panic(err)
	}

	var table string
	err = db.QueryRowContext(ctx, `
		SELECT name
		FROM sqlite_master
		WHERE type = 'table' AND name = 'users'
	`).Scan(&table)
	if err != nil {
		panic(err)
	}

	fmt.Println(table)

	// Output:
	// users
}

//nolint:testableexamples // Needs a real Postgres instance.
func ExampleUpURLPostgres() {
	ctx := context.Background()

	err := migrate.UpURLPostgres(ctx, migrate.UpURLParams{
		ConnURL:      "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
		MigrationsFS: migrationsFS,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("migrations applied")
}
