// Package testutil provides common utilities to support testing.
package testutil

import (
	"context"
	"testing"
	"time"
)

// TestCtxWithTimeout creates a new context with a timeout, suitable to use in
// tests to set quick timeouts.
func TestCtxWithTimeout(t *testing.T, timeout time.Duration) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)

	return ctx
}

// DefaultTestCtx is [TestCtxWithTimeout] with a default timeout of 10 seconds.
func DefaultTestCtx(t *testing.T) context.Context {
	t.Helper()

	return TestCtxWithTimeout(t, 10*time.Second)
}
