package storage

import (
	"context"

	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/pkg/database/psql"
	"github.com/jmoiron/sqlx"
)

func New(
	ctx context.Context,
	cfg *config.Config,
) (*sqlx.DB, error) {
	return psql.New(
		ctx,
		psql.WithHost(cfg.Storage.Host),
		psql.WithUser(cfg.Storage.User),
		psql.WithPassword(cfg.Storage.Password),
		psql.WithPort(cfg.Storage.Port),
		psql.WithConnLifetime(cfg.Storage.ConnMaxLifetime),
		psql.WithMaxConns(cfg.Storage.MaxOpenConns),
		psql.WithSSL(cfg.Storage.SSLMode),
	)
}
