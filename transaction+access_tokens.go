package main

import (
	"database/sql"
	"fmt"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

func (transaction *transaction) createAccessTokensTable() error {
	sql := `CREATE TABLE accessTokens 
		(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL, 
			token TEXT NOT NULL UNIQUE, 
			FOREIGN KEY(userId) REFERENCES users(id)
		)`

	_, err := transaction.exec(sql)
	if err != nil {
		return errors.J(err, "create table failed")
	}

	return nil
}

func (transaction *transaction) createAccessToken(userId int64, accessToken string) error {
	query := "INSERT INTO accessTokens (userId, token) VALUES (?, ?)"

	_, err := transaction.exec(query, userId, accessToken)
	if err != nil {
		return errors.J(err, fmt.Sprintf("insert failed userId=%v accessToken=%v", userId, accessToken))
	}

	return nil
}

func (transaction *transaction) findUserIdForAccessToken(accessToken string) (*int64, error) {
	query := "SELECT userId FROM accessTokens WHERE token = ?"

	userIds, err := txQuery[[]int64](transaction, query, []any{accessToken}, func(rows *sql.Rows) (*[]int64, error) {
		var userIds []int64
		for rows.Next() {
			var userId int64
			err := rows.Scan(&userId)
			if err != nil {
				return nil, errors.J(err, "scanning rows failed")
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
		return &firstUserId, errors.E("multiple users for this access token - this is bad")
	}
}
