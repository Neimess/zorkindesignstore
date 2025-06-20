package logger

import (
	"log/slog"
	"os"

	"github.com/MatusOllah/slogcolor"
)

const (
	ENVLocal = "local"
	ENVDev   = "development"
	ENVProd  = "production"
)

func MustInitLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case ENVLocal:
		log = slog.New(
			slogcolor.NewHandler(
				os.Stdout,
				&slogcolor.Options{
					Level:      slog.LevelDebug,
					TimeFormat: slogcolor.DefaultOptions.TimeFormat,
				},
			))
	case ENVDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case ENVProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("unknown environment: " + env)
	}

	return log
}

func WithComponent(logger *slog.Logger, name string) *slog.Logger {
	if logger == nil {
		logger = slog.Default()
	}
	return logger.With("component", name)
}
