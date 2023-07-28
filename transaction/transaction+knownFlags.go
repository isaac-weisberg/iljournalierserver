package transaction

import (
	"database/sql"
	"strings"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

func (transaction *Transaction) CreateKnownFlagsTable() error {
	query := `CREATE TABLE knownFlags (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		userId INTEGER NOT NULL,
		flagName TEXT NOT NULL,
		FOREIGN KEY (userId) REFERENCES users(id)
	)`

	_, err := transaction.Exec(query)
	if err != nil {
		return errors.J(err, "create table failed")
	}

	return nil
}

func (transaction *Transaction) AddKnownFlags(userId int64, flagNames []string) error {
	if len(flagNames) == 0 {
		return nil
	}
	var firstFlagName = flagNames[0]

	var args = make([]any, len(flagNames))
	var builder strings.Builder
	builder.WriteString("INSERT INTO knownFlags (userId, flagName) VALUES (?, ?)")
	args = append(args, userId, firstFlagName)

	var remainingFlagNames = flagNames[1:]

	for _, remainingFlagName := range remainingFlagNames {
		builder.WriteString(", (?, ?)")
		args = append(args, userId, remainingFlagName)
	}

	var query = builder.String()

	_, err := transaction.Exec(query, args...)
	if err != nil {
		return errors.J(err, "insert failed")
	}

	return nil
}

func (transaction *Transaction) GetKnownFlagIdsForUser(userId int64) ([]int64, error) {
	query := "SELECT (id) FROM knownFlags WHERE userId = ?"
	args := []any{userId}

	userIds, err := TxQuery[[]int64](transaction, query, args, func(rows *sql.Rows) (*[]int64, error) {
		var flagIds []int64

		for rows.Next() {
			var flagId int64
			err := rows.Scan(&flagId)
			if err != nil {
				return nil, errors.J(err, "scan failed")
			}
			flagIds = append(flagIds, flagId)
		}
		err := rows.Err()
		if err != nil {
			return nil, errors.J(err, "rows returned error")
		}

		return &flagIds, nil
	})

	if err != nil {
		return nil, errors.J(err, "txQuery failed")
	}

	return *userIds, nil
}
