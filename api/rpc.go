package api

import (
	"context"

	"github.com/xinoip/gokit/errs"
)

// RPCBody is a generic wrapper intended to be used for RPC parameters to put
// them in HTTP body.
type RPCBody[T any] struct {
	Body *T
}

// MakeRPCHandler creates the opinionated HTTP handler for an RPC endpoint.
// Passed in [rpc] handler must utilize errors defined in [errs] package.
func MakeRPCHandler[TReq, TRes any](rpc Handler[TReq, TRes]) Handler[RPCBody[TReq], RPCBody[TRes]] {
	return func(ctx context.Context, req *RPCBody[TReq]) (*RPCBody[TRes], error) {
		res, err := rpc(ctx, req.Body)
		if err != nil {
			return nil, errs.HTTPError(err)
		}

		return &RPCBody[TRes]{
			Body: res,
		}, nil
	}
}
