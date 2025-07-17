package repoError

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Neimess/zorkin-store-project/pkg/app_error"
	"github.com/jackc/pgx/v5/pgconn"
)

// MapPostgreSQLError проверяет err и возвращает только ваше приложение-определённое
// исключение, а подробности оригинальной ошибки логирует через slog.
func MapPostgreSQLError(logger *slog.Logger, err error) error {
	if err == nil {
		return nil
	}
	logger.Error("mapping PostgreSQL error",
		slog.String("error", err.Error()),
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Логируем все поля PG-ошибки
		logger.Error("PostgreSQL error",
			slog.String("code", pgErr.Code),
			slog.String("constraint", pgErr.ConstraintName),
			slog.String("table", pgErr.TableName),
			slog.String("column", pgErr.ColumnName),
			slog.String("detail", pgErr.Detail),
			slog.String("message", pgErr.Message),
		)

		switch pgErr.Code {
		case "23505":
			return app_error.ErrConflict
		case "23502":
			return app_error.ErrBadRequest
		case "23503":
			return app_error.ErrNotFound
		case "23514":
			return app_error.ErrValidation
		case "22001":
			return app_error.ErrBadRequest
		case "42P01":
			return app_error.ErrInternal
		default:
			return app_error.ErrInternal
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("no rows in result set", slog.String("error", err.Error()))
		return app_error.ErrNotFound
	}

	if errors.Is(err, context.Canceled) {
		logger.Debug("operation canceled by context", slog.String("error", err.Error()))
		return app_error.ErrCanceled
	}

	if errors.Is(err, context.DeadlineExceeded) {
		logger.Debug("operation deadline exceeded", slog.String("error", err.Error()))
		return app_error.ErrTimeout
	}

	if errors.Is(err, app_error.ErrBadRequest) {
		logger.Debug("bad request error", slog.String("error", err.Error()))
		return app_error.ErrBadRequest
	}

	if errors.Is(err, app_error.ErrNotFound) {
		logger.Debug("not found error", slog.String("error", err.Error()))
		return app_error.ErrNotFound
	}
	// Всё прочее — внутренняя ошибка хранилища
	logger.Error("unexpected database error", slog.String("error", err.Error()))
	return app_error.ErrInternal
}
