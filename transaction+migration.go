package main

import "fmt"

func (transaction transaction) createMigrationsTable() error {
	sql := "CREATE TABLE IF NOT EXISTS migrations (version TEXT NOT NULL PRIMARY KEY)"
	_, err := transaction.exec(sql)
	if err != nil {
		return j(err, "create table failed")
	}
	return err
}

func (transaction transaction) hasVersionBeenMigrated(version string) (bool, error) {
	sql := "SELECT COUNT() FROM migrations WHERE version == ?"

	row := transaction.queryRow(sql, version)

	var count int
	err := row.Scan(&count)

	if err != nil {
		return false, j(err, "counting migrations failed")
	}

	return count == 1, nil
}

func (transaction transaction) markVersionAsMigrated(version string) error {
	sql := "INSERT INTO migrations VALUES (?)"
	_, err := transaction.exec(sql, version)
	if err != nil {
		return j(err, fmt.Sprintf("insert failed %s", version))
	}
	return nil
}
