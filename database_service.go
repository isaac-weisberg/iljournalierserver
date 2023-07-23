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

	err = migrateDatabase(ctx, db)
	if err != nil {
		return nil, err
	}

	databaseService := databaseService{db}

	return &databaseService, nil
}
