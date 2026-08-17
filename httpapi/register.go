package httpapi

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

const (
	bearerScheme = "bearer"
)

// Registry holds an API definition and is needed for registering endpoints.
type Registry struct {
	// HumaAPI is the handle to the Huma backend.
	HumaAPI huma.API

	// SecureMiddleware is automatically attached to secure endpoints. See [Endpoint.Secure].
	SecureMiddleware Middleware
}

// NewRegistryParams of [NewRegistry].
type NewRegistryParams struct {
	Mux                   chi.Router
	Title                 string
	Version               string
	SecureMiddlewareMaker MiddlewareMaker
}

// NewRegistry wires up the Huma backend with the given Mux and creates a new
// [Registry] that can be used to register endpoints.
func NewRegistry(p *NewRegistryParams) *Registry {
	humaConfig := huma.DefaultConfig(p.Title, p.Version)
	humaConfig.DocsPath = "/swagger"
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		bearerScheme: {
			Type:   "http",
			Scheme: bearerScheme,
		},
	}
	humaAPI := humachi.New(p.Mux, humaConfig)

	secureMiddleware := p.SecureMiddlewareMaker(humaAPI)

	return &Registry{
		HumaAPI:          humaAPI,
		SecureMiddleware: secureMiddleware,
	}
}

// Register a new endpoint to the registry, applying additional options.
func Register[TReq, TRes any](r *Registry, pe *Endpoint, h Handler[TReq, TRes], opts ...Option) {
	e := *pe
	for _, opt := range opts {
		opt(&e)
	}

	middlewares := huma.Middlewares{}
	securitySchema := []map[string][]string{}

	if e.Secure {
		middlewares = append(middlewares, r.SecureMiddleware)
		securitySchema = []map[string][]string{
			{bearerScheme: {}},
		}
	}

	middlewares = append(middlewares, e.Middlewares...)

	//nolint:exhaustruct // [huma.Operation]
	huma.Register(r.HumaAPI, huma.Operation{
		OperationID:   e.Name,
		Method:        e.Method,
		Path:          e.Path,
		DefaultStatus: http.StatusOK,
		Middlewares:   middlewares,
		Security:      securitySchema,
	}, h)
}
