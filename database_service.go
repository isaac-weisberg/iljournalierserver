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

func (databaseService *databaseService) beginTx(ctx context.Context) (*transaction, error) {
	tx, err := databaseService.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	transaction := newTransaction(ctx, tx)

	return &transaction, nil
}
