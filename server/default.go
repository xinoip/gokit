// Package server provides default implementations for HTTP servers.
package server

import "time"

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 120 * time.Second
	maxHeaderBytes    = 1 * 1024 * 1024
	shutdownTimeout   = 10 * time.Second
)
