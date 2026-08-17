package main

import (
	handlersv1 "example/internal/handlers/v1"

	"github.com/danielgtaylor/huma/v2"
	"github.com/xinoip/gokit/httpapi"
	"github.com/xinoip/gokit/mux"
)

func createOpenAPI() error {
	r := mux.NewChi(&mux.NewChiParams{
		AllowOrigins: []string{"https://*"},
	})

	apir := httpapi.NewRegistry(&httpapi.NewRegistryParams{
		Mux:     r,
		Title:   "Notes API",
		Version: "v0.0.1",
		SecureMiddlewareMaker: func(_ huma.API) httpapi.Middleware {
			return func(ctx huma.Context, next func(huma.Context)) {
				next(ctx)
			}
		},
	})

	new(handlersv1.Handlers).Register(apir)

	err := apir.CreateJSONSpecFile("openapi.json")
	if err != nil {
		return err
	}

	return apir.CreateYAMLSpecFile("openapi.yaml")
}
