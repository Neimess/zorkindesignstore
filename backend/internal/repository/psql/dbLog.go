// helper.go
package repository

import (
	"context"
	"log/slog"
	"time"
)

func (r *ProductRepository) dbLog(ctx context.Context, query string) func(err error, extras ...slog.Attr) {
	start := time.Now()

	reqID, _ := ctx.Value("request_id").(string)

	return func(err error, extras ...slog.Attr) {
		attrs := []slog.Attr{
			slog.String("query", query),
			slog.Duration("duration", time.Since(start)),
		}
		if reqID != "" {
			attrs = append(attrs, slog.String("request_id", reqID))
		}
		attrs = append(attrs, extras...)

		level := slog.LevelDebug
		if err != nil {
			attrs = append(attrs, slog.String("error", err.Error()))
		}
		args := make([]any, 0, len(attrs))
		for _, attr := range attrs {
			args = append(args, attr)
		}
		r.log.Log(ctx, level, "db query", args...)
	}
}
