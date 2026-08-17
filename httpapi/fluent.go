package httpapi

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// Get is a fluent helper for registering a secure GET endpoint. Use [Option]
// to further customize the endpoint. These fluent helpers combined with
// available options is generally enough for most API registries but if you
// need more control over the endpoint checkout [Register] and [Endpoint].
func Get[TReq, TRes any](r *Registry, path string, handler Handler[TReq, TRes], name string, opts ...Option) {
	Register(r, &Endpoint{
		Name:        name,
		Secure:      true,
		Method:      http.MethodGet,
		Path:        path,
		Middlewares: huma.Middlewares{},
	}, handler, opts...)
}

// Post is same as [GET] but for POST.
func Post[TReq, TRes any](r *Registry, path string, handler Handler[TReq, TRes], name string, opts ...Option) {
	Register(r, &Endpoint{
		Name:        name,
		Secure:      true,
		Method:      http.MethodPost,
		Path:        path,
		Middlewares: huma.Middlewares{},
	}, handler, opts...)
}

// Put is same as [GET] but for PUT.
func Put[TReq, TRes any](r *Registry, path string, handler Handler[TReq, TRes], name string, opts ...Option) {
	Register(r, &Endpoint{
		Name:        name,
		Secure:      true,
		Method:      http.MethodPut,
		Path:        path,
		Middlewares: huma.Middlewares{},
	}, handler, opts...)
}

// Delete is same as [GET] but for DELETE.
func Delete[TReq, TRes any](r *Registry, path string, handler Handler[TReq, TRes], name string, opts ...Option) {
	Register(r, &Endpoint{
		Name:        name,
		Secure:      true,
		Method:      http.MethodDelete,
		Path:        path,
		Middlewares: huma.Middlewares{},
	}, handler, opts...)
}

// Patch is same as [GET] but for PATCH.
func Patch[TReq, TRes any](r *Registry, path string, handler Handler[TReq, TRes], name string, opts ...Option) {
	Register(r, &Endpoint{
		Name:        name,
		Secure:      true,
		Method:      http.MethodPatch,
		Path:        path,
		Middlewares: huma.Middlewares{},
	}, handler, opts...)
}
