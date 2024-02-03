package go_http_utils

import (
	"context"
	"net/http"
	"time"
)

type Listener interface {
	Listen(context.Context, int, http.Handler) error
}

type Mux interface {
	http.Handler
	Handle(string, http.Handler)
}

type Metrics interface {
	Gauge(string, interface{})
	Duration(string, time.Duration)
}
