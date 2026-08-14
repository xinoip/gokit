package api_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xinoip/gokit/api"
)

type rpcRequest struct {
	Value string `json:"value"`
}

type rpcResponse struct {
	Value string `json:"value"`
}

func TestMakeRPCHandler(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	handler := api.MakeRPCHandler(func(_ context.Context, req *rpcRequest) (*rpcResponse, error) {
		return &rpcResponse{Value: req.Value + "!"}, nil
	})

	res, err := handler(ctx, &api.RPCBody[rpcRequest]{Body: &rpcRequest{Value: "ok"}})
	require.NoError(t, err)
	require.Equal(t, "ok!", res.Body.Value)
}
