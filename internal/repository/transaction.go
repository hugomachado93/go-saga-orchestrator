package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) *Transaction {
	return &Transaction{db: db}
}

type txKey struct{}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// extractTx extracts transaction from context
func extractTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}

func (t *Transaction) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	var err error
	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})

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
