package route

import (
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	"github.com/go-chi/chi/v5"
)

func registerProductAdminRoutes(r chi.Router, h *restHTTP.ProductHandler) {
	r.Route("/product", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Post("/detailed", h.CreateWithAttributes)
		// r.Put("/{id}",    h.Update)
		// r.Delete("/{id}", h.Delete)
		r.Get("/category/{id}", h.ListByCategory)
		r.Get("/{id}", h.GetDetailed)
	})
}
