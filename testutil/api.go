package testutil

import (
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/require"

	"github.com/xinoip/gokit/api"
)

// NewTestAPIRegistry creates a testing instance of [api.Registry] with support
// for [humatest].
//
//nolint:ireturn // [humatest.TestAPI] is intended to be an interface.
func NewTestAPIRegistry(t *testing.T) (*api.Registry, humatest.TestAPI) {
	t.Helper()

	_, humaAPI := humatest.New(t)
	return &api.Registry{
		HumaAPI: humaAPI,
		SecureMiddleware: func(ctx huma.Context, next func(huma.Context)) {
			slog.Warn("SecureMiddleware is not implemented for testing", "case", t.Name())

			// Tests can register secure endpoints without implementing authentication.
			next(ctx)
		},
	}, humaAPI
}

// MarshalJSON is a testing helper to unmarshal a JSON string into a
// map[string]any.
func MarshalJSON(t *testing.T, s string) map[string]any {
	t.Helper()

	var m map[string]any
	err := json.Unmarshal([]byte(s), &m)
	require.NoError(t, err)

	return m
}

// UnmarshalJSON unmarshals JSON data into the requested type.
//
//nolint:ireturn // Generic test helpers intentionally return the requested type.
func UnmarshalJSON[T any](t *testing.T, data []byte) T {
	t.Helper()

	var value T
	err := json.Unmarshal(data, &value)
	require.NoError(t, err)

	return value
}

// UnmarshalResponseJSON unmarshals a test HTTP response body into the requested type.
//
//nolint:ireturn // Generic test helpers intentionally return the requested type.
func UnmarshalResponseJSON[T any](t *testing.T, res *httptest.ResponseRecorder) T {
	t.Helper()

	return UnmarshalJSON[T](t, res.Body.Bytes())
}
