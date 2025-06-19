package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/Neimess/zorkin-store-project/internal/config"
	storage "github.com/Neimess/zorkin-store-project/internal/storage/psql"
	"github.com/Neimess/zorkin-store-project/migrations"
	"github.com/Neimess/zorkin-store-project/pkg/args"
	"github.com/Neimess/zorkin-store-project/pkg/migrator"
)

func runMigrations(cfg *config.Config, log *slog.Logger) error {
	start := time.Now()
	provider := storage.NewProvider(cfg)

	m, err := migrator.New(
		migrator.WithDatabase(provider),
		migrator.WithSource(migrator.EmbeddedSource{FS: migrations.FS, Directory: "."}),
	)
	if err != nil {
		log.Error("migrator failed", slog.Any("error", err))
		os.Exit(1)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("db ping timeout",
				slog.String("host", cfg.Storage.Host),
				slog.Int("port", cfg.Storage.Port),
				slog.Duration("timeout", cfg.Storage.BusyTimeout),
				slog.Any("error", err),
			)
		} else {
			log.Error("migrations failed", slog.Any("error", err))
		}
		return err
	}
	v, _ := m.Version()
	log.Info("migrations applied",
		slog.Uint64("version", uint64(v)),
		slog.Duration("took", time.Since(start)),
	)
	return nil
}

func runForce(cfg *config.Config, log *slog.Logger, version int) error {
	provider := storage.NewProvider(cfg)

	m, err := migrator.New(
		migrator.WithDatabase(provider),
		migrator.WithSource(migrator.EmbeddedSource{
			FS: migrations.FS, Directory: ".",
		}),
	)
	if err != nil {
		log.Error("migrator init failed", slog.Any("error", err))
		return err
	}

	if err := m.Force(version); err != nil {
		log.Error("migrator force failed", slog.Any("error", err))
		return err
	}
	log.Info("migration force successful", slog.Int("version", version))
	return nil
}

func handleMigrations(a *args.Args, cfg *config.Config, log *slog.Logger) error {
	if a.Force != 0 {
		if err := runForce(cfg, log, a.Force); err != nil {
			return err
		}
		if !a.Migrate {
			log.Info("force only requested - exiting")
			return nil
		}
	}

	if a.Migrate {
		return runMigrations(cfg, log)
	}
	return nil
}
