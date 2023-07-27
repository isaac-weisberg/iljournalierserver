package main

import "database/sql"

func (transaction *transaction) createKnownFlagsTable() error {
	query := `CREATE TABLE knownFlags (
		id INTEGER NOT NULL PRIMARY KEY,
		userId INTEGER NOT NULL,
		flagName TEXT NOT NULL,
		FOREIGN KEY (userId) REFERENCES users(id)
	)`

	_, err := transaction.exec(query)
	if err != nil {
		return j(err, "create table failed")
	}

	return nil
}

func (transaction *transaction) addKnownFlag(userId int64, text string) error {
	query := `INSERT INTO knownFlags (userId, flagName) VALUES (?, ?)`

	_, err := transaction.exec(query, userId, text)
	if err != nil {
		return j(err, "insert failed")
	}

	return nil
}

func (transaction *transaction) getKnownFlagIdsForUser(userId int64) ([]int64, error) {
	query := "SELECT (id) FROM knownFlags WHERE userId = ?"
	args := []any{userId}

	userIds, err := txQuery[[]int64](transaction, query, args, func(rows *sql.Rows) (*[]int64, error) {
		var flagIds []int64

		for rows.Next() {
			var flagId int64
			err := rows.Scan(&flagId)
			if err != nil {
				return nil, j(err, "scan failed")
			}
			flagIds = append(flagIds, flagId)
		}
		err := rows.Err()
		if err != nil {
			return nil, j(err, "rows returned error")
		}

		return &flagIds, nil
	})

	if err != nil {
		return nil, j(err, "txQuery failed")
	}

	return *userIds, nil
}
