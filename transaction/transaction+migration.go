package transaction

import (
	"fmt"

	"caroline-weisberg.fun/iljournalierserver/errors"
)

func (transaction *Transaction) CreateMigrationsTable() error {
	query := "CREATE TABLE IF NOT EXISTS migrations (version TEXT NOT NULL PRIMARY KEY)"
	_, err := transaction.Exec(query)
	if err != nil {
		return errors.J(err, "create table failed")
	}
	return err
}

func (transaction *Transaction) HasVersionBeenMigrated(version string) (bool, error) {
	query := "SELECT COUNT() FROM migrations WHERE version == ?"

	row := transaction.QueryRow(query, version)

	var count int
	err := row.Scan(&count)

	if err != nil {
		return false, errors.J(err, "counting migrations failed")
	}

	return count == 1, nil
}

func (transaction *Transaction) MarkVersionAsMigrated(version string) error {
	query := "INSERT INTO migrations VALUES (?)"
	_, err := transaction.Exec(query, version)
	if err != nil {
		return errors.J(err, fmt.Sprintf("insert failed %s", version))
	}
	return nil
}
