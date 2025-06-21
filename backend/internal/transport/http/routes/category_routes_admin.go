package route

import (
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	"github.com/go-chi/chi/v5"
)

func registerCategoryAdminRoutes(r chi.Router, h *restHTTP.CategoryHandler, cah *restHTTP.CategoryAttributeHandler) {
	r.Route("/category", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)

		r.Route("/{categoryID}/attribute", func(r chi.Router) {
			r.Get("/", cah.GetByCategory)
			r.Post("/", cah.Create)
			r.Put("/{attributeID}", cah.Update)
			r.Delete("/{attributeID}", cah.Delete)
		})
	})
}
