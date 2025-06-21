package route

import (
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	"github.com/go-chi/chi/v5"
)

func registerProductPublicRoutes(r chi.Router, h *restHTTP.ProductHandler) {
	r.Route("/product", func(r chi.Router) {
		r.Get("/category/{id}", h.ListByCategory)
		r.Get("/{id}", h.GetDetailed)
	})
}
