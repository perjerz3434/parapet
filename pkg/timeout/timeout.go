package timeout

import (
	"net/http"
	"time"
)

// Timeout middleware
type Timeout struct {
	Duration time.Duration
	Message  string
}

// ServeHandler implements middleware interface
func (m *Timeout) ServeHandler(h http.Handler) http.Handler {
	if m.Duration <= 0 {
		return h
	}

	return http.TimeoutHandler(h, m.Duration, m.Message)
}