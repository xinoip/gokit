package notes

import (
	"context"
)

type RPCListResult struct {
	Notes []Note `json:"notes"`
}

func (r *RPC) List(ctx context.Context, _ *struct{}) (*RPCListResult, error) {
	notes, err := r.store.ListNotes(ctx)
	if err != nil {
		return nil, err
	}

	return &RPCListResult{
		Notes: notes,
	}, nil
}
