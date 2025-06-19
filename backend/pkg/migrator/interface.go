package migrator

import (
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
)

type Source interface {
	Driver() (source.Driver, string, error)
}

type Database interface {
	Driver() (database.Driver, string, func() error, error)
}
