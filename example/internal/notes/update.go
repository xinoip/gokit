package notes

import (
	"context"

	"github.com/google/uuid"
)

type noteUpdater interface {
	UpdateNote(ctx context.Context, params UpdateNoteParams) (*Note, error)
}

type UpdateCommand struct {
	NoteID uuid.UUID
	Title  string
	Body   string
}

type UpdateCommandHandler struct {
	store noteUpdater
	cache noteCacheSetter
}

func NewUpdateCommandHandler(store noteUpdater, cache noteCacheSetter) *UpdateCommandHandler {
	return &UpdateCommandHandler{
		store: store,
		cache: cache,
	}
}

func (h *UpdateCommandHandler) Handle(ctx context.Context, cmd *UpdateCommand) (*Note, error) {
	note, err := h.store.UpdateNote(ctx, UpdateNoteParams{
		Title:  cmd.Title,
		Body:   cmd.Body,
		NoteID: cmd.NoteID,
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
