package route

import (
    "github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
    "github.com/go-chi/chi/v5"
)

func registerCategoryPublicRoutes(r chi.Router, h *restHTTP.CategoryHandler) {
    r.Route("/category", func(r chi.Router) {
        r.Get("/{id}", h.Get)
        r.Get("/", h.List)
    })
}
