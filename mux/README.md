# Pio's Go Kit: Mux

Provides read to use implementations for various mux backends. It currently only
supports [chi](https://github.com/go-chi/chi).

## Example

```go
myChiMux := mux.NewChi(&mux.NewChiParams{
    AllowOrigins: []string{"http://*", "https://*"},
})

// It has all the sane defaults and middlewares setup.
// Use it as you would any other chi router.
```
