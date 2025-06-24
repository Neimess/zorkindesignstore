package route

import (
	categoryH "github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/category"
	"github.com/go-chi/chi/v5"
)

func registerCategoryPublicRoutes(r chi.Router, h *categoryH.CategoryHandler) {
	r.Route("/category", func(r chi.Router) {
		r.Get("/{id}", h.Get)
		r.Get("/", h.List)
	})
}
