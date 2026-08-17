package v1

import (
	"example/internal/notes"

	"github.com/xinoip/gokit/httpapi"
	"github.com/xinoip/gokit/httprpc"
)

type Handlers struct {
	RPCNotes *notes.RPC
}

func NewHandlers(rpc *notes.RPC) *Handlers {
	return &Handlers{
		RPCNotes: rpc,
	}
}

func (h *Handlers) Register(r *httpapi.Registry) {
	httprpc.Handle(r, h.RPCNotes.List, "list_notes", httpapi.WithInsecure())
	httprpc.Handle(r, h.RPCNotes.Create, "create_note", httpapi.WithInsecure())
	httprpc.Handle(r, h.RPCNotes.Get, "get_note", httpapi.WithInsecure())
	httprpc.Handle(r, h.RPCNotes.Update, "update_note", httpapi.WithInsecure())
	httprpc.Handle(r, h.RPCNotes.Delete, "delete_note", httpapi.WithInsecure())
}
