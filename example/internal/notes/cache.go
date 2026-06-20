package notes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const noteCacheTTL = 10 * time.Minute

var errCacheMiss = errors.New("note cache miss")

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{
		rdb: rdb,
	}
}

func cacheKey(noteID uuid.UUID) string {
	return fmt.Sprintf("notes:%s", noteID)
}

func (c *RedisCache) GetNote(ctx context.Context, noteID uuid.UUID) (*Note, error) {
	if c.rdb == nil {
		return nil, errCacheMiss
	}

	data, err := c.rdb.Get(ctx, cacheKey(noteID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, errCacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("redis get note: %w", err)
	}

	var note Note
	err = json.Unmarshal(data, &note)
	if err != nil {
		return nil, fmt.Errorf("unmarshal cached note: %w", err)
	}

	return &note, nil
}

func (c *RedisCache) SetNote(ctx context.Context, note *Note) error {
	if c.rdb == nil {
		return nil
	}

	data, err := json.Marshal(note)
	if err != nil {
		return fmt.Errorf("marshal cached note: %w", err)
	}

	err = c.rdb.Set(ctx, cacheKey(note.NoteID), data, noteCacheTTL).Err()
	if err != nil {
		return fmt.Errorf("redis set note: %w", err)
	}

	return nil
}

func (c *RedisCache) DeleteNote(ctx context.Context, noteID uuid.UUID) error {
	if c.rdb == nil {
		return nil
	}

	err := c.rdb.Del(ctx, cacheKey(noteID)).Err()
	if err != nil {
		return fmt.Errorf("redis delete note: %w", err)
	}

	return nil
}
