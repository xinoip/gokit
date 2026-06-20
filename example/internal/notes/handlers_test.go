package notes_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"example/internal/notes"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateCommandHandlerReturnsCacheError(t *testing.T) {
	t.Parallel()

	errCache := errors.New("cache failed")
	handler := notes.NewCreateCommandHandler(
		noteCreatorFunc(func(_ context.Context, params notes.CreateNoteParams) (*notes.Note, error) {
			return testNote(params.NoteID, params.Title, params.Body), nil
		}),
		noteCacheSetterFunc(func(_ context.Context, _ *notes.Note) error {
			return errCache
		}),
	)

	_, err := handler.Handle(t.Context(), &notes.CreateCommand{
		Title: "title",
		Body:  "body",
	})
	require.ErrorIs(t, err, errCache)
}

func TestCreateCommandHandlerReturnsStoreError(t *testing.T) {
	t.Parallel()

	errStore := errors.New("store failed")
	cacheCalled := false
	handler := notes.NewCreateCommandHandler(
		noteCreatorFunc(func(_ context.Context, _ notes.CreateNoteParams) (*notes.Note, error) {
			return nil, errStore
		}),
		noteCacheSetterFunc(func(_ context.Context, _ *notes.Note) error {
			cacheCalled = true
			return nil
		}),
	)

	_, err := handler.Handle(t.Context(), &notes.CreateCommand{
		Title: "title",
		Body:  "body",
	})
	require.ErrorIs(t, err, errStore)
	require.False(t, cacheCalled)
}

func TestGetQueryHandlerReturnsCachedNote(t *testing.T) {
	t.Parallel()

	note := testNote(uuid.New(), "cached title", "cached body")
	storeCalled := false
	cacheSetCalled := false
	handler := notes.NewGetQueryHandler(
		noteReaderFunc(func(_ context.Context, _ uuid.UUID) (*notes.Note, error) {
			storeCalled = true
			return testNote(uuid.New(), "store title", "store body"), nil
		}),
		noteCacheGetterSetterFake{
			getNote: func(_ context.Context, _ uuid.UUID) (*notes.Note, error) {
				return note, nil
			},
			setNote: func(_ context.Context, _ *notes.Note) error {
				cacheSetCalled = true
				return nil
			},
		},
	)

	got, err := handler.Handle(t.Context(), &notes.GetQuery{NoteID: note.NoteID})
	require.NoError(t, err)
	require.Equal(t, note, got)
	require.False(t, storeCalled)
	require.False(t, cacheSetCalled)
}

func TestGetQueryHandlerFallsBackToStoreOnCacheMiss(t *testing.T) {
	t.Parallel()

	note := testNote(uuid.New(), "store title", "store body")
	handler := notes.NewGetQueryHandler(
		noteReaderFunc(func(_ context.Context, noteID uuid.UUID) (*notes.Note, error) {
			require.Equal(t, note.NoteID, noteID)
			return note, nil
		}),
		notes.NewRedisCache(nil),
	)

	got, err := handler.Handle(t.Context(), &notes.GetQuery{NoteID: note.NoteID})
	require.NoError(t, err)
	require.Equal(t, note, got)
}

func TestGetQueryHandlerReturnsCacheError(t *testing.T) {
	t.Parallel()

	errCache := errors.New("cache failed")
	storeCalled := false
	handler := notes.NewGetQueryHandler(
		noteReaderFunc(func(_ context.Context, _ uuid.UUID) (*notes.Note, error) {
			storeCalled = true
			return testNote(uuid.New(), "store title", "store body"), nil
		}),
		noteCacheGetterSetterFake{
			getNote: func(_ context.Context, _ uuid.UUID) (*notes.Note, error) {
				return nil, errCache
			},
			setNote: func(_ context.Context, _ *notes.Note) error {
				return nil
			},
		},
	)

	_, err := handler.Handle(t.Context(), &notes.GetQuery{NoteID: uuid.New()})
	require.ErrorIs(t, err, errCache)
	require.False(t, storeCalled)
}

