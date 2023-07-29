package transaction

import (
	"context"
	"database/sql"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

type Transaction struct {
	tx  *sql.Tx
	ctx context.Context
}

func NewTransaction(ctx context.Context, tx *sql.Tx) Transaction {
	return Transaction{ctx: ctx, tx: tx}
}

func (transaction *Transaction) Exec(query string, args ...any) (sql.Result, error) {
	return transaction.tx.ExecContext(transaction.ctx, query, args...)
}

func (transaction *Transaction) QueryRow(query string, args ...any) *sql.Row {
	return transaction.tx.QueryRowContext(transaction.ctx, query, args...)
}

func (transaction *Transaction) Commit() error {
	return transaction.tx.Commit()
}

func (transaction *Transaction) Rollback() error {
	return transaction.tx.Rollback()
}

func TxQuery[R interface{}](transaction *Transaction, query string, args []any, block func(rows *sql.Rows) (*R, error)) (*R, error) {
	rows, err := transaction.tx.QueryContext(transaction.ctx, query, args...)
	if err != nil {
		return nil, errors.J(err, "tx query context creation failed")
	}
	defer rows.Close()

	res, err := block(rows)

	if err != nil {
		return res, errors.J(err, "query block failed")
	}

	return res, nil
}
