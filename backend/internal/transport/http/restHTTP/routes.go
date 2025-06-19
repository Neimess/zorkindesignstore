package restHTTP

import (
	"log/slog"
	"net/http"
	"net/http/pprof"
	"time"

	_ "github.com/Neimess/zorkin-store-project/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *Handlers) RegisterRoutes(r chi.Router, logger *slog.Logger) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Compress(5))
	r.Use(httplog.RequestLogger(logger, &httplog.Options{
		RecoverPanics:      true,
		LogRequestHeaders:  []string{"Origin"},
		LogResponseHeaders: []string{},
	}))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		})
	})

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Mount("/debug/pprof", profiler())
}

func profiler() http.Handler {
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
}
