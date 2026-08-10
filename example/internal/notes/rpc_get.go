package notes

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type RPCGetParams struct {
	NoteID uuid.UUID `json:"noteId"`
}

type RPCGetResult struct {
	Note *Note `json:"note"`
}

func (r *RPC) Get(ctx context.Context, p *RPCGetParams) (*RPCGetResult, error) {
	cachedNote, err := r.cache.GetNote(ctx, p.NoteID)
	if err == nil {
		return &RPCGetResult{
			Note: cachedNote,
		}, nil
	}
	if !errors.Is(err, errCacheMiss) {
		return nil, err
	}

	note, err := r.store.GetNote(ctx, p.NoteID)
	if err != nil {
		return nil, err
	}

	err = r.cache.SetNote(ctx, note)
	if err != nil {
		return nil, err
	}

	return &RPCGetResult{
		Note: note,
	}, nil
}
