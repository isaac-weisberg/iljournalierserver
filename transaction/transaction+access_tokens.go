package transaction

import (
	"database/sql"
	"fmt"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

func (transaction *Transaction) CreateAccessTokensTable() error {
	sql := `CREATE TABLE accessTokens 
		(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userId INTEGER NOT NULL, 
			token TEXT NOT NULL UNIQUE,
			isActive INTEGER NOT NULL,
			FOREIGN KEY(userId) REFERENCES users(id)
		)`

	_, err := transaction.Exec(sql)
	if err != nil {
		return errors.J(err, "create table failed")
	}

	return nil
}

func (transaction *Transaction) CreateAccessToken(userId int64, accessToken string) error {
	query := "INSERT INTO accessTokens (userId, token, isActive) VALUES (?, ?, ?)"

	_, err := transaction.Exec(query, userId, accessToken, 1)
	if err != nil {
		return errors.J(err, fmt.Sprintf("insert failed userId=%v accessToken=%v", userId, accessToken))
	}

	return nil
}

func (transaction *Transaction) FindUserIdForAccessToken(accessToken string) (*int64, error) {
	query := "SELECT userId FROM accessTokens WHERE token = ? AND isActive = 1"

	userIds, err := TxQuery[[]int64](transaction, query, []any{accessToken}, func(rows *sql.Rows) (*[]int64, error) {
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
