package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Neimess/zorkin-store-project/internal/config"
	storage "github.com/Neimess/zorkin-store-project/internal/storage/psql"
	"github.com/Neimess/zorkin-store-project/migrations"
	"github.com/Neimess/zorkin-store-project/pkg/args"
	"github.com/Neimess/zorkin-store-project/pkg/migrator"
	"github.com/golang-migrate/migrate/v4"
)

func runMigrations(cfg *config.Config, log *slog.Logger) error {
	start := time.Now()
	provider := storage.NewProvider(cfg)

	m, err := migrator.New(
		migrator.WithDatabase(provider),
		migrator.WithSource(migrator.EmbeddedSource{FS: migrations.FS, Directory: "."}),
		migrator.WithLogger(log),
	)
	if err != nil {
		log.Error("migrator failed", slog.Any("error", err))
		os.Exit(1)
	}

	err = m.Up()
	switch {
	case err == nil:
		v, _ := m.Version()
		log.Info("migrations applied",
			slog.Uint64("version", uint64(v)),
			slog.Duration("took", time.Since(start)),
		)

	case errors.Is(err, migrate.ErrNoChange):
		return nil

	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("db ping timeout: %w", err)

	default:
		return fmt.Errorf("migrations failed: %w", err)
	}
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

func runDown(cfg *config.Config, log *slog.Logger) error {
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

	if err := m.Down(); err != nil {
		log.Error("migrator down failed", slog.Any("error", err))
		return err
	}
	log.Info("migration rolled back successfully")
	return nil
}

func handleMigrations(a *args.Args, cfg *config.Config, log *slog.Logger) error {
	switch {
	case a.Force != 0:
		if err := runForce(cfg, log, a.Force); err != nil {
			return err
		}
		if !a.Up {
			log.Info("force only requested - exiting")
			return nil
		}
		return runMigrations(cfg, log)

	case a.Down:
		return runDown(cfg, log)

	case a.Up:
		return runMigrations(cfg, log)

	default:
		log.Info("no migration flag provided - nothing to do")
		return nil
	}
}
