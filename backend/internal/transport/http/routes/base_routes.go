package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func registerBaseRoutes(r chi.Router) {
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}
