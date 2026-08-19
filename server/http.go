package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// HTTPConfig configures [ServeHTTP]. Start with [DefaultHTTPConfig] and
// override the values needed by the application.
type HTTPConfig struct {
	Addr              string
	Handler           http.Handler
	Logger            *slog.Logger
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	MaxHeaderBytes    int
}

// DefaultHTTPConfig returns the production-oriented server defaults.
func DefaultHTTPConfig(addr string, handler http.Handler) HTTPConfig {
	return HTTPConfig{
		Addr:              addr,
		Handler:           handler,
		Logger:            slog.Default(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		ShutdownTimeout:   shutdownTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
}

// ServeHTTP serves an HTTP server with passed in [config] and graceful
// shut down support.
func ServeHTTP(ctx context.Context, config HTTPConfig) error {
	s := &http.Server{
		Addr:              config.Addr,
		Handler:           config.Handler,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		MaxHeaderBytes:    config.MaxHeaderBytes,
	}

	errChan := make(chan error, 1)

	go func() {
		config.Logger.Info("Starting HTTP server", "addr", config.Addr)

		err := s.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		errChan <- err
	}()

	select {
	case err := <-errChan:
		if err != nil {
			config.Logger.Error("Error encountered while serving HTTP server", "addr", config.Addr, "err", err)
			return fmt.Errorf("serve HTTP: %w", err)
		}

		return nil
	case <-ctx.Done():
		config.Logger.Info("Shutting down HTTP server", "addr", config.Addr, "timeout", config.ShutdownTimeout)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		defer cancel()

		//nolint:contextcheck // shutdownCtx uses new background context since ctx is already done.
		shutdownErr := s.Shutdown(shutdownCtx)
		if shutdownErr != nil {
			shutdownErr = fmt.Errorf("shut down HTTP server: %w", shutdownErr)

			closeErr := s.Close()
			if closeErr != nil {
				closeErr = fmt.Errorf("force close HTTP server: %w", closeErr)
			}

			return errors.Join(shutdownErr, closeErr)
		}

		return nil
	}
}
