package main

func (transaction transaction) createMigrationsTable() error {
	wrapError := createErrorWrapper("tx.createMigrationsTable")
	sql := "CREATE TABLE IF NOT EXISTS migrations (version TEXT NOT NULL PRIMARY KEY)"
	_, err := transaction.exec(sql)
	if err != nil {
		return wrapError(err)
	}
	return err
}

func (transaction transaction) hasVersionBeenMigrated(version string) (bool, error) {
	wrapError := createErrorWrapper("tx.hasVersionBeenMigrated failed")
	sql := "SELECT COUNT() FROM migrations WHERE version == ?"

	row := transaction.queryRow(sql, version)

	var count int
	err := row.Scan(&count)

	if err != nil {
		return false, wrapError(err)
	}

	return count == 1, nil
}

func (transaction transaction) markVersionAsMigrated(version string) error {
	wrapError := createErrorWrapper("tx.markVersionAsMigrated failed")
	sql := "INSERT INTO migrations VALUES (?)"
	_, err := transaction.exec(sql, version)
	if err != nil {
		return wrapError(err)
	}
	return nil
}
