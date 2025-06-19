package route

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
)

type Deps struct {
	Router   chi.Router
	Logger   *slog.Logger
	Handlers *restHTTP.Handlers
}
