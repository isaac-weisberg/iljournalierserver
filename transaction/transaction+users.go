package transaction

import (
	"database/sql"
	"fmt"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

func (transaction *Transaction) CreateUsersTable() error {
	sql := "CREATE TABLE users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, magicKey TEXT NOT NULL UNIQUE)"

	_, err := transaction.Exec(sql)
	if err != nil {
		return errors.J(err, "create table failed")
	}

	return nil
}

func (transaction *Transaction) CreateUser(magicKey string) (*int64, error) {
	sql := "INSERT INTO users (magicKey) VALUES (?)"

	result, err := transaction.Exec(sql, magicKey)
	if err != nil {
		return nil, errors.J(err, fmt.Sprintf("insert failed %s", magicKey))
	}

	lastIndertedId, err := result.LastInsertId()
	if err != nil {
		return nil, errors.J(err, "last inserted id failed")
	}

	return &lastIndertedId, nil
}

func (transaction *Transaction) FindUserForMagicKey(magicKey string) (*int64, error) {
	query := "SELECT (id) FROM users WHERE magicKey = ?"

	userIds, err := TxQuery[[]int64](transaction, query, []any{magicKey}, func(rows *sql.Rows) (*[]int64, error) {
		var userIds []int64
		for rows.Next() {
			var userId int64
			err := rows.Scan(&userId)
			if err != nil {
				return nil, errors.J(err, "scanning row failed")
			}
			userIds = append(userIds, userId)
		}

		err := rows.Err()
		if err != nil {
			return nil, errors.J(err, "rows returned error")
		}

		return &userIds, nil
	})

	if err != nil {
		return nil, errors.J(err, "txQuery failed")
	}

	switch len(*userIds) {
	case 0:
		return nil, nil
	case 1:
		firstUserId := (*userIds)[0]
		return &firstUserId, nil
	default:
		firstUserId := (*userIds)[0]
		return &firstUserId, errors.E("multiple users found for magic key, which is unexpected")
	}
}
