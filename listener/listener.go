package listener

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"time"

	"go.uber.org/zap"
)

const metricPanicCounter = "service.api.panic.total"

type HTTPListener struct {
	options
}

func New(opts ...Option) *HTTPListener {
	o := options{
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &HTTPListener{
		options: o,
	}
}

func (s *HTTPListener) Listen(ctx context.Context, port int, handler http.Handler) error {
	for _, mw := range s.mw {
		handler = mw(handler)
	}
	handler = panicMW(s.logger, s.metrics, handler)

	server := &http.Server{
		Addr:           ":" + strconv.Itoa(port),
		Handler:        handler,
		ReadTimeout:    defaultReadTimeout,
		WriteTimeout:   defaultWriteTimeout,
		IdleTimeout:    defaultIdleTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ErrorLog:       errorLog(s.logger),
	}

	for _, opt := range s.server {
		opt(server)
	}

	chErrors := make(chan error)
	chSignals := make(chan os.Signal, 2)
	signal.Notify(chSignals, shutdownSignals...)

	go listen(chErrors, server)
	s.logger.Info(fmt.Sprintf("serving app on %s", server.Addr))

	var err error
	select {
	case err = <-chErrors:
		_ = shutdown(server, s.shutdownTimeout)
	case <-chSignals:
		signal.Stop(chSignals)
		err = shutdown(server, s.shutdownTimeout)
	case <-ctx.Done():
		err = ctx.Err()
		if e := shutdown(server, s.shutdownTimeout); e != nil {
			err = e
		}
	}

	close(chErrors)
	close(chSignals)

	return err
}

func listen(ch chan error, server *http.Server) {
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		ch <- err
	}
}

func shutdown(server *http.Server, timeout time.Duration) error {
	var cancel func()
	ctx := context.Background()

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	return server.Shutdown(ctx)
}

func panicMW(logger *zap.Logger, metrics Metrics, handler http.Handler) http.Handler {
	if logger == nil {
		logger = zap.NewNop()
	}
	if metrics == nil {
		metrics = stubMetrics{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(
					fmt.Sprintf("handler panic: %+v\n%s", err, string(debug.Stack())),
					zap.Error(err.(error)),
					zap.String("url", r.URL.String()),
					zap.String("method", r.Method),
				)
				metrics.Increment(metricPanicCounter)
				panic(err)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

func errorLog(logger *zap.Logger) *log.Logger {
	if logger == nil {
		return log.New(os.Stderr, "", log.LstdFlags)
	}

	return log.New(writerFunc(func(p []byte) (n int, err error) {
		l := len(p)

		if bytes.HasPrefix(p, []byte("http: panic serving ")) {
			// Skip logging of panic, handled by server itself.
			// This kind of errors would be logged by middlewares
			return l, nil
		}

		p = bytes.TrimRight(p, "\n")

		logger.Error(string(p))

		return l, nil
	}), "", 0)
}
