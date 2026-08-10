package testutil

import (
	"io/fs"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/goosemigrator"
	"github.com/stretchr/testify/require"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	postgresUser     = "postgres"
	postgresPassword = "postgres"
	postgresHost     = "localhost"
	postgresPort     = "5432"
)

// NewPostgres creates a new database for testing using [pgtestdb]. It is basically a
// template based, migrations supported, fast Postgres instance.
func NewPostgres(t *testing.T, migrationsFS fs.FS) *pgx.Conn {
	t.Helper()
	ctx := t.Context()

	dbconf := pgtestdb.Config{
		DriverName:                "pgx",
		User:                      postgresUser,
		Password:                  postgresPassword,
		Host:                      postgresHost,
		Port:                      postgresPort,
		Options:                   "sslmode=disable",
		Database:                  "postgres",
		TestRole:                  nil,
		ForceTerminateConnections: false,
	}

	m := goosemigrator.New(
		"migrations",
		goosemigrator.WithFS(migrationsFS),
	)

	conf := pgtestdb.Custom(t, dbconf, m)
	require.NotEqual(t, dbconf, conf)

	conn, err := pgx.Connect(ctx, conf.URL())
	require.NoError(t, err)
	t.Cleanup(func() {
		err := conn.Close(ctx)
		require.NoError(t, err)
	})

	return conn
}
