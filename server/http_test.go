package server_test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xinoip/gokit/server"
)

func TestServeHTTPReturnsListenError(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	listener, err := new(net.ListenConfig).Listen(ctx, "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, listener.Close())
	})

	err = server.ServeHTTP(
		ctx,
		server.DefaultHTTPConfig(listener.Addr().String(), http.NotFoundHandler()),
	)
	require.Error(t, err)
	assert.ErrorContains(t, err, "serve HTTP")
}

func TestServeHTTPShutsDownWhenContextIsCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	result := make(chan error, 1)
	go func() {
		result <- server.ServeHTTP(
			ctx,
			server.DefaultHTTPConfig("127.0.0.1:0", http.NotFoundHandler()),
		)
	}()

	select {
	case err := <-result:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("ServeHTTP did not return after cancellation")
	}
}
