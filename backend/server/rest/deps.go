package rest

import (
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/internal/transport/http/restHTTP"
)

type Deps struct {
	Server   config.HTTPServer
	Config   *config.Config
	Logger   *slog.Logger
	Handlers *restHTTP.Handlers
}
