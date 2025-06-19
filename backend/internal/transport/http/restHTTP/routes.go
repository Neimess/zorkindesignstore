package restHTTP

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *Handlers) RegisterRoutes(r chi.Router, logger *slog.Logger) {
	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(logger, &httplog.Options{
		RecoverPanics:      true,
		LogRequestHeaders:  []string{"Origin"},
		LogResponseHeaders: []string{},
	}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// OpenAPI
	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		})
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// r.Route("/products", func(r chi.Router) {
	// 	r.Get("/", h.ProductHandler.List)
	// 	r.Get("/{id}", h.ProductHandler.Get)
	// 	r.Post("/", h.ProductHandler.Create)
	// 	// ...и так далее
	// })

	// r.Route("/categories", func(r chi.Router) {
	// 	r.Get("/", h.CategoryHandler.List)
	// 	r.Get("/{id}", h.CategoryHandler.Get)
	// 	r.Post("/", h.CategoryHandler.Create)
	// 	// ...
	// })
}
