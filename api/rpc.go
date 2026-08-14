package api

import (
	"context"
	"errors"
	"fmt"
)

// RPCBody is a generic wrapper intended to be used for RPC parameters to put
// them in HTTP body.
type RPCBody[T any] struct {
	Body *T
}

// RPCError indicates that the RPC handler failed. These errors send out as
// HTTP 400 response.
var RPCError = errors.New("rpc error")

// RPCErrorf wraps a formatted error with [RPCError].
func RPCErrorf(format string, a ...any) error {
	return fmt.Errorf("%w: %w", RPCError, fmt.Errorf(format, a...))
}

// RPCErrorw wraps an error with [RPCError].
func RPCErrorw(err error) error {
	return fmt.Errorf("%w: %w", RPCError, err)
}

// RPCErrorToHTTPError maps known errors to HTTP errors.
func RPCErrorToHTTPError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, RPCError):
		return HTTP400(err)
	default:
		return HTTP500(err)
	}
}

// MakeRPCHandler creates the opinionated HTTP handler for an RPC endpoint.
// Resulting endpoint handler will return these HTTP responses:
//
// - 400 Bad Request if [rpc] handler returns [RPCError].
// - 500 Internal Server Error if [rpc] handler returns any error that is not [RPCError].
// - 200 OK if [rpc] handler does not return any error.
func MakeRPCHandler[TReq, TRes any](rpc Handler[TReq, TRes]) Handler[RPCBody[TReq], RPCBody[TRes]] {
	return func(ctx context.Context, req *RPCBody[TReq]) (*RPCBody[TRes], error) {
		res, err := rpc(ctx, req.Body)
		if err != nil {
			return nil, RPCErrorToHTTPError(err)
		}

		return &RPCBody[TRes]{
			Body: res,
		}, nil
	}
}
