package route

import (
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	"github.com/go-chi/chi/v5"
)

func registerPresetAdminRoutes(r chi.Router, h *restHTTP.PresetHandler) {
	r.Route("/presets", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Delete("/{id}", h.Delete)
	})
}
