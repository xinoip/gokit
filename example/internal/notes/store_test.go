package notes_test

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"example/internal/gen"
	"example/internal/notes"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/xinoip/gokit/testutil"
)

func TestPostgresStoreCRUD(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	store := newTestPostgresStore(t)

	first, err := store.CreateNote(ctx, notes.CreateNoteParams{
		NoteID: uuid.New(),
		Title:  "first title",
		Body:   "first body",
	})
	require.NoError(t, err)
	require.Equal(t, "first title", first.Title)
	require.Equal(t, "first body", first.Body)

	time.Sleep(time.Millisecond)

	second, err := store.CreateNote(ctx, notes.CreateNoteParams{
		NoteID: uuid.New(),
		Title:  "second title",
		Body:   "second body",
	})
	require.NoError(t, err)

	noteList, err := store.ListNotes(ctx)
	require.NoError(t, err)
	require.Len(t, noteList, 2)
	require.Equal(t, second.NoteID, noteList[0].NoteID)
	require.Equal(t, first.NoteID, noteList[1].NoteID)

	got, err := store.GetNote(ctx, first.NoteID)
	require.NoError(t, err)
	require.Equal(t, first.NoteID, got.NoteID)
	require.Equal(t, first.Title, got.Title)
	require.Equal(t, first.Body, got.Body)

	updated, err := store.UpdateNote(ctx, notes.UpdateNoteParams{
		NoteID: first.NoteID,
		Title:  "updated title",
		Body:   "updated body",
	})
	require.NoError(t, err)
	require.Equal(t, "updated title", updated.Title)
	require.Equal(t, "updated body", updated.Body)
	require.False(t, updated.UpdatedAt.Before(updated.CreatedAt))

	err = store.DeleteNote(ctx, first.NoteID)
	require.NoError(t, err)

	_, err = store.GetNote(ctx, first.NoteID)
	require.ErrorIs(t, err, notes.ErrNotFound)
}

func TestPostgresStoreNotFound(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	store := newTestPostgresStore(t)
	noteID := uuid.New()

	_, err := store.GetNote(ctx, noteID)
	require.ErrorIs(t, err, notes.ErrNotFound)

	_, err = store.UpdateNote(ctx, notes.UpdateNoteParams{
		NoteID: noteID,
		Title:  "missing title",
		Body:   "missing body",
	})
	require.ErrorIs(t, err, notes.ErrNotFound)
}

func TestNoteHandlersUseRedisCache(t *testing.T) {
	t.Parallel()

	ctx := t.Context()
	store := newTestPostgresStore(t)
	rdb := testutil.NewRedis(t)
	cache := notes.NewRedisCache(rdb)

	created, err := store.CreateNote(ctx, notes.CreateNoteParams{
		NoteID: uuid.New(),
		Title:  "cache title",
		Body:   "cache body",
	})
	require.NoError(t, err)
	requireCacheExists(t, created.NoteID, rdb, 0)

	getNote := notes.NewGetQueryHandler(store, cache)
	got, err := getNote.Handle(ctx, &notes.GetQuery{NoteID: created.NoteID})
	require.NoError(t, err)
	require.Equal(t, created.NoteID, got.NoteID)
	requireCacheExists(t, created.NoteID, rdb, 1)

	updateNote := notes.NewUpdateCommandHandler(store, cache)
	updated, err := updateNote.Handle(ctx, &notes.UpdateCommand{
		NoteID: created.NoteID,
		Title:  "cached update",
		Body:   "updated body",
	})
	require.NoError(t, err)
	require.Equal(t, "cached update", updated.Title)

	cached, err := cache.GetNote(ctx, created.NoteID)
	require.NoError(t, err)
	require.Equal(t, "cached update", cached.Title)

	deleteNote := notes.NewDeleteCommandHandler(store, cache)
	err = deleteNote.Handle(ctx, &notes.DeleteCommand{NoteID: created.NoteID})
	require.NoError(t, err)
	requireCacheExists(t, created.NoteID, rdb, 0)
}

func newTestPostgresStore(t *testing.T) *notes.PostgresStore {
	t.Helper()

	conn := testutil.NewPostgres(t, migrationsFS(t))

	return notes.NewPostgresStore(gen.New(conn))
}

func migrationsFS(t *testing.T) fs.FS {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)

	return os.DirFS(filepath.Join(filepath.Dir(filename), "..", ".."))
}

func requireCacheExists(t *testing.T, noteID uuid.UUID, rdb redisExistser, expected int64) {
	t.Helper()

	exists, err := rdb.Exists(t.Context(), "notes:"+noteID.String()).Result()
	require.NoError(t, err)
	require.Equal(t, expected, exists)
}

type redisExistser interface {
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
}
