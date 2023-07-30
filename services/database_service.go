package services

import (
	"context"
	"database/sql"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/intake"
	"caroline-weisberg.fun/iljournalierserver/transaction"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseService struct {
	db *sql.DB
}

func NewDatabaseService(intakeConfig *intake.IntakeConfiguration) (*DatabaseService, error) {
	db, err := sql.Open("sqlite3", intakeConfig.DbPath)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	databaseService := DatabaseService{db}

	return &databaseService, nil
}

func BeginTxBlockVoid(databaseService *DatabaseService, ctx context.Context, block func(tx *transaction.Transaction) error) error {
	_, err := BeginTxBlock[interface{}](databaseService, ctx, func(tx *transaction.Transaction) (*interface{}, error) {
		return nil, block(tx)
	})
	return err
}

func BeginTxBlock[R interface{}](databaseService *DatabaseService, ctx context.Context, block func(tx *transaction.Transaction) (*R, error)) (*R, error) {
	tx, err := databaseService.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, errors.J(err, "begin tx failed")
	}

	transaction := transaction.NewTransaction(ctx, tx)

	res, err := block(&transaction)

	if err != nil {
		var blockError = errors.J(err, "transaction block failed")
		rollbackError := transaction.Rollback()
		if rollbackError != nil {
			return nil, errors.Js(rollbackError, blockError)
		}
		return nil, blockError
	}

	err = transaction.Commit()
	if err != nil {
		return res, errors.J(err, "commit failed")
	}

	return res, nil
}
