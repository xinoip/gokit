package notes

var (
	_ noteCreator           = (*PostgresStore)(nil)
	_ noteUpdater           = (*PostgresStore)(nil)
	_ noteDeleteStore       = (*PostgresStore)(nil)
	_ noteReader            = (*PostgresStore)(nil)
	_ noteLister            = (*PostgresStore)(nil)
	_ noteCacheSetter       = (*RedisCache)(nil)
	_ noteCacheDeleter      = (*RedisCache)(nil)
	_ noteCacheGetterSetter = (*RedisCache)(nil)
)
