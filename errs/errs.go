// Package errs provides a unified base for error definitions and utilities.
//
//nolint:revive,errname,gochecknoglobals,staticcheck // This is full of custom conventions.
package errs

import (
	"errors"
	"fmt"
)

var (
	// Failure is used when things break because of outside factors. These
	// errors are safe to be exposed to public.
	Failure = errors.New("failure")

	// NotFound is used when a thing is failed to be found. These errors are
	// safe to be exposed to public.
	NotFound = errors.New("not found")

	// Internal is used when things break because of us. All unclassified
	// errors should be treated as internal. These errors are not safe to be
	// exposed to public and should be treated as confidential bugs.
	Internal = errors.New("internal error")
)

// Failuref wraps a formatted error with [Failure].
func Failuref(format string, a ...any) error {
	return fmt.Errorf("%w: %w", Failure, fmt.Errorf(format, a...))
}

// Failurew wraps an error with [Failure].
func Failurew(err error) error {
	return fmt.Errorf("%w: %w", Failure, err)
}

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

// IsFailure checks if the given error is [Failure].
func IsFailure(err error) bool {
	return errors.Is(err, Failure)
}

// IsNotFound checks if the given error is [NotFound].
func IsNotFound(err error) bool {
	return errors.Is(err, NotFound)
}

// IsInternal checks if the given error is [Internal].
func IsInternal(err error) bool {
	return errors.Is(err, Internal)
}
