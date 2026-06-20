package notes

import (
	"context"

	"github.com/google/uuid"
)

type noteDeleteStore interface {
	GetNote(ctx context.Context, noteID uuid.UUID) (*Note, error)
	DeleteNote(ctx context.Context, noteID uuid.UUID) error
}

type noteCacheDeleter interface {
	DeleteNote(ctx context.Context, noteID uuid.UUID) error
}

type DeleteCommand struct {
	NoteID uuid.UUID
}

type DeleteCommandHandler struct {
	store noteDeleteStore
	cache noteCacheDeleter
}

func NewDeleteCommandHandler(store noteDeleteStore, cache noteCacheDeleter) *DeleteCommandHandler {
	return &DeleteCommandHandler{
		store: store,
		cache: cache,
	}
}

func (h *DeleteCommandHandler) Handle(ctx context.Context, cmd *DeleteCommand) error {
	_, err := h.store.GetNote(ctx, cmd.NoteID)
	if err != nil {
		return err
	}

	err = h.store.DeleteNote(ctx, cmd.NoteID)
	if err != nil {
		return err
	}

	return h.cache.DeleteNote(ctx, cmd.NoteID)
}
