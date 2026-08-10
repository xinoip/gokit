package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	ServeAddr       string `env:"SERVE_ADDR"`
	PostgresConnURL string `env:"POSTGRES_CONN_URL"`
	RedisConnURL    string `env:"REDIS_CONN_URL"`
}

func main() {
	err := Run()
	if err != nil {
		slog.Error("Fatal Error:", "err", err.Error())
		os.Exit(1)
	}
}

func Run() error {
	ctx, cancel := signal.NotifyContext(context.Background())
	defer cancel()

	if len(os.Args) != 2 {
		return errors.New("invalid number of arguments")
	}

	var c Config
	err := envconfig.Process(ctx, &c)
	if err != nil {
		return fmt.Errorf("failed to load config from env: %w", err)
	}

	switch os.Args[1] {
	case "serve":
		return serve(ctx, &c)
	case "migrate":
		return runMigrate(&c)
	case "openapi":
		return createOpenAPI()
	default:
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}
