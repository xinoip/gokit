package notes

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func wrapNotFound(action string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return fmt.Errorf("%s: %w", action, err)
}
