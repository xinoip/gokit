package errs

import (
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
)

// HTTP400 is a convenience function for creating a 400 Bad Request error
// targeting [huma] library.
func HTTP400(err ...error) error {
	return huma.Error400BadRequest("bad request", err...)
}

// HTTP401 is same as [HTTP400] but for 401 Unauthorized.
func HTTP401(err error) error {
	slog.Debug("unauthorized", "error", err.Error())
	return huma.Error401Unauthorized("unauthorized")
}

// HTTP403 is same as [HTTP400] but for 403 Forbidden.
func HTTP403(err error) error {
	slog.Debug("forbidden", "error", err.Error())
	return huma.Error403Forbidden("forbidden")
}

// HTTP404 is same as [HTTP400] but for 404 Not Found.
func HTTP404(err error) error {
	return huma.Error404NotFound("not found", err)
}

// HTTP500 is same as [HTTP400] but for 500 Internal Server Error.
func HTTP500(err error) error {
	slog.Error("internal server error", "error", err.Error())
	return huma.Error500InternalServerError("internal server error")
}

// HTTPError is a convenience function that maps known errors to HTTP errors.
func HTTPError(err error) error {
	switch {
	case err == nil:
		return nil
	case IsFailure(err):
		return HTTP400(err)
	case IsNotFound(err):
		return HTTP404(err)
	case IsInternal(err):
		fallthrough
	default:
		return HTTP500(err)
	}
}
