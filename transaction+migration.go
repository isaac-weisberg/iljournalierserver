package main

import (
	"fmt"
	"strings"
)

var _migrationsTableName = "migrations"

func migrationsTableSql(sql string) string {
	if strings.Count(sql, "%s") != 1 {
		panic("migrationsTableSql Invalid SQL")
	}
	return fmt.Sprintf(sql, _migrationsTableName)
}

func (transaction transaction) createMigrationsTable() error {
	sql := migrationsTableSql("CREATE TABLE IF NOT EXISTS %s (version TEXT NOT NULL PRIMARY KEY)")
	_, err := transaction.exec(sql)
	if err != nil {
		return j(e("createMigrationsTable"), err)
	}
	return err
}

func (transaction transaction) hasVersionBeenMigrated(version string) (bool, error) {
	sql := migrationsTableSql("SELECT COUNT() FROM %s WHERE version == ?")

	row := transaction.queryRow(sql, version)

	var count int
	err := row.Scan(&count)

	if err != nil {
		return false, j(
			e("hasVersionBeenMigrated failed"),
			err,
		)
	}

	return count == 1, nil
}

func (transaction transaction) markVersionAsMigrated(version string) error {
	sql := migrationsTableSql("INSERT INTO %s VALUES (?)")
	_, err := transaction.exec(sql, version)
	if err != nil {
		return j(
			e("markVersionAsMigrated failed"),
			err,
		)
	}
	return nil
}
