package logutil

import (
	"context"
	"net"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/groob/finalizer"
)

// NewHTTPLogger returns a Logger from a Go-Kit Logger.
func NewHTTPLogger(logger log.Logger) *HTTPLogger {
	return &HTTPLogger{logger: logger}
}

// HTTPLogger wraps the Go-Kit Logger to return a logger which implements a
// ServerFinalizerFunc.
// The ServerFinalizerFunc can be passed to finalizer.Middleware or a Go-Kit
// Server to create a structured HTTP Logger.
type HTTPLogger struct {
	logger log.Logger
}

// LoggingFinalizer is a finalizer.ServerFinalizerFunc which logs information about a completed
// HTTP Request.
func (l *HTTPLogger) LoggingFinalizer(ctx context.Context, code int, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	keyvals := []interface{}{
		"method", r.Method,
		"status", code,
		"proto", r.Proto,
		"host", host,
		"user_agent", r.UserAgent(),
	}

	if referer := r.Referer(); referer != "" {
		keyvals = append(keyvals, "referer", referer)
	}

	// check both the finalizer context key and the go-kit one.
	if size, ok := finalizer.ResponseSize(ctx); ok {
		keyvals = append(keyvals, "response_size", size)
	}

	l.logger.Log(keyvals...)
}
