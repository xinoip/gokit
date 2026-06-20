# Pio's Go Kit: API

Provides convenience functions for defining an API that is supporting OpenAPI
specification. It is powered by [Huma](https://github.com/danielgtaylor/huma) and [Chi](https://github.com/go-chi/chi).

## Example

```go
myChiMux := ... // Prepare a mux with Chi beforehand.
r := api.NewRegistry(api.NewRegistryParams{
    Mux:         myChiMux,
    Title:       "My API",
    Version:     "v0.0.1",
    SecureMiddlewareMaker: func(humaAPI huma.API) Middleware {
        return func(ctx huma.Context, next func(huma.Context)) {
            // You can use passed humaAPI in MiddlewareMaker to create your middlewares.

            // Authenticate the request and set the user in the context.
            next(ctx)
        }
    }
})

api.Get(r, "/things", GetThingsHandler, "get-things", api.WithInsecure())
api.Post(r, "/things", PostThingsHandler, "post-things")
api.Put(r, "/things/{thing-id}", PutThingsHandler, "put-things-by-thing-id", api.WithMiddlewares(m1, m2))
api.Delete(r, "/things/{thing-id}", DeleteThingsHandler, "delete-things-by-thing-id")
api.Patch(r, "/things/{thing-id}", PatchThingsHandler, "patch-things-by-thing-id")

err := r.CreateJSONSpecFile("openapi.json")
err := r.CreateYAMLSpecFile("openapi.yaml")

// Serve your myChiMux as you wish. It now has all the endpoints registered.
```
