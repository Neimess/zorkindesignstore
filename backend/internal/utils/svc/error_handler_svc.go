package utils

import (
	"errors"
	"fmt"
	"log/slog"
)

// ErrorHandler - универсальный обработчик ошибок для сервисов
func ErrorHandler(log *slog.Logger, op string, err error, mapping map[error]error) error {
	for src, dst := range mapping {
		if errors.Is(err, src) {
			log.Warn("handled error", slog.String("op", op), slog.Any("error", err))
			return dst
		}
	}
	log.Error("unhandled error", slog.String("op", op), slog.Any("error", err))
	return fmt.Errorf("%s: %w", op, err)
}
