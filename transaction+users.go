package main

import "fmt"

func (transaction *transaction) createUsersTable() error {
	sql := "CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, magicKey TEXT NOT NULL UNIQUE)"

	_, err := transaction.exec(sql)
	if err != nil {
		return j(err, "create table failed")
	}

	return nil
}

func (transaction *transaction) createUser(magicKey string) (*int64, error) {
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

func (transaction *transaction) findUserForMagicKey(magicKey string) (*int64, error) {
	sql := "SELECT (id) FROM users WHERE magicKey = ?"

	rows, err := transaction.query(sql, magicKey)
	if err != nil {
		return nil, j(err, "selecting failed")
	}

	var userIds []int64
	for rows.Next() {
		var userId int64
		err = rows.Scan(&userId)
		if err != nil {
			return nil, j(err, "scanning row failed")
		}
		userIds = append(userIds, userId)
	}

	err = rows.Err()
	if err != nil {
		return nil, j(err, "rows returned error")
	}

	switch len(userIds) {
	case 0:
		return nil, nil
	case 1:
		return &userIds[0], nil
	default:
		return &userIds[0], e("multiple users found for magic key, which is unexpected")
	}
}
