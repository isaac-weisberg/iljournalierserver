package main

import (
	"fmt"
)

func (transaction *transaction) createAccessTokensTable() error {
	sql := `CREATE TABLE IF NOT EXISTS accessTokens 
		(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL, 
			token TEXT NOT NULL UNIQUE, 
			FOREIGN KEY(userId) REFERENCES users(id)
		)`

	_, err := transaction.exec(sql)
	if err != nil {
		return j(err, "create table failed")
	}

	return nil
}

func (transaction *transaction) createAccessToken(userId int64, accessToken string) error {
	sql := "INSERT INTO accessTokens (userId, token) VALUES (?, ?)"

	_, err := transaction.exec(sql, userId, accessToken)
	if err != nil {
		return j(err, fmt.Sprintf("insert failed userId=%v accessToken=%v", userId, accessToken))
	}

	return nil
}

func (transaction *transaction) findUserIdForAccessToken(accessToken string) (*int64, error) {
	sql := "SELECT userId FROM accessTokens WHERE token = ?"

	rows, err := transaction.query(sql, accessToken)
	if err != nil {
		return nil, j(err, "selecting userId by token failed")
	}
	defer rows.Close()

	var userIds []int64
	for rows.Next() {
		var userId int64
		err := rows.Scan(&userId)
		if err != nil {
			return nil, j(err, "scanning rows failed")
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
		return &userIds[0], e("multiple users for this access token - this is bad")
	}
}
