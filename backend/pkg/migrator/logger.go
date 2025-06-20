package migrator

import (
	"fmt"
	"log/slog"
)

type migrateLogger struct {
	log *slog.Logger
}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	l.log.Info("migration step", slog.String("message", fmt.Sprintf(format, v...)))
}

func (l *migrateLogger) Verbose() bool {
	return true
}
