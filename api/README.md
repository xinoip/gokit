# Pio's Go Kit: API

Provides convenience functions for defining APIs with OpenAPI support. It is
powered by [Huma](https://github.com/danielgtaylor/huma) and [Chi](https://github.com/go-chi/chi).

See the package documentation for runnable examples covering registry creation,
endpoint registration, and OpenAPI generation.

## RPC

This package also provides an opinionated way to create and handle JSON-RPC-like
endpoints. These endpoints are not fully compliant with the JSON-RPC
specification. They accept parameters in a JSON body and return results in a
JSON body. Returned HTTP status codes are highly opinionated.
