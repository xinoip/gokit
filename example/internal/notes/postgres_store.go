package notes

import (
	"context"
	"errors"
	"fmt"

	"example/internal/gen"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/xinoip/gokit/errs"
)

type CreateNoteParams struct {
	NoteID uuid.UUID
	Title  string
	Body   string
}

type UpdateNoteParams struct {
	NoteID uuid.UUID
	Title  string
	Body   string
}

type PostgresStore struct {
	db *gen.Queries
}

func NewPostgresStore(db *gen.Queries) *PostgresStore {
	return &PostgresStore{
		db: db,
	}
}

func (s *PostgresStore) CreateNote(ctx context.Context, params CreateNoteParams) (*Note, error) {
	dbNote, err := s.db.NotesCreate(ctx, gen.NotesCreateParams{
		NoteID: params.NoteID,
		Title:  params.Title,
		Body:   params.Body,
	})
	if err != nil {
		return nil, fmt.Errorf("create note: %w", err)
	}

	return noteFromDB(dbNote), nil
}

func (s *PostgresStore) UpdateNote(ctx context.Context, params UpdateNoteParams) (*Note, error) {
	dbNote, err := s.db.NotesUpdate(ctx, gen.NotesUpdateParams{
		Title:  params.Title,
		Body:   params.Body,
		NoteID: params.NoteID,
	})
	if err != nil {
		return nil, wrapNotFound("update note", err)
	}

	return noteFromDB(dbNote), nil
}

func (s *PostgresStore) DeleteNote(ctx context.Context, noteID uuid.UUID) error {
	err := s.db.NotesDelete(ctx, noteID)
	if err != nil {
		return fmt.Errorf("delete note: %w", err)
	}

	return nil
}

func (s *PostgresStore) GetNote(ctx context.Context, noteID uuid.UUID) (*Note, error) {
	dbNote, err := s.db.NotesGet(ctx, noteID)
	if err != nil {
		return nil, wrapNotFound("get note", err)
	}

	return noteFromDB(dbNote), nil
}

func (s *PostgresStore) ListNotes(ctx context.Context) ([]Note, error) {
	dbNotes, err := s.db.NotesList(ctx)
	if err != nil {
		return nil, fmt.Errorf("list notes: %w", err)
	}

	notes := make([]Note, 0, len(dbNotes))
	for _, dbNote := range dbNotes {
		notes = append(notes, *noteFromDB(dbNote))
	}

	return notes, nil
}

func noteFromDB(note gen.Notes) *Note {
	return &Note{
		NoteID:    note.NoteID,
		Title:     note.Title,
		Body:      note.Body,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func wrapNotFound(action string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return errs.NotFoundf("%s: %w", action, err)
	}

	return fmt.Errorf("%s: %w", action, err)
}
