package transaction

import (
	"database/sql"
	"fmt"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

func (transaction *Transaction) CreateUsersTable() error {
	query := "CREATE TABLE users (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, publicId TEXT NOT NULL UNIQUE, magicKey TEXT NOT NULL UNIQUE, iv TEXT NOT NULL UNIQUE)"

	_, err := transaction.Exec(query)
	if err != nil {
		return errors.J(err, "create table failed")
	}

	return nil
}

func (transaction *Transaction) CreateUser(publicId string, magicKey string, iv string) (*int64, error) {
	query := "INSERT INTO users (publicId, magicKey, iv) VALUES (?, ?, ?)"

	result, err := transaction.Exec(query, publicId, magicKey, iv)
	if err != nil {
		return nil, errors.J(err, fmt.Sprintf("insert failed %s, %s, %s", publicId, magicKey, iv))
	}

	lastIndertedId, err := result.LastInsertId()
	if err != nil {
		return nil, errors.J(err, "last inserted id failed")
	}

	return &lastIndertedId, nil
}

type UserForMagicKey struct {
	Id       int64
	PublicId string
	Iv       string
}

func (transaction *Transaction) FindUserForMagicKey(magicKey string) (*UserForMagicKey, error) {
	query := "SELECT id, publicId, iv FROM users WHERE magicKey = ?"

	users, err := TxQuery[[]UserForMagicKey](transaction, query, []any{magicKey}, func(rows *sql.Rows) (*[]UserForMagicKey, error) {
		var users []UserForMagicKey
		for rows.Next() {
			var userId int64
			var publicId string
			var iv string
			err := rows.Scan(&userId, &publicId, &iv)
			if err != nil {
				return nil, errors.J(err, "scanning row failed")
			}
			users = append(users, UserForMagicKey{
				userId,
				publicId,
				iv,
			})
		}

		err := rows.Err()
		if err != nil {
			return nil, errors.J(err, "rows returned error")
		}

		return &users, nil
	})

	if err != nil {
		return nil, errors.J(err, "txQuery failed")
	}

	switch len(*users) {
	case 0:
		return nil, nil
	case 1:
		firstUserId := (*users)[0]
		return &firstUserId, nil
	default:
		firstUserId := (*users)[0]
		return &firstUserId, errors.E("multiple users found for magic key, which is unexpected")
	}
}
