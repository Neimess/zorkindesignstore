package route

import (
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/coefficients"
	"github.com/go-chi/chi/v5"
)

func registerCoefficientsAdminRoutes(r chi.Router, h *coefficients.Handler) {
	r.Route("/coefficients", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
