package repoError

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	der "github.com/Neimess/zorkin-store-project/internal/domain/error"
	"github.com/jackc/pgx/v5/pgconn"
)

func MapPostgreSQLError(err error) error {
	if err == nil {
		return nil
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {

		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("%w: %s", der	.ErrConflict, pgErr.Detail)
		case "23502": // not_null_violation
			return fmt.Errorf("%w: %s", der.ErrBadRequest, pgErr.ColumnName)
		case "23503": // foreign_key_violation
			return fmt.Errorf("%w: %s", der.ErrNotFound, pgErr.Detail)
		case "23514": // check_violation
			return fmt.Errorf("%w: %s", der.ErrValidation, pgErr.Detail)
		case "22001": // string_data_right_truncation
			return fmt.Errorf("%w: %s", der.ErrBadRequest, pgErr.ColumnName)
		case "42P01": // undefined_table
			return fmt.Errorf("%w: relation does not exist: %s", der.ErrInternal, pgErr.Message)
		default:
			return fmt.Errorf("%w: %s", der.ErrInternal, pgErr.Message)
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: record not found", der.ErrNotFound)
	}
	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("%w: operation canceled", der.ErrCanceled)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: operation timed out", der.ErrTimeout)
	}

	return fmt.Errorf("%w: %v", der.ErrInternal, err)
}
