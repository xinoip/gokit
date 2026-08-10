package notes

import (
	"example/internal/gen"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type RPC struct {
	store *PostgresStore
	cache *RedisCache
}

func NewRPC(pgdb *pgxpool.Pool, rdb *redis.Client) *RPC {
	store := NewPostgresStore(gen.New(pgdb))
	cache := NewRedisCache(rdb)
	return &RPC{
		store: store,
		cache: cache,
	}
}
