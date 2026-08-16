# Pio's Go Kit: HTTP RPC

This package can't be used without `gokit/api` package. It provides convenience
functions for defining opinionated JSON-RPC like HTTP APIs.

## Principles

It's built upon these principles, where an RPC endpoint is:

- A remote procedure call in disguise.
- Only `POST` request.
- Parameters received from JSON encoded HTTP body.
- Responses returned in JSON encoded HTTP body.
- Can only return 200, 400 or 500 HTTP status codes.
- 200 is returned when no error returned from handler.
- 400 is returned when `rpc.Failure` is returned.
- Each HTTP endpoint is a remote procedure call in disguise.
- RPC endpoints can only return 200, 400 or 500.
- URL paths are in form of `/rpc/<name_of_procedure>`.
