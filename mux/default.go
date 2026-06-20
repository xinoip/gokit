// Package mux provides ready to use mux implementations with library backends
// like chi and similar.
package mux

import "time"

const (
	requestTimeout           = 30 * time.Second
	corsMaxAge               = 600
	maxRequestCountPerMinute = 250
)
