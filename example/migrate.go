package main

import (
	"context"
	"embed"

	"github.com/xinoip/gokit/migrate"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func runMigrate(ctx context.Context, c *Config) error {
	return migrate.UpURLPostgres(ctx, migrate.UpURLParams{
		ConnURL:      c.PostgresConnURL,
		MigrationsFS: MigrationsFS,
	})
}
