package main

import (
	"context"
	"database/sql"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

type transaction struct {
	tx  *sql.Tx
	ctx context.Context
}

func newTransaction(ctx context.Context, tx *sql.Tx) transaction {
	return transaction{ctx: ctx, tx: tx}
}

func (transaction *transaction) exec(query string, args ...any) (sql.Result, error) {
	return transaction.tx.ExecContext(transaction.ctx, query, args...)
}

func (transaction *transaction) queryRow(query string, args ...any) *sql.Row {
	return transaction.tx.QueryRowContext(transaction.ctx, query, args...)
}

func txQuery[R interface{}](transaction *transaction, query string, args []any, block func(rows *sql.Rows) (*R, error)) (*R, error) {
	rows, err := transaction.tx.QueryContext(transaction.ctx, query, args...)
	if err != nil {
		return nil, errors.J(err, "query context failed")
	}
	defer rows.Close()

	res, err := block(rows)

	if err != nil {
		return res, errors.J(err, "query block failed")
	}

	return res, nil
}
