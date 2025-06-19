package storage

import (
	"context"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/internal/config"
	"github.com/Neimess/zorkin-store-project/pkg/database/psql"
	migrateDB "github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/pgx"
)

const (
	driverName = "pgx"
)

type Provider struct {
	cfg *config.Config
}

func NewProvider(cfg *config.Config) *Provider { return &Provider{cfg: cfg} }

func (p *Provider) Driver() (driver migrateDB.Driver, name string, cleanup func() error, err error) {
	db, err := psql.New(
		context.Background(),
		psql.WithHost(p.cfg.Storage.Host),
		psql.WithPort(p.cfg.Storage.Port),
		psql.WithUser(p.cfg.Storage.User),
		psql.WithPassword(p.cfg.Storage.Password),
		psql.WithDB(p.cfg.Storage.DBName),
		psql.WithSSL(p.cfg.Storage.SSLMode),
		psql.WithMaxConns(p.cfg.Storage.MaxOpenConns),
		psql.WithMinConns(p.cfg.Storage.MaxIdleConns),
		psql.WithConnLifetime(p.cfg.Storage.ConnMaxLifetime),
		psql.WithConnectTimeout(p.cfg.Storage.BusyTimeout),
	)
	if err != nil {
		return nil, "", nil, err
	}

	driver, err = pgx.WithInstance(db.DB, &pgx.Config{})
	if err != nil {
		if cerr := db.Close(); cerr != nil {
			slog.Warn("failed to close DB after driver init failure", slog.Any("error", cerr))
		}
		return nil, "", nil, err
	}

	cleanup = func() error {
		if err := db.Close(); err != nil {
			slog.Warn("failed to close DB", slog.Any("error", err))
			return err
		}
		return nil
	}
	return driver, driverName, cleanup, nil
}
