package v1

import (
	"time"

	"github.com/google/uuid"
)

type NoteBody struct {
	NoteID    uuid.UUID `json:"noteId"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NoteWriteBody struct {
	Title string `json:"title" maxLength:"200"  minLength:"1"`
	Body  string `json:"body"  maxLength:"5000" minLength:"1"`
}

type ListNotesRes struct {
	Body []NoteBody
}

type CreateNoteReq struct {
	Body NoteWriteBody
}

type CreateNoteRes struct {
	Body NoteBody
}

type GetNoteReq struct {
	NoteID uuid.UUID `doc:"Note ID" path:"note-id"`
}

type GetNoteRes struct {
	Body NoteBody
}

type UpdateNoteReq struct {
	NoteID uuid.UUID `doc:"Note ID" path:"note-id"`
	Body   NoteWriteBody
}

type UpdateNoteRes struct {
	Body NoteBody
}

type DeleteNoteReq struct {
	NoteID uuid.UUID `doc:"Note ID" path:"note-id"`
}

type DeleteNoteBody struct {
	Deleted bool `json:"deleted"`
}

type DeleteNoteRes struct {
	Body DeleteNoteBody
}
