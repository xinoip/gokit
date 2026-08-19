package notes

import (
	"context"

	"github.com/google/uuid"
)

type RPCDeleteParams struct {
	NoteID uuid.UUID `json:"noteId"`
}

func (r *RPC) Delete(ctx context.Context, p *RPCDeleteParams) (*struct{}, error) {
	_, err := r.store.GetNote(ctx, p.NoteID)
	if err != nil {
		return nil, err
	}

	err = r.store.DeleteNote(ctx, p.NoteID)
	if err != nil {
		return nil, err
	}

	logCacheError("delete note", r.cache.DeleteNote(ctx, p.NoteID))

	return &struct{}{}, nil
}
