package app

import (
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/config"
)

type Deps struct {
	Config *config.Config
	Logger *slog.Logger
}
