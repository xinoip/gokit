# Pio's Go Kit: Server

Provides default implementations for HTTP servers with graceful shutdown support.

## Example

```go
myRouter := ... // Setup any mux here.
myServer := server.NewHTTPServer(&server.NewHTTPServerParams{
    Addr: ":8080",
    Handler: myRouter,
})
```
