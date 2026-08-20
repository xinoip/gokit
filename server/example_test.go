package server_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/xinoip/gokit/server"
)

func ExampleServeHTTP() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	listener, err := new(net.ListenConfig).Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	addr := listener.Addr().String()
	err = listener.Close()
	if err != nil {
		panic(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, writeErr := io.WriteString(w, "hello from gokit")
		if writeErr != nil {
			panic(writeErr)
		}
	})
	config := server.DefaultHTTPConfig(addr, handler)
	config.Logger = slog.New(slog.DiscardHandler)

	err = server.ServeHTTP(ctx, config)
	if err != nil {
		panic(err)
	}

	fmt.Println("Shutting down HTTP server")

	// Output:
	// Shutting down HTTP server
}
