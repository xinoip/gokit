package api_test

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
	"github.com/xinoip/gokit/api"
)

type exampleRequest struct{}

type exampleResponse struct{}

func exampleHandler(context.Context, *exampleRequest) (*exampleResponse, error) {
	return &exampleResponse{}, nil
}

func exampleMiddleware(ctx huma.Context, next func(huma.Context)) {
	// Dummy middleware that does nothing.
	next(ctx)
}

func ExampleNewRegistry() {
	r := api.NewRegistry(&api.NewRegistryParams{
		Mux:     chi.NewRouter(),
		Title:   "Things API",
		Version: "v1.0.0",
		SecureMiddlewareMaker: func(_ huma.API) api.Middleware {
			return func(ctx huma.Context, next func(huma.Context)) {
				// Implement your authentication middleware here.
				next(ctx)
			}
		},
	})

	// Supports all HTTP methods.
	api.Get(r, "/things", exampleHandler, "get-things", api.WithInsecure())
	api.Post(r, "/things", exampleHandler, "create-things", api.WithInsecure())
	api.Put(r, "/things", exampleHandler, "update-things", api.WithInsecure())
	api.Delete(r, "/things", exampleHandler, "delete-things", api.WithInsecure())
	api.Patch(r, "/things", exampleHandler, "patch-things", api.WithInsecure())

	// Can modify with options.
	api.Get(r, "/things/custom", exampleHandler, "get-things-custom", api.WithMiddlewares(exampleMiddleware))

	// Endpoints are secure by default, meaning SecureMiddleware is always active.
	api.Get(r, "/things", exampleHandler, "get-things-secure")

	// OpenAPI spec is already supported and can be generated.
	err := r.CreateJSONSpecFile("/tmp/openapi.json")
	if err != nil {
		panic(err)
	}
	err = r.CreateYAMLSpecFile("/tmp/openapi.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println(r.HumaAPI.OpenAPI().Info.Title)

	// Output:
	// Things API
}
