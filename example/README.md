# Pio's Go Kit: Example

This is an opinionated example project using `gokit`.

- Notes CRUD API.
- Postgres migrations and sqlc-generated persistence.
- Redis-backed note read cache.

## Run

```sh
make up-dependencies
make migrate
make run
```

