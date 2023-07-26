package main

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type databaseService struct {
	db *sql.DB
}

func newDatabaseService(ctx context.Context) (*databaseService, error) {
	db, err := sql.Open("sqlite3", "iljournalierAlpha")

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	databaseService := databaseService{db}

	return &databaseService, nil
}

func beginTxBlockVoid(databaseService *databaseService, ctx context.Context, block func(tx *transaction) error) error {
	_, err := beginTxBlock[interface{}](databaseService, ctx, func(tx *transaction) (*interface{}, error) {
		return nil, block(tx)
	})
	return err
}

func beginTxBlock[R interface{}](databaseService *databaseService, ctx context.Context, block func(tx *transaction) (*R, error)) (*R, error) {
	tx, err := databaseService.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, j(err, "begin tx failed")
	}

	transaction := newTransaction(ctx, tx)

	res, err := block(&transaction)

	if err != nil {
		var blockError = j(err, "transaction block failed")
		rollbackError := transaction.tx.Rollback()
		if rollbackError != nil {
			return nil, js(rollbackError, blockError)
		}
		return nil, blockError
	}

	err = transaction.tx.Commit()
	if err != nil {
		return res, j(err, "commit failed")
	}

	return res, nil
}
