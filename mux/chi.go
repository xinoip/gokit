package mux

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/httprate"
)

// NewChiParams of [NewChi].
type NewChiParams struct {
	AllowOrigins []string
}

// NewChi creates a new mux with chi library that has sane defaults and
// middlewares setup.
func NewChi(p *NewChiParams) *chi.Mux {
	r := chi.NewMux()

	r.Use(
		httplog.RequestLogger(
			slog.Default(),
			//nolint:exhaustruct // [httplog.Options]
			&httplog.Options{
				Level:         slog.LevelWarn,
				Schema:        httplog.SchemaECS,
				RecoverPanics: true,
				Skip: func(_ *http.Request, respStatus int) bool {
					return respStatus == 404 || respStatus == 401 || respStatus == 403
				},
			},
		),
	)
	r.Use(httprate.LimitByIP(maxRequestCountPerMinute, time.Minute))
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(requestTimeout))
	r.Use(
		cors.Handler(
			//nolint:exhaustruct
			cors.Options{
				AllowedOrigins: p.AllowOrigins,
				AllowedMethods: []string{
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodDelete,
					http.MethodPatch,
					http.MethodOptions,
					http.MethodHead,
				},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: true,
				MaxAge:           corsMaxAge,
				Debug:            false,
			},
		),
	)

	return r
}
