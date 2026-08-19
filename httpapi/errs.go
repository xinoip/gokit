package httpapi

import (
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
)

// HTTP400 is a convenience function for creating a 400 Bad Request error
// targeting [huma] library. Error details exposed to public.
func HTTP400(err ...error) error {
	return huma.Error400BadRequest("bad request", err...)
}

// HTTP401 is a convenience function for creating a 401 Unauthorized error
// targeting [huma] library. Error details are logged and not exposed to
// public.
func HTTP401(err error) error {
	slog.Debug("unauthorized", "error", err)
	return huma.Error401Unauthorized("unauthorized")
}

// HTTP403 is same as [HTTP401] but for 403 Forbidden.
func HTTP403(err error) error {
	slog.Debug("forbidden", "error", err)
	return huma.Error403Forbidden("forbidden")
}

// HTTP404 is same as [HTTP400] but for 404 Not Found.
func HTTP404(err error) error {
	return huma.Error404NotFound("not found", err)
}

// HTTP500 is same as [HTTP401] but for 500 Internal Server Error.
func HTTP500(err error) error {
	slog.Error("internal server error", "error", err)
	return huma.Error500InternalServerError("internal server error")
}
