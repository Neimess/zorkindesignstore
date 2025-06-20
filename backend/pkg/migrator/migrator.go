package migrator

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
)

type Opt func(*Migrator)

func WithSource(src Source) Opt    { return func(m *Migrator) { m.src = src } }
func WithDatabase(db Database) Opt { return func(m *Migrator) { m.db = db } }

type Migrator struct {
	src    Source
	db     Database
	logger *slog.Logger
}

func New(opts ...Opt) (*Migrator, error) {
	m := &Migrator{}
	for _, opt := range opts {
		opt(m)
	}

	switch {
	case m.src == nil:
		return nil, fmt.Errorf("migrator: source not provided")
	case m.db == nil:
		return nil, fmt.Errorf("migrator: database not provided")
	}
	return m, nil
}

func WithLogger(log *slog.Logger) Opt {
	return func(m *Migrator) { m.logger = log }
}

func (m *Migrator) Up() error {
	return m.runRaw(func(mg *migrate.Migrate) error { return mg.Up() })
}

func (m *Migrator) Down() error {
	return m.runRaw(func(mg *migrate.Migrate) error { return mg.Steps(-1) })
}

func (m *Migrator) Steps(n int) error {
	return m.runRaw(func(mg *migrate.Migrate) error { return mg.Steps(n) })
}

func (m *Migrator) Force(version int) error {
	return m.runRaw(func(mg *migrate.Migrate) error {
		if err := mg.Force(version); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("force failed: %w", err)
		}
		return nil
	})
}

func (m *Migrator) Version() (version uint, err error) {
	err = m.runRaw(func(mg *migrate.Migrate) error {
		version, _, err = mg.Version()
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("get version: %w", err)
	}
	return version, nil
}

func (m *Migrator) runRaw(f func(*migrate.Migrate) error) error {
	dbDrv, dbName, closeFn, err := m.db.Driver()
	if err != nil {
		return fmt.Errorf("db driver: %w", err)
	}
	if closeFn != nil {
		defer closeFn()
	}

	srcDrv, srcName, err := m.src.Driver()
	if err != nil {
		return fmt.Errorf("source driver: %w", err)
	}

	mg, err := migrate.NewWithInstance(srcName, srcDrv, dbName, dbDrv)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	if m.logger != nil {
		mg.Log = &migrateLogger{log: m.logger}
	}

	return f(mg)
}
