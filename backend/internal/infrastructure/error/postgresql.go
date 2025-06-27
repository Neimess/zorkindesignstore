package repoError

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Neimess/zorkin-store-project/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

func MapPostgreSQLError(err error) error {
	if err == nil {
		return nil
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {

		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("%w: %s", domain.ErrConflict, pgErr.Detail)
		case "23502": // not_null_violation
			return fmt.Errorf("%w: %s", domain.ErrBadRequest, pgErr.ColumnName)
		case "23503": // foreign_key_violation
			return fmt.Errorf("%w: %s", domain.ErrNotFound, pgErr.Detail)
		case "23514": // check_violation
			return fmt.Errorf("%w: %s", domain.ErrValidation, pgErr.Detail)
		case "22001": // string_data_right_truncation
			return fmt.Errorf("%w: %s", domain.ErrBadRequest, pgErr.ColumnName)
		case "42P01": // undefined_table
			return fmt.Errorf("%w: relation does not exist: %s", domain.ErrInternal, pgErr.Message)
		default:
			return fmt.Errorf("%w: %s", domain.ErrInternal, pgErr.Message)
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: record not found", domain.ErrNotFound)
	}
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("%w: operation canceled", domain.ErrCanceled)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: operation timed out", domain.ErrTimeout)
	}

	return fmt.Errorf("%w: %v", domain.ErrInternal, err)
}
