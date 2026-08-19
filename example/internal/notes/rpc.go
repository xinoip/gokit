package notes

import (
	"context"
	"errors"
	"log/slog"

	"example/internal/gen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type RPC struct {
	store NoteStore
	cache NoteCache
}

// NoteStore persists notes for RPC operations.
type NoteStore interface {
	CreateNote(ctx context.Context, params CreateNoteParams) (*Note, error)
	UpdateNote(ctx context.Context, params UpdateNoteParams) (*Note, error)
	DeleteNote(ctx context.Context, noteID uuid.UUID) error
	GetNote(ctx context.Context, noteID uuid.UUID) (*Note, error)
	ListNotes(ctx context.Context) ([]Note, error)
}

// NoteCache caches notes for RPC operations.
type NoteCache interface {
	GetNote(ctx context.Context, noteID uuid.UUID) (*Note, error)
	SetNote(ctx context.Context, note *Note) error
	DeleteNote(ctx context.Context, noteID uuid.UUID) error
}

func logCacheError(action string, err error) {
	if err != nil {
		slog.Warn("Cache operation failed", "action", action, "err", err)
	}
}

// NewRPC constructs the notes RPC service from database clients.
func NewRPC(pgdb *pgxpool.Pool, rdb *redis.Client) (*RPC, error) {
	if pgdb == nil {
		return nil, errors.New("create notes RPC: Postgres pool must not be nil")
	}

	store := NewPostgresStore(gen.New(pgdb))
	cache := NewRedisCache(rdb)

	return NewRPCWithDependencies(store, cache)
}

// NewRPCWithDependencies constructs an RPC service from explicit dependencies.
func NewRPCWithDependencies(store NoteStore, cache NoteCache) (*RPC, error) {
	if store == nil {
		return nil, errors.New("create notes RPC: store must not be nil")
	}
	if cache == nil {
		return nil, errors.New("create notes RPC: cache must not be nil")
	}

	return &RPC{
		store: store,
		cache: cache,
	}, nil
}
