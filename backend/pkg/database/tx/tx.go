package tx

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

// TxFunc defines a function to run within a transaction that returns a value of type T.
// Already exists:
// type TxFunc[T any] func(*sqlx.Tx) (T, error)

// TxAction defines a function to run within a transaction that only returns an error.
type TxAction func(*sqlx.Tx) error

// RunInTx executes a TxFunc within a transaction, returning its result and any error.
func RunInTx[T any](ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) (T, error)) (result T, err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return result, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	result, err = fn(tx)
	return result, err
}

// RunInTxAction executes a TxAction within a transaction, discarding any result.
// It wraps the generic RunInTx for void operations.
func RunInTxAction(
	ctx context.Context,
	db *sqlx.DB,
	action TxAction,
) error {
	// Use struct{} as a placeholder type
	_, err := RunInTx[struct{}](ctx, db, func(tx *sqlx.Tx) (struct{}, error) {
		// Execute action and return empty struct
		err := action(tx)
		return struct{}{}, err
	})
	return err
}

// MustRunInTxAction is like RunInTxAction but panics on error.
// Useful for setup or teardown where error handling is not desired.
func MustRunInTxAction(
	ctx context.Context,
	db *sqlx.DB,
	action TxAction,
) {
	err := RunInTxAction(ctx, db, action)
	if err != nil {
		slog.With("error", err).Error("transaction failed")
	}
}

// ExecInTx is an alias for RunInTxAction.
// Provided for readability when you think in terms of Exec (no return) vs Query (with return).
func ExecInTx(
	ctx context.Context,
	db *sqlx.DB,
	execFn TxAction,
) error {
	return RunInTxAction(ctx, db, execFn)
}
