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

// ChiConfig configures [NewChi].
type ChiConfig struct {
	AllowedOrigins      []string
	Logger              *slog.Logger
	RequestTimeout      time.Duration
	RateLimit           int
	RateWindow          time.Duration
	RateLimitMiddleware func(http.Handler) http.Handler
	CORSMaxAge          int
	AllowCredentials    bool
}

// DefaultChiConfig returns opinionated, production oriented [ChiConfig].
func DefaultChiConfig() ChiConfig {
	return ChiConfig{
		AllowedOrigins:      []string{"https://*"},
		Logger:              slog.Default(),
		RequestTimeout:      requestTimeout,
		RateLimit:           maxRequestCountPerMinute,
		RateWindow:          time.Minute,
		RateLimitMiddleware: nil,
		CORSMaxAge:          corsMaxAge,
		AllowCredentials:    true,
	}
}

// NewChi creates a new mux with chi library using passed in config.
func NewChi(p ChiConfig) *chi.Mux {
	r := chi.NewMux()

	r.Use(
		httplog.RequestLogger(
			p.Logger,
			//nolint:exhaustruct // [httplog.Options]
			&httplog.Options{
				Level:         slog.LevelWarn,
				Schema:        httplog.SchemaECS,
				RecoverPanics: true,
				Skip: func(_ *http.Request, respStatus int) bool {
					return respStatus == http.StatusNotFound
				},
			},
		),
	)
	if p.RateLimitMiddleware != nil {
		r.Use(p.RateLimitMiddleware)
	} else if p.RateLimit > 0 && p.RateWindow > 0 {
		r.Use(httprate.LimitByIP(p.RateLimit, p.RateWindow))
	}
	if p.RequestTimeout > 0 {
		r.Use(chimiddleware.Timeout(p.RequestTimeout))
	}
	r.Use(
		cors.Handler(
			//nolint:exhaustruct
			cors.Options{
				AllowedOrigins: p.AllowedOrigins,
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
				AllowCredentials: p.AllowCredentials,
				MaxAge:           p.CORSMaxAge,
				Debug:            false,
			},
		),
	)

	return r
}
