package migrator

import (
	"io/fs"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type EmbeddedSource struct {
	FS        fs.FS
	Directory string
}

func (e EmbeddedSource) Driver() (source.Driver, string, error) {
	drv, err := iofs.New(e.FS, e.Directory)
	return drv, "iofs", err
}
