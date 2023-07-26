package main

import (
	"database/sql"
	"fmt"
)

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
	query := "SELECT (id) FROM users WHERE magicKey = ?"

	userIds, err := txQuery[[]int64](transaction, query, []any{magicKey}, func(rows *sql.Rows) (*[]int64, error) {
		var userIds []int64
		for rows.Next() {
			var userId int64
			err := rows.Scan(&userId)
			if err != nil {
				return nil, j(err, "scanning row failed")
			}
			userIds = append(userIds, userId)
		}

		err := rows.Err()
		if err != nil {
			return nil, j(err, "rows returned error")
		}

		return &userIds, nil
	})

	if err != nil {
		return nil, j(err, "txQuery failed")
	}

	switch len(*userIds) {
	case 0:
		return nil, nil
	case 1:
		firstUserId := (*userIds)[0]
		return &firstUserId, nil
	default:
		firstUserId := (*userIds)[0]
		return &firstUserId, e("multiple users found for magic key, which is unexpected")
	}
}
