package main

import (
	handlersv1 "example/internal/handlers/v1"

	"github.com/danielgtaylor/huma/v2"
	"github.com/xinoip/gokit/api"
	"github.com/xinoip/gokit/mux"
)

func createOpenAPI() error {
	r := mux.NewChi(&mux.NewChiParams{
		AllowOrigins: []string{"https://*"},
	})

	apir := api.NewRegistry(&api.NewRegistryParams{
		Mux:     r,
		Title:   "Notes API",
		Version: "v0.0.1",
		SecureMiddlewareMaker: func(_ huma.API) api.Middleware {
			return func(ctx huma.Context, next func(huma.Context)) {
				next(ctx)
			}
		},
	})

	new(handlersv1.Handlers).Register(apir)

	return apir.CreateSpecFiles()
}
