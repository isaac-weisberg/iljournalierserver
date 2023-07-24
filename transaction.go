package main

import (
	"context"
	"database/sql"
)

type transaction struct {
	tx  *sql.Tx
	ctx context.Context
}

func newTransaction(ctx context.Context, tx *sql.Tx) transaction {
	return transaction{ctx: ctx, tx: tx}
}

func (transaction *transaction) rollBack() error {
	return transaction.tx.Rollback()
}

func (transaction *transaction) exec(query string, args ...any) (sql.Result, error) {
	return transaction.tx.ExecContext(transaction.ctx, query, args...)
}

func (transaction *transaction) queryRow(query string, args ...any) *sql.Row {
	return transaction.tx.QueryRowContext(transaction.ctx, query, args...)
}

func (transaction *transaction) query(query string, args ...any) (*sql.Rows, error) {
	return transaction.tx.QueryContext(transaction.ctx, query, args...)
}

func (transaction *transaction) commit() error {
	return transaction.tx.Commit()
}
