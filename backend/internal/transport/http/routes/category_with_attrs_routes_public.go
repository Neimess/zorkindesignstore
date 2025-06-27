package route

import (
	attribute "github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/attribute"
	category "github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/category"
	"github.com/go-chi/chi/v5"
)

func registerCategoryWithAttrsPublicRoutes(r chi.Router, h *category.Handler, ah *attribute.Handler) {
	r.Route("/category", func(r chi.Router) {
		r.Get("/{id}", h.GetCategory)
		r.Get("/", h.ListCategories)

		r.Route("/{categoryID}/attribute", func(r chi.Router) {
			r.Get("/", ah.ListAttributes)
			r.Get("/{id}", ah.GetAttribute)
		})
	})
}
