package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseService struct {
	db *sql.DB
}

func NewDatabaseService() (DatabaseService, error) {
	db, err := sql.Open("go-sqllite3", "iljournalierAlpha")

	if err != nil {
		return DatabaseService{}, err
	}

	databaseService := DatabaseService{db}

	return databaseService, nil
}
