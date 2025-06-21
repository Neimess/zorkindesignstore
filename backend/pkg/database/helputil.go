package database

import (
	"context"
	"log/slog"
	"time"
)

func WithQuery(ctx context.Context, log *slog.Logger, query string, fn func() error, extras ...slog.Attr) error {
	start := time.Now()

	err := fn()

	attrs := []slog.Attr{
		slog.String("query", query),
		slog.Duration("duration", time.Since(start)),
	}
	attrs = append(attrs, extras...)
	if err != nil {
		attrs = append(attrs, slog.String("error", err.Error()))
	}

	log.LogAttrs(ctx, slog.LevelDebug, "db query", attrs...)
	return err
}
