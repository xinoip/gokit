package notes

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("note not found")

type Note struct {
	NoteID    uuid.UUID `json:"noteId"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
