package route

import (
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/service"
	"github.com/go-chi/chi/v5"
)

func registerServicePublicRoutes(r chi.Router, h *service.Handler) {
	r.Route("/services", func(r chi.Router) {
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
	})
}
