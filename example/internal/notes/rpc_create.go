package notes

import (
	"context"

	"github.com/google/uuid"
)

type RPCCreateParams struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type RPCCreateResult struct {
	Note *Note `json:"note"`
}

func (r *RPC) Create(ctx context.Context, p *RPCCreateParams) (*RPCCreateResult, error) {
	note, err := r.store.CreateNote(ctx, CreateNoteParams{
		NoteID: uuid.New(),
		Title:  p.Title,
		Body:   p.Body,
	})
	if err != nil {
		return nil, err
	}

	err = r.cache.SetNote(ctx, note)
	if err != nil {
		return nil, err
	}

	return &RPCCreateResult{
		Note: note,
	}, nil
}
