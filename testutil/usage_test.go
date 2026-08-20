package testutil_test

import (
	"context"
	"net/http"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
	"github.com/xinoip/gokit/httpapi"
	"github.com/xinoip/gokit/testutil"
)

type exampleBody struct {
	Message string `json:"message"`
}

type exampleResponse struct {
	Body exampleBody
}

func TestHTTPAPIUsage(t *testing.T) {
	t.Parallel()

	registry, api := testutil.NewTestAPIRegistry(t)
	httpapi.Get(
		registry,
		"/greet",
		func(context.Context, *struct{}) (*exampleResponse, error) {
			return &exampleResponse{
				Body: exampleBody{Message: "hello"},
			}, nil
		},
		"greet",
	)

	response := api.Get("/greet")
	require.Equal(t, http.StatusOK, response.Code)

	body := testutil.UnmarshalResponseJSON[map[string]any](t, response)
	require.Equal(t, "hello", body["message"])
}

func TestContextUsage(t *testing.T) {
	t.Parallel()

	ctx := testutil.DefaultTestCtx(t)
	require.NoError(t, ctx.Err())
}

func TestPostgresUsage(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	migrationsFS := fstest.MapFS{
		"migrations/001_create_widgets.sql": {
			Data: []byte(`
-- +goose Up
CREATE TABLE widgets (id BIGINT PRIMARY KEY);

-- +goose Down
DROP TABLE widgets;
`),
		},
	}
	conn, connURL := testutil.NewPostgresWithConnURL(t, migrationsFS)
	require.NotEmpty(t, connURL)

	var table string
	err := conn.QueryRow(ctx, `
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'widgets'
    `).Scan(&table)
	require.NoError(t, err)
	require.Equal(t, "widgets", table)
}

func TestRedisUsage(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	client := testutil.NewRedis(t)
	err := client.Set(ctx, t.Name(), "cached", 0).Err()
	require.NoError(t, err)

	value, err := client.Get(ctx, t.Name()).Result()
	require.NoError(t, err)
	require.Equal(t, "cached", value)

	err = client.Del(ctx, t.Name()).Err()
	require.NoError(t, err)
}
