package v1

import (
	"example/internal/notes"

	"github.com/xinoip/gokit/api"
)

type Handlers struct {
	CreateNoteCommandHandler *notes.CreateCommandHandler
	UpdateNoteCommandHandler *notes.UpdateCommandHandler
	DeleteNoteCommandHandler *notes.DeleteCommandHandler
	ListNotesQueryHandler    *notes.ListQueryHandler
	GetNoteQueryHandler      *notes.GetQueryHandler
}

func (h *Handlers) Register(r *api.Registry) {
	api.Get(r, "/v1/notes", h.listNotes, "list-notes", api.WithInsecure())
	api.Post(r, "/v1/notes", h.createNote, "create-note", api.WithInsecure())
	api.Get(r, "/v1/notes/{note-id}", h.getNote, "get-note", api.WithInsecure())
	api.Put(r, "/v1/notes/{note-id}", h.updateNote, "update-note", api.WithInsecure())
	api.Delete(r, "/v1/notes/{note-id}", h.deleteNote, "delete-note", api.WithInsecure())
}
