package route

import (
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
	"github.com/go-chi/chi/v5"
)

type Deps struct {
	Router   chi.Router
	Logger   *slog.Logger
	Handlers *restHTTP.Handlers
}
