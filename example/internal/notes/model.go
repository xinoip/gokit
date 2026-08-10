package notes

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	NoteID    uuid.UUID `json:"noteId"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
