package notes

import (
	"context"

	"github.com/google/uuid"
)

type noteCreator interface {
	CreateNote(ctx context.Context, params CreateNoteParams) (*Note, error)
}

type noteCacheSetter interface {
	SetNote(ctx context.Context, note *Note) error
}

type CreateCommand struct {
	Title string
	Body  string
}

type CreateCommandHandler struct {
	store noteCreator
	cache noteCacheSetter
}

func NewCreateCommandHandler(store noteCreator, cache noteCacheSetter) *CreateCommandHandler {
	return &CreateCommandHandler{
		store: store,
		cache: cache,
	}
}

func (h *CreateCommandHandler) Handle(ctx context.Context, cmd *CreateCommand) (*Note, error) {
	note, err := h.store.CreateNote(ctx, CreateNoteParams{
		NoteID: uuid.New(),
		Title:  cmd.Title,
		Body:   cmd.Body,
	})
	if err != nil {
		return nil, err
	}

	err = h.cache.SetNote(ctx, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}
