package v1

import (
	"example/internal/notes"

	"github.com/xinoip/gokit/api"
)

type Handlers struct {
	RPCNotes *notes.RPC
}

func NewHandlers(rpc *notes.RPC) *Handlers {
	return &Handlers{
		RPCNotes: rpc,
	}
}

func (h *Handlers) Register(r *api.Registry) {
	api.RPC(r, h.RPCNotes.List, "list_notes", api.WithInsecure())
	api.RPC(r, h.RPCNotes.Create, "create_note", api.WithInsecure())
	api.RPC(r, h.RPCNotes.Get, "get_note", api.WithInsecure())
	api.RPC(r, h.RPCNotes.Update, "update_note", api.WithInsecure())
	api.RPC(r, h.RPCNotes.Delete, "delete_note", api.WithInsecure())
}
