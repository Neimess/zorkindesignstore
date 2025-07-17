package migrator

import (
	"database/sql"
	"fmt"

	"github.com/Neimess/zorkin-store-project/migrations"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type Mode int

const (
	Up Mode = iota
	DownOne
	ForceTo
)

type Options struct {
	Mode    Mode
	Version int
	Verbose bool
}

func Run(dsn string, opts Options) (err error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			if err != nil {
				err = fmt.Errorf("run error: %v; also closing db: %w", err, cerr)
			} else {
				err = fmt.Errorf("closing db: %w", cerr)
			}
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("driver: %w", err)
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	switch opts.Mode {
	case Up:
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return err
		}
	case DownOne:
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			return err
		}
	case ForceTo:
		if err := m.Force(opts.Version); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown migration mode: %v", opts.Mode)
	}

	return nil
}
