package notes

import (
	"context"

	"github.com/google/uuid"
)

type RPCUpdateParams struct {
	NoteID uuid.UUID `json:"noteID"`
	Title  string    `json:"title"`
	Body   string    `json:"body"`
}

type RPCUpdateResult struct {
	Note *Note `json:"note"`
}

func (r *RPC) Update(ctx context.Context, p *RPCUpdateParams) (*RPCUpdateResult, error) {
	note, err := r.store.UpdateNote(ctx, UpdateNoteParams{
		Title:  p.Title,
		Body:   p.Body,
		NoteID: p.NoteID,
	})
	if err != nil {
		return nil, err
	}

	err = r.cache.SetNote(ctx, note)
	if err != nil {
		return nil, err
	}

	return &RPCUpdateResult{
		Note: note,
	}, nil
}
