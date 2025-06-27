package route

import (
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/product"
	"github.com/go-chi/chi/v5"
)

func registerProductPublicRoutes(r chi.Router, h *product.Handler) {
	r.Route("/product", func(r chi.Router) {
		r.Get("/category/{id}", h.ListByCategory)
		r.Get("/{id}", h.GetDetailed)
	})
}
