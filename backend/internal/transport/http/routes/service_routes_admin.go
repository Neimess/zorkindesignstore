package route

import (
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/service"
	"github.com/go-chi/chi/v5"
)

func registerServiceAdminRoutes(r chi.Router, h *service.Handler) {
	r.Route("/services", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
