package notes_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"example/internal/notes"
)

var errTestCache = errors.New("cache unavailable")

type stubStore struct {
	note *notes.Note
}

func (s *stubStore) CreateNote(context.Context, notes.CreateNoteParams) (*notes.Note, error) {
	return s.note, nil
}

func (s *stubStore) UpdateNote(context.Context, notes.UpdateNoteParams) (*notes.Note, error) {
	return s.note, nil
}

func (*stubStore) DeleteNote(context.Context, uuid.UUID) error {
	return nil
}

func (s *stubStore) GetNote(context.Context, uuid.UUID) (*notes.Note, error) {
	return s.note, nil
}

func (s *stubStore) ListNotes(context.Context) ([]notes.Note, error) {
	return []notes.Note{*s.note}, nil
}

type failingCache struct{}

func (*failingCache) GetNote(context.Context, uuid.UUID) (*notes.Note, error) {
	return nil, errTestCache
}

func (*failingCache) SetNote(context.Context, *notes.Note) error {
	return errTestCache
}

func (*failingCache) DeleteNote(context.Context, uuid.UUID) error {
	return errTestCache
}

func TestMutationsSucceedWhenCacheIsUnavailable(t *testing.T) {
	t.Parallel()

	noteID := uuid.New()
	note := &notes.Note{
		NoteID:    noteID,
		Title:     "title",
		Body:      "body",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	rpc, err := notes.NewRPCWithDependencies(&stubStore{note: note}, &failingCache{})
	require.NoError(t, err)

	created, err := rpc.Create(t.Context(), &notes.RPCCreateParams{Title: note.Title, Body: note.Body})
	require.NoError(t, err)
	assert.Equal(t, note, created.Note)

	updated, err := rpc.Update(t.Context(), &notes.RPCUpdateParams{
		NoteID: noteID,
		Title:  note.Title,
		Body:   note.Body,
	})
	require.NoError(t, err)
	assert.Equal(t, note, updated.Note)

	deleted, err := rpc.Delete(t.Context(), &notes.RPCDeleteParams{NoteID: noteID})
	require.NoError(t, err)
	assert.NotNil(t, deleted)
}

func TestNewRPCWithDependenciesValidatesDependencies(t *testing.T) {
	t.Parallel()

	rpc, err := notes.NewRPCWithDependencies(nil, &failingCache{})
	require.Error(t, err)
	assert.Nil(t, rpc)

	rpc, err = notes.NewRPCWithDependencies(&stubStore{note: nil}, nil)
	require.Error(t, err)
	assert.Nil(t, rpc)
}
