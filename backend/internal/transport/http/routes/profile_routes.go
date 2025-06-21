package route

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
)

func profiler(env string) http.Handler {
	switch env {
	case "development", "test", "local":
		r := chi.NewRouter()
		r.Handle("/", http.HandlerFunc(pprof.Index))
		r.Handle("/cmdline", http.HandlerFunc(pprof.Cmdline))
		r.Handle("/profile", http.HandlerFunc(pprof.Profile))
		r.Handle("/symbol", http.HandlerFunc(pprof.Symbol))
		r.Handle("/trace", http.HandlerFunc(pprof.Trace))
		r.Handle("/allocs", pprof.Handler("allocs"))
		r.Handle("/block", pprof.Handler("block"))
		r.Handle("/goroutine", pprof.Handler("goroutine"))
		r.Handle("/heap", pprof.Handler("heap"))
		r.Handle("/mutex", pprof.Handler("mutex"))
		r.Handle("/threadcreate", pprof.Handler("threadcreate"))
		return r
	default:
		return http.NotFoundHandler()
	}
}
