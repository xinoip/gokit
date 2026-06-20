package main

import (
	"embed"

	"github.com/xinoip/gokit/migrate"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func runMigrate(c *Config) error {
	return migrate.UpURL(c.PostgresConnURL, MigrationsFS)
}
