package listener

import (
	"net/http"
	"time"
)

type Option func(*options)

type options struct {
	server          []func(server *http.Server)
	mw              []func(handler http.Handler) http.Handler
	shutdownTimeout time.Duration
	logger          *zap.Logger
	metrics         Metrics
}

// WithReadTimeout sets the maximum duration for reading the entire request, including the body.
func WithReadTimeout(timeout time.Duration) Option {
	return func(o *options) { o.server = append(o.server, func(s *http.Server) { s.ReadTimeout = timeout }) }
}

// WithWriteTimeout sets the maximum duration before timing out writes of the response.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(o *options) { o.server = append(o.server, func(s *http.Server) { s.WriteTimeout = timeout }) }
}

// WithIdleTimeout sets the maximum amount of time to wait for the next request.
func WithIdleTimeout(timeout time.Duration) Option {
	return func(o *options) { o.server = append(o.server, func(s *http.Server) { s.IdleTimeout = timeout }) }
}

// WithShutdownTimeout sets the maximum duration for graceful shutdown (0 is no timeout).
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(o *options) { o.shutdownTimeout = timeout }
}

// WithMW sets middleware.
// Order of middleware: third(second(first(handler)))
// The last one will be called first for any incoming request.
func WithMW(mw ...func(handler http.Handler) http.Handler) Option {
	return func(o *options) { o.mw = append(o.mw, mw...) }
}

// WithLogger enables logging for MW.
func WithLogger(logger *zap.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// WithMetrics enables metrics aggregation.
func WithMetrics(metrics Metrics) Option {
	return func(o *options) { o.metrics = metrics }
}
