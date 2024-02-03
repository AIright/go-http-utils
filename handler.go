package go_http_utils

import (
	"net/http"
	"net/http/pprof"
	rpprof "runtime/pprof"
)

func handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// goroutine, threadcreate, heap, allocs, block, mutex
	for _, p := range rpprof.Profiles() {
		mux.Handle("/debug/pprof/"+p.Name(), pprof.Handler(p.Name()))
	}

	// readiness
	mux.Handle(readinessProbeEndpoint, readinessProbe())

	return mux
}
