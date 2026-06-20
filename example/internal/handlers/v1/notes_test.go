package v1_test

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"example/internal/gen"
	v1 "example/internal/handlers/v1"
	"example/internal/notes"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/xinoip/gokit/testutil"
)

func TestNotesAPIIntegration(t *testing.T) {
	t.Parallel()

	testRegistry, testAPI := testutil.NewTestAPIRegistry(t)
	conn := testutil.NewPostgres(t, migrationsFS(t))
	rdb := testutil.NewRedis(t)
	store := notes.NewPostgresStore(gen.New(conn))
	cache := notes.NewRedisCache(rdb)

	handlers := v1.Handlers{
		CreateNoteCommandHandler: notes.NewCreateCommandHandler(store, cache),
		UpdateNoteCommandHandler: notes.NewUpdateCommandHandler(store, cache),
		DeleteNoteCommandHandler: notes.NewDeleteCommandHandler(store, cache),
		ListNotesQueryHandler:    notes.NewListQueryHandler(store),
		GetNoteQueryHandler:      notes.NewGetQueryHandler(store, cache),
	}
	handlers.Register(testRegistry)

	createRes := testAPI.Post("/v1/notes", noteWriteJSON{
		Title: "api title",
		Body:  "api body",
	})
	require.Equal(t, http.StatusOK, createRes.Code)
	created := testutil.UnmarshalResponseJSON[v1.NoteBody](t, createRes)
	require.NotEqual(t, uuid.Nil, created.NoteID)
	require.Equal(t, "api title", created.Title)
	require.Equal(t, "api body", created.Body)

	getRes := testAPI.Get("/v1/notes/" + created.NoteID.String())
	require.Equal(t, http.StatusOK, getRes.Code)
	got := testutil.UnmarshalResponseJSON[v1.NoteBody](t, getRes)
	require.Equal(t, created.NoteID, got.NoteID)
	require.Equal(t, created.Title, got.Title)
	require.Equal(t, created.Body, got.Body)

	listRes := testAPI.Get("/v1/notes")
	require.Equal(t, http.StatusOK, listRes.Code)
	listed := testutil.UnmarshalResponseJSON[[]v1.NoteBody](t, listRes)
	require.Len(t, listed, 1)
	require.Equal(t, created.NoteID, listed[0].NoteID)

	updateRes := testAPI.Put("/v1/notes/"+created.NoteID.String(), noteWriteJSON{
		Title: "updated api title",
		Body:  "updated api body",
	})
	require.Equal(t, http.StatusOK, updateRes.Code)
	updated := testutil.UnmarshalResponseJSON[v1.NoteBody](t, updateRes)
	require.Equal(t, created.NoteID, updated.NoteID)
	require.Equal(t, "updated api title", updated.Title)
	require.Equal(t, "updated api body", updated.Body)

	deleteRes := testAPI.Delete("/v1/notes/" + created.NoteID.String())
	require.Equal(t, http.StatusOK, deleteRes.Code)
	deleted := testutil.UnmarshalResponseJSON[v1.DeleteNoteBody](t, deleteRes)
	require.True(t, deleted.Deleted)

	missingRes := testAPI.Get("/v1/notes/" + created.NoteID.String())
	require.Equal(t, http.StatusNotFound, missingRes.Code)
}

type noteWriteJSON struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func migrationsFS(t *testing.T) fs.FS {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)

	return os.DirFS(filepath.Join(filepath.Dir(filename), "..", "..", ".."))
}
