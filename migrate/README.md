# Pio's Go Kit: Migrate

Provides quick and easy way to run database migrations. Powered by [Goose](https://github.com/pressly/goose).

## Example

```go
// Use already opened database connection.
err := migrate.Up(db, migrationsFS)

// Or create a new connection on the fly.
err := UpURL(connURL, migrationsFS)
```
