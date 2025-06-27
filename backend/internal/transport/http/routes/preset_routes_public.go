package route

import (
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP/preset"
	"github.com/go-chi/chi/v5"
)

func registerPresetPublicRoutes(r chi.Router, h *preset.Handler) {
	r.Route("/presets", func(r chi.Router) {
		r.Get("/", h.ListShort)
		r.Get("/detailed", h.ListDetailed)
		r.Get("/{id}", h.Get)
	})
}
