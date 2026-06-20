// Package api provides a fluent and convenient API for defining endpoints,
// using Huma and Chi as the backend.
package api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

// Endpoint represents a customizable HTTP API path. This can be  used to have
// full control over how the endpoint is configured.
type Endpoint struct {
	// Name is important because it will be used as the OperationID in OpenAPI
	// specification. It needs to be unique but also brief because code
	// generators use this to generate resulting function names for the
	// endpoint client.
	Name string

	// Secure endpoints are endpoints that has an authentication requirement.
	// They support bearer auth by default. [Registry.SecureMiddleware] is
	// automatically attached to these endpoints ensuring an authentication
	// flow is in place.
	Secure bool

	Method      string
	Path        string
	Middlewares huma.Middlewares
}

// Handler is the function that can handle the request in a type-safe manner.
type Handler[TReq, TRes any] func(context.Context, *TReq) (*TRes, error)

// Middleware is the function that can intercept a request and modify the context.
type Middleware func(ctx huma.Context, next func(huma.Context))

// MiddlewareMaker is a function that can create a middleware from a Huma API.
type MiddlewareMaker func(humaAPI huma.API) Middleware
