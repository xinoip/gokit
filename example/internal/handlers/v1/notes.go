package v1

import (
	"context"
	"errors"
	"fmt"

	"example/internal/notes"

	"github.com/danielgtaylor/huma/v2"
)

func (h *Handlers) listNotes(ctx context.Context, _ *struct{}) (*ListNotesRes, error) {
	noteList, err := h.ListNotesQueryHandler.Handle(ctx)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}

	res := &ListNotesRes{
		Body: make([]NoteBody, 0, len(noteList)),
	}
	for _, note := range noteList {
		res.Body = append(res.Body, noteBodyFromFeature(note))
	}

	return res, nil
}

func (h *Handlers) createNote(ctx context.Context, req *CreateNoteReq) (*CreateNoteRes, error) {
	note, err := h.CreateNoteCommandHandler.Handle(ctx, &notes.CreateCommand{
		Title: req.Body.Title,
		Body:  req.Body.Body,
	})
	if err != nil {
		return nil, fmt.Errorf("create note: %w", err)
	}

	return &CreateNoteRes{Body: noteBodyFromFeature(*note)}, nil
}

func (h *Handlers) getNote(ctx context.Context, req *GetNoteReq) (*GetNoteRes, error) {
	note, err := h.GetNoteQueryHandler.Handle(ctx, &notes.GetQuery{NoteID: req.NoteID})
	if err != nil {
		return nil, noteError("get note", err)
	}

	return &GetNoteRes{Body: noteBodyFromFeature(*note)}, nil
}

func (h *Handlers) updateNote(ctx context.Context, req *UpdateNoteReq) (*UpdateNoteRes, error) {
	note, err := h.UpdateNoteCommandHandler.Handle(ctx, &notes.UpdateCommand{
		Title:  req.Body.Title,
		Body:   req.Body.Body,
		NoteID: req.NoteID,
	})
	if err != nil {
		return nil, noteError("update note", err)
	}

	return &UpdateNoteRes{Body: noteBodyFromFeature(*note)}, nil
}

func (h *Handlers) deleteNote(ctx context.Context, req *DeleteNoteReq) (*DeleteNoteRes, error) {
	err := h.DeleteNoteCommandHandler.Handle(ctx, &notes.DeleteCommand{NoteID: req.NoteID})
	if err != nil {
		return nil, noteError("delete note", err)
	}

	return &DeleteNoteRes{Body: DeleteNoteBody{Deleted: true}}, nil
}

func noteError(action string, err error) error {
	if errors.Is(err, notes.ErrNotFound) {
		return huma.Error404NotFound("note not found")
	}

	return fmt.Errorf("%s: %w", action, err)
}

func noteBodyFromFeature(note notes.Note) NoteBody {
	return NoteBody{
		NoteID:    note.NoteID,
		Title:     note.Title,
		Body:      note.Body,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
