-- name: NotesCreate :one
INSERT INTO notes (
    note_id,
    title,
    body
) VALUES (
    @note_id,
    @title,
    @body
) RETURNING *;

-- name: NotesList :many
SELECT * FROM notes ORDER BY created_at DESC;

-- name: NotesGet :one
SELECT * FROM notes WHERE note_id = @note_id;

-- name: NotesUpdate :one
UPDATE notes
SET
    title = @title,
    body = @body,
    updated_at = CURRENT_TIMESTAMP
WHERE note_id = @note_id
RETURNING *;

-- name: NotesDelete :exec
DELETE FROM notes WHERE note_id = @note_id;
