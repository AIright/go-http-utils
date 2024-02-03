package go_http_utils

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ServeHealthcheck serves pprof and readiness endpoints.
func ServeHealthcheck(ctx context.Context, log *zap.Logger) {
	srv := http.Server{
		Addr:         ":" + strconv.Itoa(envInt(envReadinessPort, defaultReadinessPort)),
		Handler:      handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Minute,
		IdleTimeout:  60 * time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()
	log.Info(fmt.Sprintf("serving readiness probe and profiler on %s", srv.Addr))

	<-ctx.Done()
	srv.SetKeepAlivesEnabled(false)
	_ = srv.Shutdown(context.Background())
}

// Listen runs http server with graceful shutdown.
func Listen(ctx context.Context, httpListener Listener, handler Mux) error {
	return httpListener.Listen(
		ctx,
		envInt(envPort, defaultPort),
		handler,
	)
}
