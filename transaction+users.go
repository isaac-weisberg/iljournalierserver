package main

import "fmt"

func (transaction transaction) createUsersTable() error {
	sql := "CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, magicKey TEXT NOT NULL UNIQUE)"

	_, err := transaction.exec(sql)
	if err != nil {
		return j(err, "create table failed")
	}

	return nil
}

func (transaction transaction) createUser(magicKey string) (*int64, error) {
	sql := "INSERT INTO users (magicKey) VALUES (?)"

	result, err := transaction.exec(sql, magicKey)
	if err != nil {
		return nil, j(err, fmt.Sprintf("insert failed %s", magicKey))
	}

	lastIndertedId, err := result.LastInsertId()
	if err != nil {
		return nil, j(err, "last inserted id failed")
	}

	return &lastIndertedId, nil
}
