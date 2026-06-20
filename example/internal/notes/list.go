package notes

import (
	"context"
)

type noteLister interface {
	ListNotes(ctx context.Context) ([]Note, error)
}

type ListQueryHandler struct {
	store noteLister
}

func NewListQueryHandler(store noteLister) *ListQueryHandler {
	return &ListQueryHandler{
		store: store,
	}
}

func (h *ListQueryHandler) Handle(ctx context.Context) ([]Note, error) {
	return h.store.ListNotes(ctx)
}
