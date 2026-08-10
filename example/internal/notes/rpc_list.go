package notes

import (
	"context"
)

type RPCListResult struct {
	Notes []Note `json:"notes"`
}

func (r *RPC) List(ctx context.Context, p *struct{}) (*RPCListResult, error) {
	notes, err := r.store.ListNotes(ctx)
	if err != nil {
		return nil, err
	}

	return &RPCListResult{
		Notes: notes,
	}, nil
}
