package migrator

import (
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/file"
)

type PathSource struct {
	Path string
}

func (p PathSource) Driver() (source.Driver, string, error) {
	file := &file.File{}
	d, err := file.Open(p.Path)
	return d, "file", err
}
