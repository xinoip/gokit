package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
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
		"bearer": {
			Type:   "http",
			Scheme: "bearer",
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

	middlewares := e.Middlewares
	securitySchema := []map[string][]string{}

	if e.Secure {
		middlewares = append(middlewares, r.SecureMiddleware)
		securitySchema = []map[string][]string{
			{"bearer": {}},
		}
	}

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

// CreateSpecFiles creates the OpenAPI spec files in the current working
// directory.
func (r *Registry) CreateSpecFiles() error {
	jsonData, err := r.HumaAPI.OpenAPI().MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal OpenAPI to JSON: %w", err)
	}

	var fileMode os.FileMode = 0644

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	jsonFilePath := filepath.Join(wd, "openapi.json")
	err = os.WriteFile(jsonFilePath, jsonData, fileMode)
	if err != nil {
		return fmt.Errorf("failed to write openapi.json: %w", err)
	}

	yamlData, err := r.HumaAPI.OpenAPI().YAML()
	if err != nil {
		return fmt.Errorf("failed to marshal OpenAPI to YAML: %w", err)
	}

	yamlFilePath := filepath.Join(wd, "openapi.yaml")
	err = os.WriteFile(yamlFilePath, yamlData, fileMode)
	if err != nil {
		return fmt.Errorf("failed to write openapi.yaml: %w", err)
	}

	slog.Info("OpenAPI spec files created", "json", jsonFilePath, "yaml", yamlFilePath)

	return nil
}
