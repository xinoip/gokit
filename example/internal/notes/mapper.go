package notes

import "example/internal/gen"

func noteFromDB(note gen.Notes) *Note {
	return &Note{
		NoteID:    note.NoteID,
		Title:     note.Title,
		Body:      note.Body,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
