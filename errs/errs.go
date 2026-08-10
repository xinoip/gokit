// Package errs provides a unified base for error definitions and utilities.
//
//nolint:revive,errname,gochecknoglobals,staticcheck // This is full of custom conventions.
package errs

import (
	"errors"
	"fmt"
)

var (
	// NotFound is used when a thing is failed to be found.
	NotFound = errors.New("not found")

	// Internal is used when things break because of us.
	Internal = errors.New("internal error")
)

// NotFoundf wraps a formatted error with [NotFound].
func NotFoundf(format string, a ...any) error {
	return fmt.Errorf("%w: %w", NotFound, fmt.Errorf(format, a...))
}

// NotFoundw wraps an error with [NotFound].
func NotFoundw(err error) error {
	return fmt.Errorf("%w: %w", NotFound, err)
}

// Internalf wraps a formatted error with [Internal].
func Internalf(format string, a ...any) error {
	return fmt.Errorf("%w: %w", Internal, fmt.Errorf(format, a...))
}

// Internalw wraps an error with [Internal].
func Internalw(err error) error {
	return fmt.Errorf("%w: %w", Internal, err)
}

// IsNotFound checks if the given error is [NotFound].
func IsNotFound(err error) bool {
	return errors.Is(err, NotFound)
}

// IsInternal checks if the given error is [Internal].
func IsInternal(err error) bool {
	return errors.Is(err, Internal)
}
