package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

const (
	redisDatabaseNumber = 15
)

// NewRedis creates a Redis client for tests using local Redis database 15.
func NewRedis(t *testing.T) *redis.Client {
	t.Helper()

	client := redis.NewClient(
		//nolint:exhaustruct // [redis.Options]
		&redis.Options{
			Addr: "localhost:6379",
			DB:   redisDatabaseNumber,
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	require.NoError(t, client.Ping(ctx).Err())

	t.Cleanup(func() {
		closeErr := client.Close()
		require.NoError(t, closeErr)
	})

	return client
}
