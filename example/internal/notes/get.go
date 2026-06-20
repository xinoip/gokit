package notes

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type noteReader interface {
	GetNote(ctx context.Context, noteID uuid.UUID) (*Note, error)
}

type noteCacheGetterSetter interface {
	GetNote(ctx context.Context, noteID uuid.UUID) (*Note, error)
	SetNote(ctx context.Context, note *Note) error
}

type GetQuery struct {
	NoteID uuid.UUID
}

type GetQueryHandler struct {
	store noteReader
	cache noteCacheGetterSetter
}

func NewGetQueryHandler(store noteReader, cache noteCacheGetterSetter) *GetQueryHandler {
	return &GetQueryHandler{
		store: store,
		cache: cache,
	}
}

func (h *GetQueryHandler) Handle(ctx context.Context, query *GetQuery) (*Note, error) {
	cachedNote, err := h.cache.GetNote(ctx, query.NoteID)
	if err == nil {
		return cachedNote, nil
	}
	if !errors.Is(err, errCacheMiss) {
		return nil, err
	}

	note, err := h.store.GetNote(ctx, query.NoteID)
	if err != nil {
		return nil, err
	}

	err = h.cache.SetNote(ctx, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}
