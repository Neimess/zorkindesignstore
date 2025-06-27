package tx

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type TxFunc[T any] func(*sqlx.Tx) (T, error)

func RunInTx[T any](
	ctx context.Context,
	db *sqlx.DB,
	fn TxFunc[T],
) (result T, err error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return result, err
	}

	defer func() {
		if p := recover(); p != nil {
			if rerr := tx.Rollback(); rerr != nil {
				slog.Error("tx rollback failed during panic", slog.Any("error", rerr))
			}
			panic(p)
		} else if err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				slog.Error("tx rollback failed", slog.Any("error", rerr))
			}
		}
	}()

	result, err = fn(tx)
	if err != nil {
		return result, err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("tx commit failed", slog.Any("error", err))
		return result, err
	}

	return result, nil
}
