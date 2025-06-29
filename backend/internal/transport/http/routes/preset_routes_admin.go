package route

import (
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/preset"
	"github.com/go-chi/chi/v5"
)

func registerPresetAdminRoutes(r chi.Router, h *preset.Handler) {
	r.Route("/presets", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Delete("/{id}", h.Delete)
		r.Put("/{id}", h.Update)
	})
}
