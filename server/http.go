package server

import (
	"context"
	"log/slog"
	"net/http"
)

// ServeHTTPParams of [ServeHTTP].
type ServeHTTPParams struct {
	Addr    string
	Handler http.Handler
}

// ServeHTTP serves an HTTP server with sane defaults and graceful shutdown support.
func ServeHTTP(ctx context.Context, p *ServeHTTPParams) error {
	s := &http.Server{
		Addr:              p.Addr,
		Handler:           p.Handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	errChan := make(chan error, 1)
	go func() {
		slog.Info("Started HTTP server", "addr", p.Addr)
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		slog.Error("Error encountered while serving HTTP server", "addr", p.Addr, "err", err.Error())
		return err
	case <-ctx.Done():
		slog.Info("Shutting down HTTP server", "addr", p.Addr, "timeout", shutdownTimeout)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		//nolint:contextcheck // ctx is already done, so use background context for shutdown.
		err := s.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}
		<-errChan
		return nil
	}
}
