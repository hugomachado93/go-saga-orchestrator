package repository

import (
	"context"
	"database/sql"
)

type Transaction struct {
	db *sql.DB
}

func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{db: db}
}

type txKey struct{}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func a[T any]() {
	return
}

func (t *Transaction) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	var err error
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})

	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(injectTx(ctx, tx))
	return err
}
