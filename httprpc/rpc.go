// Package httprpc provides an opinionated way to create and handle JSON-RPC-like
// endpoints using the [httpapi] package.
package httprpc

import (
	"context"
	"errors"

	"github.com/xinoip/gokit/httpapi"
)

// Request represents the request for an RPC endpoint.
type Request[T any] struct {
	Body *T
}

// Response represents the response for an RPC endpoint.
type Response[T any] struct {
	Body *T
}

// Failure indicates an expected error happened during RPC handler execution.
// These errors are treated safe to expose to public and will be translated to
// HTTP 400 responses.
//
//nolint:errname // Failure is the RPC domain term and part of the public API.
type Failure string

// Error implements the error interface.
func (f Failure) Error() string {
	return string(f)
}

func errToHTTPError(err error) error {
	if err == nil {
		return nil
	}

	// Only extract failure from error so that we don't expose any other errors.
	failure, ok := errors.AsType[Failure](err)
	if ok {
		return httpapi.HTTP400(failure)
	}

	return httpapi.HTTP500(err)
}

// MakeHandler creates the opinionated HTTP handler for an RPC endpoint.
// Resulting endpoint handler will return these HTTP responses:
//
// - 400: Handler returns [Failure].
// - 500: Handler returns any error that is not [Failure].
// - 200: Handler does not return any error.
//
// Requests and responses are JSON-encoded and expected to be in the HTTP body.
func MakeHandler[TReq, TRes any](rpc httpapi.Handler[TReq, TRes]) httpapi.Handler[Request[TReq], Response[TRes]] {
	return func(ctx context.Context, req *Request[TReq]) (*Response[TRes], error) {
		res, err := rpc(ctx, req.Body)
		if err != nil {
			return nil, errToHTTPError(err)
		}

		if res == nil {
			return nil, nil
		}

		return &Response[TRes]{
			Body: res,
		}, nil
	}
}
