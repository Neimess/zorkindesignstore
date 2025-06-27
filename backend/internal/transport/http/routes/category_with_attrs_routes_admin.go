package route

import (
	attributeH "github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/attribute"
	categoryH "github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/category"
	"github.com/go-chi/chi/v5"
)

func registerCategoryWithAttrsAdminRoutes(r chi.Router, h *categoryH.Handler, ah *attributeH.Handler) {
	r.Route("/category", func(r chi.Router) {
		r.Post("/", h.CreateCategory)
		r.Put("/{id}", h.UpdateCategory)
		r.Delete("/{id}", h.DeleteCategory)
		r.Get("/{id}", h.GetCategory)
		r.Get("/", h.ListCategories)

		r.Route("/{categoryID}/attribute", func(r chi.Router) {
			r.Post("/", ah.CreateAttribute)
			r.Get("/", ah.ListAttributes)
			r.Get("/{id}", ah.GetAttribute)
			r.Put("/{id}", ah.UpdateAttribute)
			r.Delete("/{id}", ah.DeleteAttribute)
			r.Post("/batch", ah.CreateAttribute)
		})
	})

}