func TestDeleteCommandHandlerReturnsStoreError(t *testing.T) {
	t.Parallel()

	errStore := errors.New("store failed")
	storeDeleteCalled := false
	cacheDeleteCalled := false
	handler := notes.NewDeleteCommandHandler(
		noteDeleteStoreFake{
			getNote: func(_ context.Context, _ uuid.UUID) (*notes.Note, error) {
				return nil, errStore
			},
			deleteNote: func(_ context.Context, _ uuid.UUID) error {
				storeDeleteCalled = true
				return nil
			},
		},
		noteCacheDeleterFunc(func(_ context.Context, _ uuid.UUID) error {
			cacheDeleteCalled = true
			return nil
		}),
	)

	err := handler.Handle(t.Context(), &notes.DeleteCommand{NoteID: uuid.New()})
	require.ErrorIs(t, err, errStore)
	require.False(t, storeDeleteCalled)
	require.False(t, cacheDeleteCalled)
}

func TestDeleteCommandHandlerReturnsCacheError(t *testing.T) {
	t.Parallel()

	note := testNote(uuid.New(), "title", "body")
	errCache := errors.New("cache failed")
	handler := notes.NewDeleteCommandHandler(
		noteDeleteStoreFake{
			getNote: func(_ context.Context, _ uuid.UUID) (*notes.Note, error) {
				return note, nil
			},
			deleteNote: func(_ context.Context, _ uuid.UUID) error {
				return nil
			},
		},
		noteCacheDeleterFunc(func(_ context.Context, _ uuid.UUID) error {
			return errCache
		}),
	)

	err := handler.Handle(t.Context(), &notes.DeleteCommand{NoteID: note.NoteID})
	require.ErrorIs(t, err, errCache)
}

func testNote(noteID uuid.UUID, title string, body string) *notes.Note {
	now := time.Now()

	return &notes.Note{
		NoteID:    noteID,
		Title:     title,
		Body:      body,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type noteCreatorFunc func(ctx context.Context, params notes.CreateNoteParams) (*notes.Note, error)

func (f noteCreatorFunc) CreateNote(ctx context.Context, params notes.CreateNoteParams) (*notes.Note, error) {
	return f(ctx, params)
}

type noteCacheSetterFunc func(ctx context.Context, note *notes.Note) error

func (f noteCacheSetterFunc) SetNote(ctx context.Context, note *notes.Note) error {
	return f(ctx, note)
}

type noteReaderFunc func(ctx context.Context, noteID uuid.UUID) (*notes.Note, error)

func (f noteReaderFunc) GetNote(ctx context.Context, noteID uuid.UUID) (*notes.Note, error) {
	return f(ctx, noteID)
}

type noteCacheGetterSetterFake struct {
	getNote func(ctx context.Context, noteID uuid.UUID) (*notes.Note, error)
	setNote func(ctx context.Context, note *notes.Note) error
}

func (f noteCacheGetterSetterFake) GetNote(ctx context.Context, noteID uuid.UUID) (*notes.Note, error) {
	return f.getNote(ctx, noteID)
}

func (f noteCacheGetterSetterFake) SetNote(ctx context.Context, note *notes.Note) error {
	return f.setNote(ctx, note)
}

type noteDeleteStoreFake struct {
	getNote    func(ctx context.Context, noteID uuid.UUID) (*notes.Note, error)
	deleteNote func(ctx context.Context, noteID uuid.UUID) error
}

func (f noteDeleteStoreFake) GetNote(ctx context.Context, noteID uuid.UUID) (*notes.Note, error) {
	return f.getNote(ctx, noteID)
}

func (f noteDeleteStoreFake) DeleteNote(ctx context.Context, noteID uuid.UUID) error {
	return f.deleteNote(ctx, noteID)
}

type noteCacheDeleterFunc func(ctx context.Context, noteID uuid.UUID) error

func (f noteCacheDeleterFunc) DeleteNote(ctx context.Context, noteID uuid.UUID) error {
	return f(ctx, noteID)
}
