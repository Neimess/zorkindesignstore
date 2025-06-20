package repository

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Deps struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}
