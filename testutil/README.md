# Pio's Go Kit: Test Util

Test utilities supporting the packages in `gokit`.

`NewTestAPIRegistry` creates an in-memory Huma API suitable for endpoint tests.
`NewPostgres` and `NewRedis` require the local services used by the repository's
Docker Compose setup. PostgreSQL databases are isolated per test; Redis uses
database 15 and should not be used concurrently with unrelated test suites.
