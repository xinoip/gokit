package main

import (
	"context"
	"fmt"
	"log/slog"

	handlersv1 "example/internal/handlers/v1"
	"example/internal/notes"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/xinoip/gokit/httpapi"
	"github.com/xinoip/gokit/mux"
	"github.com/xinoip/gokit/server"
)

func serve(ctx context.Context, c *Config) error {
	r := mux.NewChi(&mux.NewChiParams{
		AllowOrigins: []string{"https://*"},
	})

	pgdb, err := pgxpool.New(ctx, c.PostgresConnURL)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer pgdb.Close()

	err = pgdb.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: c.RedisConnURL,
	})
	defer func() {
		err := rdb.Close()
		if err != nil {
			slog.Warn("Failed to close redis client", "err", err)
		}
	}()

	err = rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	apir := httpapi.NewRegistry(&httpapi.NewRegistryParams{
		Mux:     r,
		Title:   "Notes API",
		Version: "v0.0.1",
		SecureMiddlewareMaker: func(_ huma.API) httpapi.Middleware {
			return func(ctx huma.Context, next func(huma.Context)) {
				next(ctx)
			}
		},
	})

	v1Handlers := handlersv1.Handlers{
		RPCNotes: notes.NewRPC(pgdb, rdb),
	}
	v1Handlers.Register(apir)

	return server.ServeHTTP(ctx, &server.ServeHTTPParams{
		Addr:    c.ServeAddr,
		Handler: r,
	})
}
