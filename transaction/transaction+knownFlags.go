package transaction

import (
	"database/sql"
	"strings"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/models"
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

func (transaction *Transaction) AddKnownFlags(userId int64, flagNames []string) (*[]int64, error) {
	if len(flagNames) == 0 {
		return &[]int64{}, nil
	}
	var firstFlagName = flagNames[0]

	var args = make([]any, 0, len(flagNames)*2)
	var builder strings.Builder
	builder.WriteString("INSERT INTO knownFlags (userId, flagName) VALUES (?, ?)")
	args = append(args, userId, firstFlagName)

	var remainingFlagNames = flagNames[1:]

	for _, remainingFlagName := range remainingFlagNames {
		builder.WriteString(", (?, ?)")
		args = append(args, userId, remainingFlagName)
	}

	builder.WriteString(" RETURNING id")

	var query = builder.String()

	flagIds, err := TxQuery[[]int64](transaction, query, args, func(rows *sql.Rows) (*[]int64, error) {
		var flagIds = make([]int64, 0, len(flagNames))
		for rows.Next() {
			var flagId int64
			err := rows.Scan(&flagId)
			if err != nil {
				return nil, errors.J(err, "scaning flagId failed")
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
		return nil, errors.J(err, "insert failed")
	}

	return flagIds, nil
}

func (transaction *Transaction) GetKnownFlagsForUser(userId int64) (*[]models.FlagModel, error) {
	var query = "SELECT id, flagName FROM knownFlags WHERE userId = ?"
	var args = []any{userId}

	flagModels, err := TxQuery[[]models.FlagModel](transaction, query, args, func(rows *sql.Rows) (*[]models.FlagModel, error) {
		var flagModels []models.FlagModel
		for rows.Next() {
			var flagId int64
			var flagName string

			var err = rows.Scan(&flagId, &flagName)
			if err != nil {
				return nil, errors.J(err, "scan thrown error")
			}
			flagModels = append(flagModels, models.NewFlagModel(flagId, flagName))
		}
		var err = rows.Err()
		if err != nil {
			return nil, errors.J(err, "rows returned error")
		}

		return &flagModels, nil
	})

	if err != nil {
		return nil, errors.J(err, "tx query failed")
	}

	return flagModels, nil
}

func (transaction *Transaction) GetKnownFlagIdsForUser(userId int64) ([]int64, error) {
	query := "SELECT id FROM knownFlags WHERE userId = ?"
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
