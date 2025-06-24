package tx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TxFunc[T any] func(*sqlx.Tx) (T, error)

func RunInTx[T any](
	ctx context.Context,
	db *sqlx.DB,
	fn TxFunc[T],
) (T, error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	var zero T
	if err != nil {
		return zero, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	res, err := fn(tx)
	if err != nil {
		tx.Rollback()
		return zero, err
	}

	if err := tx.Commit(); err != nil {
		return zero, err
	}
	return res, nil
}
